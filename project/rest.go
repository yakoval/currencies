package project

import (
	"database/sql"
	"encoding/json"
	config2 "github.com/yakoval/currencies/config"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

// RestServer предоставляет http-эндпоинты.
type RestServer struct {
	server *http.Server
	config *config2.WebConfig
	logger *logrus.Logger
	db     *Database
}

// StartWebServer запускает REST-сервер.
func StartWebServer(logger *logrus.Logger, config *config2.WebConfig, db *Database) (*RestServer, error) {
	rs := &RestServer{
		config: config,
		logger: logger,
		db:     db,
	}
	mux := http.NewServeMux()
	server := &http.Server{Addr: config.Host, Handler: mux}

	// Отображение списка валют.
	mux.HandleFunc("/currencies/", rs.listHandler)

	// Отображение одной валюты.
	mux.HandleFunc("/currency/", rs.detailHandler)

	rs.server = server

	go func() {
		rs.logger.Info("starting http server on http://" + rs.server.Addr)
		err := rs.server.ListenAndServe()
		Check(rs.logger, err, "starting http server")
	}()

	return rs, nil
}

// Close освобождает ресурсы сервера.
func (rs *RestServer) Close() error {
	return rs.server.Close()
}

func (rs *RestServer) detailHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	Check(rs.logger, err, "getting GET-param ID")

	currency, err := rs.db.GetByID(id)
	if err == sql.ErrNoRows {
		err = rs.writeDataToResponse(ErrorResponse{Error: "No currencies by this id"}, w)
		return
	}
	Check(rs.logger, err, "getting currency by ID")

	err = rs.writeDataToResponse(currency, w)
	Check(rs.logger, err, "rest-handle currency request")
}

func (rs *RestServer) listHandler(w http.ResponseWriter, r *http.Request) {
	var page int
	var err error
	var currencies []Currency

	if rs.config.ByPages {
		if r.URL.Query().Get("page") == "" {
			page = 1
		} else {
			page, err = strconv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				err = rs.writeDataToResponse(ErrorResponse{Error: "page input error"}, w)
				return
			}
		}
		currencies, err = rs.db.GetAllForPage(page, rs.config.ItemsPerPage)
	} else {
		currencies, err = rs.db.GetAll()
	}
	Check(rs.logger, err, "getting currencies by page")

	err = rs.writeDataToResponse(currencies, w)
	Check(rs.logger, err, "output currencies")
}

func (rs *RestServer) writeDataToResponse(data interface{}, w http.ResponseWriter) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := w.Write(res); err != nil {
		return err
	}
	return nil
}
