package exporter

import (
	"time"

	"github.com/Avash027/Intra-Page-Ranker/matrix"
)

type RexOpts struct {
	Type     string
	Dir      string
	FileName string
}

type ExporterOpts struct {
	Type             string
	Dir              string
	FileName         string
	Data             []matrix.RankData
	Domain           string
	CreatedAt        time.Time
	UserAgent        string
	NumOfURLsCrawled int
}
