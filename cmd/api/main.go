package main

import (
	"fmt"
	"github.com/MetaLoan/pp/internal/config"
	"github.com/MetaLoan/pp/internal/core"
	"github.com/MetaLoan/pp/internal/handler"
	"github.com/MetaLoan/pp/internal/service"
	"github.com/MetaLoan/pp/internal/worker"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 1. Load Config
	config.LoadConfig()

	// 2. Init Database
	core.InitDB()

	// 3. Init Market Engine
	core.InitMarket()

	// 4. Init Gin
	if config.GlobalConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Add CORS Middleware
	r.Use(handler.CORSMiddleware())

	// 3. Setup Routes
	authHandler := handler.NewAuthHandler()
	tradeService := service.NewTradeService()
	tradeHandler := handler.NewTradeHandler(tradeService)
	fundsHandler := handler.NewFundsHandler(service.NewFundsService())
	marketHandler := handler.NewMarketHandler()
	worker.StartSettlementWorker(tradeService)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected Routes
		protected := api.Group("/")
		protected.Use(handler.AuthMiddleware())
		{
			protected.POST("/trade/order", tradeHandler.PlaceOrder)
			protected.GET("/trade/orders/active", tradeHandler.GetActiveOrders) // Add this
			protected.GET("/trade/orders/history", tradeHandler.GetOrderHistory)
			protected.GET("/wallet/balance", tradeHandler.GetBalance)
			protected.POST("/funds/deposit", fundsHandler.Deposit)
			protected.POST("/funds/withdraw", fundsHandler.Withdraw)
		}

		admin := api.Group("/admin")
		admin.Use(handler.AuthMiddleware(), handler.AdminOnly())
		{
			admin.POST("/funds/review/:id", fundsHandler.Review)
		}

		// Public market data
		api.GET("/market/candles", marketHandler.GetCandles)
		// Debug: current last price for a symbol
		api.GET("/market/price", marketHandler.GetPrice)
	}

	// WebSocket Route
	r.GET("/ws", handler.WSHandler)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"version": "1.0.0",
		})
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// 4. Start Server
	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
