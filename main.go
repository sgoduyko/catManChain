package main

import (
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func getLogFile(logFileName string) *os.File {
	err := os.MkdirAll("temp", 0755)
	if err != nil {
		log.Fatalf("Didn't create path with err %s", err)
		return nil
	}
	file, err := os.OpenFile(fmt.Sprintf("temp/%s", logFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Didn't open log file with err %s", err)
	}
	return file
}

func cleanupLogFile(logFile *os.File, logFileName string) {
	sourceName := fmt.Sprintf("temp/%s", logFileName)
	closedNodeName := fmt.Sprintf("temp/%s%s", "closed_node_", logFileName)
	err := os.Rename(sourceName, closedNodeName)
	if err != nil {
		log.Fatalf("We got err when try to rename log file %s", err)
	}
	logFile.Close()
}

func main() {

	port := flag.Int("port", 8083, "Port for node")
	//nodeName := flag.String("nodeName", "jghmajvhkajhfksnaskd", "Node name for generate public key")
	flag.Parse()

	h, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *port)))
	if err != nil {
		log.Fatalf("Not created Node with err %s", err)
		return
	}
	defer h.Close()
	log.Printf("Node with peer id=\"%s\" and listening on addr=\"%s\"", h.ID(), h.Addrs())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	logFileName := fmt.Sprintf("%s.log", h.ID())
	logFile := getLogFile(logFileName)
	if logFile == nil {
		log.Fatal("Log file didn't get")
		return
	}
	log.SetOutput(logFile)

	defer func() {
		// catch panic if we have
		if r := recover(); r != nil {
			log.Printf("Panic err \"%v\"", r)
		}
	}()
	defer cleanupLogFile(logFile, logFileName)

	done := make(chan bool, 1)

	go func() {
		sig := <-signalChan
		log.Printf("Got sys signal %d", sig)
		done <- true
	}()
	<-done
	fmt.Printf("node with peer id=\"%s\" shutdown", h.ID())
}
