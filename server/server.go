package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type User struct {
	Name string
	Msg  chan string
	Conn net.Conn
}

var onlineUserMap = make(map[string]*User)
var broadcastChan = make(chan string)
var lock sync.RWMutex

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
			Name: conn.RemoteAddr().String(),
			Msg:  make(chan string),
			Conn: conn,
		}
		broadcastChan <- fmt.Sprintf("user:%v go online", user.Name)

		lock.Lock()
		onlineUserMap[user.Name] = &user
		lock.Unlock()

		go ClientRead(&user)
		go ClientWrite(&user)
	}
}
func Broadcast() {
	for {
		msg := <-broadcastChan
		lock.RLock()
		for k := range onlineUserMap {
			onlineUserMap[k].Msg <- msg
		}
		lock.RUnlock()
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
func ClientRead(user *User) {
	defer func() {
		user.Conn.Close()
	}()
	buf := make([]byte, 1024)
	for {
		n, err := user.Conn.Read(buf)
		if err != nil {
			fmt.Printf("user:%v read err:%v\n", user.Name, err)

			lock.Lock()
			delete(onlineUserMap, user.Name)
			lock.Unlock()

			broadcastChan <- fmt.Sprintf("user:%v go offline", user.Name)
			return
		}
		msg := string(buf[:n])

		if len(msg) > 9 && msg[:9] == "--rename=" { //--rename=ikun
			lock.Lock()
			_, ok := onlineUserMap[msg[9:]]
			if ok {
				user.Msg <- fmt.Sprintf("name: %v has existed !", msg[9:])
			} else {
				broadcastChan <- fmt.Sprintf("%v rename: %v", user.Name, msg[9:])
				delete(onlineUserMap, user.Name)
				user.Name = msg[9:]
				onlineUserMap[user.Name] = user //update name
			}
			lock.Unlock()

		} else if msg == "--who" { //--who
			users := make([]string, 0)

			lock.RLock()
			for k := range onlineUserMap {
				users = append(users, onlineUserMap[k].Name)
			}
			lock.RUnlock()

			msg = strings.Join(users, "\n")
			user.Msg <- msg
		} else if len(msg) > 7 && msg[:7] == "--name=" { //private --name=ikun hello world
			msg1 := msg[7:]                  //ikun hello world
			msg2 := strings.Split(msg1, " ") //ikun hello world

			lock.RLock()
			v, ok := onlineUserMap[msg2[0]] //ikun
			lock.RUnlock()

			if ok && len(msg2) > 1 {
				msg3 := strings.Join(msg2[1:], " ") //hello world
				v.Msg <- fmt.Sprintf("%v: %v", user.Name, msg3)
			} else {
				user.Msg <- fmt.Sprintf("name: %v not find or msg is null", msg2[0])
			}

		} else { //public
			fmt.Printf("receive %v msg:%v\n", user.Name, msg)
			broadcastChan <- fmt.Sprintf("%v: %v", user.Name, msg)
		}

	}

}
func ClientWrite(user *User) {
	defer func() {
		user.Conn.Close()
	}()
	for {
		msg := <-user.Msg
		_, err := user.Conn.Write([]byte(msg))
		if err != nil {
			fmt.Printf("user:%v write err:%v\n", user.Name, err)

			lock.Lock()
			delete(onlineUserMap, user.Name)
			lock.Unlock()

			broadcastChan <- fmt.Sprintf("user:%v go offline", user.Name)
			return
		}
	}
}
