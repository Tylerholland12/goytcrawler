package main

import (
	"crypto/tls"
	"fmt"
	// "io/ioutil"
	"net/http"
	"os"

	"github.com/steelx/extractlinks"
)

func main() {
	baseURL := "https://www.youtube.com/"

	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: config,
	}
	netClient := &http.Client{
		Transport: transport,
	}

	response, err := netClient.Get(baseURL)
	checkErr(err)

	defer response.Body.Close()

	

	links, err := extractlinks.All(response.Body)
	checkErr(err)

	for i, link := range links {
		fmt.Printf("index %v -- link %v \n", i, link.Href)
	}

	fmt.Println(links)

	response.Body.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}