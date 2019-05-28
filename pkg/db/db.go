// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package db

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opencensus.io/trace"
)

// ErrInvalidDBProvided is returned in the event that an uninitialized db is
// used to perform actions against.
var ErrInvalidDBProvided = errors.New("invalid DB provided")

// DB is a collection of support for different DB technologies. Currently
// only MongoDB has been implemented. We want to be able to access the raw
// database support for the given DB so an interface does not work. Each
// database is too different.
type DB struct {

	// MongoDB Support.
	database *mongo.Database
}

// New returns a new DB value for use with MongoDB based on a registered
// master session.
func New(url string, database string, timeout time.Duration) (*DB, error) {

	// Set the default timeout for the session.
	// if timeout == 0 {
	// 	timeout = 60 * time.Second
	// }

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	// ctx, _ := context.WithTimeout(context.Background(), timeout)
	ctx := context.Background()
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	err = client.Connect(ctx)

	if err != nil {
		return nil, err
	}

	db := DB{
		database: client.Database(database),
		// session:  ses,
	}

	return &db, nil
}

// Close closes a DB value being used with MongoDB.
func (db *DB) Close() error {
	return db.database.Client().Disconnect(context.Background())
}

// Copy returns a new DB value for use with MongoDB based on master session.
// func (db *DB) Copy() *DB {
// 	ses := db.session.Copy()
func (db *DB) Copy() {
	// db.database.Client().UseSession()
}

// Execute is used to execute MongoDB commands.
func (db *DB) Execute(ctx context.Context, collName string, f func(*mongo.Collection) error) error {
	ctx, span := trace.StartSpan(ctx, "pkg.DB.Execute")
	defer span.End()

	if db == nil { //|| db.session == nil {
		return errors.Wrap(ErrInvalidDBProvided, "db == nil")
	}

	return f(db.database.Collection(collName))
}

// ExecuteTimeout is used to execute MongoDB commands with a timeout.
func (db *DB) ExecuteTimeout(ctx context.Context, timeout time.Duration, collName string, f func(*mongo.Collection) error) error {
	ctx, span := trace.StartSpan(ctx, "pkg.DB.ExecuteTimeout")
	defer span.End()

	if db == nil { //|| db.session == nil {
		return errors.Wrap(ErrInvalidDBProvided, "db == nil")
	}

	// db.session.SetSocketTimeout(timeout)

	return f(db.database.Collection(collName))
}

// StatusCheck validates the DB status good.
func (db *DB) StatusCheck(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "pkg.DB.StatusCheck")
	defer span.End()

	return nil
}

// Query provides a string version of the value
func Query(value interface{}) string {
	json, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return string(json)
}

// Valid returns true if a given id is a valid mongo id
func Valid(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	return true
}

type ListOpts struct {
	Limit   int
	Offset  int
	Org     string
	Deleted bool
	Status  string
	// Q      ListFunc
}

func Limit(i int) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Limit = i
	}
}

func Offset(i int) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Offset = i
	}
}

func Org(i string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Org = i
	}
}
func Deleted() func(*ListOpts) {
	return func(l *ListOpts) {
		l.Deleted = true
	}
}

func Status(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Status = s
	}
}

func DefaultOpts() ListOpts {
	l := ListOpts{}
	l.Offset = 0
	l.Limit = 24
	l.Deleted = false
	l.Status = ""
	return l
}

func NewListOpts() []func(*ListOpts) {
	return make([]func(*ListOpts), 0)
}

// // ListFunc type of
// type ListFunc func() bson.M

// // Nil queiry
// func Nil() ListFunc {
// 	return func() bson.M {
// 		return nil
// 	}
// }

// ID query by id
// func ID(id string) ListFunc {
// 	return func() bson.M {
// 		return bson.M{"_id": bson.ObjectIdHex(id)}
// 	}
// }

// SimpleM one key bson query
// func SimpleM(key string, value interface{}) ListFunc {
// 	q := bson.M{key: value}
// 	return func() bson.M {
// 		return q
// 	}
// }
