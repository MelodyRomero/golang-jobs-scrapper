package main

import (
	"fmt"
	"golang-jobs-scrapper/internal/model"
	service "golang-jobs-scrapper/internal/services"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

var URLs []string

func main() {

	for _, portal := range Webs {
		URLs = append(URLs, portal.BaseURL)
	}
	fmt.Println(URLs)
	// Instantiate default collector
	c := colly.NewCollector(
		// Attach a debugger to the collector
		colly.Debugger(&debug.LogDebugger{}),
	// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
	//	colly.AllowedDomains("welovegolang.com"),
	// MaxDepth is 1, so only the links on the scraped page
	// is visited, and no further links are followed
	//	colly.MaxDepth(1),

	// Visit only root url and urls which start with "e" or "h" on httpbin.org
	//	colly.URLFilters(
	//		regexp.MustCompile("http://httpbin\\.org/(|e.+)$"),
	//		regexp.MustCompile("http://httpbin\\.org/h.+"),
	//	),
	)

	service := service.NewScrappingService(c, Webs)

	jobs, err := service.GetJobsLinks()

	if err != nil {
		fmt.Println(err)
	}
	for _, job := range jobs {
		fmt.Printf(
			"\n*****Job Offer*****\n %+v\n",
			job,
		)
	}
}

var Webs = []*model.Portal{
	{
		Name:     "We Love Golang",
		BaseURL:  "https://www.welovegolang.com",
		JobsURL:  "https://www.welovegolang.com",
		Keywords: []string{"/job"},
		Exclude:  []string{"/job-qa"},
	},
	{
		Name:     "Golang Projects .com",
		BaseURL:  "https://www.golangprojects.com",
		JobsURL:  "https://www.golangprojects.com/golang-remote-jobs.html",
		Keywords: []string{"/golang-go-job"},
		Exclude:  []string{""},
	},
}
