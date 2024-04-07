package main

import (
	"context"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/library/wsdl/transaction"
	"github.com/hooklift/gowsdl/soap"
	"log"
	"net/http"
	"strconv"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/sirupsen/logrus"
)

const (
	ServerHostEnv      = "SERVER_HOST"
	ServerPortEnv      = "SERVER_PORT"
	TransactionHostEnv = "TRANSACTION_HOST"
	TransactionPortEnv = "TRANSACTION_PORT"
)

var (
	err                        error
	ctx                        context.Context
	httpServerHost             = common.GetEnv(ServerHostEnv, "localhost")
	httpServerPort, _          = strconv.Atoi(common.GetEnv(ServerPortEnv, "50000"))
	nosqlNodeNames             = common.SplitCSV(common.GetEnv(common.NOSQLNodeNamesEnv, "localhost"))
	nosqlNodePorts             = common.SplitCSV(common.GetEnv(common.NOSQLNodePortsEnv, "50003"))
	sqlNodeNames               = common.SplitCSV(common.GetEnv(common.SQLNodeNamesEnv, "localhost"))
	sqlNodePorts               = common.SplitCSV(common.GetEnv(common.SQLNodePortsEnv, "50002"))
	sqlRPCHost, sqlRPCPort     = getSQLHostNameAndPort()
	nosqlRPCHost, nosqlRPCPort = getNOSQLHostNameAndPort()
	transactionSoapHost        = common.GetEnv(TransactionHostEnv, "localhost")
	transactionSoapPort, _     = strconv.Atoi(common.GetEnv(TransactionPortEnv, "50003"))
	nodeName                   = common.GetEnv(common.NodeNameEnv, "server")
	transactionService         transaction.TransactionServicePortType
)

func getSQLHostNameAndPort() (string, int) {
	sqlNodeName, sqlNodePort := common.GetRandomHostAndPort(sqlNodeNames, sqlNodePorts)
	logrus.Infof("getSQLHostName: HostName: %s, Port: %d\n", sqlNodeName, sqlNodePort)
	return sqlNodeName, sqlNodePort
}

func getNOSQLHostNameAndPort() (string, int) {
	nosqlNodeName, nosqlNodePort := common.GetRandomHostAndPort(nosqlNodeNames, nosqlNodePorts)
	logrus.Infof("getNOSQLHostNameAndPort: HostName: %s, Port: %d\n", nosqlNodeName, nosqlNodePort)
	return nosqlNodeName, nosqlNodePort
}

func initializeTransactionServiceClient() {
	logrus.Infof("initializeTransactionServiceClient: Initializing...\n")
	client := soap.NewClient(fmt.Sprintf("http://%s:%d", transactionSoapHost, transactionSoapPort))
	transactionService = transaction.NewTransactionServicePortType(client)
	logrus.Infof("initializeTransactionServiceClient: Initialized Successfully\n")
}

func initializeHTTPRouter(ctx context.Context) error {
	initializeHttpRoutes(ctx)
	if httpRouter == nil {
		return fmt.Errorf("http router not initialized")
	}
	return nil
}

func initialize() error {
	logrus.Infof("initialize: Initializating...\n")
	ctx = context.Background()

	if err := initializeHTTPRouter(ctx); err != nil {
		err = fmt.Errorf("exception while initializing HTTP Router. %v", err)
		logrus.Errorf("initialize: %v\n", err)
		return err
	}

	initializeTransactionServiceClient()

	logrus.Infof("initialize: Initialization completed Successfully!\n")
	return nil
}

func main() {
	if err := initialize(); err != nil {
		err = fmt.Errorf("exception while initializing.... %v", err)
		logrus.Panicf("main: %v\n", err)
	}
	log.Println("Server Listening ...")
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", httpServerHost, httpServerPort), httpRouter))

}
