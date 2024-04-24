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
	SessionTableName = "session_data"
)

type SessionTableOps interface {
	CreateSession(ctx context.Context) (int, error)
	GetSessionByID(ctx context.Context) (int, error)
	GetSessionByUserID(ctx context.Context) (int, error)
	DeleteSessionByID(ctx context.Context) (int, error)
}

func (session *SessionModel) CreateSession(ctx context.Context) (int, error) {
	protoModel := convertSessionModelToProtoSessionModel(ctx, session)
	request := &proto.CreateSessionRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("CreateSession: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.CreateSession(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", SessionTableName, err)
		logrus.Errorf("CreateSession: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copySessionObj(response.ResponseModel, session)
	logrus.Infof("CreateSession: Successfully created account for userID %s\n", session.UserID)
	return http.StatusOK, nil
}

func (session *SessionModel) GetSessionByID(ctx context.Context) (int, error) {
	protoModel := convertSessionModelToProtoSessionModel(ctx, session)
	request := &proto.GetSessionByIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetSessionByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.GetSessionByID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", SessionTableName, err)
		logrus.Errorf("GetSessionByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copySessionObj(response.ResponseModel, session)
	return http.StatusOK, nil
}

func (session *SessionModel) GetSessionByUserID(ctx context.Context) (int, error) {
	protoModel := convertSessionModelToProtoSessionModel(ctx, session)
	request := &proto.GetSessionByUserIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetSessionByUserID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := sqlDBClient.GetSessionByUserID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", SessionTableName, err)
		logrus.Errorf("GetSessionByUserID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copySessionObj(response.ResponseModel, session)
	return http.StatusOK, nil
}

func (session *SessionModel) DeleteSessionByID(ctx context.Context) (int, error) {
	protoModel := convertSessionModelToProtoSessionModel(ctx, session)
	request := &proto.DeleteSessionByIDRequest{
		RequestModel: protoModel,
	}
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("DeleteSessionByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	if _, err := sqlDBClient.DeleteSessionByID(ctx, request); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", SessionTableName, err)
		logrus.Errorf("DeleteSessionByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func copySessionObj(from *proto.SessionModel, to *SessionModel) {
	to.ID = from.ID
	to.UserID = from.UserID
	to.UserType = common.UserType(from.UserType)
	to.Version = int(from.Version)
	to.CreatedAt = from.CreatedAt.AsTime()
	to.UpdatedAt = from.UpdatedAt.AsTime()
}

func convertSessionModelToProtoSessionModel(ctx context.Context, model *SessionModel) *proto.SessionModel {
	return &proto.SessionModel{
		ID:        model.ID,
		UserID:    model.UserID,
		UserType:  proto.USERTYPE(model.UserType),
		Version:   int32(model.Version),
		CreatedAt: timestamppb.New(model.CreatedAt),
		UpdatedAt: timestamppb.New(model.CreatedAt),
	}
}
