package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adarshsrinivasan/DS_S24/libraries/common"
	"github.com/sirupsen/logrus"
)

type TransactionModel struct {
	ID        string    `json:"id,omitempty" bson:"id" bun:"id,pk"`
	CartID    string    `json:"cartID,omitempty" bson:"cartID" bun:"cartID,notnull"`
	ProductID string    `json:"productID,omitempty" bson:"productID" bun:"productID,notnull"`
	BuyerID   string    `json:"buyerID,omitempty" bson:"buyerID" bun:"buyerID,notnull"`
	SellerID  string    `json:"sellerID,omitempty" bson:"sellerID" bun:"sellerID,notnull"`
	Quantity  int       `json:"quantity,omitempty" bson:"quantity" bun:"quantity,notnull"`
	Price     float32   `json:"price,omitempty" bson:"price,omitempty" bun:"quantity,notnull"`
	CreatedAt time.Time `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}

func createTransaction(ctx context.Context, transactionModel *TransactionModel) (int, error) {
	if err := validateTransactionModel(ctx, transactionModel, true); err != nil {
		err := fmt.Errorf("exception while validating Transaction data. %v", err)
		logrus.Errorf("createTransaction: %v\n", err)
		return http.StatusBadRequest, err
	}
	transactionTableModelObj := convertTransactionModelToTransactionTableModel(ctx, transactionModel)
	transactionTableModelObj.ID = ""
	return transactionTableModelObj.CreateTransaction(ctx)
}

func getTransactionListByCartID(ctx context.Context, cartID string) ([]TransactionModel, int, error) {
	transactionTableModelObj := TransactionTableModel{CartID: cartID}
	transactionTableModels, statusCode, err := transactionTableModelObj.ListTransactionsByCartID(ctx)
	if err != nil {
		err := fmt.Errorf("exception while reading Transaction by cartID %s. %v", cartID, err)
		logrus.Errorf("getTransactionByCartID: %v\n", err)
		return nil, statusCode, err
	}

	transactionsModels := []TransactionModel{}

	for _, t := range transactionTableModels {
		transactionModelObj := convertTransactionTableModelToTransactionModel(ctx, &t)
		transactionsModels = append(transactionsModels, *transactionModelObj)
	}

	return transactionsModels, http.StatusOK, nil
}

func getTransactionListByBuyerID(ctx context.Context, sessionID string) ([]TransactionModel, int, error) {

	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("getTransactionListByBuyerID: %v\n", err)
		return nil, statusCode, err
	}
	if userType != common.Buyer {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("getTransactionListByBuyerID: %v\n", err)
		return nil, http.StatusBadRequest, err
	}

	transactionTableModelObj := TransactionTableModel{BuyerID: userID}
	transactionTableModels, statusCode, err := transactionTableModelObj.ListTransactionsByBuyerID(ctx)
	if err != nil {
		err := fmt.Errorf("exception while reading Transaction by buyerID %s. %v", userID, err)
		logrus.Errorf("getTransactionListByBuyerID: %v\n", err)
		return nil, statusCode, err
	}

	transactionsModels := []TransactionModel{}

	for _, t := range transactionTableModels {
		transactionModelObj := convertTransactionTableModelToTransactionModel(ctx, &t)
		transactionsModels = append(transactionsModels, *transactionModelObj)
	}

	return transactionsModels, http.StatusOK, nil
}

func getTransactionListBySellerID(ctx context.Context, sessionID string) ([]TransactionModel, int, error) {
	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("getTransactionListBySellerID: %v\n", err)
		return nil, statusCode, err
	}
	if userType != common.Seller {
		err := fmt.Errorf("user not a buyer type: %s", userID)
		logrus.Errorf("getTransactionListBySellerID: %v\n", err)
		return nil, http.StatusBadRequest, err
	}

	transactionTableModelObj := TransactionTableModel{SellerID: userID}
	transactionTableModels, statusCode, err := transactionTableModelObj.ListTransactionsBySellerID(ctx)
	if err != nil {
		err := fmt.Errorf("exception while reading Transaction by sellerID %s. %v", userID, err)
		logrus.Errorf("getTransactionListByBuyerID: %v\n", err)
		return nil, statusCode, err
	}

	transactionsModels := []TransactionModel{}

	for _, t := range transactionTableModels {
		transactionModelObj := convertTransactionTableModelToTransactionModel(ctx, &t)
		transactionsModels = append(transactionsModels, *transactionModelObj)
	}

	return transactionsModels, http.StatusOK, nil
}

func deleteTransactionByCartID(ctx context.Context, cartID string) (int, error) {
	transactionTableModelObj := TransactionTableModel{CartID: cartID}
	if statusCode, err := transactionTableModelObj.DeleteTransactionsByCartID(ctx); err != nil {
		err := fmt.Errorf("exception while delete Transaction by cartID %s. %v", cartID, err)
		logrus.Errorf("deleteTransactionByCartID: %v\n", err)
		return statusCode, err
	}
	return http.StatusOK, nil
}

func deleteTransactionByBuyerID(ctx context.Context, buyerID string) (int, error) {
	transactionTableModelObj := TransactionTableModel{BuyerID: buyerID}
	if statusCode, err := transactionTableModelObj.DeleteTransactionsByBuyerID(ctx); err != nil {
		err := fmt.Errorf("exception while delete Transaction by buyerID %s. %v", buyerID, err)
		logrus.Errorf("deleteTransactionByBuyerID: %v\n", err)
		return statusCode, err
	}
	return http.StatusOK, nil
}

func deleteTransactionBySellerID(ctx context.Context, sellerID string) (int, error) {
	transactionTableModelObj := TransactionTableModel{SellerID: sellerID}
	if statusCode, err := transactionTableModelObj.DeleteTransactionsBySellerID(ctx); err != nil {
		err := fmt.Errorf("exception while delete Transaction by sellerID %s. %v", sellerID, err)
		logrus.Errorf("deleteTransactionBySellerID: %v\n", err)
		return statusCode, err
	}
	return http.StatusOK, nil
}

func validateTransactionModel(ctx context.Context, transactionModel *TransactionModel, create bool) error {

	if !create && transactionModel.ID == "" {
		err := fmt.Errorf("invalid Transaction data. ID field is empty")
		logrus.Errorf("validateTransactionModel: %v\n", err)
		return err
	}

	if transactionModel.CartID == "" {
		err := fmt.Errorf("invalid Transaction data. cartID field is empty")
		logrus.Errorf("validateTransactionModel: %v\n", err)
		return err
	}

	if transactionModel.BuyerID == "" {
		err := fmt.Errorf("invalid Transaction data. BuyerID field is empty")
		logrus.Errorf("validateTransactionModel: %v\n", err)
		return err
	}

	if transactionModel.SellerID == "" {
		err := fmt.Errorf("invalid Transaction data. SellerID field is empty")
		logrus.Errorf("validateTransactionModel: %v\n", err)
		return err
	}

	if transactionModel.ProductID == "" {
		err := fmt.Errorf("invalid Transaction data. ProductID field is empty")
		logrus.Errorf("validateTransactionModel: %v\n", err)
		return err
	}

	if transactionModel.Quantity <= 0 {
		err := fmt.Errorf("invalid Transaction data. Quantity field is empty")
		logrus.Errorf("validateTransactionModel: %v\n", err)
		return err
	}

	if transactionModel.Price <= 0 {
		err := fmt.Errorf("invalid Transaction data. Price field is empty")
		logrus.Errorf("validateTransactionModel: %v\n", err)
		return err
	}

	return nil
}

func convertTransactionModelToTransactionTableModel(ctx context.Context, transactionModel *TransactionModel) *TransactionTableModel {
	return &TransactionTableModel{
		ID:        transactionModel.ID,
		CartID:    transactionModel.CartID,
		ProductID: transactionModel.ProductID,
		BuyerID:   transactionModel.BuyerID,
		SellerID:  transactionModel.SellerID,
		Quantity:  transactionModel.Quantity,
		Price:     transactionModel.Price,
		CreatedAt: transactionModel.CreatedAt,
		UpdatedAt: transactionModel.UpdatedAt,
	}
}

func convertTransactionTableModelToTransactionModel(ctx context.Context, transactionTableModel *TransactionTableModel) *TransactionModel {
	return &TransactionModel{
		ID:        transactionTableModel.ID,
		CartID:    transactionTableModel.CartID,
		ProductID: transactionTableModel.ProductID,
		BuyerID:   transactionTableModel.BuyerID,
		SellerID:  transactionTableModel.SellerID,
		Quantity:  transactionTableModel.Quantity,
		Price:     transactionTableModel.Price,
		CreatedAt: transactionTableModel.CreatedAt,
		UpdatedAt: transactionTableModel.UpdatedAt,
	}
}
