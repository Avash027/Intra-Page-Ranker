package exporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/Avash027/Intra-Page-Ranker/matrix"
)

type JSONExporter struct {
	Dir      string
	FileName string
	Data     []matrix.RankData
	MetaData MetaData
}

type MetaData struct {
	Domain           string
	CreatedAt        time.Time
	UserAgent        string
	NumOfURLsCrawled int
}

func (jex *JSONExporter) Export() error {

	mainData := make(map[string]interface{}, 0)
	mainData["Domain"] = jex.MetaData.Domain
	mainData["Crawled At"] = jex.MetaData.CreatedAt
	mainData["UserAgent"] = jex.MetaData.UserAgent
	mainData["Number of URLs crawled"] = jex.MetaData.NumOfURLsCrawled

	jsonData := make([]map[string]interface{}, 0)

	for _, rd := range jex.Data {
		if rd.Url == "" {
			continue
		}
		jsonData = append(jsonData, map[string]interface{}{
			"url":    rd.Url,
			"rank":   rd.Rank,
			"weight": rd.Weight,
		})
	}

	mainData["Url Data"] = jsonData

	jsonBytes, err := json.MarshalIndent(mainData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	filePath := filepath.Join(jex.Dir, jex.FileName)
	err = ioutil.WriteFile(filePath, jsonBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON data to file: %v", err)
	}

	return nil
}

func (jx *JSONExporter) Init(ex ExporterOpts) *JSONExporter {
	return &JSONExporter{
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
func (cx *JSONExporter) GetDir() string {
	return cx.Dir
}

func (cx *JSONExporter) GetFileName() string {
	return cx.FileName
}
