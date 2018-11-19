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

func listallusers(Users map[string]User, userlist string, newuser string) (listofusers_result string) {
	//userlist = userlist + cur_list
	var curlist string
	for _, u := range Users {
		curlist = (" " + u.Name + " ")
	}
	curlist = curlist + " " + newuser
	listofusers_result = userlist + curlist
	//io.WriteString(conn, listofusers)
	//fmt.Println(listofusersresult)
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

func handleConn(srvr *Server, conn net.Conn, Users map[string]User, userlist string) {
	defer conn.Close()
	guest := "Guest"
	user := User{
		Name:   guestassignname(guest),
		Output: make(chan Message, 10),
	}
<<<<<<< HEAD
	newuser := user.Name
	srvr.Join <- user
	listofusers := listallusers(Users, userlist, newuser)
	io.WriteString(conn, listofusers)
=======
	//newuser := user.Name
	srvr.Join <- user
	//listofusers := listallusers(Users, userlist, newuser)
	//io.WriteString(conn, listofusers)
>>>>>>> 64410827b1dcc3149a5ca2fd80a9dacd2a152e67

	scanner := bufio.NewScanner(conn)

	defer func() {
		srvr.Leave <- user
	}()

	//read from conn
	go func() {
		//io.WriteString(conn, listofusers)
		for scanner.Scan() {
			ln := scanner.Text()
			nickprefix := `/nick`
			if strings.HasPrefix(ln, nickprefix) {
				nn := changeNick(ln, nickprefix)
				io.WriteString(conn, "Changed nickname. ")
				io.WriteString(conn, "Nickname changed to: ")
				io.WriteString(conn, nn)
				user.Name = nn
<<<<<<< HEAD
				listofusers = listallusers(Users, userlist, nn)
=======
				//listofusers = listallusers(Users, userlist, nn)
				//fmt.Println(listofusers)
>>>>>>> 64410827b1dcc3149a5ca2fd80a9dacd2a152e67
				//io.WriteString(conn, listofusers)
				//listofusers = listallusers(Users)
				//else if strings.HasPrefix(ln, "/register") {
				//call register fn
				//rr := registerNick(ln, registrationPrefix, registrationPassword)
				//io.WriteString
<<<<<<< HEAD
				io.WriteString(conn, "register nick")
			} else if strings.HasPrefix(ln, "/login") {
				//ll := loginNick(ln, loginPrefix, loginPassword)
				io.WriteString(conn, "login")
=======
				//io.WriteString(conn, "register nick")
				//	} //else if strings.HasPrefix(ln, "/login") {
				//		ll :=loginNick(ln, loginPrefix, loginPassword)
				//io.WriteString
				//	}
>>>>>>> 64410827b1dcc3149a5ca2fd80a9dacd2a152e67
			} else {
				srvr.Input <- Message{user.Name, ln}
			}
		}
		//	}
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
	//var listofuserst string
	//listofuserst = listallusers(mainServer.Users, listofuserst)

	go mainServer.Run()

	userlist := "Users online: "

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleConn(mainServer, conn, mainServer.Users, userlist)
	}
}
