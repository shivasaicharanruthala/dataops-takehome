package store

import "github.com/shivasaicharanruthala/dataops-takehome/model"

type Login interface {
	Get(filter *model.Filter) ([]model.Response, error)
}
