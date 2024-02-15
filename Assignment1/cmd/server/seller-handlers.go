package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/adarshsrinivasan/DS_S24/library/common"
)

func sellerCreateAccountHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {

	var sellerModel SellerModel
	if err := json.Unmarshal(r.Body, &sellerModel); err != nil {
		common.TCPRespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("sellerCreateAccountHandler: exception while parsing request. %v", err))
		return
	}

	if createdObj, statusCode, err := createSellerAccount(ctx, &sellerModel); err != nil {
		common.TCPRespondWithError(conn, statusCode, fmt.Sprintf("sellerCreateAccountHandler: exception while creating user. %v", err))
		return
	} else {
		common.TCPRespondWithJSON(conn, http.StatusCreated, r.SessionID, createdObj)
	}
}

func sellerLoginHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {

	var sellerModel SellerModel
	if err := json.Unmarshal(r.Body, &sellerModel); err != nil {
		common.TCPRespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("sellerLoginHandler: exception while parsing request. %v", err))
		return
	}

	if session, statusCode, err := sellerLogin(ctx, sellerModel.UserName, sellerModel.Password); err != nil {
		common.TCPRespondWithError(conn, statusCode, fmt.Sprintf("sellerLoginHandler: exception while Logging in user. %v", err))
		return
	} else {
		common.TCPRespondWithJSON(conn, http.StatusCreated, session, map[string]string{"sessionID": session})
	}
}

func sellerLogoutHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {

	if !validateSessionID(r.SessionID) {
		common.TCPRespondWithError(conn, http.StatusForbidden, fmt.Sprintf("sellerLogoutHandler: Invalid session. Please login again"))
		return
	}

	if statusCode, err := sellerLogout(ctx, r.SessionID); err != nil {
		common.TCPRespondWithError(conn, statusCode, fmt.Sprintf("sellerLogoutHandler: exception while Logging out user. %v", err))
		return
	} else {
		common.TCPRespondWithStatusCode(conn, http.StatusOK, r.SessionID)
	}
}

func sellerGetRatingHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {

	if !validateSessionID(r.SessionID) {
		common.TCPRespondWithError(conn, http.StatusForbidden, fmt.Sprintf("sellerGetRatingHandler: Invalid session. Please login again"))
		return
	}

	if seller, statusCode, err := getSellerRating(ctx, r.SessionID); err != nil {
		common.TCPRespondWithError(conn, statusCode, fmt.Sprintf("sellerGetRatingHandler: exception while fetching rating of seller. %v", err))
		return
	} else {
		common.TCPRespondWithJSON(conn, http.StatusOK, r.SessionID, seller)
	}
}

func sellerCreateItemHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {

	log.Println("Printing client request", r.SessionID)
	if !validateSessionID(r.SessionID) {
		common.TCPRespondWithError(conn, http.StatusForbidden, fmt.Sprintf("sellerCreateItemHandler: Invalid session. Please login again"))
		return
	}

	var productModel ProductModel
	if err := json.Unmarshal(r.Body, &productModel); err != nil {
		common.TCPRespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("sellerCreateItemHandler: exception while parsing request. %v", err))
		return
	}

	if product, statusCode, err := createProduct(ctx, &productModel, r.SessionID); err != nil {
		common.TCPRespondWithError(conn, statusCode, fmt.Sprintf("sellerCreateItemHandler: exception while creating item. %v", err))
		return
	} else {
		common.TCPRespondWithJSON(conn, http.StatusCreated, r.SessionID, product)
	}
}

func sellerUpdateItemSalePriceHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {

	if !validateSessionID(r.SessionID) {
		common.TCPRespondWithError(conn, http.StatusForbidden, fmt.Sprintf("sellerUpdateItemSalePriceHandler: Invalid session. Please login again"))
		return
	}

	var productModel ProductModel
	if err := json.Unmarshal(r.Body, &productModel); err != nil {
		common.TCPRespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("sellerUpdateItemSalePriceHandler: exception while parsing request. %v", err))
		return
	}

	log.Println(productModel)
	if product, statusCode, err := changeItemSalePrice(ctx, &productModel, r.SessionID); err != nil {
		common.TCPRespondWithError(conn, statusCode, fmt.Sprintf("sellerUpdateItemSalePriceHandler: exception while updating item. %v", err))
		return
	} else {
		common.TCPRespondWithJSON(conn, http.StatusOK, r.SessionID, product)
	}
}

func sellerRemoveItemHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	if !validateSessionID(r.SessionID) {
		common.TCPRespondWithError(conn, http.StatusForbidden, fmt.Sprintf("sellerUpdateItemSalePriceHandler: Invalid session. Please login again"))
		return
	}

	var productModel ProductModel
	if err := json.Unmarshal(r.Body, &productModel); err != nil {
		common.TCPRespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("sellerUpdateItemSalePriceHandler: exception while parsing request. %v", err))
		return
	}

	if product, statusCode, err := removeItemFromSale(ctx, &productModel, r.SessionID); err != nil {
		common.TCPRespondWithError(conn, statusCode, fmt.Sprintf("sellerUpdateItemSalePriceHandler: exception while updating item. %v", err))
		return
	} else {
		common.TCPRespondWithJSON(conn, http.StatusOK, r.SessionID, product)
	}
}

func sellerGetSellerItemsHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	if !validateSessionID(r.SessionID) {
		common.TCPRespondWithError(conn, http.StatusForbidden, fmt.Sprintf("sellerGetSellerItemsHandler: Invalid session. Please login again"))
		return
	}

	if products, statusCode, err := getSellerProducts(ctx, r.SessionID); err != nil {
		common.TCPRespondWithError(conn, statusCode, fmt.Sprintf("sellerGetSellerItemsHandler: exception while fetching seller items. %v", err))
		return
	} else {
		common.TCPRespondWithJSON(conn, http.StatusOK, r.SessionID, products)
	}
}

func sellerGetSoldItemsHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {

	if !validateSessionID(r.SessionID) {
		common.TCPRespondWithError(conn, http.StatusForbidden, fmt.Sprintf("sellerGetSoldItemsHandler: Invalid session. Please login again"))
		return
	}

	if transactions, statusCode, err := getTransactionListBySellerID(ctx, r.SessionID); err != nil {
		common.TCPRespondWithError(conn, statusCode, fmt.Sprintf("sellerGetSoldItemsHandler: exception while fetching seller sold items. %v", err))
		return
	} else {
		common.TCPRespondWithJSON(conn, http.StatusOK, r.SessionID, transactions)
	}
}

func listOfSellerHandlers(ctx context.Context, conn net.Conn, req common.ClientRequest) {

	switch req.Service {
	case "0":
		sellerCreateAccountHandler(ctx, conn, req)
	case "1":
		sellerLoginHandler(ctx, conn, req)
	case "2":
		sellerLogoutHandler(ctx, conn, req)
	case "3":
		sellerGetRatingHandler(ctx, conn, req)
	case "4":
		sellerCreateItemHandler(ctx, conn, req)
	case "5":
		sellerUpdateItemSalePriceHandler(ctx, conn, req)
	case "6":
		sellerRemoveItemHandler(ctx, conn, req)
	case "7":
		sellerGetSellerItemsHandler(ctx, conn, req)
	}
}
