package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	store, err := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("redis_secret"))
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions("session", store))

	r.GET("/login/facebook/auth", LoginHandler)
	r.GET("/login/facebook/auth/callback", CallbackHandler)

	r.Run()
}
