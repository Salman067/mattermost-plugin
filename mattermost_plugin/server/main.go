package main

import (
	mattermostPlugin "github.com/example/my-plugin/server/plugin"
	"github.com/mattermost/mattermost/server/public/plugin"
)

func main() {
	// config.SetConfig()
	plugin.ClientMain(&mattermostPlugin.Plugin{})
}
