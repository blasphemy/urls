package main

import (
	"log"
	"time"
)

func RunJobs() {
	for {
		time.Sleep(time.Minute * 10)
		log.Print("Running Jobs")
		log.Print("Running Total Links Update")
		t := time.Now()
		SetTotalUrls()
		t2 := time.Since(t)
		log.Print("Total links update complete, took: ", t2)
	}
}
