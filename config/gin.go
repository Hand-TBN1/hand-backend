package config

import "github.com/gin-gonic/gin"

func NewGin() *gin.Engine {
	r := gin.Default()

	return r
}
