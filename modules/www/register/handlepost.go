package register

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

	// Check if the user already exists
	existingUser, err := models.UserByEmail(r.Context(), c.svcCtx.SqlxDB, c.form.Email)
	if err == nil && existingUser != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			"User already exists",
		))
		return
	}

	// Create the new user
	user := &models.User{
		Email:    c.form.Email,
		Password: c.form.Password,
	}
	err = user.HashPassword()
	if err != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			"Failed to create user",
		))
		return
	}
	err = models.CreateUser(r.Context(), c.svcCtx.SqlxDB, user)
	if err != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			"Failed to create user",
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

	// Redirect to the home page
	htmx.Redirect(w, r, "/app")
}
