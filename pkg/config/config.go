package config

import (
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Config struct holds all the configuration elements for our apps
type Config struct {
	Env string `envconfig:"ENV" default:"development"`
	Log struct {
		Level string `envconfig:"LOG_LEVEL" default:"debug"`
		Dev   bool   `envconfig:"Dev" default:"false"`
		Debug bool   `envconfig:"DEBUG" default:"true"`
	}
	Web struct {
		Host            string        `default:"localhost" envconfig:"HOST"`
		Port            string        `default:":8000" envconfig:"PORT"`
		ReadTimeout     time.Duration `default:"10s" envconfig:"READ_TIMEOUT"`
		WriteTimeout    time.Duration `default:"20s" envconfig:"WRITE_TIMEOUT"`
		ShutdownTimeout time.Duration `default:"10s" envconfig:"SHUTDOWN_TIMEOUT"`
		Debug           bool          `default:"true" envconfig:"DEBUG"`
		DomainName      string        `default:"mediawatch.io" envconfig:"DOMAIN_NAME"`
	}
	Streamer struct {
		Host   string   `default:"localhost" envconfig:"STREAMER_HOST"`
		Port   string   `default:"50050" envconfig:"STREAMER_PORT"`
		Follow []string `default:"" envconfig:"FOLLOW"`
		Track  []string `default:"" envconfig:"TRACK"`
	}
	Twitter struct {
		TwitterConsumerKey       string `envconfig:"TWITTER_CONSUMER_KEY" default:""`
		TwitterConsumerSecret    string `envconfig:"TWITTER_CONSUMER_SECRET" default:""`
		TwitterAccessToken       string `envconfig:"TWITTER_ACCESS_TOKEN" default:""`
		TwitterAccessTokenSecret string `envconfig:"TWITTER_ACCESS_TOKEN_SECRET" default:""`
		TwitterAuthCallBack      string `envconfig:"TWITTER_AUTH_CB" default:"http://localhost:8000/api/auth/twitter/callback?provider=twitter"`
		ClientAuthCallBack       string `envconfig:"CLIENT_AUTH_CB_URL" default:"http://localhost:8080?provider=twitter"`
	}
	Mongo struct {
		Host        string        `envconfig:"MONGO_HOST" default:"localhost"`
		Port        string        `envconfig:"MONGO_PORT" default:"27017"`
		Path        string        `envconfig:"MONGO_PATH" default:"elections"`
		User        string        `envconfig:"MONGO_USER" default:""`
		Pass        string        `envconfig:"MONGO_PASS" default:""`
		DialTimeout time.Duration `envconfig:"DIAL_TIMEOUT" default:"5s"`
	}
	Elasticsearch struct {
		Host string `envconfig:"ES_HOST" default:"localhost"`
		Port string `envconfig:"ES_PORT" default:"9200"`
		User string `envconfig:"ES_USER" default:""`
		Pass string `envconfig:"ES_USER" default:""`
	}
	SMTP struct {
		Server   string `envconfig:"SMTP_SERVER" default:"smtp"`
		Port     int    `envconfig:"SMTP_PORT" default:"587"`
		User     string `envconfig:"SMTP_USER" default:"no-reply@mediawatch.io"`
		From     string `envconfig:"SMTP_FROM" default:"no-reply@mediawatch.io"`
		FromName string `envconfig:"SMTP_FROM_NAME" default:"MediaWatch"`
		Pass     string `envconfig:"SMTP_PASS" default:""`
		Reply    string `envconfig:"SMTP_REPLY" default:"press@mediawatch.io"`
	}
	Twillio struct {
		SID   string `envconfig:"TWILIO_SID"`
		Token string `envconfig:"TWILIO_TOKEN"`
	}
	Auth struct {
		Domain         string `envconfig:"DOMAIN_NAME" default:"mediawatch.io"`
		Hash           string `envconfig:"HASH" default:"123"`
		KeyID          string `envconfig:"KEY_ID" default:"0123456789abcdef"`
		PrivateKeyFile string `envconfig:"PRIVATE_KEY_FILE" default:"private.pem"`
		Algorithm      string `envconfig:"ALGORITHM" default:"RS256"`
	}
	Google struct {
		CallBackURL  string `envconfig:"GOOGLE_AUTH_CB_URL" default:"http://localhost:8000/auth/authorize/google/callback"`
		ClientID     string `envconfig:"GOOGLE_AUTH_CLIENT_ID" default:"1"`
		ClientSecret string `envconfig:"GOOGLE_AUTH_CLIENT_SECRET" default:"1"`
	}
	Github struct {
		CallBackURL  string `envconfig:"GITHUB_AUTH_CB_URL" default:"http://localhost:8000/auth/authorize/github/callback"`
		ClientID     string `envconfig:"GITHUB_AUTH_CLIENT_ID" default:""`
		ClientSecret string `envconfig:"GITHUB_AUTH_CLIENT_SECRET" default:""`
	}
}

// New : Create new config struct
func New() *Config {
	return new(Config)
}

// MongoURL : Format MongoURL
func (c *Config) MongoURL() string {
	return fmt.Sprintf("mongodb://%s:%s", c.Mongo.Host, c.Mongo.Port)
}

// ExternalAuths : Enable externa OAuth providers
func (c *Config) ExternalAuths() (map[string]*oauth2.Config, error) {
	var auths = make(map[string]*oauth2.Config)
	google, err := c.EnableGoogle()
	if err == nil {
		auths["google"] = google
	}

	// github, err := c.EnableGithub()
	// if err == nil {
	// 	auths["github"] = github
	// }

	return auths, nil

}

// EnableGoogle : Enables Google OAuth provider
func (c *Config) EnableGoogle() (*oauth2.Config, error) {
	if c.Google.ClientID == "" || c.Google.ClientSecret == "" {
		return nil, fmt.Errorf("Invalid google creds")
	}
	return &oauth2.Config{
		RedirectURL:  c.Google.CallBackURL,
		ClientID:     c.Google.ClientID,
		ClientSecret: c.Google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"openid",
		},
		Endpoint: google.Endpoint,
	}, nil
}

// EnableGithub : Enables Github OAuth provider
// func (c *Config) EnableGithub() (*oauth2.Config, error) {
// 	if c.Github.ClientID == "" || c.Github.ClientSecret == "" {
// 		return nil, fmt.Errorf("Invalid Github creds")
// 	}
// 	return &oauth2.Config{
// 		RedirectURL:  c.Github.CallBackURL,
// 		ClientID:     c.Github.ClientID,
// 		ClientSecret: c.Github.ClientSecret,
// 		Scopes: []string{
// 			"https://www.github.com/auth/userinfo.profile",
// 			"https://www.github.com/auth/userinfo.email",
// 			"openid",
// 		},
// 		Endpoint: github.Endpoint,
// 	}, nil
// }
