package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	apiMap["/api"] = func(r *gin.RouterGroup) {
		r.GET("/seq", apiSequence)
	}
}

func apiSequence(m *gin.Context) {
	id, err := seq.NextId()
	if err != nil {
		m.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	m.JSON(http.StatusOK, id)
}
