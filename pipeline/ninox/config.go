package ninox

import (
	"fmt"
	"log"
	"os"
)

var ninoxAPIKey = "MISSING"

const CoveBasicURL = "https://api.ninoxdb.de/v1/teams/JaSodfHneNLbnZKHb/databases/bhdh22vn3oqj/tables/A/records"
const CoveBasicExlusionURL = "https://api.ninoxdb.de/v1/teams/JaSodfHneNLbnZKHb/databases/bhdh22vn3oqj/tables/B/records"

const ClinicaltrialsURL = "https://api.ninoxdb.de/v1/teams/JaSodfHneNLbnZKHb/databases/ogt4txmvycpz/tables/C/records"
const IctrpURL = "https://api.ninoxdb.de/v1/teams/JaSodfHneNLbnZKHb/databases/ogt4txmvycpz/tables/E/records"
const MedrxivURL = "https://api.ninoxdb.de/v1/teams/JaSodfHneNLbnZKHb/databases/ogt4txmvycpz/tables/F/records"
const SwissethicsURL = "https://api.ninoxdb.de/v1/teams/JaSodfHneNLbnZKHb/databases/ogt4txmvycpz/tables/G/records"

const CoveBasicTable = "covebasic"
const CoveBasicExlusionsTable = "exclusions"

func init() {
	fromEnv := os.Getenv("NINOX_API_KEY")
	if fromEnv != "" {
		ninoxAPIKey = fromEnv
		return
	}

	fmt.Println("Enter ninox api key:")
	var apiKey string
	_, err := fmt.Scanln(&apiKey)
	if err != nil {
		log.Fatalln("Could not read api key")
	}

	if apiKey == "" {
		log.Fatalln("Ninox Api Key is required")
	}
	ninoxAPIKey = apiKey
}
