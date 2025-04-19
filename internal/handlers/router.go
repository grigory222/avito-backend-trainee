package handlers

import (
	"github.com/grigory222/avito-backend-trainee/config"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/middlewares"
	"github.com/grigory222/avito-backend-trainee/pkg/logger"
	"net/http"
	"strconv"
)

func NewServer(cfg *config.Config, log *logger.Logger, ps ProductService) *http.Server {
	mux := http.NewServeMux()

	ph := NewProductHandler(ps, log)

	// Публичные маршруты
	//mux.HandleFunc("/hello", ph.Hello)

	employeeHandler := middlewares.AuthorizeMiddleware(dto.RoleEmployee)
	//moderatorHandler := middlewares.AuthorizeMiddleware(dto.RoleModerator)
	mux.Handle("/hello", employeeHandler(http.HandlerFunc(ph.Hello)))

	// AuthenticateMiddleware на всё (исключения внутри)
	handler := middlewares.AuthenticateMiddleware("test_secret_test_secret_test_secret")(mux)

	srv := &http.Server{
		Addr:    cfg.HTTP.Host + ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler: handler,
	}
	return srv
}
