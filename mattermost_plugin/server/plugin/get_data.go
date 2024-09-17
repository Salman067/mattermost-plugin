package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (p *Plugin) HandleHello(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			p.API.LogError("handle hello panicked", "error", r)
		}
	}()
	_, err := w.Write([]byte("Hello, I am hello route !!!!!!"))
	if err != nil {
		p.API.LogError("Failed to write hello message", "err", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			p.API.LogError("handle get message  panicked", "error", r)
		}
	}()

	pageParam := r.URL.Query().Get("page")
	perPageParam := r.URL.Query().Get("per_page")
	fromDateParam := r.URL.Query().Get("from_date")
	toDateParam := r.URL.Query().Get("to_date")

	page := 1
	perPage := 10

	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}
	if perPageParam != "" {
		if pp, err := strconv.Atoi(perPageParam); err == nil && pp > 0 {
			perPage = pp
		}
	}

	query := `
           SELECT user_id, username, channel_id, channel_name, team_id, team_name, message, timestamp
           FROM plugins
            WHERE 1=1
         `
	countQuery := `
        SELECT COUNT(*) 
        FROM plugins 
        WHERE 1=1
`

	var args []interface{}
	var countArgs []interface{}
	argIndex := 1

	configuration := p.API.GetPluginConfig()
	for key, value := range configuration {
		channelName, ok := value.(string)
		if key == "channelname" && channelName != "" {
			if ok {
				query += fmt.Sprintf(" AND channel_name = $%d", argIndex)
				countQuery += fmt.Sprintf(" AND channel_name = $%d", argIndex)
				args = append(args, channelName)
				countArgs = append(countArgs, channelName)
				argIndex++
			}
		}
	}

	if fromDateParam != "" {
		fromDate, err := time.Parse("2006-01-02", fromDateParam)
		if err != nil {
			p.API.LogError("Invalid from_date format", "from_date", fromDateParam, "err", err.Error())
			http.Error(w, "Bad Request: Invalid from_date format", http.StatusBadRequest)
			return
		}

		query += fmt.Sprintf(" AND timestamp >= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND timestamp >= $%d", argIndex)
		args = append(args, fromDate.Unix()*1000)
		countArgs = append(countArgs, fromDate.Unix()*1000)
		argIndex++
	}

	if toDateParam != "" {
		toDate, err := time.Parse("2006-01-02", toDateParam)
		if err != nil {
			p.API.LogError("Invalid to_date format", "to_date", toDateParam, "err", err.Error())
			http.Error(w, "Bad Request: Invalid to_date format", http.StatusBadRequest)
			return
		}

		if fromDateParam != "" && fromDateParam == toDateParam {
			toDate = toDate.Add(24 * time.Hour).Add(-time.Millisecond)
		}

		query += fmt.Sprintf(" AND timestamp <= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND timestamp <= $%d", argIndex)
		args = append(args, toDate.Unix()*1000)
		countArgs = append(countArgs, toDate.Unix()*1000)
		argIndex++
	}

	if fromDateParam == "" && toDateParam == "" {
		todayStart := time.Now().Truncate(24 * time.Hour).UTC()
		todayEnd := todayStart.Add(24 * time.Hour).Add(-time.Millisecond).UTC()

		query += fmt.Sprintf(" AND timestamp >= $%d AND timestamp <= $%d", argIndex, argIndex+1)
		countQuery += fmt.Sprintf(" AND timestamp >= $%d AND timestamp <= $%d", argIndex, argIndex+1)
		args = append(args, todayStart.Unix()*1000, todayEnd.Unix()*1000)
		countArgs = append(countArgs, todayStart.Unix()*1000, todayEnd.Unix()*1000)
		argIndex += 2
	}

	db, dbErr := p.ConnectDB()
	if dbErr != nil {
		p.API.LogError("Failed to connect to the database 11", "dbErr", dbErr.Error())
		http.Error(w, "Internal Server Error 1 "+dbErr.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var totalCount int
	err := db.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		p.API.LogError("Failed to query total message count", "err", err.Error())
		http.Error(w, "Internal Server Error 2 "+err.Error(), http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * perPage
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, perPage, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		p.API.LogError("Failed to query messages", "err", err.Error())
		http.Error(w, "Internal Server Error 3 "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var userID, username, channelID, channelName, teamID, teamName, message string
		var timestamp int64

		if err := rows.Scan(&userID, &username, &channelID, &channelName, &teamID, &teamName, &message, &timestamp); err != nil {
			p.API.LogError("Failed to scan message row", "err", err.Error())
			continue
		}

		messages = append(messages, map[string]interface{}{
			"user_id":      userID,
			"username":     username,
			"channel_id":   channelID,
			"channel_name": channelName,
			"team_id":      teamID,
			"team_name":    teamName,
			"message":      message,
			"timestamp":    timestamp,
		})
	}

	response := map[string]interface{}{
		"messages":    messages,
		"total_count": totalCount,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		p.API.LogError("Failed to encode messages to JSON", "err", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
