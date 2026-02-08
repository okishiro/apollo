package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" //_ is used bec this is used behind the scenes
	"github.com/okishiro/pidgey/internal/config"
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
	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS _%d (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, time STRING, comment STRING)",
		id,
	)

	_, err = db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (s *SQL) CreateMovie(accountname string, name string, today time.Time, comment string) (int64, error) {
	query1 := fmt.Sprintf(
		"SELECT id FROM Accounts WHERE name = '%s' ", accountname,
	)
	smt, err := s.Db.Prepare(query1) //to protect us from sql injections
	if err != nil {
		return 0, err
	}
	defer smt.Close()
	id, err := smt.Exec(query1)
	if err != nil {
		return 0, err
	}
	fmt.Printf("id: %d\n", id)
	query := fmt.Sprintf(
		"INSERT INTO _%d(name, today, comment) VALUES(?,?,? )", id,
	)
	smt, err = s.Db.Prepare(query) //to protect us from sql injections
	if err != nil {
		return 0, err
	}
	defer smt.Close()

	result, err := smt.Exec(name, today, comment)
	if err != nil {
		return 0, err
	}

	lastid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastid, nil
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
