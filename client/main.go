package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"sync"
)
type Args struct {
	Token  string
	String string
}

type UserNameArgs struct {
	Name  string
}
const (
	CMD_CREATE =  "create"
	CMD_LIST   =  "list"
	CMD_JOIN   =  "join"
	CMD_LEAVE  =  "leave"
	CMD_SWITCH  =  "switch"
	CMD_COMMANDS   =  "commands"
	CMD_QUIT   =  "quit"

	MSG_COMMANDS = "\nCommands:\n[" +
		CMD_LIST + "] to see the list of all rooms\n[" +
		CMD_JOIN + " room_name] to join an existing room\n[" +
		CMD_CREATE + " room_name] to create a new room\n["+
		CMD_SWITCH + "room_name] to switch to another room\n[" +
		CMD_LEAVE + "] to leave the existing room\n["+
		CMD_COMMANDS+"] to see the available Commands\n["+
		CMD_QUIT+"] to quit from the system\n\n"


	MSG_CONNECT = "Welcome to USASK ChatRoom.\nEnter the name you want to display:\n>"
	MSG_DISCONNECT = "Thanks for Using USASK Chatroom Application. See you Next time!!!\n"
)

var token string
var client *rpc.Client
var waitGroup sync.WaitGroup

// Adds strings from stdin to the server.
func Input() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		str, err := reader.ReadString('\n')
		if err != nil {
			waitGroup.Done()
			break
		}
		Parse(str)
	}
}

func UserNameInput() string{
	var str string
	fmt.Scanf("%s", &str)
	return str
}

func Parse(str string) (err error) {
	switch {
	default:
		err = client.Call("Receiver.SendMessage", Args{token, str}, nil)
	case strings.HasPrefix(str, CMD_CREATE):
		name := strings.TrimSuffix(strings.TrimPrefix(str, CMD_CREATE+" "), "\n")
		err = client.Call("Receiver.CreateChatRoom", Args{token, name}, nil)
		err = client.Call("Receiver.JoinChatRoom", Args{token, name}, nil)
	case strings.HasPrefix(str, CMD_LIST):
		err = client.Call("Receiver.ListChatRooms", &token, nil)
	case strings.HasPrefix(str, CMD_JOIN):
		name := strings.TrimSuffix(strings.TrimPrefix(str, CMD_JOIN+" "), "\n")
		err = client.Call("Receiver.JoinChatRoom", Args{token, name}, nil)
	case strings.HasPrefix(str, CMD_SWITCH):
		name := strings.TrimSuffix(strings.TrimPrefix(str, CMD_SWITCH+" "), "\n")
		err = client.Call("Receiver.SwitchChatRoom", Args{token, name}, nil)
	case strings.HasPrefix(str, CMD_LEAVE):
		err = client.Call("Receiver.LeaveChatRoom", &token, nil)
	case strings.HasPrefix(str, CMD_COMMANDS):
		fmt.Print(MSG_COMMANDS)
	case strings.HasPrefix(str, CMD_QUIT):
		err = client.Call("Receiver.Quit", &token, nil)
		waitGroup.Done()
	}
	return err
}

// Requests strings from the server and outputs them to stdout.
func Output() {
	for {
		var message string
		err := client.Call("Receiver.ReceiveMessage", &token, &message)
		if err != nil {
			waitGroup.Done()
			break
		}
		fmt.Print(message)
	}
}



//Main Function
func main() {
	waitGroup.Add(1)

	var err error
	client, err = rpc.DialHTTP("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(MSG_CONNECT)
	var username = UserNameInput()
	err = client.Call("Receiver.Connect", &UserNameArgs{username}, &token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(MSG_COMMANDS)

	go Input()
	go Output()

	waitGroup.Wait()
	fmt.Print(MSG_DISCONNECT)
}
