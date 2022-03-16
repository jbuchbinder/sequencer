package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	port = flag.Int("port", 8000, "")

	seq sequencer
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	// Initialize sequencer
	var err error
	seq, err = NewSequencer()
	if err != nil {
		panic(err)
	}
	log.Printf("NODE ID : %d", seq.nodeId)

	gin.SetMode(gin.ReleaseMode)
	//gin.SetMode(gin.DebugMode)

	g := gin.New()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	InitAPI(g)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        g,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	panic(srv.ListenAndServe())
}
