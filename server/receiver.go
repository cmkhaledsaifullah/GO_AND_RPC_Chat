package main

import (
	"fmt"
	"crypto/rand"
	"log"
	"encoding/base64"
)


type Args struct {
	Token  string
	String string
}

type UserNameArgs struct {
	Name  string
}

type Receiver int

func (r *Receiver) Connect(args *UserNameArgs, token *string) error {
	log.Println(args.Name+ " is now connected!!!")
	//go func() { client.Outgoing <- MSG_CONNECT }()
	*token = randomString(64)
	client := NewClient(*token, args.Name)
	err := AddClient(client)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}


func (r *Receiver) SendMessage(args *Args, _ *struct{}) error {
	client, err := GetClient(args.Token)
	log.Println(client.Name+ "Send Messages")
	if err != nil {
		log.Println(err)
		return err
	}
	client.Mutex.RLock()
	defer client.Mutex.RUnlock()
	if client.ChatRoom == nil {
		client.Outgoing <- ERROR_SEND
		return nil
	}
	client.ChatRoom.Incoming <- fmt.Sprintf("%s: %s", client.Name, args.String)
	return nil
}
func (r *Receiver) ReceiveMessage(token *string, message *string) error {
	log.Println("ReceiveMessage")
	client, err := GetClient(*token)
	if err != nil {
		return err
	}
	*message = <-client.Outgoing
	return nil
}
func (r *Receiver) ListChatRooms(token *string, _ *struct{}) error {
	log.Println("ListChatRooms")
	client, err := GetClient(*token)
	if err != nil {
		log.Println(err)
		return err
	}
	chatRoomNames := GetChatRoomNames()
	chatList := "\nChatRooms:\n"
	for _, chatRoomName := range chatRoomNames {
		chatList += chatRoomName + "\n"
	}
	chatList += "\n"
	client.Outgoing <- chatList
	return nil
}

func (r *Receiver) JoinChatRoom(args Args, _ *struct{}) error {
	log.Println("JoinChatRoom")
	client, err := GetClient(args.Token)
	if err != nil {
		log.Println(err)
		return err
	}
	chatRoom, err := GetChatRoom(args.String)
	if err != nil {
		client.Outgoing <- err.Error()
		log.Println(err)
		return err
	}
	client.Mutex.RLock()
	oldChatRoom := client.ChatRoom
	client.Mutex.RUnlock()
	if oldChatRoom != nil {
		oldChatRoom.Leave <- client
	}

	chatRoom.Join <- client
	return nil
}
func (r *Receiver) CreateChatRoom(args Args, _ *struct{}) error {
	log.Println("CreateChatRoom")
	client, err := GetClient(args.Token)
	if err != nil {
		log.Println(err)
		return err
	}
	chatRoom := NewChatRoom(args.String)
	err = AddChatRoom(chatRoom)
	if err != nil {
		client.Outgoing <- err.Error()
		log.Println(err)
		return err
	}
	client.Outgoing <- fmt.Sprintf(NOTICE_PERSONAL_CREATE, chatRoom.Name)
	return nil
}

func (r *Receiver) SwitchChatRoom(args Args, _ *struct{}) error {
	log.Println("Switch Chat Room")
	client, err := GetClient(args.Token)
	if err != nil {
		log.Println(err)
		return err
	}
	//leaving the existing chatroom
	client.Mutex.RLock()
	oldChatRoom := client.ChatRoom
	client.Mutex.RUnlock()
	if oldChatRoom != nil {
		oldChatRoom.Leave <- client
	}
	chatRoom, err := FindChatRoom(args.String)
	if err != nil {
		//creating new chatroom
		newchatRoom := NewChatRoom(args.String)
		err = AddChatRoom(newchatRoom)
		if err != nil {
			client.Outgoing <- err.Error()
			log.Println(err)
			return err
		}
		client.Outgoing <- fmt.Sprintf(NOTICE_PERSONAL_CREATE, newchatRoom.Name)
		newchatRoom.Join <- client
	}else {
		//switching to another existing chatroom
		chatRoom.Join <- client
	}
	return nil
}

func (r *Receiver) LeaveChatRoom(token *string, _ *struct{}) error {
	log.Println("LeaveChatRoom")
	client, err := GetClient(*token)
	if err != nil {
		log.Println(err)
		return err
	}
	client.Mutex.RLock()
	defer client.Mutex.RUnlock()
	client.ChatRoom.Leave <- client
	return nil
}
func (r *Receiver) Quit(token *string, _ *struct{}) error {
	log.Println("Quit")
	client, err := GetClient(*token)
	if err != nil {
		log.Println(err)
		return err
	}
	if client.ChatRoom != nil{
		client.Mutex.RLock()
		defer client.Mutex.RUnlock()
		client.ChatRoom.Leave <- client
	}
	err = RemoveClient(*token)
	if err != nil {
		return err
	}
	return nil
}

//rANDOM nUMBER gENERATOR FOR ASSIGING TOKEN TO EACH USER/CLIENT
func randomString(length int) (str string) {
	arraybyte := make([]byte, length)
	rand.Read(arraybyte)
	return base64.StdEncoding.EncodeToString(arraybyte)
}
