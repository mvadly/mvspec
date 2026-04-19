package main

import (
	"fmt"
	"net/http"
	"strconv"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	mvdocs "github.com/mvadly/mvspec/mv-docs"
)

// Dummy data
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}

// User represents a user model
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// UserRequest represents create/update request
type UserRequest struct {
	Name  string `json:"name" form:"name" binding:"required"`
	Email string `json:"email" form:"email" binding:"required,email"`
	Age   int    `json:"age" form:"age" binding:"required"`
}

// TextRequest represents plain text request
type TextRequest struct {
	Content string `json:"content"`
}

// UploadResponse represents file upload response
type UploadResponse struct {
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
	Message  string `json:"message"`
}

// FormResponse represents form data response
type FormResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Category string `json:"category"`
}

// MultipartForm represents multipart form without file upload
type MultipartForm struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Category string `form:"category"`
}

// UploadWithForm represents file upload with additional form fields
type UploadWithForm struct {
	File        *multipart.FileHeader `form:"file" binding:"required"`
	Description string                `form:"description"`
	Title       string                `form:"title"`
}

// UrlencodedForm represents urlencoded form data
type UrlencodedForm struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Age      int    `form:"age" binding:"required"`
	Category string `form:"category"`
}

// APIResponse represents generic API response
type APIResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	r := gin.Default()
	r.GET("/mvdocs/*path", gin.WrapH(mvdocs.MvHandler()))

	// Health check
	// @Summary Health check
	// @Description Returns API health status
	// @Tags health
	// @Produce json
	// @Success 200 {object} APIResponse
	// @Router /api/health [get]
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, APIResponse{
			Code:    "00",
			Message: "API is running",
		})
	})

	// User endpoints group
	users := r.Group("/api/features/users")
	{
		// @Summary Get all users
		// @Description Returns list of all users
		// @Tags users
		// @Accept json
		// @Produce json
		// @Success 200 {array} User
		// @Router /api/features/users [get]
		users.GET("", getUsers)

		// @Summary Get user by ID
		// @Description Returns a single user by ID
		// @Tags users
		// @Accept json
		// @Produce json
		// @Param id path int true "User ID"
		// @Success 200 {object} User
		// @Failure 404 {object} APIResponse
		// @Router /api/features/users/{id} [get]
		users.GET("/:id", getUserByID)

		// @Summary Create new user
		// @Description Creates a new user
		// @Tags users
		// @Accept application/json
		// @Produce json
		// @Param body body UserRequest true "User data"
		// @Success 201 {object} User
		// @Failure 400 {object} APIResponse
		// @Router /api/features/users [post]
		users.POST("", createUser)

		// @Summary Update user
		// @Description Updates an existing user
		// @Tags users
		// @Accept application/json
		// @Produce json
		// @Param id path int true "User ID"
		// @Param body body UserRequest true "User data"
		// @Success 200 {object} User
		// @Failure 400 {object} APIResponse
		// @Failure 404 {object} APIResponse
		// @Router /api/features/users/{id} [put]
		users.PUT("/:id", updateUser)

		// @Summary Patch user
		// @Description Partially updates a user
		// @Tags users
		// @Accept application/json
		// @Produce json
		// @Param id path int true "User ID"
		// @Param body body UserRequest true "User data"
		// @Success 200 {object} User
		// @Failure 400 {object} APIResponse
		// @Failure 404 {object} APIResponse
		// @Router /api/features/users/{id} [patch]
		users.PATCH("/:id", patchUser)

		// @Summary Delete user
		// @Description Deletes a user
		// @Tags users
		// @Param id path int true "User ID"
		// @Success 204
		// @Failure 404 {object} APIResponse
		// @Router /api/features/users/{id} [delete]
		users.DELETE("/:id", deleteUser)
	}

	// File upload endpoint (individual param)
	// @Summary Upload file
	// @Description Handles file multipart upload
	// @Tags upload
	// @Accept multipart/form-data
	// @Produce json
	// @Param file formData file true "File to upload"
	// @Success 200 {object} UploadResponse
	// @Failure 400 {object} APIResponse
	// @Router /api/features/upload [post]
	r.POST("/api/features/upload", uploadFile)

	// Multipart form with file upload (struct body)
	// @Summary Upload file with form data
	// @Description Handles file multipart upload with additional form fields
	// @Tags upload
	// @Accept multipart/form-data
	// @Produce json
	// @Param body formData UploadWithForm true "Upload with form data"
	// @Success 200 {object} UploadResponse
	// @Failure 400 {object} APIResponse
	// @Router /api/features/upload-struct [post]
	r.POST("/api/features/upload-struct", uploadWithForm)

	// Multipart form without file (struct body)
	// @Summary Submit multipart form
	// @Description Handles multipart form data without file upload
	// @Tags form
	// @Accept multipart/form-data
	// @Produce json
	// @Param body formData MultipartForm true "Multipart form data"
	// @Success 200 {object} FormResponse
	// @Failure 400 {object} APIResponse
	// @Router /api/features/form-multipart [post]
	r.POST("/api/features/form-multipart", handleMultipartForm)

	// URLEncoded form (struct body)
	// @Summary Submit urlencoded form
	// @Description Handles form-urlencoded data
	// @Tags form
	// @Accept application/x-www-form-urlencoded
	// @Produce json
	// @Param body formData UrlencodedForm true "URLEncoded form data"
	// @Success 200 {object} FormResponse
	// @Failure 400 {object} APIResponse
	// @Router /api/features/form-urlencoded [post]
	r.POST("/api/features/form-urlencoded", handleUrlencodedForm)

	// Form endpoint (individual params - urlencoded)
	// @Summary Submit form data
	// @Description Handles form-urlencoded data
	// @Tags form
	// @Accept application/x-www-form-urlencoded
	// @Produce json
	// @Param name formData string true "Name"
	// @Param email formData string true "Email"
	// @Param age formData int true "Age"
	// @Param category formData string false "Category"
	// @Success 200 {object} FormResponse
	// @Failure 400 {object} APIResponse
	// @Router /api/features/form [post]
	r.POST("/api/features/form", handleForm)

	// Plain text endpoint
	// @Summary Handle plain text
	// @Description Handles plain text request
	// @Tags text
	// @Accept text/plain
	// @Produce json
	// @Param content body string true "Text content"
	// @Success 200 {object} APIResponse
	// @Failure 400 {object} APIResponse
	// @Router /api/features/text [post]
	r.POST("/api/features/text", handleText)

	r.Run(":8080")
}

// @Summary Get all users
// @Description Returns list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /api/features/users [get]
func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

// @Summary Get user by ID
// @Description Returns a single user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} APIResponse
// @Router /api/features/users/{id} [get]
func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "Invalid ID format",
		})
		return
	}

	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}

	c.JSON(http.StatusNotFound, APIResponse{
		Code:    "44",
		Message: "User not found",
	})
}

// @Summary Create new user
// @Description Creates a new user
// @Tags users
// @Accept application/json
// @Produce json
// @Param body body UserRequest true "User data"
// @Success 201 {object} User
// @Failure 400 {object} APIResponse
// @Router /api/features/users [post]
func createUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	newID := len(users) + 1
	newUser := User{
		ID:    newID,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	users = append(users, newUser)

	c.JSON(http.StatusCreated, newUser)
}

// @Summary Update user
// @Description Updates an existing user
// @Tags users
// @Accept application/json
// @Produce json
// @Param id path int true "User ID"
// @Param body body UserRequest true "User data"
// @Success 200 {object} User
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /api/features/users/{id} [put]
func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "Invalid ID format",
		})
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	for i, user := range users {
		if user.ID == id {
			users[i].Name = req.Name
			users[i].Email = req.Email
			users[i].Age = req.Age
			c.JSON(http.StatusOK, users[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, APIResponse{
		Code:    "44",
		Message: "User not found",
	})
}

// @Summary Patch user
// @Description Partially updates a user
// @Tags users
// @Accept application/json
// @Produce json
// @Param id path int true "User ID"
// @Param body body UserRequest true "User data"
// @Success 200 {object} User
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /api/features/users/{id} [patch]
func patchUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "Invalid ID format",
		})
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	for i, user := range users {
		if user.ID == id {
			if req.Name != "" {
				users[i].Name = req.Name
			}
			if req.Email != "" {
				users[i].Email = req.Email
			}
			if req.Age > 0 {
				users[i].Age = req.Age
			}
			c.JSON(http.StatusOK, users[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, APIResponse{
		Code:    "44",
		Message: "User not found",
	})
}

// @Summary Delete user
// @Description Deletes a user
// @Tags users
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {object} APIResponse
// @Router /api/features/users/{id} [delete]
func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "Invalid ID format",
		})
		return
	}

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, APIResponse{
		Code:    "44",
		Message: "User not found",
	})
}

// @Summary Upload file
// @Description Handles file multipart upload
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} UploadResponse
// @Failure 400 {object} APIResponse
// @Router /api/features/upload [post]
func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "File is required",
		})
		return
	}

	filename := file.Filename
	size := file.Size

	c.JSON(http.StatusOK, UploadResponse{
		FileName: filename,
		Size:     size,
		Message:  "File uploaded successfully",
	})
}

// @Summary Submit form data
// @Description Handles form-urlencoded data
// @Tags form
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param name formData string true "Name"
// @Param email formData string true "Email"
// @Param age formData int true "Age"
// @Param category formData string false "Category"
// @Success 200 {object} FormResponse
// @Failure 400 {object} APIResponse
// @Router /api/features/form [post]
func handleForm(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	ageStr := c.PostForm("age")
	category := c.DefaultPostForm("category", "general")

	if name == "" || email == "" || ageStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "name, email, and age are required",
		})
		return
	}

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "Invalid age format",
		})
		return
	}

	c.JSON(http.StatusOK, FormResponse{
		Name:     name,
		Email:    email,
		Age:      age,
		Category: category,
	})
}

// @Summary Handle plain text
// @Description Handles plain text request
// @Tags text
// @Accept text/plain
// @Produce json
// @Param content body string true "Text content"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Router /api/features/text [post]
func handleText(c *gin.Context) {
	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "content is required",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    "00",
		Message: "Text received",
		Data: map[string]string{
			"content": content,
		},
	})
}

// @Summary Upload file with form data
// @Description Handles file multipart upload with additional form fields
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param body formData UploadWithForm true "Upload with form data"
// @Success 200 {object} UploadResponse
// @Failure 400 {object} APIResponse
// @Router /api/features/upload-struct [post]
func uploadWithForm(c *gin.Context) {
	var form UploadWithForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: fmt.Sprintf("Binding error: %s", err.Error()),
		})
		return
	}

	file, err := form.File.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: "Failed to open file",
		})
		return
	}
	defer file.Close()

	c.JSON(http.StatusOK, UploadResponse{
		FileName: form.File.Filename,
		Size:     form.File.Size,
		Message:  "File uploaded with: " + form.Description,
	})
}

// @Summary Submit multipart form
// @Description Handles multipart form data without file upload
// @Tags form
// @Accept multipart/form-data
// @Produce json
// @Param body formData MultipartForm true "Multipart form data"
// @Success 200 {object} FormResponse
// @Failure 400 {object} APIResponse
// @Router /api/features/form-multipart [post]
func handleMultipartForm(c *gin.Context) {
	var form MultipartForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, FormResponse{
		Name:     form.Name,
		Email:    form.Email,
		Category: form.Category,
	})
}

// @Summary Submit urlencoded form
// @Description Handles form-urlencoded data
// @Tags form
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param body formData UrlencodedForm true "URLEncoded form data"
// @Success 200 {object} FormResponse
// @Failure 400 {object} APIResponse
// @Router /api/features/form-urlencoded [post]
func handleUrlencodedForm(c *gin.Context) {
	var form UrlencodedForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    "99",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, FormResponse{
		Name:     form.Name,
		Email:    form.Email,
		Age:      form.Age,
		Category: form.Category,
	})
}
