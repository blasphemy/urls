package main

import (
	"log"
	"time"
)

func RunJobs() {
	for {
		time.Sleep(config.GetJobInvertal())
		log.Print("Running Jobs")
		log.Print("Running Total Links Update")
		t := time.Now()
		SetTotalUrls()
		t2 := time.Since(t)
		log.Print("Total links update complete, took: ", t2)
		log.Print("Running Total Clicks Update")
		t = time.Now()
		SetTotalClicks()
		t2 = time.Since(t)
		log.Print("Total clicks update complete,  took: ", t2)
	}
}
