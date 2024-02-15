package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"github.com/MartinEllegard/tibber-go"
)

func main() {
	bearerToken := flag.String("bearer", "5K4MVS-OjfWhK_4yrjOlFe1F6kJXPVf7eQYggo8ebAE", "Tibber Api token")
	apiUrl := flag.String("api-url", "https://localhost:8080/api/powerusage", "Api url to post this data too")
	apiToken := flag.String("api-token", "api-token", "Api token used for autherization")
	flag.Parse()

	//Log flags to remove error
	fmt.Println(apiUrl)
	fmt.Println(apiToken)

	tibberClient := tibber.CreateTibberClient(*bearerToken, "tibber-harvester")
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

	for _, home := range validHomes {
		wg.Add(1)
		messageChannel := make(chan tibber.LiveMeasurement)
		homeId := home.ID
		go ProccessMessages(messageChannel)
		go func() {
			tibberClient.StartSubscription(homeId, messageChannel)
			defer wg.Done()
		}()
	}

	wg.Wait()
}

func ProccessMessages(channel chan tibber.LiveMeasurement) {
	for message := range channel {
		jsonMessage, _ := json.Marshal(message)
		fmt.Println(string(jsonMessage))
	}
}
