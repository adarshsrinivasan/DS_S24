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
	CreatedAt        time.Time `json:"createdAt"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt" bun:"updatedAt"`
}

func CreateTransactionTable(ctx context.Context) error {

	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", TransactionTableName, err)
		logrus.Errorf("CreateTransactionTable: %v\n", err)
		return err
	}

	tableSchemaPtr := reflect.New(reflect.TypeOf(TransactionTableModel{}))
	createTableQuery := db.SqlDBClient.NewCreateTable().
		Model(tableSchemaPtr.Interface()).
		IfNotExists().
		ForeignKey(`("cartID") REFERENCES "cart_data" ("id") ON DELETE CASCADE`).
		ForeignKey(`("sellerID") REFERENCES "seller_data" ("id") ON DELETE CASCADE`).
		ForeignKey(`("buyerID") REFERENCES "buyer_data" ("id") ON DELETE CASCADE`)

	_, err := createTableQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("exception while creaiting event table %s. %v", err, TransactionTableName)
		logrus.Errorf("CreateTransactionTable: %v\n", err)
		return err
	}

	return nil
}

func (transaction *TransactionTableModel) CreateTransaction(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("CreateTransaction: %v\n", err)
		return http.StatusInternalServerError, err
	}

	transaction.ID = uuid.New().String()
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	if _, err := db.SqlDBClient.NewInsert().Model(transaction).Exec(ctx); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", TransactionTableName, err)
		logrus.Errorf("CreateTransaction: %v\n", err)
		return http.StatusInternalServerError, err
	}
	logrus.Infof("CreateTransaction: Successfully Recorded transaction for product %s Cart %s Seller %s Buyer %s %s\n", transaction.ProductID, transaction.CartID, transaction.SellerID, transaction.BuyerID)
	return http.StatusOK, nil
}

func (transaction *TransactionTableModel) ListTransactionsByCartID(ctx context.Context) ([]TransactionTableModel, int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("ListTransactionsByCartID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}

	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "cartID",
			RelationType: db.EQUAL,
			ColumnValue:  transaction.CartID,
		},
	}
	var result []TransactionTableModel
	_, statusCode, err := db.ReadUtil(ctx, TransactionTableName, nil, whereClause, nil, nil, nil, true, &result)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("ListTransactionsByCartID: %v\n", err)
		return nil, statusCode, err
	}

	return result, http.StatusOK, nil
}

func (transaction *TransactionTableModel) ListTransactionsByBuyerID(ctx context.Context) ([]TransactionTableModel, int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("ListTransactionsByBuyerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}

	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "buyerID",
			RelationType: db.EQUAL,
			ColumnValue:  transaction.BuyerID,
		},
	}
	var result []TransactionTableModel
	_, statusCode, err := db.ReadUtil(ctx, TransactionTableName, nil, whereClause, nil, nil, nil, true, &result)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("ListTransactionsByBuyerID: %v\n", err)
		return nil, statusCode, err
	}

	return result, http.StatusOK, nil
}

func (transaction *TransactionTableModel) ListTransactionsBySellerID(ctx context.Context) ([]TransactionTableModel, int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("ListTransactionsBySellerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}

	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "sellerID",
			RelationType: db.EQUAL,
			ColumnValue:  transaction.SellerID,
		},
	}
	var result []TransactionTableModel
	_, statusCode, err := db.ReadUtil(ctx, TransactionTableName, nil, whereClause, nil, nil, nil, true, &result)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("ListTransactionsBySellerID: %v\n", err)
		return nil, statusCode, err
	}

	return result, http.StatusOK, nil
}

func (transaction *TransactionTableModel) DeleteTransactionsByCartID(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("DeleteTransactionsByCartID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	if statusCode, err := transaction.deleteTransactionByCartID(ctx); err != nil {
		err := fmt.Errorf("unable to delete Transaction with with cartID: %s. %v", transaction.CartID, err)
		logrus.Errorf("DeleteTransactionsByCartID: %v\n", err)
		return statusCode, err
	}

	return http.StatusOK, nil
}

func (transaction *TransactionTableModel) DeleteTransactionsBySellerID(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("DeleteTransactionsBySellerID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	if statusCode, err := transaction.deleteTransactionBySellerID(ctx); err != nil {
		err := fmt.Errorf("unable to delete Transaction with with sellerID: %s. %v", transaction.SellerID, err)
		logrus.Errorf("DeleteTransactionsBySellerID: %v\n", err)
		return statusCode, err
	}

	return http.StatusOK, nil
}

func (transaction *TransactionTableModel) DeleteTransactionsByBuyerID(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("DeleteTransactionsByBuyerID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	if statusCode, err := transaction.deleteTransactionByBuyerID(ctx); err != nil {
		err := fmt.Errorf("unable to delete Transaction with with buyerID: %s. %v", transaction.BuyerID, err)
		logrus.Errorf("DeleteTransactionsByBuyerID: %v\n", err)
		return statusCode, err
	}

	return http.StatusOK, nil
}

func (transaction *TransactionTableModel) getByColumn(ctx context.Context, columnName string, columnValue interface{}) (*TransactionTableModel, int, error) {
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
	result := TransactionTableModel{}
	_, statusCode, err := db.ReadUtil(ctx, TransactionTableName, nil, whereClause, nil, nil, nil, true, &result)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", TransactionTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, statusCode, err
	}
	return &result, http.StatusOK, nil
}

func (transaction *TransactionTableModel) getTransactionByCartID(ctx context.Context) (*TransactionTableModel, int, error) {
	return transaction.getByColumn(ctx, "cartID", transaction.CartID)
}

func (transaction *TransactionTableModel) getTransactionBySellerID(ctx context.Context) (*TransactionTableModel, int, error) {
	return transaction.getByColumn(ctx, "sellerID", transaction.SellerID)
}

func (transaction *TransactionTableModel) getTransactionByBuyerID(ctx context.Context) (*TransactionTableModel, int, error) {
	return transaction.getByColumn(ctx, "buyerID", transaction.BuyerID)
}

func (transaction *TransactionTableModel) deleteByColumn(ctx context.Context, columnName string, columnValue interface{}) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("deleteByColumn: %v\n", err)
		return http.StatusInternalServerError, err
	}
	whereClauses := []db.WhereClauseType{
		{
			ColumnName:   columnName,
			RelationType: db.EQUAL,
			ColumnValue:  columnValue,
		},
	}

	deleteQuery := db.SqlDBClient.NewDelete().
		Model(transaction)

	// prepare whereClause.
	queryStr, vals, err := db.CreateWhereClause(ctx, whereClauses)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", TransactionTableName, err)
		logrus.Errorf("deleteByColumn: %v\n", err)
		return http.StatusInternalServerError, err
	}
	deleteQuery = deleteQuery.Where(queryStr, vals...)

	if _, err := deleteQuery.Exec(ctx); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", TransactionTableName, err)
		logrus.Errorf("deleteByColumn: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (transaction *TransactionTableModel) deleteTransactionByCartID(ctx context.Context) (int, error) {
	return transaction.deleteByColumn(ctx, "cartID", transaction.CartID)
}

func (transaction *TransactionTableModel) deleteTransactionBySellerID(ctx context.Context) (int, error) {
	return transaction.deleteByColumn(ctx, "sellerID", transaction.SellerID)
}

func (transaction *TransactionTableModel) deleteTransactionByBuyerID(ctx context.Context) (int, error) {
	return transaction.deleteByColumn(ctx, "buyerID", transaction.BuyerID)
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
