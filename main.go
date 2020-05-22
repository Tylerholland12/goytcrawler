package main

import (
	"crypto/tls"
	"fmt"
	// "io/ioutil"
	"net/http"
	"os"

	"github.com/steelx/extractlinks"
)

var (
	config = &tls.Config{
		InsecureSkipVerify: true,
	}

	transport = &http.Transport{
		TLSClientConfig: config,
	}
	netClient = &http.Client{
		Transport: transport,
	}
	queue = make(chan string)
	hasVisited = make(map[string]bool)
)

func main() {
	arguements := os.Args[1:]
	
	if len(arguements) == 0 {
		fmt.Println("Missing URL, e.g. go-webscrapper http://js.org/")
		os.Exit(1)
	}
	
	baseURL := arguements[0]

	// keep function concurrent to not exhaust all resources visiting a single link
	go func () {
		queue <- baseURL
	}()

	// crawl url when it is recieved  
	for href := range queue && isSameDomain(href, baseURL){
		if !hasVisited[href] {
			crawlURL(href)
		}
	}	
}

func isSameDomain () {
	uri, err := url.Parse(href) 
	if err != nil {
		return false
	}
	parentURI, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	if uri.Host != parentURI.Host {
		return false
	}

	return true
}

func crawlURL(href string) {
	hasVisited[href] = true
	fmt.Printf("Crawling url -> %v \n", href)
	response, err := netClient.Get(baseURL)
	checkErr(err)
	defer response.Body.Close()


	links, err := extractlinks.All(response.Body)
	checkErr(err)

	for _, link := range links {
		absluteURL := toFixedURL(link.Href, href)
		// so that we are not pushing links faster than we receive them 
		go func () {
			queue <- absluteURL
	}()
}

	fmt.Println(links)

	response.Body.Close()

}

func toFixedURL(href, baseURL string) {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}


	base, err := url.Parse(baseURL) 
	if err != nil {
		return ""
	}

	// host from base
	// path from uri 
	// has its own host
	toFixedURI := base.ResolveReference(uri)
	return toFixedURI.String()


}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}