package ninox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// FetchRecords will fetch all record from the given ninox table
func FetchRecords(endpoint string, filters string) ([]Record, error) {

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
		return nil, fmt.Errorf("could not parse ninox url: %s, %w", endpoint, err)
	}

	params := url.Values{}
	// TODO: explore use of pagination
	params.Add("perPage", "9000")

	// add filters to the query
	if filters != "" {
		params.Add("filters", filters)
	}

	endpointURL.RawQuery = params.Encode()

	// define a new request with corresponding authentication header
	req, err := http.NewRequest("GET", endpointURL.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ninoxAPIKey))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}
	defer resp.Body.Close() // nolint:errcheck

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
	}

	// parse the content
	var records []Record
	err = json.Unmarshal(content, &records)
	if err != nil {
		return nil, fmt.Errorf("cound not parse response: %w,\n%s", err, content)
	}

	return records, nil

}
