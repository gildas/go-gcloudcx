package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ChatRoutes declares the routes for the chat server
func ChatRoutes(router *mux.Router) {
	router.Path("/chat/ws/{chatid}").HandlerFunc(chatSocketHandler)
	router.Methods("POST").Path("/chat").HandlerFunc(createChatHandler)
}

// createChatHandler creates a new chat
func createChatHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.Must(logger.FromContext(r.Context())).Child(nil, "create_chat")
	config := core.Must(ConfigFromContext(r.Context()))

	log.Debugf("Request Headers: %#v", r.Header)

	// Analyzing the body
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed reading the request body", err)
		core.RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	log.Infof("Received body: %s", string(body[:]))

	chatConfig := struct{
		Account    string `json:"account"`
		Secret     string `json:"secret"`
		UserID     string `json:"userId"`
	}{}

	if err := json.Unmarshal(body, &chatConfig); err != nil {
		log.Errorf("Failed unmarshaling message", err)
		core.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	chat := config.ChatServer.CreateChat(chatConfig.UserID)

	log.Infof("Returning HTTP 200")
	core.RespondWithJSON(w, http.StatusOK, struct {Path string `json:"path"`}{Path: "/chat/ws/" + chat.ID.String()})
}

// chatSocketHandler starts the websocket that handles chats
func chatSocketHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.Must(logger.FromContext(r.Context())).Child(nil, "create_websocket")
	config := core.Must(ConfigFromContext(r.Context()))

	log.Debugf("Request Headers: %#v", r.Header)
	params := mux.Vars(r)

	value, found := params["chatid"]
	if !found {
		log.Errorf("Request parameter chatid is missing")
		core.RespondWithError(w, http.StatusBadRequest, errors.ArgumentMissing.With("chatid"))
		return
	}
	chatid, err := uuid.Parse(value)
	if err != nil {
		log.Errorf("Chat %s is not uuid", value, err)
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error while upgrading to websocket", err)
		core.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	chat, err := config.ChatServer.FindChatByID(chatid)
	if err != nil {
		log.Errorf("Chat %s not found", chatid.String())
		core.RespondWithError(w, http.StatusNotFound, err)
		return
	}

	chat.Serve(connection)
}