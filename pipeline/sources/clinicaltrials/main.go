package main

import (
	"log"
)

func main() {

	// query := `(wuhan AND (coronavirus OR corona virus OR pneumonia virus)) OR COVID19 OR COVID-19 OR COVID 19 OR coronavirus 2019 OR corona virus 2019 OR SARS-CoV-2 OR SARSCoV2 OR SARS2 OR SARS-2 OR 2019 nCoV OR ((novel coronavirus OR novel corona virus) AND 2019)`

	// filename := fmt.Sprintf("./exports/clinicaltrials_%s", time.Now().Format("2006-01-02-150405"))
	// fmt.Println("FILE:" filename)

	// // fetch data from clinicaltrials.gov
	// err := Fetch(query, filename)
	// if err != nil {
	// 	log.Fatalf("could not fetch data: %+v", err)
	// }

	// // parse the data into a csv file
	// err = Parse(filename)
	// if err != nil {
	// 	log.Fatalf("could not parse data: %+v", err)
	// }

	// // compare the data with the data from ninox
	// err = Compare(filename)
	// if err != nil {
	// 	log.Fatalf("could not compare data with ninox table: %+v", err)
	// }

	// filename := "./exports/clinicaltrials_2020-06-12-061937"
	// fmt.Println(filename)

	// // import the data into ninox
	// err := Import(filename)
	// if err != nil {
	// 	log.Fatalf("could not import data into ninox: %+v", err)
	// }

	// ToCovebasic()

	log.Println("finished")

}
