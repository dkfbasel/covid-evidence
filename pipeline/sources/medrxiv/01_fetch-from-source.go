package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type exportData struct {
	Gname string `json:"gname"`
	GrpID string `json:"grp_id"`
	Rels  []struct {
		RelTitle      string `json:"rel_title"`
		RelDoi        string `json:"rel_doi"`
		RelLink       string `json:"rel_link"`
		RelAbs        string `json:"rel_abs"`
		RelNumAuthors int    `json:"rel_num_authors"`
		RelAuthors    []struct {
			AuthorName string `json:"author_name"`
			AuthorInst string `json:"author_inst"`
		} `json:"rel_authors"`
		RelDate string `json:"rel_date"`
		RelSite string `json:"rel_site"`
	} `json:"rels"`
}

func Convert() {

	// download data from medrxiv

	timestamp := time.Now().Format("20060102-150405")

	response, err := http.Get("https://connect.medrxiv.org/relate/collection_json.php?grp=181")
	if err != nil {
		log.Println("could not fetch data")
		return
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("could not read content")
		return
	}

	ioutil.WriteFile(fmt.Sprintf("./exports/data_%s.json", timestamp), content, 0644)
	if err != nil {
		log.Println("coud not open file content")
		return
	}

	out, err := os.Create(fmt.Sprintf("./exports/data_%s.csv", timestamp))
	if err != nil {
		log.Println("could not open output file")
		return
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	writer.Comma = ';'

	writer.Write([]string{
		"rel_title", "rel_doi", "rel_link", "rel_abs", "rel_num_authors",
		"rel_authors", "rel_date", "rel_site",
	})

	// parse content
	var dta exportData
	err = json.Unmarshal(content, &dta)
	if err != nil {
		log.Println("could not parse content")
		return
	}

	for _, item := range dta.Rels {

		row := make([]string, 8)

		row[0] = item.RelTitle
		row[1] = item.RelDoi
		row[2] = item.RelLink
		row[3] = item.RelAbs
		row[4] = strconv.Itoa(item.RelNumAuthors)

		authors := []string{}

		for _, a := range item.RelAuthors {
			entry := fmt.Sprintf("%s (%s)", a.AuthorName, a.AuthorInst)
			authors = append(authors, entry)
		}

		row[5] = strings.Join(authors, "; ")

		row[6] = item.RelDate
		row[7] = item.RelSite

		writer.Write(row)

	}

	writer.Flush()

}
