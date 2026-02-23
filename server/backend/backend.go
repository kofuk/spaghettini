package backend

import "github.com/kofuk/spaghettini/server/types"

type Backend interface {
	Handle(request *types.Request) ([]byte, error)
}
