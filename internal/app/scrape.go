package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
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

	var testing string
	var panels []*cdp.Node
	for j := 2; j < 10; j++ {
		fmt.Printf("Currently Fetching Page: %d\n\n", j-1)
		sel := `//div[contains(concat(' ',normalize-space(@class),' '),' jobsearch-SerpJobCard ')]`
		if err := chromedp.Run(ctx,
			chromedp.EmulateViewport(1920, 2000),
			chromedp.WaitVisible(sel),
			chromedp.Nodes(sel, &panels),
		); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d Nodes Found.\n", len(panels))
		sel = fmt.Sprintf(`//div[@class="pagination"]/a/span[text()="%d"]`, j)

		for index, node := range panels {

			if err := chromedp.Run(ctx,
				chromedp.EmulateViewport(1920, 2000),
				chromedp.MouseClickNode(node),
				chromedp.WaitVisible("#vjs-jobtitle", chromedp.ByID),
				chromedp.Text("#vjs-jobtitle", &testing, chromedp.ByID),
			); err != nil {
				log.Printf("Node %d failed.", index)
			}

			fmt.Printf("Result: %s\n", testing)
		}
		chromedp.Run(ctx,
			chromedp.WaitVisible("#vjs-desc", chromedp.ByID),
			chromedp.Click(sel, chromedp.NodeVisible),
			chromedp.WaitNotPresent("body", chromedp.BySearch),
		)
		crawlerCapture(&ctx)
	}

}

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
