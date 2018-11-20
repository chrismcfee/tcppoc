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
	"math/rand"
	"net"
	"strconv"
	"strings"
)

var ServerName string = "Server"
var defaultName string = "GuestNNN"

var UserSlice = make([]string, 0, 999)

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

func addusertolist(addedname string) {
	UserSlice = append(UserSlice, addedname)
}

func listallusers(Users map[string]User) (listofusers_result string) {
	var curlist string
	for _, u := range Users {
		curlist = (" " + u.Name + " ")
	}
	listofusers_result = curlist
	return listofusers_result
}

func guestassignname(guestname string) (guest string) {
	prefixguestname := "Guest"
	guestid := rand.Intn(999)
	guestname = (prefixguestname + strconv.Itoa(guestid))
	return guestname
}

func changeNick(input string, nickprefix string) (changednick string) {
	changednick = strings.TrimPrefix(input, "/nick ")
	return changednick
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

func handleConn(srvr *Server, conn net.Conn, Users map[string]User) {
	defer conn.Close()
	guest := "Guest"
	user := User{
		Name:   guestassignname(guest),
		Output: make(chan Message, 10),
	}
	srvr.Join <- user
	addusertolist(user.Name)

	for index, each := range UserSlice {
		fmt.Printf("value [%d] is [%s]\n", index, each)
	}

	scanner := bufio.NewScanner(conn)

	defer func() {
		srvr.Leave <- user
	}()

	//read from conn
	go func() {
		for scanner.Scan() {
			ln := scanner.Text()
			nickprefix := `/nick`
			addusertolist(user.Name)
			if strings.HasPrefix(ln, nickprefix) {
				nn := changeNick(ln, nickprefix)
				io.WriteString(conn, "Changed nickname. ")
				io.WriteString(conn, "Nickname changed to: ")
				io.WriteString(conn, nn)
				addusertolist(nn)
				user.Name = nn
			} else if strings.HasPrefix(ln, "/register") {
				io.WriteString(conn, "register nick")
			} else if strings.HasPrefix(ln, "/login") {
				io.WriteString(conn, "login")
			} else {
				srvr.Input <- Message{user.Name, ln}
			}
		}
	}()
	for msg := range user.Output {
		if msg.Nickname != user.Name {
			_, err := io.WriteString(conn, msg.Nickname+": "+msg.Msgtext+"\n")
			if err != nil {
				break
			}
		}
	}
}

//func register(something1 string, something2 string) (something3 string) {
//	registerednick := strings.TrimPrefix(input, "/register ")
//fmt.Println(input)
//newname = strings.TrimPrefix(input, "/nick ")
//n := Newname{}
//n.SetName(newname)
//nn := n.GetName()
//fmt.Println(nn)
//io.WriteString(conn, newname)
//fmt.Println(newname)

//registration: (ideas?)
//	changednick = strings.TrimPrefix(input, "/nick ")
//	return changednick
//}

func main() {
	server, err := net.Listen("tcp", ":9009")
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

	UserSlice = append(UserSlice, "Users Online: ")

	go mainServer.Run()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleConn(mainServer, conn, mainServer.Users)
	}
}
