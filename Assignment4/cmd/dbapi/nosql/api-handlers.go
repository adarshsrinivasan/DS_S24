package main

import (
	"context"
	"time"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	libProto "github.com/adarshsrinivasan/DS_S24/library/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type noSQLServer struct {
	libProto.UnimplementedNOSQLServiceServer
}

func (server *noSQLServer) GetLeader(ctx context.Context, request *libProto.GetLeaderRequest) (*libProto.GetLeaderResponse, error) {
	for raftServer.cm.leaderID == "" {
		log.Warnf("GetLeader(%s): LeaderID Empty. Waiting for leader to be assigned.", nodeName)
		time.Sleep(100 * time.Millisecond)
	}
	return &libProto.GetLeaderResponse{
		LeaderNodeName: raftServer.cm.leaderID,
		Err:            nil,
	}, nil
}

func (server *noSQLServer) CreateProduct(ctx context.Context, request *libProto.CreateProductRequest) (*libProto.CreateProductResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := CreateProduct
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := noSQLServerHandlers{}
	return handler.CreateProduct(ctx, request)
}

func (server *noSQLServer) GetProductByID(ctx context.Context, request *libProto.GetProductByIDRequest) (*libProto.GetProductByIDResponse, error) {
	handler := noSQLServerHandlers{}
	return handler.GetProductByID(ctx, request)
}
func (server *noSQLServer) ListProductsByKeyWordsAndCategory(ctx context.Context, request *libProto.ListProductsByKeyWordsAndCategoryRequest) (*libProto.ListProductsByKeyWordsAndCategoryResponse, error) {
	handler := noSQLServerHandlers{}
	return handler.ListProductsByKeyWordsAndCategory(ctx, request)
}
func (server *noSQLServer) ListProductsBySellerID(ctx context.Context, request *libProto.ListProductsBySellerIDRequest) (*libProto.ListProductsBySellerIDResponse, error) {
	handler := noSQLServerHandlers{}
	return handler.ListProductsBySellerID(ctx, request)
}
func (server *noSQLServer) UpdateProductByID(ctx context.Context, request *libProto.UpdateProductByIDRequest) (*libProto.UpdateProductByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := UpdateProductByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := noSQLServerHandlers{}
	return handler.UpdateProductByID(ctx, request)
}
func (server *noSQLServer) DeleteProductByID(ctx context.Context, request *libProto.DeleteProductByIDRequest) (*libProto.DeleteProductByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := DeleteProductByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := noSQLServerHandlers{}
	return handler.DeleteProductByID(ctx, request)
}

type noSQLServerHandlers struct {
}

func (server *noSQLServerHandlers) CreateProduct(ctx context.Context, request *libProto.CreateProductRequest) (*libProto.CreateProductResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateProduct(ctx)
	response := &libProto.CreateProductResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertProductTableModelToProtoProductModel(ctx, tableModel),
	}
	return response, err
}
func (server *noSQLServerHandlers) GetProductByID(ctx context.Context, request *libProto.GetProductByIDRequest) (*libProto.GetProductByIDResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetProductByID(ctx)
	response := &libProto.GetProductByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertProductTableModelToProtoProductModel(ctx, tableModel),
	}
	return response, err
}
func (server *noSQLServerHandlers) ListProductsByKeyWordsAndCategory(ctx context.Context, request *libProto.ListProductsByKeyWordsAndCategoryRequest) (*libProto.ListProductsByKeyWordsAndCategoryResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListProductsByKeyWordsAndCategory(ctx)
	var listProtoResponse []*libProto.ProductModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertProductTableModelToProtoProductModel(ctx, &resp))
		}
	}
	response := &libProto.ListProductsByKeyWordsAndCategoryResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *noSQLServerHandlers) ListProductsBySellerID(ctx context.Context, request *libProto.ListProductsBySellerIDRequest) (*libProto.ListProductsBySellerIDResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListProductsBySellerID(ctx)
	var listProtoResponse []*libProto.ProductModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertProductTableModelToProtoProductModel(ctx, &resp))
		}
	}
	response := &libProto.ListProductsBySellerIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *noSQLServerHandlers) UpdateProductByID(ctx context.Context, request *libProto.UpdateProductByIDRequest) (*libProto.UpdateProductByIDResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.UpdateProductByID(ctx)
	response := &libProto.UpdateProductByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertProductTableModelToProtoProductModel(ctx, tableModel),
	}
	return response, err
}
func (server *noSQLServerHandlers) DeleteProductByID(ctx context.Context, request *libProto.DeleteProductByIDRequest) (*libProto.DeleteProductByIDResponse, error) {
	tableModel := convertProtoProductModelToProductTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteProductByID(ctx)
	response := &libProto.DeleteProductByIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}

func convertProtoProductModelToProductTableModel(ctx context.Context, protoProductModel *libProto.ProductModel) *ProductTableModel {
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

func convertProductTableModelToProtoProductModel(ctx context.Context, productTableModel *ProductTableModel) *libProto.ProductModel {
	return &libProto.ProductModel{
		ID:                 productTableModel.ID,
		Name:               productTableModel.Name,
		Category:           libProto.CATEGORY(productTableModel.Category),
		Keywords:           productTableModel.Keywords,
		Condition:          libProto.CONDITION(productTableModel.Condition),
		SalePrice:          productTableModel.SalePrice,
		SellerID:           productTableModel.SellerID,
		Quantity:           int32(productTableModel.Quantity),
		FeedBackThumbsUp:   int32(productTableModel.FeedBackThumbsUp),
		FeedBackThumbsDown: int32(productTableModel.FeedBackThumbsDown),
		CreatedAt:          timestamppb.New(productTableModel.CreatedAt),
		UpdatedAt:          timestamppb.New(productTableModel.UpdatedAt),
	}
}
