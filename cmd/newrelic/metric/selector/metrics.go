package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	log "github.com/sirupsen/logrus"
)

func fetchPrometheusMetrics(client *newrelic.NewRelic) ([]string, error) {
	resp := struct {
		Actor struct {
			Account struct {
				ID int `json:"id"`
			} `json:"account"`
			Nrql struct {
				Results []struct {
					Uniques []string `json:"uniques.metricName"`
				} `json:"results"`
			} `json:"nrql"`
		} `json:"actor"`
	}{}

	accountIDStr := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	if accountIDStr == "" {
		return nil, fmt.Errorf("NEW_RELIC_ACCOUNT_ID environment variable is not set")
	}

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		log.Error("Invalid account ID: ", err)
		return nil, err
	}

	vars := map[string]interface{}{
		"accountId": accountID,
	}

	var metrics []string
	limit := 1000
	lastFetchedMetric := ""

	for {
		whereClause := ""
		if lastFetchedMetric != "" {
			whereClause = fmt.Sprintf("AND metricName > '%s'", lastFetchedMetric)
		}

		query := fmt.Sprintf("SELECT uniques(metricName) FROM Metric WHERE (instrumentation.name = 'remote-write') AND (instrumentation.provider = 'prometheus') %s LIMIT %d", whereClause, limit)

		log.Debug("Executing NRQL query to fetch Prometheus metrics")
		err = client.NerdGraph.QueryWithResponse(`query($accountId: Int!) {
			actor {
				account(id: $accountId) {
					id
				}
				nrql(query: "`+query+`", accounts: [$accountId]) {
					results
				}
			}
		}`, vars, &resp)
		if err != nil {
			log.Error("Error executing NRQL query: ", err)
			return nil, err
		}

		fetchedMetrics := 0
		for _, result := range resp.Actor.Nrql.Results {
			metrics = append(metrics, result.Uniques...)
			fetchedMetrics += len(result.Uniques)
		}

		if fetchedMetrics < limit {
			break
		}

		lastFetchedMetric = metrics[len(metrics)-1]
		log.Debugf("Fetching next page of results with metricName greater than: %s", lastFetchedMetric)
	}

	log.Debugf("Fetched %d unique Prometheus metrics", len(metrics))
	return metrics, nil
}
