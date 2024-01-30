package common

type UserType int

const (
	Buyer UserType = iota
	Seller
)

var UserTypeToString = map[UserType]string{
	Buyer:  "Buyer",
	Seller: "Seller",
}

var StringToUserType = map[string]UserType{
	"Buyer":  Buyer,
	"Seller": Seller,
}
