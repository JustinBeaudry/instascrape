package services

import (
	"github.com/gocolly/colly"
	"instascrape/models"
	"encoding/json"
	"crypto/md5"
	"strings"
	"fmt"
	"os"
	"log"
)

const instagramUrl = "https://instagram.com/"

func Scrape(userSlug string, force bool, cacheDir string) (*models.InstagramPageData, error) {

	var results *models.InstagramPageData
	var collectorError error
	var requestURL = instagramUrl + userSlug
	var collector *colly.Collector
	var uaString = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"

	if force {
		removeErr := os.RemoveAll(cacheDir)
		if removeErr != nil {
			log.Fatal(removeErr)
		}
	}

	// for now to force cache reset we just clear the cache entirely
	collector = colly.NewCollector(
		colly.CacheDir(cacheDir),
		colly.UserAgent(uaString),
	)

	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("Referrer", requestURL)
		if r.Ctx.Get("gis") != "" {
			gis := fmt.Sprintf("%s:%s", r.Ctx.Get("gis"), r.Ctx.Get("variables"))
			h := md5.New()
			h.Write([]byte(gis))
			gisHash := fmt.Sprintf("%x", h.Sum(nil))
			r.Headers.Set("X-Instagram-GIS", gisHash)
		}
	})

	collector.OnHTML("body > script:first-of-type", func(e *colly.HTMLElement) {

		jsonData := e.Text[strings.Index(e.Text, "{") : len(e.Text)-1]
		data := &models.InstagramPageData{}

		err := json.Unmarshal([]byte(jsonData), data)
		if err != nil {
			collectorError = err
		}

		results = data
	})

	collector.OnError(func(_ *colly.Response, err error) {
		collectorError = err
	})

	collector.Visit(requestURL)

	collector.Wait()

	return results, collectorError
}
