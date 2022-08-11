package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

type Server struct {
	clients       map[string]*websocket.Conn
	handleMessage func(message []byte) // хандлер новых сообщений
}

func StartServer(handleMessage func(message []byte)) *Server {
	server := Server{
		make(map[string]*websocket.Conn),
		handleMessage,
	}

	http.HandleFunc("/", server.echo)
	go http.ListenAndServe(":8084", nil) // Уводим http сервер в горутину

	return &server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)

	token := r.FormValue("token")
	defer func() {
		connection.Close()
		fmt.Println("close connect " + token)
	}()
	server.clients[token] = connection  // Сохраняем соединение, используя его как ключ
	defer delete(server.clients, token) // Удаляем соединение

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
		}

		go server.handleMessage(message)
	}
}

func (server *Server) GetUserTokens() (list []string) {
	for token, _ := range server.clients {
		list = append(list, token)
	}

	return
}

func (server *Server) WriteMessage(userTokens []string, message []byte) {
	for _, token := range userTokens {
		if v, ok := server.clients[token]; ok {

			v.WriteMessage(websocket.TextMessage, message)
		}

	}
}
