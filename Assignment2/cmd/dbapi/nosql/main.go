package main

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/adarshsrinivasan/DS_S24/library/db/nosql"
	"github.com/adarshsrinivasan/DS_S24/library/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	ServiceName   = "server"
	ServerHostEnv = "SERVER_HOST"
	ServerPortEnv = "SERVER_PORT"
)

var (
	err           error
	ctx           context.Context
	serverHost    = common.GetEnv(ServerHostEnv, "localhost")
	serverPort, _ = strconv.Atoi(common.GetEnv(ServerPortEnv, "50001"))
)

func initializeNOSQLDB(ctx context.Context, serviceName, schemaName string) error {
	log.Infof("initializeNOSQLDB: Initializating NOSQLDB...\n")
	nosql.Client, err = nosql.NewNoSQLClient(ctx, serviceName, schemaName)
	if err != nil {
		err = fmt.Errorf("exception while initializing NOSQLDB buyer. %v", err)
		log.Errorf("initializeNOSQLDB: %v\n", err)
		return err
	}
	if err := CreateProductTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating product tabel. %v", err)
		log.Errorf("initializeNOSQLDB: %v\n", err)
		return err
	}
	log.Infof("initializeNOSQLDB: Initialized NOSQLDB Successfully!\n")
	return nil
}

func initializeDB(ctx context.Context, serviceName, schemaName string) error {
	log.Infof("initializeDB: Initializating DB...\n")
	if err := initializeNOSQLDB(ctx, serviceName, schemaName); err != nil {
		err = fmt.Errorf("exception while initializing NOSQLDB. %v", err)
		log.Errorf("initializeDB: %v\n", err)
		return err
	}
	log.Infof("initializeDB: Initialized all DB Successfully!\n")
	return nil
}

func initialize(ctx context.Context, serviceName, schemaName string) error {
	log.Infof("initialize: Initializating...\n")
	ctx = context.Background()

	if err := initializeDB(ctx, serviceName, schemaName); err != nil {
		err = fmt.Errorf("exception while initializing DBs. %v", err)
		log.Errorf("initialize: %v\n", err)
		return err
	}

	log.Infof("initialize: Initialization completed Successfully!\n")
	return nil
}

func main() {
	log.Println("Server Listening ...")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverHost, serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterNOSQLServiceServer(server, &noSQLServer{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
