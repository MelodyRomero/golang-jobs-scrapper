package main

import (
	"encoding/json"
	"fmt"
	"golang-jobs-scrapper/internal/model"
	"net/http"

	"github.com/BoseCorp/pester"
	"github.com/PuerkitoBio/goquery"
)

// Custom user agent.
const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/53.0.2785.143 " +
		"Safari/537.36"
)

func MakeRequest(url string) (string, int, error) {
	// Open url.
	// Need to use http.Client in order to set a custom user agent:
	client := pester.New()
	// set retry policy to 5 max retries
	client.MaxRetries = 5
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", http.StatusBadRequest, err
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {

		return "", http.StatusBadGateway, err
	}
	// Close the request
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	// Create array for jobs
	var jobs []*model.JobOffer

	// Find all songs on page and parse string into artist and song
	doc.Find(".media").Each(func(_ int, s *goquery.Selection) {
		var job model.JobOffer
		job.Description = s.Find(".summary").Text()
		job.Date = s.Find("time").First().Text()
		job.Location = s.Find(".location").Find("span").First().Text()
		job.Organization = s.Find(".company").Find("span").Text()
		//fmt.Println("*************")
		//fmt.Println(s.Children().Html())
		//fmt.Println("*************")
		s.Find(".media-body").Each(func(_ int, t *goquery.Selection) {
			t.Find(".media-heading").Each(func(index int, j *goquery.Selection) {

				jobDetailUrl, exist := j.Find("a").Attr("href")
				if exist != true {
					fmt.Print(err.Error())
				}
				job.URL = url + jobDetailUrl
				job.Title = j.Find("span").Text()
				jobs = append(jobs, &job)
			})
		})

	})

	raw, _ := json.Marshal(jobs)

	return string(raw), http.StatusOK, nil
}

// fetchUrl opens a url with GET method and sets a custom user agent.
// If url cannot be opened, then log it to a dedicated channel.
func fetchUrl(url string, chFailedUrls chan string, chIsFinished chan bool) {

	response, statuscode, err := MakeRequest(url)
	// Inform the channel chIsFinished that url fetching is done (no
	// matter whether successful or not). Defer triggers only once
	// we leave fetchUrl():
	defer func() {
		chIsFinished <- true
	}()

	// If url could not be opened, we inform the channel chFailedUrls:
	if err != nil || statuscode != 200 {
		chFailedUrls <- url
		return
	}

	fmt.Printf("Response %+v\n", response)
}
func main() {

	urlsList := []string{
		//	"https://www.golangprojects.com/golang-remote-jobs.html", //https://www.golangprojects.com/golang-remote-jobs.html
		"https://www.welovegolang.com", //https://www.welovegolang.com/jobs/remote
		/*	"https://remoteok.io/remote-golang-jobs",
			"https://stackoverflow.com/jobs/remote-developer-jobs-using-go",
			"https://www.indeed.com/q-Golang-l-Remote-jobs.html",
			"https://golang.cafe/Remote-Golang-Developer-Jobs",
			"https://nodesk.co/remote-jobs/golang/",
			"https://www.ziprecruiter.com/Jobs/Remote-Golang",
			"https://www.workingnomads.co/remote-golang-jobs",
			"https://www.glassdoor.com/Job/golang-developer-remote-jobs-SRCH_KO0,23.htm",
			"https://remote4me.com/remote-go-golang-jobs",
			"https://x-team.com/remote-go-developer-jobs", */
	}

	// Create 2 channels, 1 to track urls we could not open
	// and 1 to inform url fetching is done:
	chFailedUrls := make(chan string)
	chIsFinished := make(chan bool)

	// Open all urls concurrently using the 'go' keyword:
	for _, url := range urlsList {
		go fetchUrl(url, chFailedUrls, chIsFinished)
	}

	// Receive messages from every concurrent goroutine. If
	// an url fails, we log it to failedUrls array:
	failedUrls := make([]string, 0)
	for i := 0; i < len(urlsList); {
		select {
		case url := <-chFailedUrls:
			failedUrls = append(failedUrls, url)
		case <-chIsFinished:
			i++
		}
	}
	// Print all urls we could not open:
	if len(failedUrls) > 0 {
		fmt.Println("Could not fetch these urls: ", failedUrls)
	}

}
