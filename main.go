package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app"
)

var (
	Version = "Unknown"
	Build   = "Unknown"
)

var (
	//go:embed static templates
	assets embed.FS
)

func main() {
	log.Println("ImageServer", Version, "Build", Build)

	if "Unknown" != Version {
		gin.SetMode(gin.ReleaseMode)

		// Only check updates in release mode
		go app.CheckUpdate(Version)
	}

	app.InitFlag()

	appRouter := app.InitServer(assets)

	go func() {
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", *app.Port),
			Handler: appRouter,
		}
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: %s\n", err)
		}
	}()

	listenOn := app.GetIPAddress()
	if len(listenOn) > 0 {
		log.Println("Listening on these addresses:")
		for _, addr := range listenOn {
			if addr.To4() != nil {
				log.Printf("\thttp://%s:%d\n", addr, *app.Port)
			} else {
				log.Printf("\thttp://[%s]:%d\n", addr, *app.Port)
			}
		}
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Bye~")
}
