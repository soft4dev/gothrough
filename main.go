package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

func main() {
	var config Config
	configFile, _ := os.ReadFile("config.yml")
	yaml.Unmarshal(configFile, &config)
	var mainWaitGroup sync.WaitGroup
	mainWaitGroup.Add(1)
	go L4Handler(config)
	mainWaitGroup.Wait()
}

func L4Handler(config Config) {
	// ---------- PRINT TABLE ----------
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "LISTEN PORT\tBACKENDS")
	fmt.Fprintln(w, "-----------\t--------")

	for _, server := range config.L4.Servers {

		var backendList string
		for i, b := range server.Backends {
			if i > 0 {
				backendList += ", "
			}
			backendList += fmt.Sprintf("%s:%d", b.Addr, b.Port)
		}

		fmt.Fprintf(w, "%d\t%s\n", server.Listen, backendList)
	}

	w.Flush()
	// ----------------------------------
	var listenersWaitGroup sync.WaitGroup
	for _, server := range config.L4.Servers {
		listenersWaitGroup.Add(1)
		go L4Listener(&listenersWaitGroup, server)
	}
	listenersWaitGroup.Wait()
}

func L4Listener(wg *sync.WaitGroup, serverConfig Server) {
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", serverConfig.Listen))
	for {
		conn, _ := listener.Accept()
		go tcpConnHandler(conn, serverConfig.Backends)
	}
}

func tcpConnHandler(conn net.Conn, configBackends []Backend) {
	// TODO update to round robin requests to backends
	// TODO health check and forward to healthy backend
	backendConn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", configBackends[0].Addr, configBackends[0].Port))
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
