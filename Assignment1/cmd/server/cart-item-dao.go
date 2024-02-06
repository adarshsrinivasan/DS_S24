package main

import (
	"context"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/db"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/db/sql"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun/schema"
)

const (
	CartItemTableName      = "cartitem_data"
	CartItemTableAliasName = "cartitem"
)

type CartItemOps interface {
	CreateCartItem(ctx context.Context) (int, error)
	GetCartItemByID(ctx context.Context) (int, error)
	GetCartItemByCartIDAndProductID(ctx context.Context) (int, error)
	ListCartItemByCartID(ctx context.Context) ([]CartItemTableModel, int, error)
	UpdateCartItem(ctx context.Context) (int, error)
	DeleteCartItemByCartIDAndProductID(ctx context.Context) (int, error)
	DeleteCartItemByCartID(ctx context.Context) (int, error)
}

type CartItemTableModel struct {
	schema.BaseModel `bun:"table:cartitem_data,alias:cartitem"`
	ID               string    `json:"id" bson:"id" bun:"id,pk"`
	CartID           string    `json:"cartID" bson:"cartID" bun:"cartID,notnull"`
	ProductID        string    `json:"productID" bson:"productID" bun:"productID,notnull"`
	SellerID         string    `json:"sellerID" bson:"sellerID" bun:"sellerID,notnull"`
	Quantity         int       `json:"quantity" bson:"quantity" bun:"quantity,notnull"`
	Price            float32   `json:"price" bson:"price,omitempty" bun:"price,notnull"`
	Version          int       `json:"version" bson:"version" bun:"version,notnull"`
	CreatedAt        time.Time `json:"createdAt"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt" bun:"updatedAt"`
}

func CreateCartItemTable(ctx context.Context) error {

	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateCartItemTable: %v\n", err)
		return err
	}
	defer client.Close(ctx)

	tableSchemaPtr := reflect.New(reflect.TypeOf(CartItemTableModel{}))

	foreignKeys := []db.ForeignKey{
		{
			ColumnName:    "cartID",
			SrcColumnName: "id",
			SrcTableName:  CartTableName,
			CascadeDelete: true,
		},
		{
			ColumnName:    "sellerID",
			SrcColumnName: "id",
			SrcTableName:  SellerTableName,
			CascadeDelete: true,
		},
	}

	if err := client.CreateTable(ctx, tableSchemaPtr.Interface(), CartItemTableName, foreignKeys); err != nil {
		err := fmt.Errorf("exception while creating table %s. %v", err, CartItemTableName)
		logrus.Errorf("CreateCartItemTable: %v\n", err)
		return err
	}

	return nil
}

func (cartItem *CartItemTableModel) CreateCartItem(ctx context.Context) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateCartItem: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	cartItem.ID = uuid.New().String()
	cartItem.Version = 0
	cartItem.CreatedAt = time.Now()
	cartItem.UpdatedAt = time.Now()

	if err := client.Insert(ctx, cartItem, CartItemTableName); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", CartItemTableName, err)
		logrus.Errorf("CreateCartItem: %v\n", err)
		return http.StatusInternalServerError, err
	}
	logrus.Infof("CreateCartItem: Successfully created cartItem for productID %s\n", cartItem.ProductID)
	return http.StatusOK, nil
}

func (cartItem *CartItemTableModel) GetCartItemByID(ctx context.Context) (int, error) {
	var (
		existingCartItem *CartItemTableModel
		err              error
	)

	if existingCartItem, _, err = cartItem.getByColumn(ctx, "id", cartItem.ID); err != nil || existingCartItem.ID != cartItem.ID {
		err := fmt.Errorf("unable to find cartItem with with id: %s. %v", cartItem.ID, err)
		logrus.Errorf("GetCartItemByCartID: %v\n", err)
		return http.StatusBadRequest, err
	}

	copyCartItemObj(existingCartItem, cartItem)

	return http.StatusOK, nil
}

func (cartItem *CartItemTableModel) GetCartItemByCartIDAndProductID(ctx context.Context) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("GetCartItemByCartIDAndProductID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "cartID",
			RelationType: db.EQUAL,
			ColumnValue:  cartItem.CartID,
		},
		{
			ColumnName:   "productID",
			RelationType: db.EQUAL,
			ColumnValue:  cartItem.ProductID,
		},
	}
	result := CartItemTableModel{}
	if _, err := client.Read(ctx, CartItemTableName, nil, whereClause, nil, nil, nil, true, &result); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("GetCartItemByCartIDAndProductID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	copyCartItemObj(&result, cartItem)

	return http.StatusOK, nil
}

func (cartItem *CartItemTableModel) ListCartItemByCartID(ctx context.Context) ([]CartItemTableModel, int, error) {

	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("ListCartItemByCartID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "cartID",
			RelationType: db.EQUAL,
			ColumnValue:  cartItem.CartID,
		},
	}
	var result []CartItemTableModel
	if _, err := client.Read(ctx, CartItemTableName, nil, whereClause, nil, nil, nil, true, &result); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("ListCartItemByCartID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func (cartItem *CartItemTableModel) UpdateCartItem(ctx context.Context) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("UpdateCartItem: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	cartItem.UpdatedAt = time.Now()

	if err := client.Update(ctx, cartItem, CartItemTableName, true); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", CartItemTableName, err)
		logrus.Errorf("UpdateCartItem: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (cartItem *CartItemTableModel) DeleteCartItemByCartIDAndProductID(ctx context.Context) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("DeleteCartItemByCartIDAndProductID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	whereClauses := []db.WhereClauseType{
		{
			ColumnName:   "cartID",
			RelationType: db.EQUAL,
			ColumnValue:  cartItem.ID,
		},
		{
			ColumnName:   "productID",
			RelationType: db.EQUAL,
			ColumnValue:  cartItem.ID,
		},
	}

	if err := client.Delete(ctx, cartItem, CartItemTableName, whereClauses); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartItemTableName, err)
		logrus.Errorf("DeleteCartItemByID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (cartItem *CartItemTableModel) DeleteCartItemByCartID(ctx context.Context) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("DeleteCartItemByCartID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)
	whereClauses := []db.WhereClauseType{
		{
			ColumnName:   "cartID",
			RelationType: db.EQUAL,
			ColumnValue:  cartItem.CartID,
		},
	}

	if err := client.Delete(ctx, cartItem, CartItemTableName, whereClauses); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartItemTableName, err)
		logrus.Errorf("DeleteCartItemByCartID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (cartItem *CartItemTableModel) getByColumn(ctx context.Context, columnName string, columnValue interface{}) (*CartItemTableModel, int, error) {
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
	result := CartItemTableModel{}

	if _, err := client.Read(ctx, CartItemTableName, nil, whereClause, nil, nil, nil, true, &result); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	return &result, http.StatusOK, nil
}

func copyCartItemObj(from, to *CartItemTableModel) {
	to.ID = from.ID
	to.CartID = from.CartID
	to.ProductID = from.ProductID
	to.SellerID = from.SellerID
	to.Quantity = from.Quantity
	to.Price = from.Price
	to.Version = from.Version
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}
