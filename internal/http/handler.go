package handler

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	lgg "github.com/ruziba3vich/prodonik_lgger"
	"github.com/ruziba3vich/tokenizer/internal/models"
	"github.com/ruziba3vich/tokenizer/internal/service"
)

type Handler struct {
	service    *service.Service
	logger     *lgg.Logger
	privateKey *rsa.PrivateKey
}

func NewHandler(service *service.Service, privateKey *rsa.PrivateKey, logger *lgg.Logger) *Handler {
	return &Handler{
		service:    service,
		logger:     logger,
		privateKey: privateKey,
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
// @Success 200 {object} SignedResponse
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

	_, err := h.service.GetUserByEmailAndPassword(payload.Email, payload.Password)
	if err != nil {
		h.logger.Println("Login failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// 2. Create signed payload
	responsePayload := ResponsePayload{
		Status: "APPROVED",
	}
	toSign := struct {
		Status string `json:"status"`
		Vhid   string `json:"vhid"`
	}{
		Status: responsePayload.Status,
		Vhid:   payload.Vhid,
	}

	dataToSign, err := json.Marshal(toSign)
	if err != nil {
		h.logger.Println("Failed to marshal payload:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal payload"})
		return
	}

	hash := sha256.Sum256(dataToSign)
	signature, err := rsa.SignPKCS1v15(nil, h.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		h.logger.Println("Failed to sign payload:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create signature"})
		return
	}

	sigBase64 := base64.StdEncoding.EncodeToString(signature)

	// 3. Return response
	response := SignedResponse{
		Payload:   responsePayload,
		Signature: sigBase64,
	}

	c.JSON(http.StatusOK, response)
}

type ResponsePayload struct {
	Status string `json:"status"`
}

type Request struct {
	Vhid string `json:"vhid"`
}

type SignedResponse struct {
	Payload   ResponsePayload `json:"payload"`
	Signature string          `json:"signature"`
}
