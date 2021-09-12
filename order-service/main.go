package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	HOST := os.Getenv("POSTGRES_HOST")
	// HOST := "localhost"
	PASS := os.Getenv("POSTGRES_PASSWORD")
	// PASS := "docker"
	USER := os.Getenv("POSTGRES_USER")
	// USER := "docker"
	DBNAME := os.Getenv("POSTGRES_DB")
	// DBNAME := "orders"
	PORT := "5432"

	dsn := fmt.Sprintf("host=%s port=%s user=%s database=%s password=%s sslmode=disable", HOST, PORT, USER, DBNAME, PASS)
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		fmt.Println(os.Getenv("POSTGRES_HOST"))
		fmt.Println(os.Getenv("POSTGRES_PASSWORD"))
		fmt.Println(os.Getenv("POSTGRES_USER"))
		fmt.Println(os.Getenv("POSTGRES_DB"))
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&Order{})

	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)
	router := gin.Default()
	router.POST("/orders", handler.CreateOrder)
	router.Run(":8081")
}

// MODEL
type Order struct {
	ID           int
	BookID       int
	BookName     string
	CustomerName string
	Status       string
	Price        int
	Quantity     int
	TotalPrice   int
	CreatedAt    time.Time
}

// INPUT
type OrderInput struct {
	BookID       int    `json:"book_id" binding:"required"`
	CustomerName string `json:"customer_name" binding:"required"`
	Quantity     int    `json:"quantity" binding:"required"`
	Price        int    `json:"price"`
	BookName     string `json:"book_name"`
}

/*
{
	"book_id": 1,
	"customer_name": "Arief",
	"quantity": 2
}
*/

// REPOSITORY
type Repository interface {
	CreateOrder(order Order) (Order, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateOrder(order Order) (Order, error) {
	err := r.db.Create(&order).Error
	if err != nil {
		return order, err
	}
	return order, nil
}

// SERVICE
type Service interface {
	CreateOrder(input OrderInput) (Order, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateOrder(input OrderInput) (Order, error) {
	var order Order
	order.BookID = input.BookID
	order.CustomerName = input.CustomerName
	order.Price = input.Price
	order.Status = "Pending"
	order.BookName = input.BookName
	order.Quantity = input.Quantity
	order.TotalPrice = input.Quantity * input.Price

	newOrder, err := s.repository.CreateOrder(order)
	if err != nil {
		return newOrder, err
	}
	return newOrder, nil
}

// FORMATTER
type OrderFormatter struct {
	ID           int       `json:"id"`
	BookID       int       `json:"book_id"`
	CustomerName string    `json:"customer_name"`
	Price        int       `json:"price"`
	Quantity     int       `json:"quantity"`
	TotalPrice   int       `json:"total_price"`
	Status       string    `json:"status"`
	BookName     string    `json:"book_name"`
	CreatedAt    time.Time `json:"created_at"`
}

func FormatOrder(order Order) OrderFormatter {
	var orderFormatter OrderFormatter
	orderFormatter.ID = order.ID
	orderFormatter.BookID = order.BookID
	orderFormatter.CustomerName = order.CustomerName
	orderFormatter.Price = order.Price
	orderFormatter.Status = order.Status
	orderFormatter.CreatedAt = order.CreatedAt
	orderFormatter.BookName = order.BookName
	orderFormatter.Quantity = order.Quantity
	orderFormatter.TotalPrice = order.TotalPrice

	return orderFormatter
}

// HANDLER
type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service}
}

func (h *handler) CreateOrder(c *gin.Context) {
	var input OrderInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	client := &http.Client{}
	requestUrl := fmt.Sprintf("http://%s/books/%d", os.Getenv("BOOK_SERVICE_HOST"), input.BookID)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, _ := client.Do(req)
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var responseObject ResponseBook
	json.Unmarshal(bodyBytes, &responseObject)

	if responseObject.Data.ID == 0 {
		response := gin.H{"status": "error", "message": "Book Not Found"}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	input.BookName = responseObject.Data.Name
	input.Price = responseObject.Data.Price

	order, err := h.service.CreateOrder(input)
	if err != nil {
		response := gin.H{"status": "error", "message": err.Error()}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := gin.H{"status": "success", "data": FormatOrder(order)}
	c.JSON(http.StatusOK, response)
}

type ResponseBook struct {
	Status string           `json:"status"`
	Data   ResponseBookData `json:"data"`
}

type ResponseBookData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}
