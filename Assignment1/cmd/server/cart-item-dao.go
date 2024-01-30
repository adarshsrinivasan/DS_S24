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
	CreatedAt        time.Time `json:"createdAt"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt" bun:"updatedAt"`
}

func CreateCartItemTable(ctx context.Context) error {

	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", CartItemTableName, err)
		logrus.Errorf("CreateCartItemTable: %v\n", err)
		return err
	}

	tableSchemaPtr := reflect.New(reflect.TypeOf(CartItemTableModel{}))
	createTableQuery := db.SqlDBClient.NewCreateTable().
		Model(tableSchemaPtr.Interface()).
		IfNotExists().
		ForeignKey(`("cartID") REFERENCES "cart_data" ("id") ON DELETE CASCADE`).
		ForeignKey(`("sellerID") REFERENCES "seller_data" ("id") ON DELETE CASCADE`)

	_, err := createTableQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("exception while creaiting event table %s. %v", err, CartItemTableName)
		logrus.Errorf("CreateCartItemTable: %v\n", err)
		return err
	}

	return nil
}

func (cartItem *CartItemTableModel) CreateCartItem(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("CreateCartItem: %v\n", err)
		return http.StatusInternalServerError, err
	}

	cartItem.ID = uuid.New().String()
	cartItem.CreatedAt = time.Now()
	cartItem.UpdatedAt = time.Now()

	if _, err := db.SqlDBClient.NewInsert().Model(cartItem).Exec(ctx); err != nil {
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

	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("GetCartItemByCartID: %v\n", err)
		return http.StatusInternalServerError, err
	}

	if existingCartItem, _, err = cartItem.getByCartItemID(ctx); err != nil || existingCartItem.ID != cartItem.ID {
		err := fmt.Errorf("unable to find cartItem with with id: %s. %v", cartItem.ID, err)
		logrus.Errorf("GetCartItemByCartID: %v\n", err)
		return http.StatusBadRequest, err
	}

	copyCartItemObj(existingCartItem, cartItem)

	return http.StatusOK, nil
}

func (cartItem *CartItemTableModel) GetCartItemByCartIDAndProductID(ctx context.Context) (int, error) {
	var (
		err error
	)

	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("GetCartItemByCartIDAndProductID: %v\n", err)
		return http.StatusInternalServerError, err
	}

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
	_, statusCode, err := db.ReadUtil(ctx, CartItemTableName, nil, whereClause, nil, nil, nil, true, &result)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("GetCartItemByCartIDAndProductID: %v\n", err)
		return statusCode, err
	}

	copyCartItemObj(&result, cartItem)

	return http.StatusOK, nil
}

func (cartItem *CartItemTableModel) ListCartItemByCartID(ctx context.Context) ([]CartItemTableModel, int, error) {

	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("GetCartItemByCartID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}

	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "cartID",
			RelationType: db.EQUAL,
			ColumnValue:  cartItem.CartID,
		},
	}
	var result []CartItemTableModel
	_, statusCode, err := db.ReadUtil(ctx, CartItemTableName, nil, whereClause, nil, nil, nil, true, &result)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, statusCode, err
	}

	return result, http.StatusOK, nil
}

func (cartItem *CartItemTableModel) UpdateCartItem(ctx context.Context) (int, error) {
	if err = db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("UpdateCartItem: %v\n", err)
		return http.StatusInternalServerError, err
	}
	cartItem.UpdatedAt = time.Now()
	oldVersion := 0
	updateQuery := db.PrepareUpdateQuery(ctx, &oldVersion, cartItem, false, true)
	logrus.Infof("UpdateCartItem: Update query: %v", updateQuery.String())
	_, err := updateQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", CartItemTableName, err)
		logrus.Errorf("UpdateCartItem: %v\n", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (cartItem *CartItemTableModel) DeleteCartItemByCartIDAndProductID(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("DeleteCartItemByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
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

	deleteQuery := db.SqlDBClient.NewDelete().
		Model(cartItem)

	// prepare whereClause.
	queryStr, vals, err := db.CreateWhereClause(ctx, whereClauses)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartItemTableName, err)
		logrus.Errorf("DeleteCartItemByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	deleteQuery = deleteQuery.Where(queryStr, vals...)

	if _, err := deleteQuery.Exec(ctx); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartItemTableName, err)
		logrus.Errorf("DeleteCartItemByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (cartItem *CartItemTableModel) DeleteCartItemByCartID(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("DeleteCartItemByCartID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	whereClauses := []db.WhereClauseType{
		{
			ColumnName:   "cartID",
			RelationType: db.EQUAL,
			ColumnValue:  cartItem.CartID,
		},
	}

	deleteQuery := db.SqlDBClient.NewDelete().
		Model(cartItem)

	// prepare whereClause.
	queryStr, vals, err := db.CreateWhereClause(ctx, whereClauses)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartItemTableName, err)
		logrus.Errorf("DeleteCartItemByCartID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	deleteQuery = deleteQuery.Where(queryStr, vals...)

	if _, err := deleteQuery.Exec(ctx); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", CartItemTableName, err)
		logrus.Errorf("DeleteCartItemByCartID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (cartItem *CartItemTableModel) getByColumn(ctx context.Context, columnName string, columnValue interface{}) (*CartItemTableModel, int, error) {
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
	result := CartItemTableModel{}
	_, statusCode, err := db.ReadUtil(ctx, CartItemTableName, nil, whereClause, nil, nil, nil, true, &result)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartItemTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, statusCode, err
	}
	return &result, http.StatusOK, nil
}

func (cartItem *CartItemTableModel) getByCartItemID(ctx context.Context) (*CartItemTableModel, int, error) {
	return cartItem.getByColumn(ctx, "id", cartItem.ID)
}

func copyCartItemObj(from, to *CartItemTableModel) {
	to.ID = from.ID
	to.CartID = from.CartID
	to.ProductID = from.ProductID
	to.SellerID = from.SellerID
	to.Quantity = from.Quantity
	to.Price = from.Price
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}
