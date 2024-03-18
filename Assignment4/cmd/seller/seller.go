package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/nexidian/gocliselect"
	"github.com/sirupsen/logrus"
)

const (
	ServiceName   = "seller"
	ServerHostEnv = "SERVER_HOST"
	ServerPortEnv = "SERVER_PORT"
)

var (
	err               error
	ctx               context.Context
	sessionID         string
	httpServerHost    = common.GetEnv(ServerHostEnv, "localhost")
	httpServerPort, _ = strconv.Atoi(common.GetEnv(ServerPortEnv, "50000"))
	baseURL           = fmt.Sprintf("http://%s:%d/api/v1/marketplace/%s", httpServerHost, httpServerPort, ServiceName)
)

func userSellerOptions() error {
	menu := gocliselect.NewMenu("Welcome! \nSelect an option: ")

	menu.AddItem("Create an account", "0")
	menu.AddItem("Login", "1")
	menu.AddItem("Logout", "2")
	menu.AddItem("Get seller rating", "3")
	menu.AddItem("Put an item for sale", "4")
	menu.AddItem("Change the sale price of an item", "5")
	menu.AddItem("Remove an item from sale", "6")
	menu.AddItem("Display items currently on sale put up by this seller", "7")
	menu.AddItem("Display sold items by this seller", "8")
	menu.AddItem("Exit", "9")

	choice := menu.Display()

	reader := bufio.NewReader(os.Stdin)
	switch choice {
	case "0":
		// Create Account
		fmt.Println("Enter the name")
		name, _ := common.ReadTrimString(reader)
		fmt.Println("Enter the username")
		username, _ := common.ReadTrimString(reader)
		fmt.Println("Enter the Password")
		password, _ := common.ReadTrimString(reader)
		request := SellerModel{
			Name:     name,
			UserName: username,
			Password: password,
		}
		url := fmt.Sprintf("%s/create", baseURL)
		_, err := common.MakeHTTPRequest[SellerModel, SellerModel](ctx, "POST", url, "", request, true)
		if err != nil {
			return err
		}

	case "1":
		// Login
		fmt.Println("Enter the username")
		username, _ := common.ReadTrimString(reader)
		fmt.Println("Enter the Password")
		password, _ := common.ReadTrimString(reader)
		request := SellerModel{
			UserName: username,
			Password: password,
		}
		url := fmt.Sprintf("%s/login", baseURL)
		resp, err := common.MakeHTTPRequest[SellerModel, SessionModel](ctx, "POST", url, "", request, true)
		if err != nil {
			return err
		}
		sessionID = resp.SessionID
	case "2":
		// Logout
		request := SellerModel{}
		url := fmt.Sprintf("%s/logout", baseURL)
		_, err := common.MakeHTTPRequest[SellerModel, string](ctx, "POST", url, sessionID, request, true)
		if err != nil {
			return err
		}
	case "3":
		// get seller ratings
		url := fmt.Sprintf("%s/getRating", baseURL)
		_, err := common.MakeHTTPRequest[any, SellerModel](ctx, "GET", url, sessionID, nil, true)
		if err != nil {
			return err
		}
	case "4": // put an item for sale
		fmt.Println("Enter item name")
		name, _ := common.ReadTrimString(reader)

		fmt.Println("Enter item category")
		category, _ := common.ReadTrimString(reader)

		fmt.Println("Enter item keywords")
		keywords, _ := common.ReadTrimString(reader)

		fmt.Println("Enter item condition")
		condition, _ := common.ReadTrimString(reader)

		fmt.Println("Enter item price")
		price, _ := common.ReadTrimString(reader)
		priceNew := strings.Split(strings.TrimSpace(price), "\n")[0]
		floatPrice, _ := strconv.ParseFloat(priceNew, 32)

		fmt.Println("Enter item quantity")
		quantity, _ := common.ReadTrimString(reader)
		intQuantity, _ := strconv.ParseInt(strings.Split(strings.TrimSpace(quantity), "\n")[0], 10, 32)

		request := ProductModel{
			Name:      name,
			Category:  StringToCategory[category],
			Keywords:  strings.Split(keywords, ","),
			Condition: StringToCondition[condition],
			SalePrice: float32(floatPrice),
			Quantity:  int(intQuantity),
		}
		url := fmt.Sprintf("%s/createItem", baseURL)
		_, err := common.MakeHTTPRequest[ProductModel, ProductModel](ctx, "POST", url, sessionID, request, true)
		if err != nil {
			return err
		}
	case "5": // change the sale price of an item
		fmt.Println("Enter item id")
		itemId, _ := common.ReadTrimString(reader)
		fmt.Println("Enter new sale price")
		price, _ := common.ReadTrimString(reader)
		floatPrice, _ := strconv.ParseFloat(price, 32)

		request := ProductModel{
			ID:        itemId,
			SalePrice: float32(floatPrice),
		}
		url := fmt.Sprintf("%s/updateItemSalePrice", baseURL)
		_, err := common.MakeHTTPRequest[ProductModel, ProductModel](ctx, "PUT", url, sessionID, request, true)
		if err != nil {
			return err
		}
	case "6": // remove an item from sale
		fmt.Println("Enter item id")
		itemId, _ := common.ReadTrimString(reader)
		fmt.Println("Enter item quantity")
		quantity, _ := common.ReadTrimString(reader)
		intQuantity, _ := strconv.ParseInt(quantity, 10, 32)

		request := ProductModel{
			ID:       itemId,
			Quantity: int(intQuantity),
		}
		url := fmt.Sprintf("%s/removeItem", baseURL)
		_, err := common.MakeHTTPRequest[ProductModel, ProductModel](ctx, "PUT", url, sessionID, request, true)
		if err != nil {
			return err
		}

	case "7":
		// get items
		url := fmt.Sprintf("%s/getItems", baseURL)
		_, err := common.MakeHTTPRequest[any, []ProductModel](ctx, "GET", url, sessionID, nil, true)
		if err != nil {
			return err
		}
	case "8":
		// get sold items
		url := fmt.Sprintf("%s/getSoldItems", baseURL)
		_, err := common.MakeHTTPRequest[any, []TransactionModel](ctx, "GET", url, sessionID, nil, true)
		if err != nil {
			return err
		}

	case "9":
		return nil
	default:
		fmt.Printf("unknown choice: %s\n", choice)
		return nil
	}
	return nil
}

func main() {
	log.Println("Initializing seller buyer ...")
	ctx = context.Background()
	for {
		if err = userSellerOptions(); err != nil {
			logrus.Error(err)
			break
		}
		time.Sleep(1 * time.Second)
	}

	log.Fatal("Closing connection. Exiting...")
}
