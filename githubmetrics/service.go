package githubmetrics

import (
	"context"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/service"
)

type GithubStats struct {
	StargazersCount  int `json:"stargazers_count"`
	SubscribersCount int `json:"subscribers_count"`
	ForksCount       int `json:"forks"`
}

type Service struct {
	app *gofr.Gofr
}

func New(app *gofr.Gofr) Service {
	return struct{ app *gofr.Gofr }{app: app}
}

func (s *Service) Process() {
	svc := service.NewHTTPServiceWithOptions("https://api.github.com/", s.app.Logger,
		&service.Options{
			SurgeProtectorOption: &service.SurgeProtectorOption{
				Disable: true,
			},
		})

	data := getGofrRepoStats(svc)
	s.app.Logger.Debugf("Got stargazer count: %d, subscribers: %d, Forks: %d",
		data.StargazersCount, data.SubscribersCount, data.ForksCount)

	err := s.app.Metric.SetGauge("gofr_repo_stargazers", float64(data.StargazersCount))
	if err != nil {
		s.app.Logger.Error(err)
	}

	err = s.app.Metric.SetGauge("gofr_repo_subscribers", float64(data.SubscribersCount))
	if err != nil {
		s.app.Logger.Error(err)
	}

	err = s.app.Metric.SetGauge("gofr_repo_forks", float64(data.ForksCount))
	if err != nil {
		s.app.Logger.Error(err)
	}
}

func getGofrRepoStats(svc service.HTTP) (data GithubStats) {
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
