package cmd

import (
	pubapi "github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/navigator-go/publishers/protocol"
)

type publisher struct {
	client *mssqlClient
}

func (p *publisher) Publish(request protocol.PublishRequest, transport pubapi.DataTransport) {
	panic("not implemented")
}

func (p *publisher) TestConnection(request protocol.TestConnectionRequest) (protocol.TestConnectionResponse, error) {
	panic("not implemented")
}

func (p *publisher) DiscoverShapes(request protocol.DiscoverShapesRequest) (protocol.DiscoverShapesResponse, error) {
	panic("not implemented")
}
