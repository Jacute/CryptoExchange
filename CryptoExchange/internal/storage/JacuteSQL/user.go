package jacutesql

import "fmt"

func (s *Storage) SaveUser(username string, token string) (string, error) {
	const op = "storage.JacuteSQL.SaveUser"

	err := s.Exec("INSERT INTO user VALUES ('?', '?')", username, token)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	data, err := s.Query("SELECT user.user_pk FROM user WHERE user.token = '?'", token)
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

func (s *Storage) SaveLot(userID string, lotID string, quantity string) error {
	const op = "storage.JacuteSQL.SaveLot"

	err := s.Exec("INSERT INTO lot  VALUES ('?', '?', '?')", userID, lotID, quantity)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
