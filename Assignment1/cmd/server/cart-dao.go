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

type CartTableModel struct {
	schema.BaseModel `bun:"table:cart_data,alias:cart"`
	ID               string    `json:"id" bson:"id" bun:"id,pk"`
	BuyerID          string    `json:"buyerID" bson:"buyerID" bun:"buyerID,notnull,unique"`
	Saved            bool      `json:"saved" bson:"saved" bun:"saved,notnull"`
	CreatedAt        time.Time `json:"createdAt"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt" bun:"updatedAt"`
}

func CreateCartTable(ctx context.Context) error {

	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", CartTableName, err)
		logrus.Errorf("CreateCartTable: %v\n", err)
		return err
	}

	tableSchemaPtr := reflect.New(reflect.TypeOf(CartTableModel{}))
	createTableQuery := db.SqlDBClient.NewCreateTable().
		Model(tableSchemaPtr.Interface()).
		IfNotExists().
		ForeignKey(`("buyerID") REFERENCES "buyer_data" ("id") ON DELETE CASCADE`)

	_, err := createTableQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("exception while creaiting event table %s. %v", err, CartTableName)
		logrus.Errorf("CreateCartTable: %v\n", err)
		return err
	}

	return nil
}

func (cart *CartTableModel) CreateCart(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("CreateCart: %v\n", err)
		return http.StatusInternalServerError, err
	}

	cart.ID = uuid.New().String()
	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()

	if _, err := db.SqlDBClient.NewInsert().Model(cart).Exec(ctx); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", CartTableName, err)
		logrus.Errorf("CreateCart: %v\n", err)
		return http.StatusInternalServerError, err
	}
	logrus.Infof("CreateCart: Successfully created cart for buyerID %s\n", cart.BuyerID)
	return http.StatusOK, nil
}

func (cart *CartTableModel) GetCartByID(ctx context.Context) (int, error) {
	var (
		existingCart *CartTableModel
		err          error
	)

	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("GetCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	if existingCart, _, err = cart.getByCartID(ctx); err != nil || existingCart.ID != cart.ID {
		err := fmt.Errorf("unable to find cart with with id: %s. %v", cart.ID, err)
		logrus.Errorf("GetCartByID: %v\n", err)
		return http.StatusBadRequest, err
	}

	copyCartObj(existingCart, cart)

	return http.StatusOK, nil
}

func (cart *CartTableModel) GetCartByBuyerID(ctx context.Context) (int, error) {
	var (
		existingCart *CartTableModel
		err          error
	)

	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("GetCartByBuyerID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	if existingCart, _, err = cart.getByCartBuyerID(ctx); err != nil || existingCart.BuyerID != cart.BuyerID {
		err := fmt.Errorf("unable to find cart with with buyerID: %s. %v", cart.BuyerID, err)
		logrus.Errorf("GetCartByBuyerID: %v\n", err)
		return http.StatusBadRequest, err
	}

	copyCartObj(existingCart, cart)

	return http.StatusOK, nil
}

func (cart *CartTableModel) UpdateCartByID(ctx context.Context) (int, error) {
	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("UpdateCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	cart.UpdatedAt = time.Now()
	oldVersion := 0
	updateQuery := db.PrepareUpdateQuery(ctx, &oldVersion, cart, false, true)
	logrus.Infof("UpdateCartByID: Update query: %v", updateQuery.String())
	_, err := updateQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", CartTableName, err)
		logrus.Errorf("UpdateCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (cart *CartTableModel) DeleteCartByID(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("DeleteCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	whereClauses := []db.WhereClauseType{
		{
			ColumnName:   "id",
			RelationType: db.EQUAL,
			ColumnValue:  cart.ID,
		},
	}

	deleteQuery := db.SqlDBClient.NewDelete().
		Model(cart)

	// prepare whereClause.
	queryStr, vals, err := db.CreateWhereClause(ctx, whereClauses)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartTableName, err)
		logrus.Errorf("DeleteCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	deleteQuery = deleteQuery.Where(queryStr, vals...)

	if _, err := deleteQuery.Exec(ctx); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartTableName, err)
		logrus.Errorf("DeleteCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (cart *CartTableModel) getByColumn(ctx context.Context, columnName string, columnValue interface{}) (*CartTableModel, int, error) {
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
	result := CartTableModel{}
	_, statusCode, err := db.ReadUtil(ctx, CartTableName, nil, whereClause, nil, nil, nil, true, &result)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, statusCode, err
	}
	return &result, http.StatusOK, nil
}

func (cart *CartTableModel) getByCartID(ctx context.Context) (*CartTableModel, int, error) {
	return cart.getByColumn(ctx, "id", cart.ID)
}

func (cart *CartTableModel) getByCartBuyerID(ctx context.Context) (*CartTableModel, int, error) {
	return cart.getByColumn(ctx, "buyerID", cart.BuyerID)
}

func copyCartObj(from, to *CartTableModel) {
	to.ID = from.ID
	to.BuyerID = from.BuyerID
	to.Saved = from.Saved
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}
