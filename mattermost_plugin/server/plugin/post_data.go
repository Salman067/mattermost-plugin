package plugin

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type message struct {
	UserID      string
	Username    string
	ChannelID   string
	ChannelName string
	TeamID      string
	TeamName    string
	Message     string
	Timestamp   int64
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	defer func() {
		if r := recover(); r != nil {
			p.API.LogError("message post panic", "error", r)
		}
	}()

	user, err := p.API.GetUser(post.UserId)
	if err != nil {
		p.API.LogError("Failed to get user", "user_id", post.UserId, "err", err.Error())
		return
	}

	channel, err := p.API.GetChannel(post.ChannelId)
	if err != nil {
		p.API.LogError("Failed to get channel", "channel_id", post.ChannelId, "err", err.Error())
		return
	}
	team, err := p.API.GetTeam(channel.TeamId)
	if err != nil {
		p.API.LogError("Failed to get team", "team_id", channel.TeamId, "err", err.Error())
		return
	}

	db, dbErr := p.ConnectDB()
	if dbErr != nil {
		p.API.LogError("Failed to connect to database", "err", dbErr.Error())
		return
	}
	defer db.Close()

	if strings.Contains(strings.ToLower(post.Message), "wfh") {
		reqBody := &message{
			UserID:      user.Id,
			Username:    user.Username,
			ChannelID:   channel.Id,
			ChannelName: channel.Name,
			TeamID:      team.Id,
			TeamName:    team.Name,
			Message:     post.Message,
			Timestamp:   post.CreateAt,
		}

		_, errDB := db.Exec(`
        INSERT INTO plugins (user_id, username, channel_id, channel_name, team_id, team_name, message, timestamp)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			reqBody.UserID, reqBody.Username, reqBody.ChannelID, reqBody.ChannelName, reqBody.TeamID, reqBody.TeamName, reqBody.Message, reqBody.Timestamp)

		if errDB != nil {
			p.API.LogError("Failed to insert message into database", "err", errDB.Error())
		}
	}
}

// only for testing purposes
func (p *Plugin) Post(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			p.API.LogError("post panic", "error", r)
		}
	}()
	reqBody := &message{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Failed to bind request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	db, dbErr := p.ConnectDB()
	if dbErr != nil {
		p.API.LogError("Failed to connect to database", "err", dbErr.Error())
		http.Error(w, "Failed to connect to database: "+dbErr.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, errDB := db.Exec(`
        INSERT INTO plugins (user_id, username, channel_id, channel_name, team_id, team_name, message, timestamp)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		reqBody.UserID, reqBody.Username, reqBody.ChannelID, reqBody.ChannelName, reqBody.TeamID, reqBody.TeamName, reqBody.Message, reqBody.Timestamp)

	if errDB != nil {
		p.API.LogError("Failed to insert message into database", "err", errDB.Error())
		http.Error(w, "Failed to insert message into database: "+errDB.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Message inserted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		p.API.LogError("Failed to encode response to JSON", "err", err.Error())
		http.Error(w, "Failed to encode response to JSON: "+err.Error(), http.StatusInternalServerError)
	}
}
