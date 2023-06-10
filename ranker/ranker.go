package ranker

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Avash027/Intra-Page-Ranker/crawler"
	"github.com/Avash027/Intra-Page-Ranker/exporter"
	"github.com/Avash027/Intra-Page-Ranker/logger"
	"github.com/Avash027/Intra-Page-Ranker/matrix"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Ranker struct {
	Crawler      *crawler.WebCrawler
	Markov       *matrix.MarkovMatrix
	Exporter     exporter.RexOpts
	Result       map[string]float64
	NumOfResults int
}

type RankerOpts struct {
	Crawler      *crawler.WebCrawlerOpts
	Markov       *matrix.MarkovMatrixOpts
	Exporter     *exporter.RexOpts
	NumOfResults int
}

func CreateRanker(opts RankerOpts) *Ranker {

	if opts.NumOfResults == 0 {
		opts.NumOfResults = math.MaxInt32
	}

	return &Ranker{
		Crawler:      crawler.CreateWebCrawler(*opts.Crawler),
		Markov:       matrix.CreateMarkovMatrix(*opts.Markov),
		Exporter:     *opts.Exporter,
		NumOfResults: opts.NumOfResults,
	}
}

func (r *Ranker) GetPageRank() {
	logger.Blue("[-] Crawling has started")

	r.Crawler.Crawl()

	logger.Green("[+] Crawing is completed")
	logger.Purple(fmt.Sprintf("[+] A total of %d links were recorded", len(r.Crawler.Links)))

	logger.Blue("[-] Sanitizing Links")

	r.Markov.InitalizeLinks(r.Crawler.Links)

	logger.Green("[+] Completed santizing links. Creating Transition Matrix matrix")

	r.Markov.CreateTransitionMatrix()

	logger.Blue("[+] Transition Matrix Created. Computing Page rank!")

	result := r.Markov.ComputePageRank()

	logger.Green("[+] Page rank computed. Printing results")

	r.printResults(result)

}

func (r *Ranker) printResults(result []matrix.RankData) {
	sort.Slice(result, func(i, j int) bool {
		if result[i].Weight == result[j].Weight {
			return len(result[i].Url) < len(result[j].Url)
		}
		return result[i].Weight > result[j].Weight
	})

	count := 0

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleBold)
	t.AppendHeader(table.Row{"Rank", "Url", "Weight"})

	for i := 0; i < len(result); i++ {
		// result[i].Url = removeDomain(result[i].Url)
		if i == 0 {
			result[0].Rank = 1
		} else if result[i].Weight == result[i-1].Weight {
			result[i].Rank = result[i-1].Rank
		} else {
			result[i].Rank = result[i-1].Rank + 1
		}
	}

	for i := 0; i < len(result) && count < r.NumOfResults; i++ {
		result[i].Url = removeDomain(result[i].Url)
		if i == 0 {
			count++
			t.AppendSeparator()
		} else if result[i].Weight != result[i-1].Weight {
			count++
			t.AppendSeparator()
		}
		if result[i].Url == "" {
			continue
		}
		t.AppendRow([]interface{}{result[i].Rank, result[i].Url, result[i].Weight})

	}

	func() {
		exOpts := exporter.ExporterOpts{
			FileName:         r.Exporter.FileName,
			Dir:              r.Exporter.Dir,
			Data:             result,
			CreatedAt:        time.Now(),
			Domain:           r.Crawler.DomainName,
			NumOfURLsCrawled: r.Crawler.Count,
			UserAgent:        r.Crawler.Spider.UserAgent,
		}

		var err error

		switch r.Exporter.Type {
		case "json":
			jx := exporter.JSONExporter{}
			err = jx.Init(exOpts).Export()
		case "csv":
			cx := exporter.CSVExporter{}
			err = cx.Init(exOpts).Export()
		}
		if err != nil {
			fmt.Println(err)
		}
	}()

	t.Render()

}

func removeDomain(url string) string {
	if url == "" {
		return url
	}
	if url[len(url)-1] != '/' {
		url = url + "/"
	}

	parts := strings.SplitN(url, "/", 4)

	if len(parts) < 4 {
		return "/"
	}

	return "/" + parts[3]
}
