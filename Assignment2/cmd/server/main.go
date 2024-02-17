package main

import (
	"context"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/library/proto"
	"github.com/adarshsrinivasan/DS_S24/library/wsdl/transaction"
	"github.com/hooklift/gowsdl/soap"
	"log"
	"net/http"
	"strconv"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/sirupsen/logrus"
)

const (
	ServiceName        = "server"
	ServerHostEnv      = "SERVER_HOST"
	ServerPortEnv      = "SERVER_PORT"
	SQLRPCHostEnv      = "SQL_RPC_HOST"
	SQLRPCPortEnv      = "SQL_RPC_PORT"
	NOSQLRPCHostEnv    = "NOSQL_RPC_HOST"
	NOSQLRPCPortEnv    = "NOSQL_RPC_PORT"
	SQLSchemaName      = "marketplace"
	NOSQLSchemaNameEnv = "MONGO_DB"
	TransactionHostEnv = "TRANSACTION_HOST"
	TransactionPortEnv = "TRANSACTION_PORT"
)

var (
	err               error
	ctx               context.Context
	nosqlSchemaName   = common.GetEnv(NOSQLSchemaNameEnv, "marketplace")
	httpServerHost    = common.GetEnv(ServerHostEnv, "localhost")
	httpServerPort, _ = strconv.Atoi(common.GetEnv(ServerPortEnv, "50000"))
	sqlRPCHost        = common.GetEnv(SQLRPCHostEnv, "localhost")
	sqlRPCPort, _     = strconv.Atoi(common.GetEnv(SQLRPCPortEnv, "50002"))
	nosqlRPCHost      = common.GetEnv(NOSQLRPCHostEnv, "localhost")
	nosqlRPCPort, _   = strconv.Atoi(common.GetEnv(NOSQLRPCPortEnv, "50001"))
	transactionSoapHost      = common.GetEnv(TransactionHostEnv, "localhost")
	transactionSoapPort, _   = strconv.Atoi(common.GetEnv(TransactionPortEnv, "50003"))
	transactionService transaction.TransactionServicePortType
)

func initializeTransactionServiceClient()  {
	logrus.Infof("initializeTransactionServiceClient: Initializing...\n")
	client := soap.NewClient(fmt.Sprintf("http://%s:%d", transactionSoapHost, transactionSoapPort))
	transactionService = transaction.NewTransactionServicePortType(client)
	logrus.Infof("initializeTransactionServiceClient: Initialized Successfully\n")
}

func initializeSQLDB(ctx context.Context) error {
	logrus.Infof("initializeSQLDB: Initializating SQLDB...\n")
	// Set up a connection to the server.
	sqlDBClient, conn, err := common.NewSQLRPCClient(ctx, sqlRPCHost, sqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to NOSQLDB RPC server. %v", err)
		logrus.Errorf("initializeNOSQLDB: %v\n", err)
		return err
	}
	defer conn.Close()
	initializeRequest := &proto.InitializeRequest{
		ServiceName:   ServiceName,
		SQLSchemaName: SQLSchemaName,
	}

	if _, err := sqlDBClient.Initialize(ctx, initializeRequest); err != nil {
		err = fmt.Errorf("exception while initializing SQLDB RPC client. %v", err)
		logrus.Errorf("initializeSQLDB: %v\n", err)
		return err
	}

	logrus.Infof("initializeSQLDB: Initialized SQLDB Successfully!\n")
	return nil
}

func initializeNOSQLDB(ctx context.Context) error {
	logrus.Infof("initializeNOSQLDB: Initializating NOSQLDB...\n")
	nosqlDBClient, conn, err := common.NewNOSQLRPCClient(ctx, nosqlRPCHost, nosqlRPCPort)
	if err != nil {
		err = fmt.Errorf("exception while connecting to NOSQLDB RPC server. %v", err)
		logrus.Errorf("initializeNOSQLDB: %v\n", err)
		return err
	}
	defer conn.Close()
	initializeRequest := &proto.InitializeRequest{
		ServiceName:   ServiceName,
		SQLSchemaName: nosqlSchemaName,
	}

	if _, err := nosqlDBClient.Initialize(ctx, initializeRequest); err != nil {
		err = fmt.Errorf("exception while initializing NOSQLDB RPC client. %v", err)
		logrus.Errorf("initializeNOSQLDB: %v\n", err)
		return err
	}
	logrus.Infof("initializeNOSQLDB: Initialized NOSQLDB Successfully!\n")
	return nil
}

func initializeDB(ctx context.Context) error {
	logrus.Infof("initializeDB: Initializating DB...\n")
	if err := initializeSQLDB(ctx); err != nil {
		err = fmt.Errorf("exception while initializing SQLDB. %v", err)
		logrus.Errorf("initializeDB: %v\n", err)
		return err
	}
	if err := initializeNOSQLDB(ctx); err != nil {
		err = fmt.Errorf("exception while initializing NOSQLDB. %v", err)
		logrus.Errorf("initializeDB: %v\n", err)
		return err
	}
	logrus.Infof("initializeDB: Initialized all DB Successfully!\n")
	return nil
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

	if err := initializeDB(ctx); err != nil {
		err = fmt.Errorf("exception while initializing DBs. %v", err)
		logrus.Errorf("initialize: %v\n", err)
		return err
	}
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
