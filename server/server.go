package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type User struct {
	ID   string
	Name string
	Msg  chan string
	Conn net.Conn
}

var OnlineUserMap = make(map[string]User)
var UserName2IdMap = make(map[string]string)
var BroadcastChan = make(chan string)

func main() {
	//1.listen
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("listen err:%v\n", err)
	}
	defer listen.Close()
	go Broadcast()
	for {
		//2.accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept err:%v\n", err)
			continue
		}
		user := User{
			ID:   conn.RemoteAddr().String(),
			Name: conn.RemoteAddr().String(),
			Msg:  make(chan string),
			Conn: conn,
		}
		BroadcastChan <- fmt.Sprintf("user:%v go online", user.Name)
		OnlineUserMap[user.ID] = user
		UserName2IdMap[user.Name] = user.ID
		go ClientRead(user)
		go ClientWrite(user)
	}
}
func Broadcast() {
	for {
		msg := <-BroadcastChan
		for k := range OnlineUserMap {
			OnlineUserMap[k].Msg <- msg
		}
	}
}

// func ServerRead() {

// 	buf := make([]byte, 1024)
// 	for {
// 		n, err := conn.Read(buf)
// 		if err != nil {
// 			fmt.Printf("read err:%v\n", err)
// 			return
// 		}
// 		fmt.Printf("receive msg:%v\n", string(buf[:n-2]))
// 		Broadcast <- string(buf[:n-2])
// 	}

// }
//
//	func ServerWrite() {
//		for{
//			msg:=<-Broadcast
//			_, err := conn.Write([]byte(msg))
//			if err != nil {
//				fmt.Printf("write err:%v\n", err)
//			}
//		}
//	}
func ClientRead(user User) {
	defer user.Conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := user.Conn.Read(buf)
		if err != nil {
			fmt.Printf("user:%v read err:%v\n", user.Name, err)
			delete(OnlineUserMap, user.ID)
			delete(UserName2IdMap, user.Name)
			BroadcastChan <- fmt.Sprintf("user:%v go offline", user.Name)
			return
		}
		msg := string(buf[:n])

		if len(msg) > 9 && msg[:9] == "--rename=" { //--rename=ikun
			BroadcastChan <- fmt.Sprintf("%v rename: %v", user.Name, msg[9:])
			user.Name = msg[9:]
			OnlineUserMap[user.ID] = user
			delete(UserName2IdMap, user.Name)
			UserName2IdMap[user.Name] = user.ID
		} else if msg == "--who" { //--who
			users := make([]string, 0)
			for k := range OnlineUserMap {
				users = append(users, OnlineUserMap[k].Name)
			}
			msg = strings.Join(users, "\n")
			user.Msg <- msg
		} else if len(msg) > 7 && msg[:7] == "--name=" { //private --name=ikun hello world
			msg1 := msg[7:]                  //ikun hello world
			msg2 := strings.Split(msg1, " ") //ikun hello world

			v, ok := UserName2IdMap[msg2[0]] //ikun
			if ok && len(msg2) > 1 {
				msg3 := strings.Join(msg2[1:], " ") //hello world
				OnlineUserMap[v].Msg <- fmt.Sprintf("%v: %v", user.Name, msg3)
			} else {
				user.Msg <- fmt.Sprintf("name: %v not find or msg is null", msg2[0])
			}

		} else { //public
			fmt.Printf("receive %v msg:%v\n", user.Name, msg)
			BroadcastChan <- fmt.Sprintf("%v: %v", user.Name, msg)
		}

	}

}
func ClientWrite(user User) {
	defer user.Conn.Close()
	for {
		msg := <-user.Msg
		_, err := user.Conn.Write([]byte(msg))
		if err != nil {
			fmt.Printf("user:%v write err:%v\n", user.Name, err)
			delete(OnlineUserMap, user.ID)
			BroadcastChan <- fmt.Sprintf("user:%v go offline", user.Name)
			return
		}
	}
}
