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
			go func() {
				srvr.Input <- Message{
					Nickname: "GuestNNN",
					Msgtext: fmt.Sprintf("%s has joined", user.Name),
				}
			}()
		case user := <- srvr.Leave:
			delete(srvr.Users, user.Name)
			go func() {
				srvr.Input <- Output{
					Nickname: "GuestNNN",
					Msgtext: fmt.Sprintf("%s has left", user.Name),
				}

			}()

		case msg := <- srvr.Input: 
			for _, user := range srvr.Users{
				select{
				case user.Output <- msg:
				default:
				}
			}
		}
	}
}

func handleConn(srvr *Server, conn  net.Conn) (
	defer conn.Close()


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

