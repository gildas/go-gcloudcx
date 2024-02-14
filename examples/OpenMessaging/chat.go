package main

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Chat struct {
	ID     uuid.UUID
	UserID string

	server       *ChatServer
	connection   *websocket.Conn
	messages     []*ChatMessage
	send         chan []byte
	readTimeout  time.Duration
	writeTimeout time.Duration
	Logger       *logger.Logger
}

type ChatMessageError struct {
	ID        string `json:"messageId"`
	RequestID string `json:"reqid"`
	Error     string `json:"error"`
}

func NewChat(server *ChatServer, userID string) *Chat {
	chatid := uuid.New()
	return &Chat{
		ID:           chatid,
		UserID:       userID,
		server:       server,
		messages:     []*ChatMessage{},
		send:         make(chan []byte, 256),
		readTimeout:  60 * time.Second,
		writeTimeout: 10 * time.Second,
		Logger:       server.Logger.Child("chat", "chat", "chat", chatid.String()),
	}
}

func (chat *Chat) Serve(connection *websocket.Conn) {
	chat.connection = connection
	go chat.readMessageLoop()
	go chat.writeMessageLoop()
}

func (chat Chat) String() string {
	return chat.ID.String()
}

func (chat *Chat) readMessageLoop() {
	log := chat.Logger.Child(nil, "readloop")

	defer func() { // tell the server we are done and game over
		chat.server.unregister <- chat
		log.Debugf("Guest closing chat connection")
		if chat.connection != nil {
			chat.connection.Close()
			log.Infof("Chat connection closed")
		}
		log.Debugf("Closing send channel")
		close(chat.send)
		log.Infof("Send Channel closed")
	}()

	// Setting read limits and timeouts
	chat.connection.SetReadLimit(1024)
	if err := chat.connection.SetReadDeadline(time.Now().Add(chat.readTimeout)); err != nil {
		log.Errorf("Failed setting read deadline", err)
	}
	chat.connection.SetPongHandler(func(data string) error {
		return chat.connection.SetReadDeadline(time.Now().Add(chat.readTimeout))
	})

	for {
		log.Tracef("Waiting for messages from Chat GUI")
		messageType, payload, err := chat.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Errorf("Chat error", err)
			} else {
				log.Infof("Chat was closed by the GUI, ending read loop")
			}
			break
		}

		log.Record("type", messageType).Debugf("Received new message: %s", string(payload))
		message := &ChatMessage{}
		if err := json.Unmarshal(payload, &message); err != nil {
			log.Errorf("Failed to unmarshal data: \"%s\"", string(payload), err)
			response, _ := json.Marshal(ChatMessageError{Error: err.Error()})
			chat.send <- response
			continue
		}
		message.UserID = chat.UserID
		log.Record("message", message).Infof("Received Chat Message")
		chat.messages = append(chat.messages, message)
		message.Chat = chat
		chat.server.sendCX <- message
	}
}

func (chat *Chat) writeMessageLoop() {
	log := chat.Logger.Child(nil, "writeloop")
	ticker := time.NewTicker((chat.readTimeout * 9) / 10) // 90% of the pong/read wait
	defer func() {
		log.Debugf("Stopping ticker")
		if ticker != nil {
			ticker.Stop()
			log.Infof("Ticker Stopped")
		}
	}()

	for {
		select {
		case message, ok := <-chat.send:
			if err := chat.connection.SetWriteDeadline(time.Now().Add(chat.writeTimeout)); err != nil {
				log.Errorf("Failed setting write deadline", err)
			}
			if !ok {
				if err := chat.connection.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Errorf("Failed sending Close Message", err)
				}
				return
			}

			writer, err := chat.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Errorf("Failed getting a websocket writer", err)
				return
			}
			count, err := writer.Write(message)
			if err != nil {
				log.Errorf("Failed writing message", err)
				break
			}
			log.Infof("Written %d bytes", count)

		case <-ticker.C: // the ticker is used to send pings to the websocket
			if err := chat.connection.SetWriteDeadline(time.Now().Add(chat.writeTimeout)); err != nil {
				log.Errorf("Failed setting write deadline", err)
			}
			if err := chat.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Errorf("Failed to send ping", err)
				} else {
					log.Infof("Chat is closed, ending write message loop")
				}
				return
			}
		}
	}
}
