package main

import (
	"context"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/adarshsrinivasan/DS_S24/library/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

const (
	TransactionTableName      = "transaction_data"
	TransactionTableAliasName = "transaction"
)

type TransactionTableOps interface {
	CreateTransaction(ctx context.Context) (int, error)
	ListTransactionsBySellerID(ctx context.Context) ([]TransactionModel, int, error)
	ListTransactionsByBuyerID(ctx context.Context) ([]TransactionModel, int, error)
	ListTransactionsByCartID(ctx context.Context) ([]TransactionModel, int, error)
	DeleteTransactionsByCartID(ctx context.Context) (int, error)
	DeleteTransactionsByBuyerID(ctx context.Context) (int, error)
	DeleteTransactionsBySellerID(ctx context.Context) (int, error)
}

func (transaction *TransactionModel) CreateTransaction(ctx context.Context) (int, error) {
	protoModel := convertTransactionModelToProtoTransactionModel(ctx, transaction)
	request := &proto.CreateTransactionRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("CreateTransaction: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.CreateTransaction(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", TransactionTableName, err)
		logrus.Errorf("CreateTransaction: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyTransactionObj(response.ResponseModel, transaction)
	logrus.Infof("CreateTransaction: Successfully Recorded transaction for product %s Cart %s Seller %s Buyer %s\n", transaction.ProductID, transaction.CartID, transaction.SellerID, transaction.BuyerID)
	return http.StatusOK, nil
}

func (transaction *TransactionModel) ListTransactionsByCartID(ctx context.Context) ([]TransactionModel, int, error) {
	protoModel := convertTransactionModelToProtoTransactionModel(ctx, transaction)
	request := &proto.ListTransactionsByCartIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("ListTransactionsByCartID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.ListTransactionsByCartID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", TransactionTableName, err)
		logrus.Errorf("ListTransactionsByCartID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	var result []TransactionModel
	for _, resp := range response.ResponseModel {
		result = append(result, *convertProtoTransactionModelToTransactionModel(ctx, resp))
	}
	return result, http.StatusOK, nil
}

func (transaction *TransactionModel) ListTransactionsByBuyerID(ctx context.Context) ([]TransactionModel, int, error) {
	protoModel := convertTransactionModelToProtoTransactionModel(ctx, transaction)
	request := &proto.ListTransactionsByBuyerIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("ListTransactionsByBuyerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.ListTransactionsByBuyerID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", TransactionTableName, err)
		logrus.Errorf("ListTransactionsByBuyerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	var result []TransactionModel
	for _, resp := range response.ResponseModel {
		result = append(result, *convertProtoTransactionModelToTransactionModel(ctx, resp))
	}
	return result, http.StatusOK, nil
}

func (transaction *TransactionModel) ListTransactionsBySellerID(ctx context.Context) ([]TransactionModel, int, error) {
	protoModel := convertTransactionModelToProtoTransactionModel(ctx, transaction)
	request := &proto.ListTransactionsBySellerIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("ListTransactionsBySellerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.ListTransactionsBySellerID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", TransactionTableName, err)
		logrus.Errorf("ListTransactionsBySellerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	var result []TransactionModel
	for _, resp := range response.ResponseModel {
		result = append(result, *convertProtoTransactionModelToTransactionModel(ctx, resp))
	}
	return result, http.StatusOK, nil
}

func (transaction *TransactionModel) DeleteTransactionsByCartID(ctx context.Context) (int, error) {
	protoModel := convertTransactionModelToProtoTransactionModel(ctx, transaction)
	request := &proto.DeleteTransactionsByCartIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("DeleteTransactionsByCartID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	if _, err := sqlDBClient.DeleteTransactionsByCartID(ctx, request); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", TransactionTableName, err)
		logrus.Errorf("DeleteTransactionsByCartID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (transaction *TransactionModel) DeleteTransactionsBySellerID(ctx context.Context) (int, error) {
	protoModel := convertTransactionModelToProtoTransactionModel(ctx, transaction)
	request := &proto.DeleteTransactionsBySellerIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("DeleteTransactionsBySellerID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	if _, err := sqlDBClient.DeleteTransactionsBySellerID(ctx, request); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", TransactionTableName, err)
		logrus.Errorf("DeleteTransactionsBySellerID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (transaction *TransactionModel) DeleteTransactionsByBuyerID(ctx context.Context) (int, error) {
	protoModel := convertTransactionModelToProtoTransactionModel(ctx, transaction)
	request := &proto.DeleteTransactionsByBuyerIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("DeleteTransactionsByBuyerID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	if _, err := sqlDBClient.DeleteTransactionsByBuyerID(ctx, request); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", TransactionTableName, err)
		logrus.Errorf("DeleteTransactionsByBuyerID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func copyTransactionObj(from *proto.TransactionModel, to *TransactionModel) {
	to.ID = from.ID
	to.CartID = from.CartID
	to.ProductID = from.ProductID
	to.BuyerID = from.BuyerID
	to.SellerID = from.SellerID
	to.Quantity = int(from.Quantity)
	to.Price = from.Price
	to.Version = int(from.Version)
	to.CreatedAt = from.CreatedAt.AsTime()
	to.UpdatedAt = from.UpdatedAt.AsTime()
}

func convertTransactionModelToProtoTransactionModel(ctx context.Context, model *TransactionModel) *proto.TransactionModel {
	return &proto.TransactionModel{
		ID:        model.ID,
		CartID:    model.CartID,
		ProductID: model.ProductID,
		BuyerID:   model.BuyerID,
		SellerID:  model.SellerID,
		Quantity:  int32(model.Quantity),
		Price:     model.Price,
		Version:   int32(model.Version),
		CreatedAt: timestamppb.New(model.CreatedAt),
		UpdatedAt: timestamppb.New(model.CreatedAt),
	}
}

func convertProtoTransactionModelToTransactionModel(ctx context.Context, protoModel *proto.TransactionModel) *TransactionModel {
	return &TransactionModel{
		ID:        protoModel.ID,
		CartID:    protoModel.CartID,
		ProductID: protoModel.ProductID,
		BuyerID:   protoModel.BuyerID,
		SellerID:  protoModel.SellerID,
		Quantity:  int(protoModel.Quantity),
		Price:     protoModel.Price,
		Version:   int(protoModel.Version),
		CreatedAt: protoModel.CreatedAt.AsTime(),
		UpdatedAt: protoModel.UpdatedAt.AsTime(),
	}
}
