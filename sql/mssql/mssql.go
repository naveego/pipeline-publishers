package mssql

import (
	"github.com/naveego/api/types/pipeline"
	"github.com/naveego/pipeline-api/publisher"
)

type Publisher struct{}

func (p *Publisher) Shapes(ctx publisher.Context) (map[string]pipeline.Shape, error) {
	return nil, nil
}

func (p *Publisher) Publish(ctx publisher.Context, dataTransport publisher.DataTransport) {
}
