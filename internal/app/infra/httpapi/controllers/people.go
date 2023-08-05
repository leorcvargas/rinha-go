package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
)

type PeopleController struct {
	createPeople *people.CreatePeople
}

func (p *PeopleController) GetAll(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (p *PeopleController) Create(c echo.Context) error {
	person := new(people.Person)
	if err := c.Bind(person); err != nil {
		return err
	}

	result, err := p.createPeople.Execute(person)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Location", "/pessoas/"+result.ID)

	return c.JSON(http.StatusCreated, result)
}

func NewPeopleController(
	createPeople *people.CreatePeople,
) *PeopleController {
	return &PeopleController{
		createPeople: createPeople,
	}
}
