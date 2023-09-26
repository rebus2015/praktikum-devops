package agent

import (
	"context"
	"net/http"

	"github.com/rebus2015/praktikum-devops/internal/model"
)

type ExecutionFn func(ctx context.Context, args Args) error

type Result struct {
	Err        error
	Descriptor int
}

type Args struct {
	Client  *http.Client
	Metrics []model.Metrics
	Config  *Config
}

type Job struct {
	Descriptor int
	ExecFn     ExecutionFn
	Args       Args
}

func (j Job) execute(ctx context.Context) Result {
	err := j.ExecFn(ctx, j.Args)
	if err != nil {
		return Result{
			Err:        err,
			Descriptor: j.Descriptor,
		}
	}

	return Result{
		Descriptor: j.Descriptor,
	}
}
