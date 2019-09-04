package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//IndeedCrawler struct implements Crawler interface
type IndeedCrawler struct {
	Data []map[string]string
}

//Scrape function for indeed.com
func (ic *IndeedCrawler) Scrape(pr *PrinterRequest) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)

	defer cancel()

	fmt.Println(pr.JobTitle)

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://de.indeed.com"),
		chromedp.WaitVisible(`#text-input-what`, chromedp.ByID),
		chromedp.SetValue(`#text-input-what`, pr.JobTitle, chromedp.ByID),
		chromedp.SetValue(`#text-input-where`, pr.JobLocation, chromedp.ByID),
		chromedp.Submit(`icl-Button icl-Button--primary icl-Button--md icl-WhatWhere-button`),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println("test")
	foundNum := 0
	var panels []*cdp.Node
	for j := 2; j < 4; j++ {
		fmt.Printf("Currently Fetching Page: %d\n\n", j-1)
		sel := `//div[contains(concat(' ',normalize-space(@class),' '),' jobsearch-SerpJobCard ')]`
		if err := chromedp.Run(ctx,
			chromedp.EmulateViewport(1920, 2000),
			chromedp.WaitReady(sel),
			chromedp.Nodes(sel, &panels),
		); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d Nodes Found.\n", len(panels))
		sel = fmt.Sprintf(`//div[@class="pagination"]/a/span[text()="%d"]`, j)

		for index, node := range panels {
			var html string
			var company string
			if err := chromedp.Run(ctx,
				chromedp.EmulateViewport(1920, 2000),
				chromedp.MouseClickNode(node),
				chromedp.WaitVisible(`//div[@id="vjs-desc"]`),
				chromedp.InnerHTML(`//div[@id="vjs-desc"]`, &html),
				chromedp.Text(`//span[@id="vjs-cn"]`, &company),
			); err != nil {
				log.Printf("Node %d failed.", index)
			}

			ic.ProcessData(company, html, &foundNum)
			fmt.Println(company)
			foundNum = foundNum + 1
		}
		chromedp.Run(ctx,
			chromedp.WaitVisible(sel),
			chromedp.Click(sel, chromedp.NodeVisible),
			chromedp.WaitNotPresent("body", chromedp.BySearch),
		)
	}

	WriteJSON(ic.Data)
}

//ProcessData adds data to Crawler
func (ic *IndeedCrawler) ProcessData(key interface{}, value interface{}, index *int) {
	if key != "" || value != "" {
		if len(ic.Data) == 0 {
			ic.Data = make([]map[string]string, 1)
		}
		if cap(ic.Data) == len(ic.Data) {
			n := make([]map[string]string, cap(ic.Data)+1)
			copy(n, ic.Data)
			ic.Data = n
		}

		ic.Data[*index] = make(map[string]string)

		ic.Data[*index][key.(string)] = value.(string)
		return
	}

	*index = *index - 1

}

//DispatchCrawlers dispatches concurrent Crawlers
func DispatchCrawlers(pr *PrinterRequest) {
	ic := IndeedCrawler{}
	fmt.Println(pr.JobTitle)

	ic.Scrape(pr)
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
