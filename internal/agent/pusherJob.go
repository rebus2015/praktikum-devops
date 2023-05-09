package agentworkerspool

import (
	"context"
	"net/http"

	"github.com/rebus2015/praktikum-devops/internal/model"
)

type ExecutionFn func(ctx context.Context, args interface{}) error

type JobDescriptor struct {
	ID int
}

type Result struct {
	Err        error
	Descriptor JobDescriptor
}

type Args struct {
	client *http.Client
	metric model.Metrics
}

type Job struct {
	Descriptor JobDescriptor
	ExecFn     ExecutionFn
	Args       interface{}
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
