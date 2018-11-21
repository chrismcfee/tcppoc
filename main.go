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
	//"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

var ServerName string = "Server"
var defaultName string = "GuestNNN"

var UserMap map[string]int

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func addToUserMap(id int, key string) {
	UserMap[key] = id
}

func delFromUserMap(delname string) {
	delete(UserMap, delname)

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

func assignid() (id int) {
	newid := rand.Intn(999)
	return newid
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
			delFromUserMap(user.Name)
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
	io.WriteString(conn, "Users Online Now: ")
	for k := range UserMap {
		io.WriteString(conn, k)
		io.WriteString(conn, ",")
	}
	//writeString
	addToUserMap(assignid(), user.Name)

	//for index, each := range UserSlice {
	//	fmt.Printf("value [%d] is [%s]\n", index, each)
	//}
	for k := range UserMap {
		fmt.Printf("key[%s] value[%s]\n", k, UserMap[k])
	}

	scanner := bufio.NewScanner(conn)

	defer func() {
		srvr.Leave <- user
	}()

	//read from conn
	go func() {
		for scanner.Scan() {
			ln := scanner.Text()
			nickprefix := `/nick `
			registerprefix := `/register `
			loginprefix := `/login `
			addToUserMap(assignid(), user.Name)
			if strings.HasPrefix(ln, nickprefix) {
				nn := changeNick(ln, nickprefix)
				delFromUserMap(user.Name)
				io.WriteString(conn, "Changed nickname. ")
				io.WriteString(conn, "Nickname changed to: ")
				io.WriteString(conn, nn)
				addToUserMap(assignid(), nn)
				user.Name = nn
			} else if strings.HasPrefix(ln, "/register") {
				register(ln, user.Name, registerprefix)
				//io.WriteString(conn, "registered name")
			} else if strings.HasPrefix(ln, "/login") {
				//io.WriteString(conn, "login")
				login(ln, user.Name, loginprefix)
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

func register(input string, username string, registerprefix string) {
	delimiter := " "
	registrationinput := strings.TrimPrefix(input, registerprefix)

	f, err := os.OpenFile("usernameregistrations.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanfile, err := os.Open("usernameregistrations.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer scanfile.Close()

	scanner := bufio.NewScanner(scanfile)
	for scanner.Scan() {
		leftOfDelimiter := strings.Split(scanner.Text(), delimiter)[0]
		if leftOfDelimiter == username {
			fmt.Println("someone attempting to hijack registered account")
			return

			//io.WriteString(conn, "user already registered...")
		}
	}

	if _, err = f.WriteString(username + " " + registrationinput + "\n"); err != nil {
		panic(err)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func login(input string, username string, loginprefix string) (loginsuccess bool) {
	delimiter := " "
	logininputpw := strings.TrimPrefix(input, loginprefix)

	f, err := os.OpenFile("usernameregistrations.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanfile, err := os.Open("usernameregistrations.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer scanfile.Close()

	scanner := bufio.NewScanner(scanfile)
	for scanner.Scan() {
		leftOfDelimiter := strings.Split(scanner.Text(), delimiter)[0]
		rightOfDelimiter := strings.Join(strings.Split(scanner.Text(), delimiter)[1:], delimiter)
		if leftOfDelimiter == username {

			if rightOfDelimiter == logininputpw {
				fmt.Println("correct password match for login")
				return true
			} else if rightOfDelimiter != logininputpw {
				fmt.Println("User trying to login with invalid password")
				return false
			} else {
				fmt.Println("unexpected thing happened")
				return false
			}
			//fmt.Println("someone attempting to hijack registered account")

			return

			//io.WriteString(conn, "user already registered...")
		}

	}

	fmt.Println("user not found so login has failed")
	//if _, err = f.WriteString(username + " " + registrationinput + "\n"); err != nil {
	//	panic(err)
	//}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return false
}

func CreatePasswordFile() {
	file, err := os.Create("usernameregistrations.txt")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer file.Close()
	len, err := file.WriteString("Password List Format: Username Password \n")
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}

	fmt.Printf("\nLength: %d bytes", len)
	fmt.Printf("\nFile Name: %s", file.Name())
}

func main() {
	if _, err := os.Stat("./usernameregistrations.txt"); os.IsNotExist(err) {
		CreatePasswordFile()
	}

	UserMap = make(map[string]int)
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
		go handleConn(mainServer, conn, mainServer.Users)
	}
}
