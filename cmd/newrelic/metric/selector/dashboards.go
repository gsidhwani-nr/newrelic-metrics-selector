package main

import (
	"strconv"

	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	log "github.com/sirupsen/logrus"
)

func fetchDashboardQueries(client *newrelic.NewRelic, maskedAPIKey, accountIDStr string) ([]string, error) {
	//log.Infof("Using API key in fetchDashboardQueries: %s", maskedAPIKey)

	query := `query($accountId: Int!, $cursor: String) {
		actor {
			entitySearch(query: "type = 'DASHBOARD'") {
				results(cursor: $cursor) {
					entities {
						... on DashboardEntityOutline {
							name
							guid
							account {
								id
							}
						}
					}
					nextCursor
				}
			}
			account(id: $accountId) {
				id
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

	var allEntities []struct {
		Name    string `json:"name"`
		GUID    string `json:"guid"`
		Account struct {
			ID int `json:"id"`
		} `json:"account"`
	}

	for {
		resp := struct {
			Actor struct {
				EntitySearch struct {
					Results struct {
						Entities []struct {
							Name    string `json:"name"`
							GUID    string `json:"guid"`
							Account struct {
								ID int `json:"id"`
							} `json:"account"`
						} `json:"entities"`
						NextCursor *string `json:"nextCursor"`
					} `json:"results"`
				} `json:"entitySearch"`
				Account struct {
					ID int `json:"id"`
				} `json:"account"`
			} `json:"actor"`
		}{}

		log.Debug("Executing GraphQL query to fetch dashboard GUIDs")
		err = client.NerdGraph.QueryWithResponse(query, vars, &resp)
		if err != nil {
			log.Error("Error executing GraphQL query: ", err)
			return nil, err
		}

		for _, entity := range resp.Actor.EntitySearch.Results.Entities {
			if entity.Account.ID == accountID {
				allEntities = append(allEntities, entity)
				log.Debugf("Fetched dashboard Name: %s : %s ", entity.Name, entity.GUID)
			}
		}
		if resp.Actor.EntitySearch.Results.NextCursor == nil {
			break
		}
		vars["cursor"] = *resp.Actor.EntitySearch.Results.NextCursor
	}
	total := len(allEntities)
	log.Infof("Fetched %d dashboards", total)

	fTotal := float64(total)
	log.Infof("Processing will take approximately %.1f minutes", fTotal*0.55/60.0)

	var queries []string
	for _, entity := range allEntities {
		dashboardQueries, err := fetchDashboardDetails(client, entity.GUID)
		if err != nil {
			log.Error("Error fetching dashboard details: ", err)
			continue
		}
		queries = append(queries, dashboardQueries...)
	}
	log.Debugf("Fetched %d dashboard queries", len(queries))

	return queries, nil
}

func fetchDashboardDetails(client *newrelic.NewRelic, guid string) ([]string, error) {
	query := `query GetDashboardEntityQuery($entityGuid: EntityGuid!) {
		actor {
			entity(guid: $entityGuid) {
				... on DashboardEntity {
					pages {
						widgets {
							rawConfiguration
						}
					}
				}
			}
		}
				}`

	vars := map[string]interface{}{
		"entityGuid": guid,
	}

	resp := struct {
		Actor struct {
			Entity struct {
				Pages []struct {
					Widgets []struct {
						RawConfiguration struct {
							NrqlQueries []struct {
								Query string `json:"query"`
							} `json:"nrqlQueries"`
						} `json:"rawConfiguration"`
					} `json:"widgets"`
				} `json:"pages"`
			} `json:"entity"`
		} `json:"actor"`
	}{}

	log.Debug("Executing GraphQL query to fetch dashboard details")
	err := client.NerdGraph.QueryWithResponse(query, vars, &resp)
	if err != nil {
		log.Error("Error executing GraphQL query: ", err)
		return nil, err
	}

	var queries []string
	for _, page := range resp.Actor.Entity.Pages {
		for _, widget := range page.Widgets {
			for _, nrqlQuery := range widget.RawConfiguration.NrqlQueries {
				queries = append(queries, nrqlQuery.Query)
			}
		}
	}
	log.Debugf("Fetched %d NRQL queries from dashboard %s", len(queries), guid)

	return queries, nil
}
