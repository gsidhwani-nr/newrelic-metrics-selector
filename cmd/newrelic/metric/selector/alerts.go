package main

import (
	"strconv"

	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	log "github.com/sirupsen/logrus"
)

func fetchAlertQueries(client *newrelic.NewRelic, accountIDStr string) ([]string, error) {
	log.Infof("Fetching alert queries")

	query := `query($accountId: Int!) {
		actor {
			account(id: $accountId) {
				alerts {
					nrqlConditionsSearch {
						nrqlConditions {
							nrql {
								query
							}
						}
					}
				}
			}
		}
	}`

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		log.Error("Invalid account ID: ", err)
		return nil, err
	}

	vars := map[string]interface{}{
		"accountId": accountID,
	}

	resp := struct {
		Actor struct {
			Account struct {
				Alerts struct {
					NrqlConditionsSearch struct {
						NrqlConditions []struct {
							Nrql struct {
								Query string `json:"query"`
							} `json:"nrql"`
						} `json:"nrqlConditions"`
					} `json:"nrqlConditionsSearch"`
				} `json:"alerts"`
			} `json:"account"`
		} `json:"actor"`
	}{}

	log.Debug("Executing GraphQL query to fetch alert queries")
	err = client.NerdGraph.QueryWithResponse(query, vars, &resp)
	if err != nil {
		log.Error("Error executing GraphQL query: ", err)
		return nil, err
	}

	var queries []string
	for _, condition := range resp.Actor.Account.Alerts.NrqlConditionsSearch.NrqlConditions {
		queries = append(queries, condition.Nrql.Query)
	}
	log.Debugf("Fetched %d alert queries", len(queries))

	return queries, nil
}
