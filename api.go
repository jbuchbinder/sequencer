package main

import (
	"github.com/gin-gonic/gin"
)

var (
	apiMap = map[string]func(*gin.RouterGroup){}
)

func InitAPI(m *gin.Engine) {
	r := m.Group("/sequencer")
	for mapping, fn := range apiMap {
		g := r.Group(mapping)
		fn(g)
	}
}
