package server

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

func CompaniesRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/{companyID}/", func(r chi.Router) {
		r.Use(CompanyMDLWR)
		r.Get("/", GetCompanyData)
		r.Put("/", UpdateCompany)
		r.Delete("/", DeleteCompany)
	})
	r.Get("/", GetCompaniesList)
	r.Post("/", CreateCompany)
	return r
}


func CompanyMDLWR(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var company *Company
		var err error
		var companyIdInt int
		companyID := chi.URLParam(r, "companyID")
		if companyID != "" {
			companyIdInt, err = strconv.Atoi(companyID)
			company, err = dbGetCompany(companyIdInt)
		} else {
			render.Render(w, r, &DefaultResponse{HTTPStatusCode: 404, StatusText: "Resource not found."})
			return
		}
		if err != nil {
			render.Render(w, r, &DefaultResponse{HTTPStatusCode: 404, StatusText: "Resource not found."})
			return
		}

		ctx := context.WithValue(r.Context(), "company", company)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

