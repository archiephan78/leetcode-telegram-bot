package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"leetcode-telegram-bot/internal/config"
	"leetcode-telegram-bot/internal/database"
	"leetcode-telegram-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents the Telegram bot
type Bot struct {
	api    *tgbotapi.BotAPI
	db     *database.DB
	config *config.Config
}

// New creates a new Telegram bot instance
func New(token string, db *database.DB, cfg *config.Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	api.Debug = false
	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		api:    api,
		db:     db,
		config: cfg,
	}, nil
}

// Start starts the bot and handles incoming messages
func (b *Bot) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			log.Println("Bot stopping...")
			return
		case update := <-updates:
			if update.Message != nil {
				go b.handleMessage(update.Message)
			}
		}
	}
}

// handleMessage processes incoming messages
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	// Save user info
	user := &models.User{
		ID:        message.From.ID,
		Username:  message.From.UserName,
		FirstName: message.From.FirstName,
		LastName:  message.From.LastName,
	}
	if err := b.db.AddUser(user); err != nil {
		log.Printf("Error saving user: %v", err)
	}

	// Handle commands
	if message.IsCommand() {
		switch message.Command() {
		case "submit":
			b.handleSubmitCommand(message)
		case "leaderboards":
			b.handleLeaderboardCommand(message)
		case "help":
			b.handleHelpCommand(message)
		case "manual":
			b.handleManualCommand(message)
		case "testreminder":
			b.handleTestReminderCommand(message)
		case "status":
			b.handleStatusCommand(message)
		case "resetday":
			b.handleResetDayCommand(message)
		default:
			b.sendMessage(message.Chat.ID, "Unknown command. Use /help to see available commands.")
		}
	}
}

// handleSubmitCommand handles the /submit command
func (b *Bot) handleSubmitCommand(message *tgbotapi.Message) {
	today := time.Now().Format("2006-01-02")

	// Check if user already submitted today
	hasSubmitted, err := b.db.HasUserSubmittedToday(message.From.ID, today)
	if err != nil {
		log.Printf("Error checking submission: %v", err)
		b.sendMessage(message.Chat.ID, "❌ An error occurred while checking your submission.")
		return
	}

	if hasSubmitted {
		b.sendMessage(message.Chat.ID, "✅ You have already submitted today's challenge!")
		return
	}

	// Get today's challenge with day number
	todaysChallenge, dayNumber, err := b.db.GetTodaysChallengeWithDay(today)
	if err != nil {
		log.Printf("Error getting today's challenge: %v", err)
		b.sendMessage(message.Chat.ID, "❌ No challenge available for today yet.")
		return
	}

	// Add submission
	submission := &models.Submission{
		UserID:    message.From.ID,
		ProblemID: todaysChallenge.ID,
		Date:      today,
	}

	if err := b.db.AddSubmission(submission); err != nil {
		log.Printf("Error adding submission: %v", err)
		b.sendMessage(message.Chat.ID, "❌ An error occurred while submitting.")
		return
	}

	responseText := fmt.Sprintf("🎉 Great job! You've successfully submitted Day %d challenge:\n\n"+
		"📝 **%s**\n"+
		"🔗 %s\n\n"+
		"Keep up the good work! 💪", dayNumber, todaysChallenge.Title, todaysChallenge.URL)

	b.sendMessage(message.Chat.ID, responseText)
}

// handleLeaderboardCommand handles the /leaderboards command
func (b *Bot) handleLeaderboardCommand(message *tgbotapi.Message) {
	leaderboard, err := b.db.GetLeaderboard(10)
	if err != nil {
		log.Printf("Error getting leaderboard: %v", err)
		b.sendMessage(message.Chat.ID, "❌ An error occurred while fetching the leaderboard.")
		return
	}

	if len(leaderboard) == 0 {
		b.sendMessage(message.Chat.ID, "📊 No submissions yet! Be the first to submit a challenge.")
		return
	}

	var responseText strings.Builder
	responseText.WriteString("🏆 **LeetCode Challenge Leaderboard** 🏆\n\n")

	for i, entry := range leaderboard {
		var emoji string
		switch i {
		case 0:
			emoji = "🥇"
		case 1:
			emoji = "🥈"
		case 2:
			emoji = "🥉"
		default:
			emoji = fmt.Sprintf("%d.", i+1)
		}

		name := entry.FirstName
		if entry.LastName != "" {
			name += " " + entry.LastName
		}
		if entry.Username != "" {
			name += fmt.Sprintf(" (@%s)", entry.Username)
		}

		responseText.WriteString(fmt.Sprintf("%s %s - %d solved\n", emoji, name, entry.TotalSolved))
	}

	responseText.WriteString("\n💪 Keep solving to climb the ranks!")

	b.sendMessage(message.Chat.ID, responseText.String())
}

// handleHelpCommand handles the /help command
func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	helpText := `🤖 **LeetCode Challenge Bot Help**

Available commands:
• /submit - Submit today's challenge
• /leaderboards - View the leaderboard
• /status - Show bot status and current day info
• /help - Show this help message

**Admin Commands (Group only):**
• /manual - Manually post daily challenge immediately
• /testreminder - Test reminder functionality
• /resetday - Reset day counter (next challenge will be Day 9)

📅 **How it works:**
- Every weekday (Monday to Friday) at 7:00 AM, I post a new LeetCode challenge
- Challenge numbering starts from Day 9
- Use /submit to mark that you've completed it
- Check /leaderboards to see who's solving the most problems
- I'll remind you at 3:00 PM and 10:00 PM if you haven't submitted yet
- No challenges on weekends (Saturday & Sunday) 🎉

Happy coding! 💻✨`

	b.sendMessage(message.Chat.ID, helpText)
}

// handleManualCommand handles the /manual command for manually posting daily challenge
func (b *Bot) handleManualCommand(message *tgbotapi.Message) {
	// Check if user is admin (you can customize this logic)
	if message.From.ID != b.config.TelegramGroupID && message.Chat.ID != b.config.TelegramGroupID {
		b.sendMessage(message.Chat.ID, "❌ This command can only be used in the main group.")
		return
	}

	b.sendMessage(message.Chat.ID, "📝 Manually posting daily challenge...")

	if err := b.PostDailyChallenge(); err != nil {
		log.Printf("Error in manual command: %v", err)
		b.sendMessage(message.Chat.ID, fmt.Sprintf("❌ Error posting manual challenge: %v", err))
	} else {
		b.sendMessage(message.Chat.ID, "✅ Manual challenge posted successfully!")
	}
}

// handleTestReminderCommand handles the /testreminder command for testing reminders
func (b *Bot) handleTestReminderCommand(message *tgbotapi.Message) {
	// Check if user is admin (you can customize this logic)
	if message.From.ID != b.config.TelegramGroupID && message.Chat.ID != b.config.TelegramGroupID {
		b.sendMessage(message.Chat.ID, "❌ This command can only be used in the main group.")
		return
	}

	b.sendMessage(message.Chat.ID, "🧪 Testing reminder...")

	if err := b.SendReminder(); err != nil {
		log.Printf("Error in test reminder command: %v", err)
		b.sendMessage(message.Chat.ID, fmt.Sprintf("❌ Error sending test reminder: %v", err))
	} else {
		b.sendMessage(message.Chat.ID, "✅ Test reminder sent successfully!")
	}
}

// handleStatusCommand handles the /status command for showing bot status
func (b *Bot) handleStatusCommand(message *tgbotapi.Message) {
	today := time.Now().Format("2006-01-02")

	// Get current day number
	currentDay, err := b.db.GetCurrentDayNumber()
	if err != nil {
		log.Printf("Error getting current day: %v", err)
		currentDay = 0
	}

	// Check if there's a challenge today
	todaysChallenge, dayNumber, err := b.db.GetTodaysChallengeWithDay(today)
	var challengeStatus string
	if err != nil {
		challengeStatus = "❌ No challenge posted today"
	} else {
		challengeStatus = fmt.Sprintf("✅ Day %d: %s", dayNumber, todaysChallenge.Title)
	}

	// Get leaderboard summary (top 3)
	leaderboard, err := b.db.GetLeaderboard(3)
	var leaderboardStatus string
	if err != nil || len(leaderboard) == 0 {
		leaderboardStatus = "No submissions yet"
	} else {
		leaderboardStatus = fmt.Sprintf("Top: %s (%d solved)", leaderboard[0].FirstName, leaderboard[0].TotalSolved)
	}

	// Get users who haven't submitted today
	usersNotSubmitted, err := b.db.GetUsersWhoDidntSubmitToday(today)
	var submissionStatus string
	if err != nil {
		submissionStatus = "Error checking submissions"
	} else {
		submissionStatus = fmt.Sprintf("%d users haven't submitted today", len(usersNotSubmitted))
	}

	statusText := fmt.Sprintf("🤖 **Bot Status** 🤖\n\n"+
		"📅 Date: %s\n"+
		"📊 Current Day Counter: %d\n"+
		"🎯 Today's Challenge: %s\n"+
		"📈 Leaderboard: %s\n"+
		"📝 Submissions: %s\n\n"+
		"⏰ Next challenge: Tomorrow 7:00 AM (Mon-Fri only)\n"+
		"🎉 Weekend: No challenges",
		time.Now().Format("January 2, 2006"),
		currentDay,
		challengeStatus,
		leaderboardStatus,
		submissionStatus)

	b.sendMessage(message.Chat.ID, statusText)
}

// handleResetDayCommand handles the /resetday command
func (b *Bot) handleResetDayCommand(message *tgbotapi.Message) {
	// Check if user is admin (you can customize this logic)
	if message.From.ID != b.config.TelegramGroupID && message.Chat.ID != b.config.TelegramGroupID {
		b.sendMessage(message.Chat.ID, "❌ This command can only be used in the main group.")
		return
	}

	if err := b.db.ResetDayNumber(); err != nil {
		log.Printf("Error resetting day number: %v", err)
		b.sendMessage(message.Chat.ID, "❌ An error occurred while resetting the day counter.")
	} else {
		log.Println("Day counter reset successfully.")
		b.sendMessage(message.Chat.ID, "✅ Day counter reset successfully! Next challenge will be Day 9.")
	}
}

// sendMessage sends a message to a chat
func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// PostDailyChallenge posts the daily challenge to the group
func (b *Bot) PostDailyChallenge() error {
	// Get a random unused problem
	problem, err := b.db.GetRandomUnusedProblem()
	if err != nil {
		return fmt.Errorf("failed to get random problem: %w", err)
	}

	// Mark problem as used
	if err := b.db.MarkProblemAsUsed(problem.ID); err != nil {
		return fmt.Errorf("failed to mark problem as used: %w", err)
	}

	// Get and increment day number
	dayNumber, err := b.db.IncrementDayNumber()
	if err != nil {
		return fmt.Errorf("failed to increment day number: %w", err)
	}

	// Add to daily challenges
	today := time.Now().Format("2006-01-02")
	challenge := &models.DailyChallenge{
		ProblemID: problem.ID,
		Date:      today,
		DayNumber: dayNumber,
	}

	if err := b.db.AddDailyChallenge(challenge); err != nil {
		return fmt.Errorf("failed to add daily challenge: %w", err)
	}

	// Create message with day number
	messageText := fmt.Sprintf("🌅 **Daily LeetCode Challenge - Day %d** 🌅\n"+
		"📅 %s\n\n"+
		"📝 **%s**\n"+
		"🏷️ Category: %s\n"+
		"🔗 %s\n\n"+
		"💪 Ready to solve it? Use /submit when you're done!\n"+
		"Good luck everyone! 🍀",
		dayNumber,
		time.Now().Format("January 2, 2006"),
		problem.Title,
		problem.Category,
		problem.URL)

	// Send to group
	b.sendMessage(b.config.TelegramGroupID, messageText)

	log.Printf("Posted daily challenge Day %d: %s", dayNumber, problem.Title)
	return nil
}

// SendReminder sends a reminder to users who haven't submitted
func (b *Bot) SendReminder() error {
	today := time.Now().Format("2006-01-02")

	// Get users who haven't submitted today
	users, err := b.db.GetUsersWhoDidntSubmitToday(today)
	if err != nil {
		return fmt.Errorf("failed to get users who didn't submit: %w", err)
	}

	if len(users) == 0 {
		log.Println("All users have submitted today!")
		return nil
	}

	// Get today's challenge with day number
	todaysChallenge, dayNumber, err := b.db.GetTodaysChallengeWithDay(today)
	if err != nil {
		return fmt.Errorf("failed to get today's challenge: %w", err)
	}

	// Create reminder message with mentions
	var mentions []string
	for _, user := range users {
		if user.Username != "" {
			mentions = append(mentions, "@"+user.Username)
		} else {
			mentions = append(mentions, user.FirstName)
		}
	}

	currentHour := time.Now().Hour()
	var reminderEmoji string
	var reminderTime string

	if currentHour == 15 {
		reminderEmoji = "⏰"
		reminderTime = "Afternoon"
	} else {
		reminderEmoji = "🌙"
		reminderTime = "Evening"
	}

	messageText := fmt.Sprintf("%s **%s Reminder** %s\n\n"+
		"Hey %s!\n\n"+
		"Don't forget about today's LeetCode challenge (Day %d):\n"+
		"📝 **%s**\n"+
		"🔗 %s\n\n"+
		"Use /submit when you're done! ⚡",
		reminderEmoji, reminderTime, reminderEmoji,
		strings.Join(mentions, ", "),
		dayNumber,
		todaysChallenge.Title,
		todaysChallenge.URL)

	// Send to group
	b.sendMessage(b.config.TelegramGroupID, messageText)

	log.Printf("Sent reminder to %d users for Day %d", len(users), dayNumber)
	return nil
}
