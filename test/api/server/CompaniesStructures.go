package server

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"net/url"
)

type CompanySerializer struct {
	*Company
}

type CompanyListSerializer struct {
	ID   uint    `json:"id"`
	Name string `json:"name"`
}

func CompanyListSerializerFunc(companies []*Company) []render.Renderer {
	var list []render.Renderer
	for _, company := range companies {
		list = append(list, &CompanyListSerializer{
			Name: company.Name,
			ID: company.ID,
		})
	}
	return list
}

func (a *CompanySerializer) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a *CompanyListSerializer) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ParsePostCompany(data url.Values) (Company, error) {
	var a Company

	if data["name"] == nil {
		return a, errors.New("missing required Company fields")
	}

	if len(data["name"]) != 1 {
		return a, errors.New("multiple forms submitted")
	}

	a.Name = data["name"][0]
	return a, nil
}
