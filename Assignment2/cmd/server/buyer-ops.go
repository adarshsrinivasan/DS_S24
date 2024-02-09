package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adarshsrinivasan/DS_S24/libraries/common"
	"github.com/sirupsen/logrus"
)

type BuyerModel struct {
	Id                     string    `json:"id,omitempty" bson:"id" bun:"id,pk"`
	Name                   string    `json:"name,omitempty" bson:"name" bun:"name,notnull"`
	NumberOfItemsPurchased int       `json:"numberOfItemsPurchased,omitempty" bson:"numberOfItemsPurchased" bun:"numberOfItemsPurchased"`
	UserName               string    `json:"userName,omitempty" bson:"userName" bun:"userName,notnull,unique"`
	Password               string    `json:"password,omitempty" bson:"password" bun:"password,notnull,unique"`
	Version                int       `json:"version,omitempty" bson:"version" bun:"version,notnull"`
	CreatedAt              time.Time `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}

func createBuyerAccount(ctx context.Context, buyerModel *BuyerModel) (BuyerModel, int, error) {
	if err := validateBuyerModel(ctx, buyerModel, true); err != nil {
		err := fmt.Errorf("exception while validating Buyer data. %v", err)
		logrus.Errorf("createSellerAccount: %v\n", err)
		return BuyerModel{}, http.StatusBadRequest, err
	}
	buyerTableModelObj := convertBuyerModelToBuyerTableModel(ctx, buyerModel)
	buyerTableModelObj.Id = ""
	if statusCode, err := buyerTableModelObj.CreateBuyer(ctx); err != nil {
		err := fmt.Errorf("exception while creating Seller. %v", err)
		logrus.Errorf("createSellerAccount: %v\n", err)
		return BuyerModel{}, statusCode, err
	}
	createdBuyerModel := convertBuyerTableModelToBuyerModel(ctx, buyerTableModelObj)
	return *createdBuyerModel, http.StatusOK, nil
}

func buyerLogin(ctx context.Context, userName, password string) (string, int, error) {
	buyerTableModelObj := BuyerTableModel{UserName: userName}
	if statusCode, err := buyerTableModelObj.GetBuyerByUserName(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Buyer by username %s. %v", userName, err)
		logrus.Errorf("buyerLogin: %v\n", err)
		return "", statusCode, err
	}

	if buyerTableModelObj.Password != password {
		err := fmt.Errorf("worng username/password for username: %s", userName)
		logrus.Errorf("buyerLogin: %v\n", err)
		return "", http.StatusForbidden, err
	}

	return createNewSession(ctx, buyerTableModelObj.Id, common.Buyer)
}

func buyerLogout(ctx context.Context, sessionID string) (int, error) {
	_, _, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("buyerLogout: %v\n", err)
		return statusCode, err
	}

	return deleteSessionByID(ctx, sessionID)
}

func getSellerRatingBySellerID(ctx context.Context, sellerID string) (int, int, int, error) {
	sellerTableModelObj := SellerTableModel{Id: sellerID}
	if statusCode, err := sellerTableModelObj.GetSellerByID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Seller by userID %s. %v", sellerID, err)
		logrus.Errorf("getSellerRatingBySellerID: %v\n", err)
		return -1, -1, statusCode, err
	}
	return sellerTableModelObj.FeedBackThumbsUp, sellerTableModelObj.FeedBackThumbsDown, http.StatusOK, nil
}

func buyerAddProductToCart(ctx context.Context, sessionID string, productModel *ProductModel) (int, error) {
	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("buyerAddProductToCart: %v\n", err)
		return statusCode, err
	}
	if userType != common.Buyer {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("buyerAddProductToCart: %v\n", err)
		return http.StatusBadRequest, err
	}
	cartTableModel := CartTableModel{
		BuyerID: userID,
	}
	if _, err := cartTableModel.GetCartByBuyerID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Cart with buyerID %s. %v. Attempting to create a new cart", userID, err)
		logrus.Errorf("buyerAddProductToCart: %v\n", err)
		if statusCode, err := cartTableModel.CreateCart(ctx); err != nil {
			err := fmt.Errorf("exception while creating Cart for buyerID %s. %v", userID, err)
			logrus.Errorf("buyerAddProductToCart: %v\n", err)
			return statusCode, err
		}
	}
	product, statusCode, err := getProductByID(ctx, productModel.ID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Product with ID %s. %v", productModel.ID, err)
		logrus.Errorf("buyerAddProductToCart: %v\n", err)
		return statusCode, err
	}

	cartItemModel := CartItemModel{
		CartID:    cartTableModel.ID,
		ProductID: product.ID,
		SellerID:  product.SellerID,
		Quantity:  productModel.Quantity,
		Price:     product.SalePrice * float32(productModel.Quantity),
	}

	exists := false
	existingCartItemModel := CartItemTableModel{
		CartID:    cartTableModel.ID,
		ProductID: product.ID,
	}
	if _, err := existingCartItemModel.GetCartItemByCartIDAndProductID(ctx); err == nil {
		cartItemModel.Quantity += existingCartItemModel.Quantity
		cartItemModel.Price += existingCartItemModel.Price
		exists = true
	}

	if product.Quantity < cartItemModel.Quantity {
		err := fmt.Errorf("attempting to add more quantity (%d) than available (%d)", cartItemModel.Quantity, product.Quantity)
		logrus.Errorf("buyerAddProductToCart: %v\n", err)
		return http.StatusBadRequest, err
	}

	if exists {
		if statusCode, err := updateCartItemByCartIDAndProductID(ctx, &cartItemModel); err != nil {
			err := fmt.Errorf("exception while adding Item %s to Cart %s. %v", productModel.ID, cartTableModel.ID, err)
			logrus.Errorf("buyerAddProductToCart: %v\n", err)
			return statusCode, err
		}
	} else {
		if statusCode, err := addProductToCart(ctx, &cartItemModel); err != nil {
			err := fmt.Errorf("exception while adding Item %s to Cart %s. %v", productModel.ID, cartTableModel.ID, err)
			logrus.Errorf("buyerAddProductToCart: %v\n", err)
			return statusCode, err
		}
	}

	return http.StatusOK, nil
}

func buyerRemoveProductToCart(ctx context.Context, sessionID string, productModel *ProductModel) (int, error) {
	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("buyerRemoveProductToCart: %v\n", err)
		return statusCode, err
	}
	if userType != common.Buyer {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("buyerRemoveProductToCart: %v\n", err)
		return http.StatusBadRequest, err
	}
	cartTableModel := CartTableModel{
		BuyerID: userID,
	}
	if _, err := cartTableModel.GetCartByBuyerID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Cart with buyerID %s. %v", userID, err)
		logrus.Errorf("buyerRemoveProductToCart: %v\n", err)
		return statusCode, err
	}
	product, statusCode, err := getProductByID(ctx, productModel.ID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Product with ID %s. %v", productModel.ID, err)
		logrus.Errorf("buyerRemoveProductToCart: %v\n", err)
		return statusCode, err
	}
	cartItemModel := CartItemModel{
		CartID:    cartTableModel.ID,
		ProductID: product.ID,
		SellerID:  product.SellerID,
		Quantity:  productModel.Quantity,
		Price:     product.SalePrice * float32(productModel.Quantity),
	}
	if statusCode, err := removeProductToCart(ctx, &cartItemModel); err != nil {
		err := fmt.Errorf("exception while removing Item %s to Cart %s. %v", productModel.ID, cartTableModel.ID, err)
		logrus.Errorf("buyerRemoveProductToCart: %v\n", err)
		return statusCode, err
	}
	return http.StatusOK, nil
}

func buyerSaveCart(ctx context.Context, sessionID string) (int, error) {
	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("buyerSaveCart: %v\n", err)
		return statusCode, err
	}
	if userType != common.Buyer {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("buyerSaveCart: %v\n", err)
		return http.StatusBadRequest, err
	}
	cartTableModel := CartTableModel{
		BuyerID: userID,
	}
	if _, err := cartTableModel.GetCartByBuyerID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Cart with buyerID %s. %v", userID, err)
		logrus.Errorf("buyerSaveCart: %v\n", err)
		return statusCode, err
	}

	if statusCode, err := saveCart(ctx, cartTableModel.ID); err != nil {
		err := fmt.Errorf("exception while saving Cart for buyerID %s. %v", userID, err)
		logrus.Errorf("buyerSaveCart: %v\n", err)
		return statusCode, err
	}
	return http.StatusOK, nil
}

func buyerClearCart(ctx context.Context, sessionID string) (int, error) {
	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("buyerClearCart: %v\n", err)
		return statusCode, err
	}
	if userType != common.Buyer {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("buyerClearCart: %v\n", err)
		return http.StatusBadRequest, err
	}
	cartTableModel := CartTableModel{
		BuyerID: userID,
	}
	if _, err := cartTableModel.GetCartByBuyerID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Cart with buyerID %s. %v", userID, err)
		logrus.Errorf("buyerClearCart: %v\n", err)
		return statusCode, err
	}

	if statusCode, err := clearCart(ctx, cartTableModel.ID); err != nil {
		err := fmt.Errorf("exception while clear Cart for buyerID %s. %v", userID, err)
		logrus.Errorf("buyerClearCart: %v\n", err)
		return statusCode, err
	}
	return http.StatusOK, nil
}

func buyerGetCart(ctx context.Context, sessionID string) (CartModel, int, error) {
	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("buyerGetCart: %v\n", err)
		return CartModel{}, statusCode, err
	}
	if userType != common.Buyer {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("buyerGetCart: %v\n", err)
		return CartModel{}, http.StatusBadRequest, err
	}
	cartTableModel := CartTableModel{
		BuyerID: userID,
	}
	if _, err := cartTableModel.GetCartByBuyerID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Cart with buyerID %s. %v", userID, err)
		logrus.Errorf("buyerGetCart: %v\n", err)
		return CartModel{}, statusCode, err
	}

	cartModel, statusCode, err := getCartByID(ctx, cartTableModel.ID)
	if err != nil {
		err := fmt.Errorf("exception while clear Cart for buyerID %s. %v", userID, err)
		logrus.Errorf("buyerGetCart: %v\n", err)
		return CartModel{}, statusCode, err
	}
	return *cartModel, http.StatusOK, nil
}

func buyerProvideProductFeedBack(ctx context.Context, sessionID, productID string, liked bool) (int, error) {
	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("buyerProvideProductFeedBack: %v\n", err)
		return statusCode, err
	}
	if userType != common.Buyer {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("buyerProvideProductFeedBack: %v\n", err)
		return http.StatusBadRequest, err
	}
	transactionModels, statusCode, err := getTransactionListByBuyerID(ctx, userID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Traansaction List for userID %s. %v", userID, err)
		logrus.Errorf("buyerProvideProductFeedBack: %v\n", err)
		return statusCode, err
	}

	found := false
	for _, t := range transactionModels {
		if t.ProductID == productID {
			found = true
			break
		}
	}
	if !found {
		err := fmt.Errorf("couldn't find product %s in user's %s purchace history", productID, userID)
		logrus.Errorf("buyerProvideProductFeedBack: %v\n", err)
		return statusCode, err
	}

	if liked {
		if statusCode, err := incrementProductRating(ctx, productID); err != nil {
			err := fmt.Errorf("exception while incrementing product %s rating. %v", productID, userID)
			logrus.Errorf("buyerProvideProductFeedBack: %v\n", err)
			return statusCode, err
		}
	} else {
		if statusCode, err := decrementProductRating(ctx, productID); err != nil {
			err := fmt.Errorf("exception while decrementing product %s rating. %v", productID, userID)
			logrus.Errorf("buyerProvideProductFeedBack: %v\n", err)
			return statusCode, err
		}
	}

	return http.StatusOK, nil
}

func convertBuyerModelToBuyerTableModel(ctx context.Context, buyerModel *BuyerModel) *BuyerTableModel {
	return &BuyerTableModel{
		Id:                     buyerModel.Id,
		Name:                   buyerModel.Name,
		NumberOfItemsPurchased: buyerModel.NumberOfItemsPurchased,
		UserName:               buyerModel.UserName,
		Password:               buyerModel.Password,
		Version:                buyerModel.Version,
		CreatedAt:              buyerModel.CreatedAt,
		UpdatedAt:              buyerModel.UpdatedAt,
	}
}

func convertBuyerTableModelToBuyerModel(ctx context.Context, buyerTableModel *BuyerTableModel) *BuyerModel {
	return &BuyerModel{
		Id:                     buyerTableModel.Id,
		Name:                   buyerTableModel.Name,
		NumberOfItemsPurchased: buyerTableModel.NumberOfItemsPurchased,
		UserName:               buyerTableModel.UserName,
		Password:               buyerTableModel.Password,
		Version:                buyerTableModel.Version,
		CreatedAt:              buyerTableModel.CreatedAt,
		UpdatedAt:              buyerTableModel.UpdatedAt,
	}
}

func validateBuyerModel(ctx context.Context, buyerModel *BuyerModel, create bool) error {

	if !create && buyerModel.Id == "" {
		err := fmt.Errorf("invalid Buyer data. ID field is empty")
		logrus.Errorf("validateSellerModel: %v\n", err)
		return err
	}

	if buyerModel.Name == "" {
		err := fmt.Errorf("invalid Buyer data. Name field is empty")
		logrus.Errorf("validateSellerModel: %v\n", err)
		return err
	}

	if buyerModel.UserName == "" {
		err := fmt.Errorf("invalid Buyer data. UserName field is empty")
		logrus.Errorf("validateSellerModel: %v\n", err)
		return err
	}

	if buyerModel.Password == "" {
		err := fmt.Errorf("invalid Buyer data. Password field is empty")
		logrus.Errorf("validateSellerModel: %v\n", err)
		return err
	}

	if create {
		buyerModel.NumberOfItemsPurchased = 0
	}
	return nil
}
