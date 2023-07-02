package main

import (
	"log"
	"os"
	"os/user"
	"sync"

	"github.com/allusernamestakenexceptthis/langchat/routes/home"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{}

func main() {
	e := echo.New()

	e.Static("/static", "assets")

	chatRoom := CreateChatRoom("main")
	e.GET("/cr/main", chatRoom.webSocket)

	//define routes
	e.GET("/", home.Home)
	//e.Static("/", "../front")
	log.SetOutput(os.Stdout)
	log.Printf("hello world")

	e.Logger.Fatal(e.Start(":1323"))
}

type chatRoom struct {
	name  string
	users UserConnections
	sync.RWMutex
}

func CreateChatRoom(name string) *chatRoom {
	return &chatRoom{
		name:  name,
		users: make(UserConnections),
	}
}

type UserConnection struct {
	conn     *websocket.Conn
	user     *user.User
	chatRoom *chatRoom
	comChan  chan []byte
}

func (uc *UserConnection) ReadMessages() {
	// read messages from user
	defer uc.chatRoom.removeUser(uc)

	for {
		msgType, msg, err := uc.conn.ReadMessage()

		if err != nil {
			log.Printf("error reading message from user: %s", err)
			break
		}

		//send message to message queue
		//send message to other users
		//uc.chatRoom.SendMessage(string(msg))
		log.Printf("msgType: %d, msg: %s", msgType, msg)

		for user := range uc.chatRoom.users {
			//if user != uc {
			user.comChan <- msg
			//}
		}
	}
}

func (uc *UserConnection) SendMessages() {
	defer uc.chatRoom.removeUser(uc)
	log.Printf("send messages")
	for {
		if msg, success := <-uc.comChan; success {
			//send message to user
			err := uc.conn.WriteMessage(websocket.TextMessage, msg)
			log.Printf("msg: %s", msg)
			if err != nil {
				log.Printf("error sending message to user: %s", err)
				return
			}
		}
	}
}

type UserConnections map[*UserConnection]bool

func NewUserConnection(u *user.User, ws *websocket.Conn, chatRoom *chatRoom) *UserConnection {
	return &UserConnection{
		conn:     ws,
		user:     u,
		chatRoom: chatRoom,
		comChan:  make(chan []byte),
	}
}

func (c *chatRoom) addUser(uc *UserConnection) {
	c.Lock()
	defer c.Unlock()

	c.users[uc] = true
}

func (c *chatRoom) removeUser(uc *UserConnection) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.users[uc]; !ok {
		return
	}

	uc.conn.Close()
	delete(c.users, uc)
}

func (c *chatRoom) webSocket(ctx echo.Context) error {
	log.SetOutput(os.Stdout)
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)

	if err != nil {
		log.Printf("error upgrading connection to websocket: %s", err)
		return err
	}

	//get message from message queue
	//send message to user

	// Write
	/*
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			ctx.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			ctx.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)*/
	user := new(user.User)
	UserConn := NewUserConnection(user, ws, c)
	c.addUser(UserConn)

	log.Printf("user connected: %s", user.Username)

	go UserConn.ReadMessages()
	go UserConn.SendMessages()

	return nil
}
