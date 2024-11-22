package jacutesql

import (
	"CryptoExchange/internal/models"
	"CryptoExchange/internal/storage"
	"fmt"
	"strconv"
)

func (s *Storage) GetUserLots(userID string) ([]*models.UserLot, error) {
	const op = "storage.JacuteSQL.GetUserLots"

	data, err := s.Query("SELECT user_lot.lot_id, user_lot.quantity FROM user_lot WHERE user_lot.user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := make([]*models.UserLot, len(data))
	for i, row := range data {
		lotID, err := strconv.Atoi(row["user_lot.lot_id"])
		if err != nil {
			continue
		}
		quantity, err := strconv.ParseFloat(row["user_lot.quantity"], 64)
		if err != nil {
			continue
		}

		result[i] = &models.UserLot{
			LotID:    lotID,
			Quantity: quantity,
		}
	}

	return result, nil
}

func (s *Storage) Pay(userID, lotID string, price float64) (string, error) {
	const op = "storage.JacuteSQL.Pay"

	// race condition
	data, err := s.Query("SELECT user_lot.quantity FROM user_lot WHERE user_lot.user_id = '?' AND user_lot.lot_id = '?'", userID, lotID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	curQuantity, err := strconv.ParseFloat(data[0]["user_lot.quantity"], 64)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if curQuantity < price {
		return "", storage.ErrNotEnoughMoney
	}
	newQuantity := curQuantity - price

	err = s.Delete("DELETE FROM user_lot WHERE user_lot.user_id = '?' AND user_lot.lot_id = '?'", userID, lotID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	userLotID, err := s.Insert("INSERT INTO user_lot VALUES ('?', '?', '?')", userID, lotID, strconv.FormatFloat(newQuantity, 'f', -1, 64))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userLotID, nil
}

func (s *Storage) AddMoney(userID, lotID string, quantity float64) (string, error) {
	const op = "storage.JacuteSQL.AddMoney"

	// race condition
	data, err := s.Query("SELECT user_lot.quantity FROM user_lot WHERE user_lot.user_id = '?' AND user_lot.lot_id = '?'", userID, lotID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	curQuantity, err := strconv.ParseFloat(data[0]["user_lot.quantity"], 64)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	newQuantity := curQuantity + quantity

	err = s.Delete("DELETE FROM user_lot WHERE user_lot.user_id = '?' AND user_lot.lot_id = '?'", userID, lotID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	userLotID, err := s.Insert("INSERT INTO user_lot VALUES ('?', '?', '?')", userID, lotID, strconv.FormatFloat(newQuantity, 'f', -1, 64))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userLotID, nil
}
