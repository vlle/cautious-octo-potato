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
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	Layout              = "2006-01-02"
	Driver              = "sqlite3"
	MsgMethodNotAllowed = `{"error": "method not allowed"}`
	MsgBadRequest       = `{"error": "bad request"}`
)

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

// type update_event struct {
// 	cache *Cache
// }
//
// func (s *update_event) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	response_method := r.Method
// 	if response_method != "POST" {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		w.Write([]byte(`{"message": "method not allowed"}`))
// 		return
// 	}
//   dateStr := parse_http_json(r)
//   date, _ := time.Parse(layout, dateStr)
//   s.cache.WriteToCache(dateStr, date)
//   s.cache.SaveEventForWeek(date)
//   w.WriteHeader(http.StatusOK)
//   w.Write([]byte(`{"status": "updated"}`))
// }
//
// type delete_event struct {
// 	cache *Cache
// }
// type events_for_day struct {
// 	cache *Cache
// }
//
// func (s *events_for_day) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	response_method := r.Method
// 	if response_method != "GET" {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		w.Write([]byte(`{"message": "method not allowed"}`))
// 		return
// 	}
// 	var sb strings.Builder
// 	query := r.URL.Query()
// 	msg := query.Get("date")
// 	events := s.cache.GetFromCache(msg)
// 	if events == "" {
// 		w.WriteHeader(http.StatusNotFound)
// 	} else {
// 		w.WriteHeader(http.StatusOK)
// 	}
// 	sb.WriteString(`{"events": [`)
// 	sb.WriteString(events)
// 	sb.WriteString(`]}`)
// 	w.Write([]byte(sb.String()))
// }
//
// type events_for_week struct {
// 	cache *Cache
// }
//
// func (s *events_for_week) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	response_method := r.Method
// 	if response_method != "GET" {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		w.Write([]byte(`{"message": "method not allowed"}`))
// 		return
// 	}
// 	query := r.URL.Query()
// 	date := query.Get("date")
// 	parsed_date, _ := time.Parse(layout, date)
// 	var sb strings.Builder
// 	events := s.cache.EventsForWeek(parsed_date)
// 	if events == "" {
// 		w.WriteHeader(http.StatusNotFound)
// 		sb.WriteString(`{"events": []}`)
// 	} else {
// 		w.WriteHeader(http.StatusOK)
// 		sb.WriteString(`{"events": [`)
// 		sb.WriteString(events)
// 		sb.WriteString(`]}`)
// 		w.Write([]byte(sb.String()))
// 	}
// }
//
// type events_for_month struct{}

func main() {
	// cache := NewCache()
	db_name := os.Getenv("DB_NAME")
	if db_name == "" {
		db_name = "events.db"
	}
	db, err := sql.Open(Driver, db_name)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db_struct := &DB{Db: db}

	s := &create_event{db: db_struct}
	http.Handle("/create_event", s)
	// u := &update_event{cache: cache}
	// http.Handle("/update_event", u)
	// e_day := &events_for_day{cache: cache}
	// http.Handle("/events_for_day", e_day)
	// e_week := &events_for_week{cache: cache}
	// http.Handle("/events_for_week", e_week)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
