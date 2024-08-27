package main

import (
	"strconv"

	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	log "github.com/sirupsen/logrus"
)

func fetchAlertQueries(client *newrelic.NewRelic, maskedAPIKey, accountIDStr string) ([]string, error) {
	log.Debugf("Using API key in fetchAlertQueries: %s", maskedAPIKey)

	query := `query($accountId: Int!, $cursor: String) {
		actor {
			account(id: $accountId) {
				alerts {
					nrqlConditionsSearch(cursor: $cursor) {
						nextCursor
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
		"cursor":    nil,
	}

	var allConditions []struct {
		Nrql struct {
			Query string `json:"query"`
		} `json:"nrql"`
	}

	for {
		resp := struct {
			Actor struct {
				Account struct {
					Alerts struct {
						NrqlConditionsSearch struct {
							NextCursor     *string `json:"nextCursor"`
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

		allConditions = append(allConditions, resp.Actor.Account.Alerts.NrqlConditionsSearch.NrqlConditions...)

		if resp.Actor.Account.Alerts.NrqlConditionsSearch.NextCursor == nil {
			break
		}
		vars["cursor"] = *resp.Actor.Account.Alerts.NrqlConditionsSearch.NextCursor
	}

	var queries []string
	for _, condition := range allConditions {
		queries = append(queries, condition.Nrql.Query)
	}
	log.Debugf("Fetched %d alert queries", len(queries))

	return queries, nil
}
