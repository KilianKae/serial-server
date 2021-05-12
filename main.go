package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/tarm/serial"
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
	c := &serial.Config{Name: "/dev/cu.usbserial-1460", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err == nil {
		go readFromSerial(s)

		http.HandleFunc("/api/write", writeHandler(s))
	}

	http.HandleFunc("/api/status", statusHandler(c, err))

	log.Printf("Server starting at %s%s", address, port)

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal(err)
	}



	//n, err := s.Write([]byte("test"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//buf := make([]byte, 128)
	//n, err = s.Read(buf)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Print("%q", buf[:n])
}

func readFromSerial(s *serial.Port) {
	for {
		buf := make([]byte, 128)
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s", buf[:n])
	}
}

type StatusResponse struct {
	Name        string 			`json:"name"`
	Baud        int 			`json:"baud"`
	ReadTimeout time.Duration 	`json:"teadTimeout"`
	Size        byte 			`json:"size"`
	Error		string			`json:"error"`
}

func statusHandler(c *serial.Config, err ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := StatusResponse{
			Name:        c.Name,
			Baud:        c.Baud,
			ReadTimeout: c.ReadTimeout,
			Size:        c.Size,
			Error: 		 fmt.Sprint(err...),
		}

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

func writeHandler(s *serial.Port) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp WriteResponse

		_, err := s.Write([]byte("test"))
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
