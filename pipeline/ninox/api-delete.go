package ninox

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// DeleteRecords will delete the given records from the given database
func DeleteRecords(url string, ids []int) {

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
