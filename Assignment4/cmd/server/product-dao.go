package main

import (
	"context"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/adarshsrinivasan/DS_S24/library/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
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

type ProductTableOps interface {
	CreateProduct(ctx context.Context) (int, error)
	GetProductByID(ctx context.Context) (int, error)
	ListProductsByKeyWordsAndCategory(ctx context.Context) ([]ProductModel, int, error)
	ListProductsBySellerID(ctx context.Context) ([]ProductModel, int, error)
	UpdateProductByID(ctx context.Context) (int, error)
	DeleteProductByID(ctx context.Context) (int, error)
}

func (product *ProductModel) CreateProduct(ctx context.Context) (int, error) {
	protoModel := convertProductModelToProtoProductModel(ctx, product)
	request := &proto.CreateProductRequest{
		RequestModel: protoModel,
	}
	nosqlDBClient, conn, err := common.NewNOSQLRPCClient(ctx, nosqlRPCHost, nosqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("CreateProduct: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := nosqlDBClient.CreateProduct(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Create", ProductTableName, err)
		logrus.Errorf("CreateProduct: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyProductModelObject(response.ResponseModel, product)
	logrus.Infof("CreateProduct: Successfully created product for ID %s\n", product.ID)
	return http.StatusOK, nil
}
func (product *ProductModel) GetProductByID(ctx context.Context) (int, error) {
	protoModel := convertProductModelToProtoProductModel(ctx, product)
	request := &proto.GetProductByIDRequest{
		RequestModel: protoModel,
	}
	nosqlDBClient, conn, err := common.NewNOSQLRPCClient(ctx, nosqlRPCHost, nosqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	leaderInfo, err := nosqlDBClient.GetLeader(ctx, &proto.GetLeaderRequest{})
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "GetLeader", ProductTableName, err)
		logrus.Errorf("GetProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	conn.Close()

	nosqlDBClient, conn, err = common.NewNOSQLRPCClient(ctx, leaderInfo.GetLeaderNodeName(), int(leaderInfo.GetLeaderNodePort()))
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("GetProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := nosqlDBClient.GetProductByID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", ProductTableName, err)
		logrus.Errorf("GetProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyProductModelObject(response.ResponseModel, product)
	return http.StatusOK, nil
}

func (product *ProductModel) ListProductsByKeyWordsAndCategory(ctx context.Context) ([]ProductModel, int, error) {
	protoModel := convertProductModelToProtoProductModel(ctx, product)
	request := &proto.ListProductsByKeyWordsAndCategoryRequest{
		RequestModel: protoModel,
	}
	nosqlDBClient, conn, err := common.NewNOSQLRPCClient(ctx, nosqlRPCHost, nosqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("ListProductsByKeyWordsAndCategory: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	leaderInfo, err := nosqlDBClient.GetLeader(ctx, &proto.GetLeaderRequest{})
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "GetLeader", ProductTableName, err)
		logrus.Errorf("ListProductsByKeyWordsAndCategory: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	conn.Close()
	nosqlDBClient, conn, err = common.NewNOSQLRPCClient(ctx, leaderInfo.GetLeaderNodeName(), int(leaderInfo.GetLeaderNodePort()))
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("ListProductsByKeyWordsAndCategory: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := nosqlDBClient.ListProductsByKeyWordsAndCategory(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", ProductTableName, err)
		logrus.Errorf("ListProductsByKeyWordsAndCategory: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	var result []ProductModel
	for _, resp := range response.ResponseModel {
		result = append(result, *convertProtoProductModelToProductModel(ctx, resp))
	}
	return result, http.StatusOK, nil
}

func (product *ProductModel) ListProductsBySellerID(ctx context.Context) ([]ProductModel, int, error) {
	protoModel := convertProductModelToProtoProductModel(ctx, product)
	request := &proto.ListProductsBySellerIDRequest{
		RequestModel: protoModel,
	}
	nosqlDBClient, conn, err := common.NewNOSQLRPCClient(ctx, nosqlRPCHost, nosqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("ListProductsBySellerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	leaderInfo, err := nosqlDBClient.GetLeader(ctx, &proto.GetLeaderRequest{})
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "GetLeader", ProductTableName, err)
		logrus.Errorf("ListProductsBySellerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	conn.Close()
	nosqlDBClient, conn, err = common.NewNOSQLRPCClient(ctx, leaderInfo.GetLeaderNodeName(), int(leaderInfo.GetLeaderNodePort()))
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("ListProductsBySellerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := nosqlDBClient.ListProductsBySellerID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", ProductTableName, err)
		logrus.Errorf("ListProductsBySellerID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	var result []ProductModel
	for _, resp := range response.ResponseModel {
		result = append(result, *convertProtoProductModelToProductModel(ctx, resp))
	}
	return result, http.StatusOK, nil
}

func (product *ProductModel) UpdateProductByID(ctx context.Context) (int, error) {
	protoModel := convertProductModelToProtoProductModel(ctx, product)
	request := &proto.UpdateProductByIDRequest{
		RequestModel: protoModel,
	}
	nosqlDBClient, conn, err := common.NewNOSQLRPCClient(ctx, nosqlRPCHost, nosqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("UpdateProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	leaderInfo, err := nosqlDBClient.GetLeader(ctx, &proto.GetLeaderRequest{})
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "GetLeader", ProductTableName, err)
		logrus.Errorf("UpdateProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	conn.Close()
	nosqlDBClient, conn, err = common.NewNOSQLRPCClient(ctx, leaderInfo.GetLeaderNodeName(), int(leaderInfo.GetLeaderNodePort()))
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("UpdateProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	response, err := nosqlDBClient.UpdateProductByID(ctx, request)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Update", ProductTableName, err)
		logrus.Errorf("UpdateProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	copyProductModelObject(response.ResponseModel, product)
	return http.StatusOK, nil
}
func (product *ProductModel) DeleteProductByID(ctx context.Context) (int, error) {
	protoModel := convertProductModelToProtoProductModel(ctx, product)
	request := &proto.DeleteProductByIDRequest{
		RequestModel: protoModel,
	}
	nosqlDBClient, conn, err := common.NewNOSQLRPCClient(ctx, nosqlRPCHost, nosqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("DeleteProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	leaderInfo, err := nosqlDBClient.GetLeader(ctx, &proto.GetLeaderRequest{})
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "GetLeader", ProductTableName, err)
		logrus.Errorf("DeleteProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	conn.Close()
	nosqlDBClient, conn, err = common.NewNOSQLRPCClient(ctx, leaderInfo.GetLeaderNodeName(), int(leaderInfo.GetLeaderNodePort()))
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("DeleteProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer conn.Close()

	if _, err := nosqlDBClient.DeleteProductByID(ctx, request); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", ProductTableName, err)
		logrus.Errorf("DeleteProductByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func copyProductModelObject(from *proto.ProductModel, to *ProductModel) {
	to.ID = from.ID
	to.Name = from.Name
	to.Category = CATEGORY(from.Category)
	to.Keywords = from.Keywords
	to.Condition = CONDITION(from.Condition)
	to.SalePrice = from.SalePrice
	to.SellerID = from.SellerID
	to.Quantity = int(from.Quantity)
	to.FeedBackThumbsUp = int(from.FeedBackThumbsUp)
	to.FeedBackThumbsDown = int(from.FeedBackThumbsDown)
	to.CreatedAt = from.CreatedAt.AsTime()
	to.UpdatedAt = from.UpdatedAt.AsTime()
}

func convertProductModelToProtoProductModel(ctx context.Context, model *ProductModel) *proto.ProductModel {
	return &proto.ProductModel{
		ID:                 model.ID,
		Name:               model.Name,
		Category:           proto.CATEGORY(model.Category),
		Keywords:           model.Keywords,
		Condition:          proto.CONDITION(model.Condition),
		SalePrice:          model.SalePrice,
		SellerID:           model.SellerID,
		Quantity:           int32(model.Quantity),
		FeedBackThumbsUp:   int32(model.FeedBackThumbsUp),
		FeedBackThumbsDown: int32(model.FeedBackThumbsDown),
		CreatedAt:          timestamppb.New(model.CreatedAt),
		UpdatedAt:          timestamppb.New(model.CreatedAt),
	}
}

func convertProtoProductModelToProductModel(ctx context.Context, protoModel *proto.ProductModel) *ProductModel {
	return &ProductModel{
		ID:                 protoModel.ID,
		Name:               protoModel.Name,
		Category:           CATEGORY(protoModel.Category),
		Keywords:           protoModel.Keywords,
		Condition:          CONDITION(protoModel.Condition),
		SalePrice:          protoModel.SalePrice,
		SellerID:           protoModel.SellerID,
		Quantity:           int(protoModel.Quantity),
		FeedBackThumbsUp:   int(protoModel.FeedBackThumbsUp),
		FeedBackThumbsDown: int(protoModel.FeedBackThumbsDown),
		CreatedAt:          protoModel.CreatedAt.AsTime(),
		UpdatedAt:          protoModel.UpdatedAt.AsTime(),
	}
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
