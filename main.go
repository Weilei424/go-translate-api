package main

import (
	"flag"
	"fmt"
	"github.com/Weilei424/go-translate-api/cli"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var sourceLang string
var targetLang string
var sourceText string

func init() {
	flag.StringVar(&sourceLang, "s", "en", "Source Language[en]")
	flag.StringVar(&targetLang, "t", "fr", "Target Language[fr]")
	flag.StringVar(&sourceText, "st", "", "Text to translate")
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Options: ")
		flag.PrintDefaults()
		os.Exit(1)
	}

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
