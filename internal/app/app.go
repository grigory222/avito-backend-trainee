package app

import (
	"context"
	"errors"
	"github.com/grigory222/avito-backend-trainee/config"
	"github.com/grigory222/avito-backend-trainee/internal/db"
	"github.com/grigory222/avito-backend-trainee/internal/handlers"
	"github.com/grigory222/avito-backend-trainee/internal/repo"
	"github.com/grigory222/avito-backend-trainee/internal/usecases"
	"github.com/grigory222/avito-backend-trainee/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	pg, err := db.CreateConnection(cfg.DB)
	if err != nil {
		l.Fatal(err)
	}
	defer func() {
		err = pg.Close()
		if err != nil {
			l.Fatal(err)
		}
	}()

	productRepo := repo.NewProductRepository(pg)
	userRepo := repo.NewUserRepository(pg)
	pvzRepo := repo.NewPVZRepository(pg)
	receptionRepo := repo.NewReceptionRepository(pg)

	l.Debug("Repositories created: ", productRepo, userRepo, pvzRepo, receptionRepo)

	// создать сервисы
	productService := usecases.NewProductService(productRepo, receptionRepo, pvzRepo)

	srv := handlers.NewServer(cfg, l, productService)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Fatal("Listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	l.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Fatal("Shutdown error: %v", err)
	}

	l.Info("Server exited gracefully")

}
