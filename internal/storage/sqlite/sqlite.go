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

//go doesnt have any constructor, so we make

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
	// 1. Correctly retrieve the ID using QueryRow and Scan
	fmt.Printf("uhh")
	var accountID int64

	// We use a '?' placeholder here to actually protect against SQL injection.
	// Your previous code used fmt.Sprintf inside Prepare, which is not secure.
	err := s.Db.QueryRow("SELECT id FROM Accounts WHERE name = ?", accountname).Scan(&accountID)
	if err != nil {
		return 0, fmt.Errorf("could not find account %s: %v", accountname, err)
	}
	fmt.Printf("uhh")
	// 2. Now accountID is safely the integer you need (e.g., 3)
	// We use Sprintf here because table names cannot be parameterized with '?'
	// 1. Properly format the string using Sprintf
	query := fmt.Sprintf("INSERT INTO well%d (name, today, comment) VALUES (?, ?, ?)", accountID)

	// 2. Pass that finished string to Prepare
	smt, err := s.Db.Prepare(query)
	if err != nil {
		// If you get Code 1 here now, it means the table 'well1' doesn't exist yet.
		return 0, fmt.Errorf("prepare failed for table well%d: %w", accountID, err)
	}
	defer smt.Close()
	fmt.Println(comment)
	fmt.Println(name)
	// 3. Execute with data
	result, err := smt.Exec(name, today, comment)
	if err != nil {
		return 0, err
	}
	fmt.Printf("uhh")
	return result.LastInsertId()
}

// 1. Change the return type to match the Interface
func (s *SQL) Getmovies(id int64) ([]types.GetDataResponse, error) {
	query := fmt.Sprintf("SELECT id, name, today, comment FROM well%d", id)

	rows, err := s.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fmt.Printf("uhh")
	// 2. Use the Storage struct here as well
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

/*
func (s *SQL) GetMovieByID(id int64) (types.Movie, error) {
	smt, err := s.Db.Prepare("SELECT * FROM movies WHERE id=?") //to protect us from sql injections
	if err != nil {
		return types.Movie{}, err
	}
	defer smt.Close()

	var moviee types.Movie

	err = smt.QueryRow(id).Scan(&moviee.Id, &moviee.Name, &moviee.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Movie{}, fmt.Errorf("NO MOVIE WITH THIS ID")
		}
		return types.Movie{}, fmt.Errorf("SOMETHING WRONG WITH RUNNING THE QUERY")
	}

	return moviee, nil
}
*/
