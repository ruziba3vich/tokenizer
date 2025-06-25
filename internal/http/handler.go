package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	lgg "github.com/ruziba3vich/prodonik_lgger"
	"github.com/ruziba3vich/tokenizer/internal/models"
	"github.com/ruziba3vich/tokenizer/internal/service"
)

type Handler struct {
	service *service.Service
	logger  *lgg.Logger
}

func NewHandler(service *service.Service, logger *lgg.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) GenerateOneTimeLink(c *gin.Context) {
	link, err := h.service.GenerateOneTimeLink()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to generate url: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": link})
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Registers a user using a one-time key. Key must not have been used before.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param payload body models.RegisterPayload true "User registration data"
// @Success 200 {object} map[string]string "message: user registered successfully"
// @Failure 400 {object} map[string]string "error: invalid request or key already used"
// @Router /register [post]
func (h *Handler) RegisterUser(c *gin.Context) {
	var payload models.RegisterPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Println("Invalid payload:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.CreateUser(&payload); err != nil {
		h.logger.Println("CreateUser error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

// Login godoc
// @Summary Login user
// @Description Authenticates a user using email and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param payload body models.LoginPayload true "Login credentials"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "error: invalid request"
// @Failure 401 {object} map[string]string "error: invalid email or password"
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var payload models.LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Println("Login payload bind failed:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.service.GetUserByEmailAndPassword(payload.Email, payload.Password)
	if err != nil {
		h.logger.Println("Login failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
