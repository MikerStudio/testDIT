package server

import (
	"errors"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func dbGetUser(id int) (*User, error) {
	user := &User{}
	var err error
	err = db.First(user, id).Error
	if err != nil {
		return user, errors.New("user matching query does not exist")
	}
	return user, err
}

func GetUserData(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	err := render.Render(w, r, &UserSerializer{User: user})
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func GetUsersList(w http.ResponseWriter, r *http.Request) {
	var users []*User
	db.Find(&users)
	if err := render.RenderList(w, r, UserListSerializer(users)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := r.ParseForm()
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user, err = ParsePostUser(r.Form)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = db.Create(&user).Error
	log.Println("User created")
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, SuccessfulyCreated())
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)
	var data map[string]interface{}

	err := r.ParseForm()
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	data, err = ParsePutData(r.Form)
	err = db.Model(&user).Updates(data).Error
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, SuccessfulyUpdated())
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)

	err := db.Delete(&user).Error
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, SuccessfulyDeleted())
}