//go:build wireinject

package container

import "github.com/google/wire"

func initializeApplication() (*Application, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
