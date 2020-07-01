package main

import (
	"log"
	"net/http"
	"time"

	"dkfbasel.ch/covid-evidence/sources/clinicaltrials"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

const currentFileName = "./exports/clinicaltrials_2020-04-17-111344"

func main() {

	var err error

	// _, err = clinicaltrials.Parse(currentFileName)
	// if err != nil {
	// 	log.Printf("could not parse data: %+v", err)
	// }

	// err = clinicaltrials.Compare(currentFileName)
	// if err != nil {
	// 	log.Printf("could not compare data: %+v", err)
	// }

	err = clinicaltrials.Import(currentFileName)
	if err != nil {
		log.Printf("could not import data: %+v", err)
	}

	// setup()f91b5fb0-7065-11ea-b81a-4182d593d91e
}

func setup() {
	// start a http server to handle requests
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/clinicaltrials/fetch", clinicaltrialsFetch)
	r.Get("/clinicaltrials/parse", clinicaltrialsParse)
	r.Get("/clinicaltrials/compare", clinicaltrialsCompare)

	http.ListenAndServe("0.0.0.0:8080", r)
}

// noop is used for all handlers that are not yet implemented
func noop(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("not yet implemented"))
}

// clinicaltrialsFetch will fetch data from clinicaltrials.gov using the
// given query
func clinicaltrialsFetch(w http.ResponseWriter, r *http.Request) {

	// define the query
	query := `(wuhan AND (coronavirus OR corona virus OR pneumonia virus)) OR COVID19 OR COVID-19 OR COVID 19 OR coronavirus 2019 OR corona virus 2019 OR SARS-CoV-2 OR SARSCoV2 OR SARS2 OR SARS-2 OR 2019 nCoV OR ((novel coronavirus OR novel corona virus) AND 2019)`

	dta, err := clinicaltrials.Fetch(query)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not fetch data", http.StatusInternalServerError)
		return
	}

	render.PlainText(w, r, dta)
}

// clinicaltrialsParse will parse the data from clinicaltrials.gov using
// the given matching table
func clinicaltrialsParse(w http.ResponseWriter, r *http.Request) {
	clinicaltrials.Parse(currentFileName)
}

// clinicaltrialsCompare will parse the data from clinicaltrials.gov using
// the given matching table
func clinicaltrialsCompare(w http.ResponseWriter, r *http.Request) {
	clinicaltrials.Compare(currentFileName)
}
