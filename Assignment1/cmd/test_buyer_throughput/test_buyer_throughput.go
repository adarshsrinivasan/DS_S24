package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	ServiceName       = "test_buyer_response"
	HttpServerHostEnv = "HTTP_SERVER_HOST"
	HttpServerPortEnv = "HTTP_SERVER_PORT"
)

var (
	err               error
	ctx               context.Context
	sessionID         string
	httpServerHost    = common.GetEnv(HttpServerHostEnv, "localhost")
	httpServerPort, _ = strconv.Atoi(common.GetEnv(HttpServerPortEnv, "50000"))
)

func initialBuyerExchange(conn net.Conn) {
	t := time.Now()
	myTime := t.Format(time.RFC3339Nano) + "\n"
	conn.Write([]byte("Hi Server. I am a buyer: " + myTime))
}

func handleConcurrentMessagesFromServer(conn net.Conn) {
	defer conn.Close()
	for {
		requestBody := make([]byte, 5000)
		var response common.ClientResponse
		if _, err := conn.Read(requestBody); err != nil {
			return
		}
		err := response.DeserializeRequest(requestBody)
		if err != nil {
			logrus.Errorf("handleConcurrentMessagesFromServer exception: %v", err)
			return
		}

		if response.SessionID != "" {
			sessionID = response.SessionID
		} else {
			sessionID = ""
		}
		//response.LogResponse()

		if strings.HasPrefix(response.Message, "Timeout: ") {
			log.Fatal(response.Message)
		}
	}
}

func createLoginPayload() ([]byte, error) {
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
		UserType:  common.Buyer,
		Body:      payload,
	}
	var serializedPayload []byte
	if serializedPayload, err = requestPayload.SerializeRequest(); err != nil {
		logrus.Errorf("createLoginPayload: exception when trying to create login payload: %v", err)
		return nil, err
	}
	return serializedPayload, nil
}

func createLogoutPayload() ([]byte, error) {
	var payload []byte
	requestPayload := common.ClientRequest{
		SessionID: sessionID,
		Service:   "2",
		UserType:  common.Buyer,
		Body:      payload,
	}
	var serializedPayload []byte
	if serializedPayload, err = requestPayload.SerializeRequest(); err != nil {
		logrus.Errorf("createLogoutPayload: exception when trying to create logout payload: %v", err)
		return nil, err
	}
	return serializedPayload, nil
}

func main() {
	log.Println("Initializing test buyer ...")
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", httpServerHost, httpServerPort))
	if err != nil {
		log.Fatal("Connection failed")
	}
	defer conn.Close()

	initialBuyerExchange(conn)
	go handleConcurrentMessagesFromServer(conn)

	var iterations int64 = 10
	var average int64 = 0
	for i := 0; i < int(iterations); i++ {
		var buffer []byte
		start := time.Now()
		duration := time.Since(start)
		for j := 0; j < 1000; j++ {
			if buffer, err = createLoginPayload(); err != nil {
				logrus.Error(err)
				break
			}
			//log.Println("Sending login buffer to server at ", time.Now().Format(time.RFC3339Nano))
			conn.Write(buffer)
			time.Sleep(1 * time.Second)
		}
		timeNanoSeconds := duration.Nanoseconds() - 1000000000
		fmt.Printf("%f\n\n", timeNanoSeconds)
		average += timeNanoSeconds

	}
	defer conn.Close()

	fmt.Printf("%f\n", float64(average/iterations))
	log.Fatal("Closing connection. Exiting...")
}
