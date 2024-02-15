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
	CartItemTableName      = "cartitem_data"
	CartItemTableAliasName = "cartitem"
)

type CartItemOps interface {
	CreateCartItem(ctx context.Context) (int, error)
	GetCartItemByID(ctx context.Context) (int, error)
	GetCartItemByCartIDAndProductID(ctx context.Context) (int, error)
	ListCartItemByCartID(ctx context.Context) ([]CartItemModel, int, error)
	UpdateCartItem(ctx context.Context) (int, error)
	DeleteCartItemByCartIDAndProductID(ctx context.Context) (int, error)
	DeleteCartItemByCartID(ctx context.Context) (int, error)
	DeleteCartItemByProductID(ctx context.Context) (int, error)
}

func (cartItem *CartItemModel) CreateCartItem(ctx context.Context) (int, error) {
	protoModel := convertCartItemModelToProtoCartItemModel(ctx, cartItem)
	request := &proto.CreateCartItemRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("CreateCartItem: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.CreateCartItem(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", CartItemTableName, err)
		logrus.Errorf("CreateCartItem: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyCartItemObj(response.ResponseModel, cartItem)
	logrus.Infof("CreateCartItem: Successfully created cartItem for cartID %s\n", cartItem.CartID)
	return http.StatusOK, nil
}

func (cartItem *CartItemModel) GetCartItemByID(ctx context.Context) (int, error) {
	protoModel := convertCartItemModelToProtoCartItemModel(ctx, cartItem)
	request := &proto.GetCartItemByIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetCartItemByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.GetCartItemByID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("GetCartItemByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyCartItemObj(response.ResponseModel, cartItem)
	return http.StatusOK, nil
}

func (cartItem *CartItemModel) GetCartItemByCartIDAndProductID(ctx context.Context) (int, error) {
	protoModel := convertCartItemModelToProtoCartItemModel(ctx, cartItem)
	request := &proto.GetCartItemByCartIDAndProductIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetCartItemByCartIDAndProductID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.GetCartItemByCartIDAndProductID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("GetCartItemByCartIDAndProductID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyCartItemObj(response.ResponseModel, cartItem)
	return http.StatusOK, nil
}

func (cartItem *CartItemModel) ListCartItemByCartID(ctx context.Context) ([]CartItemModel, int, error) {
	protoModel := convertCartItemModelToProtoCartItemModel(ctx, cartItem)
	request := &proto.ListCartItemByCartIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("ListCartItemByCartID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.ListCartItemByCartID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("ListCartItemByCartID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	var result []CartItemModel
	for _, resp := range response.ResponseModel {
		result = append(result, *convertProtoCartItemModelToCartItemModel(ctx, resp))
	}
	return result, http.StatusOK, nil
}

func (cartItem *CartItemModel) UpdateCartItem(ctx context.Context) (int, error) {
	protoModel := convertCartItemModelToProtoCartItemModel(ctx, cartItem)
	request := &proto.UpdateCartItemRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("UpdateCartItem: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.UpdateCartItem(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", CartItemTableName, err)
		logrus.Errorf("UpdateCartItem: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyCartItemObj(response.ResponseModel, cartItem)
	return http.StatusOK, nil
}

func (cartItem *CartItemModel) DeleteCartItemByCartIDAndProductID(ctx context.Context) (int, error) {
	protoModel := convertCartItemModelToProtoCartItemModel(ctx, cartItem)
	request := &proto.DeleteCartItemByCartIDAndProductIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("DeleteCartItemByCartIDAndProductID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	if _, err := sqlDBClient.DeleteCartItemByCartIDAndProductID(ctx, request); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartItemTableName, err)
		logrus.Errorf("DeleteCartItemByCartIDAndProductID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (cartItem *CartItemModel) DeleteCartItemByCartID(ctx context.Context) (int, error) {
	protoModel := convertCartItemModelToProtoCartItemModel(ctx, cartItem)
	request := &proto.DeleteCartItemByCartIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("DeleteCartItemByCartID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	if _, err := sqlDBClient.DeleteCartItemByCartID(ctx, request); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartItemTableName, err)
		logrus.Errorf("DeleteCartItemByCartID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (cartItem *CartItemModel) DeleteCartItemByProductID(ctx context.Context) (int, error) {
	protoModel := convertCartItemModelToProtoCartItemModel(ctx, cartItem)
	request := &proto.DeleteCartItemByProductIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("DeleteCartItemByProductID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	if _, err := sqlDBClient.DeleteCartItemByProductID(ctx, request); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartItemTableName, err)
		logrus.Errorf("DeleteCartItemByProductID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func copyCartItemObj(from *proto.CartItemModel, to *CartItemModel) {
	to.ID = from.ID
	to.CartID = from.CartID
	to.ProductID = from.ProductID
	to.Quantity = int(from.Quantity)
	to.Version = int(from.Version)
	to.CreatedAt = from.CreatedAt.AsTime()
	to.UpdatedAt = from.UpdatedAt.AsTime()
}

func convertCartItemModelToProtoCartItemModel(ctx context.Context, model *CartItemModel) *proto.CartItemModel {
	return &proto.CartItemModel{
		ID:        model.ID,
		CartID:    model.CartID,
		ProductID: model.ProductID,
		Quantity:  int32(model.Quantity),
		Version:   int32(model.Version),
		CreatedAt: timestamppb.New(model.CreatedAt),
		UpdatedAt: timestamppb.New(model.CreatedAt),
	}
}

func convertProtoCartItemModelToCartItemModel(ctx context.Context, protoModel *proto.CartItemModel) *CartItemModel {
	return &CartItemModel{
		ID:        protoModel.ID,
		CartID:    protoModel.CartID,
		ProductID: protoModel.ProductID,
		Quantity:  int(protoModel.Quantity),
		Version:   int(protoModel.Version),
		CreatedAt: protoModel.CreatedAt.AsTime(),
		UpdatedAt: protoModel.UpdatedAt.AsTime(),
	}
}
