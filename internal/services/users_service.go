package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"media-service/internal/model"
	"media-service/internal/repository"
	"media-service/pkg/jwt"
	"time"
)

func Authenticate(email, password, ipAddress, userAgent string) (string, string, error) {
	user, err := repository.GetUserByMail(email)

	if err != nil {
		log.Printf("Error finding user by email %s: %v", email, err)
		return "", "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	expiredAt := time.Now().Add(7 * 24 * time.Hour)
	sessionID, err := repository.CreateSession(user.ID, ipAddress, userAgent, expiredAt)
	if err != nil {
		return "", "", fmt.Errorf("could not save session: %w", err)
	}

	accessToken, refreshToken, jti, err := jwt.GenerateTokens(user.Email)
	if err != nil {
		return "", "", fmt.Errorf("could not generate tokens: %w", err)
	}

	if err = repository.CreateToken(user.ID, sessionID, jti); err != nil {
		return "", "", fmt.Errorf("could not save token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func Logout(refreshTokenString string) error {
	claims, err := jwt.ParseToken(refreshTokenString)
	if err != nil {
		return fmt.Errorf("could not parse token: %w", err)
	}

	jtiStr, ok := claims["jti"].(string)
	if !ok {
		return errors.New("jti not found in token")
	}

	jti, err := uuid.Parse(jtiStr)
	if err != nil {
		return fmt.Errorf("failed to parse jti: %w", err)
	}

	token, err := repository.GetTokenByJTI(jti)
	if err != nil {
		return fmt.Errorf("could not get token by jti: %w", err)
	}

	if err := repository.RevokeToken(jti); err != nil {
		return fmt.Errorf("could not revoke token: %w", err)
	}

	if err := repository.ExpireSession(token.SessionID); err != nil {
		return fmt.Errorf("could not expire session: %w", err)
	}

	return nil
}

func LogoutAll(email string, refreshTokenString string) error {
	claims, err := jwt.ParseToken(refreshTokenString)
	if err != nil {
		return fmt.Errorf("could not parse token: %w", err)
	}

	typetoken, ok := claims["type"].(string)
	if typetoken != "refresh" || !ok {
		return errors.New("invalid token type: expected refresh token")
	}

	jtiStr, ok := claims["jti"].(string)
	if !ok {
		return errors.New("jti not found in token")
	}

	jti, err := uuid.Parse(jtiStr)
	if err != nil {
		return fmt.Errorf("failed to parse jti: %w", err)
	}

	refreshToken, err := repository.GetTokenByJTI(jti)
	if err != nil {
		return fmt.Errorf("could not get token by jti: %w", err)
	}

	user, err := repository.GetUserByMail(email)
	if err != nil {
		return fmt.Errorf("could not find user by email %s: %v", email, err)
	}

	if user.ID != refreshToken.UserID {
		return errors.New("token does not belong to the user")
	}

	if err := repository.RevokeAllUserTokens(user.ID); err != nil {
		return fmt.Errorf("could not revoke all tokens: %w", err)
	}

	if err := repository.ExpireAllUserSessions(user.ID); err != nil {
		return fmt.Errorf("could not expire all sessions: %w", err)
	}

	return nil
}

func Refresh(refreshTokenString string) (string, string, error) {
	claims, err := jwt.ParseToken(refreshTokenString)
	if err != nil {
		return "", "", fmt.Errorf("could not parse refresh token: %w", err)
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", errors.New("invalid token type: expected refresh token")
	}

	jtiStr, ok := claims["jti"].(string)
	if !ok {
		return "", "", errors.New("jti not found in token")
	}

	jti, err := uuid.Parse(jtiStr)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse jti: %w", err)
	}

	token, err := repository.GetTokenByJTI(jti)
	if err != nil {
		return "", "", errors.New("refresh token is invalid or has been revoked")
	}

	// TODO: нужно ли делать return в случае ошибки?
	if err := repository.RevokeToken(jti); err != nil {
		log.Printf("could not revoke old refresh token: %v", err)
	}

	email, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("subject not found in token")
	}

	user, err := repository.GetUserByMail(email)
	if err != nil {
		return "", "", fmt.Errorf("user '%s' from token not found", email)
	}

	newAccessToken, newRefreshToken, newJti, err := jwt.GenerateTokens(user.Email)
	if err != nil {
		return "", "", fmt.Errorf("could not generate new tokens: %w", err)
	}

	if err = repository.CreateToken(user.ID, token.SessionID, newJti); err != nil {
		return "", "", fmt.Errorf("could not save new token: %w", err)
	}

	now := time.Now()
	newExpiresAt := now.Add(7 * 24 * time.Hour)
	if err := repository.ExtendSession(token.SessionID, newExpiresAt); err != nil {
		return "", "", fmt.Errorf("could not extend session: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}
	return string(hashedPassword), nil
}

func CreateUser(email string, password string, roles []string) (*model.User, error) {

	_, err := repository.GetUserByMail(email)

	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := repository.CreateUser(email, passwordHash, roles)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetAllUsers() ([]model.UserResponseFull, error) {
	users, err := repository.GetAllUsers()

	if err != nil {
		return nil, errors.New("could not get all users")
	}

	var usersResponse []model.UserResponseFull
	for _, u := range users {
		tzName, _ := u.CreatedAt.Zone()

		usersResponse = append(usersResponse, model.UserResponseFull{
			ID:    u.ID,
			Email: u.Email,
			Roles: u.Roles,
			CreatedAt: model.CreatedAtInfo{
				Date:         u.CreatedAt,
				TimezoneType: 3, // заглушка! в го нет аналога time_zone из РЗ
				Timezone:     tzName,
			},
		})
	}

	return usersResponse, nil
}

func GetUserSessions(userID uint, RefreshToken string) ([]model.SessionResponse, error) {
	claims, err := jwt.ParseToken(RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	jtiStr, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("jti not found in token")
	}

	jti, err := uuid.Parse(jtiStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jti: %w", err)
	}

	token, err := repository.GetTokenByJTI(jti)
	if err != nil {
		return nil, errors.New("current refresh token is invalid or has been revoked")
	}
	sessionID := token.SessionID

	sessions, err := repository.GetSessionsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var response []model.SessionResponse
	for _, s := range sessions {
		response = append(response, model.SessionResponse{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
			ExpiredAt: s.ExpiredAt,
			UserAgent: s.UserAgent,
			IPAddress: s.IPAddress,
			Expired:   s.Expired || time.Now().After(s.ExpiredAt),
			IsCurrent: s.ID == sessionID,
		})
	}

	return response, nil
}

func GetUserByMail(email string) (*model.User, error) {
	user, err := repository.GetUserByMail(email)

	if err != nil {
		return nil, errors.New("could not get the user")
	}

	return user, nil
}
