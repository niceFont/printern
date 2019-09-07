package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

//IndeedCrawler struct implements Crawler interface
type IndeedCrawler struct {
	ID   uuid.UUID
	Data map[string]int
}

//Scrape function for indeed.com
func (ic *IndeedCrawler) Scrape(pr *PrinterRequest) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)

	defer cancel()

	fmt.Println(pr.JobTitle)

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://de.indeed.com"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("Navigate Done")
			return nil
		}),
		chromedp.WaitVisible(`#text-input-what`, chromedp.ByID),
		chromedp.SetValue(`#text-input-what`, pr.JobTitle, chromedp.ByID),
		chromedp.SetValue(`#text-input-where`, pr.JobLocation, chromedp.ByID),
		chromedp.Submit(`icl-Button icl-Button--primary icl-Button--md icl-WhatWhere-button`),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println("test")
	var panels []*cdp.Node
	var html string
	var company string
	for j := 2; j < 42; j++ {
		fmt.Printf("Currently Fetching Page: %d\n\n", j-1)
		sel := `//div[contains(concat(' ',normalize-space(@class),' '),' jobsearch-SerpJobCard ')]`
		if err := chromedp.Run(ctx,
			chromedp.EmulateViewport(1920, 2000),
			chromedp.ActionFunc(func(ctx context.Context) error {
				log.Println("Waiting Nodes")
				return nil
			}),
			chromedp.WaitVisible(`//div[contains(concat(' ',normalize-space(@class),' '),' jobsearch-SerpJobCard ')]`),
			chromedp.ActionFunc(func(ctx context.Context) error {
				log.Println("Waiting Nodes Done")
				return nil
			}),
			chromedp.Nodes(sel, &panels, chromedp.NodeVisible),
		); err != nil {
			log.Println(err)

		}
		fmt.Printf("%d Nodes Found.\n", len(panels))
		sel = fmt.Sprintf(`//div[@class="pagination"]/a/span[text()="%d"]`, j)

		for index, node := range panels {
			if err := ctx.Err(); err != nil {
				return
			}
			if err := chromedp.Run(ctx,
				chromedp.EmulateViewport(1920, 2000),
				chromedp.ActionFunc(func(ctx context.Context) error {
					log.Println("Clicking Nodes")
					return nil
				}),
				chromedp.MouseClickNode(node),
				chromedp.ActionFunc(func(ctx context.Context) error {
					log.Println("Clicking Nodes Done")
					return nil
				}),
				chromedp.ActionFunc(func(ctx context.Context) error {
					log.Println("Waiting Desc")
					return nil
				}),
				chromedp.WaitReady(`//div[@id="vjs-desc"]`),
				chromedp.ActionFunc(func(ctx context.Context) error {
					log.Println("Waiting Desc Done")
					return nil
				}),
				chromedp.InnerHTML(`//div[@id="vjs-desc"]`, &html),
				chromedp.Text(`//span[@id="vjs-cn"]`, &company),
			); err != nil {
				log.Printf("Node %d failed.", index)
				continue
			}

			ic.ProcessData(html)
			fmt.Println(company)
		}
		if err := chromedp.Run(ctx,
			chromedp.WaitReady(sel),
			chromedp.Click(sel, chromedp.NodeVisible),
			chromedp.WaitNotPresent(sel),
		); err != nil {
			log.Println(err)
		}
	}

}

//ProcessData adds data to Crawler
func (ic *IndeedCrawler) ProcessData(input string) {

	data := strings.Split(input, " ")

	for _, item := range data {
		if _, ok := ic.Data[strings.ToLower(item)]; ok {
			ic.Data[strings.ToLower(item)]++

		}
	}

}

func (ic *IndeedCrawler) GetData() ScrapeResult {
	return ic.Data
}

//DispatchCrawlers dispatches concurrent Crawlers
func DispatchCrawlers(pr *PrinterRequest) map[string]int {
	ic, err := NewCrawler(0)

	if err != nil {
		log.Fatal(err)
	}

	ic.Scrape(pr)

	return ic.GetData()

}

func crawlerCapture(ctx *context.Context) {
	var buff []byte
	if err := chromedp.Run(*ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)

			if err != nil {
				return err
			}

			buff, err = page.CaptureScreenshot().
				WithQuality(90).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)

			if err != nil {
				return err
			}

			return nil
		}),
	); err != nil {
		log.Fatalln(err)
	}
	if err := ioutil.WriteFile("sc.png", buff, 0644); err != nil {
		log.Fatal(err)
	}
}

//NewCrawler creates a new Crawler based on the input given
func NewCrawler(ctype int) (Crawler, error) {

	var err error

	switch ctype {
	case 0:
		list := make(map[string]int)

		for _, keyword := range Keywords {
			list[keyword] = 0
		}
		ic := &IndeedCrawler{
			ID:   uuid.New(),
			Data: list,
		}

		return ic, nil
	default:
		return nil, err
	}
}
