/*

Task: Implement telnet chat server.
Description: Implement a TCP server application performing as a chat server. Clients should be able to connect to the listening port using plaintext protocol and be able to communicate with each other. Messages are separated by <LF>, when connected user should be presented with a list of users currently online, everyone see messages from everyone, server should support /nick <nickname> command for users to be able to redefine the default auto-assigned nickname "GuestNNN" and /register command for users to be able to protect their nickname from being taken by other user with a password.
Language choice: Go, C++.

*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

var ServerName string = "Server"
var defaultName string = "GuestNNN"

type Newname struct {
	Name string
}

type User struct {
	Name   string
	Output chan Message
}

type Message struct {
	Nickname string
	Msgtext  string
}

type Server struct {
	Users map[string]User
	Join  chan User
	Leave chan User
	Input chan Message
}

func (srvr *Server) Run() {
	for {
		select {
		case user := <-srvr.Join:
			srvr.Users[user.Name] = user
			go func() {
				srvr.Input <- Message{
					Nickname: "System Message",
					Msgtext:  fmt.Sprintf("%s has joined", user.Name),
				}
			}()
		case user := <-srvr.Leave:
			delete(srvr.Users, user.Name)
			go func() {
				srvr.Input <- Message{
					Nickname: "System Message",
					Msgtext:  fmt.Sprintf("%s has left", user.Name),
				}

			}()

		case msg := <-srvr.Input:
			for _, user := range srvr.Users {
				select {
				case user.Output <- msg:
				default:
				}
			}
		}
	}
}

func handleConn(srvr *Server, conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	//scanner.Scan()
	user := User{
		Name:   defaultName,
		Output: make(chan Message, 10),
	}
	srvr.Join <- user

	//print list of users connected to server

	defer func() {
		srvr.Leave <- user
	}()
	go func() {
		for scanner.Scan() {
			ln := scanner.Text()
			srvr.Input <- Message{user.Name, ln}
		}
	}()

	//write to connection
	for msg := range user.Output {
		if msg.Nickname != user.Name {
			_, err := io.WriteString(conn, msg.Nickname+": "+msg.Msgtext+"\n")
			if err != nil {
				break
			}
		}
	}
}

func (n *Newname) SetName(Name string) {
	n.Name = Name
	//user := User{
	//	Name:	scanner.Text(),
	//	Output:	make(chan Message, 10),
}

func (n Newname) GetName() string {
	return n.Name
}

func changeNick(newname string) {
	n := Newname{}
	n.SetName(newname)
	nn := n.GetName()
	fmt.Println(nn)
}

//func register(x,y,z?){
//}

func main() {
	server, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer server.Close()

	mainServer := &Server{
		Users: make(map[string]User),
		Join:  make(chan User),
		Leave: make(chan User),
		Input: make(chan Message),
	}

	go mainServer.Run()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleConn(mainServer, conn)
	}
}
