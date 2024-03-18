package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/adarshsrinivasan/DS_S24/library/db"
	"github.com/adarshsrinivasan/DS_S24/library/db/sql"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun/schema"
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

type SellerTableModel struct {
	schema.BaseModel   `bun:"table:seller_data,alias:seller"`
	Id                 string    `json:"id" bson:"id" bun:"id,pk"`
	Name               string    `json:"name" bson:"name" bun:"name,notnull"`
	FeedBackThumbsUp   int       `json:"feedBackThumbsUp" bson:"feedBackThumbsUp" bun:"feedBackThumbsUp"`
	FeedBackThumbsDown int       `json:"feedBackThumbsDown" bson:"feedBackThumbsDown" bun:"feedBackThumbsDown"`
	NumberOfItemsSold  int       `json:"numberOfItemsSold" bson:"numberOfItemsSold" bun:"numberOfItemsSold"`
	UserName           string    `json:"userName" bson:"userName" bun:"userName,notnull,unique"`
	Password           string    `json:"password" bson:"password" bun:"password,notnull,unique"`
	Version            int       `json:"version" bson:"version" bun:"version,notnull"`
	CreatedAt          time.Time `json:"createdAt"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" bson:"updatedAt" bun:"updatedAt"`
}

func CreateSellerTable(ctx context.Context) error {
	client, err := sql.NewSQLClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateSellerTable: %v\n", err)
		return err
	}
	defer client.Close(ctx)

	tableSchemaPtr := reflect.New(reflect.TypeOf(SellerTableModel{}))

	if err := client.CreateTable(ctx, tableSchemaPtr.Interface(), SellerTableName, nil); err != nil {
		err := fmt.Errorf("exception while creating table %s. %v", err, SellerTableName)
		logrus.Errorf("CreateSellerTable: %v\n", err)
		return err
	}

	seller := SellerTableModel{
		Name:               "admin",
		FeedBackThumbsUp:   0,
		FeedBackThumbsDown: 0,
		NumberOfItemsSold:  0,
		UserName:           "admin",
		Password:           "admin",
	}

	if statusCode, err := seller.CreateSeller(ctx); err != nil {
		err := fmt.Errorf("exception while inserting the admin Seller. %v", err)
		if statusCode == http.StatusBadRequest {
			return nil
		}
		logrus.Errorf("CreateSellerTable: %v. StatusCode: %v\n", statusCode, err)
		return err
	}

	return nil
}

func (seller *SellerTableModel) CreateSeller(ctx context.Context) (int, error) {
	client, err := sql.NewSQLClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateSellerTable: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	seller.Id = uuid.New().String()
	seller.Version = 0
	seller.CreatedAt = time.Now()
	seller.UpdatedAt = time.Now()

	if existingSeller, _, _ := seller.getByColumn(ctx, "userName", seller.UserName); existingSeller != nil && existingSeller.UserName == seller.UserName {
		err := fmt.Errorf("exception while verifying Seller data. userName alredy taken")
		logrus.Errorf("CreateSeller: %v\n", err)
		return http.StatusBadRequest, err
	}

	if err := client.Insert(ctx, seller, SellerTableName); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", SellerTableName, err)
		logrus.Errorf("CreateSeller: %v\n", err)
		return http.StatusInternalServerError, err
	}
	logrus.Infof("CreateSeller: Successfully created account for userneme %s\n", seller.UserName)
	return http.StatusOK, nil
}

func (seller *SellerTableModel) GetSellerByID(ctx context.Context) (int, error) {
	var (
		existingUser *SellerTableModel
		err          error
	)

	if existingUser, _, err = seller.getByColumn(ctx, "id", seller.Id); err != nil || existingUser.Id != seller.Id {
		err := fmt.Errorf("unable to find user with with id: %s. %v", seller.Id, err)
		logrus.Errorf("fetchSellerWithSessionID: %v\n", err)
		return http.StatusBadRequest, err
	}

	copySellerObj(existingUser, seller)

	return http.StatusOK, nil
}

func (seller *SellerTableModel) GetSellerByUserName(ctx context.Context) (int, error) {
	var (
		existingUser *SellerTableModel
		err          error
	)

	if existingUser, _, err = seller.getByColumn(ctx, "userName", seller.UserName); err != nil || existingUser.UserName != seller.UserName {
		err := fmt.Errorf("unable to find user with with userName: %s. %v", seller.UserName, err)
		logrus.Errorf("GetSellerByUserName: %v\n", err)
		return http.StatusBadRequest, err
	}

	copySellerObj(existingUser, seller)

	return http.StatusOK, nil
}

func (seller *SellerTableModel) UpdateSellerByID(ctx context.Context) (int, error) {
	client, err := sql.NewSQLClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("UpdateSellerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)
	seller.UpdatedAt = time.Now()

	if err := client.Update(ctx, seller, SellerTableName, true); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", SellerTableName, err)
		logrus.Errorf("UpdateBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (seller *SellerTableModel) getByColumn(ctx context.Context, columnName string, columnValue interface{}) (*SellerTableModel, int, error) {
	client, err := sql.NewSQLClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("UpdateSellerByID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	defer client.Close(ctx)
	whereClause := []db.WhereClauseType{
		{
			ColumnName:   columnName,
			RelationType: db.EQUAL,
			ColumnValue:  columnValue,
		},
	}
	resultSeller := SellerTableModel{}
	if _, err := client.Read(ctx, SellerTableName, nil, whereClause, nil, nil, nil, true, &resultSeller); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", SellerTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	return &resultSeller, http.StatusOK, nil
}

func copySellerObj(from, to *SellerTableModel) {
	to.Id = from.Id
	to.Name = from.Name
	to.FeedBackThumbsUp = from.FeedBackThumbsUp
	to.FeedBackThumbsDown = from.FeedBackThumbsDown
	to.NumberOfItemsSold = from.NumberOfItemsSold
	to.UserName = from.UserName
	to.Password = from.Password
	to.Version = from.Version
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}
