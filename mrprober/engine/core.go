package engine

import (
	"fmt"
	"log"
	"mrprober/conf"
	"mrprober/probes"
	"sync"
	"time"
)

// OneShotRun execute all rules as probes and return a slice of results.
// Rules are executed in parallel and the function ends only after all
// results have been collected.
func OneShotRun() chan probes.Result {

	// Read rules from configuration
	config := conf.SafeConfiguration.Get()

	// goroutines pool
	var wg sync.WaitGroup

	// channel for storing results
	results := make(chan probes.Result, len(config.Rules))

	// Read rules
	for _, rule := range config.Rules {

		wg.Add(1)

		// Generate probe from rule
		p, err := probes.New(&rule)
		if err != nil {
			log.Fatal(err)
		}

		// Wrap the worker call in a closure that makes sure to tell the WaitGroup that this worker is done.
		// This way the worker itself does not have to be aware of the concurrency primitives involved in its execution.
		go func() {
			defer wg.Done()
			// Execute probe and collect result
			results <- p.Exec()
		}()
	}

	// Block until the WaitGroup counter goes back to 0; all the workers notified they’re done.
	wg.Wait()
	// Close channel
	close(results)

	return results
}

func StartActivePolling(quit chan int) {

	config := conf.SafeConfiguration.Get()
	rate, _ := time.ParseDuration(fmt.Sprintf("%ds", config.Global.Tickrate))

	// Execute all the probes and collect results to update status for prometheus metrics.
	// This first loop is done so that we don't have to wait for the initial iteration of the loop below.
	for r := range OneShotRun() {
		r.Update()
	}

	go func() {
		tick := time.Tick(rate)
		for {
			select {
			case <-tick:
				for r := range OneShotRun() {
					r.Update()
				}
				log.Print(".")
			case <-quit:
				log.Println("Stopped active polling")
				return
			default:
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()
}
