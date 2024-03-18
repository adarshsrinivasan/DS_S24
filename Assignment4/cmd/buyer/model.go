package main

import (
	"time"
)


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

type BuyerModel struct {
	Id                     string    `json:"id,omitempty" bson:"id" bun:"id,pk"`
	Name                   string    `json:"name,omitempty" bson:"name" bun:"name,notnull"`
	UserName               string    `json:"userName,omitempty" bson:"userName" bun:"userName,notnull,unique"`
	Password               string    `json:"password,omitempty" bson:"password" bun:"password,notnull,unique"`
	Version                int       `json:"version,omitempty" bson:"version" bun:"version,notnull"`
	CreatedAt              time.Time `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}


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

type SessionModel struct {
	SessionID string          `json:"sessionID,omitempty" bson:"sessionID" bun:"sessionID,pk"`
}

type TransactionModel struct {
	ID        string    `json:"id,omitempty" bson:"id" bun:"id,pk"`
	CartID    string    `json:"cartID,omitempty" bson:"cartID" bun:"cartID,notnull"`
	ProductID string    `json:"productID,omitempty" bson:"productID" bun:"productID,notnull"`
	BuyerID   string    `json:"buyerID,omitempty" bson:"buyerID" bun:"buyerID,notnull"`
	SellerID  string    `json:"sellerID,omitempty" bson:"sellerID" bun:"sellerID,notnull"`
	Quantity  int       `json:"quantity,omitempty" bson:"quantity" bun:"quantity,notnull"`
	Price     float32   `json:"price,omitempty" bson:"price,omitempty" bun:"quantity,notnull"`
	Version   int       `json:"version,omitempty" bson:"version" bun:"version,notnull"`
	CreatedAt time.Time `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}

type PurchaseDetailsModel struct {
	Name string `json:"name,omitempty" bson:"name" bun:"name,pk"`
	CreditCardNumber string `json:"creditCardNumber,omitempty" bson:"creditCardNumber" bun:"creditCardNumber,pk"`
	Expiry string `json:"expiry,omitempty" bson:"expiry" bun:"expiry,pk"`
}

type SellerModel struct {
	Id                 string    `json:"id,omitempty" bson:"id" bun:"id,pk"`
	Name               string    `json:"name,omitempty" bson:"name" bun:"name,notnull"`
	FeedBackThumbsUp   int       `json:"feedBackThumbsUp" bson:"feedBackThumbsUp" bun:"feedBackThumbsUp"`
	FeedBackThumbsDown int       `json:"feedBackThumbsDown" bson:"feedBackThumbsDown" bun:"feedBackThumbsDown"`
	NumberOfItemsSold  int       `json:"numberOfItemsSold,omitempty" bson:"numberOfItemsSold" bun:"numberOfItemsSold"`
	UserName           string    `json:"userName,omitempty" bson:"userName" bun:"userName,notnull,unique"`
	Password           string    `json:"password,omitempty" bson:"password" bun:"password,notnull,unique"`
	Version            int       `json:"version,omitempty" bson:"version" bun:"version,notnull"`
	CreatedAt          time.Time `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}