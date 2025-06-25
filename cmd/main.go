package main

import (
	"log"

	"github.com/gin-gonic/gin"

	lgg "github.com/ruziba3vich/prodonik_lgger"
	_ "github.com/ruziba3vich/tokenizer/docs"
	handler "github.com/ruziba3vich/tokenizer/internal/http"
	"github.com/ruziba3vich/tokenizer/internal/pkg/helper"
	"github.com/ruziba3vich/tokenizer/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

func NewLogger() (*lgg.Logger, error) {
	return lgg.NewLogger("/etc/app.log")
}

func StartServer(h *handler.Handler) {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/generate-url", h.GenerateOneTimeLink)
	router.POST("/register", h.RegisterUser)
	router.POST("/login", h.Login)

	if err := router.Run(":7777"); err != nil {
		log.Fatal("failed to run server:", err)
	}
}

func main() {
	app := fx.New(
		fx.Provide(
			NewLogger,
			helper.NewDB,
			service.NewService,
			handler.NewHandler,
		),
		fx.Invoke(StartServer),
	)

	app.Run()
}
