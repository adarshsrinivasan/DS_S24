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
	CartTableName      = "cart_data"
	CartTableAliasName = "cart"
)

type CartOps interface {
	CreateCart(ctx context.Context) (int, error)
	GetCartByID(ctx context.Context) (int, error)
	GetCartByBuyerID(ctx context.Context) (int, error)
	UpdateCartByID(ctx context.Context) (int, error)
	DeleteCartByID(ctx context.Context) (int, error)
}

func (cart *CartModel) CreateCart(ctx context.Context) (int, error) {
	protoModel := convertCartModelToProtoCartModel(ctx, cart)
	request := &proto.CreateCartRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("CreateCart: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.CreateCart(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", CartTableName, err)
		logrus.Errorf("CreateCart: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyCartObj(response.ResponseModel, cart)
	logrus.Infof("CreateCart: Successfully created cart for buyerID %s\n", cart.BuyerID)
	return http.StatusOK, nil
}

func (cart *CartModel) GetCartByID(ctx context.Context) (int, error) {
	protoModel := convertCartModelToProtoCartModel(ctx, cart)
	request := &proto.GetCartByIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.GetCartByID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartTableName, err)
		logrus.Errorf("GetCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyCartObj(response.ResponseModel, cart)
	return http.StatusOK, nil
}

func (cart *CartModel) GetCartByBuyerID(ctx context.Context) (int, error) {
	protoModel := convertCartModelToProtoCartModel(ctx, cart)
	request := &proto.GetCartByBuyerIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetCartByBuyerID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.GetCartByBuyerID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartTableName, err)
		logrus.Errorf("GetCartByBuyerID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyCartObj(response.ResponseModel, cart)
	return http.StatusOK, nil
}

func (cart *CartModel) UpdateCartByID(ctx context.Context) (int, error) {
	protoModel := convertCartModelToProtoCartModel(ctx, cart)
	request := &proto.UpdateCartByIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("UpdateCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.UpdateCartByID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", CartTableName, err)
		logrus.Errorf("UpdateCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyCartObj(response.ResponseModel, cart)
	return http.StatusOK, nil
}

func (cart *CartModel) DeleteCartByID(ctx context.Context) (int, error) {
	protoModel := convertCartModelToProtoCartModel(ctx, cart)
	request := &proto.DeleteCartByIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("DeleteCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	if _, err := sqlDBClient.DeleteCartByID(ctx, request); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartTableName, err)
		logrus.Errorf("DeleteCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func copyCartObj(from *proto.CartModel, to *CartModel) {
	to.ID = from.ID
	to.BuyerID = from.BuyerID
	to.Saved = from.Saved
	to.Version = int(from.Version)
	to.CreatedAt = from.CreatedAt.AsTime()
	to.UpdatedAt = from.UpdatedAt.AsTime()
}

func convertCartModelToProtoCartModel(ctx context.Context, model *CartModel) *proto.CartModel {
	return &proto.CartModel{
		ID:        model.ID,
		BuyerID:   model.BuyerID,
		Saved:     model.Saved,
		Version:   int32(model.Version),
		CreatedAt: timestamppb.New(model.CreatedAt),
		UpdatedAt: timestamppb.New(model.CreatedAt),
	}
}
