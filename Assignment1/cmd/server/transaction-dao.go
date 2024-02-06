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
	TransactionTableName      = "transaction_data"
	TransactionTableAliasName = "transaction"
)

type TransactionTableOps interface {
	CreateTransaction(ctx context.Context) (int, error)
	ListTransactionsBySellerID(ctx context.Context) ([]TransactionTableModel, int, error)
	ListTransactionsByBuyerID(ctx context.Context) ([]TransactionTableModel, int, error)
	ListTransactionsByCartID(ctx context.Context) ([]TransactionTableModel, int, error)
	DeleteTransactionsByCartID(ctx context.Context) (int, error)
	DeleteTransactionsByBuyerID(ctx context.Context) (int, error)
	DeleteTransactionsBySellerID(ctx context.Context) (int, error)
}

type TransactionTableModel struct {
	schema.BaseModel `bun:"table:transaction_data,alias:transaction"`
	ID               string    `json:"id" bson:"id" bun:"id,pk"`
	CartID           string    `json:"cartID" bson:"cartID" bun:"cartID,notnull"`
	ProductID        string    `json:"productID" bson:"productID" bun:"productID,notnull"`
	BuyerID          string    `json:"buyerID" bson:"buyerID" bun:"buyerID,notnull"`
	SellerID         string    `json:"sellerID" bson:"sellerID" bun:"sellerID,notnull"`
	Quantity         int       `json:"quantity" bson:"quantity" bun:"quantity,notnull"`
	Price            float32   `json:"price" bson:"price,omitempty" bun:"quantity,notnull"`
	Version          int       `json:"version" bson:"version" bun:"version,notnull"`
	CreatedAt        time.Time `json:"createdAt"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt" bun:"updatedAt"`
}

func CreateTransactionTable(ctx context.Context) error {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateTransactionTable: %v\n", err)
		return err
	}
	defer client.Close(ctx)

	tableSchemaPtr := reflect.New(reflect.TypeOf(TransactionTableModel{}))

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
		{
			ColumnName:    "buyerID",
			SrcColumnName: "id",
			SrcTableName:  BuyerTableName,
			CascadeDelete: true,
		},
	}

	if err := client.CreateTable(ctx, tableSchemaPtr.Interface(), TransactionTableName, foreignKeys); err != nil {
		err := fmt.Errorf("exception while creating table %s. %v", err, TransactionTableName)
		logrus.Errorf("CreateTransactionTable: %v\n", err)
		return err
	}

	return nil
}

func (transaction *TransactionTableModel) CreateTransaction(ctx context.Context) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateTransaction: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	transaction.ID = uuid.New().String()
	transaction.Version = 0
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	if err := client.Insert(ctx, transaction, TransactionTableName); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", TransactionTableName, err)
		logrus.Errorf("CreateTransaction: %v\n", err)
		return http.StatusInternalServerError, err
	}
	logrus.Infof("CreateTransaction: Successfully Recorded transaction for product %s Cart %s Seller %s Buyer %s %s\n", transaction.ProductID, transaction.CartID, transaction.SellerID, transaction.BuyerID)
	return http.StatusOK, nil
}

func (transaction *TransactionTableModel) ListTransactionsByCartID(ctx context.Context) ([]TransactionTableModel, int, error) {
	return transaction.listByColumn(ctx, "cartID", transaction.CartID)
}

func (transaction *TransactionTableModel) ListTransactionsByBuyerID(ctx context.Context) ([]TransactionTableModel, int, error) {
	return transaction.listByColumn(ctx, "buyerID", transaction.BuyerID)
}

func (transaction *TransactionTableModel) ListTransactionsBySellerID(ctx context.Context) ([]TransactionTableModel, int, error) {
	return transaction.listByColumn(ctx, "sellerID", transaction.SellerID)
}

func (transaction *TransactionTableModel) DeleteTransactionsByCartID(ctx context.Context) (int, error) {
	return transaction.deleteByColumn(ctx, "cartID", transaction.CartID)
}

func (transaction *TransactionTableModel) DeleteTransactionsBySellerID(ctx context.Context) (int, error) {
	return transaction.deleteByColumn(ctx, "sellerID", transaction.SellerID)
}

func (transaction *TransactionTableModel) DeleteTransactionsByBuyerID(ctx context.Context) (int, error) {
	return transaction.deleteByColumn(ctx, "buyerID", transaction.BuyerID)
}

func (transaction *TransactionTableModel) listByColumn(ctx context.Context, columnName string, columnValue interface{}) ([]TransactionTableModel, int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("listByColumn: %v\n", err)
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
	var result []TransactionTableModel
	if _, err := client.Read(ctx, TransactionTableName, nil, whereClause, nil, nil, nil, false, &result); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("ListTransactionsBySellerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

func (transaction *TransactionTableModel) deleteByColumn(ctx context.Context, columnName string, columnValue interface{}) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("deleteByColumn: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)
	whereClauses := []db.WhereClauseType{
		{
			ColumnName:   columnName,
			RelationType: db.EQUAL,
			ColumnValue:  columnValue,
		},
	}

	if err := client.Delete(ctx, transaction, TransactionTableName, whereClauses); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", TransactionTableName, err)
		logrus.Errorf("deleteByColumn: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func copyTransactionObj(from, to *TransactionTableModel) {
	to.ID = from.ID
	to.CartID = from.CartID
	to.BuyerID = from.BuyerID
	to.SellerID = from.SellerID
	to.Quantity = from.Quantity
	to.Price = from.Price
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}
