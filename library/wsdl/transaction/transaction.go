// Code generated by gowsdl DO NOT EDIT.

package transaction

import (
	"context"
	"encoding/xml"
	"github.com/hooklift/gowsdl/soap"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type AnyType struct {
	InnerXML string `xml:",innerxml"`
}

type AnyURI string

type NCName string

type TransactionRequest struct {
	XMLName xml.Name `xml:"http://localhost:50003 TransactionRequest"`

	Name string `xml:"Name,omitempty" json:"Name,omitempty"`

	CreditCardDetails string `xml:"CreditCardDetails,omitempty" json:"CreditCardDetails,omitempty"`

	Expiry string `xml:"expiry,omitempty" json:"expiry,omitempty"`
}

type TransactionResponse struct {
	XMLName xml.Name `xml:"http://localhost:50003 TransactionResponse"`

	Approved bool `xml:"approved,omitempty" json:"approved,omitempty"`
}

type TransactionServicePortType interface {
	IsTransactionApproved(request *TransactionRequest) (*TransactionResponse, error)

	IsTransactionApprovedContext(ctx context.Context, request *TransactionRequest) (*TransactionResponse, error)
}

type transactionServicePortType struct {
	client *soap.Client
}

func NewTransactionServicePortType(client *soap.Client) TransactionServicePortType {
	return &transactionServicePortType{
		client: client,
	}
}

func (service *transactionServicePortType) IsTransactionApprovedContext(ctx context.Context, request *TransactionRequest) (*TransactionResponse, error) {
	response := new(TransactionResponse)
	err := service.client.CallContext(ctx, "http://localhost:50003/isTransactionApproved", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *transactionServicePortType) IsTransactionApproved(request *TransactionRequest) (*TransactionResponse, error) {
	return service.IsTransactionApprovedContext(
		context.Background(),
		request,
	)
}