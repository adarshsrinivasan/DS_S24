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
	SellerTableName      = "seller_data"
	SellerTableAliasName = "seller"
)

type SellerTableOps interface {
	CreateSeller(ctx context.Context) (int, error)
	GetSellerByID(ctx context.Context) (int, error)
	GetSellerByUserName(ctx context.Context) (int, error)
	UpdateSellerByID(ctx context.Context) (int, error)
}

func (seller *SellerModel) CreateSeller(ctx context.Context) (int, error) {
	protoModel := convertSellerModelToProtoSellerModel(ctx, seller)
	request := &proto.CreateSellerRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("CreateSeller: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.CreateSeller(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", SellerTableName, err)
		logrus.Errorf("CreateSeller: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copySellerObj(response.ResponseModel, seller)
	logrus.Infof("CreateSeller: Successfully created account for userneme %s\n", seller.UserName)
	return http.StatusOK, nil
}

func (seller *SellerModel) GetSellerByID(ctx context.Context) (int, error) {
	protoModel := convertSellerModelToProtoSellerModel(ctx, seller)
	request := &proto.GetSellerByIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetSellerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.GetSellerByID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", SellerTableName, err)
		logrus.Errorf("GetSellerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copySellerObj(response.ResponseModel, seller)
	return http.StatusOK, nil
}

func (seller *SellerModel) GetSellerByUserName(ctx context.Context) (int, error) {
	protoModel := convertSellerModelToProtoSellerModel(ctx, seller)
	request := &proto.GetSellerByUserNameRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetSellerByUserName: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.GetSellerByUserName(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", SellerTableName, err)
		logrus.Errorf("GetSellerByUserName: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copySellerObj(response.ResponseModel, seller)
	return http.StatusOK, nil
}

func (seller *SellerModel) UpdateSellerByID(ctx context.Context) (int, error) {
	protoModel := convertSellerModelToProtoSellerModel(ctx, seller)
	request := &proto.UpdateSellerByIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("UpdateSellerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.UpdateSellerByID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", SellerTableName, err)
		logrus.Errorf("UpdateSellerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copySellerObj(response.ResponseModel, seller)
	return http.StatusOK, nil
}

func copySellerObj(from *proto.SellerModel, to *SellerModel) {
	to.Id = from.ID
	to.Name = from.Name
	to.FeedBackThumbsUp = int(from.FeedBackThumbsUp)
	to.FeedBackThumbsDown = int(from.FeedBackThumbsDown)
	to.NumberOfItemsSold = int(from.NumberOfItemsSold)
	to.UserName = from.UserName
	to.Password = from.Password
	to.Version = int(from.Version)
	to.CreatedAt = from.CreatedAt.AsTime()
	to.UpdatedAt = from.UpdatedAt.AsTime()
}

func convertSellerModelToProtoSellerModel(ctx context.Context, model *SellerModel) *proto.SellerModel {
	return &proto.SellerModel{
		ID:                 model.Id,
		Name:               model.Name,
		FeedBackThumbsUp:   int32(model.FeedBackThumbsUp),
		FeedBackThumbsDown: int32(model.FeedBackThumbsDown),
		NumberOfItemsSold:  int32(model.NumberOfItemsSold),
		UserName:           model.UserName,
		Password:           model.Password,
		Version:            int32(model.Version),
		CreatedAt:          timestamppb.New(model.CreatedAt),
		UpdatedAt:          timestamppb.New(model.CreatedAt),
	}
}
