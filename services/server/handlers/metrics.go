package handlers

import (
	"context"

	"github.com/cvcio/elections-api/pkg/config"
	"github.com/cvcio/elections-api/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

// Metrics Controller
type Metrics struct {
	cfg *config.Config
	db  *db.DB
	es  *elastic.Client
}

// GetVolumeByUser ...
func (ctrl *Metrics) GetVolumeByUser(c *gin.Context) {
	log.Infof("%s", c.Param("id"))

	agg := elastic.NewDateHistogramAggregation().
		Field("crawled_at").Interval("hour")

	res, err := ctrl.es.Search().
		Index("mediawatch_twitter_elections_users").
		Type("document").Query(elastic.NewTermQuery("screen_name", c.Param("id"))).
		Aggregation("histogram", agg).
		Do(context.Background())

	if err != nil {
		ResponseError(c, 404, err.Error())
		return
	}

	var results = []elastic.Aggregations{}
	if agg, found := res.Aggregations.Terms("histogram"); found {
		for _, bucket := range agg.Buckets {
			results = append(results, bucket.Aggregations)
		}
	}

	ResponseJSON(c, results)
}
