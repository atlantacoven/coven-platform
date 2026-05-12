package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"rabidaudio.com/coven-door/server/api"
)

type PutSessionRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PutSessionResponseBody struct {
	Id         int    `json:"id"`
	UserSecret []byte `json:"user_secret,omitempty"`
}

func Router(r chi.Router) {
	r.Post("/session", postSession)
}

func postSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body := PutSessionRequestBody{}
	if ok := api.UnmarshalBody(w, r, &body); !ok {
		return
	}

	u, err := AuthenticatePassword(ctx, body.Email, body.Password)
	if err == ErrInvalidPassword {
		api.RespondBadFormat(w, err)
		return
	} else if err != nil {
		api.RespondError(w, err)
	}

	res := PutSessionResponseBody{
		Id: u.Id,
		// UserSecret: ,
	}
	api.Respond(w, &res, "OK")
}
