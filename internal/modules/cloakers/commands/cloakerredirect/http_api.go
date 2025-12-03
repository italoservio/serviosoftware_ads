package cloakerredirect

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
)

type RedirectCloakerHttpAPI struct {
	validate           validator.Validate
	RedirectCloakerCmd *RedirectCloakerCmd
}

func NewRedirectCloakerHttpAPI(
	validate *validator.Validate,
	redirectCloakerCmd *RedirectCloakerCmd,
) *RedirectCloakerHttpAPI {
	return &RedirectCloakerHttpAPI{
		validate:           *validate,
		RedirectCloakerCmd: redirectCloakerCmd,
	}
}

func (c *RedirectCloakerHttpAPI) RedirectCloaker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encodedID := vars["encodedId"]

	userIP := getClientIP(r)
	userAgent := r.UserAgent()

	input := RedirectCloakerInput{
		EncodedID: encodedID,
		UserIP:    userIP,
		UserAgent: userAgent,
	}

	if err := c.validate.Struct(input); err != nil {
		exception.NewValidatorException(err).WriteJSON(w)
		return
	}

	output, err := c.RedirectCloakerCmd.Exec(&input)
	if err != nil {
		exception.ToAppException(err).WriteJSON(w)
		return
	}

	http.Redirect(w, r, output.RedirectURL, http.StatusFound)
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	return r.RemoteAddr
}
