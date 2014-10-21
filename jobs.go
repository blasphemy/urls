package main

import (
	"log"
	"time"
)

func RunJobs() {
	for {
		log.Print("Running Jobs")
		log.Print("Running Total Links Update")
		t := time.Now()
		SetGetTotalUrlsFromScript()
		t2 := time.Since(t)
		log.Print("Total links update complete, took: ", t2)
		log.Print("Running Total Clicks Update")
		t = time.Now()
		SetGetTotalClicksFromScript()
		t2 = time.Since(t)
		log.Print("Total clicks update complete,  took: ", t2)
		log.Print("Running Clicks Per Url Update")
		t = time.Now()
		SetGetClicksPerUrl()
		t2 = time.Since(t)
		log.Print("Clicks Per Url complete, took:", t2)
		time.Sleep(config.GetJobInvertal())
	}
}
