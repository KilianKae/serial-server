package serial

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
	"go.bug.st/serial/enumerator"
)

type Service interface {
	FindPort() 		error
	Read()
	Write(string) 	error
	GetStatus()     StatusResponse
	Ports() 		[]*enumerator.PortDetails
	SetPort(int)    error
}

type utility struct {
	config 		*serial.Config
	serialPort 	*serial.Port
	scanner		*bufio.Scanner
	error 		error
}

func NewService() Service {
	return &utility{config: &serial.Config{}}
}

func (u *utility) Ports() (portList []*enumerator.PortDetails) {
	portList = make([]*enumerator.PortDetails, 0)

	ports, _ := enumerator.GetDetailedPortsList()
	for _, port := range ports {
		if port.IsUSB {
			portList = append(portList, port)
		}
	}

	return portList
}

func (u *utility) SetPort(n int) error {
	ports, _ := enumerator.GetDetailedPortsList()
	port := ports[n]
	u.config = &serial.Config{Name: port.Name, Baud: 115200}

	var err error
	u.serialPort, err = serial.OpenPort(u.config)
	if err != nil {
		return err
	}

	u.scanner = bufio.NewScanner(u.serialPort)
	u.scanner.Scan()
	if u.scanner.Text() == "Setup" {
		fmt.Println("Found port: ", port)
		return nil
	}
	return errors.New("port open error")
}

func (u *utility) FindPort() error {
	ports, _ := enumerator.GetDetailedPortsList()
	for _, port := range ports {
		fmt.Println("Checking port: ", port)
		if port.IsUSB {
			u.config = &serial.Config{Name: port.Name, Baud: 115200}

			var err error
			u.serialPort, err = serial.OpenPort(u.config)
			if err != nil {
				continue
			}

			u.scanner = bufio.NewScanner(u.serialPort)
			u.scanner.Scan()
			if u.scanner.Text() == "Setup" {
				fmt.Println("Found port: ", port)
				return nil
			}
		}
	}
	u.error = errors.New("no port found")
	return errors.New("no port found")
}

func (u *utility) Read() {
	for u.scanner.Scan() {
		log.Printf(u.scanner.Text())
	}
}

func (u *utility) Write(s string) error {
	_, err := u.serialPort.Write([]byte(s))
	return err
}

type StatusResponse struct {
	Name        string 			`json:"name"`
	Baud        int 			`json:"baud"`
	ReadTimeout time.Duration 	`json:"readTimeout"`
	Size        byte 			`json:"size"`
	Error		string			`json:"error"`
}

func (u *utility) GetStatus() StatusResponse {
	status := StatusResponse{
		Name:        u.config.Name,
		Baud:        u.config.Baud,
		ReadTimeout: u.config.ReadTimeout,
		Size:        u.config.Size,
	}
	if u.error != nil {
		status.Error = u.error.Error()
	}
	return status
}
