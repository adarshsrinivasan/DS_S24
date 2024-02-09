package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adarshsrinivasan/DS_S24/libraries/common"
	"github.com/sirupsen/logrus"
)

type SellerModel struct {
	Id                 string    `json:"id,omitempty" bson:"id" bun:"id,pk"`
	Name               string    `json:"name,omitempty" bson:"name" bun:"name,notnull"`
	FeedBackThumbsUp   int       `json:"feedBackThumbsUp" bson:"feedBackThumbsUp" bun:"feedBackThumbsUp"`
	FeedBackThumbsDown int       `json:"feedBackThumbsDown" bson:"feedBackThumbsDown" bun:"feedBackThumbsDown"`
	NumberOfItemsSold  int       `json:"numberOfItemsSold,omitempty" bson:"numberOfItemsSold" bun:"numberOfItemsSold"`
	UserName           string    `json:"userName,omitempty" bson:"userName" bun:"userName,notnull,unique"`
	Password           string    `json:"password,omitempty" bson:"password" bun:"password,notnull,unique"`
	Version            int       `json:"version,omitempty" bson:"version" bun:"version,notnull"`
	CreatedAt          time.Time `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}

func createSellerAccount(ctx context.Context, sellerModel *SellerModel) (SellerModel, int, error) {
	if err := validateSellerModel(ctx, sellerModel, true); err != nil {
		err := fmt.Errorf("exception while validating Seller data. %v", err)
		logrus.Errorf("createSellerAccount: %v\n", err)
		return SellerModel{}, http.StatusBadRequest, err
	}
	sellerTableModelObj := convertSellerModelToSellerTableModel(ctx, sellerModel)
	sellerTableModelObj.Id = ""
	if statusCode, err := sellerTableModelObj.CreateSeller(ctx); err != nil {
		err := fmt.Errorf("exception while creating Seller. %v", err)
		logrus.Errorf("createSellerAccount: %v\n", err)
		return SellerModel{}, statusCode, err
	}
	createdSellerModel := convertSellerTableModelToSellerModel(ctx, sellerTableModelObj)
	return *createdSellerModel, http.StatusOK, nil
}

func sellerLogin(ctx context.Context, userName, password string) (string, int, error) {
	sellerTableModelObj := SellerTableModel{UserName: userName}
	if statusCode, err := sellerTableModelObj.GetSellerByUserName(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Seller by username %s. %v", userName, err)
		logrus.Errorf("sellerLogin: %v\n", err)
		return "", statusCode, err
	}

	if sellerTableModelObj.Password != password {
		err := fmt.Errorf("worng username/password for username: %s", userName)
		logrus.Errorf("sellerLogin: %v\n", err)
		return "", http.StatusForbidden, err
	}

	return createNewSession(ctx, sellerTableModelObj.Id, common.Seller)
}

func sellerLogout(ctx context.Context, sessionID string) (int, error) {
	_, _, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("sellerLogout: %v\n", err)
		return statusCode, err
	}

	return deleteSessionByID(ctx, sessionID)
}

func getSellerRating(ctx context.Context, sessionID string) (SellerModel, int, error) {
	userID, _, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("sellerLogout: %v\n", err)
		return SellerModel{}, statusCode, err
	}

	sellerTableModelObj := SellerTableModel{Id: userID}
	if statusCode, err := sellerTableModelObj.GetSellerByID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Seller by userID %s. %v", userID, err)
		logrus.Errorf("sellerLogin: %v\n", err)
		return SellerModel{}, statusCode, err
	}
	sellerModel := SellerModel{
		FeedBackThumbsUp:   sellerTableModelObj.FeedBackThumbsUp,
		FeedBackThumbsDown: sellerTableModelObj.FeedBackThumbsDown,
	}
	return sellerModel, http.StatusOK, nil
}

func incrementSellerRating(ctx context.Context, sellerID string) (int, error) {
	sellerTableModelObj := SellerTableModel{Id: sellerID}
	if statusCode, err := sellerTableModelObj.GetSellerByUserName(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Seller by userID %s. %v", sellerID, err)
		logrus.Errorf("sellerLogin: %v\n", err)
		return statusCode, err
	}
	sellerTableModelObj.FeedBackThumbsUp++
	return sellerTableModelObj.UpdateSellerByID(ctx)
}

func decrementSellerRating(ctx context.Context, sellerID string) (int, error) {
	sellerTableModelObj := SellerTableModel{Id: sellerID}
	if statusCode, err := sellerTableModelObj.GetSellerByUserName(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Seller by userID %s. %v", sellerID, err)
		logrus.Errorf("sellerLogin: %v\n", err)
		return statusCode, err
	}
	sellerTableModelObj.FeedBackThumbsDown++
	return sellerTableModelObj.UpdateSellerByID(ctx)
}

func convertSellerModelToSellerTableModel(ctx context.Context, sellerModel *SellerModel) *SellerTableModel {
	return &SellerTableModel{
		Id:                 sellerModel.Id,
		Name:               sellerModel.Name,
		FeedBackThumbsUp:   sellerModel.FeedBackThumbsUp,
		FeedBackThumbsDown: sellerModel.FeedBackThumbsDown,
		NumberOfItemsSold:  sellerModel.NumberOfItemsSold,
		UserName:           sellerModel.UserName,
		Password:           sellerModel.Password,
		Version:            sellerModel.Version,
		CreatedAt:          sellerModel.CreatedAt,
		UpdatedAt:          sellerModel.UpdatedAt,
	}
}

func convertSellerTableModelToSellerModel(ctx context.Context, sellerTableModel *SellerTableModel) *SellerModel {
	return &SellerModel{
		Id:                 sellerTableModel.Id,
		Name:               sellerTableModel.Name,
		FeedBackThumbsUp:   sellerTableModel.FeedBackThumbsUp,
		FeedBackThumbsDown: sellerTableModel.FeedBackThumbsDown,
		NumberOfItemsSold:  sellerTableModel.NumberOfItemsSold,
		UserName:           sellerTableModel.UserName,
		Password:           sellerTableModel.Password,
		Version:            sellerTableModel.Version,
		CreatedAt:          sellerTableModel.CreatedAt,
		UpdatedAt:          sellerTableModel.UpdatedAt,
	}
}

func validateSellerModel(ctx context.Context, sellerModel *SellerModel, create bool) error {

	if !create && sellerModel.Id == "" {
		err := fmt.Errorf("invalid Seller data. ID field is empty")
		logrus.Errorf("validateSellerModel: %v\n", err)
		return err
	}

	if sellerModel.Name == "" {
		err := fmt.Errorf("invalid Seller data. Name field is empty")
		logrus.Errorf("validateSellerModel: %v\n", err)
		return err
	}

	if sellerModel.UserName == "" {
		err := fmt.Errorf("invalid Seller data. UserName field is empty")
		logrus.Errorf("validateSellerModel: %v\n", err)
		return err
	}

	if sellerModel.Password == "" {
		err := fmt.Errorf("invalid Seller data. Password field is empty")
		logrus.Errorf("validateSellerModel: %v\n", err)
		return err
	}

	if create {
		sellerModel.FeedBackThumbsDown = 0
		sellerModel.FeedBackThumbsUp = 0
		sellerModel.NumberOfItemsSold = 0
	}
	return nil
}
