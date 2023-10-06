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
	Config  *Config
	Metrics []model.Metrics
}

type Job struct {
	ExecFn     ExecutionFn
	Args       Args
	Descriptor int
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
