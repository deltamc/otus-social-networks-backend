package main

import (
	"encoding/json"
	"github.com/deltamc/otus-social-networks-backend/models/feed"
	post2 "github.com/deltamc/otus-social-networks-backend/models/post"
	"github.com/deltamc/otus-social-networks-backend/queue"
	"github.com/deltamc/otus-social-networks-backend/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	//обновляем кэш при запуские
	err := feed.FeedSetCache()
	if err != nil{
		log.Fatal(err.Error())
		return
	}

	//обновляем кэш если опубликовали новый пост
	rmq := queue.NewRabbitMQ()

	err = rmq.Connect()
	if err != nil{
		log.Fatal(err.Error())
		return
	}
	defer rmq.Close()
	msgs, ch, err := rmq.GetMessages(os.Getenv("RABBITMQ_QUEQE_FEED"))

	if err != nil{
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

			err = feed.FeedSetCache(post.UserId)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	routes.Feed()
	log.Fatal(http.ListenAndServe(":" + os.Getenv("FEED_PORT"), nil))
}