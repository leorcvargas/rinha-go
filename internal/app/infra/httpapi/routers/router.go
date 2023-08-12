package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Router interface {
	Load()
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func MakeRouter(
	peopleRouter *PeopleRouter,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	peopleRouter.Load(r)

	return r
}
