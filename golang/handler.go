package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
)

func LoginHandler(c *gin.Context) {
	fb := NewConfig()

	state := uuid.NewV4().String()
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()

	url := fb.AuthCodeURL(state, oauth2.AccessTypeOnline)
	c.Redirect(http.StatusMovedPermanently, url)
	return
}

func CallbackHandler(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	state := c.DefaultQuery("state", "")

	session := sessions.Default(c)
	v := session.Get("state")

	err := validateFacebookCode(code, state, v)
	if err != nil {
		c.String(http.StatusBadRequest, "%s", err)
		return
	}

	id, err := Login(&FacebookImpl{}, code)
	if err != nil {
		c.String(http.StatusBadRequest, "%s", err)
		return
	}

	c.String(http.StatusOK, "%s", id)
	return
}