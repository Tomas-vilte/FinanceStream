package main

import (
	"github.com/Tomas-vilte/FinanceStream/internal/config"
	"github.com/Tomas-vilte/FinanceStream/internal/kafka"
	"github.com/Tomas-vilte/FinanceStream/internal/realtime"
	"log"
	"time"
)

func main() {
	appConfig := config.RealTimeConfig{
		BinanceChannels: []config.ChannelConfig{
			{
				Symbol:     "btcusdt",
				Channel:    "bookTicker",
				KafkaTopic: "binanceBookTicker",
			},
			{
				Symbol:     "btcusdt",
				Channel:    "ticker",
				KafkaTopic: "binanceTrade",
			},
		},
		KafkaBroker: "localhost:9092",
	}

	kafkaConn, err := kafka.NewKafkaProducer(appConfig.KafkaBroker)
	if err != nil {
		log.Fatal("Error al crear la conexión a Kafka:", err)
		return
	}
	defer kafkaConn.Close()

	// Procesar cada configuración de canal
	for _, channelConfig := range appConfig.BinanceChannels {
		channelWS, err := realtime.NewBinanceWebSocket([]config.ChannelConfig{channelConfig})
		if err != nil {
			log.Fatalf("Error al crear la conexión WebSocket para %s: %v\n", channelConfig.Channel, err)
		}
		defer channelWS.Close()

		// Suscribirse y publicar en Kafka
		subscribeAndPublish(channelWS, kafkaConn, channelConfig.KafkaTopic)
	}

	time.Sleep(1 * time.Minute)
}
