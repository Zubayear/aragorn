package routes

import (
	"github.com/Zubayear/aragorn/api/handlers"
	"github.com/Zubayear/aragorn/pkg/sms"
	"github.com/Zubayear/aragorn/pkg/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Router(ch *handlers.CpHandler, userService user.Service, smsService sms.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Route("/rb/ecmapigw/webresources/ecmapigw.v3", func(r chi.Router) {
		r.Post("", ch.HandleRequest(userService, smsService))
	})
	return r
}
