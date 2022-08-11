package main

import (
	"encoding/json"
	"fmt"
	post2 "github.com/deltamc/otus-social-networks-backend/models/post"
	"github.com/deltamc/otus-social-networks-backend/models/users"
	"github.com/deltamc/otus-social-networks-backend/queue"
	"github.com/deltamc/otus-social-networks-backend/ws"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	server := ws.StartServer(messageHandler)

	rmq := queue.NewRabbitMQ()

	err := rmq.Connect()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer rmq.Close()
	msgs, ch, err := rmq.GetMessages(os.Getenv("RABBITMQ_QUEQE_FEED_WS"))

	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer ch.Close()

	go func() {
		for d := range msgs {
			post := post2.Post{}
			err := json.Unmarshal(d.Body, &post)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			//
			//jsonMessage, _ := json.Marshal(post)

			//ищем друзей
			user := users.User{
				Id: post.UserId,
			}
			flowers, err := user.GetFlowers()

			if err != nil {
				log.Println(err.Error())
				continue
			}

			var tokens []string
			for _, v := range flowers {
				//Для учебного проекта пойдет id в качестве токена
				//В реальном приложении нужно использовать сгенерированный токен
				tokens = append(tokens, strconv.Itoa(int(v.Id)))
			}

			fmt.Println(string(d.Body))
			server.WriteMessage(tokens, d.Body)
		}
	}()

	select {}
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}
