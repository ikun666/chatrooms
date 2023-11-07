package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	//1.listen
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("listen err:%v\n", err)
	}
	defer listen.Close()
	for {
		//2.accept
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalf("accept err:%v\n", err)
		}
		go read(conn)
		// go write(conn)
	}
}

func read(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("read err:%v\n", err)
			return
		}
		fmt.Printf("receive msg:%v\n", string(buf[:n-2]))
	}

}
func write(conn net.Conn) {
	// defer conn.Close()
	// buf := make([]byte, 1024)
	buf := "ikun666" + conn.RemoteAddr().String()
	_, err := conn.Write([]byte(buf))
	if err != nil {
		fmt.Printf("write err:%v\n", err)
	}
}
