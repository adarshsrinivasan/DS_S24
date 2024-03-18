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
	BuyerTableName      = "buyer_data"
	BuyerTableAliasName = "buyer"
)

type BuyerOps interface {
	CreateBuyer(ctx context.Context) (int, error)
	GetBuyerByID(ctx context.Context) (int, error)
	GetBuyerByUserName(ctx context.Context) (int, error)
	UpdateBuyerByID(ctx context.Context) (int, error)
}

type BuyerTableModel struct {
	schema.BaseModel `bun:"table:buyer_data,alias:buyer"`
	Id               string    `json:"id" bson:"id" bun:"id,pk"`
	Name             string    `json:"name" bson:"name" bun:"name,notnull"`
	UserName         string    `json:"userName" bson:"userName" bun:"userName,notnull,unique"`
	Password         string    `json:"password" bson:"password" bun:"password,notnull,unique"`
	Version          int       `json:"version" bson:"version" bun:"version,notnull"`
	CreatedAt        time.Time `json:"createdAt"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt" bun:"updatedAt"`
}

func CreateBuyerTable(ctx context.Context) error {
	client, err := sql.NewSQLClient(ctx, serviceName, schemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateBuyerTable: %v\n", err)
		return err
	}
	defer client.Close(ctx)

	tableSchemaPtr := reflect.New(reflect.TypeOf(BuyerTableModel{}))

	if err := client.CreateTable(ctx, tableSchemaPtr.Interface(), BuyerTableName, nil); err != nil {
		err := fmt.Errorf("exception while creating table %s. %v", err, BuyerTableName)
		logrus.Errorf("CreateBuyerTable: %v\n", err)
		return err
	}

	buyer := BuyerTableModel{
		Name:     "admin",
		UserName: "admin",
		Password: "admin",
	}

	if statusCode, err := buyer.CreateBuyer(ctx); err != nil {
		err := fmt.Errorf("exception while inserting the admin Buyer. %v", err)
		if statusCode == http.StatusBadRequest {
			return nil
		}
		logrus.Errorf("CreateBuyerTable: %v. StatusCode: %v\n", statusCode, err)
		return err
	}

	return nil
}

func (buyer *BuyerTableModel) CreateBuyer(ctx context.Context) (int, error) {
	client, err := sql.NewSQLClient(ctx, serviceName, schemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateBuyer: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	buyer.Id = uuid.New().String()
	buyer.Version = 0
	buyer.CreatedAt = time.Now()
	buyer.UpdatedAt = time.Now()

	if existingBuyer, _, _ := buyer.getByColumn(ctx, "userName", buyer.UserName); existingBuyer != nil && existingBuyer.UserName == buyer.UserName {
		err := fmt.Errorf("exception while verifying Buyer data. userName alredy taken")
		logrus.Errorf("CreateBuyer: %v\n", err)
		return http.StatusBadRequest, err
	}

	if err := client.Insert(ctx, buyer, BuyerTableName); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", BuyerTableName, err)
		logrus.Errorf("CreateBuyer: %v\n", err)
		return http.StatusInternalServerError, err
	}
	logrus.Infof("CreateBuyer: Successfully created account for userneme %s\n", buyer.UserName)
	return http.StatusOK, nil
}

func (buyer *BuyerTableModel) GetBuyerByID(ctx context.Context) (int, error) {
	var (
		existingUser *BuyerTableModel
		err          error
	)

	if existingUser, _, err = buyer.getByColumn(ctx, "id", buyer.Id); err != nil || existingUser.Id != buyer.Id {
		err := fmt.Errorf("unable to find user with with id: %s. %v", buyer.Id, err)
		logrus.Errorf("GetBuyerByID: %v\n", err)
		return http.StatusBadRequest, err
	}

	copyBuyerObj(existingUser, buyer)

	return http.StatusOK, nil
}

func (buyer *BuyerTableModel) GetBuyerByUserName(ctx context.Context) (int, error) {
	var (
		existingUser *BuyerTableModel
		err          error
	)

	if existingUser, _, err = buyer.getByColumn(ctx, "userName", buyer.UserName); err != nil || existingUser.UserName != buyer.UserName {
		err := fmt.Errorf("unable to find user with with userName: %s. %v", buyer.UserName, err)
		logrus.Errorf("GetBuyerByUserName: %v\n", err)
		return http.StatusBadRequest, err
	}

	copyBuyerObj(existingUser, buyer)

	return http.StatusOK, nil
}

func (buyer *BuyerTableModel) UpdateBuyerByID(ctx context.Context) (int, error) {
	client, err := sql.NewSQLClient(ctx, serviceName, schemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("UpdateBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	buyer.UpdatedAt = time.Now()

	if err := client.Update(ctx, buyer, BuyerTableName, true); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", BuyerTableName, err)
		logrus.Errorf("UpdateBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (buyer *BuyerTableModel) getByColumn(ctx context.Context, columnName string, columnValue interface{}) (*BuyerTableModel, int, error) {
	client, err := sql.NewSQLClient(ctx, serviceName, schemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("getByColumn: %v\n", err)
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
	resultBuyer := BuyerTableModel{}

	if _, err := client.Read(ctx, BuyerTableName, nil, whereClause, nil, nil, nil, true, &resultBuyer); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", BuyerTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	return &resultBuyer, http.StatusOK, nil
}

func copyBuyerObj(from, to *BuyerTableModel) {
	to.Id = from.Id
	to.Name = from.Name
	to.UserName = from.UserName
	to.Password = from.Password
	to.Version = from.Version
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}
