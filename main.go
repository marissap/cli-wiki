package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var Topics string

const (
	baseURL       = "https://en.wikipedia.org/w/api.php?action=query&format=json&prop=extracts&titles="
	wordSeperator = "%20"
	endingURL     = "&exintro=1"
)

type Text struct {
	PageID  int    `json:"pageid"`
	Ns      int    `json:"ns"`
	Title   string `json:"title"`
	Summary string `json:"extract"`
}

type Response struct {
	Status   string `json:"batchcomplete"`
	Warnings struct {
		Extracts struct {
			Special string `json:"*"`
		} `json:"extracts"`
	} `json:"warnings"`
	Query struct {
		Normalized []map[string]string `json:"normalized`
		Pages      map[string]Text
	} `json:"query"`
}

func init() {
	const (
		usage = "Search Wikipedia Topic (Ex. \"dark matter\")"
	)
	flag.StringVar(&Topics, "t", "", usage)
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Searching for topic: %s\n", Topics)

	SearchTopicOnWikipedia(Topics)
}

func SearchTopicOnWikipedia(query string) {
	var responseBody Response
	rawText := ""
	p := bluemonday.StripTagsPolicy()

	queryString := strings.ReplaceAll(query, " ", wordSeperator)

	resp, err := http.Get(baseURL + queryString + endingURL)
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		os.Exit(1)
	}

	if resp.StatusCode == 200 {
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed reading data from the response %s\n", err)
			os.Exit(1)
		}

		json.Unmarshal(responseData, &responseBody)

		result := responseBody.Query.Pages
		for _, item := range result {
			rawText = item.Summary
		}

		html := p.Sanitize(rawText)
		fmt.Println(html)

	} else {
		fmt.Printf("%s not found\n", query)
	}

}
