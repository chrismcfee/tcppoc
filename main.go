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
	"io/ioutil"
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

//m := make(map[int]string)
//var UserSlice = make([]string, 0, 999)

//var UserMap map[int]string

//type connID struct {
//	id int
//}

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

	//var sessions =  map[string] chan int{}
	//var sessions = map[string] chan int{};
	//sessions["moo"] = make (chan int);
	//_, ok := sessions["moo"];
	//if ok {
	//    delete(sessions, "moo");
}

//func addusertolist(addedname string) {
//	UserSlice = append(UserSlice, addedname)
//}

//func deluserfromlist(deletedname string) {
//	//a := []string{"A", "B", "C", "D", "E"}
//	i := len(UserSlice) - 1
//	UserSlice[i] = UserSlice[(len(UserSlice))-1]
//	UserSlice[len(UserSlice)-1] = ""
//	UserSlice = UserSlice[:len(UserSlice)-1]
//	// Remove the element at index i from userslice
//}

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
			nickprefix := `/nick`
			registerprefix := `/register `
			//addusertolist(user.Name)
			addToUserMap(assignid(), user.Name)
			if strings.HasPrefix(ln, nickprefix) {
				nn := changeNick(ln, nickprefix)
				delFromUserMap(user.Name)
				io.WriteString(conn, "Changed nickname. ")
				io.WriteString(conn, "Nickname changed to: ")
				io.WriteString(conn, nn)
				addToUserMap(assignid(), nn)
				//addusertolist(nn)
				user.Name = nn
			} else if strings.HasPrefix(ln, "/register") {
				//io.WriteString(conn, "registered username")
				register(ln, user.Name, registerprefix)
				io.WriteString(conn, "registered name")
			} else if strings.HasPrefix(ln, "/login") {
				io.WriteString(conn, "login")
				//} //else if strings.HasPrefix(ln, "/names") {
				//for _, each := range UserSlice {
				//	io.WriteString(conn, each+" ")
				//}
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
			//io.WriteString(conn, "user already registered...")
		} else {

			if _, err = f.WriteString(username + " " + registrationinput + "\n"); err != nil {
				panic(err)
			}

		}
		//rightOfDelimiter := strings.Join(strings.Split(scanner.Text(), delimiter)[1:], delimiter)
		//fmt.Println(leftOfDelimiter)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

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

func TestAddLine() {
	//file, err:=
	//defer file.Close()
	//len, err := file.WriteString
	//mydata := []byte("test\n")

	// the WriteFile method returns an error if unsuccessful
	//err := ioutil.WriteFile("usernameregistrations.txt", mydata, 0777)
	// handle this error
	//if err != nil {
	// print it out
	//	fmt.Println(err)
	//}

	data, err := ioutil.ReadFile("usernameregistrations.txt")
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Print(string(data))

	f, err := os.OpenFile("usernameregistrations.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString("adding line!\n"); err != nil {
		panic(err)
	}

	data, err = ioutil.ReadFile("usernameregistrations.txt")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(string(data))

}
func main() {
	if _, err := os.Stat("./usernameregistrations.txt"); os.IsNotExist(err) {
		// path/to/whatever does not exist
		CreatePasswordFile()
	}

	//if _, err := os.Stat("./usernameregistrations.txt"); err == nil {
	//TestAddLine()
	//	fmt.Println("Username registration file already exists. Continuing...")
	// path/to/whatever exists
	//}

	//testing write to file here
	//f1 := []byte("hello\ngo\n")
	//err := ioutil.WriteFile("/tmp/dat1", f1, 0644)
	//check(err)
	//f, err := os.Create("/tmp/dat2")
	//check(err)
	//defer f.Close()

	//f2 := []byte{115, 111, 109, 101, 10}
	//n2, err := f.Write(f2)
	//check(err)
	//fmt.Printf("wrote %d bytes to file\n", n2)
	//f.Sync()

	//passwordlisttitle := "Password List: "

	//fileHandle, _ := os.

	UserMap = make(map[string]int)
	//UserMap["Users Online: "] = assignid()
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

	//UserSlice = append(UserSlice, "Users Online: ")

	go mainServer.Run()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleConn(mainServer, conn, mainServer.Users)
	}
}
