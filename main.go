package main

import (
	"gofr.dev/pkg/gofr"
	"test-service/githubmetrics"
)

func main() {
	app := gofr.New()

	// Push gofr repo-related metrics to see its progress over time.
	_ = app.NewGauge("gofr_repo_stargazers", "number of stargazers of gofr repo")
	_ = app.NewGauge("gofr_repo_subscribers", "number of subscribers of gofr repo")
	_ = app.NewGauge("gofr_repo_forks", "number of forks of gofr repo")

	app.GET("", func(ctx *gofr.Context) (interface{}, error) {
		return "Hello GoFr!", nil
	})

	app.GET("/test", func(ctx *gofr.Context) (interface{}, error) {
		var resp string

		err := ctx.Redis.Get(ctx, "test").Scan(&resp)
		if err != nil {
			return nil, err
		}

		return resp, nil
	})

	app.GET("/count", func(ctx *gofr.Context) (interface{}, error) {
		var count int

		row := ctx.DB().QueryRowContext(ctx, "SELECT count(*) FROM user")

		if err := row.Scan(&count); err != nil {
			return nil, err
		}

		return count, nil
	})

	gCron := gofr.NewCron()
	svc := githubmetrics.New(app)

	go func() {
		err := gCron.AddJob("*/5 * * * *", svc.Process)
		if err != nil {
			app.Logger.Errorf("error while creating cron job for github stargazers. Error: %v", err.Error())
		}
	}()

	app.Start()
}
