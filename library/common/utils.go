package common

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"os"
	"strconv"
	"strings"

	myproto "github.com/adarshsrinivasan/DS_S24/library/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ReadTrimString(reader *bufio.Reader) (string, error) {
	str, err := reader.ReadString('\n')
	return strings.Split(strings.TrimSpace(str), "\n")[0], err
}

func ConvertInterfaceToAny(v interface{}) (*any.Any, error) {
	anyValue := &any.Any{}
	bytes, _ := json.Marshal(v)
	bytesValue := &wrappers.BytesValue{
		Value: bytes,
	}
	err := anypb.MarshalFrom(anyValue, bytesValue, proto.MarshalOptions{})
	return anyValue, err
}

func ConvertAnyToInterface(anyValue *any.Any) (interface{}, error) {
	var value interface{}
	bytesValue := &wrappers.BytesValue{}
	err := anypb.UnmarshalTo(anyValue, bytesValue, proto.UnmarshalOptions{})
	if err != nil {
		return value, err
	}
	uErr := json.Unmarshal(bytesValue.Value, &value)
	if uErr != nil {
		return value, uErr
	}
	return value, nil
}

func ConvertErrorToProtoError(err error) *myproto.Error {
	if err == nil {
		return nil
	}
	return &myproto.Error{
		Message: err.Error(),
	}
}

func NewNOSQLRPCClient(ctx context.Context, nosqlRPCHost string, nosqlRPCPort int) (myproto.NOSQLServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", nosqlRPCHost, nosqlRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		err = fmt.Errorf("exception while connecting to NOSQLDB RPC server. %v", err)
		logrus.Errorf("NewNOSQLRPCClient: %v\n", err)
		return nil, nil, err
	}
	client := myproto.NewNOSQLServiceClient(conn)
	return client, conn, err
}

func NewSQLRPCClient(ctx context.Context, sqlRPCHost string, sqlRPCPort int) (myproto.SQLServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", sqlRPCHost, sqlRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		err = fmt.Errorf("exception while connecting to SQLDB RPC server. %v", err)
		logrus.Errorf("NewSQLRPCClient: %v\n", err)
		return nil, nil, err
	}
	client := myproto.NewSQLServiceClient(conn)
	return client, conn, err
}

func ReturnTrueWithProbability(probability int) bool {
	if probability < 0 {
		probability *= -1
	}
	if probability != 100 {
		probability %= 100
	}
	return rand.Intn(100) <= (probability - 1)
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func SplitNodeNames(nodeNames string) []string {
	return strings.Split(nodeNames, ",")
}

func GetRandomHostAndPort(nodeNamesList, ports []string) (string, int) {
	ind := rand.Intn(len(nodeNamesList))
	port, _ := strconv.Atoi(ports[ind])
	return nodeNamesList[ind], port
}
