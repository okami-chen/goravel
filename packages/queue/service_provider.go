package queue

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
	"github.com/hibiken/asynq"
	"goravel/packages/queue/task"
	"time"
)

const (
	Binding = "queue.v2"
)

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		db := app.MakeConfig().GetInt("database.redis.default.database")
		host := app.MakeConfig().GetString("database.redis.default.host")
		port := app.MakeConfig().GetString("database.redis.default.port")
		pass := app.MakeConfig().GetString("database.redis.default.password")

		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: host + ":" + port, DB: db, Password: pass},
			asynq.Config{
				Concurrency: 10,
				RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
					if t.Type() == task.TypeEmailDelivery {
						return 5 * time.Second
					}
					return asynq.DefaultRetryDelayFunc(n, e, t)
				},
				Queues: map[string]int{
					"critical": 6,
					"default":  3,
					"low":      1,
				},
			},
		)
		return srv, nil
	})

	app.Bind(Binding+".client", func(app foundation.Application) (any, error) {
		db := app.MakeConfig().GetInt("database.redis.default.database")
		host := app.MakeConfig().GetString("database.redis.default.host")
		port := app.MakeConfig().GetString("database.redis.default.port")
		pass := app.MakeConfig().GetString("database.redis.default.password")
		client := asynq.NewClient(asynq.RedisClientOpt{Addr: host + ":" + port, DB: db, Password: pass})
		return client, nil
	})

}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	go func() {
		server, _ := app.Make(Binding)
		defer server.(*asynq.Server).Stop()
		mux := asynq.NewServeMux()
		mux.HandleFunc(task.TypeEmailDelivery, task.NewEmailDeliveryTaskListner)
		if err := server.(*asynq.Server).Run(mux); err != nil {
			facades.Log().Error(err.Error())
		}
	}()
}
