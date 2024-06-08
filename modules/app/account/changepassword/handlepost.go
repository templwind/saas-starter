package changepassword

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
	ctx := r.Context()
	account := middleware.AccountFromContext(r)
	user, _ := models.UserByID(ctx, c.svcCtx.SqlxDB, account.PrimaryUserID)

	changePasswordForm := ChangePasswordForm{}
	forms.Bind(r, &changePasswordForm)

	// authenticate the user
	if !user.ValidatePassword(changePasswordForm.CurrentPassword) {
		utils.Render(w, r, http.StatusOK, apperror.New(
			"Your current password is incorrect",
		))
		return
	}

	// validate the account
	err := changePasswordForm.Validate()
	if err != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			err.Error(),
		))
		return
	}

	if err := user.UpdateWithPassword(ctx, c.svcCtx.SqlxDB, changePasswordForm.NewPassword); err != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			err.Error(),
		))
		return
	}

	// fire the event trigger
	if err := htmx.Trigger(w, r, "on-change-password-success"); err != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			err.Error(),
		))
		return
	}
}
