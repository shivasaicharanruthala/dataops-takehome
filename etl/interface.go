package etl

import (
	"context"
	"github.com/shivasaicharanruthala/dataops-takehome/model"
)

type Processor interface {
	Worker(ctx context.Context, id int, results chan<- *model.Response)
	ProcessDataFromWorker(ctx context.Context, results chan *model.Response)
}

type Extract interface {
	FetchDataFromSQS() (*model.Response, error)
}

type Loader interface {
	BatchInsert(ctx context.Context, responses []*model.Response) error
	SequentialInsert(ctx context.Context, response []model.Response) error
}
