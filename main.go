package main

import (
	"flag"

	"github.com/Avash027/Intra-Page-Ranker/constants"
	"github.com/Avash027/Intra-Page-Ranker/crawler"
	"github.com/Avash027/Intra-Page-Ranker/exporter"
	"github.com/Avash027/Intra-Page-Ranker/matrix"
	"github.com/Avash027/Intra-Page-Ranker/ranker"
)

func main() {

	// Define command line flags
	domainName := flag.String("domain", constants.DOMAIN, "The domain name to crawl")
	maxDepth := flag.Int("depth", constants.DEFAULT_MAX_DEPTH, "The maximum depth to crawl")
	maxPages := flag.Int("pages", constants.DEFAULT_MAX_PAGES, "The maximum number of pages to crawl")
	crawlDelay := flag.Duration("delay", constants.DEFAULT_CRAWL_DELAY, "The delay between requests")
	userAgent := flag.String("useragent", constants.USER_AGENT, "The user agent to use for requests")
	maxParallelism := flag.Int("parallelism", constants.MAX_PARALLELISM, "The maximum number of parallel requests")
	alpha := flag.Float64("alpha", constants.ALPHA, "Constant used while calculating page rank")
	epsilon := flag.Float64("epsilon", constants.EPSILON, "Directly proportional to accuracy of results but also consumes more time")
	numOfResult := flag.Int("maxRank", constants.NUM_OF_RESULT, "Max Rank upto which you want to see the result")

	// Parse command line flags
	flag.Parse()

	r := ranker.CreateRanker(
		ranker.RankerOpts{
			Crawler: &crawler.WebCrawlerOpts{
				DomainName:     *domainName,
				MaxDepth:       *maxDepth,
				MaxPages:       *maxPages,
				CrawlDelay:     *crawlDelay,
				UserAgent:      *userAgent,
				MaxParallelism: *maxParallelism,
			},
			Markov: &matrix.MarkovMatrixOpts{
				Alpha:   *alpha,
				Epsilon: *epsilon,
			},
			Exporter: &exporter.RexOpts{
				Type:     "csv",
				Dir:      "./",
				FileName: "filename.csv",
			},
			NumOfResults: *numOfResult,
		},
	)

	r.GetPageRank()
}
