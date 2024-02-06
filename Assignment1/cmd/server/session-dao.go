package main

import (
	"context"
	"fmt"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/db/sql"
	"net/http"
	"reflect"
	"time"

	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/db"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun/schema"
)

const (
	SessionTableName = "session_data"
)

type SessionTableOps interface {
	CreateSession(ctx context.Context) (int, error)
	GetSessionByID(ctx context.Context) (int, error)
	GetSessionByUserID(ctx context.Context) (int, error)
	DeleteSessionByID(ctx context.Context) (int, error)
}

type SessionTableModel struct {
	schema.BaseModel `bun:"table:session_data,alias:session"`
	ID               string          `json:"id,omitempty" bson:"id" bun:"id,pk"`
	UserID           string          `json:"userID,omitempty" bson:"userID" bun:"userID,notnull,unique"`
	UserType         common.UserType `json:"userType,omitempty" bson:"userType"  bun:"userType,notnull"`
	Version          int             `json:"version" bson:"version" bun:"version,notnull"`
	CreatedAt        time.Time       `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}

func CreateSessionTable(ctx context.Context) error {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateSessionTable: %v\n", err)
		return err
	}
	defer client.Close(ctx)

	tableSchemaPtr := reflect.New(reflect.TypeOf(SessionTableModel{}))

	if err := client.CreateTable(ctx, tableSchemaPtr.Interface(), SessionTableName, nil); err != nil {
		err := fmt.Errorf("exception while creating table %s. %v", err, SessionTableName)
		logrus.Errorf("CreateSessionTable: %v\n", err)
		return err
	}

	return nil
}

func (session *SessionTableModel) CreateSession(ctx context.Context) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("CreateSession: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)

	session.ID = uuid.New().String()
	session.Version = 0
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	if err := client.Insert(ctx, session, SessionTableName); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", SessionTableName, err)
		logrus.Errorf("CreateSession: %v\n", err)
		return http.StatusInternalServerError, err
	}
	logrus.Infof("CreateSession: Successfully created session for userID %s\n", session.UserID)
	return http.StatusOK, nil
}

func (session *SessionTableModel) GetSessionByID(ctx context.Context) (int, error) {
	var (
		existingSession *SessionTableModel
		err             error
	)

	if existingSession, _, err = session.getByColumn(ctx, "id", session.ID); err != nil || existingSession.ID != session.ID {
		err := fmt.Errorf("unable to find session with with id: %s. %v", session.ID, err)
		logrus.Errorf("GetSessionByID: %v\n", err)
		return http.StatusBadRequest, err
	}
	copySessionObj(existingSession, session)

	return http.StatusOK, nil
}

func (session *SessionTableModel) GetSessionByUserID(ctx context.Context) (int, error) {
	var (
		existingSession *SessionTableModel
		err             error
	)

	if existingSession, _, err = session.getByColumn(ctx, "userID", session.UserID); err != nil || existingSession.UserID != session.UserID {
		err := fmt.Errorf("unable to find session with with id: %s. %v", session.ID, err)
		logrus.Errorf("GetSessionByUserID: %v\n", err)
		return http.StatusBadRequest, err
	}
	copySessionObj(existingSession, session)

	return http.StatusOK, nil
}

func (session *SessionTableModel) DeleteSessionByID(ctx context.Context) (int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("DeleteSessionByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	defer client.Close(ctx)
	whereClauses := []db.WhereClauseType{
		{
			ColumnName:   "id",
			RelationType: db.EQUAL,
			ColumnValue:  session.ID,
		},
	}

	if err := client.Delete(ctx, session, SessionTableName, whereClauses); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", SessionTableName, err)
		logrus.Errorf("DeleteSessionByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (session *SessionTableModel) getByColumn(ctx context.Context, columnName string, columnValue interface{}) (*SessionTableModel, int, error) {
	client, err := sql.NewClient(ctx, ServiceName, SQLSchemaName)
	if err != nil {
		err = fmt.Errorf("exception while creating SQLDB client. %v", err)
		logrus.Errorf("DeleteSessionByID: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	defer client.Close(ctx)
	whereClause := []db.WhereClauseType{
		{
			ColumnName:   columnName,
			RelationType: db.EQUAL,
			ColumnValue:  columnValue,
		},
	}
	resultSession := SessionTableModel{}

	if _, err := client.Read(ctx, SessionTableName, nil, whereClause, nil, nil, nil, true, &resultSession); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	return &resultSession, http.StatusOK, nil
}

func copySessionObj(from, to *SessionTableModel) {
	to.ID = from.ID
	to.UserID = from.UserID
	to.UserType = from.UserType
	to.Version = from.Version
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}
