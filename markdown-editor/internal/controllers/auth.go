package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (co *Controller) HandleLogin(c *gin.Context) {
	if c.Request.Method == "GET" {
		msg := ""
		if c.Query("login") == "failed" {
			msg = "Invalid"
		}
		co.Templates.ExecuteTemplate(c.Writer, "login.html", msg)
		return
	}
	if c.Request.Method == "POST" && co.Auth.Authenticate(c.PostForm("username"), c.PostForm("password")) {
		c.SetCookie("session_token", "logged-in", 3600*24, "/", "", false, true)
		c.Redirect(http.StatusFound, "/")
	} else {
		c.Redirect(http.StatusFound, "/login?login=failed")
	}
}

func (co *Controller) HandleLogout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func (co *Controller) Protect() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_token")
		if err != nil || cookie != "logged-in" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
