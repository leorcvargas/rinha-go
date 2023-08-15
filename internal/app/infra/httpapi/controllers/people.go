package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func (p *PeopleController) Search(c *gin.Context) {
	t := c.Query("t")

	if t == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query param 't'"})
		return
	}

	people, err := p.findPeople.Search(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]PersonResponse, 0, len(people))
	for _, person := range people {
		response = append(response, mapPersonResponse(&person))
	}

	c.JSON(http.StatusOK, response)
}

func (p *PeopleController) Get(c *gin.Context) {
	id := c.Param("id")

	person, err := p.findPeople.ByID(id)
	if err != nil {
		if errors.Is(err, people.ErrPersonNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := mapPersonResponse(person)

	c.JSON(http.StatusOK, response)
}

func (p *PeopleController) Create(c *gin.Context) {
	var dto CreatePersonRequest
	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err := dto.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	person, err := p.createPerson.Execute(
		dto.Nickname,
		dto.Name,
		dto.Birthdate,
		dto.Stack,
	)
	if err != nil {
		if errors.Is(err, people.ErrNicknameTaken) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Location", "/pessoas/"+person.ID)

	response := mapPersonResponse(person)

	c.JSON(http.StatusCreated, response)
	return
}

func (p *PeopleController) CountAll(c *gin.Context) {
	result, err := p.countPeople.CountAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
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
