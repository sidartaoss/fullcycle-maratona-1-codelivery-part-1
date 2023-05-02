package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sidartaoss/maratona-fullcycle-1/simulator/entity"
	"github.com/sidartaoss/maratona-fullcycle-1/simulator/queue"
	"github.com/streadway/amqp"
)

var active []string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	in := make(chan []byte)

	ch := queue.Connect()
	queue.StartConsuming(in, ch)

	for msg := range in {
		var order entity.Order

		err := json.Unmarshal(msg, &order)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		log.Println("New order Received: ", order.Uuid)

		start(order, ch)
	}
}

func start(order entity.Order, ch *amqp.Channel) {
	if !stringInSlice(order.Uuid, active) {
		active = append(active, order.Uuid)
		go SimulatorWorker(order, ch)
	} else {
		log.Println("Order", order.Uuid, "was already completed or is on going...")
	}
}

func SimulatorWorker(order entity.Order, ch *amqp.Channel) {

	f, err := os.Open("destinations/" + order.Destination + ".txt")
	if err != nil {
		panic(err.Error())
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")
		json := destinationToJson(order, data[0], data[1])
		time.Sleep(1 * time.Second)
		queue.Notify(string(json), ch)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	json := destinationToJson(order, "0", "0")
	queue.Notify(string(json), ch)
}

func destinationToJson(order entity.Order, lat string, lng string) []byte {
	dest := entity.Destination{
		Order: order.Uuid,
		Lat:   lat,
		Lng:   lng,
	}
	json, _ := json.Marshal(dest)
	return json
}

func stringInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
