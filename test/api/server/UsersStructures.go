package server

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"net/url"
	"strconv"
)

type UserSerializer struct {
	*User
}

func UserListSerializer(users []*User) []render.Renderer {
	var list []render.Renderer
	for _, user := range users {
		list = append(list, &UserSerializer{User: user})
	}
	return list
}

func (a *UserSerializer) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ParsePostUser(data url.Values) (User, error) {
	var a User
	var company Company
	var compId int
	var err error

	if data["surname"] == nil ||
		data["name"] == nil ||
		data["middlename"] == nil ||
		data["inn"] == nil ||
		data["sex"] == nil ||
		data["login"] == nil ||
		data["password"] == nil ||
		data["email"] == nil ||
		data["company"] == nil {
		return a, errors.New("missing required User fields")
	}

	if len(data["surname"]) != 1 ||
		len(data["name"]) != 1 ||
		len(data["middlename"]) != 1 ||
		len(data["inn"]) != 1 ||
		len(data["sex"]) != 1 ||
		len(data["login"]) != 1 ||
		len(data["password"]) != 1 ||
		len(data["email"]) != 1 ||
		len(data["company"]) != 1 {
		return a, errors.New("multiple forms submitted")
	}

	a.Surname = data["surname"][0]
	a.Name = data["name"][0]
	a.MiddleName = data["middlename"][0]

	a.INN, err = strconv.Atoi(data["inn"][0])
	if err != nil {
		return a, errors.New("invalid values submitted")
	}

	if data["sex"][0] != "M" && data["sex"][0] != "F" {
		return a, errors.New("invalid values submitted")
	}
	a.Sex = data["sex"][0]
	a.Email = data["email"][0]
	a.Login = data["login"][0]

	a.Password = hashAndSalt([]byte(data["password"][0]))

	compId, err = strconv.Atoi(data["company"][0])
	if err != nil {
		return a, errors.New("invalid values submitted")
	}
	err = db.First(&company, compId).Error
	if err != nil {
		return a, errors.New("copmany does not exist")
	}
	a.CompanyRefer = compId
	return a, nil
}
