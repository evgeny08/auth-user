package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Service.CreateUser encoders/decoders.
func encodeCreateUserRequest(_ context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = "/api/v1/user"
	req := request.(createUserRequest)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req.User); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createUserRequest
	err := json.NewDecoder(r.Body).Decode(&req.User)
	return req, err
}

func encodeCreateUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(createUserResponse)
	if res.Err != nil {
		return encodeError(w, res.Err, true)
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func decodeCreateUserResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return createUserResponse{Err: decodeError(r)}, nil
	}
	res := createUserResponse{Err: nil}
	return res, nil
}

// Service.AuthUser encoders/decoders.
func encodeAuthUserRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(authUserRequest)
	r.URL.Path = "/api/v1/login/" + url.QueryEscape(req.Login) + "/password/" + url.QueryEscape(req.Password)
	return nil
}

func decodeAuthUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	login := mux.Vars(r)["login"]
	password := mux.Vars(r)["password"]
	return authUserRequest{Login: login, Password: password}, nil
}

func encodeAuthUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(authUserResponse)
	if res.Err != nil {
		return encodeError(w, res.Err, true)
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Session)
}

func decodeAuthUserResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return authUserResponse{Err: decodeError(r)}, nil
	}
	res := authUserResponse{Err: nil}
	return res, nil
}

// errKindToStatus maps service error kinds to the HTTP response codes.
var errKindToStatus = map[ErrorKind]int{
	ErrBadParams: http.StatusBadRequest,
	ErrNotFound:  http.StatusNotFound,
	ErrConflict:  http.StatusConflict,
	ErrInternal:  http.StatusInternalServerError,
}

// encodeError writes a service error to the given http.ResponseWriter.
func encodeError(w http.ResponseWriter, err error, writeMessage bool) error {
	status := http.StatusInternalServerError
	message := err.Error()
	if err, ok := err.(*Error); ok {
		if s, ok := errKindToStatus[err.Kind]; ok {
			status = s
		}
		if err.Kind == ErrInternal {
			message = "internal error"
		} else {
			message = err.Message
		}
	}
	w.WriteHeader(status)
	if writeMessage {
		_, writeErr := io.WriteString(w, message)
		return writeErr
	}
	return nil
}

// decodeError reads a service error from the given *http.Response.
func decodeError(r *http.Response) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, io.LimitReader(r.Body, 1024)); err != nil {
		return fmt.Errorf("%d: %s", r.StatusCode, http.StatusText(r.StatusCode))
	}
	msg := strings.TrimSpace(buf.String())
	if msg == "" {
		msg = http.StatusText(r.StatusCode)
	}
	for kind, status := range errKindToStatus {
		if status == r.StatusCode {
			return &Error{Kind: kind, Message: msg}
		}
	}
	return fmt.Errorf("%d: %s", r.StatusCode, msg)
}
