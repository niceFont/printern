package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strings"

	"github.com/sclevine/agouti"

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
	fmt.Println(pr.JobLocation)
	var err error

	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	locate := agouti.ChromeOptions("args", []string{"--disable-geolocation"})

	agoutiDriver := agouti.ChromeDriver(locate)

	err = agoutiDriver.Start()

	page, err := agoutiDriver.NewPage()
	url := fmt.Sprintf("https://indeed.com/jobs?q=%s&l=%s", pr.JobTitle, pr.JobLocation)
	linkhref := fmt.Sprintf(`//p[@class="oocs"]/a[contains(text(), "%s")]`, pr.JobTitle)
	err = page.Navigate(url)
	err = page.FindByXPath(linkhref).Click()

	var sel string
	var html string
	for i := 1; i <= 2; i++ {
		sel = fmt.Sprintf("//div[contains(concat(' ',normalize-space(@class),' '),' jobsearch-SerpJobCard ')][%d]", i)
		err = page.FindByXPath(sel).Click()
		err = page.NextWindow()
		html, err = page.HTML()
		//	err = page.NextWindow()
		err = page.CloseWindow()
		ic.ProcessData(html)
	}
	//err = agoutiDriver.Stop()
}

//ProcessData adds data to Crawler
func (ic *IndeedCrawler) ProcessData(input string) {

	data := strings.Split(input, " ")
	fmt.Printf("%+v", input)
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
