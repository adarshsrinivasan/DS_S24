package main

import (
	"context"
	"github.com/golang/protobuf/proto"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	libProto "github.com/adarshsrinivasan/DS_S24/library/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sqlServer struct {
	libProto.UnimplementedSQLServiceServer
}

func (server *sqlServer) CreateBuyer(ctx context.Context, request *libProto.CreateBuyerRequest) (*libProto.CreateBuyerResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := CreateBuyer
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.CreateBuyer(ctx, request)
}
func (server *sqlServer) GetBuyerByID(ctx context.Context, request *libProto.GetBuyerByIDRequest) (*libProto.GetBuyerByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := GetBuyerByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.GetBuyerByID(ctx, request)
}
func (server *sqlServer) GetBuyerByUserName(ctx context.Context, request *libProto.GetBuyerByUserNameRequest) (*libProto.GetBuyerByUserNameResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := GetBuyerByUserName
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.GetBuyerByUserName(ctx, request)
}
func (server *sqlServer) UpdateBuyerByID(ctx context.Context, request *libProto.UpdateBuyerByIDRequest) (*libProto.UpdateBuyerByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := UpdateBuyerByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.UpdateBuyerByID(ctx, request)
}
func (server *sqlServer) CreateCart(ctx context.Context, request *libProto.CreateCartRequest) (*libProto.CreateCartResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := CreateCart
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.CreateCart(ctx, request)
}
func (server *sqlServer) GetCartByID(ctx context.Context, request *libProto.GetCartByIDRequest) (*libProto.GetCartByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := GetCartByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.GetCartByID(ctx, request)
}
func (server *sqlServer) GetCartByBuyerID(ctx context.Context, request *libProto.GetCartByBuyerIDRequest) (*libProto.GetCartByBuyerIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := GetCartByBuyerID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.GetCartByBuyerID(ctx, request)
}
func (server *sqlServer) UpdateCartByID(ctx context.Context, request *libProto.UpdateCartByIDRequest) (*libProto.UpdateCartByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := UpdateCartByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.UpdateCartByID(ctx, request)
}
func (server *sqlServer) DeleteCartByID(ctx context.Context, request *libProto.DeleteCartByIDRequest) (*libProto.DeleteCartByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := DeleteCartByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.DeleteCartByID(ctx, request)
}
func (server *sqlServer) CreateCartItem(ctx context.Context, request *libProto.CreateCartItemRequest) (*libProto.CreateCartItemResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := CreateCartItem
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.CreateCartItem(ctx, request)
}
func (server *sqlServer) GetCartItemByID(ctx context.Context, request *libProto.GetCartItemByIDRequest) (*libProto.GetCartItemByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := GetCartItemByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.GetCartItemByID(ctx, request)
}
func (server *sqlServer) GetCartItemByCartIDAndProductID(ctx context.Context, request *libProto.GetCartItemByCartIDAndProductIDRequest) (*libProto.GetCartItemByCartIDAndProductIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := GetCartItemByCartIDAndProductID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.GetCartItemByCartIDAndProductID(ctx, request)
}
func (server *sqlServer) ListCartItemByCartID(ctx context.Context, request *libProto.ListCartItemByCartIDRequest) (*libProto.ListCartItemByCartIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := ListCartItemByCartID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.ListCartItemByCartID(ctx, request)
}
func (server *sqlServer) UpdateCartItem(ctx context.Context, request *libProto.UpdateCartItemRequest) (*libProto.UpdateCartItemResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := UpdateCartItem
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.UpdateCartItem(ctx, request)
}
func (server *sqlServer) DeleteCartItemByCartIDAndProductID(ctx context.Context, request *libProto.DeleteCartItemByCartIDAndProductIDRequest) (*libProto.DeleteCartItemByCartIDAndProductIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := DeleteCartItemByCartIDAndProductID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.DeleteCartItemByCartIDAndProductID(ctx, request)
}
func (server *sqlServer) DeleteCartItemByCartID(ctx context.Context, request *libProto.DeleteCartItemByCartIDRequest) (*libProto.DeleteCartItemByCartIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := DeleteCartItemByCartID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.DeleteCartItemByCartID(ctx, request)
}
func (server *sqlServer) DeleteCartItemByProductID(ctx context.Context, request *libProto.DeleteCartItemByProductIDRequest) (*libProto.DeleteCartItemByProductIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := DeleteCartItemByProductID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.DeleteCartItemByProductID(ctx, request)
}
func (server *sqlServer) CreateSeller(ctx context.Context, request *libProto.CreateSellerRequest) (*libProto.CreateSellerResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := CreateSeller
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.CreateSeller(ctx, request)
}
func (server *sqlServer) GetSellerByID(ctx context.Context, request *libProto.GetSellerByIDRequest) (*libProto.GetSellerByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := GetSellerByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.GetSellerByID(ctx, request)
}
func (server *sqlServer) GetSellerByUserName(ctx context.Context, request *libProto.GetSellerByUserNameRequest) (*libProto.GetSellerByUserNameResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := GetSellerByUserName
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.GetSellerByUserName(ctx, request)
}
func (server *sqlServer) UpdateSellerByID(ctx context.Context, request *libProto.UpdateSellerByIDRequest) (*libProto.UpdateSellerByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := UpdateSellerByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.UpdateSellerByID(ctx, request)
}
func (server *sqlServer) CreateSession(ctx context.Context, request *libProto.CreateSessionRequest) (*libProto.CreateSessionResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := CreateSession
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.CreateSession(ctx, request)
}
func (server *sqlServer) GetSessionByID(ctx context.Context, request *libProto.GetSessionByIDRequest) (*libProto.GetSessionByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := GetSessionByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.GetSessionByID(ctx, request)
}
func (server *sqlServer) GetSessionByUserID(ctx context.Context, request *libProto.GetSessionByUserIDRequest) (*libProto.GetSessionByUserIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := GetSessionByUserID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.GetSessionByUserID(ctx, request)
}
func (server *sqlServer) DeleteSessionByID(ctx context.Context, request *libProto.DeleteSessionByIDRequest) (*libProto.DeleteSessionByIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := DeleteSessionByID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.DeleteSessionByID(ctx, request)
}
func (server *sqlServer) CreateTransaction(ctx context.Context, request *libProto.CreateTransactionRequest) (*libProto.CreateTransactionResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := CreateTransaction
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.CreateTransaction(ctx, request)
}
func (server *sqlServer) ListTransactionsBySellerID(ctx context.Context, request *libProto.ListTransactionsBySellerIDRequest) (*libProto.ListTransactionsBySellerIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := ListTransactionsBySellerID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.ListTransactionsBySellerID(ctx, request)
}
func (server *sqlServer) ListTransactionsByBuyerID(ctx context.Context, request *libProto.ListTransactionsByBuyerIDRequest) (*libProto.ListTransactionsByBuyerIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := ListTransactionsByBuyerID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.ListTransactionsByBuyerID(ctx, request)
}
func (server *sqlServer) ListTransactionsByCartID(ctx context.Context, request *libProto.ListTransactionsByCartIDRequest) (*libProto.ListTransactionsByCartIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := ListTransactionsByCartID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.ListTransactionsByCartID(ctx, request)
}
func (server *sqlServer) DeleteTransactionsByCartID(ctx context.Context, request *libProto.DeleteTransactionsByCartIDRequest) (*libProto.DeleteTransactionsByCartIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := DeleteTransactionsByCartID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.DeleteTransactionsByCartID(ctx, request)
}
func (server *sqlServer) DeleteTransactionsByBuyerID(ctx context.Context, request *libProto.DeleteTransactionsByBuyerIDRequest) (*libProto.DeleteTransactionsByBuyerIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := DeleteTransactionsByBuyerID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.DeleteTransactionsByBuyerID(ctx, request)
}
func (server *sqlServer) DeleteTransactionsBySellerID(ctx context.Context, request *libProto.DeleteTransactionsBySellerIDRequest) (*libProto.DeleteTransactionsBySellerIDResponse, error) {
	payload, _ := proto.Marshal(request)
	opsType := DeleteTransactionsBySellerID
	requestID, respChan := sendRequestToPeers(ctx, opsType, payload)
	log.Infof("%s: Waiting. for requestID: %s to complete.\n", opsTypeToStr[opsType], requestID)
	<-respChan
	log.Infof("%s: requestID: %s completed!\n", opsTypeToStr[opsType], requestID)
	handler := sqlServerHandlers{}
	return handler.DeleteTransactionsBySellerID(ctx, request)
}

type sqlServerHandlers struct {
}

func (server *sqlServerHandlers) CreateBuyer(ctx context.Context, request *libProto.CreateBuyerRequest) (*libProto.CreateBuyerResponse, error) {
	tableModel := convertProtoBuyerModelToBuyerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateBuyer(ctx)
	response := &libProto.CreateBuyerResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertBuyerTableModelToProtoBuyerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) GetBuyerByID(ctx context.Context, request *libProto.GetBuyerByIDRequest) (*libProto.GetBuyerByIDResponse, error) {
	tableModel := convertProtoBuyerModelToBuyerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetBuyerByID(ctx)
	response := &libProto.GetBuyerByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertBuyerTableModelToProtoBuyerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) GetBuyerByUserName(ctx context.Context, request *libProto.GetBuyerByUserNameRequest) (*libProto.GetBuyerByUserNameResponse, error) {
	tableModel := convertProtoBuyerModelToBuyerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetBuyerByUserName(ctx)
	response := &libProto.GetBuyerByUserNameResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertBuyerTableModelToProtoBuyerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) UpdateBuyerByID(ctx context.Context, request *libProto.UpdateBuyerByIDRequest) (*libProto.UpdateBuyerByIDResponse, error) {
	tableModel := convertProtoBuyerModelToBuyerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.UpdateBuyerByID(ctx)
	response := &libProto.UpdateBuyerByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertBuyerTableModelToProtoBuyerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) CreateCart(ctx context.Context, request *libProto.CreateCartRequest) (*libProto.CreateCartResponse, error) {
	tableModel := convertProtoCartModelToCartTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateCart(ctx)
	response := &libProto.CreateCartResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartTableModelToProtoCartModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) GetCartByID(ctx context.Context, request *libProto.GetCartByIDRequest) (*libProto.GetCartByIDResponse, error) {
	tableModel := convertProtoCartModelToCartTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetCartByID(ctx)
	response := &libProto.GetCartByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartTableModelToProtoCartModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) GetCartByBuyerID(ctx context.Context, request *libProto.GetCartByBuyerIDRequest) (*libProto.GetCartByBuyerIDResponse, error) {
	tableModel := convertProtoCartModelToCartTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetCartByBuyerID(ctx)
	response := &libProto.GetCartByBuyerIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartTableModelToProtoCartModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) UpdateCartByID(ctx context.Context, request *libProto.UpdateCartByIDRequest) (*libProto.UpdateCartByIDResponse, error) {
	tableModel := convertProtoCartModelToCartTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.UpdateCartByID(ctx)
	response := &libProto.UpdateCartByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartTableModelToProtoCartModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) DeleteCartByID(ctx context.Context, request *libProto.DeleteCartByIDRequest) (*libProto.DeleteCartByIDResponse, error) {
	tableModel := convertProtoCartModelToCartTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteCartByID(ctx)
	response := &libProto.DeleteCartByIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServerHandlers) CreateCartItem(ctx context.Context, request *libProto.CreateCartItemRequest) (*libProto.CreateCartItemResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateCartItem(ctx)
	response := &libProto.CreateCartItemResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartItemTableModelToProtoCartItemModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) GetCartItemByID(ctx context.Context, request *libProto.GetCartItemByIDRequest) (*libProto.GetCartItemByIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetCartItemByID(ctx)
	response := &libProto.GetCartItemByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartItemTableModelToProtoCartItemModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) GetCartItemByCartIDAndProductID(ctx context.Context, request *libProto.GetCartItemByCartIDAndProductIDRequest) (*libProto.GetCartItemByCartIDAndProductIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetCartItemByCartIDAndProductID(ctx)
	response := &libProto.GetCartItemByCartIDAndProductIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartItemTableModelToProtoCartItemModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) ListCartItemByCartID(ctx context.Context, request *libProto.ListCartItemByCartIDRequest) (*libProto.ListCartItemByCartIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListCartItemByCartID(ctx)
	var listProtoResponse []*libProto.CartItemModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertCartItemTableModelToProtoCartItemModel(ctx, &resp))
		}
	}
	response := &libProto.ListCartItemByCartIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *sqlServerHandlers) UpdateCartItem(ctx context.Context, request *libProto.UpdateCartItemRequest) (*libProto.UpdateCartItemResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.UpdateCartItem(ctx)
	response := &libProto.UpdateCartItemResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartItemTableModelToProtoCartItemModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) DeleteCartItemByCartIDAndProductID(ctx context.Context, request *libProto.DeleteCartItemByCartIDAndProductIDRequest) (*libProto.DeleteCartItemByCartIDAndProductIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteCartItemByCartIDAndProductID(ctx)
	response := &libProto.DeleteCartItemByCartIDAndProductIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServerHandlers) DeleteCartItemByCartID(ctx context.Context, request *libProto.DeleteCartItemByCartIDRequest) (*libProto.DeleteCartItemByCartIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteCartItemByCartID(ctx)
	response := &libProto.DeleteCartItemByCartIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServerHandlers) DeleteCartItemByProductID(ctx context.Context, request *libProto.DeleteCartItemByProductIDRequest) (*libProto.DeleteCartItemByProductIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteCartItemByProductID(ctx)
	response := &libProto.DeleteCartItemByProductIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServerHandlers) CreateSeller(ctx context.Context, request *libProto.CreateSellerRequest) (*libProto.CreateSellerResponse, error) {
	tableModel := convertProtoSellerModelToSellerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateSeller(ctx)
	response := &libProto.CreateSellerResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSellerTableModelToProtoSellerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) GetSellerByID(ctx context.Context, request *libProto.GetSellerByIDRequest) (*libProto.GetSellerByIDResponse, error) {
	tableModel := convertProtoSellerModelToSellerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetSellerByID(ctx)
	response := &libProto.GetSellerByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSellerTableModelToProtoSellerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) GetSellerByUserName(ctx context.Context, request *libProto.GetSellerByUserNameRequest) (*libProto.GetSellerByUserNameResponse, error) {
	tableModel := convertProtoSellerModelToSellerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetSellerByUserName(ctx)
	response := &libProto.GetSellerByUserNameResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSellerTableModelToProtoSellerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) UpdateSellerByID(ctx context.Context, request *libProto.UpdateSellerByIDRequest) (*libProto.UpdateSellerByIDResponse, error) {
	tableModel := convertProtoSellerModelToSellerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.UpdateSellerByID(ctx)
	response := &libProto.UpdateSellerByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSellerTableModelToProtoSellerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) CreateSession(ctx context.Context, request *libProto.CreateSessionRequest) (*libProto.CreateSessionResponse, error) {
	tableModel := convertProtoSessionModelToSessionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateSession(ctx)
	response := &libProto.CreateSessionResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSessionTableModelToProtoSessionModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) GetSessionByID(ctx context.Context, request *libProto.GetSessionByIDRequest) (*libProto.GetSessionByIDResponse, error) {
	tableModel := convertProtoSessionModelToSessionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetSessionByID(ctx)
	response := &libProto.GetSessionByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSessionTableModelToProtoSessionModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) GetSessionByUserID(ctx context.Context, request *libProto.GetSessionByUserIDRequest) (*libProto.GetSessionByUserIDResponse, error) {
	tableModel := convertProtoSessionModelToSessionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetSessionByUserID(ctx)
	response := &libProto.GetSessionByUserIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSessionTableModelToProtoSessionModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) DeleteSessionByID(ctx context.Context, request *libProto.DeleteSessionByIDRequest) (*libProto.DeleteSessionByIDResponse, error) {
	tableModel := convertProtoSessionModelToSessionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteSessionByID(ctx)
	response := &libProto.DeleteSessionByIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServerHandlers) CreateTransaction(ctx context.Context, request *libProto.CreateTransactionRequest) (*libProto.CreateTransactionResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateTransaction(ctx)
	response := &libProto.CreateTransactionResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertTransactionTableModelToProtoTransactionModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServerHandlers) ListTransactionsBySellerID(ctx context.Context, request *libProto.ListTransactionsBySellerIDRequest) (*libProto.ListTransactionsBySellerIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListTransactionsBySellerID(ctx)
	var listProtoResponse []*libProto.TransactionModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertTransactionTableModelToProtoTransactionModel(ctx, &resp))
		}
	}
	response := &libProto.ListTransactionsBySellerIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *sqlServerHandlers) ListTransactionsByBuyerID(ctx context.Context, request *libProto.ListTransactionsByBuyerIDRequest) (*libProto.ListTransactionsByBuyerIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListTransactionsByBuyerID(ctx)
	var listProtoResponse []*libProto.TransactionModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertTransactionTableModelToProtoTransactionModel(ctx, &resp))
		}
	}
	response := &libProto.ListTransactionsByBuyerIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *sqlServerHandlers) ListTransactionsByCartID(ctx context.Context, request *libProto.ListTransactionsByCartIDRequest) (*libProto.ListTransactionsByCartIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListTransactionsByCartID(ctx)
	var listProtoResponse []*libProto.TransactionModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertTransactionTableModelToProtoTransactionModel(ctx, &resp))
		}
	}
	response := &libProto.ListTransactionsByCartIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *sqlServerHandlers) DeleteTransactionsByCartID(ctx context.Context, request *libProto.DeleteTransactionsByCartIDRequest) (*libProto.DeleteTransactionsByCartIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteTransactionsByCartID(ctx)
	response := &libProto.DeleteTransactionsByCartIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServerHandlers) DeleteTransactionsByBuyerID(ctx context.Context, request *libProto.DeleteTransactionsByBuyerIDRequest) (*libProto.DeleteTransactionsByBuyerIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteTransactionsByBuyerID(ctx)
	response := &libProto.DeleteTransactionsByBuyerIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServerHandlers) DeleteTransactionsBySellerID(ctx context.Context, request *libProto.DeleteTransactionsBySellerIDRequest) (*libProto.DeleteTransactionsBySellerIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteTransactionsBySellerID(ctx)
	response := &libProto.DeleteTransactionsBySellerIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}

func convertBuyerTableModelToProtoBuyerModel(ctx context.Context, buyerTableModel *BuyerTableModel) *libProto.BuyerModel {
	return &libProto.BuyerModel{
		ID:        buyerTableModel.Id,
		Name:      buyerTableModel.Name,
		UserName:  buyerTableModel.UserName,
		Password:  buyerTableModel.Password,
		Version:   int32(buyerTableModel.Version),
		CreatedAt: timestamppb.New(buyerTableModel.CreatedAt),
		UpdatedAt: timestamppb.New(buyerTableModel.UpdatedAt),
	}
}

func convertProtoBuyerModelToBuyerTableModel(ctx context.Context, protoBuyerModel *libProto.BuyerModel) *BuyerTableModel {
	return &BuyerTableModel{
		Id:        protoBuyerModel.ID,
		Name:      protoBuyerModel.Name,
		UserName:  protoBuyerModel.UserName,
		Password:  protoBuyerModel.Password,
		Version:   int(protoBuyerModel.Version),
		CreatedAt: protoBuyerModel.CreatedAt.AsTime(),
		UpdatedAt: protoBuyerModel.UpdatedAt.AsTime(),
	}
}

func convertCartTableModelToProtoCartModel(ctx context.Context, cartTableModel *CartTableModel) *libProto.CartModel {
	return &libProto.CartModel{
		ID:        cartTableModel.ID,
		BuyerID:   cartTableModel.BuyerID,
		Saved:     cartTableModel.Saved,
		Version:   int32(cartTableModel.Version),
		CreatedAt: timestamppb.New(cartTableModel.CreatedAt),
		UpdatedAt: timestamppb.New(cartTableModel.UpdatedAt),
	}
}

func convertProtoCartModelToCartTableModel(ctx context.Context, protoCartModel *libProto.CartModel) *CartTableModel {
	return &CartTableModel{
		ID:        protoCartModel.ID,
		BuyerID:   protoCartModel.BuyerID,
		Saved:     protoCartModel.Saved,
		Version:   int(protoCartModel.Version),
		CreatedAt: protoCartModel.CreatedAt.AsTime(),
		UpdatedAt: protoCartModel.UpdatedAt.AsTime(),
	}
}

func convertCartItemTableModelToProtoCartItemModel(ctx context.Context, cartItemTableModel *CartItemTableModel) *libProto.CartItemModel {
	return &libProto.CartItemModel{
		ID:        cartItemTableModel.ID,
		CartID:    cartItemTableModel.CartID,
		ProductID: cartItemTableModel.ProductID,
		Quantity:  int32(cartItemTableModel.Quantity),
		Version:   int32(cartItemTableModel.Version),
		CreatedAt: timestamppb.New(cartItemTableModel.CreatedAt),
		UpdatedAt: timestamppb.New(cartItemTableModel.UpdatedAt),
	}
}

func convertProtoCartItemModelToCartItemTableModel(ctx context.Context, protoCartItemModel *libProto.CartItemModel) *CartItemTableModel {
	return &CartItemTableModel{
		ID:        protoCartItemModel.ID,
		CartID:    protoCartItemModel.CartID,
		ProductID: protoCartItemModel.ProductID,
		Quantity:  int(protoCartItemModel.Quantity),
		Version:   int(protoCartItemModel.Version),
		CreatedAt: protoCartItemModel.CreatedAt.AsTime(),
		UpdatedAt: protoCartItemModel.UpdatedAt.AsTime(),
	}
}

func convertSellerTableModelToProtoSellerModel(ctx context.Context, sellerTableModel *SellerTableModel) *libProto.SellerModel {
	return &libProto.SellerModel{
		ID:                 sellerTableModel.Id,
		Name:               sellerTableModel.Name,
		FeedBackThumbsUp:   int32(sellerTableModel.FeedBackThumbsUp),
		FeedBackThumbsDown: int32(sellerTableModel.FeedBackThumbsDown),
		NumberOfItemsSold:  int32(sellerTableModel.NumberOfItemsSold),
		UserName:           sellerTableModel.UserName,
		Password:           sellerTableModel.Password,
		Version:            int32(sellerTableModel.Version),
		CreatedAt:          timestamppb.New(sellerTableModel.CreatedAt),
		UpdatedAt:          timestamppb.New(sellerTableModel.UpdatedAt),
	}
}

func convertProtoSellerModelToSellerTableModel(ctx context.Context, protoSellerModel *libProto.SellerModel) *SellerTableModel {
	return &SellerTableModel{
		Id:                 protoSellerModel.ID,
		Name:               protoSellerModel.Name,
		FeedBackThumbsUp:   int(protoSellerModel.FeedBackThumbsUp),
		FeedBackThumbsDown: int(protoSellerModel.FeedBackThumbsDown),
		NumberOfItemsSold:  int(protoSellerModel.NumberOfItemsSold),
		UserName:           protoSellerModel.UserName,
		Password:           protoSellerModel.Password,
		Version:            int(protoSellerModel.Version),
		CreatedAt:          protoSellerModel.CreatedAt.AsTime(),
		UpdatedAt:          protoSellerModel.UpdatedAt.AsTime(),
	}
}

func convertSessionTableModelToProtoSessionModel(ctx context.Context, sessionTableModel *SessionTableModel) *libProto.SessionModel {
	return &libProto.SessionModel{
		ID:        sessionTableModel.ID,
		UserID:    sessionTableModel.UserID,
		UserType:  libProto.USERTYPE(sessionTableModel.UserType),
		Version:   int32(sessionTableModel.Version),
		CreatedAt: timestamppb.New(sessionTableModel.CreatedAt),
		UpdatedAt: timestamppb.New(sessionTableModel.UpdatedAt),
	}
}

func convertProtoSessionModelToSessionTableModel(ctx context.Context, protoSessionModel *libProto.SessionModel) *SessionTableModel {
	return &SessionTableModel{
		ID:        protoSessionModel.ID,
		UserID:    protoSessionModel.UserID,
		UserType:  common.UserType(protoSessionModel.UserType),
		Version:   int(protoSessionModel.Version),
		CreatedAt: protoSessionModel.CreatedAt.AsTime(),
		UpdatedAt: protoSessionModel.UpdatedAt.AsTime(),
	}
}

func convertTransactionTableModelToProtoTransactionModel(ctx context.Context, transactionTableModel *TransactionTableModel) *libProto.TransactionModel {
	return &libProto.TransactionModel{
		ID:        transactionTableModel.ID,
		CartID:    transactionTableModel.CartID,
		ProductID: transactionTableModel.ProductID,
		BuyerID:   transactionTableModel.BuyerID,
		SellerID:  transactionTableModel.SellerID,
		Quantity:  int32(transactionTableModel.Quantity),
		Price:     transactionTableModel.Price,
		Version:   int32(transactionTableModel.Version),
		CreatedAt: timestamppb.New(transactionTableModel.CreatedAt),
		UpdatedAt: timestamppb.New(transactionTableModel.UpdatedAt),
	}
}

func convertProtoTransactionModelToTransactionTableModel(ctx context.Context, protoTransactionModel *libProto.TransactionModel) *TransactionTableModel {
	return &TransactionTableModel{
		ID:        protoTransactionModel.ID,
		CartID:    protoTransactionModel.CartID,
		ProductID: protoTransactionModel.ProductID,
		BuyerID:   protoTransactionModel.BuyerID,
		SellerID:  protoTransactionModel.SellerID,
		Quantity:  int(protoTransactionModel.Quantity),
		Price:     protoTransactionModel.Price,
		Version:   int(protoTransactionModel.Version),
		CreatedAt: protoTransactionModel.CreatedAt.AsTime(),
		UpdatedAt: protoTransactionModel.UpdatedAt.AsTime(),
	}
}
