package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	telebot "gopkg.in/telebot.v3"

	lgg "github.com/ruziba3vich/prodonik_lgger"
	_ "github.com/ruziba3vich/tokenizer/docs"
	handler "github.com/ruziba3vich/tokenizer/internal/http"
	"github.com/ruziba3vich/tokenizer/internal/models"
	"github.com/ruziba3vich/tokenizer/internal/pkg/helper"
	"github.com/ruziba3vich/tokenizer/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

func NewLogger() (*lgg.Logger, error) {
	return lgg.NewLogger("./app.log")
}

func StartServer(h *handler.Handler) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	router := gin.Default()
	router.Use(cors.New(config))
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
		fx.Invoke(
			StartGoBot,
			StartServer,
		),
	)

	app.Run()
}

func StartGoBot() {
	go StartBot()
}

func StartBot() {
	fmt.Println("--------------------- THE BOT IS GOING TO START ------------------")
	pref := telebot.Settings{
		Token:  "7217692907:AAHaBAN4efeuXqbTr54POoHAo8-8LQhxdsc",
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	// Test bot connection
	me := bot.Me
	log.Printf("Bot started successfully: @%s", me.Username)

	bot.Handle("/start", func(c telebot.Context) error {
		log.Printf("request received: /start from %d %s", c.Sender().ID, c.Sender().LastName)

		// Fixed message - removed problematic characters and simplified markdown
		message := "Hello there! 👋" +
			"\nI'm your friendly bot. Here's what I can do:" +
			"\n\nCommands:" +
			"\n/start - Shows this welcome message and available commands" +
			"\n/generate_key - Generates a one-time link for something awesome" +
			"\n\nGot questions or suggestions? Feel free to ask!"

		return c.Send(message)

	})

	// --- Existing: /generate_key handler ---
	bot.Handle("/generate_key", func(c telebot.Context) error {
		log.Printf("request received: /generate_key from %d %s", c.Sender().ID, c.Sender().LastName)

		resp, err := http.Post("http://168.119.255.188:7777/generate-url", "application/json", nil)
		if err != nil {
			log.Printf("Error contacting API: %v", err)
			return c.Send("Failed to contact API: " + err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("API returned non-OK status: %d", resp.StatusCode)
			return c.Send("API returned an unexpected status code.")
		}

		var result models.GenerateResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("Error decoding API response: %v", err)
			return c.Send("Failed to parse API response.")
		}

		if result.Error != "" {
			return c.Send("❌ Error: " + result.Error)
		}

		return c.Send(fmt.Sprintf("✅ Your one-time link: ```%s```\n", result.URL))
	})

	// Add a catch-all handler for debugging
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		log.Printf("Received unhandled message: %s from %d", c.Text(), c.Sender().ID)
		return nil // Don't respond to avoid spam
	})

	log.Println("Bot is running...")
	bot.Start() // This call is blocking and keeps the bot running
}
