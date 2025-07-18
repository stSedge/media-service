package repository

import (
	"media-service/internal/database"
	"media-service/internal/model"
	"time"
)

func CreateSession(userID uint, ipAddress, userAgent string, expiredAt time.Time) (uint, error) {
	session := model.Session{
		UserID:    userID,
		ExpiredAt: expiredAt,
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}
	res := database.GormDB.Create(&session)
	if res.Error != nil {
		return 0, res.Error
	}
	return session.ID, nil
}

func ExpireSession(sessionID uint) error {
	updates := map[string]interface{}{
		"expired":    true,
		"expired_at": time.Now(),
	}
	err := database.GormDB.Model(&model.Session{}).Where("id = ?", sessionID).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}


func ExtendSession(sessionID uint, newExpiredAt time.Time) error {
	updates := map[string]interface{}{
		"expired_at": newExpiredAt,
		"expired":    false,
	}
	err := database.GormDB.Model(&model.Session{}).Where("id = ?", sessionID).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}

func GetSessionsByUserID(userID uint) ([]model.Session, error) {
	var sessions []model.Session
	err := database.GormDB.Where("user_id = ?", userID).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func ExpireAllUserSessions(userID uint) error {
	updates := map[string]interface{}{
		"expired":    true,
		"expired_at": time.Now(),
	}
	err := database.GormDB.Model(&model.Session{}).Where("user_id = ?", userID).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}
