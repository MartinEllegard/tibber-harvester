package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MartinEllegard/tibber-go"
)

type MessageHandler struct {
	endpoint    string
	contentType string
}

type PowerMetricRequest struct {
	Timestamp                      time.Time `json:"ts"`
	Currency                       string    `json:"currency"`
	HomeId                         string    `json:"homeId"`
	AccumulatedConsumption         float64   `json:"accumulatedConsumption"`
	AveragePower                   float64   `json:"averagePower"`
	MaxPower                       float64   `json:"maxPower"`
	LastMeterConsumption           float64   `json:"lastMeterConsumption"`
	MinPower                       float64   `json:"minPower"`
	AccumulatedProduction          float64   `json:"accumulatedProduction"`
	AccumulatedCost                float64   `json:"accumulatedCost"`
	AccumulatedConsumptionLastHour float64   `json:"accumulatedConsumptionLastHour"`
	AccumulatedProductionLastHour  float64   `json:"accumulatedProductionLastHour"`
	Power                          float64   `json:"power"`
	LastMeterProduction            float64   `json:"lastMeterProduction"`
}

func CreateMessageHandler(url string) MessageHandler {
	return MessageHandler{endpoint: url, contentType: "application/json"}
}

func (h *MessageHandler) HandlePowerMessage(message tibber.LiveMeasurement) {
	metric := PowerMetricRequest{
		Timestamp:                      message.Timestamp,
		Currency:                       message.Currency,
		LastMeterProduction:            message.LastMeterProduction,
		MinPower:                       message.MinPower,
		AveragePower:                   message.AveragePower,
		MaxPower:                       message.MaxPower,
		LastMeterConsumption:           message.LastMeterConsumption,
		AccumulatedConsumption:         message.AccumulatedConsumption,
		AccumulatedProduction:          message.AccumulatedProduction,
		AccumulatedCost:                message.AccumulatedCost,
		AccumulatedConsumptionLastHour: message.AccumulatedConsumptionLastHour,
		AccumulatedProductionLastHour:  message.AccumulatedProductionLastHour,
		Power:                          message.Power,
		HomeId:                         message.HomeId,
	}

	postBody, err := json.Marshal(metric)
	if err != nil {
		log.Println("Failed to convert to json")
		return
	}
	log.Println(string(postBody))
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(h.endpoint, h.contentType, responseBody)
	if err != nil {
		log.Println("Error occured dring post request", err)
	}
	defer resp.Body.Close()
}
