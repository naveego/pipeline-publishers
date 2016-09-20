package publishers

import (
	"github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/pipeline-publishers/sql/mssql"
	"github.com/naveego/pipeline-publishers/web/wellcast"
	"github.com/naveego/pipeline-publishers/xml/wellez"
)

func init() {
	publisher.RegisterFactory("mssql", mssql.NewPublisher)
	publisher.RegisterFactory("wellcast", wellcast.NewPublisher)
	publisher.RegisterFactory("wellez", wellez.NewPublisher)
}
