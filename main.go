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
	defer response.Body.Close()

	checkErr(err)

	links, err := extractlinks.All(response.Body)
	checkErr(err)

	fmt.Println(links)

	response.Body.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}