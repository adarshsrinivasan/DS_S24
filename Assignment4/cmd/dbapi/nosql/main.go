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
	ServerHostEnv      = "SERVER_HOST"
	ServerPortEnv      = "SERVER_PORT"
	NOSQLSchemaNameEnv = "MONGO_DB"

	ServiceName = "server"
)

var (
	err             error
	ctx             context.Context
	serverHost      = common.GetEnv(ServerHostEnv, "localhost")
	serverPort, _   = strconv.Atoi(common.GetEnv(ServerPortEnv, "50001"))
	nosqlSchemaName = common.GetEnv(NOSQLSchemaNameEnv, "marketplace")
	nodeName        = common.GetEnv(common.NodeNameEnv, "nosql-server")
	peerNodeNames   = common.SplitCSV(common.GetEnv(common.PeerNodeNamesEnv, "nosql-server1,nosql-server2,nosql-server3,nosql-server4,nosql-server5"))
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

func verifyNOSQLDBConnection(ctx context.Context) error {
	log.Infof("verifyNOSQLDBConnection: Verifying NOSQLDB...\n")
	if err := nosql.VerifyNoSQLConnection(ctx); err != nil {
		err := fmt.Errorf("exception while verifying SQL DB connection. %v", err)
		log.Errorf("VerifyNoSQLConnection: %v\n", err)
		return err
	}
	return nil
}

func verifyDBConnections(ctx context.Context) error {
	log.Infof("verifyDBConnections: Verifying DB connections...\n")
	if err := verifyNOSQLDBConnection(ctx); err != nil {
		err = fmt.Errorf("exception while verifying SQLDB. %v", err)
		log.Errorf("verifyDBConnections: %v\n", err)
		return err
	}
	log.Infof("verifyDBConnections: Verified DB connections Successfully!\n")
	return nil
}

func main() {
	ctx = context.Background()

	if err := verifyDBConnections(ctx); err != nil {
		err = fmt.Errorf("exception while verifying DB connections.... %v", err)
		log.Panicf("main: %v\n", err)
	}

	if err := initialize(ctx, ServiceName, nosqlSchemaName); err != nil {
		err = fmt.Errorf("exception while initializing DB.... %v", err)
		log.Panicf("main: %v\n", err)
	}

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
