package plugin

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/robfig/cron/v3"
)

type Plugin struct {
	plugin.MattermostPlugin
	Cron *cron.Cron
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/hello", p.HandleHello).Methods(http.MethodGet)

	router.HandleFunc("/list", p.HandleGetMessages).Methods(http.MethodGet)

	// only for testing purposes
	router.HandleFunc("/post", p.Post).Methods(http.MethodPost)

	router.ServeHTTP(w, r)
}

func (p *Plugin) ConnectDB() (*sql.DB, error) {
	defer func() {
		if r := recover(); r != nil {
			p.API.LogError("db connection panic", "error", r)
		}
	}()
	config := p.API.GetConfig()
	dbSettings := config.SqlSettings

	dns := "host=postgres port=5432 user=mmuser password=mmuser_password dbname=mattermost sslmode=disable"
	// dns := "host=postgres port=5432 user=wfhuser password=xZbacDFWm4hxOUDn dbname=wfhplugindb sslmode=disable"

	dbSettings.DataSource = &dns
	db, err := sql.Open(*dbSettings.DriverName, *dbSettings.DataSource)
	if err != nil {
		return nil, fmt.Errorf("connection error : %w", err)
	}
	db.SetMaxIdleConns(*dbSettings.MaxIdleConns)
	db.SetMaxOpenConns(*dbSettings.MaxOpenConns)
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS plugins (
		id SERIAL PRIMARY KEY,
		user_id TEXT,
		username TEXT,
		channel_id TEXT,
		channel_name TEXT,
		team_id TEXT,
		team_name TEXT,
		message TEXT,
		timestamp BIGINT
	)
`)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugins table: %w", err)
	}
	return db, nil
}
