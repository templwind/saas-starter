package login

import (
	"net/http"

	"github.com/templwind/sass-starter/internal/middleware"
	"github.com/templwind/sass-starter/internal/models"
	"github.com/templwind/sass-starter/internal/ui/components/apperror"
	"github.com/templwind/sass-starter/internal/utils"

	"github.com/templwind/templwind/forms"
	"github.com/templwind/templwind/htmx"
)

func (c *Controller) HandlePost(w http.ResponseWriter, r *http.Request) {
	// Bind the form using the default binder
	err := forms.Bind(r, c.form)
	if err != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			"Invalid form data",
		))
		return
	}

	// Validate the form
	err = c.form.Validate()
	if err != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			err.Error(),
		))
		return
	}

	// Find the user by email
	user, err := models.UserByEmail(r.Context(), c.svcCtx.SqlxDB, c.form.Email)
	if err != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			"User not found",
		))
		return
	}

	// Authenticate the user
	if !user.ValidatePassword(c.form.Password) {
		utils.Render(w, r, http.StatusOK, apperror.New(
			"Invalid password",
		))
		return
	}

	// Set the authentication token
	err = middleware.SetAuthToken(w, r, c.svcCtx, user)
	if err != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			err.Error(),
		))
		return
	}

	// How many accounts are associated with this?
	accounts, err := models.UserAccountsByUserID(r.Context(), c.svcCtx.SqlxDB, user.ID)
	if err != nil && err != models.ErrDoesNotExist {
		utils.Render(w, r, http.StatusOK, apperror.New(
			err.Error(),
		))
		return
	}

	if len(accounts) == 1 {
		err = middleware.SetAccountToken(w, r, c.svcCtx, accounts[0])
		if err != nil {
			utils.Render(w, r, http.StatusOK, apperror.New(
				err.Error(),
			))
			return
		}
	}

	// Add a header for a HX-Redirect
	htmx.Redirect(w, r, "/app")
}
