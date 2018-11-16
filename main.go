package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"log"
)

type User struct {
	Name string
	Output chan Msg
}

type Line struct{
	Nickname string
	Msgtext string
}

type Server struct{
	Users map[string]User
	Join chan User
	Leave chan User
	Input chan Line
}

func (srvr *Server) Run(){
	for{
		select{
		case user := <-srvr.Join:
			srvr.User[user.Name] = user
			go func(){
				srvr.Input <- Message{
					Username: "GuestNNN",
					Msgtext: fmt.Sprintf("%s joined", user.Name),
				}
			}()
		case user := < srvr.Leave:
			delete(cs.Users, user.Name)
			go func() {
				srvr.Input <- Msgtext

		case msg := <- srvr.Input:
			for _, user :=




func handleConn

func changeNick

func register

func main(){
	server, err := net.Listen("tcp", ":9011")
	if err !=nil {
		log.Fatalln(err.Error())
	}
	defer server.Close()

	mainServer := &MainServer{
		Users
		Join
		Leave
		Input

	go mainServer.Run()

	for {
		conn, err:= server



		go handleConn(chatServer, conn)
	}
}
}

