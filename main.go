package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/Weilei424/go-translate-api/cli"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var sourceLang string
var targetLang string
var sourceText string
var trends bool

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel *Channel `xml:"channel"`
}

type Channel struct {
	Title    string `xml:"title"`
	ItemList []Item `xml:"item"`
}

type Item struct {
	Title     string `xml:"title"`
	Link      string `xml:"link"`
	Traffic   string `xml:"approx_traffic"`
	NewsItems []News `xml:"news_item"`
}

type News struct {
	Headline     string `xml:"news_item_title"`
	HeadlineLink string `xml:"news_item_url"`
}

func init() {
	flag.StringVar(&sourceLang, "s", "en", "Source Language[en]")
	flag.StringVar(&targetLang, "t", "fr", "Target Language[fr]")
	flag.StringVar(&sourceText, "st", "", "Text to translate")
	flag.BoolVar(&trends, "tr", false, "Show Google Trends")
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Options: ")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if trends {
		var rss RSS
		data := readGoogleTrends()
		err := xml.Unmarshal(data, &rss)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("\n Below are all the Google Search Trends in Canada For Today! ")
		fmt.Println("---------------------------------------------------------------")

		for i := range rss.Channel.ItemList {
			rank := i + 1
			fmt.Println("#", rank)
			fmt.Println("Search Term: ", rss.Channel.ItemList[i].Title)
			fmt.Println("Link to the Trend: ", rss.Channel.ItemList[i].Link)
			fmt.Println("Headline: ", rss.Channel.ItemList[i].NewsItems[0].Headline)
			fmt.Println("---------------------------------------------------------------")
		}
	} else {
		strChan := make(chan string)
		wg.Add(1)

		requestBody := &cli.RequestBody{
			SourceLang: sourceLang,
			TargetLang: targetLang,
			SourceText: sourceText,
		}

		go cli.RequestTranslate(requestBody, strChan, &wg)

		processedStr := strings.ReplaceAll(<-strChan, " + ", " ")

		fmt.Printf("%s\n", processedStr)

		close(strChan)
		wg.Wait()
	}
}

func readGoogleTrends() []byte {
	response := getGoogleTrends()
	data, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return data
}

// get http request
func getGoogleTrends() *http.Response {
	response, err := http.Get("https://trends.google.com/trends/trendingsearches/daily/rss?geo=CA")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return response
}
