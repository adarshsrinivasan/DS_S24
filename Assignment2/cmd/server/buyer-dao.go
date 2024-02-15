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
	BuyerTableName      = "buyer_data"
	BuyerTableAliasName = "buyer"
)

type BuyerOps interface {
	CreateBuyer(ctx context.Context) (int, error)
	GetBuyerByID(ctx context.Context) (int, error)
	GetBuyerByUserName(ctx context.Context) (int, error)
	UpdateBuyerByID(ctx context.Context) (int, error)
}

func (buyer *BuyerModel) CreateBuyer(ctx context.Context) (int, error) {
	protoModel := convertBuyerModelToProtoBuyerModel(ctx, buyer)
	request := &proto.CreateBuyerRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("CreateBuyer: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.CreateBuyer(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", BuyerTableName, err)
		logrus.Errorf("CreateBuyer: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyBuyerObj(response.ResponseModel, buyer)
	logrus.Infof("CreateBuyer: Successfully created account for userneme %s\n", buyer.UserName)
	return http.StatusOK, nil
}

func (buyer *BuyerModel) GetBuyerByID(ctx context.Context) (int, error) {
	protoBuyerModel := convertBuyerModelToProtoBuyerModel(ctx, buyer)
	request := &proto.GetBuyerByIDRequest{
		RequestModel: protoBuyerModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.GetBuyerByID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", BuyerTableName, err)
		logrus.Errorf("GetBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	copyBuyerObj(response.ResponseModel, buyer)
	return http.StatusOK, nil
}

func (buyer *BuyerModel) GetBuyerByUserName(ctx context.Context) (int, error) {
	protoBuyerModel := convertBuyerModelToProtoBuyerModel(ctx, buyer)
	request := &proto.GetBuyerByUserNameRequest{
		RequestModel: protoBuyerModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetBuyerByUserName: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.GetBuyerByUserName(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", BuyerTableName, err)
		logrus.Errorf("GetBuyerByUserName: %v\n", err)
		return http.StatusInternalServerError, err
	}

	copyBuyerObj(response.ResponseModel, buyer)
	return http.StatusOK, nil
}

func (buyer *BuyerModel) UpdateBuyerByID(ctx context.Context) (int, error) {
	protoBuyerModel := convertBuyerModelToProtoBuyerModel(ctx, buyer)
	request := &proto.UpdateBuyerByIDRequest{
		RequestModel: protoBuyerModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("UpdateBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.UpdateBuyerByID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", BuyerTableName, err)
		logrus.Errorf("UpdateBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	copyBuyerObj(response.ResponseModel, buyer)
	return http.StatusOK, nil
}

func copyBuyerObj(from *proto.BuyerModel, to *BuyerModel) {
	to.Id = from.ID
	to.Name = from.Name
	to.UserName = from.UserName
	to.Password = from.Password
	to.Version = int(from.Version)
	to.CreatedAt = from.CreatedAt.AsTime()
	to.UpdatedAt = from.UpdatedAt.AsTime()
}

func convertBuyerModelToProtoBuyerModel(ctx context.Context, buyerModel *BuyerModel) *proto.BuyerModel {
	return &proto.BuyerModel{
		ID:                     buyerModel.Id,
		Name:                   buyerModel.Name,
		UserName:               buyerModel.UserName,
		Password:               buyerModel.Password,
		Version:                int32(buyerModel.Version),
		CreatedAt:              timestamppb.New(buyerModel.CreatedAt),
		UpdatedAt:              timestamppb.New(buyerModel.CreatedAt),
	}
}
