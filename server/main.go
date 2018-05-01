package main

import (
	"net"
	"net/http"
	"net/rpc"
	"time"
	"log"
)

const (
	MAX_USERS = 20

	ERROR_PREFIX   = "Error: "
	ERROR_SEND     = ERROR_PREFIX + "You Need to Join/Create a Chat Room to send a message. \n"
	ERROR_CREATE   = ERROR_PREFIX + "Duplicate Name of Chat Room. Enter another name\n"
	ERROR_JOIN     = ERROR_PREFIX + "No Such ChatRoom Exits\n"
	ERROR_SWITCH     = ERROR_PREFIX + "No Such ChatRoom Exits\n A New Chat Room is creating..."
	ERROR_LEAVE    = ERROR_PREFIX + "You are out of the chatroom already.\n"
	ERROR_TOKEN    = ERROR_PREFIX + "Internal Server Error. Please Try Again!!!.\n"
	ERROR_NO_TOKEN = ERROR_PREFIX + "No such user exits. Internal Server Error. Plese disconnect the applicationa and start again!!!\n"


	NOTICE_ROOM_JOIN       = "\"%s\" joined\n"
	NOTICE_ROOM_LEAVE      = "\"%s\" left\n"
	NOTICE_ROOM_DELETE     = "Chat room is found inactive for seven days and deleted.\n"
	NOTICE_PERSONAL_CREATE = "Welcome to Chat Room \"%s\".\n"

	EXPIRY_TIME time.Duration = 7 * 24 * time.Hour
)

func main() {
	newReciever := new(Receiver)
	rpc.Register(newReciever)
	rpc.HandleHTTP()
	t, p := net.Listen("tcp", ":9000")
	if p != nil {
		log.Fatal("listen error:", p)
	}
	http.Serve(t, nil)
}
