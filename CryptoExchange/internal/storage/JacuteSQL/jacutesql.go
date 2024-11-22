package jacutesql

import (
	"CryptoExchange/internal/storage"
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	InsertType = iota + 1
	SelectType
	DeleteType
)
const ExecutedCommandOutput = "command executed successfully"

var (
	InsertMaliciousRegexp = regexp.MustCompile("[" + regexp.QuoteMeta(",'()=") + "]")
	SelectMaliciousRegexp = regexp.MustCompile("(?i)\b(AND|OR)\b|[" + regexp.QuoteMeta("'=,") + "]")
)

var (
	InsertStartsWithRegexp = regexp.MustCompile(`(?i)^INSERT\s+INTO`)
	SelectStartsWithRegexp = regexp.MustCompile(`(?i)^SELECT`)
	DeleteStartsWithRegexp = regexp.MustCompile(`(?i)^DELETE\s+FROM`)
)

type Storage struct {
	ip   string
	port int
}

func New(ip string, port int, lots []string) *Storage {
	s := &Storage{
		ip:   ip,
		port: port,
	}

	return s
}

func (s *Storage) write(data string) (string, error) {
	const op = "storage.JacuteSQL.write"

	conn, err := s.getConn()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, storage.ErrConnect)
	}
	defer conn.Close()

	buf := make([]byte, 3) // read '>> '
	_, err = conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	n, err := conn.Write([]byte(data + "\n"))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if n == 0 {
		return "", fmt.Errorf("%s: no bytes written", op)
	}

	var outputBuilder strings.Builder
	buf = make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		outputBuilder.Write(buf[:n])
		if strings.Contains(outputBuilder.String(), ">> ") {
			break
		}
	}

	output := outputBuilder.String()

	if strings.Contains(output, "is already blocked") {
		return "", fmt.Errorf("%s: %s", op, "table blocked")
	}

	if !strings.Contains(output, "command executed successfully") {
		return "", fmt.Errorf("%s: %w", op, storage.ErrSQLExecFailed)
	}

	// if !strings.Contains(output, "Incorrect number of columns") {
	// 	return "", fmt.Errorf("%s: %w", op, storage.ErrIncorrectNumberOfColumns)
	// }

	return output, nil
}

func (s *Storage) argSanitize(query string, args ...string) (string, error) {
	var commandType int
	if InsertStartsWithRegexp.MatchString(query) {
		commandType = InsertType
	} else if SelectStartsWithRegexp.MatchString(query) {
		commandType = SelectType
	} else if DeleteStartsWithRegexp.MatchString(query) {
		commandType = DeleteType
	} else {
		return "", storage.ErrInvalidSQLCommand
	}

	if strings.Count(query, "?") != len(args) {
		fmt.Println(3)
		return "", storage.ErrInvalidSQLCommand
	}

	for _, arg := range args {
		if commandType == InsertType {
			if InsertMaliciousRegexp.MatchString(arg) {
				return "", fmt.Errorf("%w. argument: %s", storage.ErrMaliciousParameter, arg)
			}
		} else if commandType == SelectType || commandType == DeleteType {
			if SelectMaliciousRegexp.MatchString(arg) {
				return "", fmt.Errorf("%w. argument: %s", storage.ErrMaliciousParameter, arg)
			}
		}
		query = strings.Replace(query, "?", arg, 1)
	}

	return query, nil
}

func (s *Storage) Delete(query string, args ...string) error {
	const op = "storage.JacuteSQL.Exec"

	newQuery, err := s.argSanitize(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.write(newQuery)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Insert(query string, args ...string) (string, error) {
	const op = "storage.JacuteSQL.Exec"

	newQuery, err := s.argSanitize(query, args...)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	output, err := s.write(newQuery)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	splitted := strings.Split(output, "\n")
	splitted = splitted[2 : len(splitted)-1]

	if splitted[0] == "" {
		return "", fmt.Errorf("%s: %s", op, "failed to insert")
	}

	return splitted[0], nil
}

func (s *Storage) Query(query string, args ...string) ([]map[string]string, error) {
	const op = "storage.JacuteSQL.Query"

	newQuery, err := s.argSanitize(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	output, err := s.write(newQuery)
	if err != nil {
		return nil, err
	}
	rows := strings.Split(output, "\n")
	if len(rows) == 3 {
		return []map[string]string{}, nil
	}
	rows = rows[2 : len(rows)-2] // remove user-friendly interface strings
	header := strings.Split(rows[0], ",")
	rows = rows[1:] // remove header
	table := make([]map[string]string, len(rows))
	for i, row := range rows {
		cols := strings.Split(row, ",")
		for j := range cols {
			if table[i] == nil {
				table[i] = make(map[string]string)
			}
			table[i][header[j]] = cols[j]
		}
	}

	return table, nil
}

func (s *Storage) Destroy() error {
	const op = "storage.JacuteSQL.Destroy"

	err := s.Delete("DELETE FROM lot, pair, user, user_lot, order")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) getIDByParam(table []map[string]string, idName string, paramName string, paramValue string) string {
	for _, row := range table {
		if row[paramName] == paramValue {
			return row[idName]
		}
	}
	return ""
}

func (s *Storage) MakeMigrations(lots []string) {
	// create lots
	table, err := s.Query("SELECT lot.name FROM lot")
	if err != nil {
		panic(err)
	}
	if len(table) == 0 {
		for _, lot := range lots {
			_, err := s.Insert("INSERT INTO lot VALUES ('?')", lot)
			if err != nil {
				panic(err)
			}
		}
	}

	// create pairs
	table, err = s.Query("SELECT pair.pair_pk FROM pair")
	if err != nil {
		panic(err)
	}

	if len(table) == 0 {
		table, err := s.Query("SELECT lot.lot_pk, lot.name FROM lot")
		fmt.Println(table)
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(lots); i++ {
			for j := 0; j < len(lots); j++ {
				if i == j {
					continue
				}

				firstLotID := s.getIDByParam(table, "lot.lot_pk", "lot.name", lots[i])
				secondLotID := s.getIDByParam(table, "lot.lot_pk", "lot.name", lots[j])
				if firstLotID == "" {
					panic(fmt.Sprintf("can't find id for %s", lots[i]))
				}
				if secondLotID == "" {
					panic(fmt.Sprintf("can't find id for %s", lots[j]))
				}

				_, err := s.Insert("INSERT INTO pair VALUES ('?', '?')", firstLotID, secondLotID)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (s *Storage) getConn() (net.Conn, error) {
	const op = "storage.JacuteSQL.getConn"

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return conn, nil
}
