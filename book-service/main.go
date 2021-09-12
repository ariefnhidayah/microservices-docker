package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	// "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	HOST := os.Getenv("POSTGRES_HOST")
	PASS := os.Getenv("POSTGRES_PASSWORD")
	USER := os.Getenv("POSTGRES_USER")
	DBNAME := os.Getenv("POSTGRES_DB")
	PORT := "5432"
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	dsn := fmt.Sprintf("host=%s port=%s user=%s database=%s password=%s sslmode=disable", HOST, PORT, USER, DBNAME, PASS)
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		fmt.Println(os.Getenv("POSTGRES_HOST"))
		fmt.Println(os.Getenv("POSTGRES_PASSWORD"))
		fmt.Println(os.Getenv("POSTGRES_USER"))
		fmt.Println(os.Getenv("POSTGRES_DB"))
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&Book{})

	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	router := gin.Default()
	router.GET("/books", handler.GetBooks)
	router.GET("/books/:id", handler.GetBook)
	router.POST("/books", handler.Create)
	router.PUT("/books/:id", handler.Update)
	router.DELETE("/books/:id", handler.Delete)

	router.Run(":8080")
}

// ENTITY / MODEL
type Book struct {
	ID          int
	Name        string
	Description string
	Price       int
}

// INPUT
type CreateBookInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int    `json:"price" binding:"required"`
}

type GetBookInput struct {
	ID int `uri:"id" binding:"required"`
}

// REPOSITORY
type Repository interface {
	FindAll() ([]Book, error)
	FindById(ID int) (Book, error)
	Save(book Book) (Book, error)
	Update(book Book) (Book, error)
	Delete(book Book) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Book, error) {
	var books []Book
	err := r.db.Find(&books).Error
	if err != nil {
		return books, err
	}
	return books, nil
}

func (r *repository) FindById(ID int) (Book, error) {
	var book Book
	err := r.db.Where("id = ?", ID).Find(&book).Error
	if err != nil {
		return book, err
	}
	return book, nil
}

func (r *repository) Save(book Book) (Book, error) {
	err := r.db.Create(&book).Error
	if err != nil {
		return book, err
	}
	return book, nil
}

func (r *repository) Update(book Book) (Book, error) {
	err := r.db.Save(&book).Error
	if err != nil {
		return book, err
	}

	return book, nil
}

func (r *repository) Delete(book Book) error {
	err := r.db.Where("id = ?", book.ID).Delete(&Book{}).Error
	return err
}

// SERVICE
type Service interface {
	GetBooks() ([]Book, error)
	GetBook(inputID GetBookInput) (Book, error)
	CreateBook(input CreateBookInput) (Book, error)
	UpdateBook(inputID GetBookInput, input CreateBookInput) (Book, error)
	DeleteBook(inputID GetBookInput) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetBooks() ([]Book, error) {
	books, err := s.repository.FindAll()
	if err != nil {
		return books, err
	}

	return books, nil
}

func (s *service) GetBook(inputID GetBookInput) (Book, error) {
	book, err := s.repository.FindById(inputID.ID)
	if err != nil {
		return book, err
	}
	return book, nil
}

func (s *service) CreateBook(input CreateBookInput) (Book, error) {
	book := Book{}
	book.Name = input.Name
	book.Description = input.Description
	book.Price = input.Price

	newBook, err := s.repository.Save(book)
	if err != nil {
		return newBook, err
	}

	return newBook, nil
}

func (s *service) UpdateBook(inputID GetBookInput, input CreateBookInput) (Book, error) {
	book, err := s.repository.FindById(inputID.ID)
	if err != nil {
		return book, err
	}

	book.Name = input.Name
	book.Description = input.Description
	book.Price = input.Price

	updatedBook, err := s.repository.Update(book)
	if err != nil {
		return updatedBook, err
	}

	return updatedBook, nil
}

func (s *service) DeleteBook(inputID GetBookInput) error {
	book, err := s.repository.FindById(inputID.ID)
	if err != nil {
		return err
	}

	err = s.repository.Delete(book)
	return err
}

// FORMATTER
type BookFormatter struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

func FormatBook(book Book) BookFormatter {
	bookFormatter := BookFormatter{
		ID:          book.ID,
		Name:        book.Name,
		Price:       book.Price,
		Description: book.Description,
	}
	return bookFormatter
}

func FormatBooks(books []Book) []BookFormatter {
	booksFormatter := []BookFormatter{}

	for _, book := range books {
		bookFormatter := FormatBook(book)
		booksFormatter = append(booksFormatter, bookFormatter)
	}

	return booksFormatter
}

// HANDLER
type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service}
}

func (h *handler) GetBooks(c *gin.Context) {
	books, err := h.service.GetBooks()
	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := gin.H{"status": "success", "data": FormatBooks(books)}
	c.JSON(http.StatusOK, response)
}

func (h *handler) GetBook(c *gin.Context) {
	var input GetBookInput
	err := c.ShouldBindUri(&input)

	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	book, err := h.service.GetBook(input)
	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := gin.H{"status": "success", "data": FormatBook(book)}
	c.JSON(http.StatusOK, response)
}

func (h *handler) Create(c *gin.Context) {
	var input CreateBookInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newBook, err := h.service.CreateBook(input)
	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := gin.H{"status": "success", "data": FormatBook(newBook)}
	c.JSON(http.StatusOK, response)
}

func (h *handler) Update(c *gin.Context) {
	var input GetBookInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData CreateBookInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	updatedBook, err := h.service.UpdateBook(input, inputData)
	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := gin.H{"status": "success", "data": FormatBook(updatedBook)}
	c.JSON(http.StatusOK, response)
}

func (h *handler) Delete(c *gin.Context) {
	var input GetBookInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.service.DeleteBook(input)
	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := gin.H{"status": "success"}
	c.JSON(http.StatusOK, response)
}
