package main

import (
	"context"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/library/proto"
	"google.golang.org/grpc"
	"net"
	"strconv"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/adarshsrinivasan/DS_S24/library/db/sql"
	log "github.com/sirupsen/logrus"
)

const (
	ServerHostEnv    = "SERVER_HOST"
	ServerPortEnv    = "SERVER_PORT"
	SQLSchemaNameEnv = "POSTGRES_DB"

	ServiceName = "server"
)

var (
	err           error
	ctx           context.Context
	serverHost    = common.GetEnv(ServerHostEnv, "localhost")
	serverPort, _ = strconv.Atoi(common.GetEnv(ServerPortEnv, "50002"))
	nodeName      = common.GetEnv(common.NodeNameEnv, "sql-server")
	peerNodeNames = common.SplitNodeNames(common.GetEnv(common.PeerNodeNamesEnv, "sql-server1,sql-server2,sql-server3,sql-server4,sql-server5"))
	serviceName   string
	schemaName    = common.GetEnv(SQLSchemaNameEnv, "marketplace")
)

func initializeSQLDB(ctx context.Context, serviceName, schemaName string) error {
	log.Infof("initializeSQLDB: Initializating SQLDB...\n")
	client, err := sql.NewSQLClient(ctx, serviceName, schemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		log.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	defer client.Close(ctx)

	if err := client.Initialize(ctx, schemaName); err != nil {
		err = fmt.Errorf("exception while initializing SQLDB client. %v", err)
		log.Errorf("initializeSQLDB: %v\n", err)
		return err
	}

	if err := CreateSellerTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating seller tabel. %v", err)
		log.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	if err := CreateBuyerTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating buyer tabel. %v", err)
		log.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	if err := CreateSessionTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating session tabel. %v", err)
		log.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	if err := CreateCartTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating cart tabel. %v", err)
		log.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	if err := CreateCartItemTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating cartItem tabel. %v", err)
		log.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	if err := CreateTransactionTable(ctx); err != nil {
		err = fmt.Errorf("exception while creating transaction tabel. %v", err)
		log.Errorf("initializeSQLDB: %v\n", err)
		return err
	}
	log.Infof("initializeSQLDB: Initialized SQLDB Successfully!\n")
	return nil
}

func initializeDB(ctx context.Context, serviceName, schemaName string) error {
	log.Infof("initializeDB: Initializating DB...\n")
	if err := initializeSQLDB(ctx, serviceName, schemaName); err != nil {
		err = fmt.Errorf("exception while initializing SQLDB. %v", err)
		log.Errorf("initializeDB: %v\n", err)
		return err
	}
	log.Infof("initializeDB: Initialized all DB Successfully!\n")
	return nil
}

func initialize(ctx context.Context, serviceName, schemaName string) error {
	log.Infof("initialize: Initializating...\n")
	if err := initializeDB(ctx, serviceName, schemaName); err != nil {
		err = fmt.Errorf("exception while initializing DBs. %v", err)
		log.Errorf("initialize: %v\n", err)
		return err
	}

	log.Infof("initialize: Initialization completed Successfully!\n")
	return nil
}

func verifySQLDBConnection(ctx context.Context) error {
	log.Infof("verifySQLDBConnection: Verifying SQLDB...\n")
	if err := sql.VerifySQLConnection(ctx); err != nil {
		err := fmt.Errorf("exception while verifying SQL DB connection. %v", err)
		log.Errorf("verifySQLDBConnection: %v\n", err)
		return err
	}
	return nil
}

func verifyDBConnections(ctx context.Context) error {
	log.Infof("verifyDBConnections: Verifying DB connections...\n")
	if err := verifySQLDBConnection(ctx); err != nil {
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

	if err := initialize(ctx, ServiceName, schemaName); err != nil {
		err = fmt.Errorf("exception while initializing DB.... %v", err)
		log.Panicf("main: %v\n", err)
	}

	log.Println("Server Listening ...")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverHost, serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterSQLServiceServer(server, &sqlServer{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
