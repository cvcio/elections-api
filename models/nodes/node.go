package nodes

import (
	proto "github.com/cvcio/elections-api/pkg/proto"
)

// UserObj ...
type UserObj struct {
	CrawledAt       string              `json:"crawled_at"`
	Id              int64               `json:"id"`
	IdStr           string              `json:"id_str"`
	CreatedAt       string              `json:"created_at"`
	Verified        bool                `json:"verified"`
	ScreenName      string              `json:"screen_name"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	FollowersCount  int                 `json:"followers_count"`
	FriendsCount    int                 `json:"friends_count"`
	ListedCount     int64               `json:"listed_count"`
	StatusesCount   int64               `json:"statuses_count"`
	FavouritesCount int                 `json:"favourites_count"`
	ProfileImage    string              `json:"profile_image_url_https"`
	BannerImage     string              `json:"profile_banner_url"`
	UserClass       string              `json:"user_class"`
	UserClassScore  float64             `json:"user_class_score"`
	QuotedStatus    *UserObj            `json:"quoted_status,omitempty"`
	RetweetedStatus *UserObj            `json:"retweeted_status,omitempty"`
	Metrics         *proto.UserFeatures `json:"metrics"`
}

// ESUserObj ...
type ESUserObj struct {
	CrawledAt       string  `json:"crawled_at"`
	Id              int64   `json:"id"`
	IdStr           string  `json:"id_str"`
	CreatedAt       string  `json:"created_at"`
	Verified        bool    `json:"verified"`
	ScreenName      string  `json:"screen_name"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	FollowersCount  int     `json:"followers_count"`
	FriendsCount    int     `json:"friends_count"`
	ListedCount     int64   `json:"listed_count"`
	StatusesCount   int64   `json:"statuses_count"`
	FavouritesCount int     `json:"favourites_count"`
	ProfileImage    string  `json:"profile_image_url_https"`
	BannerImage     string  `json:"profile_banner_url"`
	UserClass       string  `json:"user_class"`
	UserClassScore  float64 `json:"user_class_score"`
}
