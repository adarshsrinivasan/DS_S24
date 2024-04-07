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

type msgType int

const (
	MsgType_None msgType = iota
	MsgType_Request
	MsgType_Sequence
	MsgType_Retransmit
	MsgType_ACK
)

var msgTypeToStr = map[msgType]string{
	MsgType_None:       "None",
	MsgType_Request:    "Request",
	MsgType_Sequence:   "Sequence",
	MsgType_Retransmit: "Retransmit",
	MsgType_ACK:        "ACK",
}

type ackType int

const (
	ACKType_None ackType = iota
	ACKType_Positive
	ACKType_Negavite
)

type message struct {
	ID                 string  `json:"id"`
	MsgType            msgType `json:"msgType"`
	OpsType            opsType `json:"opsType"`
	Payload            []byte  `json:"payload"`
	RequestNodeName    string  `json:"requestNodeName"`
	SequenceNodeName   string  `json:"sequenceNodeName"`
	RetransmitNodeName string  `json:"retransmitNodeName"`
	LocalSeqNum        int32   `json:"localSeqNum"`
	GlobalSeqNum       int32   `json:"globalSeqNum"`
	ACKType            ackType `json:"ackType"`
}

func (m message) toString() string {
	return string(marshallMsg(ctx, &m))
}

var (
	sentRequestMsgs                  = map[string]message{}
	sentSequenceMsgs                 = map[string]message{}
	deliveredSequenceMsgs            = map[string]bool{}
	toBeDeliveredBufferedRequestMsgs = make([]message, 0)
	outOfOrderBufferedRequestMsgs    = map[string]message{}
	outOfOrderBufferedSequenceMsgs   = map[string]message{}
	retransmitTracker                = map[string]message{}
	lastLocalSeqBuffered             = map[string]int32{}
	responseTrackers                 = map[string]chan bool{}
	localCounter, globalCounter      atomic.Int32
)

func getRequestMsgKey(ctx context.Context, requestNodeName string, localSeqNum int32) string {
	return fmt.Sprintf("%s-%d", requestNodeName, localSeqNum)
}

func getSequenceMsgKey(ctx context.Context, msg *message) string {
	return fmt.Sprintf("%d", msg.GlobalSeqNum)
}

func getRetransmitMsgKey(ctx context.Context, msg *message) string {
	key := ""
	switch msg.MsgType {
	case MsgType_Sequence:
		key = fmt.Sprintf("%s-%s-%d", msgTypeToStr[msg.MsgType], msg.ID, msg.GlobalSeqNum)
	case MsgType_Request:
		key = fmt.Sprintf("%s-%s-%s-%d", msgTypeToStr[msg.MsgType], msg.ID, msg.RetransmitNodeName, msg.LocalSeqNum)
	}
	return key
}

func recordRequestSentMsg(ctx context.Context, msg *message) {
	sentRequestMsgs[getRequestMsgKey(ctx, msg.RequestNodeName, msg.LocalSeqNum)] = *msg
}

func recordSequenceSentMsg(ctx context.Context, msg *message) {
	sentSequenceMsgs[getSequenceMsgKey(ctx, msg)] = *msg
}

func addRequestMsgToToBeDeliveredBuffer(ctx context.Context, msg *message) {
	if _, ok := deliveredSequenceMsgs[msg.ID]; ok {
		return
	}
	toBeDeliveredBufferedRequestMsgs = append(toBeDeliveredBufferedRequestMsgs, *msg)
}

func addRequestMsgToOutOfOrderBuffer(ctx context.Context, msg *message) {
	outOfOrderBufferedRequestMsgs[getRequestMsgKey(ctx, msg.RequestNodeName, msg.LocalSeqNum)] = *msg
}

func addSequenceMsgToOutOfOrderBuffer(ctx context.Context, msg *message) {
	outOfOrderBufferedSequenceMsgs[getSequenceMsgKey(ctx, msg)] = *msg
}

func addMsgToRetransmitTracker(ctx context.Context, msg *message) {
	retransmitTracker[getRetransmitMsgKey(ctx, msg)] = *msg
}

func addMsgToDeliveredSequenceMsgs(ctx context.Context, msg *message) {
	deliveredSequenceMsgs[msg.ID] = true
}

func removeMsgFromRetransmitTracker(ctx context.Context, msg *message) {
	key := getRetransmitMsgKey(ctx, msg)
	if _, ok := retransmitTracker[key]; ok {
		delete(retransmitTracker, key)
	}
}

func removeRequestMsgFromToBeDeliveredBuffered(ctx context.Context, msg *message) {
	for i := 0; i < len(toBeDeliveredBufferedRequestMsgs); i++ {
		if toBeDeliveredBufferedRequestMsgs[i].ID != msg.ID {
			toBeDeliveredBufferedRequestMsgs = append(toBeDeliveredBufferedRequestMsgs[:i], toBeDeliveredBufferedRequestMsgs[i+1:]...)
			break
		}
	}
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

func broadcastMsgToPeers(ctx context.Context, msg *message, sendLastNodeName string) {
	sendLastNodePort := ""
	for i := 0; i < len(peerNodeNames); i++ {
		if peerNodeNames[i] == sendLastNodeName {
			sendLastNodePort = peerNodePorts[i]
			continue
		}
		addr := net.JoinHostPort(peerNodeNames[i], peerNodePorts[i])
		if raddr, err := net.ResolveUDPAddr("udp", addr); err != nil {
			log.Errorf("broadcastMsgToPeers(%s): Exception while resolving addr %s. %v\n", nodeName, addr, err)
			continue
		} else {
			if conn, err := net.DialUDP("udp", nil, raddr); err != nil {
				log.Errorf("broadcastMsgToPeers(%s): Exception while dailing addr %s. %v\n", nodeName, addr, err)
				continue
			} else {
				log.Infof("broadcastMsgToPeers(%s): Sending the following msg to %s: %s\n", nodeName, addr, msg.toString())
				if _, err = conn.Write(marshallMsg(ctx, msg)); err != nil {
					log.Errorf("broadcastMsgToPeers(%s): Exception while writing msg to addr %s. %v\n", nodeName, addr, err)
					continue
				}
			}
		}
	}
	sendMsgToNode(ctx, sendLastNodeName, sendLastNodePort, msg)
	return
}

func sendMsgToNode(ctx context.Context, receiverNodeName, receiverNodePort string, msg *message) {
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
			log.Infof("listenFromPeers(%s): Received: %v\n", nodeName, parsedMsg.toString())
			handleReceivedMsg(ctx, parsedMsg)
		}
	}
}

func sendSequenceRetransmitToPeers(ctx context.Context, msg *message, from, to int32) {
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
		broadcastMsgToPeers(ctx, msg, "")
	}
	return
}

func sendRequestRetransmitToNode(ctx context.Context, msg *message, from, to int32) {
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

func sendRequestToPeers(ctx context.Context, opsType opsType, payload []byte) (string, <-chan bool) {
	requestID := common.GenerateUUID()
	responseChan := make(chan bool)

	requestMsg := &message{
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
	nextLeader := ((globalCounter.Load() + 1) % int32(len(peerNodeNames)))
	broadcastMsgToPeers(ctx, requestMsg, fmt.Sprintf("%s%d", CustomerDBNodeNameBase, nextLeader))
	return requestID, responseChan
}

func checkTurnAndSendSequenceToPeers(ctx context.Context) {
	nextLeader := ((globalCounter.Load() + 1) % int32(len(peerNodeNames)))
	if nextLeader == 0 {
		nextLeader = int32(len(peerNodeNames))
	}
	for len(toBeDeliveredBufferedRequestMsgs) > 0 {
		if _, ok := deliveredSequenceMsgs[toBeDeliveredBufferedRequestMsgs[0].ID]; ok {
			toBeDeliveredBufferedRequestMsgs = toBeDeliveredBufferedRequestMsgs[1:]
		} else {
			break
		}
	}
	if nodeName == fmt.Sprintf("%s%d", CustomerDBNodeNameBase, nextLeader) && len(toBeDeliveredBufferedRequestMsgs) > 0 {
		log.Infof("checkTurnAndSendSequenceToPeers(%s): Taking responsibility to send next sequence message.\n", nodeName)
		for len(retransmitTracker) > 0 {
			log.Infof("checkTurnAndSendSequenceToPeers(%s): Cannot deliver since I have %d retransmit requests pending. Sleeping for 2sec before recheck.\n", nodeName, len(retransmitTracker))
			time.Sleep(2 * time.Second)
		}

		nextMsg := toBeDeliveredBufferedRequestMsgs[0]
		nextMsg.MsgType = MsgType_Sequence
		nextMsg.SequenceNodeName = nodeName
		nextMsg.GlobalSeqNum = (globalCounter.Load() + 1)
		recordSequenceSentMsg(ctx, &nextMsg)
		broadcastMsgToPeers(ctx, &nextMsg, "")
	} else {
		log.Infof("checkTurnAndSendSequenceToPeers(%s): Not my responsibility to send next sequence message. Responsibility of: %s", nodeName, fmt.Sprintf("%s%d", CustomerDBNodeNameBase, nextLeader))
	}
}

func handleReceivedMsg(ctx context.Context, msg *message) {
	removeMsgFromRetransmitTracker(ctx, msg)
	switch msg.MsgType {
	case MsgType_Sequence:
		{
			if (globalCounter.Load() + 1) == msg.GlobalSeqNum {
				removeRequestMsgFromToBeDeliveredBuffered(ctx, msg)
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

func deliverSequenceMsg(ctx context.Context, msg *message) {
	addMsgToDeliveredSequenceMsgs(ctx, msg)
	//TODO: Add retry if handleRequest fails???
	if err := handleRequest(ctx, msg.ID, msg.OpsType, msg.Payload); err != nil {
		log.Errorf("deliverSequenceMsg(%s): Exception while delivering Sequence msg: SeqNo.: %d, opsType: %s, Err: %v\n", nodeName, msg.GlobalSeqNum, opsTypeToStr[msg.OpsType], err)
		return
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
