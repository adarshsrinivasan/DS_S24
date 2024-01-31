package main

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"net/http"
//
//	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
//	"github.com/gorilla/mux"
//)
//
//const (
//	ApiPrefix   = "/api/v1/marketplace"
//	IdUrlRegex  = "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}"
//	SlotIdRegex = "[0-9]+"
//	RateRegex   = "(?:liked|disliked)"
//)
//
//var (
//	httpRouter *mux.Router
//)
//
//func sellerCreateAccountHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//
//	var sellerModel SellerModel
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&sellerModel); err != nil {
//		common.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("sellerCreateAccountHandler: exception while parsing request. %v", err))
//		return
//	}
//	defer r.Body.Close()
//	if createdObj, statusCode, err := createSellerAccount(ctx, &sellerModel); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("sellerCreateAccountHandler: exception while creating user. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusCreated, r.Header.Get("User-Session-Id"), createdObj)
//	}
//}
//
//func sellerLoginHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//
//	var sellerModel SellerModel
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&sellerModel); err != nil {
//		common.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("sellerLoginHandler: exception while parsing request. %v", err))
//		return
//	}
//	defer r.Body.Close()
//	if session, statusCode, err := sellerLogin(ctx, sellerModel.UserName, sellerModel.Password); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("sellerLoginHandler: exception while Logging in user. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusCreated, r.Header.Get("User-Session-Id"), map[string]string{"sessionID": session})
//	}
//}
//
//func sellerLogoutHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("sellerLogoutHandler: Invalid session. Please login again"))
//		return
//	}
//
//	defer r.Body.Close()
//	if statusCode, err := sellerLogout(ctx, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("sellerLogoutHandler: exception while Logging out user. %v", err))
//		return
//	} else {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//}
//
//func sellerGetRatingHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("sellerGetRatingHandler: Invalid session. Please login again"))
//		return
//	}
//
//	defer r.Body.Close()
//	if seller, statusCode, err := getSellerRating(ctx, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("sellerGetRatingHandler: exception while fetching rating of seller. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusOK, r.Header.Get("User-Session-Id"), seller)
//	}
//}
//
//func sellerCreateItemHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("sellerCreateItemHandler: Invalid session. Please login again"))
//		return
//	}
//
//	var productModel ProductModel
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&productModel); err != nil {
//		common.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("sellerCreateItemHandler: exception while parsing request. %v", err))
//		return
//	}
//	defer r.Body.Close()
//
//	if product, statusCode, err := createProduct(ctx, &productModel, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("sellerCreateItemHandler: exception while creating item. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusCreated, r.Header.Get("User-Session-Id"), product)
//	}
//}
//
//func sellerUpdateItemSalePriceHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("sellerUpdateItemSalePriceHandler: Invalid session. Please login again"))
//		return
//	}
//
//	var productModel ProductModel
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&productModel); err != nil {
//		common.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("sellerUpdateItemSalePriceHandler: exception while parsing request. %v", err))
//		return
//	}
//	defer r.Body.Close()
//	if product, statusCode, err := changeItemSalePrice(ctx, &productModel, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("sellerUpdateItemSalePriceHandler: exception while updating item. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusOK, r.Header.Get("User-Session-Id"), product)
//	}
//}
//
//func sellerRemoveItemHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("sellerUpdateItemSalePriceHandler: Invalid session. Please login again"))
//		return
//	}
//
//	var productModel ProductModel
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&productModel); err != nil {
//		common.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("sellerUpdateItemSalePriceHandler: exception while parsing request. %v", err))
//		return
//	}
//	defer r.Body.Close()
//	if product, statusCode, err := removeItemFromSale(ctx, &productModel, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("sellerUpdateItemSalePriceHandler: exception while updating item. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusOK, r.Header.Get("User-Session-Id"), product)
//	}
//}
//
//func sellerGetSellerItemsHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("sellerGetSellerItemsHandler: Invalid session. Please login again"))
//		return
//	}
//
//	defer r.Body.Close()
//	if products, statusCode, err := getSellerProducts(ctx, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("sellerGetSellerItemsHandler: exception while fetching seller items. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusOK, r.Header.Get("User-Session-Id"), products)
//	}
//}
//
//func sellerGetSoldItemsHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("sellerGetSoldItemsHandler: Invalid session. Please login again"))
//		return
//	}
//
//	defer r.Body.Close()
//	if transactions, statusCode, err := getTransactionListBySellerID(ctx, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("sellerGetSoldItemsHandler: exception while fetching seller sold items. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusOK, r.Header.Get("User-Session-Id"), transactions)
//	}
//}
//
////func buyerCreateAccountHandler(w http.ResponseWriter, r *http.Request) {
////	// Stop here if its Preflighted OPTIONS request
////	if r.Method == "OPTIONS" {
////		common.RespondWithStatusCode(w, http.StatusOK, nil)
////	}
////
////	var buyerModel BuyerModel
////	decoder := json.NewDecoder(r.Body)
////	if err := decoder.Decode(&buyerModel); err != nil {
////		common.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("buyerCreateAccountHandler: exception while parsing request. %v", err))
////		return
////	}
////	defer r.Body.Close()
////	if createdObj, statusCode, err := createBuyerAccount(ctx, &buyerModel); err != nil {
////		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerCreateAccountHandler: exception while creating buyer. %v", err))
////		return
////	} else {
////		common.RespondWithJSON(w, http.StatusCreated, r.Header.Get("User-Session-Id"), createdObj)
////	}
////}
//
//func buyerLoginHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//
//	var buyerModel BuyerModel
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&buyerModel); err != nil {
//		common.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("buyerLoginHandler: exception while parsing request. %v", err))
//		return
//	}
//	defer r.Body.Close()
//	if session, statusCode, err := buyerLogin(ctx, buyerModel.UserName, buyerModel.Password); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerLoginHandler: exception while Logging in buyer. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusCreated, r.Header.Get("User-Session-Id"), map[string]string{"sessionID": session})
//	}
//}
//
//func buyerLogoutHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerLogoutHandler: Invalid session. Please login again"))
//		return
//	}
//
//	defer r.Body.Close()
//	if statusCode, err := buyerLogout(ctx, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerLogoutHandler: exception while Logging out buyer. %v", err))
//		return
//	} else {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//}
//
//func buyerSearchItemsHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerSearchItemsHandler: Invalid session. Please login again"))
//		return
//	}
//
//	var productModel ProductModel
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&productModel); err != nil {
//		common.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("buyerSearchItemsHandler: exception while parsing request. %v", err))
//		return
//	}
//	defer r.Body.Close()
//	if products, statusCode, err := searchProduct(ctx, &productModel); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerSearchItemsHandler: exception while searching item. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusOK, r.Header.Get("User-Session-Id"), products)
//	}
//}
//
//func buyerAddItemToCartHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerAddItemToCartHandler: Invalid session. Please login again"))
//		return
//	}
//
//	var productModel ProductModel
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&productModel); err != nil {
//		common.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("buyerAddItemToCartHandler: exception while parsing request. %v", err))
//		return
//	}
//	defer r.Body.Close()
//	if statusCode, err := buyerAddProductToCart(ctx, r.Header.Get("User-Session-Id"), &productModel); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerAddItemToCartHandler: exception while adding item. %v", err))
//		return
//	} else {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//}
//
//func buyerRemoveItemFromCartHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerRemoveItemFromCartHandler: Invalid session. Please login again"))
//		return
//	}
//
//	var productModel ProductModel
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&productModel); err != nil {
//		common.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("buyerRemoveItemFromCartHandler: exception while parsing request. %v", err))
//		return
//	}
//	defer r.Body.Close()
//	if statusCode, err := buyerRemoveProductToCart(ctx, r.Header.Get("User-Session-Id"), &productModel); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerRemoveItemFromCartHandler: exception while removing item. %v", err))
//		return
//	} else {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//}
//
//func buyerSaveCartHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerSaveCartHandler: Invalid session. Please login again"))
//		return
//	}
//
//	defer r.Body.Close()
//	if statusCode, err := buyerSaveCart(ctx, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerSaveCartHandler: exception while saving cart. %v", err))
//		return
//	} else {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//}
//
//func buyerClearCartHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerSaveCartHandler: Invalid session. Please login again"))
//		return
//	}
//
//	defer r.Body.Close()
//	if statusCode, err := buyerClearCart(ctx, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerSaveCartHandler: exception while clearing cart. %v", err))
//		return
//	} else {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//}
//
//func buyerGetCartHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerSaveCartHandler: Invalid session. Please login again"))
//		return
//	}
//
//	defer r.Body.Close()
//	if cart, statusCode, err := buyerGetCart(ctx, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerSaveCartHandler: exception while clearing cart. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusOK, r.Header.Get("User-Session-Id"), cart)
//	}
//}
//
//func buyerMakePurchaseHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	common.RespondWithStatusCode(w, http.StatusNotImplemented, nil)
//}
//
//func buyerProvideProductFeedBackHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerProvideProductFeedBackHandler: Invalid session. Please login again"))
//		return
//	}
//	vars := mux.Vars(r)
//
//	if vars["productID"] == "" {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerProvideProductFeedBackHandler: productID query param empty"))
//		return
//	}
//
//	if vars["rating"] == "" {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerProvideProductFeedBackHandler: rating query param empty"))
//		return
//	}
//	liked := false
//	productID := vars["productID"]
//
//	if vars["rating"] == "liked" {
//		liked = true
//	}
//
//	defer r.Body.Close()
//	if statusCode, err := buyerProvideProductFeedBack(ctx, r.Header.Get("User-Session-Id"), productID, liked); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerProvideProductFeedBackHandler: exception while adding item feedback. %v", err))
//		return
//	} else {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//}
//
//func buyerGetProductSellerRatingHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerGetProductSellerRatingHandler: Invalid session. Please login again"))
//		return
//	}
//	vars := mux.Vars(r)
//
//	if vars["productID"] == "" {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerGetProductSellerRatingHandler: productID query param empty"))
//		return
//	}
//
//	productID := vars["productID"]
//
//	defer r.Body.Close()
//	if seller, statusCode, err := getProductSellerRatingByProductID(ctx, productID); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerGetProductSellerRatingHandler: exception while fetching seller rating for productID %s. %v", productID, err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusOK, r.Header.Get("User-Session-Id"), seller)
//	}
//}
//
//func buyerGetPurchaseHistoryHandler(w http.ResponseWriter, r *http.Request) {
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		common.RespondWithStatusCode(w, http.StatusOK, nil)
//	}
//	if !validateSessionID(r.Header.Get("User-Session-Id")) {
//		common.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("buyerGetProductSellerRatingHandler: Invalid session. Please login again"))
//		return
//	}
//
//	defer r.Body.Close()
//	if transactions, statusCode, err := getTransactionListByBuyerID(ctx, r.Header.Get("User-Session-Id")); err != nil {
//		common.RespondWithError(w, statusCode, fmt.Sprintf("buyerGetProductSellerRatingHandler: exception while fetching buyer purchased items. %v", err))
//		return
//	} else {
//		common.RespondWithJSON(w, http.StatusOK, r.Header.Get("User-Session-Id"), transactions)
//	}
//}
//
//func initializeHttpRoutes(ctx context.Context) {
//	httpRouter = mux.NewRouter()
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "seller", "create"),
//		sellerCreateAccountHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "seller", "login"),
//		sellerLoginHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "seller", "logout"),
//		sellerLogoutHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "seller", "getRating"),
//		sellerGetRatingHandler).Methods("GET", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "seller", "createItem"),
//		sellerCreateItemHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "seller", "updateItemSalePrice"),
//		sellerUpdateItemSalePriceHandler).Methods("PUT", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "seller", "removeItem"),
//		sellerRemoveItemHandler).Methods("PUT", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "seller", "getItems"),
//		sellerGetSellerItemsHandler).Methods("GET", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "seller", "getSoldItems"),
//		sellerGetSoldItemsHandler).Methods("GET", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "buyer", "create"),
//		buyerCreateAccountHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "buyer", "login"),
//		buyerLoginHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "buyer", "logout"),
//		buyerLogoutHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "buyer", "searchItems"),
//		buyerSearchItemsHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "buyer", "addItemToCart"),
//		buyerAddItemToCartHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "buyer", "removeItemFromCart"),
//		buyerRemoveItemFromCartHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "buyer", "saveCart"),
//		buyerSaveCartHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "buyer", "clearCart"),
//		buyerClearCartHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "buyer", "getCart"),
//		buyerGetCartHandler).Methods("GET", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "buyer", "makePurchase"),
//		buyerMakePurchaseHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s/%s/%s", ApiPrefix, "buyer", "feedback", fmt.Sprintf("{productID:%s}", IdUrlRegex), fmt.Sprintf("{rating:%s}", RateRegex)),
//		buyerProvideProductFeedBackHandler).Methods("POST", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s/%s", ApiPrefix, "buyer", "getSellerRating", fmt.Sprintf("{productID:%s}", IdUrlRegex)),
//		buyerGetProductSellerRatingHandler).Methods("GET", "OPTIONS")
//	httpRouter.HandleFunc(fmt.Sprintf("%s/%s/%s", ApiPrefix, "buyer", "getPurchaseHistory"),
//		buyerGetPurchaseHistoryHandler).Methods("GET", "OPTIONS")
//}
