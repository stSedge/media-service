package repository

import (
	"github.com/google/uuid"
	"log"
	"media-service/internal/database"
	"media-service/internal/model"
)

func CreateToken(userID, sessionID uint, jti uuid.UUID) error {
	token := &model.Token{
		UserID: userID,
		JTI:    jti,
	}
	res := database.GormDB.Create(token)
	if res.Error != nil {
		log.Printf("Error creating token: %v", res.Error)
		return res.Error
	}
	log.Printf("Token for user %d created with JTI %s", token.UserID, token.JTI.String())
	return nil
}

func GetTokenByJTI(jti uuid.UUID) (*model.Token, error) {
	var token model.Token
	err := database.GormDB.Where("jti = ? AND is_active = ?", jti, true).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func RevokeToken(jti uuid.UUID) error {
	res := database.GormDB.Model(&model.Token{}).Where("jti = ?", jti).Update("is_active", false)
	if res.Error != nil {
		log.Printf("Error revoking token: %v", res.Error)
		return res.Error
	}

	if res.RowsAffected > 0 {
		log.Printf("Token with JTI %s revoked", jti.String())
	}

	return nil
}

func RevokeAllUserTokens(userID uint) error {
	res := database.GormDB.Model(&model.Token{}).Where("user_id = ? AND is_active = ?", userID, true).Update("is_active", false)
	if res.Error != nil {
		log.Printf("Error revoking users tokens %d: %v", userID, res.Error)
		return res.Error
	}
	log.Printf("%d token(s) for user %d revoked", res.RowsAffected, userID)
	return nil
}
