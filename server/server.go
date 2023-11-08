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
	buf := make([]byte, 1024)
	for {
		n, err := user.Conn.Read(buf)
		if err != nil {
			fmt.Printf("user:%v read err:%v\n", user.Name, err)
			delete(OnlineUserMap, user.ID)
			BroadcastChan <- fmt.Sprintf("user:%v go offline", user.Name)
			return
		}
		msg := string(buf[:n])

		if len(msg) > 9 && msg[:8] == "--rename" { //--rename ikun
			BroadcastChan <- fmt.Sprintf("%v rename: %v", user.Name, msg[9:])
			user.Name = msg[9:]
			OnlineUserMap[user.ID] = user
		} else if msg == "--who" {
			users := make([]string, 0)
			for k := range OnlineUserMap {
				users = append(users, OnlineUserMap[k].Name)
			}
			msg = strings.Join(users, "\n")
			user.Msg <- msg
		} else {
			fmt.Printf("receive %v msg:%v\n", user.Name, msg)
			BroadcastChan <- fmt.Sprintf("%v: %v", user.Name, msg)
		}

	}

}
func ClientWrite(user User) {
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
