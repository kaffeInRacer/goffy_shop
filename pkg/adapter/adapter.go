package adapter

import "context"

type Adapter interface {
	Start(context.Context) error
	Stop(context.Context) error
}
