package common

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type ClientResponse struct {
	SessionID  string
	StatusCode int
	Message    string
}

func (req *ClientResponse) SerializeRequest() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(req)
	if err != nil {
		log.Fatal("Encode error:", err)
	}
	return buffer.Bytes()
}

func (req *ClientResponse) DeserializeRequest(data []byte) {
	var buffer bytes.Buffer
	buffer.Write(data)
	decoder := gob.NewDecoder(&buffer)
	if err := decoder.Decode(&req); err != nil {
		log.Fatal("Decode error:", err)
		return
	}
}

func (req *ClientResponse) LogResponse() {
	log.Println("Code: " + strconv.Itoa(req.StatusCode) + " Body: " + req.Message)
}

func RespondWithError(conn net.Conn, code int, message string) {
	response := ClientResponse{
		SessionID:  "",
		StatusCode: code,
		Message:    message,
	}
	conn.Write(response.SerializeRequest())
}

func RespondWithJSON(conn net.Conn, code int, sessionID string, message interface{}) {
	body, _ := json.Marshal(message)
	response := ClientResponse{
		SessionID:  sessionID,
		StatusCode: code,
		Message:    string(body),
	}
	conn.Write(response.SerializeRequest())
}

func RespondWithStatusCode(conn net.Conn, code int, sessionID string) {
	response := ClientResponse{
		SessionID:  sessionID,
		StatusCode: code,
		Message:    "Success",
	}
	conn.Write(response.SerializeRequest())
}

//func RespondWithError(w http.ResponseWriter, code int, message string) {
//	log.Errorf("respondWithError: %v", message)
//	RespondWithJSON(w, code, "", map[string]string{"error": message})
//}
//
//func RespondWithJSON(w http.ResponseWriter, code int, sessionID string, payload interface{}) {
//	response, _ := json.Marshal(payload)
//
//	if len(sessionID) > 0 {
//		w.Header().Set("User-Session-Id", sessionID)
//	}
//
//	w.Header().Set("Access-Control-Expose-Headers", "*")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("Access-Control-Allow-Methods", "*")
//	w.Header().Set("Access-Control-Allow-Headers", "*")
//	w.Header().Set("Access-Control-Allow-Credentials", "true")
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(code)
//	w.Write(response)
//}

//func RespondWithStatusCode(w http.ResponseWriter, code int, headers map[string]string) {
//	if len(headers) > 0 {
//		for k, v := range headers {
//			w.Header().Set(k, v)
//		}
//	}
//	w.Header().Set("Access-Control-Expose-Headers", "*")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("Access-Control-Allow-Methods", "*")
//	w.Header().Set("Access-Control-Allow-Headers", "*")
//	w.Header().Set("Access-Control-Allow-Credentials", "true")
//	w.WriteHeader(code)
//}

func RetryOnError(retryCount int, backOffTime time.Duration, fn func() error) error {
	var err error
	for attempt := 0; attempt < retryCount; attempt++ {
		if err = fn(); err != nil {
			log.Errorf("RetryOnError: attempt number: %d. Error: %v", attempt, err)
		} else {
			log.Infof("RetryOnError: attempt number: %d. successful!", attempt)
			return nil
		}
		time.Sleep(backOffTime)
		continue
	}
	return err
}
