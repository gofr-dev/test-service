package main

import (
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.GET("/", func(ctx *gofr.Context) (interface{}, error) {
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
