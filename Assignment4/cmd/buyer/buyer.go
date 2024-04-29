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
	ServiceName   = "buyer"
	ServerHostEnv = "SERVER_HOST"
	ServerPortEnv = "SERVER_PORT"
)

var (
	err               error
	ctx               context.Context
	sessionID         string
	httpServerHost    = common.GetEnv(ServerHostEnv, "localhost")
	httpServerPort, _ = strconv.Atoi(common.GetEnv(ServerPortEnv, "50020"))
	baseURL           = fmt.Sprintf("http://%s:%d/api/v1/marketplace/%s", httpServerHost, httpServerPort, ServiceName)
)

func userBuyerOptions() error {
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
	menu.AddItem("Make Purchase", "9")
	menu.AddItem("Provide feedback", "10")
	menu.AddItem("Get seller rating", "11")
	menu.AddItem("Get buyer purchase history", "12")
	menu.AddItem("Exit", "13")

	choice := menu.Display()

	reader := bufio.NewReader(os.Stdin)
	switch choice {
	case "0":
		// Create Account
		fmt.Println("Enter the name")
		name, _ := common.ReadTrimString(reader)
		fmt.Println("Enter the username")
		username, _ := common.ReadTrimString(reader)
		fmt.Println("Enter the password")
		password, _ := common.ReadTrimString(reader)
		request := BuyerModel{
			Name:     name,
			UserName: username,
			Password: password,
		}
		url := fmt.Sprintf("%s/create", baseURL)
		_, err := common.MakeHTTPRequest[BuyerModel, BuyerModel](ctx, "POST", url, "", request, true)
		if err != nil {
			return err
		}
	case "1":
		// LOGIN
		fmt.Println("Enter the username")
		username, _ := common.ReadTrimString(reader)
		fmt.Println("Enter the Password")
		password, _ := common.ReadTrimString(reader)
		request := BuyerModel{
			UserName: username,
			Password: password,
		}
		url := fmt.Sprintf("%s/login", baseURL)
		resp, err := common.MakeHTTPRequest[BuyerModel, SessionModel](ctx, "POST", url, "", request, true)
		if err != nil {
			return err
		}
		sessionID = resp.SessionID
	case "2":
		// LOGOUT
		request := BuyerModel{}
		url := fmt.Sprintf("%s/logout", baseURL)
		_, err := common.MakeHTTPRequest[BuyerModel, string](ctx, "POST", url, sessionID, request, true)
		if err != nil {
			return err
		}
	case "3":
		// search items for sale
		fmt.Println("Enter item category")
		category, _ := common.ReadTrimString(reader)
		fmt.Println("Enter item keywords")
		keywords, _ := common.ReadTrimString(reader)

		request := ProductModel{
			Category: StringToCategory[category],
			Keywords: strings.Split(keywords, ","),
		}
		url := fmt.Sprintf("%s/searchItems", baseURL)
		_, err := common.MakeHTTPRequest[ProductModel, []ProductModel](ctx, "POST", url, sessionID, request, true)
		if err != nil {
			return err
		}
	case "4":
		// add item to shopping cart
		fmt.Println("Enter item id")
		id, _ := common.ReadTrimString(reader)

		fmt.Println("Enter item quantity")
		quantity, _ := common.ReadTrimString(reader)
		intQuantity, _ := strconv.ParseInt(quantity, 10, 32)

		request := ProductModel{
			ID:       id,
			Quantity: int(intQuantity),
		}
		url := fmt.Sprintf("%s/addItemToCart", baseURL)
		_, err := common.MakeHTTPRequest[ProductModel, string](ctx, "POST", url, sessionID, request, true)
		if err != nil {
			return err
		}

	case "5":
		// remove item to shopping cart
		fmt.Println("Enter item id")
		id, _ := common.ReadTrimString(reader)

		fmt.Println("Enter item quantity")
		quantity, _ := common.ReadTrimString(reader)
		intQuantity, _ := strconv.ParseInt(quantity, 10, 32)

		request := ProductModel{
			ID:       id,
			Quantity: int(intQuantity),
		}
		url := fmt.Sprintf("%s/removeItemFromCart", baseURL)
		_, err := common.MakeHTTPRequest[ProductModel, string](ctx, "POST", url, sessionID, request, true)
		if err != nil {
			return err
		}
	case "6":
		// save cart
		request := ProductModel{}
		url := fmt.Sprintf("%s/saveCart", baseURL)
		_, err := common.MakeHTTPRequest[ProductModel, string](ctx, "POST", url, sessionID, request, true)
		if err != nil {
			return err
		}
	case "7":
		// clear shopping cart
		request := ProductModel{}
		url := fmt.Sprintf("%s/clearCart", baseURL)
		_, err := common.MakeHTTPRequest[ProductModel, string](ctx, "POST", url, sessionID, request, true)
		if err != nil {
			return err
		}
	case "8":
		// display shopping cart
		url := fmt.Sprintf("%s/getCart", baseURL)
		_, err := common.MakeHTTPRequest[any, CartModel](ctx, "GET", url, sessionID, nil, true)
		if err != nil {
			return err
		}
	case "9":
		// Make Purchase
		fmt.Println("Enter name")
		name, _ := common.ReadTrimString(reader)
		fmt.Println("Enter card number")
		cardNumber, _ := common.ReadTrimString(reader)
		fmt.Println("Enter expiry date (MM/YYYY)")
		expiry, _ := common.ReadTrimString(reader)
		request := PurchaseDetailsModel{
			Name:             name,
			CreditCardNumber: cardNumber,
			Expiry:           expiry,
		}
		url := fmt.Sprintf("%s/makePurchase", baseURL)
		_, err := common.MakeHTTPRequest[PurchaseDetailsModel, string](ctx, "POST", url, sessionID, request, true)
		if err != nil {
			return err
		}

	case "10":
		// provide feedback
		fmt.Println("Enter item id")
		id, _ := common.ReadTrimString(reader)

		fmt.Println("Positive review?(Y/N)")
		rating, _ := common.ReadTrimString(reader)
		if rating != "Y" && rating != "N" {
			fmt.Printf("Wrong review response: %s. Expecting Y or N.", rating)
			return nil
		}
		request := SellerModel{}
		url := fmt.Sprintf("%s/feedback/%s", baseURL, id)
		switch rating {
		case "Y":
			url = fmt.Sprintf("%s/liked", url)

		case "N":
			url = fmt.Sprintf("%s/disliked", url)
		default:
			fmt.Printf("Wrong review response: %s. Expecting Y or N.\n", rating)
			return nil
		}
		_, err := common.MakeHTTPRequest[SellerModel, string](ctx, "POST", url, sessionID, request, true)
		if err != nil {
			return err
		}
	case "11":
		// get seller ratings
		fmt.Println("Enter item id")
		id, _ := common.ReadTrimString(reader)
		url := fmt.Sprintf("%s/getSellerRating/%s", baseURL, id)
		_, err := common.MakeHTTPRequest[any, SellerModel](ctx, "GET", url, sessionID, nil, true)
		if err != nil {
			return err
		}
	case "12":
		// get buyer purchase history
		url := fmt.Sprintf("%s/getPurchaseHistory", baseURL)
		_, err := common.MakeHTTPRequest[any, []TransactionModel](ctx, "GET", url, sessionID, nil, true)
		if err != nil {
			return err
		}
	case "13": //exit
		return nil
	default:
		fmt.Printf("unknown choice: %s\n", choice)
		return nil
	}
	return nil
}

func main() {
	log.Println("Initializing buyer buyer ...")
	ctx = context.Background()
	for {
		if err = userBuyerOptions(); err != nil {
			logrus.Error(err)
			break
		}

		time.Sleep(1 * time.Second)
	}

	log.Fatal("Closing connection. Exiting...")
}
