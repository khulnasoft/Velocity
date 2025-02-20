package modules

import (
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
	"go.khulnasoft.com/velocity"
)

var _ types.ModuleFunc = ConfigureVelocity

// ConfigureVelocity Configures all services required by the Velocity app.
func ConfigureVelocity(registry types.ServiceRegistry) error {
	registration.RegisterInstance(registry, velocity.Config{
		AppName:   "parsley-app-recipe",
		Immutable: true,
	})

	registry.Register(newVelocity, types.LifetimeSingleton)
	registry.RegisterModule(RegisterRouteHandlers)

	return nil
}

// newVelocity Activator method for new Velocity instances.
func newVelocity(config velocity.Config) *velocity.App {
	return velocity.New(config)
}
