package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/library/common"
	"net"
	"sync/atomic"
	"time"

	libProto "github.com/adarshsrinivasan/DS_S24/library/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type OpsType int

const (
	CreateBuyer OpsType = iota
	GetBuyerByID
	GetBuyerByUserName
	UpdateBuyerByID
	CreateCart
	GetCartByID
	GetCartByBuyerID
	UpdateCartByID
	DeleteCartByID
	CreateCartItem
	GetCartItemByID
	GetCartItemByCartIDAndProductID
	ListCartItemByCartID
	UpdateCartItem
	DeleteCartItemByCartIDAndProductID
	DeleteCartItemByCartID
	DeleteCartItemByProductID
	CreateSeller
	GetSellerByID
	GetSellerByUserName
	UpdateSellerByID
	CreateSession
	GetSessionByID
	GetSessionByUserID
	DeleteSessionByID
	CreateTransaction
	ListTransactionsBySellerID
	ListTransactionsByBuyerID
	ListTransactionsByCartID
	DeleteTransactionsByCartID
	DeleteTransactionsByBuyerID
	DeleteTransactionsBySellerID
)

var OpsTypeToStr = map[OpsType]string{
	CreateBuyer:                        "CreateBuyer",
	GetBuyerByID:                       "GetBuyerByID",
	GetBuyerByUserName:                 "GetBuyerByUserName",
	UpdateBuyerByID:                    "UpdateBuyerByID",
	CreateCart:                         "CreateCart",
	GetCartByID:                        "GetCartByID",
	GetCartByBuyerID:                   "GetCartByBuyerID",
	UpdateCartByID:                     "UpdateCartByID",
	DeleteCartByID:                     "DeleteCartByID",
	CreateCartItem:                     "CreateCartItem",
	GetCartItemByID:                    "GetCartItemByID",
	GetCartItemByCartIDAndProductID:    "GetCartItemByCartIDAndProductID",
	ListCartItemByCartID:               "ListCartItemByCartID",
	UpdateCartItem:                     "UpdateCartItem",
	DeleteCartItemByCartIDAndProductID: "DeleteCartItemByCartIDAndProductID",
	DeleteCartItemByCartID:             "DeleteCartItemByCartID",
	DeleteCartItemByProductID:          "DeleteCartItemByProductID",
	CreateSeller:                       "CreateSeller",
	GetSellerByID:                      "GetSellerByID",
	GetSellerByUserName:                "GetSellerByUserName",
	UpdateSellerByID:                   "UpdateSellerByID",
	CreateSession:                      "CreateSession",
	GetSessionByID:                     "GetSessionByID",
	GetSessionByUserID:                 "GetSessionByUserID",
	DeleteSessionByID:                  "DeleteSessionByID",
	CreateTransaction:                  "CreateTransaction",
	ListTransactionsBySellerID:         "ListTransactionsBySellerID",
	ListTransactionsByBuyerID:          "ListTransactionsByBuyerID",
	ListTransactionsByCartID:           "ListTransactionsByCartID",
	DeleteTransactionsByCartID:         "DeleteTransactionsByCartID",
	DeleteTransactionsByBuyerID:        "DeleteTransactionsByBuyerID",
	DeleteTransactionsBySellerID:       "DeleteTransactionsBySellerID",
}

type MsgType int

const (
	MsgType_None MsgType = iota
	MsgType_Request
	MsgType_Sequence
	MsgType_Retransmit
	MsgType_ACK
)

var MsgTypeToStr = map[MsgType]string{
	MsgType_None:       "None",
	MsgType_Request:    "Request",
	MsgType_Sequence:   "Sequence",
	MsgType_Retransmit: "Retransmit",
	MsgType_ACK:        "ACK",
}

type ACKType int

const (
	ACKType_None ACKType = iota
	ACKType_Positive
	ACKType_Negavite
)

type Message struct {
	ID                 string  `json:"id"`
	MsgType            MsgType `json:"msgType"`
	OpsType            OpsType `json:"opsType"`
	Payload            []byte  `json:"payload"`
	RequestNodeName    string  `json:"requestNodeName"`
	SequenceNodeName   string  `json:"sequenceNodeName"`
	RetransmitNodeName string  `json:"retransmitNodeName"`
	LocalSeqNum        int32   `json:"localSeqNum"`
	GlobalSeqNum       int32   `json:"globalSeqNum"`
	ACKType            ACKType `json:"ackType"`
}

func (m Message) ToString() string {
	return string(marshallMsg(ctx, &m))
}

var (
	sentRequestMsgs                  = map[string]Message{}
	sentSequenceMsgs                 = map[string]Message{}
	toBeDeliveredBufferedRequestMsgs = make([]Message, 0)
	outOfOrderBufferedRequestMsgs    = map[string]Message{}
	outOfOrderBufferedSequenceMsgs   = map[string]Message{}
	retransmitTracker                = map[string]Message{}
	lastLocalSeqBuffered             = map[string]int32{}
	responseTrackers                 = map[string]chan bool{}
	localCounter, globalCounter      atomic.Int32
)

func getRequestMsgKey(ctx context.Context, requestNodeName string, localSeqNum int32) string {
	return fmt.Sprintf("%s-%d", requestNodeName, localSeqNum)
}

func getSequenceMsgKey(ctx context.Context, msg *Message) string {
	return fmt.Sprintf("%d", msg.GlobalSeqNum)
}

func getRetransmitMsgKey(ctx context.Context, msg *Message) string {
	key := ""
	switch msg.MsgType {
	case MsgType_Sequence:
		key = fmt.Sprintf("%s-%s-%d", MsgTypeToStr[msg.MsgType], msg.ID, msg.GlobalSeqNum)
	case MsgType_Request:
		key = fmt.Sprintf("%s-%s-%s-%d", MsgTypeToStr[msg.MsgType], msg.ID, msg.RetransmitNodeName, msg.LocalSeqNum)
	}
	return key
}

func recordRequestSentMsg(ctx context.Context, msg *Message) {
	sentRequestMsgs[getRequestMsgKey(ctx, msg.RequestNodeName, msg.LocalSeqNum)] = *msg
}

func recordSequenceSentMsg(ctx context.Context, msg *Message) {
	sentSequenceMsgs[getSequenceMsgKey(ctx, msg)] = *msg
}

func addRequestMsgToToBeDeliveredBuffer(ctx context.Context, msg *Message) {
	toBeDeliveredBufferedRequestMsgs = append(toBeDeliveredBufferedRequestMsgs, *msg)
}

func addRequestMsgToOutOfOrderBuffer(ctx context.Context, msg *Message) {
	outOfOrderBufferedRequestMsgs[getRequestMsgKey(ctx, msg.RequestNodeName, msg.LocalSeqNum)] = *msg
}

func addSequenceMsgToOutOfOrderBuffer(ctx context.Context, msg *Message) {
	outOfOrderBufferedSequenceMsgs[getSequenceMsgKey(ctx, msg)] = *msg
}

func addMsgToRetransmitTracker(ctx context.Context, msg *Message) {
	retransmitTracker[getRetransmitMsgKey(ctx, msg)] = *msg
}

func removeMsgFromRetransmitTracker(ctx context.Context, msg *Message) {
	key := getRetransmitMsgKey(ctx, msg)
	if _, ok := retransmitTracker[key]; ok {
		delete(retransmitTracker, key)
	}
}

func removeRequestMsgFromToBeDeliveredBuffered(ctx context.Context, msg *Message) {
	i := 0
	for ; toBeDeliveredBufferedRequestMsgs[i].ID != msg.ID; i++ {
	}
	toBeDeliveredBufferedRequestMsgs = append(toBeDeliveredBufferedRequestMsgs[:i], toBeDeliveredBufferedRequestMsgs[i+1:]...)
}

func marshallMsg(ctx context.Context, msg *Message) []byte {
	data, _ := json.Marshal(*msg)
	return data
}

func unmarshallMsg(ctx context.Context, msgBytes []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func broadcastMsgToPeers(ctx context.Context, msg *Message) {
	for i := 0; i < len(peerNodeNames); i++ {
		addr := net.JoinHostPort(peerNodeNames[i], peerNodePorts[i])
		if raddr, err := net.ResolveUDPAddr("udp", addr); err != nil {
			log.Errorf("broadcastMsgToPeers(%s): Exception while resolving addr %s. %v\n", nodeName, addr, err)
			continue
		} else {
			if conn, err := net.DialUDP("udp", nil, raddr); err != nil {
				log.Errorf("broadcastMsgToPeers(%s): Exception while dailing addr %s. %v\n", nodeName, addr, err)
				continue
			} else {
				log.Infof("broadcastMsgToPeers(%s): Sending the following msg to %s: %s\n", nodeName, addr, msg.ToString())
				if _, err = conn.Write(marshallMsg(ctx, msg)); err != nil {
					log.Errorf("broadcastMsgToPeers(%s): Exception while writing msg to addr %s. %v\n", nodeName, addr, err)
					continue
				}
			}
		}
	}
	return
}

func sendMsgToNode(ctx context.Context, receiverNodeName, receiverNodePort string, msg *Message) {
	if receiverNodePort == "" {
		for i := 0; i < len(peerNodeNames); i++ {
			if peerNodeNames[i] == receiverNodeName {
				receiverNodePort = peerNodePorts[i]
				break
			}
		}
	}

	addr := net.JoinHostPort(receiverNodeName, receiverNodePort)
	if raddr, err := net.ResolveUDPAddr("udp", addr); err != nil {
		log.Errorf("sendMsgToNode(%s): Exception while resolving addr %s. %v\n", nodeName, addr, err)
		return
	} else {
		if conn, err := net.DialUDP("udp", nil, raddr); err != nil {
			log.Errorf("sendMsgToNode(%s): Exception while dailing addr %s. %v\n", nodeName, addr, err)
			return
		} else {
			if _, err = conn.Write(marshallMsg(ctx, msg)); err != nil {
				log.Errorf("sendMsgToNode(%s): Exception while writing msg to addr %s. %v\n", nodeName, addr, err)
				return
			}
		}
	}
	return
}

func listenFromPeers(ctx context.Context) {
	addr := net.JoinHostPort(syncHost, fmt.Sprintf("%d", syncPort))
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Errorf("listenFromPeers(%s): Exception while resolving addr %s. %v\n", nodeName, addr, err)
		return
	}
	log.Infof("listenFromPeers(%s): listening on addr %s...\n", nodeName, addr)

	conn, err := net.ListenUDP("udp", raddr)
	if err != nil {
		log.Fatalf("%v: ERROR: Server listening failed. %v", ServiceName, err)
	}

	for {
		responseBuf := make([]byte, 8192)
		readLen := 0

		if readLen, _, err = conn.ReadFrom(responseBuf); err != nil {
			log.Errorf("listenFromPeers(%s): exception while reading incoming msg on addr %s. %v\n", nodeName, addr, err)
			continue
		}
		if parsedMsg, err := unmarshallMsg(ctx, responseBuf[:readLen]); err != nil {
			log.Errorf("listenFromPeers(%s): exception while unmarshalling incoming msg on addr %s. %v\n", nodeName, addr, err)
		} else {
			log.Infof("listenFromPeers(%s): Received: %v\n", nodeName, parsedMsg.ToString())
			handleReceivedMsg(ctx, parsedMsg)
		}
	}
}

func sendSequenceRetransmitToPeers(ctx context.Context, msg *Message, from, to int32) {
	if to < from {
		return
	}
	msg.LocalSeqNum = -1
	msg.RetransmitNodeName = nodeName
	for i := from; i <= to; i++ {
		msg.MsgType = MsgType_Sequence
		msg.GlobalSeqNum = i
		addMsgToRetransmitTracker(ctx, msg)
		msg.MsgType = MsgType_Retransmit
		broadcastMsgToPeers(ctx, msg)
	}
	return
}

func sendRequestRetransmitToNode(ctx context.Context, msg *Message, from, to int32) {
	if to < from {
		return
	}
	msg.GlobalSeqNum = -1
	msg.RetransmitNodeName = nodeName
	for i := from; i <= to; i++ {
		msg.MsgType = MsgType_Request
		msg.LocalSeqNum = i
		addMsgToRetransmitTracker(ctx, msg)
		msg.MsgType = MsgType_Retransmit
		sendMsgToNode(ctx, msg.RequestNodeName, "", msg)
	}
	return
}

func sendRequestToPeers(ctx context.Context, opsType OpsType, payload []byte) (string, <-chan bool) {
	requestID := common.GenerateUUID()
	responseChan := make(chan bool)

	requestMsg := &Message{
		ID:              requestID,
		MsgType:         MsgType_Request,
		OpsType:         opsType,
		Payload:         payload,
		RequestNodeName: nodeName,
		LocalSeqNum:     localCounter.Add(1),
		GlobalSeqNum:    -1,
	}
	responseTrackers[requestID] = responseChan
	recordRequestSentMsg(ctx, requestMsg)
	broadcastMsgToPeers(ctx, requestMsg)
	return requestID, responseChan
}

func checkTurnAndSendSequenceToPeers(ctx context.Context) {
	nextLeader := ((globalCounter.Load() + 1) % int32(len(peerNodeNames)))
	if nodeName == fmt.Sprintf("%s%d", NodeNameBase, nextLeader) && len(toBeDeliveredBufferedRequestMsgs) > 0 {
		log.Infof("checkTurnAndSendSequenceToPeers(%s): Taking responsibility to send next sequence message.\n", nodeName)
		for len(retransmitTracker) > 0 {
			log.Infof("checkTurnAndSendSequenceToPeers(%s): Cannot deliver since I have %d retransmit requests pending. Sleeping for 2sec before recheck.\n", nodeName, len(retransmitTracker))
			time.Sleep(2 * time.Second)
		}

		nextMsg := toBeDeliveredBufferedRequestMsgs[0]
		nextMsg.MsgType = MsgType_Sequence
		nextMsg.SequenceNodeName = nodeName
		nextMsg.GlobalSeqNum = globalCounter.Add(1)
		recordSequenceSentMsg(ctx, &nextMsg)
		broadcastMsgToPeers(ctx, &nextMsg)
	} else {
		log.Infof("checkTurnAndSendSequenceToPeers(%s): Not my responsibility to send next sequence message. Responsibility of: %s\n", nodeName, fmt.Sprintf("%s%d", NodeNameBase, nextLeader))
	}
}

func handleReceivedMsg(ctx context.Context, msg *Message) {
	removeMsgFromRetransmitTracker(ctx, msg)
	switch msg.MsgType {
	case MsgType_Sequence:
		{
			if ((globalCounter.Load() + 1) == msg.GlobalSeqNum) || (globalCounter.Load() == msg.GlobalSeqNum && msg.SequenceNodeName == nodeName) {
				globalCounter.Add(1)
				deliverSequenceMsg(ctx, msg)
				for bufferedSeqMsg, ok := outOfOrderBufferedSequenceMsgs[fmt.Sprintf("%d", (globalCounter.Load()+1))]; ok; {
					deliverSequenceMsg(ctx, &bufferedSeqMsg)
					globalCounter.Add(1)
					delete(outOfOrderBufferedSequenceMsgs, fmt.Sprintf("%d", globalCounter.Load()))
				}
			} else if (globalCounter.Load() + 1) < msg.GlobalSeqNum {
				addSequenceMsgToOutOfOrderBuffer(ctx, msg)
				sendSequenceRetransmitToPeers(ctx, msg, (globalCounter.Load() + 1), (msg.GlobalSeqNum - 1))
			} else {
				log.Infof("handleReceivedMsg(%s): Received old sequence msg: %d. globalCounter: %d\n", nodeName, msg.GlobalSeqNum, globalCounter)
			}
		}
	case MsgType_Request:
		{
			if (lastLocalSeqBuffered[msg.RequestNodeName] + 1) == msg.LocalSeqNum {
				addRequestMsgToToBeDeliveredBuffer(ctx, msg)
				lastLocalSeqBuffered[msg.RequestNodeName]++
				for bufferedReqMsg, ok := outOfOrderBufferedRequestMsgs[getRequestMsgKey(ctx, msg.RequestNodeName, (lastLocalSeqBuffered[msg.RequestNodeName]+1))]; ok; {
					lastLocalSeqBuffered[msg.RequestNodeName]++
					delete(outOfOrderBufferedRequestMsgs, getRequestMsgKey(ctx, msg.RequestNodeName, lastLocalSeqBuffered[msg.RequestNodeName]))
					addRequestMsgToToBeDeliveredBuffer(ctx, &bufferedReqMsg)
				}
				checkTurnAndSendSequenceToPeers(ctx)
			} else if (lastLocalSeqBuffered[msg.RequestNodeName] + 1) < msg.LocalSeqNum {
				addRequestMsgToOutOfOrderBuffer(ctx, msg)
				sendRequestRetransmitToNode(ctx, msg, (lastLocalSeqBuffered[msg.RequestNodeName] + 1), (msg.LocalSeqNum - 1))
			} else {
				log.Infof("handleReceivedMsg(%s): Received old request msg: %s. globalCounter: %d\n", nodeName, getRequestMsgKey(ctx, msg.RequestNodeName, msg.LocalSeqNum), globalCounter)
			}
		}
	case MsgType_Retransmit:
		{
			if msg.GlobalSeqNum != -1 && msg.SequenceNodeName == nodeName {
				if sentSeqMsg, ok := sentSequenceMsgs[getSequenceMsgKey(ctx, msg)]; ok {
					log.Infof("handleReceivedMsg(%s): Retransmitting Sequence msg: %d\n", nodeName, msg.GlobalSeqNum)
					sendMsgToNode(ctx, msg.RetransmitNodeName, "", &sentSeqMsg)
				}
			} else if msg.LocalSeqNum != -1 && msg.RequestNodeName == nodeName {
				if sentReqMsg, ok := sentRequestMsgs[getRequestMsgKey(ctx, msg.RequestNodeName, msg.LocalSeqNum)]; ok {
					log.Infof("handleReceivedMsg(%s): Retransmitting Request msg: %s-%d\n", nodeName, nodeName, msg.LocalSeqNum)
					sendMsgToNode(ctx, msg.RetransmitNodeName, "", &sentReqMsg)
				}
			}
		}
	default:
		{
			log.Errorf("handleReceivedMsg(%s): Invalid msg type: %d\n", nodeName, msg.MsgType)
		}
	}
	return
}

func deliverSequenceMsg(ctx context.Context, msg *Message) {
	removeRequestMsgFromToBeDeliveredBuffered(ctx, msg)
	//TODO: Add retry if handleRequest fails???
	if err := handleRequest(ctx, msg.ID, msg.OpsType, msg.Payload); err != nil {
		log.Errorf("deliverSequenceMsg(%s): Exception while delivering Sequence msg: SeqNo.: %d, OpsType: %s, Err: %v\n", nodeName, msg.GlobalSeqNum, OpsTypeToStr[msg.OpsType], err)
		return
	}
}

func handleRequest(ctx context.Context, requestID string, opsType OpsType, payload []byte) error {
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
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateBuyer(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case GetBuyerByID:
		msg := &libProto.GetBuyerByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.GetBuyerByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case GetBuyerByUserName:
		msg := &libProto.GetBuyerByUserNameRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.GetBuyerByUserName(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case UpdateBuyerByID:
		msg := &libProto.UpdateBuyerByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.UpdateBuyerByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case CreateCart:
		msg := &libProto.CreateCartRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateCart(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case GetCartByID:
		msg := &libProto.GetCartByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.GetCartByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case GetCartByBuyerID:
		msg := &libProto.GetCartByBuyerIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.GetCartByBuyerID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case UpdateCartByID:
		msg := &libProto.UpdateCartByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.UpdateCartByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteCartByID:
		msg := &libProto.DeleteCartByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteCartByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case CreateCartItem:
		msg := &libProto.CreateCartItemRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateCartItem(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case GetCartItemByID:
		msg := &libProto.GetCartItemByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.GetCartItemByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case GetCartItemByCartIDAndProductID:
		msg := &libProto.GetCartItemByCartIDAndProductIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.GetCartItemByCartIDAndProductID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case ListCartItemByCartID:
		msg := &libProto.ListCartItemByCartIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.ListCartItemByCartID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case UpdateCartItem:
		msg := &libProto.UpdateCartItemRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.UpdateCartItem(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteCartItemByCartIDAndProductID:
		msg := &libProto.DeleteCartItemByCartIDAndProductIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteCartItemByCartIDAndProductID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteCartItemByCartID:
		msg := &libProto.DeleteCartItemByCartIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteCartItemByCartID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteCartItemByProductID:
		msg := &libProto.DeleteCartItemByProductIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteCartItemByProductID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case CreateSeller:
		msg := &libProto.CreateSellerRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateSeller(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case GetSellerByID:
		msg := &libProto.GetSellerByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.GetSellerByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case GetSellerByUserName:
		msg := &libProto.GetSellerByUserNameRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.GetSellerByUserName(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case UpdateSellerByID:
		msg := &libProto.UpdateSellerByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.UpdateSellerByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case CreateSession:
		msg := &libProto.CreateSessionRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateSession(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case GetSessionByID:
		msg := &libProto.GetSessionByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.GetSessionByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case GetSessionByUserID:
		msg := &libProto.GetSessionByUserIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.GetSessionByUserID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteSessionByID:
		msg := &libProto.DeleteSessionByIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteSessionByID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case CreateTransaction:
		msg := &libProto.CreateTransactionRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.CreateTransaction(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case ListTransactionsBySellerID:
		msg := &libProto.ListTransactionsBySellerIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.ListTransactionsBySellerID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case ListTransactionsByBuyerID:
		msg := &libProto.ListTransactionsByBuyerIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.ListTransactionsByBuyerID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case ListTransactionsByCartID:
		msg := &libProto.ListTransactionsByCartIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.ListTransactionsByCartID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteTransactionsByCartID:
		msg := &libProto.DeleteTransactionsByCartIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteTransactionsByCartID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteTransactionsByBuyerID:
		msg := &libProto.DeleteTransactionsByBuyerIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteTransactionsByBuyerID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	case DeleteTransactionsBySellerID:
		msg := &libProto.DeleteTransactionsBySellerIDRequest{}
		if err := proto.Unmarshal(payload, msg); err != nil {
			err = fmt.Errorf("exception while Unmarshalling %s Msg: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
		if _, err := sqlRPCServer.DeleteTransactionsBySellerID(ctx, msg); err != nil {
			err = fmt.Errorf("exception while invoking %s operation: %v", OpsTypeToStr[opsType], err)
			log.Errorf("handleRequest: %v\n", err)
			return err
		}
	default:
		return fmt.Errorf("handleRequest: unknown OPSType: %d", opsType)
	}
	return nil
}
