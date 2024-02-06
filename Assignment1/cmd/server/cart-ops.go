package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type CartModel struct {
	ID         string          `json:"id,omitempty" bson:"id" bun:"id,pk"`
	BuyerID    string          `json:"buyerID,omitempty" bson:"buyerID" bun:"buyerID,notnull,unique"`
	Saved      bool            `json:"saved,omitempty" bson:"saved" bun:"saved,notnull"`
	Items      []CartItemModel `json:"items,omitempty" bson:"items"`
	TotalPrice float32         `json:"totalPrice,omitempty" bson:"totalPrice,omitempty" bun:"totalPrice,notnull"`
	Version    int             `json:"version,omitempty" bson:"version" bun:"version,notnull"`
	CreatedAt  time.Time       `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}

type CartItemModel struct {
	ID        string    `json:"id,omitempty" bson:"id" bun:"id,pk"`
	CartID    string    `json:"cartID,omitempty" bson:"cartID" bun:"cartID,notnull"`
	ProductID string    `json:"productID,omitempty" bson:"productID" bun:"productID,notnull"`
	SellerID  string    `json:"sellerID,omitempty" bson:"sellerID" bun:"sellerID,notnull"`
	Quantity  int       `json:"quantity,omitempty" bson:"quantity" bun:"quantity,notnull"`
	Price     float32   `json:"price,omitempty" bson:"price,omitempty" bun:"price,notnull"`
	Version   int       `json:"version,omitempty" bson:"version" bun:"version,notnull"`
	CreatedAt time.Time `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}

func createCart(ctx context.Context, buyerID string) (int, error) {
	cartTableModel := CartTableModel{
		BuyerID: buyerID,
		Saved:   false,
	}
	return cartTableModel.CreateCart(ctx)
}

func getCartByID(ctx context.Context, cartID string) (*CartModel, int, error) {
	cartTableModel := CartTableModel{
		ID: cartID,
	}
	if statusCode, err := cartTableModel.GetCartByID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Cart with ID %s. %v", cartID, err)
		logrus.Errorf("getCartByID: %v\n", err)
		return nil, statusCode, err
	}
	cartItemTableModel := CartItemTableModel{
		CartID: cartID,
	}
	cartItemTableModelList, statusCode, err := cartItemTableModel.ListCartItemByCartID(ctx)
	if err != nil {
		err := fmt.Errorf("exception while fetching CartItems for ID %s. %v", cartID, err)
		logrus.Errorf("getCartByID: %v\n", err)
		return nil, statusCode, err
	}

	return buildCartModel(&cartTableModel, cartItemTableModelList), http.StatusOK, nil
}

func addProductToCart(ctx context.Context, cartItemModel *CartItemModel) (int, error) {
	_, _, err := getCartByID(ctx, cartItemModel.CartID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Cart with ID %s. %v", cartItemModel.CartID, err)
		logrus.Errorf("addProductToCart: %v\n", err)
		return http.StatusBadRequest, err
	}
	cartItemTableModel := CartItemTableModel{
		CartID:    cartItemModel.CartID,
		ProductID: cartItemModel.ProductID,
		SellerID:  cartItemModel.SellerID,
		Quantity:  cartItemModel.Quantity,
		Price:     cartItemModel.Price,
	}
	if statusCode, err := cartItemTableModel.CreateCartItem(ctx); err != nil {
		err := fmt.Errorf("exception while Creating CartItems for ID %s. %v", cartItemModel.CartID, err)
		logrus.Errorf("addProductToCart: %v\n", err)
		return statusCode, err
	}
	return http.StatusOK, nil
}

func updateCartItemByCartIDAndProductID(ctx context.Context, cartItemModel *CartItemModel) (int, error) {
	_, _, err := getCartByID(ctx, cartItemModel.CartID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Cart with ID %s. %v", cartItemModel.CartID, err)
		logrus.Errorf("addProductToCart: %v\n", err)
		return http.StatusBadRequest, err
	}
	cartItemTableModel := CartItemTableModel{
		CartID:    cartItemModel.CartID,
		ProductID: cartItemModel.ProductID,
	}
	if statusCode, err := cartItemTableModel.GetCartItemByCartIDAndProductID(ctx); err != nil {
		err := fmt.Errorf("exception while Fetching CartItems for ID %s. %v", cartItemModel.CartID, err)
		logrus.Errorf("addProductToCart: %v\n", err)
		return statusCode, err
	}
	cartItemTableModel.Quantity = cartItemModel.Quantity
	cartItemTableModel.Price = cartItemModel.Price
	if statusCode, err := cartItemTableModel.UpdateCartItem(ctx); err != nil {
		err := fmt.Errorf("exception while Updating CartItem with CartID %s. %v", cartItemTableModel.CartID, err)
		logrus.Errorf("addProductToCart: %v\n", err)
		return statusCode, err
	}

	return http.StatusOK, nil
}

func removeProductToCart(ctx context.Context, cartItemModel *CartItemModel) (int, error) {
	_, _, err := getCartByID(ctx, cartItemModel.CartID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Cart with ID %s. %v", cartItemModel.CartID, err)
		logrus.Errorf("removeProductToCart: %v\n", err)
		return http.StatusBadRequest, err
	}
	existingCartItemModel := CartItemTableModel{
		CartID:    cartItemModel.CartID,
		ProductID: cartItemModel.ProductID,
	}
	if statusCode, err := existingCartItemModel.GetCartItemByCartIDAndProductID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching CartItem with CartID %s. %v", existingCartItemModel.CartID, err)
		logrus.Errorf("removeProductToCart: %v\n", err)
		return statusCode, err
	}
	if existingCartItemModel.Quantity <= cartItemModel.Quantity {
		logrus.Errorf("removeProductToCart: Removing CartItem %s from Cart %s.\n", existingCartItemModel.ID, existingCartItemModel.CartID)
		if statusCode, err := existingCartItemModel.DeleteCartItemByCartID(ctx); err != nil {
			err := fmt.Errorf("exception while Deleting CartItem with CartID %s. %v", existingCartItemModel.CartID, err)
			logrus.Errorf("removeProductToCart: %v\n", err)
			return statusCode, err
		}
	} else {
		existingCartItemModel.Quantity -= cartItemModel.Quantity

		if statusCode, err := existingCartItemModel.UpdateCartItem(ctx); err != nil {
			err := fmt.Errorf("exception while Updating CartItem with CartID %s. %v", existingCartItemModel.CartID, err)
			logrus.Errorf("removeProductToCart: %v\n", err)
			return statusCode, err
		}
	}
	return http.StatusOK, nil
}

func saveCart(ctx context.Context, cartID string) (int, error) {
	cartTableModel := CartTableModel{
		ID: cartID,
	}
	if statusCode, err := cartTableModel.GetCartByID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Cart with ID %s. %v", cartID, err)
		logrus.Errorf("saveCart: %v\n", err)
		return statusCode, err
	}
	if !cartTableModel.Saved {
		cartTableModel.Saved = true
		if statusCode, err := cartTableModel.UpdateCartByID(ctx); err != nil {
			err := fmt.Errorf("exception while Updating Cart with ID %s. %v", cartTableModel.ID, err)
			logrus.Errorf("removeProductToCart: %v\n", err)
			return statusCode, err
		}
	}
	return http.StatusOK, nil
}

func clearCart(ctx context.Context, cartID string) (int, error) {
	cartTableModel := CartTableModel{
		ID: cartID,
	}
	if statusCode, err := cartTableModel.GetCartByID(ctx); err != nil {
		err := fmt.Errorf("exception while fetching Cart with ID %s. %v", cartID, err)
		logrus.Errorf("clearCart: %v\n", err)
		return statusCode, err
	}
	cartItemTableModel := CartItemTableModel{CartID: cartID}
	if statusCode, err := cartItemTableModel.DeleteCartItemByCartID(ctx); err != nil {
		err := fmt.Errorf("exception while Deleting CartItems with cartID %s. %v", cartItemTableModel.CartID, err)
		logrus.Errorf("clearCart: %v\n", err)
		return statusCode, err
	}
	if statusCode, err := cartTableModel.DeleteCartByID(ctx); err != nil {
		err := fmt.Errorf("exception while Deleting Cart with ID %s. %v", cartTableModel.ID, err)
		logrus.Errorf("clearCart: %v\n", err)
		return statusCode, err
	}

	return http.StatusOK, nil
}

func buildCartModel(cartTableModel *CartTableModel, cartItemTableModelList []CartItemTableModel) *CartModel {
	cartModel := &CartModel{
		ID:         cartTableModel.ID,
		BuyerID:    cartTableModel.BuyerID,
		Saved:      cartTableModel.Saved,
		Items:      make([]CartItemModel, 0),
		TotalPrice: 0,
		Version:    cartTableModel.Version,
		CreatedAt:  cartTableModel.CreatedAt,
		UpdatedAt:  cartTableModel.UpdatedAt,
	}

	for _, cartItemTableModel := range cartItemTableModelList {
		cartItemModel := CartItemModel{
			ID:        cartItemTableModel.ID,
			CartID:    cartItemTableModel.CartID,
			ProductID: cartItemTableModel.ProductID,
			SellerID:  cartItemTableModel.SellerID,
			Quantity:  cartItemTableModel.Quantity,
			Price:     cartItemTableModel.Price,
			Version:   cartItemTableModel.Version,
			CreatedAt: cartItemTableModel.CreatedAt,
			UpdatedAt: cartItemTableModel.UpdatedAt,
		}
		cartModel.Items = append(cartModel.Items, cartItemModel)
		cartModel.TotalPrice += cartItemModel.Price
	}
	return cartModel
}
