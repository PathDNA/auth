package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/missionMeteora/apiserv"
)

type Client struct {
	c http.Client
}

func newClient() *Client {
	var c Client
	c.c.Jar, _ = cookiejar.New(nil)
	return &c
}

func (c *Client) Signup() *apiserv.JSONResponse {
	parts := strings.Split(*signup, ":")
	if len(parts) != 5 {
		log.Panic("invalid signup string, expected username:password:name:phone:user-type")
	}
	u := newUserWithProfile()
	u.Username, u.Password = parts[0], parts[1]
	p := u.Profile.(*Profile)
	p.Name, p.Phone = parts[2], parts[3]
	if parts[4] == "agency" {
		p.Agency = &Agency{
			Fee: 0.15,
		}
	} else {
		p.Advertiser = &Advertiser{
			AgencyID: parts[4],
		}
	}
	return c.post("signup", u, nil)
}

func (c *Client) Login() *apiserv.JSONResponse {
	parts := strings.Split(*login, ":")
	if len(parts) != 2 {
		log.Panic("invalid login string, expected username:password")
	}
	var loginReq struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}
	loginReq.Username, loginReq.Password = parts[0], parts[1]

	return c.post("login", &loginReq, nil)
}

func (c *Client) Profile() *apiserv.JSONResponse {
	return c.get("profile", newUserWithProfile())
}

func (c *Client) post(ep string, data interface{}, dataValue interface{}) *apiserv.JSONResponse {
	j, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "http://"+addr+"/api/v1/"+ep, bytes.NewReader(j))
	resp, err := c.c.Do(req)
	if err != nil {
		log.Panic(err)
	}
	jresp, err := apiserv.ReadJSONResponse(resp.Body, dataValue)
	if err != nil {
		log.Panic(err)
	}
	return jresp
}

func (c *Client) get(ep string, dataValue interface{}) *apiserv.JSONResponse {
	req, _ := http.NewRequest("GET", "http://"+addr+"/api/v1/"+ep, nil)
	resp, err := c.c.Do(req)
	if err != nil {
		log.Panic(err)
	}
	jresp, err := apiserv.ReadJSONResponse(resp.Body, dataValue)
	if err != nil {
		log.Panic(err)
	}
	return jresp
}
