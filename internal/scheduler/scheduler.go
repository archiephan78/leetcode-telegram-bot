package scheduler

import (
	"io/ioutil"
	"log"
	"time"

	"leetcode-telegram-bot/internal/bot"
	"leetcode-telegram-bot/internal/config"
	"leetcode-telegram-bot/internal/database"
	"leetcode-telegram-bot/internal/models"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
)

// Scheduler handles scheduled tasks
type Scheduler struct {
	cron   *cron.Cron
	bot    *bot.Bot
	db     *database.DB
	config *config.Config
}

// New creates a new scheduler instance
func New(bot *bot.Bot, db *database.DB, cfg *config.Config) *Scheduler {
	// Create cron with timezone support
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		log.Printf("Failed to load timezone %s, using UTC: %v", cfg.Timezone, err)
		loc = time.UTC
	}

	c := cron.New(cron.WithLocation(loc))

	return &Scheduler{
		cron:   c,
		bot:    bot,
		db:     db,
		config: cfg,
	}
}

// Start starts the scheduler with all cron jobs
func (s *Scheduler) Start() {
	// Load problems from YAML file first
	if err := s.loadProblemsFromFile(); err != nil {
		log.Printf("Warning: Failed to load problems from file: %v", err)
	}

	// Schedule daily challenge posting at 7:00 AM, Monday to Friday only
	_, err := s.cron.AddFunc("0 7 * * 1-5", func() {
		log.Println("Posting daily challenge...")
		if err := s.bot.PostDailyChallenge(); err != nil {
			log.Printf("Error posting daily challenge: %v", err)
		}
	})
	if err != nil {
		log.Printf("Error scheduling daily challenge: %v", err)
	}

	// Schedule afternoon reminder at 3:00 PM, Monday to Friday only
	_, err = s.cron.AddFunc("0 15 * * 1-5", func() {
		log.Println("Sending afternoon reminder...")
		if err := s.bot.SendReminder(); err != nil {
			log.Printf("Error sending afternoon reminder: %v", err)
		}
	})
	if err != nil {
		log.Printf("Error scheduling afternoon reminder: %v", err)
	}

	// Schedule evening reminder at 10:00 PM, Monday to Friday only
	_, err = s.cron.AddFunc("0 22 * * 1-5", func() {
		log.Println("Sending evening reminder...")
		if err := s.bot.SendReminder(); err != nil {
			log.Printf("Error sending evening reminder: %v", err)
		}
	})
	if err != nil {
		log.Printf("Error scheduling evening reminder: %v", err)
	}

	// Schedule check submissions every 5 minutes
	_, err = s.cron.AddFunc("*/5 * * * *", func() {
		log.Println("Checking new submissions...")
		if err := s.bot.CheckSubmissions(); err != nil {
			log.Printf("Error while checking new submissions %v", err)
		}
	})
	if err != nil {
		log.Printf("Error scheduling check new submissions: %v", err)
	}

	// Start the cron scheduler
	s.cron.Start()
	log.Println("Scheduler started successfully - posting challenges Monday to Friday only")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

// loadProblemsFromFile loads problems from the YAML file into the database
func (s *Scheduler) loadProblemsFromFile() error {
	// Read the YAML file
	data, err := ioutil.ReadFile(s.config.ProblemsFilePath)
	if err != nil {
		return err
	}

	// Parse the YAML data
	var problemsData models.ProblemsData
	if err := yaml.Unmarshal(data, &problemsData); err != nil {
		return err
	}

	// Load problems into database
	if err := s.db.LoadProblemsFromYAML(problemsData); err != nil {
		return err
	}

	log.Printf("Successfully loaded problems from %s", s.config.ProblemsFilePath)
	return nil
}

// GetNextScheduledTimes returns information about next scheduled tasks (for debugging)
func (s *Scheduler) GetNextScheduledTimes() []time.Time {
	entries := s.cron.Entries()
	var nextTimes []time.Time

	for _, entry := range entries {
		nextTimes = append(nextTimes, entry.Next)
	}

	return nextTimes
}
