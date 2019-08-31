package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
)

//IndeedCrawler struct implements Crawler interface
type IndeedCrawler struct {
}

//Scrape function for indeed.com
func (is *IndeedCrawler) Scrape(pr *PrinterRequest) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancel()

	var result string

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://de.indeed.com"),
		chromedp.WaitVisible(`#text-input-what`, chromedp.ByID),
		chromedp.SendKeys(`#text-input-what`, pr.JobTitle),
		chromedp.SendKeys(`#text-input-where`, pr.JobLocation),
		chromedp.Click(`icl-Button icl-Button--primary icl-Button--md icl-WhatWhere-button`, chromedp.NodeVisible),
		chromedp.Text(`#searchCount`, &result, chromedp.NodeVisible, chromedp.ByID),
	)

	fmt.Println("test")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

func DispatchCrawlers(pr *PrinterRequest) {
	ic := IndeedCrawler{}

	go ic.Scrape(pr)
}
