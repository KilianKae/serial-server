package main

import (
	"encoding/json"
	"fmt"
	"github.com/KilianKae/serial-server/internal/serial"
	"log"
	"net/http"
	"time"

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

	//Define ping endpoint that responds with pong.
	http.HandleFunc("/api/ping", pingHandler())


	//Serve static files
	http.Handle("/static/", http.FileServer(appBox.HTTPBox()))
	//Serve SPA (Single Page Application)
	http.HandleFunc("/", serveAppHandler(appBox))

	// Serial
	serialService := serial.NewService()
	err = serialService.FindPort()
	if err == nil {
		go serialService.Read()

		http.HandleFunc("/api/write", writeHandler(serialService))
	}

	http.HandleFunc("/api/ports", portsHandler(serialService))

	http.HandleFunc("/api/status", statusHandler(serialService))

	log.Printf("Server starting at %s%s", address, port)

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal(err)
	}
}

func portsHandler(s serial.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := s.Ports()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "x-custom-header")
		json.NewEncoder(w).Encode(resp)
	}
}

func statusHandler(s serial.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := s.Status()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "x-custom-header")
		json.NewEncoder(w).Encode(resp)
	}
}

type WriteResponse struct {
	Success bool    `json:"success"`
	Message  string `json:"message"`
}

func writeHandler(s serial.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp WriteResponse

		err := s.Write("test")
		log.Printf("Writing %s", "test")
		if err != nil {
			resp = WriteResponse{Success: true, Message: fmt.Sprintf("Writing %s", err.Error())}
		} else {
			resp = WriteResponse{Success: true, Message: fmt.Sprintf("Writing %s", "test")}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	}
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
