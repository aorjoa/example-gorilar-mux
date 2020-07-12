package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"thaichana/logger"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	viper.SetDefault("port", "8000")
	viper.SetDefault("db.conn", "thaichana.db")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	l, _ := zap.NewDevelopment()
	defer l.Sync()

	hostname, _ := os.Hostname()
	l = l.With(zap.String("hostname", hostname))
	zap.ReplaceGlobals(l)

	db, err := sql.Open("sqlite3", viper.GetString("db.conn"))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	r := mux.NewRouter()

	r.Use(logger.MiddleWare(l))
	r.Use(SealMiddleWare())

	// This will serve files under http://localhost:8000/static/<filename>
	r.HandleFunc("/recently", Recently).Methods(http.MethodPost)
	r.HandleFunc("/checkin", CheckIn(NewInsertCheckin(db))).Methods(http.MethodPost)
	r.HandleFunc("/checkout", CheckOut).Methods(http.MethodPost)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	zap.L().Info("start...", zap.String("port", viper.GetString("port")))

	log.Fatal(srv.ListenAndServe())
}

type Check struct {
	ID      int64
	PlaceID int64
}

type Location struct {
	Lat  float64
	Long float64
}

func NewInsertCheckin(db *sql.DB) InFunc {
	return func(ID, placeID int64) (err error) {
		_, err = db.Exec("INSERT INTO visits VALUES(?, ?);", ID, placeID)
		return
	}
}

type InFunc func(id, placeID int64) error

func (fn InFunc) In(id, placeID int64) error {
	return fn(id, placeID)
}

type Iner interface {
	In(ID, plateID int64) error
}

// Recently returns currently visited
func Recently(w http.ResponseWriter, r *http.Request) {
}

// CheckIn check-in to place, returns density (ok, too much)
func CheckIn(check Iner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var chk Check

		if err := json.NewDecoder(r.Body).Decode(&chk); err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(err)
			return
		}
		defer r.Body.Close()

		if err := check.In(chk.ID, chk.PlaceID); err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(err)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "ok",
		})

		// w.WriteHeader(200)
		// json.NewEncoder(w).Encode(chk)
	}
}

// CheckOut check-out from place
func CheckOut(w http.ResponseWriter, r *http.Request) {

}

func SealMiddleWare() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
			}

			data, err := base64.StdEncoding.DecodeString(string(b))
			if err != nil {
				return
			}

			buff := bytes.NewBuffer(data)
			if err != nil {
				return
			}

			r.Body = ioutil.NopCloser(buff)

			next.ServeHTTP(&EncoderWriter{w}, r)

		})
	}
}

type EncoderWriter struct {
	http.ResponseWriter
}

func (w *EncoderWriter) Write(b []byte) (int, error) {
	str := base64.StdEncoding.EncodeToString(b)
	return w.ResponseWriter.Write([]byte(str))
}
