package main

import (
	"context"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/db"
	"net/http"
	"reflect"
	"time"

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
	schema.BaseModel       `bun:"table:buyer_data,alias:buyer"`
	Id                     string    `json:"id" bson:"id" bun:"id,pk"`
	Name                   string    `json:"name" bson:"name" bun:"name,notnull"`
	NumberOfItemsPurchased int       `json:"numberOfItemsPurchased" bson:"numberOfItemsPurchased" bun:"numberOfItemsPurchased"`
	UserName               string    `json:"userName" bson:"userName" bun:"userName,notnull,unique"`
	Password               string    `json:"password" bson:"password" bun:"password,notnull,unique"`
	CreatedAt              time.Time `json:"createdAt"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt" bson:"updatedAt" bun:"updatedAt"`
}

func CreateBuyerTable(ctx context.Context) error {

	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", BuyerTableName, err)
		logrus.Errorf("CreateBuyerTable: %v\n", err)
		return err
	}

	tableSchemaPtr := reflect.New(reflect.TypeOf(BuyerTableModel{}))
	createTableQuery := db.SqlDBClient.NewCreateTable().
		Model(tableSchemaPtr.Interface()).
		IfNotExists()

	_, err := createTableQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("exception while creaiting event table %s. %v", err, BuyerTableName)
		logrus.Errorf("CreateBuyerTable: %v\n", err)
		return err
	}

	return nil
}

func (buyer *BuyerTableModel) CreateBuyer(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("CreateBuyer: %v\n", err)
		return http.StatusInternalServerError, err
	}

	buyer.Id = uuid.New().String()
	buyer.CreatedAt = time.Now()
	buyer.UpdatedAt = time.Now()

	if existingBuyer, _, _ := buyer.getByBuyerUserName(ctx); existingBuyer != nil && existingBuyer.UserName == buyer.UserName {
		err := fmt.Errorf("exception while verifying Buyer data. userName alredy taken")
		logrus.Errorf("CreateBuyer: %v\n", err)
		return http.StatusBadRequest, err
	}

	if _, err := db.SqlDBClient.NewInsert().Model(buyer).Exec(ctx); err != nil {
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

	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("GetBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	if existingUser, _, err = buyer.getByBuyerID(ctx); err != nil || existingUser.Id != buyer.Id {
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

	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("GetBuyerByUserName: %v\n", err)
		return http.StatusInternalServerError, err
	}

	if existingUser, _, err = buyer.getByBuyerUserName(ctx); err != nil || existingUser.UserName != buyer.UserName {
		err := fmt.Errorf("unable to find user with with userName: %s. %v", buyer.UserName, err)
		logrus.Errorf("GetBuyerByUserName: %v\n", err)
		return http.StatusBadRequest, err
	}

	copyBuyerObj(existingUser, buyer)

	return http.StatusOK, nil
}

func (buyer *BuyerTableModel) UpdateBuyerByID(ctx context.Context) (int, error) {
	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("UpdateBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	buyer.UpdatedAt = time.Now()
	oldVersion := 0
	updateQuery := db.PrepareUpdateQuery(ctx, &oldVersion, buyer, false, true)
	logrus.Infof("UpdateBuyerByID: Update query: %v", updateQuery.String())
	_, err := updateQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", BuyerTableName, err)
		logrus.Errorf("UpdateBuyerByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (buyer *BuyerTableModel) getByColumn(ctx context.Context, columnName string, columnValue interface{}) (*BuyerTableModel, int, error) {
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
	resultBuyer := BuyerTableModel{}
	_, statusCode, err := db.ReadUtil(ctx, BuyerTableName, nil, whereClause, nil, nil, nil, true, &resultBuyer)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", BuyerTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, statusCode, err
	}
	return &resultBuyer, http.StatusOK, nil
}

func (buyer *BuyerTableModel) getByBuyerUserName(ctx context.Context) (*BuyerTableModel, int, error) {
	return buyer.getByColumn(ctx, "userName", buyer.UserName)
}

func (buyer *BuyerTableModel) getByBuyerID(ctx context.Context) (*BuyerTableModel, int, error) {
	return buyer.getByColumn(ctx, "id", buyer.Id)
}

func copyBuyerObj(from, to *BuyerTableModel) {
	to.Id = from.Id
	to.Name = from.Name
	to.NumberOfItemsPurchased = from.NumberOfItemsPurchased
	to.UserName = from.UserName
	to.Password = from.Password
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}
