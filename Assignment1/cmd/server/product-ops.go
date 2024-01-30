package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"github.com/sirupsen/logrus"
)

type ProductModel struct {
	ID                 string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name               string    `json:"name,omitempty" bson:"name,omitempty"`
	Category           CATEGORY  `json:"category,omitempty" bson:"category,omitempty"`
	Keywords           []string  `json:"keywords,omitempty" bson:"keywords,omitempty"`
	Condition          CONDITION `json:"condition,omitempty" bson:"condition,omitempty"`
	SalePrice          float32   `json:"salePrice,omitempty" bson:"salePrice,omitempty"`
	SellerID           string    `json:"sellerID,omitempty" bson:"sellerID,omitempty"`
	Quantity           int       `json:"quantity,omitempty" bson:"quantity"`
	FeedBackThumbsUp   int       `json:"feedBackThumbsUp,omitempty" bson:"feedBackThumbsUp"`
	FeedBackThumbsDown int       `json:"feedBackThumbsDown,omitempty" bson:"feedBackThumbsDown"`
	CreatedAt          time.Time `json:"createdAt,omitempty"  bson:"createdAt,omitempty"`
	UpdatedAt          time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func createProduct(ctx context.Context, productModel *ProductModel, sessionID string) (ProductModel, int, error) {
	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("createProduct: %v\n", err)
		return ProductModel{}, statusCode, err
	}

	if userType != common.Seller {
		err := fmt.Errorf("provided userID is not a seller %s", userID)
		logrus.Errorf("createProduct: %v\n", err)
		return ProductModel{}, http.StatusForbidden, err
	}
	productModel.SellerID = userID
	if err := validateProductModel(ctx, productModel, true); err != nil {
		err = fmt.Errorf("exception while validating Product data: %v", err)
		logrus.Errorf("createProduct: %v\n", err)
		return ProductModel{}, http.StatusBadRequest, err
	}

	if productModel.SellerID != userID {
		err := fmt.Errorf("current user's ID %s does not match with Product's SellerID %s", userID, productModel.SellerID)
		logrus.Errorf("createProduct: %v\n", err)
		return ProductModel{}, http.StatusForbidden, err
	}

	productTableModel := convertProductModelToProductTableModel(productModel)
	if statusCode, err := productTableModel.CreateProduct(ctx); err != nil {
		err = fmt.Errorf("exception while Creating Product for ID:%s. %v", productModel.ID, err)
		logrus.Errorf("changeItemSalePrice: %v\n", err)
		return ProductModel{}, statusCode, err
	}
	createdProduct := convertProductTableModelToProductModel(productTableModel)
	return *createdProduct, http.StatusOK, nil
}

func searchProduct(ctx context.Context, productModel *ProductModel) ([]ProductModel, int, error) {
	for _, keyword := range productModel.Keywords {
		if len(keyword) > 8 {
			err := fmt.Errorf("invalid Product data. Each Keywords is allowed max 8 characters: %s", keyword)
			logrus.Errorf("searchProduct: %v\n", err)
			return nil, http.StatusBadRequest, err
		}
	}

	productTableModel := convertProductModelToProductTableModel(productModel)
	productTableModels, statusCode, err := productTableModel.GetProductsByKeyWords(ctx)
	if err != nil {
		err = fmt.Errorf("exception while fetching Products data: %v", err)
		logrus.Errorf("searchProduct: %v\n", err)
		return nil, statusCode, err
	}
	productModels := []ProductModel{}
	for _, p := range productTableModels {
		productModels = append(productModels, *convertProductTableModelToProductModel(&p))
	}

	return productModels, http.StatusOK, nil
}

func getProductByID(ctx context.Context, productID string) (ProductModel, int, error) {
	productTableModel := ProductTableModel{ID: productID}
	if statusCode, err := productTableModel.GetProductByID(ctx); err != nil {
		err = fmt.Errorf("exception while fetching Product for ID:%s. %v", productID, err)
		logrus.Errorf("getProductByID: %v\n", err)
		return ProductModel{}, statusCode, err
	}
	fetchedProduct := convertProductTableModelToProductModel(&productTableModel)
	return *fetchedProduct, http.StatusOK, nil
}

func changeItemSalePrice(ctx context.Context, productModel *ProductModel, sessionID string) (ProductModel, int, error) {
	productTableModel := ProductTableModel{ID: productModel.ID}
	if statusCode, err := productTableModel.GetProductByID(ctx); err != nil {
		err = fmt.Errorf("exception while fetching Product for ID:%s. %v", productModel.ID, err)
		logrus.Errorf("changeItemSalePrice: %v\n", err)
		return ProductModel{}, statusCode, err
	}

	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("changeItemSalePrice: %v\n", err)
		return ProductModel{}, statusCode, err
	}

	if userType != common.Seller {
		err := fmt.Errorf("provided userID is not a seller %s", userID)
		logrus.Errorf("changeItemSalePrice: %v\n", err)
		return ProductModel{}, http.StatusForbidden, err
	}

	if productTableModel.SellerID != userID {
		err := fmt.Errorf("current user's ID %s does not match with Product's SellerID %s", userID, productTableModel.SellerID)
		logrus.Errorf("changeItemSalePrice: %v\n", err)
		return ProductModel{}, http.StatusForbidden, err
	}

	productTableModel.SalePrice = productModel.SalePrice
	if statusCode, err := productTableModel.UpdateProductByID(ctx); err != nil {
		err = fmt.Errorf("exception while Updating Product for ID:%s. %v", productModel.ID, err)
		logrus.Errorf("changeItemSalePrice: %v\n", err)
		return ProductModel{}, statusCode, err
	}
	updatedProduct := convertProductTableModelToProductModel(&productTableModel)
	return *updatedProduct, http.StatusOK, nil
}

func removeItemFromSale(ctx context.Context, productModel *ProductModel, sessionID string) (ProductModel, int, error) {
	productTableModel := ProductTableModel{ID: productModel.ID}
	if statusCode, err := productTableModel.GetProductByID(ctx); err != nil {
		err = fmt.Errorf("exception while fetching Product for ID:%s. %v", productModel.ID, err)
		logrus.Errorf("removeItemFromSale: %v\n", err)
		return ProductModel{}, statusCode, err
	}

	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("removeItemFromSale: %v\n", err)
		return ProductModel{}, statusCode, err
	}

	if userType != common.Seller {
		err := fmt.Errorf("provided userID is not a seller %s", userID)
		logrus.Errorf("removeItemFromSale: %v\n", err)
		return ProductModel{}, http.StatusForbidden, err
	}

	if productTableModel.SellerID != userID {
		err := fmt.Errorf("current user's ID %s does not match with Product's SellerID %s", userID, productTableModel.SellerID)
		logrus.Errorf("removeItemFromSale: %v\n", err)
		return ProductModel{}, http.StatusForbidden, err
	}

	if productTableModel.Quantity <= productModel.Quantity {
		productTableModel.Quantity = 0
		if statusCode, err := productTableModel.DeleteProductByID(ctx); err != nil {
			err = fmt.Errorf("exception while Deleting Product for ID:%s. %v", productModel.ID, err)
			logrus.Errorf("removeItemFromSale: %v\n", err)
			return ProductModel{}, statusCode, err
		}
	} else {
		productTableModel.Quantity = productModel.Quantity
		if statusCode, err := productTableModel.UpdateProductByID(ctx); err != nil {
			err = fmt.Errorf("exception while Updating Product for ID:%s. %v", productModel.ID, err)
			logrus.Errorf("removeItemFromSale: %v\n", err)
			return ProductModel{}, statusCode, err
		}
	}
	updatedProduct := convertProductTableModelToProductModel(&productTableModel)
	return *updatedProduct, http.StatusOK, nil
}

func getSellerProducts(ctx context.Context, sessionID string) ([]ProductModel, int, error) {
	userID, userType, statusCode, err := getUserIDAndTypeFromSessionID(ctx, sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("getSellerProducts: %v\n", err)
		return nil, statusCode, err
	}
	if userType != common.Seller {
		err := fmt.Errorf("provided userID is not a seller %s", userID)
		logrus.Errorf("getSellerProducts: %v\n", err)
		return nil, http.StatusForbidden, err
	}
	productTableModel := ProductTableModel{SellerID: userID}
	productTableModels, statusCode, err := productTableModel.GetProductsBySellerID(ctx)
	if err != nil {
		err = fmt.Errorf("exception while fetching Products data: %v", err)
		logrus.Errorf("getSellerProducts: %v\n", err)
		return nil, statusCode, err
	}
	productModels := []ProductModel{}
	for _, p := range productTableModels {
		productModels = append(productModels, *convertProductTableModelToProductModel(&p))
	}

	return productModels, http.StatusOK, nil
}

func incrementProductRating(ctx context.Context, productID string) (int, error) {
	productTableModel := ProductTableModel{ID: productID}
	if statusCode, err := productTableModel.GetProductByID(ctx); err != nil {
		err = fmt.Errorf("exception while fetching Product for ID:%s. %v", productID, err)
		logrus.Errorf("incrementProductRating: %v\n", err)
		return statusCode, err
	}
	productTableModel.FeedBackThumbsUp++
	if statusCode, err := productTableModel.UpdateProductByID(ctx); err != nil {
		err = fmt.Errorf("exception while Updating Product for ID:%s. %v", productID, err)
		logrus.Errorf("incrementProductRating: %v\n", err)
		return statusCode, err
	}
	if statusCode, err := incrementSellerRating(ctx, productTableModel.SellerID); err != nil {
		err = fmt.Errorf("exception while Updating Seller %s rating for productID:%s. %v", productTableModel.SellerID, productID, err)
		logrus.Errorf("incrementProductRating: %v\n", err)
		return statusCode, err
	}
	return http.StatusOK, nil
}

func decrementProductRating(ctx context.Context, productID string) (int, error) {
	productTableModel := ProductTableModel{ID: productID}
	if statusCode, err := productTableModel.GetProductByID(ctx); err != nil {
		err = fmt.Errorf("exception while fetching Product for ID:%s. %v", productID, err)
		logrus.Errorf("decrementProductRating: %v\n", err)
		return statusCode, err
	}
	productTableModel.FeedBackThumbsDown++
	if statusCode, err := productTableModel.UpdateProductByID(ctx); err != nil {
		err = fmt.Errorf("exception while Updating Product for ID:%s. %v", productID, err)
		logrus.Errorf("decrementProductRating: %v\n", err)
		return statusCode, err
	}
	if statusCode, err := decrementSellerRating(ctx, productTableModel.SellerID); err != nil {
		err = fmt.Errorf("exception while Updating Seller %s rating for productID:%s. %v", productTableModel.SellerID, productID, err)
		logrus.Errorf("decrementProductRating: %v\n", err)
		return statusCode, err
	}
	return http.StatusOK, nil
}

func getProductSellerRatingByProductID(ctx context.Context, productID string) (SellerModel, int, error) {
	productTableModel := ProductTableModel{ID: productID}
	if statusCode, err := productTableModel.GetProductByID(ctx); err != nil {
		err = fmt.Errorf("exception while fetching Product for ID:%s. %v", productID, err)
		logrus.Errorf("getProductSellerRatingByProductID: %v\n", err)
		return SellerModel{}, statusCode, err
	}
	thumbsUp, thumbsDown, statusCode, err := getSellerRatingBySellerID(ctx, productTableModel.SellerID)
	if err != nil {
		err = fmt.Errorf("exception while fetching SellerRating for ProductID:%s. %v", productID, err)
		logrus.Errorf("getProductSellerRatingByProductID: %v\n", err)
		return SellerModel{}, statusCode, err
	}
	sellerModel := SellerModel{
		FeedBackThumbsUp:   thumbsUp,
		FeedBackThumbsDown: thumbsDown,
	}
	return sellerModel, http.StatusOK, nil
}

func convertProductModelToProductTableModel(productModel *ProductModel) *ProductTableModel {
	return &ProductTableModel{
		ID:                 productModel.ID,
		Name:               productModel.Name,
		Category:           productModel.Category,
		Keywords:           productModel.Keywords,
		Condition:          productModel.Condition,
		SalePrice:          productModel.SalePrice,
		SellerID:           productModel.SellerID,
		Quantity:           productModel.Quantity,
		FeedBackThumbsUp:   productModel.FeedBackThumbsUp,
		FeedBackThumbsDown: productModel.FeedBackThumbsDown,
		CreatedAt:          productModel.CreatedAt,
		UpdatedAt:          productModel.UpdatedAt,
	}
}

func convertProductTableModelToProductModel(productTableModel *ProductTableModel) *ProductModel {
	return &ProductModel{
		ID:                 productTableModel.ID,
		Name:               productTableModel.Name,
		Category:           productTableModel.Category,
		Keywords:           productTableModel.Keywords,
		Condition:          productTableModel.Condition,
		SalePrice:          productTableModel.SalePrice,
		SellerID:           productTableModel.SellerID,
		Quantity:           productTableModel.Quantity,
		FeedBackThumbsUp:   productTableModel.FeedBackThumbsUp,
		FeedBackThumbsDown: productTableModel.FeedBackThumbsDown,
		CreatedAt:          productTableModel.CreatedAt,
		UpdatedAt:          productTableModel.UpdatedAt,
	}
}

func validateProductModel(ctx context.Context, productModel *ProductModel, create bool) error {

	if !create && productModel.ID == "" {
		err := fmt.Errorf("invalid Product data. ID field is empty")
		logrus.Errorf("validateProductModel: %v\n", err)
		return err
	}

	if productModel.Name == "" {
		err := fmt.Errorf("invalid Product data. Name field is empty")
		logrus.Errorf("validateProductModel: %v\n", err)
		return err
	}

	if len(productModel.Keywords) > 5 {
		err := fmt.Errorf("invalid Product data. Only 5 Keywords allowed")
		logrus.Errorf("validateProductModel: %v\n", err)
		return err
	}

	for _, keyword := range productModel.Keywords {
		if len(keyword) > 8 {
			err := fmt.Errorf("invalid Product data. Each Keywords is allowed max 8 characters: %s", keyword)
			logrus.Errorf("validateProductModel: %v\n", err)
			return err
		}
	}

	if productModel.SalePrice == 0 {
		err := fmt.Errorf("invalid Product data. SalePrice field is empty")
		logrus.Errorf("validateProductModel: %v\n", err)
		return err
	}

	if productModel.Quantity == 0 {
		err := fmt.Errorf("invalid Product data. Quantity field is empty")
		logrus.Errorf("validateProductModel: %v\n", err)
		return err
	}

	if create {
		productModel.FeedBackThumbsUp = 0
		productModel.FeedBackThumbsDown = 0
	}

	return nil
}
