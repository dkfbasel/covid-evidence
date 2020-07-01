package ninox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// UpdateRecords will update the given records in the ninox database
func UpdateRecords(url string, records []*Record) {

	// define transport properties
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	// initialie a new http client
	client := &http.Client{Transport: tr}

	payload, err := json.Marshal(records)
	if err != nil {
		log.Fatalf("could not encode records for import: %+v", err)
	}

	// define a new request with corresponding authentication header
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ninoxAPIKey))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		log.Fatalf("could not create post request: %+v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("could not perform post request: %+v", err)
	}
	defer resp.Body.Close() // nolint:errcheck

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("could not read response: %+v", err)
	}

	log.Printf("%s", content)

}
