package db

import (
	"context"
	"log"
	"tibber-harvester/config"
	"time"

	"github.com/MartinEllegard/tibber-go"
	qdb "github.com/questdb/go-questdb-client/v3"
)

func CreateDbHandler() *QuestDbHandler {
	return &QuestDbHandler{
		PowerChannel: make(chan tibber.LiveMeasurement),
		PriceChannel: make(chan tibber.PriceInfo),
	}
}

type QuestDbHandler struct {
	PowerChannel chan tibber.LiveMeasurement
	PriceChannel chan tibber.PriceInfo
}

func (h *QuestDbHandler) SetupHandler() error {
	connectionString := config.Config("QDBADDR")

	sender, err := qdb.NewLineSender(context.Background(), qdb.WithAddress(connectionString), qdb.WithTcp())
	if err != nil {
		return err
	}

	go func() {
		qdbPowerHandler(h.PowerChannel, sender)
	}()

	return nil
}

func qdbPowerHandler(rx <-chan tibber.LiveMeasurement, lineSender qdb.LineSender) {
	defer lineSender.Close(context.Background())
	log.Println("Started Influx proto writer")

	running := true

	go func() {
		for running {
			time.Sleep(time.Millisecond * 1000)
			lineSender.Flush(context.Background())
		}
	}()

	for message := range rx {
		if (message != tibber.LiveMeasurement{}) {
			log.Println("Tibber message received")
			err := lineSender.
				Table("sm_powermeter").
				Symbol("home_id", message.HomeId).
				TimestampColumn("ts", message.Timestamp).
				Float64Column("power", message.Power).
				Float64Column("min_power", message.MinPower).
				Float64Column("max_power", message.MaxPower).
				Float64Column("average_power", message.AveragePower).
				Float64Column("last_meter_consumption", message.LastMeterConsumption).
				Float64Column("last_meter_production", message.LastMeterProduction).
				Float64Column("accumulated_consumption", message.AccumulatedConsumption).
				Float64Column("accumulated_production", message.AccumulatedProduction).
				Float64Column("accumulated_cost", message.AccumulatedCost).
				Float64Column("accumulated_consumption_last_hour", message.AccumulatedConsumptionLastHour).
				Float64Column("accumulated_production_last_hour", message.AccumulatedProductionLastHour).
				StringColumn("currency", message.Currency).
				At(context.Background(), message.Timestamp)
			if err != nil {
				log.Println("Failed to write message to questdb")
			}
		}
	}
	log.Println("Stopped Influx proto writer")
}
