package publishers


import (
    "github.com/naveego/pipeline-api/publisher"
    "github.com/naveego/pipeline-publishers/sql/mssql"
    "github.com/naveego/pipeline-publishers/web/wellcast"
)

func init(){
    publisher.RegisterFactory("mssql", func() publisher.Factory { return mssql.Publisher{} })
    publisher.RegisterFactory("wellcast", func() publisher.Publisher { return wellcast.Publisher{} )
}