package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"

	"github.com/KilianKae/serial-server/internal/serial"

	rice "github.com/GeertJohan/go.rice"
)

const(
	address = "localhost"
	port = ":8080"
)

func main() {
	//Define the rice box with the frontend client static files.
	appBox, err := rice.FindBox("./client/build")
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()


	//Serve static react files
	staticFileServer := http.FileServer(appBox.HTTPBox())
	e.GET("/*", echo.WrapHandler(staticFileServer))

	// Serial

	serialService := serial.NewService()
	serialHandler := serial.NewSerialHandler(serialService)
	err = serialService.FindPort()
	if err == nil {
		go serialService.Read()

		e.POST("/api/write", serialHandler.Write)
	}

	e.GET("/api/ports", serialHandler.GetPorts)

	e.POST("/api/port", serialHandler.SetPort)

	e.GET("/api/status", serialHandler.GetStatus)

	log.Printf("Server starting at %s%s", address, port)

	e.Logger.Fatal(e.Start(":8080"))
}


func serveAppHandler(app *rice.Box) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := app.Open("index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, "index.html", time.Time{}, indexFile)
	}
}
