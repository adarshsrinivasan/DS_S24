package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adarshsrinivasan/DS_S24/library/common"
	"github.com/sirupsen/logrus"
)

type SessionModel struct {
	ID        string          `json:"id,omitempty" bson:"id" bun:"id,pk"`
	UserID    string          `json:"userID,omitempty" bson:"userID" bun:"userID,notnull,unique"`
	UserType  common.UserType `json:"userType,omitempty" bson:"userType"  bun:"userType,notnull"`
	Version   int             `json:"version" bson:"version" bun:"version,notnull"`
	CreatedAt time.Time       `json:"createdAt,omitempty"  bson:"createdAt" bun:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt,omitempty" bson:"updatedAt" bun:"updatedAt"`
}

func createNewSession(ctx context.Context, userID string, userType common.UserType) (string, int, error) {
	sessionDBObj := SessionModel{UserID: userID, UserType: userType}

	if sessionID, _, err := getSessionByUserID(ctx, userID); err == nil {
		return sessionID, http.StatusOK, nil
	}

	statusCode, err := sessionDBObj.CreateSession(ctx)
	if err != nil {
		err := fmt.Errorf("exception while creating session.%v", err)
		logrus.Errorf("createNewSession: %v\n", err)
		return "", statusCode, err
	}
	return sessionDBObj.ID, http.StatusOK, nil
}

func getUserIDAndTypeFromSessionID(ctx context.Context, sessionID string) (string, common.UserType, int, error) {
	sessionDBObj := SessionModel{ID: sessionID}
	statusCode, err := sessionDBObj.GetSessionByID(ctx)
	if err != nil {
		err := fmt.Errorf("exception while fetching session: %s. %v", sessionID, err)
		logrus.Errorf("getUserIDAndTypeFromSessionID: %v\n", err)
		return "", 0, statusCode, err
	}
	return sessionDBObj.UserID, sessionDBObj.UserType, http.StatusOK, nil
}

func getSessionByUserID(ctx context.Context, userID string) (string, int, error) {
	sessionDBObj := SessionModel{UserID: userID}
	statusCode, err := sessionDBObj.GetSessionByUserID(ctx)
	if err != nil {
		err := fmt.Errorf("exception while fetching session by UserID: %s. %v", userID, err)
		logrus.Errorf("getSessionByUserID: %v\n", err)
		return "", statusCode, err
	}
	return sessionDBObj.ID, http.StatusOK, nil
}

func deleteSessionByID(ctx context.Context, sessionID string) (int, error) {
	sessionDBObj := SessionModel{ID: sessionID}
	return sessionDBObj.DeleteSessionByID(ctx)
}

func validateSessionID(sessionID string) bool {
	_, _, _, err := getUserIDAndTypeFromSessionID(context.TODO(), sessionID)
	if err != nil {
		err := fmt.Errorf("exception while fetching Session with ID %s. %v", sessionID, err)
		logrus.Errorf("validateSessionID: %v\n", err)
		return false
	}
	return true
}
