package publishers

import (
	"github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/pipeline-publishers/web/wellcast"
	"github.com/naveego/pipeline-publishers/xml/wellez"
)

func init() {
	publisher.RegisterFactory("wellcast", wellcast.NewPublisher)
	publisher.RegisterFactory("wellez", wellez.NewPublisher)
}
