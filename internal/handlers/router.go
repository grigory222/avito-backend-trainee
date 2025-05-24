package handlers

import (
	"github.com/grigory222/avito-backend-trainee/config"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/dto"
	"github.com/grigory222/avito-backend-trainee/internal/handlers/middlewares"
	"github.com/grigory222/avito-backend-trainee/internal/myjwt"
	"github.com/grigory222/avito-backend-trainee/pkg/logger"
	"net/http"
	"strconv"
	"strings"
)

func NewServer(cfg *config.Config, log *logger.Logger, provider *myjwt.Provider, ps ProductService, pvzs PVZService, rs ReceptionService, us UserService) *http.Server {
	mux := http.NewServeMux()

	ph := NewProductHandler(ps, log)
	rh := NewReceptionHandler(rs, log)
	pvzh := NewPVZHandler(pvzs, log)
	uh := NewUserHandler(us, provider, log)

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

	pvzHandler := moderatorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			moderatorHandler(
				middlewares.DecoderMiddleware[dto.PVZDto]()(
					http.HandlerFunc(pvzh.AddPVZ),
				),
			).ServeHTTP(w, r)
		case http.MethodGet:
			moderatorHandler(http.HandlerFunc(pvzh.GetPVZ)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.Handle("/pvz", pvzHandler)

	// /pvz/{pvzId}/
	mux.Handle("/pvz/", postHandler(moderatorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.Method == http.MethodPost && strings.HasSuffix(path, "/close_last_reception") {
			pvzh.CloseLastReception(w, r)
			return
		}
		if r.Method == http.MethodPost && strings.HasSuffix(path, "/delete_last_product") {
			pvzh.DeleteLastProduct(w, r)
			return
		}
		http.NotFound(w, r)
	}))))

	mux.Handle("/dummyLogin", postHandler(middlewares.DecoderMiddleware[dto.RoleStruct]()(http.HandlerFunc(uh.DummyLogin))))
	mux.Handle("/register", postHandler(middlewares.DecoderMiddleware[dto.UserCreateRequest]()(http.HandlerFunc(uh.Register))))
	mux.Handle("/login", postHandler(middlewares.DecoderMiddleware[dto.UserLoginRequest]()(http.HandlerFunc(uh.Login))))

	// AuthenticateMiddleware на всё (исключения внутри)
	handler := middlewares.AuthenticateMiddleware(provider)(mux)

	srv := &http.Server{
		Addr:    cfg.HTTP.Host + ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler: handler,
	}
	return srv
}
