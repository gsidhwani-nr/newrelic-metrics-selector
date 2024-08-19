package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Define CLI flags
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		fmt.Println("Usage: nrms [options]")
		fmt.Println("Options:")
		fmt.Println("  --help     Show help message")
		fmt.Println("Environment Variables:")
		fmt.Println("  NEW_RELIC_API_KEY     Your New Relic API key")
		fmt.Println("  NEW_RELIC_ACCOUNT_ID  Your New Relic account ID")
		return
	}

	// Ensure the API key is set
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	if apiKey == "" {
		log.Fatal("NEW_RELIC_API_KEY environment variable is not set")
	}

	// Ensure the account ID is set
	accountIDStr := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	if accountIDStr == "" {
		log.Fatal("NEW_RELIC_ACCOUNT_ID environment variable is not set")
	}

	// Mask the API key for logging
	maskedAPIKey := maskAPIKey(apiKey)
	log.Infof("Using API key: %s", maskedAPIKey)

	// Print the Account ID
	fmt.Printf("Account ID: %s\n", accountIDStr)

	// Initialize the client with the production API base URL.
	client, err := newrelic.New(
		newrelic.ConfigPersonalAPIKey(apiKey),
		// newrelic.ConfigBaseURL("https://staging-api.newrelic.com/graphql"), // Uncomment this line for staging
	)
	if err != nil {
		log.Fatal("error initializing client:", err)
	}
	log.Info("New Relic client initialized successfully")

	// Start a spinner to indicate processing
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	s.Start()                                                   // Start the spinner

	// Step 1: Fetch all Prometheus metrics
	prometheusMetrics, err := fetchPrometheusMetrics(client)
	if err != nil {
		log.Fatal("error fetching Prometheus metrics:", err)
	}
	log.Infof("Fetched %d Prometheus metrics", len(prometheusMetrics))

	// Step 2: Fetch all dashboard definitions
	dashboardQueries, err := fetchDashboardQueries(client, accountIDStr)
	if err != nil {
		log.Fatal("error fetching dashboard queries:", err)
	}
	log.Infof("Fetched %d dashboard queries", len(dashboardQueries))

	// Step 3: Fetch all alert definitions
	alertQueries, err := fetchAlertQueries(client, accountIDStr)
	if err != nil {
		log.Fatal("error fetching alert queries:", err)
	}
	log.Infof("Fetched %d alert queries", len(alertQueries))

	// Step 4: Identify used and unused metrics
	usedMetrics := make(map[string]bool)
	for _, query := range append(dashboardQueries, alertQueries...) {
		for _, metric := range prometheusMetrics {
			if strings.Contains(query, metric) {
				usedMetrics[metric] = true
			}
		}
	}

	// Get current timestamp
	timestamp := time.Now().Format("20060102150405")

	// Step 5: Output used and unused metrics to files
	usedFilename := fmt.Sprintf("%s_used_%s.txt", accountIDStr, timestamp)
	unusedFilename := fmt.Sprintf("%s_unused_%s.txt", accountIDStr, timestamp)

	writeMetricsToFile(usedFilename, prometheusMetrics, usedMetrics, true)
	writeMetricsToFile(unusedFilename, prometheusMetrics, usedMetrics, false)

	// Stop the spinner
	s.Stop()

	fmt.Println("Processing complete. Check the output files for details.")
}

func writeMetricsToFile(filename string, metrics []string, usedMetrics map[string]bool, used bool) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating file %s: %v", filename, err)
	}
	defer file.Close()

	for _, metric := range metrics {
		if usedMetrics[metric] == used {
			_, err := file.WriteString(metric + "\n")
			if err != nil {
				log.Fatalf("Error writing to file %s: %v", filename, err)
			}
		}
	}

	log.Infof("Metrics written to %s", filename)
}
