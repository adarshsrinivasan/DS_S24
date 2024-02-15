package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type ResponseType struct {
	StatusCode int `json:"statusCode"`
	Response string `json:"response,omitempty"`
}

func MakeHTTPRequest[T any, U any](ctx context.Context, method, url, sessionID string, data T, printResponse bool) (U, error) {
	var resp U
	switch method {
	case "POST": return post[T, U](ctx, url, sessionID, data, printResponse)
	case "PUT": return put[T, U](ctx, url, sessionID, data, printResponse)
	case "GET": return get[U](ctx, url, sessionID, printResponse)
	default:
		return resp, fmt.Errorf("unknown method: %s. supported methods:[POST, PUT, GET]", method)
	}
}

func get[U any](ctx context.Context, url, sessionID string, printResponse bool) (U, error) {
	var m U
	r, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return m, err
	}

	if len(sessionID) > 0 {
		r.Header.Add("User-Session-Id", sessionID)
	}

	r.Header.Add("Access-Control-Expose-Headers", "*")
	r.Header.Add("Access-Control-Allow-Origin", "*")
	r.Header.Add("Access-Control-Allow-Methods", "*")
	r.Header.Add("Access-Control-Allow-Headers", "*")
	r.Header.Add("Access-Control-Allow-Credentials", "true")

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return m, err
	}
	statusCode := res.StatusCode
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return m, err
	}
	return parseJSON[U](body, statusCode, printResponse)
}

func put[T any, U any](ctx context.Context, url, sessionID string, data T, printResponse bool) (U, error) {
	var m U
	b, err := toJSON(data)
	if err != nil {
		return m, err
	}
	byteReader := bytes.NewReader(b)
	r, err := http.NewRequestWithContext(ctx, "PUT", url, byteReader)
	if err != nil {
		return m, err
	}
	// Important to set
	if len(sessionID) > 0 {
		r.Header.Add("User-Session-Id", sessionID)
	}

	r.Header.Add("Access-Control-Expose-Headers", "*")
	r.Header.Add("Access-Control-Allow-Origin", "*")
	r.Header.Add("Access-Control-Allow-Methods", "*")
	r.Header.Add("Access-Control-Allow-Headers", "*")
	r.Header.Add("Access-Control-Allow-Credentials", "true")
	r.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return m, err
	}
	statusCode := res.StatusCode
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return m, err
	}
	return parseJSON[U](body, statusCode, printResponse)
}

func post[T any, U any](ctx context.Context, url, sessionID string, data T, printResponse bool) (U, error) {
	var m U
	b, err := toJSON(data)
	if err != nil {
		return m, err
	}
	byteReader := bytes.NewReader(b)
	r, err := http.NewRequestWithContext(ctx, "POST", url, byteReader)
	if err != nil {
		return m, err
	}
	// Important to set
	if len(sessionID) > 0 {
		r.Header.Add("User-Session-Id", sessionID)
	}

	r.Header.Add("Access-Control-Expose-Headers", "*")
	r.Header.Add("Access-Control-Allow-Origin", "*")
	r.Header.Add("Access-Control-Allow-Methods", "*")
	r.Header.Add("Access-Control-Allow-Headers", "*")
	r.Header.Add("Access-Control-Allow-Credentials", "true")
	r.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return m, err
	}
	statusCode := res.StatusCode
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return m, err
	}
	return parseJSON[U](body, statusCode, printResponse)
}

func parseJSON[U any](s []byte, statusCode int, print bool) (U, error) {
	var r U
	if len(s) > 0 && reflect.TypeOf(r) != reflect.TypeOf("")  {
		if err := json.Unmarshal(s, &r); err != nil {
			return r, err
		}
	}

	if print {
		if len(s) > 0 && reflect.TypeOf(r) != reflect.TypeOf("") {
			b, err := json.MarshalIndent(r, "", "  ")
			if err != nil {
				return r, err
			}
			fmt.Printf("StatusCode: %d\nResponse:%s\n", statusCode, string(b))
		} else {
			fmt.Printf("StatusCode: %d\nResponse: %s \n", statusCode, string(s))
		}

	}
	return r, nil
}

func toJSON(T any) ([]byte, error) {
	return json.Marshal(T)
}
