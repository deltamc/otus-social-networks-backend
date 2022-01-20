package queue

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"testing"
)
type Message struct {
	Name string
	Body string
}

func TestRabbitMQ(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil{
		fmt.Println(err.Error())
	}
fmt.Println(os.Getenv("RABBITMQ_USERNAME"))
	rmq := NewRabbitMQ()
	err = rmq.Connect()
	if err != nil{
		fmt.Println("conn " + err.Error())
	}
	defer rmq.Close()

	err = rmq.Push("test", Message{
		"Вася",
		"Привет, как дела?",
	})

	if err != nil{
		fmt.Println(err.Error())
	}

}