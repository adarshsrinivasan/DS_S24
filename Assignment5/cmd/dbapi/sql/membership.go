package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	libProto "github.com/adarshsrinivasan/DS_S24/library/proto"
	"github.com/hashicorp/memberlist"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

var (
	node                        *Node
	responseTrackers            = map[string]chan bool{}
	localCounter, globalCounter atomic.Int32
)

type opsType int

const (
	CreateBuyer opsType = iota
	UpdateBuyerByID
	CreateCart
	UpdateCartByID
	DeleteCartByID
	CreateCartItem
	UpdateCartItem
	DeleteCartItemByCartIDAndProductID
	DeleteCartItemByCartID
	DeleteCartItemByProductID
	CreateSeller
	UpdateSellerByID
	CreateSession
	DeleteSessionByID
	CreateTransaction
	DeleteTransactionsByCartID
	DeleteTransactionsByBuyerID
	DeleteTransactionsBySellerID
)

var opsTypeToStr = map[opsType]string{
	CreateBuyer:                        "CreateBuyer",
	UpdateBuyerByID:                    "UpdateBuyerByID",
	CreateCart:                         "CreateCart",
	UpdateCartByID:                     "UpdateCartByID",
	DeleteCartByID:                     "DeleteCartByID",
	CreateCartItem:                     "CreateCartItem",
	UpdateCartItem:                     "UpdateCartItem",
	DeleteCartItemByCartIDAndProductID: "DeleteCartItemByCartIDAndProductID",
	DeleteCartItemByCartID:             "DeleteCartItemByCartID",
	DeleteCartItemByProductID:          "DeleteCartItemByProductID",
	CreateSeller:                       "CreateSeller",
	UpdateSellerByID:                   "UpdateSellerByID",
	CreateSession:                      "CreateSession",
	DeleteSessionByID:                  "DeleteSessionByID",
	CreateTransaction:                  "CreateTransaction",
	DeleteTransactionsByCartID:         "DeleteTransactionsByCartID",
	DeleteTransactionsByBuyerID:        "DeleteTransactionsByBuyerID",
	DeleteTransactionsBySellerID:       "DeleteTransactionsBySellerID",
}

type message struct {
	ID              string  `json:"id"`
	OpsType         opsType `json:"opsType"`
	Payload         []byte  `json:"payload"`
	RequestNodeName string  `json:"requestNodeName"`
	Sequence        int32   `json:"sequence"`
}

func (m message) toString() string {
	return string(marshallMsg(ctx, &m))
}

type Node struct {
	memberlist  *memberlist.Memberlist
	host        string
	port        string
	requestPort string
	started     bool
}

type metadata struct {
	RequestPort    string `json:"request_port"`
	SequenceNumber int32  `json:"sequence_number"`
}

func (m metadata) toString() string {
	data, _ := json.Marshal(m)
	return string(data)
}

func (m metadata) fromString(data string) error {
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return err
	}
	return nil
}

func (n *Node) ListenToNodesData() {

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", n.memberlist.LocalNode().Addr.String(), n.requestPort))
	if err != nil {
		log.Fatalf("ListenToNodesData: Cannot start the cluster node: %v", err)
	}
	defer l.Close()
	log.Printf("ListenToNodesData: Started the cluster node on: %s:%s", n.memberlist.LocalNode().Addr.String(), n.requestPort)
	for {
		conn, err := l.Accept() // wait for connection from other nodes and this connection the other nodes will be pushing the data
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func (n *Node) ListenNodeLeave() {
	// Create a channel to listen for exit signals
	stop := make(chan os.Signal, 1)
	// Register the signals we want to be notified, these 3 indicate exit
	// signals, similar to CTRL+C
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-stop
	// Leave the cluster with a 5 second timeout. If leaving takes more than 5
	// seconds we return.
	if err := n.memberlist.Leave(time.Second * 5); err != nil {
		panic(err)
	}
	os.Exit(1)
}

func (n *Node) PushDataToOtherNodes(data []byte) {
	for _, m := range n.memberlist.Members() { //iterate over the member list   and push data to other cluster nodes.
		if m == n.memberlist.LocalNode() { //its the localnode so we don’t want to use this data
			continue
		}
		var conn net.Conn
		var err error
		mdata := metadata{}
		if err := mdata.fromString(string(m.Meta)); err != nil {
			log.Errorf("PushDataToOtherNodes: exception while parsing metadata of server %s:%s: %v", m.Address(), mdata.RequestPort, err)
			continue
		}
		log.Infof("PushDataToOtherNodes: Dialing: %s:%s", m.Address(), mdata.RequestPort)
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", m.Addr.String(), mdata.RequestPort))
		if err != nil {
			log.Errorf("PushDataToOtherNodes: could not connect to server %s:%s: %v", m.Address(), mdata.RequestPort, err)
			continue
		}
		_, err = conn.Write(data) //send the data we received from client to other nodes
		if err != nil {
			log.Errorf("Cluster down with address: %s:%s", m.Address(), mdata.RequestPort)
		}
	}
}

func initMemberNode(host, port, masterHost, masterPort, requestPort, clusterKey string) error {
	if node != nil {
		return fmt.Errorf("initMemberNode: Node already initialized")
	}
	//clusterKey := make([]byte, 32)
	//_, err := rand.Read(clusterKey)
	//if err != nil {
	//	panic(err)
	//}

	config := memberlist.DefaultLANConfig()
	config.BindAddr = host
	config.BindPort, _ = strconv.Atoi(port)
	config.SecretKey, _ = base64.StdEncoding.DecodeString(clusterKey)

	ml, err := memberlist.Create(config)
	if err != nil {
		err = fmt.Errorf("initMemberNode: Exception while creating cluster obj: %v", err)
		return err
	}
	mdata := metadata{
		RequestPort:    requestPort,
		SequenceNumber: globalCounter.Add(1),
	}
	ml.LocalNode().Meta = []byte(mdata.toString())

	node = &Node{
		memberlist:  ml,
		host:        host,
		port:        port,
		requestPort: requestPort,
		started:     false,
	}

	if nodeName != peerNodeNames[0] {
		_, err = ml.Join([]string{fmt.Sprintf("%s:%s", masterHost, masterPort)})
		if err != nil {
			err = fmt.Errorf("initMemberNode: Failed to join cluster:: %v", err)
			return err
		}
	} else {
		log.Infof("initMemberNode: new cluster created. key: %s\n", clusterKey)
	}

	go node.ListenNodeLeave()
	go node.ListenToNodesData()
	node.started = true
	return nil
}

func sendRequestToPeers(ctx context.Context, opsType opsType, payload []byte) string {
	requestID := common.GenerateUUID()
	requestMsg := &message{
		ID:              requestID,
		OpsType:         opsType,
		Payload:         payload,
		RequestNodeName: nodeName,
		Sequence:        globalCounter.Load(),
	}
	node.PushDataToOtherNodes(marshallMsg(ctx, requestMsg))
	return requestID
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
		log.Infof("handleConnection(%s): Client %s disconnected", nodeName, conn.RemoteAddr().String())
	}(conn)

	clientReader := bufio.NewReader(conn)

	for {
		responseBuf := make([]byte, 8192)
		readLen := 0
		err := conn.SetDeadline(time.Now().Add(time.Minute * 10))
		readLen, err = clientReader.Read(responseBuf)
		if errors.Is(err, os.ErrDeadlineExceeded) {
			log.Warnf("handleConnection(%s): Client %s is inactive. Maybe its dead. Exiting", nodeName, conn.RemoteAddr().String())
			return // TODO: continue? or return?
		}
		if err != nil {
			log.Errorf("handleConnection(%s): exception while handling incoming message: %v", nodeName, err)
			return
		}
		if parsedMsg, err := unmarshallMsg(ctx, responseBuf[:readLen]); err != nil {
			log.Errorf("handleConnection(%s): exception while unmarshalling incoming msg: %v\n", nodeName, err)
			continue
		} else {
			log.Infof("handleConnection(%s): Received: %v\n", nodeName, parsedMsg.toString())
			if err := handleRequest(ctx, parsedMsg.ID, parsedMsg.OpsType, parsedMsg.Payload); err != nil {
				log.Errorf("handleConnection(%s): Exception while delivering msg: %s. Err: %v\n", nodeName, parsedMsg.toString(), err)
				continue
			}
		}
	}
}

func handleRequest(ctx context.Context, requestID string, opsType opsType, payload []byte) error {
	var sqlRPCServer sqlServerHandlers
	if val, ok := responseTrackers[requestID]; ok {
		val <- true
		delete(responseTrackers, requestID)
		return nil
	}
	switch opsType {
	case CreateBuyer:
		msg := &libProto.CreateBuyerRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateBuyer(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case UpdateBuyerByID:
		msg := &libProto.UpdateBuyerByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.UpdateBuyerByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case CreateCart:
		msg := &libProto.CreateCartRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateCart(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case UpdateCartByID:
		msg := &libProto.UpdateCartByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.UpdateCartByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteCartByID:
		msg := &libProto.DeleteCartByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteCartByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case CreateCartItem:
		msg := &libProto.CreateCartItemRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateCartItem(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case UpdateCartItem:
		msg := &libProto.UpdateCartItemRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.UpdateCartItem(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteCartItemByCartIDAndProductID:
		msg := &libProto.DeleteCartItemByCartIDAndProductIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteCartItemByCartIDAndProductID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteCartItemByCartID:
		msg := &libProto.DeleteCartItemByCartIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteCartItemByCartID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteCartItemByProductID:
		msg := &libProto.DeleteCartItemByProductIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteCartItemByProductID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case CreateSeller:
		msg := &libProto.CreateSellerRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateSeller(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case UpdateSellerByID:
		msg := &libProto.UpdateSellerByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.UpdateSellerByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case CreateSession:
		msg := &libProto.CreateSessionRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateSession(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteSessionByID:
		msg := &libProto.DeleteSessionByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteSessionByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case CreateTransaction:
		msg := &libProto.CreateTransactionRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateTransaction(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteTransactionsByCartID:
		msg := &libProto.DeleteTransactionsByCartIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteTransactionsByCartID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteTransactionsByBuyerID:
		msg := &libProto.DeleteTransactionsByBuyerIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteTransactionsByBuyerID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteTransactionsBySellerID:
		msg := &libProto.DeleteTransactionsBySellerIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteTransactionsBySellerID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", opsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	default:
		return fmt.Errorf("handleRequest: unknown OPSType: %d", opsType)
	}
	return nil
}

func marshallMsg(ctx context.Context, msg *message) []byte {
	data, _ := json.Marshal(*msg)
	return data
}

func unmarshallMsg(ctx context.Context, msgBytes []byte) (*message, error) {
	var msg message
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
