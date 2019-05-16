package es

import (
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
