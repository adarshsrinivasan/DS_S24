package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adarshsrinivasan/DS_S24/Assignment1/libraries/common"
	"github.com/sirupsen/logrus"
)

func createNewSession(ctx context.Context, userID string, userType common.UserType) (string, int, error) {
	sessionDBObj := SessionDBModel{UserID: userID, UserType: userType}

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
	sessionDBObj := SessionDBModel{ID: sessionID}
	statusCode, err := sessionDBObj.GetSessionByID(ctx)
	if err != nil {
		err := fmt.Errorf("exception while fetching session: %s. %v", sessionID, err)
		logrus.Errorf("getUserIDAndTypeFromSessionID: %v\n", err)
		return "", 0, statusCode, err
	}
	return sessionDBObj.UserID, sessionDBObj.UserType, http.StatusOK, nil
}

func getSessionByUserID(ctx context.Context, userID string) (string, int, error) {
	sessionDBObj := SessionDBModel{UserID: userID}
	statusCode, err := sessionDBObj.GetSessionByUserID(ctx)
	if err != nil {
		err := fmt.Errorf("exception while fetching session by UserID: %s. %v", userID, err)
		logrus.Errorf("getSessionByUserID: %v\n", err)
		return "", statusCode, err
	}
	return sessionDBObj.ID, http.StatusOK, nil
}

func deleteSessionByID(ctx context.Context, sessionID string) (int, error) {
	sessionDBObj := SessionDBModel{ID: sessionID}
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
