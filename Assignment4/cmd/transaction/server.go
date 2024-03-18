package main

import (
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/library/common"
	"log"
	"net/http"
	"strconv"

	"github.com/adarshsrinivasan/DS_S24/library/wsdl/transaction"
)

const (
	ServiceName   = "transaction"
	ServerHostEnv = "SERVER_HOST"
	ServerPortEnv = "SERVER_PORT"
)

var (
	httpServerHost    = common.GetEnv(ServerHostEnv, "localhost")
	httpServerPort, _ = strconv.Atoi(common.GetEnv(ServerPortEnv, "50003"))
	nodeName          = common.GetEnv(common.NodeNameEnv, "transaction-server")
	baseURL           = fmt.Sprintf("%s:%d", httpServerHost, httpServerPort)
)

func main() {
	http.HandleFunc("/", transaction.Endpoint)
	log.Println("Server Listening ...")
	log.Fatal(http.ListenAndServe(baseURL, nil))
}
