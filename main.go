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

//var guest string = "Guest"

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

type MapTwo map[string]int

func (m MapTwo) keys() []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func listallusers(Users map[string]User) (listofusers string) {
	listofusers = "Users online: "
	//Users: make(map[string]User),
	//var m map[int]string
	//var keys []int
	//for l := range Users {
	//	keys = append(keys, k)
	//	}
	for _, u := range Users {
		//io.WriteString("User:", k, "Value:", m[k])
		listofusers = listofusers + u
	}
	//for _, k := range keys {
	//	fmt.Println("Key:", k, "Value:", m[k])
	//	}
	io.WriteString(conn, listofusers)
	return listofusers
}

func guestassignname(guestname string) (guest string) {
	prefixguestname := "Guest"
	guestid := rand.Intn(999)
	guestname = (prefixguestname + strconv.Itoa(guestid))
	return guestname
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

func swap(x, y string) (string, string) {
	return y, x
}

func changeNick(input string, nickprefix string) (changednick string) {
	//fmt.Println(input)
	//newname = strings.TrimPrefix(input, "/nick ")
	//n := Newname{}
	//n.SetName(newname)
	//nn := n.GetName()
	//fmt.Println(nn)
	//io.WriteString(conn, newname)
	//fmt.Println(newname)
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

func handleConn(srvr *Server, conn net.Conn) {
	defer conn.Close()
	//customusername := defaultName
	//io.WriteString(conn, "Use default userame? Y/N")
	//duscanner := bufio.NewScanner(conn)
	//duscanner.Scan()
	//if duscanner.Text == "Y"{
	//	customusername = defaultName
	//}
	//else if duscanner.Text == "N"{
	//	io.WriteString(conn, "Enter username:")
	//ryb enteeruseername fn
	//}
	//else {
	//	io.WriteString(conn, "Not valid entry so defaulting to default username")
	//		customusername = defaultName
	//
	//io.WriteString(conn, "Enter username:")
	guest := "Guest"
	//scanner := bufio.NewScanner(conn)
	//scanner.Scan()
	//if len(scanner.Text()) == 0 {
	//	user := User{
	//		Name:   defaultName,
	//		Output: make(chan Message, 10),
	//	}
	//} else {
	user := User{
		Name:   guestassignname(guest),
		Output: make(chan Message, 10),
	}

	//}
	srvr.Join <- user

	//print all users
	//for _, users
	listallusers(Users)
	//for all strings in the map of users (user has a string called name and we need to print all names of the map)
	//for
	//case msg := <-srvr.Input:
	//	for _, user := range srvr.Users {
	//		select {
	//		case user.Output <- msg:
	//		default:

	scanner := bufio.NewScanner(conn)
	//scanner.Scan()

	defer func() {
		srvr.Leave <- user
	}()

	//read from conn
	go func() {
		for scanner.Scan() {
			ln := scanner.Text()
			nickprefix := `/nick`
			//fmt.Println(ln)
			if strings.HasPrefix(ln, nickprefix) {
				nn := changeNick(ln, nickprefix)
				io.WriteString(conn, "Changed nickname. ")
				io.WriteString(conn, "Nickname changed to: ")
				io.WriteString(conn, nn)
				user.Name = nn
			} else if strings.HasPrefix(ln, "/register") {

				//call register fn
				io.WriteString(conn, "register nick")
			} else {
				srvr.Input <- Message{user.Name, ln}
			}
		}
	}()

	//write to connection
	//for users.names := range user.Output {
	//	io.WriteString(usernames)
	//}
	for msg := range user.Output {
		if msg.Nickname != user.Name {
			_, err := io.WriteString(conn, msg.Nickname+": "+msg.Msgtext+"\n")
			if err != nil {
				break
			}
		}
	}
}

func register(something1 string, something2 string) (something3 string) {
	registerednick := strings.TrimPrefix(input, "/register ")
	//fmt.Println(input)
	//newname = strings.TrimPrefix(input, "/nick ")
	//n := Newname{}
	//n.SetName(newname)
	//nn := n.GetName()
	//fmt.Println(nn)
	//io.WriteString(conn, newname)
	//fmt.Println(newname)

	//registration: (ideas?)
	changednick = strings.TrimPrefix(input, "/nick ")
	return changednick
}

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

	go mainServer.Run()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleConn(mainServer, conn)
	}
}
