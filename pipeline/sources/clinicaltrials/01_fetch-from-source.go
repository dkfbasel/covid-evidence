package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type fetchResponse struct {
	FullStudiesResponse struct {
		APIVrs           string
		DataVrs          string
		Expression       string
		NStudiesAvail    int
		NStudiesFound    int
		MinRank          int
		MaxRank          int
		NStudiesReturned int
		FullStudies      []json.RawMessage
	}
}

// Fetch will fetch data from clinicaltrials.gov with the given query paramenter
// and return the data as string matrix (i.e. for csv tables)
func Fetch(query string, filename string) error {

	filename = fmt.Sprintf("%s.json", filename)

	// open a file to save the results
	output, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could open file for export: %w", err)
	}
	defer output.Close()

	// prepare the api url
	apiUrl, err := url.Parse("https://www.clinicaltrials.gov/api/query/full_studies")
	if err != nil {
		return fmt.Errorf("could not parse url: %w", err)
	}

	// add additional query params
	params := url.Values{}
	params.Add("expr", query)
	params.Add("fmt", "JSON")

	// save all studies exported from clinicaltrials.gov
	studies := []json.RawMessage{}

	// do not make more then 100 fetch attempts (just to ensure that we do not,
	// overload the clinicaltrials.gov api if something goes wrong on our side)
	currentFetchAttempt := 0
	maxFetchAttempts := 100

	// initialize variables to select the min and max ranking of the results
	// (required for pagination of the clinicaltrials.gov results)
	var currentRank, maxRank int

	for {

		// fetch the next batch of studies (maximum 100 per batch)
		currentRank = currentFetchAttempt*100 + 1
		maxRank = currentRank + 99
		params.Set("min_rnk", strconv.Itoa(currentRank))
		params.Set("max_rnk", strconv.Itoa(maxRank))

		// encode the url parameters
		apiUrl.RawQuery = params.Encode()

		// perform the request and parse the response
		response, err := performRequest(apiUrl.String())
		if err != nil {
			return fmt.Errorf("could not fetch data: %s, %w", apiUrl.String(), err)
		}

		// append the results to the previously fetched studies
		for _, study := range response.FullStudiesResponse.FullStudies {
			studies = append(studies, study)
		}

		// log some information
		log.Printf("MaxRank: % 4d, NStudiesAvailable: % 4d\n",
			response.FullStudiesResponse.MaxRank,
			response.FullStudiesResponse.NStudiesFound,
		)

		// stop fetching if the rank is higher
		if response.FullStudiesResponse.MaxRank >= response.FullStudiesResponse.NStudiesFound {
			break
		}

		currentFetchAttempt++
		if currentFetchAttempt > maxFetchAttempts {
			break
		}

	}

	// convert the information to json
	asJSON, err := json.Marshal(studies)
	if err != nil {
		return fmt.Errorf("could not marshal studies")
	}

	// write the information into the file
	_, err = output.Write(asJSON)
	if err != nil {
		return fmt.Errorf("could not write result to output")
	}
	output.Close() // nolint:errcheck

	return nil
}

// performRequest will perform the request for clinicaltrials.gov and return
// the parsed results
func performRequest(url string) (*fetchResponse, error) {

	request, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch data: %w", err)
	}

	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
	}
	request.Body.Close()

	var response fetchResponse
	err = json.Unmarshal(content, &response)
	if err != nil {
		return nil, fmt.Errorf("could not parse response: %w", err)
	}

	return &response, nil

}
