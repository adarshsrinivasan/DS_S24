package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/sirupsen/logrus"
)

const (
	ServiceName   = "test-latency"
	ServerHostEnv = "SERVER_HOST"
	ServerPortEnv = "SERVER_PORT"
)

var (
	err               error
	httpServerHost    = common.GetEnv(ServerHostEnv, "localhost")
	httpServerPort, _ = strconv.Atoi(common.GetEnv(ServerPortEnv, "50000"))
)

func initialBuyerExchange(conn net.Conn) {
	t := time.Now()
	myTime := t.Format(time.RFC3339Nano) + "\n"
	conn.Write([]byte("Hi Server. I am a buyer: " + myTime))
}

func readResponse(serverReader *bufio.Reader) {
	requestBody := make([]byte, 5000)
	var response common.ClientResponse
	if _, err := serverReader.Read(requestBody); err != nil {
		logrus.Errorf("readResponse read exception: %v", err)
		return
	}
	err := response.DeserializeRequest(requestBody)
	if err != nil {
		logrus.Errorf("readResponse exception: %v", err)
		return
	}

	//response.LogResponse()

	if strings.HasPrefix(response.Message, "Timeout: ") {
		log.Fatal(response.Message)
	}
}

func createLoginPayload(user int) ([]byte, error) {
	var payload []byte
	if payload, err = json.Marshal(&common.Credentials{
		UserName: "admin",
		Password: "admin",
	}); err != nil {
		logrus.Errorf("createLoginPayload: exception when trying to create login payload: %v", err)
		return nil, err
	}

	requestPayload := common.ClientRequest{
		SessionID: "",
		Service:   "1",
		Body:      payload,
	}
	switch user {
	case 1:
		requestPayload.UserType = common.SELLER
	default:
		requestPayload.UserType = common.BUYER
	}
	var serializedPayload []byte
	if serializedPayload, err = requestPayload.SerializeRequest(); err != nil {
		logrus.Errorf("createLoginPayload: exception when trying to create login payload: %v", err)
		return nil, err
	}
	return serializedPayload, nil
}

func runOperation(file *os.File, wg *sync.WaitGroup, thread, user int) {
	// Decrement the counter when the goroutine completes.
	defer wg.Done()
	fmt.Printf("Thread %d: Start\n", thread)
	var (
		iterations int64   = 10
		average    float64 = 0
		conn       net.Conn
		err        error
	)
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", httpServerHost, httpServerPort))
	if err != nil {
		log.Fatal("Connection failed")
	}
	defer conn.Close()
	initialBuyerExchange(conn)
	serverReader := bufio.NewReader(conn)
	var buffer []byte
	if buffer, err = createLoginPayload(user); err != nil {
		logrus.Error(err)
		return
	}
	for i := 0; i < int(iterations); i++ {
		start := time.Now()
		conn.Write(buffer)
		readResponse(serverReader)
		duration := time.Since(start)
		average += float64(duration.Milliseconds())
	}
	if _, err := file.WriteString(fmt.Sprintf("%f\n", (average / float64(iterations)))); err != nil {
		fmt.Printf("Error while writing to file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Thread %d: End\n", thread)
}

func main() {
	// Check if the user has provided the number of threads
	if len(os.Args) < 3 {
		fmt.Println("Usage:", os.Args[0], "<thread_count> <user-{0: Buyer, 1: Seller}>")
		os.Exit(1)
	}

	threadCount := 10
	threadCount, _ = strconv.Atoi(os.Args[1])

	user := 0
	user, _ = strconv.Atoi(os.Args[2])

	filePath := fmt.Sprintf("latency-test-%d-%d-%s", user, threadCount, time.Now().Format(time.RFC3339))

	// Write output to the file
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file for writing: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	var wg sync.WaitGroup
	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		go runOperation(f, &wg, i, user)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
}
