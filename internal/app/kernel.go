package app

import (
	"context"
	"fmt"

	"easygin/internal/routers"
	"easygin/pkg/config"
	"easygin/pkg/db"
	"easygin/pkg/ent"
	"easygin/pkg/logging"
	pkg_otel "easygin/pkg/otel"

	"net/http"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	StateStopped int = iota
	StateStarting
	StateRunning
	StateStopping
)

const (
	serveHost = "0.0.0.0"
)

type Kernel struct {
	server      *http.Server
	config      *config.Config
	dbClient    *ent.Client
	state       int
	wg          *sync.WaitGroup
	context     cancelContext
	shutdownFns []func() error
}

func New(configPath string) (*Kernel, error) {
	config := config.LoadConfig(configPath)
	tracer, shutdownTraceFunc, err := pkg_otel.NewTracer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to new tracer")
	}
	database, err := db.CreateDBClient(config, tracer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}
	ctx, cancel := context.WithCancel(context.Background())

	s := newServer(config, database)

	app := &Kernel{
		config:   config,
		dbClient: database,
		server:   s,
		state:    StateStarting,
		shutdownFns: []func() error{
			shutdownTraceFunc,
		},
		wg:      &sync.WaitGroup{},
		context: cancelContext{cancel: cancel, ctx: ctx},
	}

	app.state = StateRunning

	return app, nil
}

func newServer(
	config *config.Config,
	db *ent.Client,
) *http.Server {
	r := routers.SetupRouter(config.Server.RouterPrefix, db)
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", serveHost, config.Server.Port),
		Handler: r,
	}
}

type cancelContext struct {
	cancel context.CancelFunc
	ctx    context.Context
}

func (k *Kernel) ListenAndServe(stopChan <-chan struct{}) {
	serverErrCh := make(chan error, 1)
	logger := logging.GetGlobalLogger()

	go func() {
		description := fmt.Sprintf("starting server at %v", k.config.Server.Port)
		logger.Info(description)
		if err := k.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Warn("error in start server", zap.Error(err))
			serverErrCh <- err
		}
	}()

	select {
	case err := <-serverErrCh:
		logger.Warn("server stopped with error", zap.Error(err))
	case <-stopChan:
		logger.Info("received stop signal, shutting down server")
	}
}

func (k *Kernel) Shutdown(ctx context.Context) error {
	logger := logging.GetGlobalLogger()
	if k.state != StateRunning {
		logger.Warn("Application cannot be shutdown since current state is not 'running'")
		return nil
	}

	k.state = StateStopping
	defer func() {
		k.state = StateStopped
	}()

	if k.server != nil {
		if err := k.server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown error", zap.Error(err))
		} else {
			logger.Info("server stopped")
		}
	}

	if k.server != nil {
		if err := k.server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown error", zap.Error(err))
		} else {
			logger.Info("server stopped")
		}
	}

	k.context.cancel()
	done := make(chan struct{})
	go func() {
		k.wg.Wait()
		close(done)
	}()

	for _, fn := range k.shutdownFns {
		shutdownErr := fn()
		if shutdownErr != nil {
			logger.Error("shutdown function returned error", zap.Error(shutdownErr))
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	return k.dbClient.Close()
}
