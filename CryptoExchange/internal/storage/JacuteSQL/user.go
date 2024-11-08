package jacutesql

import (
	"JacuteCE/internal/storage"
	"fmt"
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

	// TODO: fix the race condition here

	err = s.Exec("INSERT INTO user VALUES ('?', '?')", username, token)
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

func (s *Storage) AddLots(userID string, quantity string) error {
	const op = "storage.JacuteSQL.SaveLot"

	data, err := s.Query("SELECT lot.lot_pk FROM lot")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	fmt.Println(data)

	for _, row := range data {
		err := s.Exec("INSERT INTO user_lot VALUES ('?', '?', '?')", userID, row["lot.lot_pk"], quantity)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
