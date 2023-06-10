package exporter

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Avash027/Intra-Page-Ranker/matrix"
)

type CSVExporter struct {
	Dir      string
	FileName string
	Data     []matrix.RankData
	MetaData MetaData
}

func (cex *CSVExporter) Export() error {
	fmt.Println(cex)
	filePath := filepath.Join(cex.Dir, cex.FileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write metadata rows
	err = writer.Write([]string{"Domain", cex.MetaData.Domain})
	if err != nil {
		return fmt.Errorf("failed to write CSV metadata: %v", err)
	}
	err = writer.Write([]string{"Crawled At", cex.MetaData.CreatedAt.String()})
	if err != nil {
		return fmt.Errorf("failed to write CSV metadata: %v", err)
	}
	err = writer.Write([]string{"UserAgent", cex.MetaData.UserAgent})
	if err != nil {
		return fmt.Errorf("failed to write CSV metadata: %v", err)
	}
	err = writer.Write([]string{"Number of URLs crawled", fmt.Sprintf("%d", cex.MetaData.NumOfURLsCrawled)})
	if err != nil {
		return fmt.Errorf("failed to write CSV metadata: %v", err)
	}

	// Write header row
	header := []string{"URL", "Rank", "Weight"}
	err = writer.Write(header)
	if err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write data rows
	for _, rd := range cex.Data {
		if rd.Url == "" {
			continue
		}
		row := []string{rd.Url, fmt.Sprintf("%d", rd.Rank), fmt.Sprintf("%f", rd.Weight)}
		err = writer.Write(row)
		if err != nil {
			return fmt.Errorf("failed to write CSV row: %v", err)
		}
	}

	return nil
}

func (cx *CSVExporter) Init(ex ExporterOpts) *CSVExporter {
	return &CSVExporter{
		Dir:      ex.Dir,
		FileName: ex.FileName,
		Data:     ex.Data,
		MetaData: MetaData{
			Domain:           ex.Domain,
			CreatedAt:        ex.CreatedAt,
			UserAgent:        ex.UserAgent,
			NumOfURLsCrawled: ex.NumOfURLsCrawled,
		},
	}
}

func (cx *CSVExporter) GetDir() string {
	return cx.Dir
}

func (cx *CSVExporter) GetFileName() string {
	return cx.FileName
}
