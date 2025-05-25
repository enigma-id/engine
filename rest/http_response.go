package rest

import (
	"net/http"

	"github.com/enigma-id/engine/validate"
	"github.com/labstack/echo/v4"
)

type (
	ResponseBody struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
		Total   int64       `json:"total,omitempty"`
	}
)

func (r *ResponseBody) Body(data interface{}, total int64) {
	r.Data = data
	r.Total = total
}

func NewActionError(msg string, e error) *echo.HTTPError {
	return &echo.HTTPError{
		Code:     http.StatusBadRequest,
		Message:  msg,
		Internal: e,
	}
}

func NewResponseBody() *ResponseBody {
	return &ResponseBody{
		Message: "success",
	}
}

func Response(c echo.Context, rb interface{}, e error) error {
	if e != nil {
		return e
	}

	r, ok := rb.(*ResponseBody)
	if !ok {
		r = &ResponseBody{
			Data:    rb,
			Message: "success",
		}
	}

	return c.JSON(200, r)
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	he, ok := err.(*echo.HTTPError)
	if !ok {
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	var resp interface{}

	if RestConfig.IsDev && he.Internal != nil {
		resp = echo.Map{"message": he.Message, "error": he.Internal.Error()}
	} else {
		resp = echo.Map{"message": he.Message}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(he.Code)
		} else {
			// check if error validations
			if ev, ok := err.(*validate.Response); ok {
				resp = echo.Map{"message": "Your input has a problem", "error": ev.GetMessages()}
				err = c.JSON(422, resp)
			} else {
				err = c.JSON(he.Code, resp)
			}
		}

		if err != nil {
			Server.Logger.Error(err)
		}
	}
}
