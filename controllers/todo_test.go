package controllers

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/ichtrojan/go-todo/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var mock sqlmock.Sqlmock

func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	database = db
	return db, mock
}

func TestShow(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "item", "completed"}).
		AddRow(1, "Buy groceries", 0).
		AddRow(2, "Walk the dog", 1)

	mock.ExpectQuery("SELECT (.+) FROM todos").WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Show)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAdd(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	mock.ExpectExec("INSERT INTO todos").
		WithArgs("Test todo").
		WillReturnResult(sqlmock.NewResult(1, 1))

	form := strings.NewReader("item=Test todo")
	req, err := http.NewRequest("POST", "/add", form)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Add)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMovedPermanently)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM todos").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("GET", "/delete/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Delete)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMovedPermanently)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestComplete(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	mock.ExpectExec("UPDATE todos SET completed").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("GET", "/complete/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Complete)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMovedPermanently)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}