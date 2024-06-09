package etl

import (
	"context"
	"github.com/shivasaicharanruthala/dataops-takehome/model"
	"sync"
)

type Processor interface {
	Worker(ctx context.Context, id int, wg *sync.WaitGroup, results chan<- *model.Response)
	ProcessDataFromWorker(ctx context.Context, results chan *model.Response)
}

type Extract interface {
	FetchDataFromSQS() (*model.Response, error)
}

type Loader interface {
	BatchInsert(ctx context.Context, responses []*model.Response) error
}
