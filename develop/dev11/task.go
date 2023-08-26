package main

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API:
POST /create_event
POST /update_event
POST /delete_event

GET /events_for_day
GET /events_for_week
GET /events_for_month

Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	Layout              = "2006-01-02"
	Driver              = "sqlite3"
	MsgMethodNotAllowed = `{"error": "method not allowed"}`
	MsgBadRequest       = `{"error": "bad request"}`
)


func main() {
	db_name := os.Getenv("DB_NAME")
	if db_name == "" {
		db_name = "events.db"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port
	db, err := sql.Open(Driver, db_name)
	if err != nil {
		log.Fatal(err)
	}
  _, err = db.Exec("CREATE TABLE IF NOT EXISTS events(id INTEGER PRIMARY KEY AUTOINCREMENT, day INTEGER, month INTEGER, week INTEGER, year INTEGER, name TEXT)")
  if err != nil {
    log.Fatal(err)
  }
	defer db.Close()
	db_struct := &DB{Db: db}

	s := &create_event{db: db_struct}
	http.Handle("/create_event", s)
	u := &update_event{db: db_struct}
	http.Handle("/update_event", u)
	d := &delete_event{db: db_struct}
	http.Handle("/delete_event", d)
	e_day := &events_for_day{db: db_struct}
	http.Handle("/events_for_day", e_day)
	e_week := &events_for_week{db: db_struct}
	http.Handle("/events_for_week", e_week)
	e_month := &events_for_month{db: db_struct}
	http.Handle("/events_for_month", e_month)
	log.Fatal(http.ListenAndServe(port, nil))
}

type DB struct {
	Db *sql.DB
}

func ParseTime(date string) (time.Time, error) {
	return time.Parse(Layout, date)
}

func (d *DB) CreateEvent(date time.Time, name string) error {
	day := date.Day()
	month := int(date.Month())
	year, week := date.ISOWeek()
	stmt, err := d.Db.Prepare("INSERT INTO events(day, month, week, year, name) values(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(day, month, week, year, name)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) UpdateByID(id int, name string) error {
	stmt, err := d.Db.Prepare("UPDATE events SET name = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) DeleteByID(id int) error {
	stmt, err := d.Db.Prepare("DELETE FROM events WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) GetByID(id int) (string, error) {
	stmt, err := d.Db.Prepare("SELECT name FROM events WHERE id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	name := ""
	err = stmt.QueryRow(id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (d *DB) GetByDate(date time.Time) ([]string, error) {
	stmt, err := d.Db.Prepare("SELECT name FROM events WHERE day = ? AND month = ? AND year = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	name := ""
	rows, err := stmt.Query(date.Day(), int(date.Month()), date.Year())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

func (d *DB) GetByWeek(date time.Time) ([]string, error) {
	stmt, err := d.Db.Prepare("SELECT name FROM events WHERE week = ? AND year = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	name := ""
	year, week := date.ISOWeek()
	rows, err := stmt.Query(week, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

func (d *DB) GetByMonth(date time.Time) ([]string, error) {
	stmt, err := d.Db.Prepare("SELECT name FROM events WHERE month = ? AND year = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	name := ""
	rows, err := stmt.Query(int(date.Month()), date.Year())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

type event_json struct {
	Date string `json:"date"`
	Name string `json:"name"`
}

type request_answer struct {
	Result string `json:"result"`
}

type request_error struct {
	Error string `json:"error"`
}

type create_event struct {
	db *DB
}

func (s *create_event) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encoder.Encode(request_error{Error: "method not allowed"})
		return
	}
	decoder := json.NewDecoder(r.Body)
	event := event_json{}
	err := decoder.Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}
	if event.Date == "" || event.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}
	w.WriteHeader(http.StatusOK)
	time, err := ParseTime(event.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}
	err = s.db.CreateEvent(time, event.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(request_error{Error: "internal server error"})
		return
	}
	encoder.Encode(request_answer{Result: "created"})
	log.Println("Create event: ", event)
}

type update_event struct {
	db *DB
}

func (s *update_event) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encoder.Encode(request_error{Error: "method not allowed"})
		return
	}
	decoder := json.NewDecoder(r.Body)
	idus := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{}
	err := decoder.Decode(&idus)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error 1: ", err)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}

	int_id, err := strconv.Atoi(idus.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error 2: ", err)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}
	if int_id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error 2: ", err)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}

	err = s.db.UpdateByID(int_id, idus.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(request_error{Error: "internal server error"})
		return
	}

	encoder.Encode(request_answer{Result: "updated"})
	log.Println("Update event: ", idus)

}

type delete_event struct {
	db *DB
}

func (s *delete_event) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encoder.Encode(request_error{Error: "method not allowed"})
		return
	}
	decoder := json.NewDecoder(r.Body)
	idus := struct {
		ID string `json:"id"`
	}{}
	err := decoder.Decode(&idus)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error 1: ", err)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}

	int_id, err := strconv.Atoi(idus.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error 2: ", err)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}
	if int_id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error 3: ", err)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}

	err = s.db.DeleteByID(int_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(request_error{Error: "internal server error"})
		return
	}

	encoder.Encode(request_answer{Result: "deleted"})
	log.Println("Delete event: ", idus)

}

type events_for_day struct {
	db *DB
}

func (s *events_for_day) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encoder.Encode(request_error{Error: "method not allowed"})
		return
	}
	query := r.URL.Query()
	date := query.Get("date")
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}
	time, err := ParseTime(date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}
	events, err := s.db.GetByDate(time)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		encoder.Encode(request_error{Error: "internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	events_joined := strings.Join(events, ", ")
	encoder.Encode(request_answer{Result: events_joined})
  log.Println("Events for day: ", events_joined)
}

type events_for_week struct {
	db *DB
}

func (s *events_for_week) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encoder.Encode(request_error{Error: "method not allowed"})
		return
	}
	query := r.URL.Query()
	date := query.Get("date")
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}

	time, err := ParseTime(date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}
	events, err := s.db.GetByWeek(time)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		encoder.Encode(request_error{Error: "internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	events_joined := strings.Join(events, ", ")
	encoder.Encode(request_answer{Result: events_joined})
  log.Println("Events for week: ", events_joined)
}

type events_for_month struct {
	db *DB
}

func (s *events_for_month) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		encoder.Encode(request_error{Error: "method not allowed"})
		return
	}
	query := r.URL.Query()
	date := query.Get("date")
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}
	time, err := ParseTime(date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(request_error{Error: "bad request"})
		return
	}
	events, err := s.db.GetByMonth(time)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		encoder.Encode(request_error{Error: "internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	events_joined := strings.Join(events, ", ")
	encoder.Encode(request_answer{Result: events_joined})
  log.Println("Events for month: ", events_joined)
}
