package beego

import (
	beegoWeb "github.com/beego/beego/v2/server/web"
	"std-library/app/web/errors"
)

type errorController struct {
	beegoWeb.Controller
}

func (c *errorController) Error404() {
	errors.NotFoundError(404)
}

func (c *errorController) Error501() {
	errors.InternalError(501)
}

func (c *errorController) ErrorDb() {
	errors.InternalError(500)
}
