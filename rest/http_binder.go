package rest

import (
	"net/http"
	"sync"

	"github.com/enigma-id/engine/validate"
	"github.com/labstack/echo/v4"
)

type (
	CustomBinder struct {
		once      sync.Once
		validator *validate.Validator
	}
)

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
	db := new(echo.DefaultBinder)

	if err = db.Bind(i, c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if errv := cb.Validate(i); !errv.Valid {
		return errv
	}

	return
}

// Validate the request when binding
func (cb *CustomBinder) Validate(obj interface{}) (resp *validate.Response) {
	cb.lazyinit()

	if vr, ok := obj.(validate.Request); ok {
		resp = cb.validator.Request(vr)
	} else {
		resp = cb.validator.Struct(obj)
	}

	return
}

// lazyinit initialing validator instances for one of time only.
func (cb *CustomBinder) lazyinit() {
	cb.once.Do(func() {
		cb.validator = validate.New()
	})
}
