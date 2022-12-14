package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Recover(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil { // 捕获 panic
				debug.PrintStack()
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": fmt.Sprintf("internal error: %v", err),
				})
			}
		}()

		c.Next()
	}
}

func NewGin(l *logrus.Logger) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	e := gin.Default()

	e.Use(cors.New(cors.Config{
		AllowOrigins:    []string{"*"},
		AllowHeaders:    []string{"*"},
		AllowWebSockets: true,
		ExposeHeaders:   nil,
	}))

	e.Use(Recover(l))

	return e
}
