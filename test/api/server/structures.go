package server

import (
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func StringArrToInt(t []string) ([]int, error) {
	var t2 = []int{}
	var err error
	var j int
	for _, i := range t {
		j, err = strconv.Atoi(i)
		if err != nil {
			return t2, err
		}
		t2 = append(t2, j)
	}
	return t2, nil
}

func hashAndSalt(pwd []byte) string {

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ParsePutData(data url.Values) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	for k, v := range data {
		res[k] = v[0]
	}

	return res, nil
}

type DefaultResponse struct {
	Result         string `json:"-"` // low-level runtime error
	HTTPStatusCode int    `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *DefaultResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrRender(err error) render.Renderer {
	return &DefaultResponse{
		Result:         err.Error(),
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response",
		ErrorText:      err.Error(),
	}
}

func ErrInvalidRequest(err error) render.Renderer {
	return &DefaultResponse{
		Result:         err.Error(),
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}

func SuccessfulyCreated() render.Renderer {
	return &DefaultResponse{
		Result:         "Instance created successfully",
		HTTPStatusCode: 201,
		StatusText:     "Created",
		ErrorText:      "",
	}
}

func SuccessfulyUpdated() render.Renderer {
	return &DefaultResponse{
		Result:         "Instance updated successfully",
		HTTPStatusCode: 200,
		StatusText:     "Updated",
		ErrorText:      "",
	}
}

func SuccessfulyDeleted() render.Renderer {
	return &DefaultResponse{
		Result:         "Instance deleted successfully",
		HTTPStatusCode: 200,
		StatusText:     "Deleted",
		ErrorText:      "",
	}
}
