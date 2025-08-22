package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"leetcode-telegram-bot/internal/bot"
	"leetcode-telegram-bot/internal/config"
	"leetcode-telegram-bot/internal/database"
	"leetcode-telegram-bot/internal/scheduler"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize bot
	telegramBot, err := bot.New(cfg.TelegramBotToken, db, cfg)
	if err != nil {
		log.Fatal("Failed to initialize bot:", err)
	}

	// Initialize scheduler
	scheduler := scheduler.New(telegramBot, db, cfg)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start bot
	go telegramBot.Start(ctx)

	// Start scheduler
	go scheduler.Start()

	log.Println("LeetCode Telegram Bot started successfully!")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	cancel()
	scheduler.Stop()
} 