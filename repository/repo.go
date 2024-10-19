package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	taskservice "test/task-service"
	"time"

	_ "modernc.org/sqlite"
)

type Repository struct {
	Repo *sql.DB
}

type RepositoryProcesser interface {
	AddTask(task taskservice.Task) (string, error)
}

func NewRepo() (*Repository, error) {

	dbFile := os.Getenv("TODO_DFILE")
	if dbFile == "" {
		dbFile = dbCheck()
	}

	fmt.Println(dbFile)

	repo := Repository{}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT,  title TEXT, comment TEXT, repeat TEXT)")
	if err != nil {
		panic(err)
	}
	repo.Repo = db
	if err = db.Ping(); err != nil {
		panic(err)
	}
	return &repo, nil
}

func dbCheck() string {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(appPath)
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	if err != nil {
		os.Create("scheduler.db")
	}

	return dbFile
}

func (repo *Repository) AddTask(task taskservice.Task) (string, error) {
	nextDate, err := task.GetNextRepeatDate()
	if err != nil {
		return "", err
	}

	fmt.Println(nextDate)

	task.NextDate = nextDate

	if task.Title == "" {
		return "", errors.New("no title")
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102") // Присваиваем текущую дату
	}

	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		return "", err
	}
	fmt.Println(date, "1")
	// Здесь можно проверять, если дата уже прошла
	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		fmt.Println(date.Before(time.Now()))
		if task.Repeat == "" {
			date = time.Now() // Если нет повторения, ставим текущую дату
		} else {
			nextDateParsed, err := time.Parse("20060102", task.NextDate)
			if err != nil {
				return "", err
			}
			date = nextDateParsed // Иначе используем следующую дату
		}
	}
	fmt.Println(date, "2")
	task.Date = date.Format("20060102") // Устанавливаем отформатированную дату

	res, err := repo.Repo.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return "", err
	}

	id, _ := res.LastInsertId()
	strid := strconv.Itoa(int(id))
	fmt.Println(strid)
	return strid, nil
}
