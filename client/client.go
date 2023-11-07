package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	//1.connection
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatalf("dial err:%v\n", err)
	}
	defer conn.Close()

	go read(conn)
	write(conn)
}
func read(conn net.Conn) {
	// defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalf("read err:%v\n", err)
		}
		fmt.Printf("receive msg:%v\n", buf[:n])
	}

}
func write(conn net.Conn) {
	// defer conn.Close()
	inputStream := bufio.NewReader(os.Stdin)
	for {
		input, err := inputStream.ReadString('\n')
		if err != nil {
			log.Fatalf("input err:%v\n", err)
		}
		n, err := conn.Write([]byte(input))
		if err != nil {
			log.Fatalf("write err:%v\n", err)
		}
		fmt.Printf("write msg:%v\n", input[:n-2])
	}

}
