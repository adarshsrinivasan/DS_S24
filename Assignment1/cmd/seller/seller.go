package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"github.com/nexidian/gocliselect"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ServiceName       = "seller"
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

var ConditionToString = map[CONDITION]string{
	NEW:  "NEW",
	USED: "USED",
}

var StringToCondition = map[string]CONDITION{
	"NEW":  NEW,
	"USED": USED,
}

func initialSellerExchange(conn net.Conn) {
	t := time.Now()
	myTime := t.Format(time.RFC3339Nano) + "\n"
	conn.Write([]byte("Hi Server. I am a seller: " + myTime))
}

func userSellerOptions() []byte {
	menu := gocliselect.NewMenu("Welcome! \nSelect an option: ")

	menu.AddItem("Create an account", "0")
	menu.AddItem("Login", "1")
	menu.AddItem("Logout", "2")
	menu.AddItem("Get seller rating", "3")
	menu.AddItem("Put an item for sale", "4")
	menu.AddItem("Change the sale price of an item", "5")
	menu.AddItem("Remove an item from sale", "6")
	menu.AddItem("Display items currently on sale put up by this seller", "7")
	menu.AddItem("Exit", "8")

	choice := menu.Display()

	var body []byte
	reader := bufio.NewReader(os.Stdin)
	switch choice {
	case "0":
		// Create Account
		fmt.Println("Enter the name")
		name, _ := reader.ReadString('\n')
		fmt.Println("Enter the username")
		username, _ := reader.ReadString('\n')
		fmt.Println("Enter the Password")
		pwd, _ := reader.ReadString('\n')
		body, _ = json.Marshal(&common.Credentials{
			Name:     name,
			UserName: username,
			Password: pwd,
		})
	case "1":
		// Login
		fmt.Println("Enter the username")
		username, _ := reader.ReadString('\n')
		fmt.Println("Enter the Password")
		pwd, _ := reader.ReadString('\n')
		body, _ = json.Marshal(&common.Credentials{
			UserName: username,
			Password: pwd,
		})
	case "2": // logout
	case "3": // get seller ratings
	case "4": // put an item for sale
		fmt.Println("Enter item name")
		name, _ := reader.ReadString('\n')

		fmt.Println("Enter item category")
		category, _ := reader.ReadString('\n')

		fmt.Println("Enter item keywords")
		keywords, _ := reader.ReadString('\n')

		fmt.Println("Enter item condition")
		condition, _ := reader.ReadString('\n')

		fmt.Println("Enter item price")
		price, _ := reader.ReadString('\n')
		priceNew := strings.Split(strings.TrimSpace(price), "\n")[0]
		floatPrice, _ := strconv.ParseFloat(priceNew, 32)

		fmt.Println("Enter item quantity")
		quantity, _ := reader.ReadString('\n')
		intQuantity, _ := strconv.ParseInt(strings.Split(strings.TrimSpace(quantity), "\n")[0], 10, 32)

		body, _ = json.Marshal(&ProductModel{
			Name:      name,
			Category:  StringToCategory[category],
			Keywords:  strings.Split(keywords, ","),
			Condition: StringToCondition[condition],
			SalePrice: float32(floatPrice),
			Quantity:  int(intQuantity),
		})
	case "5": // change the sale price of an item
		fmt.Println("Enter item id")
		itemId, _ := reader.ReadString('\n')
		fmt.Println("Enter new sale price")
		price, _ := reader.ReadString('\n')
		priceNew := strings.Split(strings.TrimSpace(price), "\n")[0]
		floatPrice, _ := strconv.ParseFloat(priceNew, 32)

		body, _ = json.Marshal(&ProductModel{
			ID:        strings.Split(strings.TrimSpace(itemId), "\n")[0],
			SalePrice: float32(floatPrice),
		})
	case "6": // remove an item from sale
		fmt.Println("Enter item id")
		itemId, _ := reader.ReadString('\n')
		fmt.Println("Enter item quantity")
		quantity, _ := reader.ReadString('\n')
		intQuantity, _ := strconv.ParseInt(strings.Split(strings.TrimSpace(quantity), "\n")[0], 10, 32)

		body, _ = json.Marshal(&ProductModel{
			ID:       strings.Split(strings.TrimSpace(itemId), "\n")[0],
			Quantity: int(intQuantity),
		})
	case "7": // display
	case "8":
		return nil
	}
	requestPayload := common.ClientRequest{
		SessionID: sessionID,
		Service:   choice,
		UserType:  common.Seller,
		Body:      body,
	}

	payload := requestPayload.SerializeRequest()
	return payload
}

func handleConcurrentMessagesFromServer(conn net.Conn) {
	defer conn.Close()
	for {
		requestBody := make([]byte, 5000)
		var response common.ClientResponse
		if _, err := conn.Read(requestBody); err != nil {
			return
		}
		response.DeserializeRequest(requestBody)

		if response.SessionID != "" {
			sessionID = response.SessionID
		}
		response.LogResponse()

		if strings.HasPrefix(response.Message, "Timeout: ") {
			log.Fatal(response.Message)
		}
	}
}

func main() {
	log.Println("Initializing seller buyer ...")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", httpServerHost, httpServerPort))
	if err != nil {
		log.Fatal("Connection failed")
	}
	defer conn.Close()

	initialSellerExchange(conn)
	go handleConcurrentMessagesFromServer(conn)

	for {
		buffer := userSellerOptions()
		if buffer == nil {
			break
		} else {
			log.Println("Sending buffer to server at ", time.Now().Format(time.RFC3339Nano))
			conn.Write(buffer)
		}

		time.Sleep(2 * time.Second)
	}
	defer conn.Close()

	log.Fatal("Closing connection. Exiting...")
}
