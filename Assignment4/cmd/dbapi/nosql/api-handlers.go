package main

import (
	"context"
	"fmt"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/adarshsrinivasan/DS_S24/library/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type noSQLServer struct {
	proto.UnimplementedNOSQLServiceServer
}

func (server *noSQLServer) Initialize(ctx context.Context, request *proto.InitializeRequest) (*proto.InitializeResponse, error) {
	if err := initialize(ctx, request.ServiceName, request.SQLSchemaName); err != nil {
		err = fmt.Errorf("exception while initializing.... %v", err)
		log.Panicf("Initialize: %v\n", err)
	}
	response := &proto.InitializeResponse{Err: common.ConvertErrorToProtoError(err)}
	return response, err
}

func (server *noSQLServer) CreateProduct(ctx context.Context, request *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateProduct(ctx)
	response := &proto.CreateProductResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertProductTableModelToProtoProductModel(ctx, tableModel),
	}
	return response, err
}
func (server *noSQLServer) GetProductByID(ctx context.Context, request *proto.GetProductByIDRequest) (*proto.GetProductByIDResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetProductByID(ctx)
	response := &proto.GetProductByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertProductTableModelToProtoProductModel(ctx, tableModel),
	}
	return response, err
}
func (server *noSQLServer) ListProductsByKeyWordsAndCategory(ctx context.Context, request *proto.ListProductsByKeyWordsAndCategoryRequest) (*proto.ListProductsByKeyWordsAndCategoryResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListProductsByKeyWordsAndCategory(ctx)
	var listProtoResponse []*proto.ProductModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertProductTableModelToProtoProductModel(ctx, &resp))
		}
	}
	response := &proto.ListProductsByKeyWordsAndCategoryResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *noSQLServer) ListProductsBySellerID(ctx context.Context, request *proto.ListProductsBySellerIDRequest) (*proto.ListProductsBySellerIDResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListProductsBySellerID(ctx)
	var listProtoResponse []*proto.ProductModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertProductTableModelToProtoProductModel(ctx, &resp))
		}
	}
	response := &proto.ListProductsBySellerIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *noSQLServer) UpdateProductByID(ctx context.Context, request *proto.UpdateProductByIDRequest) (*proto.UpdateProductByIDResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.UpdateProductByID(ctx)
	response := &proto.UpdateProductByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertProductTableModelToProtoProductModel(ctx, tableModel),
	}
	return response, err
}
func (server *noSQLServer) DeleteProductByID(ctx context.Context, request *proto.DeleteProductByIDRequest) (*proto.DeleteProductByIDResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteProductByID(ctx)
	response := &proto.DeleteProductByIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}

func convertProtoProductModelToProductTableModel(ctx context.Context, protoProductModel *proto.ProductModel) *ProductTableModel {
	return &ProductTableModel{
		ID:                 protoProductModel.ID,
		Name:               protoProductModel.Name,
		Category:           CATEGORY(protoProductModel.Category),
		Keywords:           protoProductModel.Keywords,
		Condition:          CONDITION(protoProductModel.Condition),
		SalePrice:          protoProductModel.SalePrice,
		SellerID:           protoProductModel.SellerID,
		Quantity:           int(protoProductModel.Quantity),
		FeedBackThumbsUp:   int(protoProductModel.FeedBackThumbsUp),
		FeedBackThumbsDown: int(protoProductModel.FeedBackThumbsDown),
		CreatedAt:          protoProductModel.CreatedAt.AsTime(),
		UpdatedAt:          protoProductModel.UpdatedAt.AsTime(),
	}
}

func convertProductTableModelToProtoProductModel(ctx context.Context, productTableModel *ProductTableModel) *proto.ProductModel {
	return &proto.ProductModel{
		ID:                 productTableModel.ID,
		Name:               productTableModel.Name,
		Category:           proto.CATEGORY(productTableModel.Category),
		Keywords:           productTableModel.Keywords,
		Condition:          proto.CONDITION(productTableModel.Condition),
		SalePrice:          productTableModel.SalePrice,
		SellerID:           productTableModel.SellerID,
		Quantity:           int32(productTableModel.Quantity),
		FeedBackThumbsUp:   int32(productTableModel.FeedBackThumbsUp),
		FeedBackThumbsDown: int32(productTableModel.FeedBackThumbsDown),
		CreatedAt:          timestamppb.New(productTableModel.CreatedAt),
		UpdatedAt:          timestamppb.New(productTableModel.UpdatedAt),
	}
}
