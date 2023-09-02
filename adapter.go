package ws

import (
	"github.com/pawelgaczynski/gain"
)

type adapter struct {
}

func (a *adapter) StartAsMainProcess(address string) error {
	panic("you should not use this method,please do this in your application")
}

func (a *adapter) Start(address string) error {
	//TODO implement me
	panic("implement me")
}

func (a *adapter) Shutdown() {

}

func (a *adapter) AsyncShutdown() {
	//TODO implement me
	panic("implement me")
}

func (a *adapter) ActiveConnections() int {
	//TODO implement me
	panic("implement me")
}

func (a *adapter) IsRunning() bool {
	//TODO implement me
	panic("implement me")
}

var _ gain.Server = (*adapter)(nil)
