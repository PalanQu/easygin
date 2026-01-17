package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"easygin/internal/app"
	"easygin/pkg/logging"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func initServe() error {
	rootCmd.AddCommand(serveCmd)
	return nil
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Server",
	Long:  `This command boots the web server and serves the application to the local network.`,
	Run:   runServer,
}

func runServer(_ *cobra.Command, _ []string) {
	kernel := boot()

	stopChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		kernel.ListenAndServe(stopChan)
	}()

	graceful(kernel, 30*time.Second, stopChan, &wg)
}

func graceful(
	kernel *app.Kernel,
	timeout time.Duration,
	stopChan chan struct{},
	wg *sync.WaitGroup,
) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger := logging.GetGlobalLogger()
	logger.Info("Shutdown signal received, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	close(stopChan)

	if err := kernel.Shutdown(ctx); err != nil {
		logger.Fatal("Shutdown error", zap.Error(err))
	}

	wg.Wait()
	logger.Info("Server stopped gracefully")
}
