package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/tarm/serial"
	"go.bug.st/serial/enumerator"
)

const(
	address = "localhost"
	port = ":8080"
)

func main() {
	fmt.Println("Looking for ports")
	ports, _ := enumerator.GetDetailedPortsList()
	fmt.Println("Ports: ",ports)

	var c *serial.Config
	var scanner *bufio.Scanner
	var serialPort *serial.Port
	for _, port := range ports {
		fmt.Println("Checking port: ", port)
		if port.IsUSB {
			c = &serial.Config{Name: port.Name, Baud: 115200}

			var err error
			serialPort, err = serial.OpenPort(c)
			if err != nil {
				continue
			}

			scanner = bufio.NewScanner(serialPort)
			scanner.Scan()
			if scanner.Text() == "Setup" {
				fmt.Println("Found port: ", port)
				break
			}
		}
	}

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
	if err == nil {
		go readFromSerial(scanner)

		http.HandleFunc("/api/write", writeHandler(serialPort))
	}

	http.HandleFunc("/api/status", statusHandler(c, err))

	log.Printf("Server starting at %s%s", address, port)

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal(err)
	}
}

func readFromSerial(s *bufio.Scanner) {
	for s.Scan() {
		log.Printf(s.Text())
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
