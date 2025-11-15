package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var config Config

func main() {
	configFile, _ := os.ReadFile("config.yml")
	yaml.Unmarshal(configFile, &config)
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", config.Server.Listen[0]))
	fmt.Printf("listening: %d", config.Server.Listen[0])
	for {
		conn, _ := listener.Accept()
		go tcpConnHandler(conn)
	}
}

func tcpConnHandler(conn net.Conn) {
	backendConn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Backends[0].Addr, config.Backends[0].Port))
	var wg sync.WaitGroup

	wg.Add(2)
	go tcpPipe(conn, backendConn, &wg)
	go tcpPipe(backendConn, conn, &wg)
	wg.Wait()
	backendConn.Close()
	conn.Close()
}

func tcpPipe(src, dst net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	io.Copy(dst, src)
	if tcp, ok := dst.(*net.TCPConn); ok {
		tcp.CloseWrite()
	}
}
