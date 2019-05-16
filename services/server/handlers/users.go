package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/cvcio/elections-api/models/user"
	"github.com/cvcio/elections-api/pkg/auth"
	"github.com/cvcio/elections-api/pkg/config"
	"github.com/cvcio/elections-api/pkg/db"
	"github.com/cvcio/elections-api/pkg/mailer"
	"github.com/gin-gonic/gin"
	gothic "github.com/markbates/goth/gothic"
	"github.com/plagiari-sm/mediawatch/pkg/es"
	log "github.com/sirupsen/logrus"
)

// Users Controller
type Users struct {
	cfg           *config.Config
	db            *db.DB
	es            *es.ES
	authenticator *auth.Authenticator
	mail          *mailer.Mailer
}

// OAuthTwitter Atuhorizes Twitter using gothic
func (ctrl *Users) OAuthTwitter(c *gin.Context) {
	c.Request = c.Request.WithContext(
		context.WithValue(
			c.Request.Context(),
			"provider",
			c.Param("provider"),
		),
	)
	if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		log.Debugf("%v", gothUser)
	} else {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

// OAuthTwitterCB Twitter Authorization Callback using gothic
func (ctrl *Users) OAuthTwitterCB(c *gin.Context) {
	c.Request = c.Request.WithContext(
		context.WithValue(
			c.Request.Context(),
			"provider",
			c.Param("provider"),
		),
	)

	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		ResponseError(c, 401, err.Error())
		return
	}

	exists, _ := user.ByScreenNanme(ctrl.db, gothUser.NickName)
	if exists != nil {
		log.Infof("EXISTS ------ %s", exists.ScreenName)
		// update tokens
		exists.TwitterAccessToken = &gothUser.AccessToken
		exists.TwitterAccessTokenSecret = &gothUser.AccessTokenSecret
		exists.ProfileImageURL = &gothUser.AvatarURL
		// no need to check for errors
		user.Update(ctrl.db, exists.ID.Hex(), exists)

		if exists.Status == "" {
			c.Redirect(
				http.StatusMovedPermanently,
				ctrl.cfg.Twitter.ClientAuthCallBack+"&method=create&twitterId="+
					exists.UserID+"&screenName="+exists.ScreenName+"&idStr="+exists.IDStr)
			return
		}
		if exists.Status == "create" {
			c.Redirect(
				http.StatusMovedPermanently,
				ctrl.cfg.Twitter.ClientAuthCallBack+"&method=verify&idStr="+exists.IDStr+"&mobile="+exists.Mobile)
			return
		}

		// TODO: redirect to home...
		c.Redirect(
			http.StatusMovedPermanently,
			ctrl.cfg.Twitter.ClientAuthCallBack+"&method=login&twitterId="+
				exists.UserID+"&screenName="+exists.ScreenName+"&idStr="+exists.IDStr+"&twitterAccessToken="+*exists.TwitterAccessToken)
		return
	}

	var u user.User

	u.UserID = gothUser.UserID
	u.ScreenName = gothUser.NickName
	u.TwitterAccessToken = &gothUser.AccessToken
	u.TwitterAccessTokenSecret = &gothUser.AccessTokenSecret
	u.ProfileImageURL = &gothUser.AvatarURL
	u.Pin = user.RandStringBytes(8)

	created, err := user.Create(ctrl.db, &u)
	if err != nil {
		ResponseError(c, 401, err.Error())
		return
	}

	c.Redirect(
		http.StatusMovedPermanently,
		ctrl.cfg.Twitter.ClientAuthCallBack+"&method=create&twitterId="+
			created.UserID+"&screenName="+created.ScreenName+"&idStr="+created.IDStr)
}

// Update User Handler
func (ctrl *Users) Update(c *gin.Context) {
	var u *user.User
	if err := c.Bind(&u); err != nil {
		ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := user.Update(ctrl.db, u.IDStr, u)
	if err != nil {
		ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ret, _ := user.Get(ctrl.db, u.IDStr)
	if ret.Status == "create" {
		go user.SendOTP(ctrl.cfg.Twillio.SID, ctrl.cfg.Twillio.Token, ret.Mobile, ret.Pin)
	}
	ResponseJSON(c, &ret)
}

// Verify OTP
func (ctrl *Users) Verify(c *gin.Context) {
	type Req struct {
		Pin   string `json:"pin"`
		IDStr string `json:"idStr"`
	}
	var req *Req
	if err := c.Bind(&req); err != nil {
		ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	u, err := user.Get(ctrl.db, req.IDStr)
	if err != nil {
		ResponseError(c, 404, "Account not found")
		return
	}
	if u.Pin != req.Pin {
		ResponseError(c, 403, "Authorization code doesn't match")
		return
	}

	u.Pin = ""
	u.Status = "verify"
	_, err = user.Update(ctrl.db, u.IDStr, u)
	if err != nil {
		ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseJSON(c, &u)
}

// SendPin Send New OTP
func (ctrl *Users) SendPin(c *gin.Context) {
	type Req struct {
		IDStr string `json:"idStr"`
	}
	var req *Req
	if err := c.Bind(&req); err != nil {
		ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	u, err := user.Get(ctrl.db, req.IDStr)
	if err != nil {
		ResponseError(c, 404, "Account not found")
		return
	}

	u.Pin = user.RandStringBytes(8)
	_, err = user.Update(ctrl.db, u.IDStr, u)
	if err != nil {
		ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	go user.SendOTP(ctrl.cfg.Twillio.SID, ctrl.cfg.Twillio.Token, u.Mobile, u.Pin)

	ResponseJSON(c, &u)
}

// Token Authorize App
func (ctrl *Users) Token(c *gin.Context) {
	type Req struct {
		IDStr              string `json:"idStr"`
		TwitterAccessToken string `json:"twitterAccessToken"`
	}
	var req *Req
	if err := c.Bind(&req); err != nil {
		ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	u, err := user.Get(ctrl.db, req.IDStr)
	if err != nil {
		ResponseError(c, 404, "Account not found")
		return
	}

	if *u.TwitterAccessToken != req.TwitterAccessToken {
		ResponseError(c, http.StatusUnauthorized, "Access Token Error")
		return
	}

	tkn, err := user.Authenticate(ctrl.authenticator, time.Now(), u)
	if err != nil {
		log.Debug(err)
		ResponseError(c, http.StatusUnauthorized, "Access Token Error")
		return
	}

	c.Header("Authorization", tkn.Token)

	ResponseJSON(c, &u)
}
