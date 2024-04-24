package common

type UserType int

const (
	NodeNameEnv       = "NODE_NAME"
	SQLNodeNamesEnv   = "SQL_NODE_NAMES"
	SQLNodePortsEnv   = "SQL_NODE_PORTS"
	NOSQLNodeNamesEnv = "NOSQL_NODE_NAMES"
	NOSQLNodePortsEnv = "NOSQL_NODE_PORTS"
	PeerNodeNamesEnv  = "PEER_NODE_NAMES"
	PeerNodePortsEnv  = "PEER_NODE_PORTS"
	SyncHostEnv       = "SYNC_HOST"
	SyncPortEnv       = "SYNC_PORT"
	RequestPortEnv    = "REQUEST_PORT"
	ClusterKeyEnv     = "CLUSTER_KEY"
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
