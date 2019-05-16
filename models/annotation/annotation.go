package annotation

import (
	"context"
	"fmt"
	"time"

	"github.com/cvcio/elections-api/pkg/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

const annotationsCollection = "annotations"

// Annotation : Annotation Schema model
type Annotation struct {
	ID                   primitive.ObjectID `bson:"_id" json:"id"`
	IDStr                string             `bson:"idStr" json:"idStr"`
	CreatedAt            time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt            time.Time          `json:"updatedAt" bson:"updatedAt"`
	Deleted              bool               `json:"-" bson:"deleted"`
	AnnotatorIDStr       string             `bson:"annotatorIdStr" json:"annotatorIdStr"`
	AnnotatorScreenName  string             `bson:"annotatorScreenName" json:"annotatorScreenName"`
	UserIDStr            string             `bson:"userIdStr" json:"userIdStr"`
	UserScreenName       string             `bson:"userScreenName" json:"userScreenName"`
	AccountType          string             `json:"accountType" bson:"accountType"`
	PoliticalOrientation string             `json:"politicalOrientation" bson:"politicalOrientation"`
	Context              []string           `json:"context" bson:"context"`
	Note                 string             `json:"note" bson:"note"`
}

// Create inserts a new user into the database.
func Create(dbConn *db.DB, a *Annotation) (*Annotation, error) {
	ctx, span := trace.StartSpan(context.Background(), "models.annotation.Create")
	defer span.End()

	// Mongo truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now := time.Now().Truncate(time.Millisecond)

	a.ID = primitive.NewObjectID()
	a.IDStr = a.ID.Hex()
	a.CreatedAt = now
	a.UpdatedAt = now

	f := func(collection *mongo.Collection) error {
		_, err := collection.InsertOne(ctx, &a)
		return err
	}
	if err := dbConn.Execute(ctx, annotationsCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.annotations.insert(%s)", db.Query(&a)))
	}

	return a, nil
}
