package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"net"
	"net/http"
)

func buyerCreateAccountHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	var buyerModel BuyerModel
	if err := json.Unmarshal(r.Body, &buyerModel); err != nil {
		common.RespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("buyerCreateAccountHandler: exception while parsing request. %v", err))
		return
	}

	if createdObj, _, err := createBuyerAccount(ctx, &buyerModel); err != nil {
		common.RespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("buyerCreateAccountHandler: exception while creating buyer. %v", err))
		return
	} else {
		common.RespondWithJSON(conn, http.StatusCreated, "", createdObj)
	}
}
func buyerLoginHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	var buyerModel BuyerModel
	if err := json.Unmarshal(r.Body, &buyerModel); err != nil {
		common.RespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("buyerLoginHandler: exception while parsing request. %v", err))
		return
	}

	if session, statusCode, err := buyerLogin(ctx, buyerModel.UserName, buyerModel.Password); err != nil {
		common.RespondWithError(conn, statusCode, fmt.Sprintf("buyerLoginHandler: exception while Logging in buyer. %v", err))
		return
	} else {
		common.RespondWithJSON(conn, http.StatusCreated, session, map[string]string{"sessionID": session})
	}
}

func buyerLogoutHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	if !validateSessionID(r.SessionID) {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerLogoutHandler: Invalid session. Please login again"))
		return
	}

	if statusCode, err := buyerLogout(ctx, r.SessionID); err != nil {
		common.RespondWithError(conn, statusCode, fmt.Sprintf("buyerLogoutHandler: exception while Logging out buyer. %v", err))
		return
	} else {
		common.RespondWithStatusCode(conn, http.StatusOK, r.SessionID)
	}
}

func buyerSearchItemsHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {

	if !validateSessionID(r.SessionID) {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerSearchItemsHandler: Invalid session. Please login again"))
		return
	}

	var productModel ProductModel
	if err := json.Unmarshal(r.Body, &productModel); err != nil {
		common.RespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("buyerSearchItemsHandler: exception while parsing request. %v", err))
		return
	}

	if products, statusCode, err := searchProduct(ctx, &productModel); err != nil {
		common.RespondWithError(conn, statusCode, fmt.Sprintf("buyerSearchItemsHandler: exception while searching item. %v", err))
		return
	} else {
		common.RespondWithJSON(conn, http.StatusOK, r.SessionID, products)
	}
}

func buyerAddItemToCartHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {

	if !validateSessionID(r.SessionID) {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerAddItemToCartHandler: Invalid session. Please login again"))
		return
	}

	var productModel ProductModel
	if err := json.Unmarshal(r.Body, &productModel); err != nil {
		common.RespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("buyerAddItemToCartHandler: exception while parsing request. %v", err))
		return
	}
	if statusCode, err := buyerAddProductToCart(ctx, r.SessionID, &productModel); err != nil {
		common.RespondWithError(conn, statusCode, fmt.Sprintf("buyerAddItemToCartHandler: exception while adding item. %v", err))
		return
	} else {
		common.RespondWithStatusCode(conn, http.StatusOK, r.SessionID)
	}
}

func buyerRemoveItemFromCartHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	if !validateSessionID(r.SessionID) {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerRemoveItemFromCartHandler: Invalid session. Please login again"))
		return
	}

	var productModel ProductModel
	if err := json.Unmarshal(r.Body, &productModel); err != nil {
		common.RespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("buyerRemoveItemFromCartHandler: exception while parsing request. %v", err))
		return
	}

	if statusCode, err := buyerRemoveProductToCart(ctx, r.SessionID, &productModel); err != nil {
		common.RespondWithError(conn, statusCode, fmt.Sprintf("buyerRemoveItemFromCartHandler: exception while removing item. %v", err))
		return
	} else {
		common.RespondWithStatusCode(conn, http.StatusOK, r.SessionID)
	}
}

func buyerSaveCartHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	if !validateSessionID(r.SessionID) {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerSaveCartHandler: Invalid session. Please login again"))
		return
	}

	if statusCode, err := buyerSaveCart(ctx, r.SessionID); err != nil {
		common.RespondWithError(conn, statusCode, fmt.Sprintf("buyerSaveCartHandler: exception while saving cart. %v", err))
		return
	} else {
		common.RespondWithStatusCode(conn, http.StatusOK, r.SessionID)
	}
}

func buyerClearCartHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	if !validateSessionID(r.SessionID) {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerSaveCartHandler: Invalid session. Please login again"))
		return
	}

	if statusCode, err := buyerClearCart(ctx, r.SessionID); err != nil {
		common.RespondWithError(conn, statusCode, fmt.Sprintf("buyerSaveCartHandler: exception while clearing cart. %v", err))
		return
	} else {
		common.RespondWithStatusCode(conn, http.StatusOK, r.SessionID)
	}
}

func buyerGetCartHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	if !validateSessionID(r.SessionID) {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerSaveCartHandler: Invalid session. Please login again"))
		return
	}

	if cart, statusCode, err := buyerGetCart(ctx, r.SessionID); err != nil {
		common.RespondWithError(conn, statusCode, fmt.Sprintf("buyerSaveCartHandler: exception while clearing cart. %v", err))
		return
	} else {
		common.RespondWithJSON(conn, http.StatusOK, r.SessionID, cart)
	}
}

func buyerMakePurchaseHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	common.RespondWithStatusCode(conn, http.StatusNotImplemented, r.SessionID)
}

func buyerProvideProductFeedBackHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	if !validateSessionID(r.SessionID) {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerProvideProductFeedBackHandler: Invalid session. Please login again"))
		return
	}
	var vars map[string]string
	if err := json.Unmarshal(r.Body, &vars); err != nil {
		common.RespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("buyerProvideProductFeedBackHandler: exception while unmarshalling the request. %v", err))
		return
	}

	if vars["productID"] == "" {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerProvideProductFeedBackHandler: productID query param empty"))
		return
	}

	if vars["rating"] == "" {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerProvideProductFeedBackHandler: rating query param empty"))
		return
	}
	liked := false
	productID := vars["productID"]

	if vars["rating"] == "liked" {
		liked = true
	}

	if statusCode, err := buyerProvideProductFeedBack(ctx, r.SessionID, productID, liked); err != nil {
		common.RespondWithError(conn, statusCode, fmt.Sprintf("buyerProvideProductFeedBackHandler: exception while adding item feedback. %v", err))
		return
	} else {
		common.RespondWithStatusCode(conn, http.StatusOK, r.SessionID)
	}
}

func buyerGetProductSellerRatingHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {

	if !validateSessionID(r.SessionID) {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerGetProductSellerRatingHandler: Invalid session. Please login again"))
		return
	}
	var vars map[string]string
	if err := json.Unmarshal(r.Body, &vars); err != nil {
		common.RespondWithError(conn, http.StatusBadRequest, fmt.Sprintf("buyerGetProductSellerRatingHandler: exception while unmarshalling the request. %v", err))
		return
	}

	if vars["productID"] == "" {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerGetProductSellerRatingHandler: productID query param empty"))
		return
	}

	productID := vars["productID"]

	if seller, statusCode, err := getProductSellerRatingByProductID(ctx, productID); err != nil {
		common.RespondWithError(conn, statusCode, fmt.Sprintf("buyerGetProductSellerRatingHandler: exception while fetching seller rating for productID %s. %v", productID, err))
		return
	} else {
		common.RespondWithJSON(conn, http.StatusOK, r.SessionID, seller)
	}
}

func buyerGetPurchaseHistoryHandler(ctx context.Context, conn net.Conn, r common.ClientRequest) {
	if !validateSessionID(r.SessionID) {
		common.RespondWithError(conn, http.StatusForbidden, fmt.Sprintf("buyerGetProductSellerRatingHandler: Invalid session. Please login again"))
		return
	}

	if transactions, statusCode, err := getTransactionListByBuyerID(ctx, r.SessionID); err != nil {
		common.RespondWithError(conn, statusCode, fmt.Sprintf("buyerGetProductSellerRatingHandler: exception while fetching buyer purchased items. %v", err))
		return
	} else {
		common.RespondWithJSON(conn, http.StatusOK, r.SessionID, transactions)
	}
}

func listOfBuyerHandlers(ctx context.Context, conn net.Conn, req common.ClientRequest) {

	switch req.Service {
	case "0":
		buyerCreateAccountHandler(ctx, conn, req)
	case "1":
		buyerLoginHandler(ctx, conn, req)
	case "2":
		buyerLogoutHandler(ctx, conn, req)
	case "3":
		buyerSearchItemsHandler(ctx, conn, req)
	case "4":
		buyerAddItemToCartHandler(ctx, conn, req)
	case "5":
		buyerRemoveItemFromCartHandler(ctx, conn, req)
	case "6":
		buyerSaveCartHandler(ctx, conn, req)
	case "7":
		buyerClearCartHandler(ctx, conn, req)
	case "8":
		buyerGetCartHandler(ctx, conn, req)
	case "9":
		buyerProvideProductFeedBackHandler(ctx, conn, req)
	case "10":
		buyerGetProductSellerRatingHandler(ctx, conn, req)
	case "11":
		buyerGetPurchaseHistoryHandler(ctx, conn, req)
	}
}
