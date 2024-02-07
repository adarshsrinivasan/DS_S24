package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/db/nosql"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/db/sql"
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

func initialize() error {
	logrus.Infof("initialize: Initializating...\n")
	ctx = context.Background()

	if err := initializeDB(ctx); err != nil {
		err = fmt.Errorf("exception while initializing DBs. %v", err)
		logrus.Errorf("initialize: %v\n", err)
		return err
	}

	logrus.Infof("initialize: Initialization completed Successfully!\n")
	return nil
}

func initialExchange(clientReader *bufio.Reader) error {
	t := time.Now()
	myTime := t.Format(time.RFC3339Nano)
	clientMsg, err := clientReader.ReadString('\n')
	if err != nil {
		return err
	}
	log.Println("Message from client at ", myTime, ": ", clientMsg)
	return nil
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
		log.Println("Client " + conn.RemoteAddr().String() + " disconnected")
	}(conn)

	warning := false
	clientReader := bufio.NewReader(conn)

	if err := initialExchange(clientReader); err != nil {
		logrus.Errorf("Unable to read Client %s initial message. Logging out the user", conn.RemoteAddr().String())
		common.RespondWithError(conn, http.StatusGatewayTimeout, "Timeout: Logging you out!\n")
		return
	}

	for {
		requestBody := make([]byte, 1000)
		err := conn.SetDeadline(time.Now().Add(time.Minute * 4))

		clientRequest := common.ClientRequest{}

		_, err = clientReader.Read(requestBody)

		if errors.Is(err, os.ErrDeadlineExceeded) {
			log.Println("Warning: Sending a inactivity warning to the buyer")

			_ = conn.SetDeadline(time.Now().Add(time.Minute))
			warning = true
			common.RespondWithError(conn, http.StatusContinue, "Session timeout warning: You will be automatically logged out in the next minute\n")

			_, err = clientReader.Read(requestBody)

			if errors.Is(err, os.ErrDeadlineExceeded) {
				if warning {
					_ = conn.SetDeadline(time.Now().Add(time.Second * 10))
					log.Printf("Client %s is inactive. Logging out the user\n", conn.RemoteAddr().String())
					common.RespondWithError(conn, http.StatusGatewayTimeout, "Timeout: Logging you out!\n")
					return
				}
			}
		}

		if err != nil {
			return
		}

		t := time.Now()
		myTime := t.Format(time.RFC3339Nano)

		clientRequest.DeserializeRequest(requestBody)
		log.Printf("Received request: %s at: %v", clientRequest.String(), myTime)
		if clientRequest.UserType == common.Seller {
			listOfSellerHandlers(ctx, conn, clientRequest)
		} else {
			listOfBuyerHandlers(ctx, conn, clientRequest)
		}
	}
}

func main() {
	if err := initialize(); err != nil {
		err = fmt.Errorf("exception while initializing.... %v", err)
		logrus.Panicf("main: %v\n", err)
	}
	log.Println("Server Listening ...")
	if l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", httpServerHost, httpServerPort)); err != nil {
		logrus.Fatal("ERROR: Server listening failed.")
	} else {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go handleConnection(conn)
		}
	}

}
