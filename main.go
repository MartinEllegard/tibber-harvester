package main

import (
	"flag"
	"fmt"
	"sync"
	"tibber-harvester/config"
	"tibber-harvester/handlers"

	"github.com/MartinEllegard/tibber-go"
)

func main() {
	bearerToken := config.Config("BEARER")
	apiUrl := config.Config("API")
	apiToken := flag.String("api-token", "api-token", "Api token used for autherization")
	flag.Parse()

	// Setup DB
	// dbHandler := db.CreateDbHandler()
	// err := dbHandler.SetupHandler()
	// if err != nil {
	// 	panic(err)
	// }

	// Log flags to remove error
	fmt.Println(apiUrl)
	fmt.Println(apiToken)

	tibberClient := tibber.CreateTibberClient(bearerToken, "tibber-harvester")
	viewer, err := tibberClient.GetHomes()
	if err != nil {
		panic(err)
	}

	validHomes := []tibber.Home{}
	for i := range viewer.Viewer.Homes {
		if viewer.Viewer.Homes[i].Features.RealTimeConsumptionEnabled {
			validHomes = append(validHomes, viewer.Viewer.Homes[i])
		}
	}

	var wg sync.WaitGroup

	messageHandler := handlers.CreateMessageHandler(apiUrl)

	messageChannel := make(chan tibber.LiveMeasurement)
	go ProccessMessages(messageHandler, messageChannel)

	for _, home := range validHomes {
		wg.Add(1)
		homeId := home.ID
		go func() {
			tibberClient.StartSubscription(homeId, messageChannel)
			defer wg.Done()
		}()
	}

	wg.Wait()
}

func ProccessMessages(messageHandler handlers.MessageHandler, channel chan tibber.LiveMeasurement) {
	for message := range channel {
		messageHandler.HandlePowerMessage(message)
	}
}
