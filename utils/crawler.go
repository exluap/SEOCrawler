/**
    * @project SEOCrawler
    * @date 02.04.2018 19:29
    * @author Nikita Zaytsev (exluap) <nickzaytsew@gmail.com>
    * @twitter https://twitter.com/exluap
    * @keybase https://keybase.io/exluap
*/

package utils

import (
"crypto/tls"
"fmt"
"github.com/jackdanger/collectlinks"
"net/http"
"net/url"
"github.com/badoux/goscraper"
)


func StartHere(url string) {


	queue := make(chan string)
	filteredQueue := make(chan string)

	go func() { queue <- url }()
	go filterQueue(queue, filteredQueue)


	done := make(chan bool)


	for i := 0; i < 5; i++ {
		go func() {
			for uri := range filteredQueue {
				enqueue(uri, queue)
			}
			done <- true
		}()
	}
	<-done
}

func filterQueue(in chan string, out chan string) {
	var seen = make(map[string]bool)
	for val := range in {
		if !seen[val] {
			seen[val] = true
			out <- val
		}
	}
}

func enqueue(uri string, queue chan string) {
	checkDecription(uri)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: transport}
	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)

	for _, link := range links {
		absolute := fixUrl(link, uri)
		if uri != "" {
			go func() { queue <- absolute }()
		}
	}
}

func fixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}

func checkDecription(url string) {
	s, err := goscraper.Scrape(url,5)

	if err != nil {
		fmt.Println(err)
		return
	}

	if s.Preview.Description == "" {
		fmt.Printf("Title: %s\n", s.Preview.Title)
		fmt.Printf("Url: %s\n",s.Preview.Link)
	}

}
