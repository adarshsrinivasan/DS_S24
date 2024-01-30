package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/db"
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
	CreatedAt          time.Time `json:"createdAt"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" bson:"updatedAt" bun:"updatedAt"`
}

func CreateSellerTable(ctx context.Context) error {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", SellerTableName, err)
		logrus.Errorf("CreateSellerTable: %v\n", err)
		return err
	}

	tableSchemaPtr := reflect.New(reflect.TypeOf(SellerTableModel{}))
	createTableQuery := db.SqlDBClient.NewCreateTable().
		Model(tableSchemaPtr.Interface()).
		IfNotExists()

	_, err := createTableQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("exception while creaiting event table %s. %v", err, SellerTableName)
		logrus.Errorf("CreateSellerTable: %v\n", err)
		return err
	}

	return nil
}

func (seller *SellerTableModel) CreateSeller(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("CreateSeller: %v\n", err)
		return http.StatusInternalServerError, err
	}

	seller.Id = uuid.New().String()
	seller.CreatedAt = time.Now()
	seller.UpdatedAt = time.Now()

	if existingSeller, _, _ := seller.getBySellerUserName(ctx); existingSeller != nil && existingSeller.UserName == seller.UserName {
		err := fmt.Errorf("exception while verifying Seller data. userName alredy taken")
		logrus.Errorf("CreateSeller: %v\n", err)
		return http.StatusBadRequest, err
	}

	if _, err := db.SqlDBClient.NewInsert().Model(seller).Exec(ctx); err != nil {
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

	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("Logout: %v\n", err)
		return http.StatusInternalServerError, err
	}

	if existingUser, _, err = seller.getBySellerID(ctx); err != nil || existingUser.Id != seller.Id {
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

	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("GetSellerByUserName: %v\n", err)
		return http.StatusInternalServerError, err
	}

	if existingUser, _, err = seller.getBySellerUserName(ctx); err != nil || existingUser.UserName != seller.UserName {
		err := fmt.Errorf("unable to find user with with userName: %s. %v", seller.UserName, err)
		logrus.Errorf("GetSellerByUserName: %v\n", err)
		return http.StatusBadRequest, err
	}

	copySellerObj(existingUser, seller)

	return http.StatusOK, nil
}

func (seller *SellerTableModel) UpdateSellerByID(ctx context.Context) (int, error) {
	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("UpdateBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	seller.UpdatedAt = time.Now()
	oldVersion := 0
	updateQuery := db.PrepareUpdateQuery(ctx, &oldVersion, seller, false, true)
	logrus.Infof("UpdateBuyerByID: Update query: %v", updateQuery.String())
	_, err := updateQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", SellerTableName, err)
		logrus.Errorf("UpdateBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (seller *SellerTableModel) getByColumn(ctx context.Context, columnName string, columnValue interface{}) (*SellerTableModel, int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	whereClause := []db.WhereClauseType{
		{
			ColumnName:   columnName,
			RelationType: db.EQUAL,
			ColumnValue:  columnValue,
		},
	}
	resultSeller := SellerTableModel{}
	_, statusCode, err := db.ReadUtil(ctx, SellerTableName, nil, whereClause, nil, nil, nil, true, &resultSeller)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", SellerTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, statusCode, err
	}
	return &resultSeller, http.StatusOK, nil
}

func (seller *SellerTableModel) getBySellerUserName(ctx context.Context) (*SellerTableModel, int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "userName",
			RelationType: db.EQUAL,
			ColumnValue:  seller.UserName,
		},
	}
	resultSeller := SellerTableModel{}
	_, statusCode, err := db.ReadUtil(ctx, SellerTableName, nil, whereClause, nil, nil, nil, true, &resultSeller)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", SellerTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, statusCode, err
	}
	return &resultSeller, http.StatusOK, nil
	//return seller.getByColumn(ctx, "userName", seller.UserName)
}

func (seller *SellerTableModel) getBySellerID(ctx context.Context) (*SellerTableModel, int, error) {
	return seller.getByColumn(ctx, "id", seller.Id)
}

func copySellerObj(from, to *SellerTableModel) {
	to.Id = from.Id
	to.Name = from.Name
	to.FeedBackThumbsUp = from.FeedBackThumbsUp
	to.FeedBackThumbsDown = from.FeedBackThumbsDown
	to.NumberOfItemsSold = from.NumberOfItemsSold
	to.UserName = from.UserName
	to.Password = from.Password
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}
