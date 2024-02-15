package main

import (
	"context"
	"fmt"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/adarshsrinivasan/DS_S24/library/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sqlServer struct {
	proto.UnimplementedSQLServiceServer
}

func (server *sqlServer) Initialize(ctx context.Context, request *proto.InitializeRequest) (*proto.InitializeResponse, error) {
	if err := initialize(ctx, request.ServiceName, request.SQLSchemaName); err != nil {
		err = fmt.Errorf("exception while initializing.... %v", err)
		log.Panicf("Initialize: %v\n", err)
	}
	response := &proto.InitializeResponse{Err: common.ConvertErrorToProtoError(err)}
	return response, err
}

func (server *sqlServer) CreateBuyer(ctx context.Context, request *proto.CreateBuyerRequest) (*proto.CreateBuyerResponse, error) {
	tableModel := convertProtoBuyerModelToBuyerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateBuyer(ctx)
	response := &proto.CreateBuyerResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertBuyerTableModelToProtoBuyerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) GetBuyerByID(ctx context.Context, request *proto.GetBuyerByIDRequest) (*proto.GetBuyerByIDResponse, error) {
	tableModel := convertProtoBuyerModelToBuyerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetBuyerByID(ctx)
	response := &proto.GetBuyerByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertBuyerTableModelToProtoBuyerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) GetBuyerByUserName(ctx context.Context, request *proto.GetBuyerByUserNameRequest) (*proto.GetBuyerByUserNameResponse, error) {
	tableModel := convertProtoBuyerModelToBuyerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetBuyerByUserName(ctx)
	response := &proto.GetBuyerByUserNameResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertBuyerTableModelToProtoBuyerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) UpdateBuyerByID(ctx context.Context, request *proto.UpdateBuyerByIDRequest) (*proto.UpdateBuyerByIDResponse, error) {
	tableModel := convertProtoBuyerModelToBuyerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.UpdateBuyerByID(ctx)
	response := &proto.UpdateBuyerByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertBuyerTableModelToProtoBuyerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) CreateCart(ctx context.Context, request *proto.CreateCartRequest) (*proto.CreateCartResponse, error) {
	tableModel := convertProtoCartModelToCartTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateCart(ctx)
	response := &proto.CreateCartResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartTableModelToProtoCartModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) GetCartByID(ctx context.Context, request *proto.GetCartByIDRequest) (*proto.GetCartByIDResponse, error) {
	tableModel := convertProtoCartModelToCartTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetCartByID(ctx)
	response := &proto.GetCartByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartTableModelToProtoCartModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) GetCartByBuyerID(ctx context.Context, request *proto.GetCartByBuyerIDRequest) (*proto.GetCartByBuyerIDResponse, error) {
	tableModel := convertProtoCartModelToCartTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetCartByBuyerID(ctx)
	response := &proto.GetCartByBuyerIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartTableModelToProtoCartModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) UpdateCartByID(ctx context.Context, request *proto.UpdateCartByIDRequest) (*proto.UpdateCartByIDResponse, error) {
	tableModel := convertProtoCartModelToCartTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.UpdateCartByID(ctx)
	response := &proto.UpdateCartByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartTableModelToProtoCartModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) DeleteCartByID(ctx context.Context, request *proto.DeleteCartByIDRequest) (*proto.DeleteCartByIDResponse, error) {
	tableModel := convertProtoCartModelToCartTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteCartByID(ctx)
	response := &proto.DeleteCartByIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServer) CreateCartItem(ctx context.Context, request *proto.CreateCartItemRequest) (*proto.CreateCartItemResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateCartItem(ctx)
	response := &proto.CreateCartItemResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartItemTableModelToProtoCartItemModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) GetCartItemByID(ctx context.Context, request *proto.GetCartItemByIDRequest) (*proto.GetCartItemByIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetCartItemByID(ctx)
	response := &proto.GetCartItemByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartItemTableModelToProtoCartItemModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) GetCartItemByCartIDAndProductID(ctx context.Context, request *proto.GetCartItemByCartIDAndProductIDRequest) (*proto.GetCartItemByCartIDAndProductIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetCartItemByCartIDAndProductID(ctx)
	response := &proto.GetCartItemByCartIDAndProductIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartItemTableModelToProtoCartItemModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) ListCartItemByCartID(ctx context.Context, request *proto.ListCartItemByCartIDRequest) (*proto.ListCartItemByCartIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListCartItemByCartID(ctx)
	var listProtoResponse []*proto.CartItemModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertCartItemTableModelToProtoCartItemModel(ctx, &resp))
		}
	}
	response := &proto.ListCartItemByCartIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *sqlServer) UpdateCartItem(ctx context.Context, request *proto.UpdateCartItemRequest) (*proto.UpdateCartItemResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.UpdateCartItem(ctx)
	response := &proto.UpdateCartItemResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertCartItemTableModelToProtoCartItemModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) DeleteCartItemByCartIDAndProductID(ctx context.Context, request *proto.DeleteCartItemByCartIDAndProductIDRequest) (*proto.DeleteCartItemByCartIDAndProductIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteCartItemByCartIDAndProductID(ctx)
	response := &proto.DeleteCartItemByCartIDAndProductIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServer) DeleteCartItemByCartID(ctx context.Context, request *proto.DeleteCartItemByCartIDRequest) (*proto.DeleteCartItemByCartIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteCartItemByCartID(ctx)
	response := &proto.DeleteCartItemByCartIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServer) DeleteCartItemByProductID(ctx context.Context, request *proto.DeleteCartItemByProductIDRequest) (*proto.DeleteCartItemByProductIDResponse, error) {
	tableModel := convertProtoCartItemModelToCartItemTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteCartItemByProductID(ctx)
	response := &proto.DeleteCartItemByProductIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServer) CreateSeller(ctx context.Context, request *proto.CreateSellerRequest) (*proto.CreateSellerResponse, error) {
	tableModel := convertProtoSellerModelToSellerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateSeller(ctx)
	response := &proto.CreateSellerResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSellerTableModelToProtoSellerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) GetSellerByID(ctx context.Context, request *proto.GetSellerByIDRequest) (*proto.GetSellerByIDResponse, error) {
	tableModel := convertProtoSellerModelToSellerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetSellerByID(ctx)
	response := &proto.GetSellerByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSellerTableModelToProtoSellerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) GetSellerByUserName(ctx context.Context, request *proto.GetSellerByUserNameRequest) (*proto.GetSellerByUserNameResponse, error) {
	tableModel := convertProtoSellerModelToSellerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetSellerByUserName(ctx)
	response := &proto.GetSellerByUserNameResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSellerTableModelToProtoSellerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) UpdateSellerByID(ctx context.Context, request *proto.UpdateSellerByIDRequest) (*proto.UpdateSellerByIDResponse, error) {
	tableModel := convertProtoSellerModelToSellerTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.UpdateSellerByID(ctx)
	response := &proto.UpdateSellerByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSellerTableModelToProtoSellerModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) CreateSession(ctx context.Context, request *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	tableModel := convertProtoSessionModelToSessionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateSession(ctx)
	response := &proto.CreateSessionResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSessionTableModelToProtoSessionModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) GetSessionByID(ctx context.Context, request *proto.GetSessionByIDRequest) (*proto.GetSessionByIDResponse, error) {
	tableModel := convertProtoSessionModelToSessionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetSessionByID(ctx)
	response := &proto.GetSessionByIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSessionTableModelToProtoSessionModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) GetSessionByUserID(ctx context.Context, request *proto.GetSessionByUserIDRequest) (*proto.GetSessionByUserIDResponse, error) {
	tableModel := convertProtoSessionModelToSessionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.GetSessionByUserID(ctx)
	response := &proto.GetSessionByUserIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertSessionTableModelToProtoSessionModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) DeleteSessionByID(ctx context.Context, request *proto.DeleteSessionByIDRequest) (*proto.DeleteSessionByIDResponse, error) {
	tableModel := convertProtoSessionModelToSessionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteSessionByID(ctx)
	response := &proto.DeleteSessionByIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServer) CreateTransaction(ctx context.Context, request *proto.CreateTransactionRequest) (*proto.CreateTransactionResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.CreateTransaction(ctx)
	response := &proto.CreateTransactionResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: convertTransactionTableModelToProtoTransactionModel(ctx, tableModel),
	}
	return response, err
}
func (server *sqlServer) ListTransactionsBySellerID(ctx context.Context, request *proto.ListTransactionsBySellerIDRequest) (*proto.ListTransactionsBySellerIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListTransactionsBySellerID(ctx)
	var listProtoResponse []*proto.TransactionModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertTransactionTableModelToProtoTransactionModel(ctx, &resp))
		}
	}
	response := &proto.ListTransactionsBySellerIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *sqlServer) ListTransactionsByBuyerID(ctx context.Context, request *proto.ListTransactionsByBuyerIDRequest) (*proto.ListTransactionsByBuyerIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListTransactionsByBuyerID(ctx)
	var listProtoResponse []*proto.TransactionModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertTransactionTableModelToProtoTransactionModel(ctx, &resp))
		}
	}
	response := &proto.ListTransactionsByBuyerIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *sqlServer) ListTransactionsByCartID(ctx context.Context, request *proto.ListTransactionsByCartIDRequest) (*proto.ListTransactionsByCartIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	listResponse, statusCode, err := tableModel.ListTransactionsByCartID(ctx)
	var listProtoResponse []*proto.TransactionModel
	if err == nil {
		for _, resp := range listResponse {
			listProtoResponse = append(listProtoResponse, convertTransactionTableModelToProtoTransactionModel(ctx, &resp))
		}
	}
	response := &proto.ListTransactionsByCartIDResponse{
		StatusCode:    int32(statusCode),
		Err:           common.ConvertErrorToProtoError(err),
		ResponseModel: listProtoResponse,
	}
	return response, err
}
func (server *sqlServer) DeleteTransactionsByCartID(ctx context.Context, request *proto.DeleteTransactionsByCartIDRequest) (*proto.DeleteTransactionsByCartIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteTransactionsByCartID(ctx)
	response := &proto.DeleteTransactionsByCartIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServer) DeleteTransactionsByBuyerID(ctx context.Context, request *proto.DeleteTransactionsByBuyerIDRequest) (*proto.DeleteTransactionsByBuyerIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteTransactionsByBuyerID(ctx)
	response := &proto.DeleteTransactionsByBuyerIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}
func (server *sqlServer) DeleteTransactionsBySellerID(ctx context.Context, request *proto.DeleteTransactionsBySellerIDRequest) (*proto.DeleteTransactionsBySellerIDResponse, error) {
	tableModel := convertProtoTransactionModelToTransactionTableModel(ctx, request.RequestModel)
	statusCode, err := tableModel.DeleteTransactionsBySellerID(ctx)
	response := &proto.DeleteTransactionsBySellerIDResponse{
		StatusCode: int32(statusCode),
		Err:        common.ConvertErrorToProtoError(err),
	}
	return response, err
}

func convertBuyerTableModelToProtoBuyerModel(ctx context.Context, buyerTableModel *BuyerTableModel) *proto.BuyerModel {
	return &proto.BuyerModel{
		ID:                     buyerTableModel.Id,
		Name:                   buyerTableModel.Name,
		UserName:               buyerTableModel.UserName,
		Password:               buyerTableModel.Password,
		Version:                int32(buyerTableModel.Version),
		CreatedAt:              timestamppb.New(buyerTableModel.CreatedAt),
		UpdatedAt:              timestamppb.New(buyerTableModel.UpdatedAt),
	}
}

func convertProtoBuyerModelToBuyerTableModel(ctx context.Context, protoBuyerModel *proto.BuyerModel) *BuyerTableModel {
	return &BuyerTableModel{
		Id:                     protoBuyerModel.ID,
		Name:                   protoBuyerModel.Name,
		UserName:               protoBuyerModel.UserName,
		Password:               protoBuyerModel.Password,
		Version:                int(protoBuyerModel.Version),
		CreatedAt:              protoBuyerModel.CreatedAt.AsTime(),
		UpdatedAt:              protoBuyerModel.UpdatedAt.AsTime(),
	}
}

func convertCartTableModelToProtoCartModel(ctx context.Context, cartTableModel *CartTableModel) *proto.CartModel {
	return &proto.CartModel{
		ID:        cartTableModel.ID,
		BuyerID:   cartTableModel.BuyerID,
		Saved:     cartTableModel.Saved,
		Version:   int32(cartTableModel.Version),
		CreatedAt: timestamppb.New(cartTableModel.CreatedAt),
		UpdatedAt: timestamppb.New(cartTableModel.UpdatedAt),
	}
}

func convertProtoCartModelToCartTableModel(ctx context.Context, protoCartModel *proto.CartModel) *CartTableModel {
	return &CartTableModel{
		ID:        protoCartModel.ID,
		BuyerID:   protoCartModel.BuyerID,
		Saved:     protoCartModel.Saved,
		Version:   int(protoCartModel.Version),
		CreatedAt: protoCartModel.CreatedAt.AsTime(),
		UpdatedAt: protoCartModel.UpdatedAt.AsTime(),
	}
}

func convertCartItemTableModelToProtoCartItemModel(ctx context.Context, cartItemTableModel *CartItemTableModel) *proto.CartItemModel {
	return &proto.CartItemModel{
		ID:        cartItemTableModel.ID,
		CartID:    cartItemTableModel.CartID,
		ProductID: cartItemTableModel.ProductID,
		Quantity:  int32(cartItemTableModel.Quantity),
		Version:   int32(cartItemTableModel.Version),
		CreatedAt: timestamppb.New(cartItemTableModel.CreatedAt),
		UpdatedAt: timestamppb.New(cartItemTableModel.UpdatedAt),
	}
}

func convertProtoCartItemModelToCartItemTableModel(ctx context.Context, protoCartItemModel *proto.CartItemModel) *CartItemTableModel {
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

func convertSellerTableModelToProtoSellerModel(ctx context.Context, sellerTableModel *SellerTableModel) *proto.SellerModel {
	return &proto.SellerModel{
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

func convertProtoSellerModelToSellerTableModel(ctx context.Context, protoSellerModel *proto.SellerModel) *SellerTableModel {
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

func convertSessionTableModelToProtoSessionModel(ctx context.Context, sessionTableModel *SessionTableModel) *proto.SessionModel {
	return &proto.SessionModel{
		ID:        sessionTableModel.ID,
		UserID:    sessionTableModel.UserID,
		UserType:  proto.USERTYPE(sessionTableModel.UserType),
		Version:   int32(sessionTableModel.Version),
		CreatedAt: timestamppb.New(sessionTableModel.CreatedAt),
		UpdatedAt: timestamppb.New(sessionTableModel.UpdatedAt),
	}
}

func convertProtoSessionModelToSessionTableModel(ctx context.Context, protoSessionModel *proto.SessionModel) *SessionTableModel {
	return &SessionTableModel{
		ID:        protoSessionModel.ID,
		UserID:    protoSessionModel.UserID,
		UserType:  common.UserType(protoSessionModel.UserType),
		Version:   int(protoSessionModel.Version),
		CreatedAt: protoSessionModel.CreatedAt.AsTime(),
		UpdatedAt: protoSessionModel.UpdatedAt.AsTime(),
	}
}

func convertTransactionTableModelToProtoTransactionModel(ctx context.Context, transactionTableModel *TransactionTableModel) *proto.TransactionModel {
	return &proto.TransactionModel{
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

func convertProtoTransactionModelToTransactionTableModel(ctx context.Context, protoTransactionModel *proto.TransactionModel) *TransactionTableModel {
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
