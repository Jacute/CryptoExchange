package jacutesql

import (
	"CryptoExchange/internal/models"
	"CryptoExchange/internal/storage"
	"fmt"
	"strconv"
)

func (s *Storage) SaveUser(username string, token string) (string, error) {
	const op = "storage.JacuteSQL.SaveUser"

	// check if the user exists
	data, err := s.Query("SELECT user.user_pk FROM user WHERE user.username = '?'", username)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if len(data) != 0 {
		return "", fmt.Errorf("%s: %w", op, storage.ErrUserExists)
	}

	// race condition

	_, err = s.Insert("INSERT INTO user VALUES ('?', '?')", username, token)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	data, err = s.Query("SELECT user.user_pk FROM user WHERE user.token = '?'", token)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if len(data) == 0 {
		return "", fmt.Errorf("%s: can't find user with token '%s'", op, token)
	}
	id, ok := data[0]["user.user_pk"]

	if !ok {
		return "", fmt.Errorf("%s: can't find user_pk", op)
	}

	return id, nil
}

func (s *Storage) GetUserByToken(token string) (*models.User, error) {
	const op = "storage.JacuteSQL.GetUserByToken"

	data, err := s.Query("SELECT user.user_pk, user.username, user.token FROM user WHERE user.token = '?'", token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	}
	id, err := strconv.Atoi(data[0]["user.user_pk"])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.User{
		ID:       id,
		Username: data[0]["user.username"],
		Token:    data[0]["user.token"],
	}, nil
}

func (s *Storage) EnoughQuantity(userID, lotID int, neededQuantity float64) (bool, error) {
	data, _ := s.Query("SELECT user_lot.quantity FROM user_lot WHERE user.user_id = '?' AND user.lot_id = '?'", strconv.Itoa(userID), strconv.Itoa(lotID))
	if len(data) == 0 {
		return false, fmt.Errorf("can't find lot")
	}

	curQuantity, err := strconv.ParseFloat(data[0]["user_lot.quantity"], 64)
	if err != nil {
		return false, fmt.Errorf("can't parse user lot quantity: %w", err)
	}

	if curQuantity < neededQuantity {
		return false, nil
	}

	return true, nil
}
