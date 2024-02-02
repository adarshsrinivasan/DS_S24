package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/db"
	"github.com/google/uuid"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	uid  uuid.UUID
	addr string
}

var activeClients map[string]Client

const (
	ServiceName       = "server"
	HttpServerHostEnv = "HTTP_SERVER_HOST"
	HttpServerPortEnv = "HTTP_SERVER_PORT"
)

var (
	err               error
	ctx               context.Context
	httpServerHost    = common.GetEnv(HttpServerHostEnv, "localhost")
	httpServerPort, _ = strconv.Atoi(common.GetEnv(HttpServerPortEnv, "50000"))
)

func initializeSQLDB(ctx context.Context) error {
	logrus.Infof("initializeSQLDB: Initializating SQLDB...\n")
	db.SqlDBClient, err = db.NewSQLClient(ctx, ServiceName)
	if err != nil {
		err = fmt.Errorf("exception while initializing SQLDB buyer. %v", err)
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
	db.NoSQLClient, err = db.NewNoSQLClient(ctx, ServiceName)
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

func initialExchange(conn net.Conn) {
	t := time.Now()
	myTime := t.Format(time.RFC3339Nano)
	clientMsg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}
	log.Println("Message from buyer at ", myTime, ": ", clientMsg)
}

func handleClient(conn net.Conn) {
	if activeClients == nil {
		activeClients = make(map[string]Client, 1000)
	}
	_, ok := activeClients[conn.RemoteAddr().String()]

	if ok == false {
		activeClients[conn.RemoteAddr().String()] = Client{
			uid:  uuid.New(),
			addr: conn.RemoteAddr().String(),
		}
	}
	log.Printf("Client %s connected\n", conn.RemoteAddr().String())
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
		delete(activeClients, conn.RemoteAddr().String())
		log.Println("Client " + conn.RemoteAddr().String() + " disconnected")
	}(conn)

	warning := false
	for {
		requestBody := make([]byte, 1000)
		err := conn.SetDeadline(time.Now().Add(time.Minute * 4))

		clientRequest := common.ClientRequest{}
		_, err = conn.Read(requestBody)

		if errors.Is(err, os.ErrDeadlineExceeded) {
			log.Println("Warning: Sending a inactivity warning to the buyer")

			_ = conn.SetDeadline(time.Now().Add(time.Minute))
			warning = true
			conn.Write([]byte("Session timeout warning: You will be automatically logged out in the next minute\n"))

			_, err = conn.Read(requestBody)

			if errors.Is(err, os.ErrDeadlineExceeded) {
				if warning {
					_ = conn.SetDeadline(time.Now().Add(time.Second * 10))
					log.Printf("Client %s is inactive. Logging out the user\n", conn.RemoteAddr().String())
					conn.Write([]byte("Timeout: Logging you out!\n"))
					return
				}
			}
		}

		if err != nil {
			return
		}

		t := time.Now()
		myTime := t.Format(time.RFC3339Nano)
		log.Println("Received request at", myTime)

		clientRequest.DeserializeRequest(requestBody)
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
			conn, _ := l.Accept()
			handleClient(conn)
			initialExchange(conn)
			go handleConnection(conn)
		}
	}

}
