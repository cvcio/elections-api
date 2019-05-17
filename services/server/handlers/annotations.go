package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/cvcio/elections-api/models/annotation"
	"github.com/cvcio/elections-api/pkg/config"
	"github.com/cvcio/elections-api/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
)

// Annotations Controller
type Annotations struct {
	cfg *config.Config
	db  *db.DB
	es  *elastic.Client
}

// Create New Annotation
func (ctrl *Annotations) Create(c *gin.Context) {
	var a annotation.Annotation
	if err := c.Bind(&a); err != nil {
		ResponseError(c, 406, err.Error())
		return
	}
	res, err := annotation.Create(ctrl.db, &a)
	if err != nil {
		ResponseError(c, 401, err.Error())
		return
	}
	ResponseJSON(c, res.IDStr)
}

// GetRandom Gets a random document from ES
func (ctrl *Annotations) GetRandom(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())

	q := elastic.NewFunctionScoreQuery().
		AddScoreFunc(elastic.NewRandomFunction().Seed(rand.Intn(1000000))).
		ScoreMode("avg")

	res, err := ctrl.es.Search().
		Index("mediawatch_twitter_elections_tweets").
		Type("document").Query(q).Size(1).
		Do(context.Background())

	if err != nil {
		ResponseError(c, 404, err.Error())
		return
	}

	if len(res.Hits.Hits) > 0 {
		type T struct {
			Text      string   `json:"full_text"`
			CreatedAt string   `json:"created_at"`
			Media     []string `json:"media"`
			Urls      []string `json:"urls"`
		}
		var tweet *anaconda.Tweet
		err := json.Unmarshal(*res.Hits.Hits[0].Source, &tweet)
		if err != nil {
			ResponseError(c, 500, err.Error())
		}
		response := &T{
			Text:      tweet.FullText,
			CreatedAt: tweet.CreatedAt,
		}
		for u := range tweet.Entities.Urls {
			response.Urls = append(response.Urls, tweet.Entities.Urls[u].Expanded_url)
		}
		for m := range tweet.Entities.Media {
			response.Media = append(response.Media, tweet.Entities.Media[m].Media_url_https)
		}
		ResponseJSON(c, response)
		return
	}
	ResponseJSON(c, res)
}
