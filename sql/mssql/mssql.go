package mssql

import (
	"github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/api/types/pipeline"
)

type Publisher struct{}

// NewPublisher creates a new MSSQL publisher instance
func NewPublisher() publisher.Publisher {
	return &Publisher{}
}

func (p *Publisher) Shapes(ctx publisher.Context) (map[string]pipeline.Shape, error) {
	return nil, nil
}

func (p *Publisher) Publish(ctx publisher.Context, dataTransport publisher.DataTransport) {
}
