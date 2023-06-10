package crawler

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Avash027/Intra-Page-Ranker/logger"
	"github.com/gocolly/colly"
)

type Link struct {
	From string
	To   string
}

type WebCrawler struct {
	DomainName   string
	MaxDepth     int
	MaxPages     int
	Spider       colly.Collector
	CrawlDelayIn time.Duration
	Links        []Link
	Visited      map[string]bool
	Count        int
}

type WebCrawlerOpts struct {
	DomainName     string
	MaxDepth       int
	MaxPages       int
	CrawlDelay     time.Duration
	UserAgent      string
	MaxParallelism int
}

func CreateWebCrawler(webCrawlerOpts WebCrawlerOpts) *WebCrawler {
	var wc WebCrawler
	wc.DomainName = webCrawlerOpts.DomainName
	wc.MaxDepth = webCrawlerOpts.MaxDepth
	wc.MaxPages = webCrawlerOpts.MaxPages
	wc.CrawlDelayIn = webCrawlerOpts.CrawlDelay
	wc.Visited = make(map[string]bool)
	wc.Count = 0

	wc.Spider = *colly.NewCollector(
		colly.AllowedDomains(wc.DomainName),
		colly.MaxDepth(wc.MaxDepth),
		colly.Async(true),
		colly.UserAgent(webCrawlerOpts.UserAgent),
		colly.URLFilters(
			regexp.MustCompile(fmt.Sprintf("^https?://%s/", wc.DomainName)),
		),
	)

	wc.Spider.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: webCrawlerOpts.MaxParallelism,
		Delay:       wc.CrawlDelayIn,
	})

	return &wc
}

func (wc *WebCrawler) Crawl() {

	wc.Spider.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))

		if strings.HasPrefix(link, "#") {
			return
		}

		link = removeParamsFromUrl(link)

		if wc.Visited[link] {
			e.Request.Abort()
			return
		}

		if !isValidDomain(link, wc.DomainName) {
			e.Request.Abort()
			return
		}

		wc.Visited[link] = true
		wc.Count++

		if wc.Count > wc.MaxPages {
			e.Request.Abort()
			return
		}

		if wc.Spider.MaxDepth > 0 && e.Request.Depth >= wc.Spider.MaxDepth {
			e.Request.Abort()
			return
		}

		if wc.Count%10 == 0 {
			go logger.InPlaceCyan(fmt.Sprintf("\r [>] %d Pages crawled", wc.Count))
		}

		linkEdge := Link{
			From: e.Request.URL.String(),
			To:   link,
		}

		wc.Links = append(wc.Links, linkEdge)
		wc.Spider.Visit(link)
	})

	wc.Spider.OnRequest(func(r *colly.Request) {

	})

	wc.Spider.OnError(func(r *colly.Response, err error) {
		logger.Red(fmt.Sprintf("Error: %v", err))
	})

	wc.Visited[fmt.Sprintf("https://%s/", wc.DomainName)] = true
	wc.Count++
	wc.Spider.Visit(fmt.Sprintf("https://%s/", wc.DomainName))
	wc.Spider.Wait()
	logger.InPlaceCyan("\n")
}

func removeParamsFromUrl(url string) string {
	res := strings.Split(url, "?")[0]

	if len(res) == 0 {
		return res
	}

	if res[len(res)-1] != '/' {
		res = res + "/"
	}

	return res
}

func isValidDomain(AbsoluteURL, domainName string) bool {
	re := regexp.MustCompile(`^https?://(?:www\.)?` + domainName)
	return re.MatchString(AbsoluteURL)
}
