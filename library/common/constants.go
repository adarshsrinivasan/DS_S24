package common

type UserType int

const (
	BUYER UserType = iota
	SELLER
)

var UserTypeToString = map[UserType]string{
	BUYER:  "Buyer",
	SELLER: "Seller",
}

var StringToUserType = map[string]UserType{
	"Buyer":  BUYER,
	"Seller": SELLER,
}
