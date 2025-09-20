package services

import (
	"log"
	"reading-microservices/user-service/models"
	"reading-microservices/user-service/repositories"
)

type LoginLogger struct {
	userRepo repositories.UserRepository
}

func NewLoginLogger(userRepo repositories.UserRepository) *LoginLogger {
	return &LoginLogger{userRepo: userRepo}
}

func (l *LoginLogger) LogSuccess(userID string, sessionID *string, platform, deviceID string) {
	logEntry := &models.LoginLog{
		UserID:    userID,
		SessionID: sessionID,
		LoginType: "password",
		Platform:  platform,
		IsSuccess: true,
	}
	if deviceID != "" {
		logEntry.DeviceID = &deviceID
	}
	if err := l.userRepo.CreateLoginLog(logEntry); err != nil {
		log.Printf("warning: create login log failed: %v", err)
	}
}

func (l *LoginLogger) LogFailure(userID, reason, platform, deviceID string) {
	logEntry := &models.LoginLog{
		UserID:        userID,
		LoginType:     "password",
		Platform:      platform,
		IsSuccess:     false,
		FailureReason: &reason,
	}
	if deviceID != "" {
		logEntry.DeviceID = &deviceID
	}
	if err := l.userRepo.CreateLoginLog(logEntry); err != nil {
		log.Printf("warning: create login log failed: %v", err)
	}
}
