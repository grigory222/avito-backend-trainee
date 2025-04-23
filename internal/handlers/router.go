package handlers

import (
	"github.com/grigory222/avito-backend-trainee/config"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/middlewares"
	myjwt "github.com/grigory222/avito-backend-trainee/internal/jwt"
	"github.com/grigory222/avito-backend-trainee/pkg/logger"
	"net/http"
	"strconv"
)

func NewServer(cfg *config.Config, log *logger.Logger, provider *myjwt.Provider, ps ProductService, pvzs PVZService, rs ReceptionService, us UserService) *http.Server {
	mux := http.NewServeMux()

	ph := NewProductHandler(ps, log)
	rh := NewReceptionHandler(rs, log)
	pvzh := NewPVZHandler(pvzs, log)

	employeeHandler := middlewares.AuthorizeMiddleware(dto.RoleEmployee)
	moderatorHandler := middlewares.AuthorizeMiddleware(dto.RoleModerator)
	postHandler := middlewares.MethodMiddleware(http.MethodPost)
	getHandler := middlewares.MethodMiddleware(http.MethodGet)

	mux.Handle("/hello", getHandler(employeeHandler(http.HandlerFunc(ph.Hello))))

	productHandlerChain := postHandler(
		moderatorHandler(
			middlewares.DecoderMiddleware[dto.AddProductRequestDto]()(
				http.HandlerFunc(ph.AddProductToReception),
			),
		),
	)
	mux.Handle("/products", productHandlerChain)

	receptionHandlerChain := postHandler(
		employeeHandler(
			middlewares.DecoderMiddleware[dto.CreateReceptionRequestDto]()(
				http.HandlerFunc(rh.AddReception),
			),
		),
	)
	mux.Handle("/receptions", receptionHandlerChain)

	pvzHandlerChain := postHandler(
		moderatorHandler(
			middlewares.DecoderMiddleware[dto.PVZDto]()(
				http.HandlerFunc(pvzh.AddPVZ)),
		),
	)
	mux.Handle("/pvz", pvzHandlerChain)

	// AuthenticateMiddleware на всё (исключения внутри)
	handler := middlewares.AuthenticateMiddleware(provider)(mux)

	srv := &http.Server{
		Addr:    cfg.HTTP.Host + ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler: handler,
	}
	return srv
}
