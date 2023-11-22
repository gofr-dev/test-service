package main

import (
	"context"
	"gofr.dev/pkg/gofr"
	gofrLogger "gofr.dev/pkg/log"
	"gofr.dev/pkg/service"
	"time"
)

func main() {

	app := gofr.New()

	// Push some gofr repo related metrics to see its progress over time.
	_ = app.NewGauge("gofr_repo_stargazers", "number of stargazers of gofr repo")
	_ = app.NewGauge("gofr_repo_subscribers", "number of subscribers of gofr repo")
	_ = app.NewGauge("gofr_repo_forks", "number of forks of gofr repo")
	// Start a go routine to fill in the counters every 5 min
	go func() {
		for {
			data := getGofrRepoStats(app.Logger)
			app.Logger.Debugf("Got stargazer count: %d, subscribers: %d, Forks: %d",
				data.StargazersCount, data.SubscribersCount, data.ForksCount)
			err := app.Metric.SetGauge("gofr_repo_stargazers", float64(data.StargazersCount))
			if err != nil {
				app.Logger.Error(err)
			}
			err = app.Metric.SetGauge("gofr_repo_subscribers", float64(data.SubscribersCount))
			if err != nil {
				app.Logger.Error(err)
			}
			err = app.Metric.SetGauge("gofr_repo_forks", float64(data.ForksCount))
			if err != nil {
				app.Logger.Error(err)
			}
			time.Sleep(5 * time.Minute)
		}
	}()

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

	app.Start()
}

type GithubStats struct {
	StargazersCount  int `json:"stargazers_count"`
	SubscribersCount int `json:"subscribers_count"`
	ForksCount       int `json:"forks"`
}

func getGofrRepoStats(logger gofrLogger.Logger) (data GithubStats) {
	svc := service.NewHTTPServiceWithOptions("https://api.github.com/", logger,
		&service.Options{
			SurgeProtectorOption: &service.SurgeProtectorOption{
				Disable: true,
			},
		})
	res, err := svc.Get(context.Background(), "/repos/gofr-dev/gofr", nil)
	if err != nil {
		return
	}

	err = svc.Bind(res.Body, &data)

	if err != nil {
		return
	}

	return data
}
