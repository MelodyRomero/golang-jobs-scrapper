package service

import (
	"fmt"
	model "golang-jobs-scrapper/internal/model"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type ScrappingService interface {
	GetJobsLinks() ([]*model.JobOffer, error)
}

type scrapper struct {
	Collector *colly.Collector
	Sites     []*model.Portal
}

func NewScrappingService(c *colly.Collector, sites []*model.Portal) ScrappingService {
	return &scrapper{
		Collector: c,
		Sites:     sites,
	}
}

func (s *scrapper) GetJobsLinks() ([]*model.JobOffer, error) {

	var jobsOffers []*model.JobOffer

	for _, site := range s.Sites {

		// On every a element which has href attribute call callback
		s.Collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
			var jobOffer model.JobOffer
			link := e.Attr("href")
			if link != "" && link != "/" {
				if valid := validateUrl(link, site.Keywords, site.Exclude); valid == true {
					s.Collector.Visit(e.Request.AbsoluteURL(link))
					jobOffer.URL = e.Request.AbsoluteURL(link)
					s.Collector.OnHTML(".title", func(e *colly.HTMLElement) {
						jobOffer.Title = e.Text
					})
					s.Collector.OnHTML("time[itemprop=validThrough]", func(e *colly.HTMLElement) {
						jobOffer.Date = e.Text
					})
					jobsOffers = append(jobsOffers, &jobOffer)
				}
			}

		})

		// Before making a request print "Visiting ..."
		s.Collector.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})

		// attach callbacks after login
		s.Collector.OnResponse(func(r *colly.Response) {
			log.Println("response received", string(r.StatusCode))
		})

		// Start scraping on https://hackerspaces.org
		s.Collector.Visit(site.BaseURL)
	}

	return jobsOffers, nil
}

func validateUrl(url string, keywords, exclude []string) bool {
	valid := true
	for _, keyword := range keywords {
		for _, excludeEntry := range exclude {
			if !strings.Contains(url, keyword) || strings.Contains(url, excludeEntry) {
				valid = false
			}
		}
	}
	return valid
}
