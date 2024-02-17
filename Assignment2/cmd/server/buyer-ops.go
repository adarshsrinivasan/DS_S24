package main

import (
	"context"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/library/wsdl/transaction"
	"net/http"
	"time"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/sirupsen/logrus"
)

type BuyerModel struct {
	Id                     string    `json:"id,omitempty" bson:"id" bun:"id,pk"`
	Name                   string    `json:"name,omitempty" bson:"name" bun:"name,notnull"`
	UserName               string    `json:"userName,omitempty" bson:"userName" bun:"userName,notnull,unique"`
	Password               string    `json:"password,omitempty" bson:"password" bun:"password,notnull,unique"`
	Version                int       `json:"version,omitempty" bson:"version" bun:"version,notnull"`
	CreatedAt              time.Time `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}

type PurchaseDetailsModel struct {
	Name string `json:"name,omitempty" bson:"name" bun:"name,pk"`
	CreditCardNumber string `json:"creditCardNumber,omitempty" bson:"creditCardNumber" bun:"creditCardNumber,pk"`
	Expiry string `json:"expiry,omitempty" bson:"expiry" bun:"expiry,pk"`
}

func createBuyerAccount(ctx context.Context, buyerModel *BuyerModel) (BuyerModel, int, error) {
	if err := validateBuyerModel(ctx, buyerModel, true); err != nil {
		err := fmt.Errorf("exception while validating Buyer data. %v", err)
		logrus.Errorf("createSellerAccount: %v\n", err)
		return BuyerModel{}, http.StatusBadRequest, err
	}
	buyerModel.Id = ""
	if statusCode, err := buyerModel.CreateBuyer(ctx); err != nil {
		err := fmt.Errorf("exception while creating Seller. %v", err)
		logrus.Errorf("createSellerAccount: %v\n", err)
		return BuyerModel{}, statusCode, err
	}
	return *buyerModel, http.StatusOK, nil
}

func buyerLogin(ctx context.Context, userName, password string) (string, int, error) {
	buyerModelObj := BuyerModel{UserName: userName}
	if statusCode, err := buyerModelObj.GetBuyerByUserName(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Buyer by username %s. %v", userName, err)
		logrus.Errorf("buyerLogin: %v\n", err)
		return "", statusCode, err
	}

	if buyerModelObj.Password != password {
		err := fmt.Errorf("worng username/password for username: %s", userName)
		logrus.Errorf("buyerLogin: %v\n", err)
		return "", http.StatusForbidden, err
	}

	return createNewSession(ctx, buyerModelObj.Id, common.BUYER)
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

func buyerAddProductToCart(ctx context.Context, sessionID string, productModel *ProductModel) (int, error) {
	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("buyerAddProductToCart: %v\n", err)
		return statusCode, err
	}
	if userType != common.BUYER {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("buyerAddProductToCart: %v\n", err)
		return http.StatusBadRequest, err
	}
	cartTableModel := CartModel{
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
	}

	exists := false
	existingCartItemModel := CartItemModel{
		CartID:    cartTableModel.ID,
		ProductID: product.ID,
	}
	if _, err := existingCartItemModel.GetCartItemByCartIDAndProductID(ctx); err == nil {
		cartItemModel.Quantity += existingCartItemModel.Quantity
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
	if userType != common.BUYER {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("buyerRemoveProductToCart: %v\n", err)
		return http.StatusBadRequest, err
	}
	cartTableModel := CartModel{
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
	}
	if statusCode, err := removeProductFromCart(ctx, &cartItemModel); err != nil {
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
	if userType != common.BUYER {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("buyerSaveCart: %v\n", err)
		return http.StatusBadRequest, err
	}
	cartTableModel := CartModel{
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
	if userType != common.BUYER {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("buyerClearCart: %v\n", err)
		return http.StatusBadRequest, err
	}
	cartTableModel := CartModel{
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
	if userType != common.BUYER {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("buyerGetCart: %v\n", err)
		return CartModel{}, http.StatusBadRequest, err
	}
	cartTableModel := CartModel{
		BuyerID: userID,
	}
	if _, err := cartTableModel.GetCartByBuyerID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Cart with buyerID %s. %v", userID, err)
		logrus.Errorf("buyerGetCart: %v\n", err)
		return CartModel{}, statusCode, err
	}

	cartModel, statusCode, err := getCartByID(ctx, cartTableModel.ID)
	if err != nil {
		err := fmt.Errorf("exception while get Cart for buyerID %s. %v", userID, err)
		logrus.Errorf("buyerGetCart: %v\n", err)
		return CartModel{}, statusCode, err
	}
	return *cartModel, http.StatusOK, nil
}

func buyerProvideProductFeedBack(ctx context.Context, sessionID, productID string, liked bool) (int, error) {
	transactionModels, statusCode, err := getTransactionListByBuyerID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Traansaction List for sessionID %s. %v", sessionID, err)
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
		err := fmt.Errorf("couldn't find product %s in user's %s purchace history", productID, transactionModels[0].BuyerID)
		logrus.Errorf("buyerProvideProductFeedBack: %v\n", err)
		return statusCode, err
	}

	if liked {
		if statusCode, err := incrementProductRating(ctx, productID); err != nil {
			err := fmt.Errorf("exception while incrementing product %s rating. %v", productID, err)
			logrus.Errorf("buyerProvideProductFeedBack: %v\n", err)
			return statusCode, err
		}
	} else {
		if statusCode, err := decrementProductRating(ctx, productID); err != nil {
			err := fmt.Errorf("exception while decrementing product %s rating. %v", productID, err)
			logrus.Errorf("buyerProvideProductFeedBack: %v\n", err)
			return statusCode, err
		}
	}

	return http.StatusOK, nil
}

func buyerMakeTransaction(ctx context.Context, sessionID string, purchaseDetailsModel PurchaseDetailsModel) (int, error) {
	resp, err := transactionService.IsTransactionApproved(&transaction.TransactionRequest{
		Name:              purchaseDetailsModel.Name,
		CreditCardDetails: purchaseDetailsModel.CreditCardNumber,
		Expiry:            purchaseDetailsModel.Expiry,
	})
	if err != nil {
		err = fmt.Errorf("exception while contacting transaction server. %v\n", err)
		logrus.Errorf("buyerMakeTransaction: %v\n", err)
		return http.StatusInternalServerError, err
	}

	if resp.Approved {
		userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
		if err != nil {
			err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
			logrus.Errorf("buyerGetCart: %v\n", err)
			return statusCode, err
		}
		if userType != common.BUYER {
			err := fmt.Errorf("user not a buyer type: %s", userID)
			logrus.Errorf("buyerGetCart: %v\n", err)
			return http.StatusBadRequest, err
		}
		cartModel, statusCode, err := buyerGetCart(ctx, sessionID)
		if err != nil {
			err := fmt.Errorf("exception while Fetching Cart by sessionID %s. %v", sessionID, err)
			logrus.Errorf("makeTransaction: %v\n", err)
			return statusCode, err
		}

		for _, item := range cartModel.Items {
			product, _, err := getProductByID(ctx, item.ProductID)
			if err != nil {
				err := fmt.Errorf("exception while fetching product with ID %s. %v", item.ProductID, err)
				logrus.Errorf("buyerMakeTransaction: %v\n", err)
				continue
			}

			if product.Quantity <= item.Quantity {
				logrus.Infof("buyerMakeTransaction: Attemting to buy %d count of %s product, while only %d count left in stock. Changing purchase quantity to %d.", item.Quantity, product.ID, product.Quantity, product.Quantity)
				product.Quantity = 0
				item.Quantity = product.Quantity
				if statusCode, err := product.DeleteProductByID(ctx); err != nil {
					err = fmt.Errorf("exception while Deleting Product for ID:%s. %v", product.ID, err)
					logrus.Errorf("buyerMakeTransaction: %v\n", err)
					return statusCode, err
				}
			} else {
				product.Quantity -= item.Quantity
				if statusCode, err := product.UpdateProductByID(ctx); err != nil {
					err = fmt.Errorf("exception while Updating Product for ID:%s. %v", product.ID, err)
					logrus.Errorf("buyerMakeTransaction: %v\n", err)
					return statusCode, err
				}
			}

			transaction := TransactionModel{
				CartID:    item.CartID,
				ProductID: item.ProductID,
				BuyerID:   userID,
				SellerID:  item.SellerID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}

			if _, err := transaction.CreateTransaction(ctx); err != nil {
				err := fmt.Errorf("exception while creating Transaction for CartItemID %s. %v", item.ID, err)
				logrus.Errorf("buyerMakeTransaction: %v\n", err)
				continue
			}
		}
		if statusCode, err := clearCart(ctx, cartModel.ID); err != nil {
			err := fmt.Errorf("exception while clear Cart for buyerID %s. %v", userID, err)
			logrus.Errorf("makeTransaction: %v\n", err)
			return statusCode, err
		}
		return http.StatusOK, nil
	} else {
		return http.StatusBadRequest, fmt.Errorf("Transaction failed.")
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

	return nil
}
