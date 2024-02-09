package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adarshsrinivasan/DS_S24/libraries/db"
	"github.com/adarshsrinivasan/DS_S24/libraries/db/nosql"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	ProductTableName = "product_data"
)

type CATEGORY int

const (
	ZERO CATEGORY = iota
	ONE
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
)

var CategoryToString = map[CATEGORY]string{
	ZERO:  "ZERO",
	ONE:   "ONE",
	TWO:   "TWO",
	THREE: "THREE",
	FOUR:  "FOUR",
	FIVE:  "FIVE",
	SIX:   "SIX",
	SEVEN: "SEVEN",
	EIGHT: "EIGHT",
	NINE:  "NINE",
}

var StringToCategory = map[string]CATEGORY{
	"ZERO":  ZERO,
	"ONE":   ONE,
	"TWO":   TWO,
	"THREE": THREE,
	"FOUR":  FOUR,
	"FIVE":  FIVE,
	"SIX":   SIX,
	"SEVEN": SEVEN,
	"EIGHT": EIGHT,
	"NINE":  NINE,
}

type CONDITION int

const (
	NEW CONDITION = iota
	USED
)

var ConditionToString = map[CONDITION]string{
	NEW:  "NEW",
	USED: "USED",
}

var StringToCondition = map[string]CONDITION{
	"NEW":  NEW,
	"USED": USED,
}

type ProductTableModel struct {
	ID                 string    `json:"id" bson:"_id,omitempty"`
	Name               string    `json:"name" bson:"name,omitempty"`
	Category           CATEGORY  `json:"category" bson:"category,omitempty"`
	Keywords           []string  `json:"keywords" bson:"keywords,omitempty"`
	Condition          CONDITION `json:"condition" bson:"condition,omitempty"`
	SalePrice          float32   `json:"salePrice" bson:"salePrice,omitempty"`
	SellerID           string    `json:"sellerID" bson:"sellerID,omitempty"`
	Quantity           int       `json:"quantity" bson:"quantity"`
	FeedBackThumbsUp   int       `json:"feedBackThumbsUp" bson:"feedBackThumbsUp"`
	FeedBackThumbsDown int       `json:"feedBackThumbsDown" bson:"feedBackThumbsDown"`
	CreatedAt          time.Time `json:"createdAt"  bson:"createdAt,omitempty"`
	UpdatedAt          time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type ProductTableOps interface {
	CreateProduct(ctx context.Context) (int, error)
	GetProductByID(ctx context.Context) (int, error)
	GetProductsByKeyWordsAndCategory(ctx context.Context) ([]ProductTableModel, int, error)
	GetProductsBySellerID(ctx context.Context) ([]ProductTableModel, int, error)
	UpdateProductByID(ctx context.Context) (int, error)
	DeleteProductByID(ctx context.Context) (int, error)
}

func CreateProductTable(ctx context.Context) error {

	if err := nosql.VerifyNOSQLDatabaseConnection(ctx, nosql.Client); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", ProductTableName, err)
		logrus.Errorf("CreateProductTable: %v\n", err)
		return err
	}

	return nosql.Client.CreateCollection(ctx, ProductTableName)
}

func (product *ProductTableModel) CreateProduct(ctx context.Context) (int, error) {
	if err := nosql.VerifyNOSQLDatabaseConnection(ctx, nosql.Client); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", ProductTableName, err)
		logrus.Errorf("CreateProduct: %v\n", err)
		return http.StatusInternalServerError, err
	}

	product.ID = uuid.New().String()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	return nosql.Client.InsertOne(ctx, ProductTableName, *product)
}
func (product *ProductTableModel) GetProductByID(ctx context.Context) (int, error) {
	if err := nosql.VerifyNOSQLDatabaseConnection(ctx, nosql.Client); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", ProductTableName, err)
		logrus.Errorf("GetProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "_id",
			RelationType: db.EQUAL,
			ColumnValue:  product.ID,
		},
	}
	var result ProductTableModel

	if statusCode, err := nosql.Client.FindOne(ctx, ProductTableName, whereClause, &result); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", ProductTableName, err)
		logrus.Errorf("GetProductByID: %v\n", err)
		return statusCode, err
	}
	copyProductTableModelObject(&result, product)
	return http.StatusOK, nil

}

func (product *ProductTableModel) GetProductsByKeyWordsAndCategory(ctx context.Context) ([]ProductTableModel, int, error) {
	if err := nosql.VerifyNOSQLDatabaseConnection(ctx, nosql.Client); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", ProductTableName, err)
		logrus.Errorf("GetProductsByKeyWordsAndCategory: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "keywords",
			RelationType: db.IN,
			ColumnValue:  product.Keywords,
		},
		{
			ColumnName:   "category",
			RelationType: db.EQUAL,
			ColumnValue:  product.Category,
		},
	}
	var result []ProductTableModel

	if statusCode, err := nosql.Client.FindMany(ctx, ProductTableName, whereClause, &result); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", ProductTableName, err)
		logrus.Errorf("GetProductsByKeyWordsAndCategory: %v\n", err)
		return nil, statusCode, err
	}
	return result, http.StatusOK, nil
}

func (product *ProductTableModel) GetProductsBySellerID(ctx context.Context) ([]ProductTableModel, int, error) {
	if err := nosql.VerifyNOSQLDatabaseConnection(ctx, nosql.Client); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", ProductTableName, err)
		logrus.Errorf("GetProductBySellerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "sellerID",
			RelationType: db.EQUAL,
			ColumnValue:  product.SellerID,
		},
	}
	var result []ProductTableModel

	if statusCode, err := nosql.Client.FindMany(ctx, ProductTableName, whereClause, &result); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", ProductTableName, err)
		logrus.Errorf("GetProductsBySellerID: %v\n", err)
		return nil, statusCode, err
	}
	return result, http.StatusOK, nil
}
func (product *ProductTableModel) UpdateProductByID(ctx context.Context) (int, error) {
	if err := nosql.VerifyNOSQLDatabaseConnection(ctx, nosql.Client); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", ProductTableName, err)
		logrus.Errorf("UpdateProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "_id",
			RelationType: db.EQUAL,
			ColumnValue:  product.ID,
		},
	}
	product.UpdatedAt = time.Now()
	if statusCode, err := nosql.Client.UpdateOne(ctx, ProductTableName, whereClause, *product); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", ProductTableName, err)
		logrus.Errorf("GetProductBySellerID: %v\n", err)
		return statusCode, err
	}
	return http.StatusOK, nil
}
func (product *ProductTableModel) DeleteProductByID(ctx context.Context) (int, error) {
	if err := nosql.VerifyNOSQLDatabaseConnection(ctx, nosql.Client); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", ProductTableName, err)
		logrus.Errorf("DeleteProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	whereClause := []db.WhereClauseType{
		{
			ColumnName:   "_id",
			RelationType: db.EQUAL,
			ColumnValue:  product.ID,
		},
	}
	product.UpdatedAt = time.Now()
	if statusCode, err := nosql.Client.DeleteOne(ctx, ProductTableName, whereClause); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", ProductTableName, err)
		logrus.Errorf("DeleteProductByID: %v\n", err)
		return statusCode, err
	}
	return http.StatusOK, nil
}

func copyProductTableModelObject(from, to *ProductTableModel) {
	to.ID = from.ID
	to.Name = from.Name
	to.Category = from.Category
	to.Keywords = from.Keywords
	to.Condition = from.Condition
	to.SalePrice = from.SalePrice
	to.SellerID = from.SellerID
	to.Quantity = from.Quantity
	to.FeedBackThumbsUp = from.FeedBackThumbsUp
	to.FeedBackThumbsDown = from.FeedBackThumbsDown
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}

//func (product *ProductModel) updateProductByID() (int, error) {
//	if err := db.VerifyNOSQLDatabaseConnection(db.Client); err != nil {
//		return http.StatusInternalServerError, fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", ProductTableName, err)
//	}
//	product.UpdatedAt = time.Now()
//	bsonCourt, _ := toDoc(product)
//	update := bson.D{{"$set", bsonCourt}}
//	if _, err := db.Client.UpdateByID(common.Ctx, product.ID, update); err != nil {
//		return http.StatusInternalServerError, fmt.Errorf("exception while performing %s Operation on Table: %s. %v", "Update", ProductTableName, err)
//	}
//	return http.StatusOK, nil
//}
//
//
//func toDoc(v interface{}) (doc *bson.D, err error) {
//	data, err := bson.Marshal(v)
//	if err != nil {
//		return
//	}
//
//	err = bson.Unmarshal(data, &doc)
//	return
//}
