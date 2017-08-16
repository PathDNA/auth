package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Path94/auth"
	"github.com/missionMeteora/apiserv"
)

type Profile struct {
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`

	Agency     *Agency     `json:"agency,omitempty"`
	Advertiser *Advertiser `json:"advertiser,omitempty"`
}

type Agency struct {
	Fee float64 `json:"fee,omitempty"`
}

type Advertiser struct {
	AgencyID string `json:"agencyID,omitempty"`
}

type Server struct {
	s      *apiserv.Server
	a      *auth.Auth
	tokens sync.Map
	dbPath string
}

func newServer() *Server {
	var (
		s   Server
		err error
	)
	if s.dbPath, err = ioutil.TempDir("", "auth_webapp_demo"); err != nil {
		log.Panic(err)
	}
	if s.a, err = auth.New(s.dbPath); err != nil {
		log.Panic(err)
	}

	s.a.ProfileFn = func() interface{} { return &Profile{} } // this allows proper unmarshaling of users

	s.s = apiserv.New(apiserv.SetNoCatchPanics(true))

	if *debug {
		s.s.Use(apiserv.LogRequests(true))
	}

	s.s.POST("/api/v1/signup", s.signup)
	s.s.POST("/api/v1/login", s.login)
	g := s.s.Group("/api/v1", s.checkUser)
	g.GET("/profile", s.getProfile)
	return &s
}

func (s *Server) signup(ctx *apiserv.Context) apiserv.Response {
	u := newUserWithProfile()
	if err := ctx.BindJSON(u); err != nil {
		return apiserv.NewJSONErrorResponse(400, err)
	}

	u.Status = auth.StatusActive

	id, err := s.a.CreateUser(u, u.Password)
	if err != nil {
		return apiserv.NewJSONErrorResponse(400, err)
	}

	return apiserv.NewJSONResponse(id)
}

var respBadUserPass = apiserv.NewJSONErrorResponse(http.StatusBadRequest, "invalid username and/or password")

func (s *Server) login(ctx *apiserv.Context) apiserv.Response {
	var loginReq struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}
	if err := ctx.BindJSON(&loginReq); err != nil {
		return apiserv.NewJSONErrorResponse(400, err)
	}
	u, err := s.a.GetUserByName(loginReq.Username)
	if err != nil {
		return respBadUserPass
	}
	if !u.PasswordsMatch(loginReq.Password) {
		return respBadUserPass
	}
	// TODO do proper tokens with MAC
	token := auth.RandomToken(16, true)
	s.tokens.Store(token, u.ID)
	ctx.SetCookie("token", token, "", false, time.Hour)
	return apiserv.NewJSONResponse(u.ID)
}

func (s *Server) checkUser(ctx *apiserv.Context) apiserv.Response {
	tok, ok := ctx.GetCookie("token")
	if !ok {
		return apiserv.RespForbidden
	}

	uid, ok := s.tokens.Load(tok)
	if !ok {
		return apiserv.RespForbidden
	}

	u, err := s.a.GetUserByID(uid.(string))
	if err != nil {
		return apiserv.RespForbidden
	}
	ctx.Set("user", u)
	return nil
}

func (s *Server) getProfile(ctx *apiserv.Context) apiserv.Response {
	return apiserv.NewJSONResponse(ctx.Get("user"))
}

// this is used for BindJSON mostly, the auth lib handles the rest
func newUserWithProfile() *auth.User {
	return &auth.User{
		Profile: &Profile{},
	}
}
