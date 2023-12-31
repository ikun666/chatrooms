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
		fmt.Printf("dial err:%v\n", err)
		return
	}
	defer conn.Close()

	go read(conn)
	write(conn)
}
func read(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalf("read err:%v\n", err)
		}
		fmt.Printf("%v\n", string(buf[:n]))
	}

}
func write(conn net.Conn) {
	inputStream := bufio.NewReader(os.Stdin)
	for {
		input, err := inputStream.ReadString('\n')
		if err != nil {
			log.Fatalf("input err:%v\n", err)
		}
		msg := []byte(input)
		//delete /r/n
		_, err = conn.Write(msg[:len(msg)-2])
		if err != nil {
			log.Fatalf("write err:%v\n", err)
		}
		// fmt.Printf("write msg:%v\n", string(input[:n-2]))
	}

}
