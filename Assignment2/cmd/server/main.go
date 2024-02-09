package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/adarshsrinivasan/DS_S24/libraries/common"
	"github.com/adarshsrinivasan/DS_S24/libraries/db/nosql"
	"github.com/adarshsrinivasan/DS_S24/libraries/db/sql"
	"github.com/sirupsen/logrus"
)

const (
	ServiceName   = "server"
	ServerHostEnv = "SERVER_HOST"
	ServerPortEnv = "SERVER_PORT"
	SQLSchemaName = "marketplace"
)

var (
	err               error
	ctx               context.Context
	httpServerHost    = common.GetEnv(ServerHostEnv, "localhost")
	httpServerPort, _ = strconv.Atoi(common.GetEnv(ServerPortEnv, "50000"))
)

func initializeSQLDB(ctx context.Context) error {
	logrus.Infof("initializeSQLDB: Initializating SQLDB...\n")
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	defer client.Close(ctx)

	if err := client.Initialize(ctx, SQLSchemaName); err != nil {
		err = fmt.Errorf("exception while initializing SQLDB client. %v", err)
		logrus.Errorf("initializeSQLDB: %v\n", err)
		return err
	}

	if err := CreateSellerTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating seller tabel. %v", err)
		logrus.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	if err := CreateBuyerTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating buyer tabel. %v", err)
		logrus.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	if err := CreateSessionTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating session tabel. %v", err)
		logrus.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	if err := CreateCartTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating cart tabel. %v", err)
		logrus.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	if err := CreateCartItemTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating cartItem tabel. %v", err)
		logrus.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	if err := CreateTransactionTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating transaction tabel. %v", err)
		logrus.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	logrus.Infof("initializeSQLDB: Initialized SQLDB Successfully!\n")
	return nil
}

func initializeNOSQLDB(ctx context.Context) error {
	logrus.Infof("initializeNOSQLDB: Initializating NOSQLDB...\n")
	nosql.Client, err = nosql.NewNoSQLClient(ctx, ServiceName)
	if err != nil {
		err = fmt.Errorf("exception while initializing NOSQLDB buyer. %v", err)
		logrus.Errorf("initializeNOSQLDB: %v\n", err)
		return err
	}
	if err := CreateProductTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating product tabel. %v", err)
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
