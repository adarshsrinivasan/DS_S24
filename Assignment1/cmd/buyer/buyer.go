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
	ServiceName       = "buyer"
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

func userBuyerOptions() []byte {
	menu := gocliselect.NewMenu("Welcome! \nSelect an option: ")

	menu.AddItem("Create an account", "0")
	menu.AddItem("Login", "1")
	menu.AddItem("Logout", "2")
	menu.AddItem("Search items for sale", "3")
	menu.AddItem("Add item to the shopping cart", "4")
	menu.AddItem("Remove item from the shopping cart", "5")
	menu.AddItem("Save the shopping cart", "6")
	menu.AddItem("Clear the shopping cart", "7")
	menu.AddItem("Display shopping cart", "8")
	menu.AddItem("Provide feedback", "9")
	menu.AddItem("Get seller rating", "10")
	menu.AddItem("Get buyer purchase history", "11")
	menu.AddItem("Exit", "12")

	choice := menu.Display()

	var body []byte
	reader := bufio.NewReader(os.Stdin)
	switch choice {
	case "0": // Create Account
		fmt.Println("Enter the name")
		name, _ := reader.ReadString('\n')
		fmt.Println("Enter the username")
		username, _ := reader.ReadString('\n')
		fmt.Println("Enter the password")
		pwd, _ := reader.ReadString('\n')
		body, _ = json.Marshal(&common.Credentials{
			Name:     name,
			UserName: username,
			Password: pwd,
		})
	case "1": // LOGIN
		fmt.Println("Enter the username")
		username, _ := reader.ReadString('\n')
		fmt.Println("Enter the Password")
		pwd, _ := reader.ReadString('\n')
		body, _ = json.Marshal(&common.Credentials{
			UserName: username,
			Password: pwd,
		})
	case "2": // LOGOUT
	case "3":
		// search items for sale
		fmt.Println("Enter item category")
		category, _ := reader.ReadString('\n')
		fmt.Println("Enter item keywords")
		keywords, _ := reader.ReadString('\n')
		body, _ = json.Marshal(&ProductModel{
			Category: StringToCategory[category],
			Keywords: strings.Split(keywords, ","),
		})
	case "4":
		// add item to shopping cart
		fmt.Println("Enter item id")
		id, _ := reader.ReadString('\n')

		fmt.Println("Enter item quantity")
		quantity, _ := reader.ReadString('\n')
		intQuantity, _ := strconv.ParseInt(strings.Split(strings.TrimSpace(quantity), "\n")[0], 10, 32)

		body, _ = json.Marshal(&ProductModel{
			ID:       strings.Split(strings.TrimSpace(id), "\n")[0],
			Quantity: int(intQuantity),
		})
	case "5":
		// remove item to shopping cart
		fmt.Println("Enter item id")
		id, _ := reader.ReadString('\n')

		fmt.Println("Enter item quantity")
		quantity, _ := reader.ReadString('\n')
		intQuantity, _ := strconv.ParseInt(strings.Split(strings.TrimSpace(quantity), "\n")[0], 10, 32)

		body, _ = json.Marshal(&ProductModel{
			ID:       strings.Split(strings.TrimSpace(id), "\n")[0],
			Quantity: int(intQuantity),
		})
	case "6":
		// save shopping cart
	case "7":
		// clear shopping cart
	case "8":
		// display shopping cart
	case "9":
		// provide feedback
		fmt.Println("Enter item id")
		id, _ := reader.ReadString('\n')

		fmt.Println("Enter item rating")
		rating, _ := reader.ReadString('\n')

		body, _ = json.Marshal(map[string]string{
			"productID": strings.Split(strings.TrimSpace(id), "\n")[0],
			"rating":    strings.Split(strings.TrimSpace(rating), "\n")[0],,
		})
	case "10":
		// get seller ratings
		fmt.Println("Enter item id")
		id, _ := reader.ReadString('\n')
		body, _ = json.Marshal(map[string]string{
			"productID": strings.Split(strings.TrimSpace(id), "\n")[0],
		})
	case "11": // get buyer purchase history
	case "12": //exit
		return nil
	}
	requestPayload := common.ClientRequest{
		SessionID: sessionID,
		Service:   choice,
		UserType:  common.Buyer,
		Body:      body,
	}

	payload := requestPayload.SerializeRequest()
	return payload
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
	log.Println("Initializing buyer buyer ...")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", httpServerHost, httpServerPort))
	if err != nil {
		log.Fatal("Connection failed")
	}
	defer conn.Close()

	initialBuyerExchange(conn)
	go handleConcurrentMessagesFromServer(conn)

	for {
		buffer := userBuyerOptions()
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
