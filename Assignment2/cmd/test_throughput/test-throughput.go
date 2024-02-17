package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/adarshsrinivasan/DS_S24/library/common"
)

const (
	ServiceName   = "test-throughput"
	ServerHostEnv = "SERVER_HOST"
	ServerPortEnv = "SERVER_PORT"
)

var (
	err               error
	httpServerHost    = common.GetEnv(ServerHostEnv, "localhost")
	httpServerPort, _ = strconv.Atoi(common.GetEnv(ServerPortEnv, "50000"))
)

type BuyerModel struct {
	Id                     string    `json:"id,omitempty" bson:"id" bun:"id,pk"`
	Name                   string    `json:"name,omitempty" bson:"name" bun:"name,notnull"`
	UserName               string    `json:"userName,omitempty" bson:"userName" bun:"userName,notnull,unique"`
	Password               string    `json:"password,omitempty" bson:"password" bun:"password,notnull,unique"`
	Version                int       `json:"version,omitempty" bson:"version" bun:"version,notnull"`
	CreatedAt              time.Time `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}

type SellerModel struct {
	Id                 string    `json:"id,omitempty" bson:"id" bun:"id,pk"`
	Name               string    `json:"name,omitempty" bson:"name" bun:"name,notnull"`
	FeedBackThumbsUp   int       `json:"feedBackThumbsUp" bson:"feedBackThumbsUp" bun:"feedBackThumbsUp"`
	FeedBackThumbsDown int       `json:"feedBackThumbsDown" bson:"feedBackThumbsDown" bun:"feedBackThumbsDown"`
	NumberOfItemsSold  int       `json:"numberOfItemsSold,omitempty" bson:"numberOfItemsSold" bun:"numberOfItemsSold"`
	UserName           string    `json:"userName,omitempty" bson:"userName" bun:"userName,notnull,unique"`
	Password           string    `json:"password,omitempty" bson:"password" bun:"password,notnull,unique"`
	Version            int       `json:"version,omitempty" bson:"version" bun:"version,notnull"`
	CreatedAt          time.Time `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}

type SessionModel struct {
	SessionID string          `json:"sessionID,omitempty" bson:"sessionID" bun:"sessionID,pk"`
}

func fireLoginAPI(ctx context.Context, user int)  {
	switch user {
	case 1: {
		request := SellerModel{
			UserName:           "admin",
			Password:           "admin",
		}
		url := fmt.Sprintf("http://%s:%d/api/v1/marketplace/%s/login", httpServerHost, httpServerPort, "seller")
		common.MakeHTTPRequest[SellerModel, SessionModel](ctx, "POST", url, "", request, true)
	}
	default: {
		request := BuyerModel{
			UserName:           "admin",
			Password:           "admin",
		}
		url := fmt.Sprintf("http://%s:%d/api/v1/marketplace/%s/login", httpServerHost, httpServerPort, "buyer")
		common.MakeHTTPRequest[BuyerModel, SessionModel](ctx, "POST", url, "", request, true)
	}
	}
}

func runOperation(ctx context.Context, file *os.File, wg *sync.WaitGroup, thread, user int) {
	// Decrement the counter when the goroutine completes.
	defer wg.Done()
	fmt.Printf("Thread %d: Start\n", thread)
	var (
		iterations int     = 10
		opCount    int     = 1000
		average    float64 = 0
	)

	for i := 0; i < iterations; i++ {
		start := time.Now()
		for i := 0; i < opCount; i++ {
			fireLoginAPI(ctx, user)
		}
		duration := time.Since(start)
		average += (float64(opCount) / duration.Seconds())
	}
	if _, err := file.WriteString(fmt.Sprintf("%f\n", (average / float64(iterations)))); err != nil {
		fmt.Printf("Error while writing to file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Thread %d: End\n", thread)
}

func main() {
	// Check if the user has provided the number of threads
	if len(os.Args) < 3 {
		fmt.Println("Usage:", os.Args[0], "<thread_count> <user: {0: Buyer, 1: Seller}>")
		os.Exit(1)
	}

	threadCount := 10
	threadCount, _ = strconv.Atoi(os.Args[1])

	user := 0
	user, _ = strconv.Atoi(os.Args[2])

	filePath := fmt.Sprintf("throughout-test-%d-%d-%s", user, threadCount, time.Now().Format(time.RFC3339))

	// Write output to the file
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file for writing: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	var wg sync.WaitGroup
	ctx := context.Background()
	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		go runOperation(ctx, f, &wg, i, user)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
}
