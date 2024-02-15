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
	Version          int       `json:"version" bson:"version" bun:"version,notnull"`
	CreatedAt        time.Time `json:"createdAt"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt" bun:"updatedAt"`
}

func CreateCartTable(ctx context.Context) error {

	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateCartTable: %v\n", err)
		return err
	}
	defer client.Close(ctx)

	tableSchemaPtr := reflect.New(reflect.TypeOf(CartTableModel{}))

	foreignKeys := []db.ForeignKey{
		{
			ColumnName:    "buyerID",
			SrcColumnName: "id",
			SrcTableName:  BuyerTableName,
			CascadeDelete: true,
		},
	}

	if err := client.CreateTable(ctx, tableSchemaPtr.Interface(), CartTableName, foreignKeys); err != nil {
		err := fmt.Errorf("exception while creating table %s. %v", err, CartTableName)
		logrus.Errorf("CreateCartTable: %v\n", err)
		return err
	}

	return nil
}

func (cart *CartTableModel) CreateCart(ctx context.Context) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateCart: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	cart.ID = uuid.New().String()
	cart.Version = 0
	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()

	if err := client.Insert(ctx, cart, CartTableName); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", BuyerTableName, err)
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

	if existingCart, _, err = cart.getByColumn(ctx, "id", cart.ID); err != nil || existingCart.ID != cart.ID {
		err := fmt.Errorf("unable to find cart with with id: %s. %v", cart.ID, err)
		logrus.Errorf("GetCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	copyCartObj(existingCart, cart)

	return http.StatusOK, nil
}

func (cart *CartTableModel) GetCartByBuyerID(ctx context.Context) (int, error) {
	var (
		existingCart *CartTableModel
		err          error
	)

	if existingCart, _, err = cart.getByColumn(ctx, "buyerID", cart.BuyerID); err != nil || existingCart.BuyerID != cart.BuyerID {
		err := fmt.Errorf("unable to find cart with with buyerID: %s. %v", cart.BuyerID, err)
		logrus.Errorf("GetCartByBuyerID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	copyCartObj(existingCart, cart)

	return http.StatusOK, nil
}

func (cart *CartTableModel) UpdateCartByID(ctx context.Context) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("UpdateCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	cart.UpdatedAt = time.Now()

	if err := client.Update(ctx, cart, CartTableName, true); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", CartTableName, err)
		logrus.Errorf("UpdateCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (cart *CartTableModel) DeleteCartByID(ctx context.Context) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("DeleteCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	whereClauses := []db.WhereClauseType{
		{
			ColumnName:   "id",
			RelationType: db.EQUAL,
			ColumnValue:  cart.ID,
		},
	}

	if err := client.Delete(ctx, cart, CartTableName, whereClauses); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartTableName, err)
		logrus.Errorf("DeleteCartByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (cart *CartTableModel) getByColumn(ctx context.Context, columnName string, columnValue interface{}) (*CartTableModel, int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
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
	result := CartTableModel{}

	if _, err := client.Read(ctx, CartTableName, nil, whereClause, nil, nil, nil, true, &result); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	return &result, http.StatusOK, nil
}

func copyCartObj(from, to *CartTableModel) {
	to.ID = from.ID
	to.BuyerID = from.BuyerID
	to.Saved = from.Saved
	to.Version = from.Version
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}
