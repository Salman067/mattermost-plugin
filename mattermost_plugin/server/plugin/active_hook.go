package plugin

import (
	"time"

	"github.com/robfig/cron/v3"
)

func (p *Plugin) OnActivate() error {
	location, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		p.API.LogError("Failed to load location", "error", err)
		return err
	}

	p.Cron = cron.New(cron.WithLocation(location))

	_, err = p.Cron.AddFunc("0 9 * * *", func() {
		p.RunCurlCommand()
	})
	if err != nil {
		p.API.LogError("Failed to schedule cron job", "error", err)
		return err
	}

	p.Cron.Start()

	return nil
}
func (p *Plugin) OnDeactivate() error {
	if p.Cron != nil {
		p.Cron.Stop()
	}
	return nil
}
