package es

import (
	"context"
	"time"

	"github.com/olivere/elastic"
)

// NewElasticsearch Client
func NewElasticsearch(host, port, user, pass string) (*elastic.Client, error) {
	url := "http://" + host + ":" + port
	esclient, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetBasicAuth(user, pass),
		elastic.SetHealthcheckTimeoutStartup(15*time.Second),
	)
	if err != nil {
		return nil, err
	}
	return esclient, nil
}

// StartElastic runs a elasticsearch container.
func CreateElasticIndex(client *elastic.Client) error {
	// log.Println("Init indexes")
	mappings := map[string]string{
		"mediawatch_twitter_elections_tweets": indexTweets,
		// "articles/_mapping/document": "elasticsearch/mapping.articles.json",
		"mediawatch_twitter_elections_users": indexUsers,
		// "relationships/_mapping/document": "elasticsearch/mapping.relationships.json",
	}

	ctx := context.Background()
	for k, v := range mappings {
		createIndex, err := client.CreateIndex(k).BodyJson(v).Do(ctx)
		if err != nil {
			// log.Printf("Error creating mapping %s from file %s: %v", k, v, err)
			continue
		}
		if !createIndex.Acknowledged {
			// log.Printf("Error mapping %s from file %s not acknowledged", k, v)
			continue
		}
	}

	return nil
}
