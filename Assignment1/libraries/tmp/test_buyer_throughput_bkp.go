package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"github.com/sirupsen/logrus"
)

const (
	ServiceName   = "test_latency"
	ServerHostEnv = "SERVER_HOST"
	ServerPortEnv = "SERVER_PORT"
)

var (
	err               error
	ctx               context.Context
	sessionID         string
	httpServerHost    = common.GetEnv(ServerHostEnv, "localhost")
	httpServerPort, _ = strconv.Atoi(common.GetEnv(ServerPortEnv, "50000"))
)

func initialBuyerExchange(conn net.Conn) {
	t := time.Now()
	myTime := t.Format(time.RFC3339Nano) + "\n"
	conn.Write([]byte("Hi Server. I am a buyer: " + myTime))
}

func handleConcurrentMessagesFromServer(conn net.Conn, responseChannel chan bool) {
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
		// Send a signal that the response is received
		responseChannel <- true
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

func main() {
	//log.Println("Initializing test buyer ...")
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", httpServerHost, httpServerPort))
	if err != nil {
		log.Fatal("Connection failed")
	}
	defer conn.Close()

	responseChannel := make(chan bool)

	initialBuyerExchange(conn)
	go handleConcurrentMessagesFromServer(conn, responseChannel)

	var iterations1 int64 = 10
	var iterations2 int64 = 1000
	var average float64 = 0
	for i := 0; i < int(iterations1); i++ {
		var buffer []byte
		start := time.Now()
		for j := 0; j < int(iterations2); j++ {
			if buffer, err = createLoginPayload(); err != nil {
				logrus.Error(err)
				break
			}
			//log.Println("Sending login buffer to server at ", time.Now().Format(time.RFC3339Nano))
			conn.Write(buffer)
			// wait for the response
			<-responseChannel
		}
		duration := time.Since(start)
		timeMillisecond := duration.Milliseconds()
		//fmt.Printf("%d\n\n", timeMillisecond)
		average += (float64(iterations2) / float64(timeMillisecond/1000))

	}
	defer conn.Close()

	fmt.Printf("%f\n", (average / float64(iterations1)))
	//log.Fatal("Closing connection. Exiting...")
}
