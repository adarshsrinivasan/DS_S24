package common

type UserType int

const (
	NodeNameEnv       = "NODE_NAME"
	PeerNodeNamesEnv  = "PEER_NODE_NAMES"
	SQLNodeNamesEnv   = "SQL_NODE_NAMES"
	SQLNodePortsEnv   = "SQL_NODE_PORTS"
	NOSQLNodeNamesEnv = "NOSQL_NODE_NAMES"
	NOSQLNodePortsEnv = "NOSQL_NODE_PORTS"
)
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
