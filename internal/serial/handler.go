package serial

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPHandler interface {
	GetStatus(c echo.Context) error
	GetPorts(c echo.Context) error
	SetPort(c echo.Context) error
	Write(c echo.Context) error
}

type serialHandler struct {
	Service Service
}

func NewSerialHandler(service Service) HTTPHandler {
	return &serialHandler{Service: service}
}

func (s *serialHandler) GetPorts(c echo.Context) error {
	ports := s.Service.Ports()

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Allow-Methods", "GET")
	c.Response().Header().Set("Access-Control-Allow-Headers", "x-custom-header")

	return c.JSON(http.StatusOK, ports)
}

type setPortRequest struct {
	port int    `json:"port"`
}

func (s *serialHandler) SetPort(c echo.Context) error {
	body, err := c.Request().GetBody()
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(body)
	var portReq setPortRequest
	err = decoder.Decode(&portReq)
	if err != nil {
		return err
	}
	//TODO on error
	s.Service.SetPort(portReq.port)
	return nil
}

func (s *serialHandler) GetStatus(c echo.Context) error {
		status := s.Service.GetStatus()

		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET")
		c.Response().Header().Set("Access-Control-Allow-Headers", "x-custom-header")

		return c.JSON(http.StatusOK, status)
}

type WriteResponse struct {
	Success bool    `json:"success"`
	Message  string `json:"message"`
}

//TODO
func (s *serialHandler) Write(c echo.Context) error {
	var resp WriteResponse

	log.Printf("Writing")
	err := s.Service.Write("test")
	if err != nil {
		resp = WriteResponse{Success: true, Message: fmt.Sprintf("Writing %s", err.Error())}
	} else {
		resp = WriteResponse{Success: true, Message: fmt.Sprintf("Writing %s", "test")}
	}

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, resp)
}