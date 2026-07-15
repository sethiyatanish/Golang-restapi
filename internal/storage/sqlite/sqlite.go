package sqlite

import (
	"database/sql"
	"fmt"

	config "github.com/sethiyatanish/student-api/internal"
	"github.com/sethiyatanish/student-api/internal/types"

	"github.com/go-sql-driver/mysql"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	var db *sql.DB
	var err error

	if cfg.StoragePath != "" {
		// Use SQLite
		db, err = sql.Open("sqlite3", cfg.StoragePath)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
		)`)
		if err != nil {
			return nil, err
		}
	} else if cfg.Database_Url != "" {
		// Use MySQL
		mysqlCfg, err := mysql.ParseDSN(cfg.Database_Url)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mysql dsn: %w", err)
		}

		dbname := mysqlCfg.DBName
		mysqlCfg.DBName = ""
		baseDSN := mysqlCfg.FormatDSN()

		tempDb, err := sql.Open("mysql", baseDSN)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to mysql server: %w", err)
		}

		_, err = tempDb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbname))
		tempDb.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to create database %s: %w", dbname, err)
		}

		db, err = sql.Open("mysql", cfg.Database_Url)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100),
		age INT
		)`)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("neither storage_path nor database_url is provided in configuration")
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) UpdateStudent(id int64, name string, email string, age int) error {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no student found with id %d to update", id)
	}

	return nil
}

func (s *Sqlite) DeleteStudent(id int64) error {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no student found with id %d to delete", id)
	}

	return nil
}