package controllers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
)

var InvalidDtoErr = errors.New("invalid dto")

type CreatePersonRequest struct {
	Nickname  string   `json:"apelido" validate:"required,max=32"`
	Name      string   `json:"nome" validate:"required,max=100"`
	Birthdate string   `json:"nascimento" validate:"required,datetime=2006-01-02"`
	Stack     []string `json:"stack" validate:"dive,max=32"`
}

func (c *CreatePersonRequest) Validate() error {
	if len(c.Nickname) > 32 {
		return InvalidDtoErr
	}

	if len(c.Name) > 100 {
		return InvalidDtoErr
	}

	dateLayout := "2006-01-02"
	if _, err := time.Parse(dateLayout, c.Birthdate); err != nil {
		return InvalidDtoErr
	}

	for _, tech := range c.Stack {
		if len(tech) > 32 {
			return InvalidDtoErr
		}
	}

	return nil
}

type PersonResponse struct {
	ID        string   `json:"id"`
	Nickname  string   `json:"apelido"`
	Name      string   `json:"nome"`
	Birthdate string   `json:"nascimento"`
	Stack     []string `json:"stack"`
}

type PeopleController struct {
	createPerson *people.CreatePerson
	findPeople   *people.FindPeople
	countPeople  *people.CountPeople
}

func mapPersonResponse(person *people.Person) PersonResponse {
	return PersonResponse{
		ID:        person.ID,
		Nickname:  person.Nickname,
		Name:      person.Name,
		Birthdate: person.Birthdate,
		Stack:     person.Stack,
	}
}

func (p *PeopleController) Search(c *fiber.Ctx) error {
	t := c.Query("t")

	if t == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing query param 't'",
		})
	}

	people, err := p.findPeople.Search(t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	response := make([]PersonResponse, 0, len(people))
	log.Infof("people: %+v", people)

	for _, person := range people {
		response = append(response, mapPersonResponse(&person))
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (p *PeopleController) Get(c *fiber.Ctx) error {
	id := c.Params("id")

	person, err := p.findPeople.ByID(id)
	if err != nil {
		if errors.Is(err, people.ErrPersonNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	response := mapPersonResponse(person)

	return c.Status(fiber.StatusOK).JSON(response)
}

func (p *PeopleController) Create(c *fiber.Ctx) error {
	var dto CreatePersonRequest

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	if err := dto.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	person, err := p.createPerson.Execute(
		dto.Nickname,
		dto.Name,
		dto.Birthdate,
		dto.Stack,
	)
	if err != nil {
		if errors.Is(err, people.ErrNicknameTaken) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		log.Error(err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	c.Set("Location", "/pessoas/"+person.ID)

	response := mapPersonResponse(person)

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (p *PeopleController) CountAll(c *fiber.Ctx) error {
	result, err := p.countPeople.CountAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func NewPeopleController(
	createPerson *people.CreatePerson,
	countPeople *people.CountPeople,
	findPeople *people.FindPeople,
) *PeopleController {
	return &PeopleController{
		createPerson: createPerson,
		findPeople:   findPeople,
		countPeople:  countPeople,
	}
}
