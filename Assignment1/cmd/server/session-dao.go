package main

import (
	"context"
	"fmt"
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

type SessionDBOps interface {
	CreateSession(ctx context.Context) (int, error)
	GetSessionByID(ctx context.Context) (int, error)
	GetSessionByUserID(ctx context.Context) (int, error)
	DeleteSessionByID(ctx context.Context) (int, error)
}

type SessionDBModel struct {
	schema.BaseModel `bun:"table:session_data,alias:session"`
	ID               string          `json:"id,omitempty" bson:"id" bun:"id,pk"`
	UserID           string          `json:"userID,omitempty" bson:"userID" bun:"userID,notnull,unique"`
	UserType         common.UserType `json:"userType,omitempty" bson:"userType"  bun:"userType,notnull"`
	CreatedAt        time.Time       `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt,omitempty"`
	UpdatedAt        time.Time       `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt,omitempty"`
}

func CreateSessionTable(ctx context.Context) error {

	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while creating %s table. %v", SessionTableName, err)
		logrus.Errorf("CreateSessionTable: %v\n", err)
		return err
	}

	tableSchemaPtr := reflect.New(reflect.TypeOf(SessionDBModel{}))
	createTableQuery := db.SqlDBClient.NewCreateTable().
		Model(tableSchemaPtr.Interface()).
		IfNotExists()

	_, err := createTableQuery.Exec(ctx)
	if err != nil {
		err := fmt.Errorf("exception while creaiting event table %s. %v", err, SessionTableName)
		logrus.Errorf("CreateSessionTable: %v\n", err)
		return err
	}

	return nil
}

func (session *SessionDBModel) CreateSession(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("CreateSession: %v\n", err)
		return http.StatusInternalServerError, err
	}

	session.ID = uuid.New().String()
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	if _, err := db.SqlDBClient.NewInsert().Model(session).Exec(ctx); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Insert", SessionTableName, err)
		logrus.Errorf("CreateSession: %v\n", err)
		return http.StatusInternalServerError, err
	}
	logrus.Infof("CreateSession: Successfully created session for userID %s\n", session.UserID)
	return http.StatusOK, nil
}

func (session *SessionDBModel) GetSessionByID(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("GetSessionByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	var (
		existingSession *SessionDBModel
		err             error
	)

	if existingSession, _, err = session.getByColumn(ctx, "id", session.ID); err != nil || existingSession.ID != session.ID {
		err := fmt.Errorf("unable to find session with with id: %s. %v", session.ID, err)
		logrus.Errorf("GetSessionByID: %v\n", err)
		return http.StatusBadRequest, err
	}
	session.ID = existingSession.ID
	session.UserID = existingSession.UserID
	session.UserType = existingSession.UserType
	session.CreatedAt = existingSession.CreatedAt
	session.UpdatedAt = existingSession.UpdatedAt

	return http.StatusOK, nil
}

func (session *SessionDBModel) GetSessionByUserID(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("GetSessionByUserID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	var (
		existingSession *SessionDBModel
		err             error
	)

	if existingSession, _, err = session.getByColumn(ctx, "userID", session.UserID); err != nil || existingSession.UserID != session.UserID {
		err := fmt.Errorf("unable to find session with with id: %s. %v", session.ID, err)
		logrus.Errorf("GetSessionByUserID: %v\n", err)
		return http.StatusBadRequest, err
	}
	session.ID = existingSession.ID
	session.UserID = existingSession.UserID
	session.UserType = existingSession.UserType
	session.CreatedAt = existingSession.CreatedAt
	session.UpdatedAt = existingSession.UpdatedAt

	return http.StatusOK, nil
}

func (session *SessionDBModel) DeleteSessionByID(ctx context.Context) (int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("DeleteSessionByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	whereClauses := []db.WhereClauseType{
		{
			ColumnName:   "id",
			RelationType: db.EQUAL,
			ColumnValue:  session.ID,
		},
	}

	deleteQuery := db.SqlDBClient.NewDelete().
		Model(session)

	// prepare whereClause.
	queryStr, vals, err := db.CreateWhereClause(ctx, whereClauses)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", SessionTableName, err)
		logrus.Errorf("DeleteSessionByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	deleteQuery = deleteQuery.Where(queryStr, vals...)

	if _, err := deleteQuery.Exec(ctx); err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Delete", SessionTableName, err)
		logrus.Errorf("DeleteSessionByID: %v\n", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (session *SessionDBModel) getByColumn(ctx context.Context, columnName string, columnValue interface{}) (*SessionDBModel, int, error) {
	if err := db.VerifySQLDatabaseConnection(ctx, db.SqlDBClient); err != nil {
		err := fmt.Errorf("exception while verifying DB connection. %v", err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, http.StatusInternalServerError, err
	}
	whereClause := []db.WhereClauseType{
		{
			ColumnName:   columnName,
			RelationType: db.EQUAL,
			ColumnValue:  columnValue,
		},
	}
	resultSession := SessionDBModel{}
	_, statusCode, err := db.ReadUtil(ctx, SessionTableName, nil, whereClause, nil, nil, nil, true, &resultSession)
	if err != nil {
		err := fmt.Errorf("unable to Perform %s Operation on Table: %s. %v", "Read", CartTableName, err)
		logrus.Errorf("getByColumn: %v\n", err)
		return nil, statusCode, err
	}
	return &resultSession, http.StatusOK, nil
}
