package server

import (
	"errors"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func dbGetCompany(id int) (*Company, error) {
	company := &Company{}
	var err error
	err = db.Preload("Users").First(company, id).Error
	if err != nil {
		return company, errors.New("company matching query does not exist")
	}
	return company, err
}

func GetCompanyData(w http.ResponseWriter, r *http.Request) {
	company := r.Context().Value("company").(*Company)
	err := render.Render(w, r, &CompanySerializer{Company: company})
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func GetCompaniesList(w http.ResponseWriter, r *http.Request) {
	var companies []*Company
	db.Find(&companies)
	if err := render.RenderList(w, r, CompanyListSerializerFunc(companies)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}


func CreateCompany(w http.ResponseWriter, r *http.Request) {
	var company Company
	err := r.ParseForm()

	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	company, err = ParsePostCompany(r.Form)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	log.Println(db)
	err = db.Create(&company).Error
	log.Println("Company created")
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, SuccessfulyCreated())
}

func UpdateCompany(w http.ResponseWriter, r *http.Request) {
	company := r.Context().Value("company").(*Company)
	var data map[string]interface{}

	err := r.ParseForm()
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	log.Println(company.Name)
	data, err = ParsePutData(r.Form)
	err = db.Model(&company).Updates(data).Error
	log.Println(company.Name)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, SuccessfulyUpdated())
}

func DeleteCompany(w http.ResponseWriter, r *http.Request) {
	company := r.Context().Value("company").(*Company)

	err := db.Delete(&company).Error
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, SuccessfulyDeleted())
}