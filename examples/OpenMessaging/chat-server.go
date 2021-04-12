package main

import (
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-purecloud"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type ChatServer struct {
	router      *mux.Router
	chats       map[uuid.UUID]*Chat
	messages    map[string]*ChatMessage
	register    chan *Chat
	unregister  chan *Chat
	sendCX      chan *ChatMessage
	Logger      *logger.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewChatServer(router *mux.Router, log *logger.Logger) *ChatServer {
	return &ChatServer{
		router:      router,
		chats:       map[uuid.UUID]*Chat{},
		messages:    map[string]*ChatMessage{},
		register:    make(chan *Chat),
		unregister:  make(chan *Chat),
		sendCX:      make(chan *ChatMessage),
		Logger:      log.Child("chatserver", "chatserver"),
	}
}

func (server *ChatServer) CreateChat(userID string) *Chat {
	log := server.Logger.Child(nil, "create_chat")

	log.Infof("Creating a new chat with user %s", userID)
	chat := NewChat(server, userID)
	log = log.Record("chat", chat.ID)
	log.Tracef("registering chat")
	server.register <- chat
	log.Debugf("chat registered")
	return chat
}

func (server *ChatServer) FindChatByID(id uuid.UUID) (*Chat, error) {
	log := server.Logger.Child(nil, "findchatbyid", "chat", id)

	log.Debugf("Looking for Chat %s", id)
	if chat, found := server.chats[id]; found {
		return chat, nil
	}
	return nil, errors.NotFound.With("chat", id.String()).WithStack()
}

func (server *ChatServer) FindChatByUserID(userID string) (*Chat, error) {
	log := server.Logger.Child(nil, "findchatbyuserid", "user", userID)

	log.Debugf("Looking for Chat with User %s", userID)
	for _, chat := range server.chats {
		if userID == chat.UserID {
			return chat, nil
		}
	}
	return nil, errors.NotFound.With("user", userID).WithStack()
}

func (server *ChatServer) FindChatByMessageID(id string) (*Chat, error) {
	log := server.Logger.Child(nil, "findchatbyid", "message", id)

	log.Debugf("Looking for Chat with Message %s", id)
	if message, found := server.messages[id]; found {
		return message.Chat, nil
	}
	return nil, errors.NotFound.With("message", id).WithStack()
}

func (server *ChatServer) Start(config *Config) {
	for {
		select {
		case chat := <-server.register:
			server.Logger.Record("chat", chat).Infof("Registering new chat %s", chat)
			server.chats[chat.ID] = chat
		case chat := <- server.unregister:
			server.Logger.Record("chat", chat).Infof("Unregistering chat %s", chat)
			delete(server.chats, chat.ID)
		case message := <-server.sendCX:
			server.messages[message.ID] = message

			// Here we need to send the message to GENESYS Cloud
			go func(){
				log := server.Logger.Child(nil, "sendCX", "chat", message.Chat.ID, "message", message.ID)

				log.Debugf("Sending message to GENESYS Cloud")
				inboundResult, err := config.Integration.SendInboundTextMessage(
					&purecloud.OpenMessageFrom{
						ID:        message.UserID,
						Type:      "email",
						Firstname: "Bob",
						Lastname:  "Minion",
						Nickname:  "",
						ImageURL:  core.Must(url.Parse("https://gravatar.com/avatar/97959eb8244f0cb560e2d30b2075f013?s=400&d=robohash&r=x")).(*url.URL),
					},
					message.ID,
					message.Content,
				)
				if err != nil {
					Log.Errorf("Failed to send inbound", err)
				} else {
					Log.Record("message", inboundResult.ID).Record("result", inboundResult).Infof("Message sent successfully")
				}

				log.Infof("Message sent to GENESYS Cloud")
			}()
		}
	}
}