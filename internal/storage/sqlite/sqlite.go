package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" //_ is used bec this is used behind the scenes
	"github.com/okishiro/pidgey/internal/config"
	"github.com/okishiro/pidgey/internal/types"
)

type SQL struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*SQL, error) {
	db, err := sql.Open("sqlite3", cfg.Storage_path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Accounts(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)`)

	if err != nil {
		return nil, err
	}

	return &SQL{
		Db: db,
	}, nil

}

/*
In Go, this s is called a Receiver. It is what makes this function a Method rather than a standalone function.
Inm languages like Python, Java, or C#, s is the equivalent of self or this.
By putting (s *SQL) before the function name, we are telling Go: "This function belongs to the SQL struct." we can only call this function if we have an instance of SQL.
*/
func (s *SQL) CreateAccount(name string) (int64, error) {
	smt, err := s.Db.Prepare("INSERT INTO Accounts(name) VALUES(? )")
	if err != nil {
		return 0, err
	}

	result, err := smt.Exec(name)
	if err != nil {
		return 0, err
	}

	LastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return LastId, nil
}

func (s *SQL) CreateTable(id int64, path string) error {

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS well%d (id INTEGER PRIMARY KEY, name TEXT, today DATETIME, comment TEXT)", id)

	_, err = db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (s *SQL) CreateMovieEntry(accountname string, name string, today time.Time, comment string) (int64, error) {
	var accountID int64

	// We use a '?' placeholder here to actually protect against SQL injection.
	err := s.Db.QueryRow("SELECT id FROM Accounts WHERE name = ?", accountname).Scan(&accountID)
	if err != nil {
		return 0, fmt.Errorf("could not find account %s: %v", accountname, err)
	}
	fmt.Printf("uhh")
	query := fmt.Sprintf("INSERT INTO well%d (name, today, comment) VALUES (?, ?, ?)", accountID)
	smt, err := s.Db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("prepare failed for table well%d: %w", accountID, err)
	}
	defer smt.Close()
	fmt.Println(comment)
	fmt.Println(name)

	result, err := smt.Exec(name, today, comment)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *SQL) Getmovies(id int64) ([]types.GetDataResponse, error) {
	query := fmt.Sprintf("SELECT id, name, today, comment FROM well%d", id)

	rows, err := s.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []types.GetDataResponse

	for rows.Next() {
		var movie types.GetDataResponse
		err := rows.Scan(&movie.ID, &movie.MovieName, &movie.Timestamp, &movie.Comment)
		if err != nil {
			return nil, err
		}
		results = append(results, movie)
	}

	return results, nil
}
