package user

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cvcio/elections-api/pkg/db"
	"github.com/pkg/errors"
	"github.com/plagiari-sm/mediawatch/pkg/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opencensus.io/trace"
)

const usersCollection = "users"

var (
	// ErrNotFound abstracts the  not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	ErrAuthenticationFailure = errors.New("Authentication failed")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// EnsureIndex fix the indexes in the account collections
func EnsureIndex(ctx context.Context, dbConn *db.DB) error {
	// index := .Index{
	// 	Key:    []string{"email"},
	// 	Unique: true,
	// }
	// mongo.IndexView.CreateOne()

	index := []mongo.IndexModel{mongo.IndexModel{
		Keys:    "screenName",
		Options: options.Index().SetUnique(true), // {Unique: true},
		// Keys: bson.D
	},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	f := func(collection *mongo.Collection) error {
		_, err := collection.Indexes().CreateMany(ctx, index, opts) //EnsureIndex(index)
		return err
	}
	if err := dbConn.Execute(ctx, usersCollection, f); err != nil {
		return errors.Wrap(err, "db.users.ensureIndex()")
	}
	return nil
}

// TokenGenerator is the behavior we need in our Authenticate to generate
// tokens for authenticated users.
type TokenGenerator interface {
	GenerateToken(auth.Claims) (string, error)
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Token that can be used to authenticate in the future.
//
// The key, keyID, and alg are required for generating the token.
func Authenticate(ctx context.Context, tknGen TokenGenerator, now time.Time, u *User) (Token, error) {
	_, span := trace.StartSpan(ctx, "models.user.Authenticate")
	defer span.End()

	// q := bson.M{"email": email}

	// var u *Account
	// f := func(collection *.Collection) error {
	// 	return collection.Find(q).One(&u)
	// }

	// if err := dbConn.Execute(ctx, accountsCollection, f); err != nil {

	// 	// Normally we would return ErrNotFound in this scenario but we do not want
	// 	// to leak to an unauthenticated user which emails are in the system.
	// 	if err == .ErrNotFound {
	// 		return Token{}, ErrAuthenticationFailure
	// 	}
	// 	return Token{}, errors.Wrap(err, fmt.Sprintf("db.accounts.find(%s)", db.Query(q)))
	// }

	// // Compare the provided password with the saved hash. Use the bcrypt
	// // comparison function so it is cryptographically secure.
	// if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
	// 	return Token{}, ErrAuthenticationFailure
	// }

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.NewClaims(u.ID.Hex(), u.Email, u.Roles, u.OrgName, now, 24*time.Hour)

	tkn, err := tknGen.GenerateToken(claims)
	if err != nil {
		return Token{}, errors.Wrap(err, "generating token")
	}

	return Token{Token: tkn}, nil
}

// User : User Schema model
type User struct {
	ID                       primitive.ObjectID `bson:"_id" json:"id"`
	IDStr                    string             `bson:"idStr" json:"idStr"`
	UserID                   string             `bson:"userId" json:"userId"`
	CreatedAt                time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt                time.Time          `json:"updatedAt" bson:"updatedAt"`
	Deleted                  bool               `json:"-" bson:"deleted"`
	Roles                    []string           `bson:"roles" json:"-"`
	Status                   string             `json:"status" bson:"status"`
	ScreenName               string             `bson:"screenName" json:"screenName"`
	FirstName                string             `bson:"firstName" json:"firstName"`
	LastName                 string             `bson:"lastName" json:"lastName"`
	Email                    string             `bson:"email" json:"email"`
	Mobile                   string             `bson:"mobile" json:"mobile"`
	Profession               string             `bson:"profession" json:"profession"`
	Pin                      string             `bson:"pin" json:"-"`
	ProfileImageURL          *string            `bson:"profileImageURL" json:"profileImageURL"`
	TwitterAccessToken       *string            `bson:"twitterAccessToken" json:"-"`
	TwitterAccessTokenSecret *string            `bson:"twitterAccessTokenSecret" json:"-"`

	OrgName      string `bson:"orgname" json:"orgname"`
	OrgURL       string `bson:"orgurl" json:"orgurl"`
	OrgFirstName string `bson:"orgfirstName" json:"orgfirstName"`
	OrgLastName  string `bson:"orglastName" json:"orglastName"`
	OrgEmail     string `bson:"orgemail" json:"orgemail"`
	OrgMobile    string `bson:"orgmobile" json:"orgmobile"`
}

// Token is the payload we deliver to users when they authenticate.
type Token struct {
	Token string `json:"token"`
}

// Create inserts a new user into the database.
func Create(dbConn *db.DB, nu *User) (*User, error) {
	ctx, span := trace.StartSpan(context.Background(), "models.user.Create")
	defer span.End()

	// Mongo truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now := time.Now().Truncate(time.Millisecond)

	if len(nu.Roles) == 0 {
		nu.Roles = []string{auth.RoleUser}
	}

	nu.ID = primitive.NewObjectID()
	nu.IDStr = nu.ID.Hex()
	nu.CreatedAt = now
	nu.UpdatedAt = now
	nu.Status = "authorize"

	f := func(collection *mongo.Collection) error {
		_, err := collection.InsertOne(ctx, &nu) // (&u)
		return err
	}
	if err := dbConn.Execute(ctx, usersCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.users.insert(%s)", db.Query(&nu)))
	}

	return nu, nil
}

// Update replaces a user document in the database.
func Update(dbConn *db.DB, id string, upd *User) (*User, error) {
	ctx, span := trace.StartSpan(context.Background(), "models.user.Update")
	defer span.End()

	// Mongo truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now := time.Now().Truncate(time.Millisecond)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	/*
		cachedUser, err := Get(ctx, dbConn, id)
		if err != nil {
			return ErrNotFound
		}
	*/

	// log.Debugf("\n\n%+v\n\n", upd)

	fields := make(bson.M)

	if &upd.Status != nil {
		fields["status"] = upd.Status
	}

	if &upd.Status != nil && upd.Status == "create" {
		fields["pin"] = RandStringBytes(8)
	}

	if &upd.FirstName != nil {
		fields["firstName"] = upd.FirstName
	}
	if &upd.LastName != nil {
		fields["lastName"] = upd.LastName
	}
	if &upd.Email != nil {
		fields["email"] = upd.Email
	}
	if &upd.Mobile != nil {
		fields["mobile"] = upd.Mobile
	}
	if &upd.Profession != nil {
		fields["profession"] = upd.Profession
	}
	if &upd.OrgName != nil {
		fields["orgname"] = upd.OrgName
	}
	if &upd.OrgURL != nil {
		fields["orgurl"] = upd.OrgURL
	}
	if &upd.OrgFirstName != nil {
		fields["orgfirstName"] = upd.OrgFirstName
	}
	if &upd.OrgLastName != nil {
		fields["orglastName"] = upd.OrgLastName
	}
	if &upd.OrgEmail != nil {
		fields["orgemail"] = upd.OrgEmail
	}
	if &upd.OrgMobile != nil {
		fields["orgmobile"] = upd.OrgMobile
	}

	if upd.ProfileImageURL != nil {
		fields["profileImageURL"] = upd.ProfileImageURL
	}
	if upd.TwitterAccessToken != nil {
		fields["twitterAccessToken"] = upd.TwitterAccessToken
	}
	if upd.TwitterAccessTokenSecret != nil {
		fields["twitterAccessTokenSecret"] = upd.TwitterAccessTokenSecret
	}

	// If there's nothing to update we can quit early.
	if len(fields) == 0 {
		return nil, nil
	}

	fields["updatedAt"] = now

	update := bson.M{"$set": fields}
	filter := bson.M{"_id": oid, "deleted": false} // bson.ObjectIdHex(id)}
	res, err := dbConn.Database.Collection(usersCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments { //ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.account.update(%s, %s)", db.Query(filter), db.Query(update)))
	}
	/*
		f := func(collection *mongo.Collection) (*mongo.UpdateResult, error) {
			res, err := collection.UpdateOne(ctx, filter, update)
			return res, err
		}

		res, err := dbConn.Execute(ctx, usersCollection, f)
		if err != nil {
			if err == mongo.ErrNoDocuments { //ErrNotFound {
				return nil, ErrNotFound
			}
			return nil, errors.Wrap(err, fmt.Sprintf("db.account.update(%s, %s)", db.Query(filter), db.Query(update)))
		}
	*/
	fmt.Printf("%+v", res)

	return nil, nil
}

// ByScreenNanme : Returns User by Twitter Screen Name
func ByScreenNanme(dbConn *db.DB, screenName string) (*User, error) {
	ctx, span := trace.StartSpan(context.Background(), "models.users.ByScreenNanme")
	defer span.End()

	filter := bson.M{"screenName": screenName, "deleted": false}

	var u *User

	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&u) // (q).One(&u)
	}

	if err := dbConn.Execute(ctx, usersCollection, f); err != nil {
		// Normally we would return ErrNotFound in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.accounts.find(%s)", db.Query(filter)))
	}
	return u, nil
}

// Get gets the specified user from the database.
func Get(dbConn *db.DB, id string) (*User, error) {
	ctx, span := trace.StartSpan(context.Background(), "models.users.Get")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	filter := bson.M{"_id": oid, "deleted": false}

	var u *User

	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&u) // Find(q).One(&u)
	}

	if err := dbConn.Execute(ctx, usersCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { // .ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.users.find(%s)", db.Query(filter)))
	}

	return u, nil
}

// SendOTP to User with Twilio
func SendOTP(accountSid string, authToken string, mobile string, pin string) {
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", "+"+mobile)
	msgData.Set("From", "MediaWatch")
	msgData.Set("Body", "Your MediaWatch Verification Code is "+pin+".")
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
}

// RandStringBytes Generates Random String by length
func RandStringBytes(n int) string {
	const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return strings.ToUpper(string(b))
}
