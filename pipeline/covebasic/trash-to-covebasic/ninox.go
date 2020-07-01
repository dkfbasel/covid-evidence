package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var ninoxAPIKey = "MISSING"

const ninoxBasicURL = "https://api.ninoxdb.de/v1/teams/JaSodfHneNLbnZKHb/databases/bhdh22vn3oqj/tables/A/records"
const ninoxExlusionURL = "https://api.ninoxdb.de/v1/teams/JaSodfHneNLbnZKHb/databases/bhdh22vn3oqj/tables/B/records"

// NinoxRecord is used to handle data from ninoxDB
type NinoxRecord struct {
	BasicID    int
	ID         int                    `json:"id,omitempty"`
	Sequence   int                    `json:"sequence,omitempty"`
	CreatedAt  string                 `json:"createdAt,omitempty"`
	CreatedBy  string                 `json:"createdBy,omitempty"`
	ModifiedAt string                 `json:"modifiedAt,omitempty"`
	ModifiedBy string                 `json:"modifiedBy,omitempty"`
	Fields     map[string]interface{} `json:"fields"`
}

func init() {
	fromEnv := os.Getenv("NINOX_API_KEY")
	if fromEnv != "" {
		ninoxAPIKey = fromEnv
		return
	}

	fmt.Println("Enter ninox api key:")
	var apiKey string
	fmt.Scanln(&apiKey)

	if apiKey == "" {
		log.Fatalln("Ninox Api Key is required")
	}
	ninoxAPIKey = apiKey
}

func fetchRecords(endpoint string) []NinoxRecord {

	// define transport properties
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	// initialie a new http client
	client := &http.Client{Transport: tr}

	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalf("could not parse ninox url: %s, %+v", endpoint, err)
	}

	params := url.Values{}
	params.Add("perPage", "5000")

	endpointURL.RawQuery = params.Encode()

	// define a new request with corresponding authentication header
	req, err := http.NewRequest("GET", endpointURL.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ninoxAPIKey))
	if err != nil {
		log.Fatalf("could not create request: %+v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("could not perform request: %+v", err)
	}
	defer resp.Body.Close() // nolint:errcheck

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("could not read response: %+v", err)
	}

	// parse the content
	var records []NinoxRecord
	err = json.Unmarshal(content, &records)
	if err != nil {
		log.Fatalf("cound not parse response: %+v\n%s", err, content)
	}

	return records

}

func updateRecords(url string, records []*NinoxRecord) {

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

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("could not read response: %+v", err)
	}

	// log.Printf("%s", content)

}

// deleteRecords will delete the given records from the given database
func deleteRecords(url string, ids []int) {

	// define transport properties
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	// initialize a new http client
	client := &http.Client{Transport: tr}

	for _, id := range ids {

		// each record must be deleted individually
		deleteURL := fmt.Sprintf("%s/%d", url, id)

		log.Printf("DELETE: %s", deleteURL)

		// define a new request with corresponding authentication header
		req, err := http.NewRequest("DELETE", deleteURL, nil)
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

}
