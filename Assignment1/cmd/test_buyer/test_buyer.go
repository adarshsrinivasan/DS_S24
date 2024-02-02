package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ServiceName       = "test_buyer"
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

type ProductModel struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	Category  CATEGORY  `json:"category,omitempty" bson:"category,omitempty"`
	Keywords  []string  `json:"keywords,omitempty" bson:"keywords,omitempty"`
	Condition CONDITION `json:"condition,omitempty" bson:"condition,omitempty"`
	SalePrice float32   `json:"salePrice,omitempty" bson:"salePrice,omitempty"`
	Quantity  int       `json:"quantity,omitempty" bson:"quantity"`
}

type CONDITION int

const (
	NEW CONDITION = iota
	USED
)

type CATEGORY int

const (
	ZERO CATEGORY = iota
	ONE
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
)

var CategoryToString = map[CATEGORY]string{
	ZERO:  "ZERO",
	ONE:   "ONE",
	TWO:   "TWO",
	THREE: "THREE",
	FOUR:  "FOUR",
	FIVE:  "FIVE",
	SIX:   "SIX",
	SEVEN: "SEVEN",
	EIGHT: "EIGHT",
	NINE:  "NINE",
}

var StringToCategory = map[string]CATEGORY{
	"ZERO":  ZERO,
	"ONE":   ONE,
	"TWO":   TWO,
	"THREE": THREE,
	"FOUR":  FOUR,
	"FIVE":  FIVE,
	"SIX":   SIX,
	"SEVEN": SEVEN,
	"EIGHT": EIGHT,
	"NINE":  NINE,
}

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
		response.LogResponse()

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

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter server host")
	httpServerHost, _ := common.ReadTrimString(reader)
	fmt.Println("Enter server port")
	httpServerPortString, _ := common.ReadTrimString(reader)
	httpServerPort, _ := strconv.Atoi(httpServerPortString)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", httpServerHost, httpServerPort))
	if err != nil {
		log.Fatal("Connection failed")
	}
	defer conn.Close()

	initialBuyerExchange(conn)
	go handleConcurrentMessagesFromServer(conn)

	fmt.Println("Enter the no of iterations")
	iterationsString, _ := common.ReadTrimString(reader)
	iterations, _ := strconv.Atoi(iterationsString)

	for i := 0; i < iterations; i++ {
		var buffer []byte
		if buffer, err = createLoginPayload(); err != nil {
			logrus.Error(err)
			break
		}

		log.Println("Sending login buffer to server at ", time.Now().Format(time.RFC3339Nano))
		conn.Write(buffer)

		time.Sleep(1 * time.Second)

		if buffer, err = createLogoutPayload(); err != nil {
			logrus.Error(err)
			break
		}

		log.Println("Sending logout buffer to server at ", time.Now().Format(time.RFC3339Nano))
		conn.Write(buffer)

		time.Sleep(1 * time.Second)

	}
	defer conn.Close()

	log.Fatal("Closing connection. Exiting...")
}
