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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	layout = "2006-01-02"
)

type event struct {
	Date string `json:"date"`
}

func IsoWeek(t time.Time) (int, int) {
	isoYear, isoWeek := t.ISOWeek()
	return isoYear, isoWeek
}

func StringIsoWeek(t time.Time) string {
	isoYear, isoWeek := t.ISOWeek()
	return fmt.Sprintf("%s%s", fmt.Sprint(isoYear), fmt.Sprint(isoWeek))
}

type Cache struct {
	cache map[string]string
	mutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]string),
	}
}

func (c *Cache) WriteToCache(key string, value time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[key] = value.String()
}

func (c *Cache) SaveEventForWeek(value time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	key := StringIsoWeek(value)
	c.cache[key] += value.String() + ","

	fmt.Println(c.cache[key])
	fmt.Println(key)
}

func (c *Cache) GetFromCache(key string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.cache[key]
}

func (c *Cache) EventsForWeek(value time.Time) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	key := StringIsoWeek(value)
	fmt.Println(key)
	return c.cache[key]
}

func (c *Cache) EventsForMonth(month string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.cache[month]
}

type create_event struct {
	cache *Cache
}

func parse_http_json(r *http.Request) string {
	var event event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Fatal(err)
	}
	return event.Date
}

func (s *create_event) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// user_id=3&date=2019-09-09
	w.Header().Set("Content-Type", "application/json")
	log.Println(r.Body)
	response_method := r.Method
	if response_method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"message": "method not allowed"}`))
		return
	}
	layout := "2006-01-02"
	dateStr := parse_http_json(r)
	date, _ := time.Parse(layout, dateStr)
	s.cache.WriteToCache(dateStr, date)
	s.cache.SaveEventForWeek(date)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "created"}`))
}

type update_event struct {
	cache *Cache
}

func (s *update_event) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response_method := r.Method
	if response_method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"message": "method not allowed"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}

type delete_event struct {
	cache *Cache
}
type events_for_day struct {
	cache *Cache
}

func (s *events_for_day) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response_method := r.Method
	if response_method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"message": "method not allowed"}`))
		return
	}
	var sb strings.Builder
	query := r.URL.Query()
	log.Println(query)
	msg := query.Get("date")
	log.Println(msg)
	events := s.cache.GetFromCache(msg)
	if events == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	sb.WriteString(`{"events": [`)
	sb.WriteString(events)
	sb.WriteString(`]}`)
	w.Write([]byte(sb.String()))
}

type events_for_week struct {
	cache *Cache
}

func (s *events_for_week) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response_method := r.Method
	if response_method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"message": "method not allowed"}`))
		return
	}
	query := r.URL.Query()
	date := query.Get("date")
	parsed_date, _ := time.Parse(layout, date)
	fmt.Println(parsed_date)
	var sb strings.Builder
	events := s.cache.EventsForWeek(parsed_date)
	if events == "" {
		w.WriteHeader(http.StatusNotFound)
		sb.WriteString(`{"events": []}`)
	} else {
		w.WriteHeader(http.StatusOK)
		sb.WriteString(`{"events": [`)
		sb.WriteString(events)
		sb.WriteString(`]}`)
		w.Write([]byte(sb.String()))
	}
}

type events_for_month struct{}

func main() {
	cache := NewCache()
	s := &create_event{cache: cache}
	http.Handle("/create_event", s)
	u := &update_event{cache: cache}
	http.Handle("/update_event", u)
	e_day := &events_for_day{cache: cache}
	http.Handle("/events_for_day", e_day)
	e_week := &events_for_week{cache: cache}
	http.Handle("/events_for_week", e_week)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
