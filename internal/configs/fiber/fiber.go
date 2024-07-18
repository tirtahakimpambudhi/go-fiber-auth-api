package fiber

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"syscall"
)

// Fiber struct represents a custom struct that embeds fiber.App
type Fiber struct {
	App    *fiber.App   // Fiber struct embedding fiber.App
	Config *FiberConfig // Config field of type *FiberConfig
}

// NewFiber used for initialize Fiber
func NewFiber(config *FiberConfig) *Fiber {
	app := fiber.New(config.ToFiberAppConfig())
	return &Fiber{App: app, Config: config}
}

// Serve method starts serving requests without TLS
func (appFiber *Fiber) Serve() error {
	// Serve method to start serving requests
	// without tls
	if appFiber.Config.SSL.CertFile == "" || appFiber.Config.SSL.KeyFile == "" {
		return appFiber.App.Listen(fmt.Sprintf("%s:%s", appFiber.Config.Host, appFiber.Config.Port))
	}
	return appFiber.App.ListenTLS(fmt.Sprintf("%s:%s", appFiber.Config.Host, appFiber.Config.Port), appFiber.Config.SSL.CertFile, appFiber.Config.SSL.KeyFile)
}

// ServeWithGraceful method starts serving requests gracefully
func (appFiber *Fiber) ServeWithGraceful() error {
	// ServeWithGraceful method to start serving requests gracefully
	errChan := make(chan error, 1)                       // errChan for handling errors
	exit := make(chan os.Signal, 1)                      // exit channel for signals
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM) // Notify for signals

	go func() {
		defer close(errChan)
		address := fmt.Sprintf("%s:%s", appFiber.Config.Host, appFiber.Config.Port) // address to listen on
		var err error

		if appFiber.Config.SSL.CertFile == "" || appFiber.Config.SSL.KeyFile == "" {
			err = appFiber.App.Listen(address)
		} else {
			err = appFiber.App.ListenTLS(address, appFiber.Config.SSL.CertFile, appFiber.Config.SSL.KeyFile)
		}

		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-exit:
		close(exit)
		return appFiber.App.Shutdown()
	case err := <-errChan:
		return err
	}
}
