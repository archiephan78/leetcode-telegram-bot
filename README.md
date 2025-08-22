# LeetCode Telegram Bot

A Telegram bot written in Go to manage daily LeetCode challenges in Telegram groups.

## Features

- 🌅 **Daily Challenge**: Automatically posts a random LeetCode challenge at 7 AM on weekdays (Monday-Friday)
- 📊 **Day Counter**: Counts challenge days starting from Day 9
- 📝 **Submit Command**: Allows users to submit when they complete a challenge
- 🏆 **Leaderboard**: Displays rankings based on the number of completed challenges
- ⏰ **Reminders**: Automatically reminds users who haven't submitted at 3 PM and 10 PM (weekdays only)
- 🎉 **Weekend Break**: No challenges on Saturday and Sunday
- 🗄️ **SQLite Database**: Stores user information, challenges, and submissions
- 🐳 **Docker Support**: Easy deployment with Docker

## Commands

- `/submit` - Submit today's challenge
- `/leaderboards` - View the leaderboard
- `/help` - Display help information

## Setup

### Requirements

- Go 1.21+
- SQLite3
- Telegram Bot Token
- Telegram Group ID

### Step 1: Clone repository

```bash
git clone <repository-url>
cd leetcode-telegram-bot
```

### Step 2: Configure environment variables

Copy the `env.example` file to `.env` and update the values:

```bash
cp env.example .env
```

Edit the `.env` file:

```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_GROUP_ID=your_group_id_here
DATABASE_PATH=leetcode_bot.db
PROBLEMS_FILE_PATH=problem_deduplicated.yaml
TIMEZONE=Asia/Ho_Chi_Minh
```

### Step 3: Run with Docker (Recommended)

**Option 1: Docker Compose (Easiest)**
```bash
# Create .env file
echo "TELEGRAM_BOT_TOKEN=your_token_here" > .env
echo "TELEGRAM_GROUP_ID=your_group_id" >> .env

# Build and run
docker-compose up -d

# View logs
docker-compose logs -f
```

**Option 2: Multi-Architecture Build**
```bash
# For Apple Silicon (ARM) and Intel (AMD64)
chmod +x build-multiarch.sh
./build-multiarch.sh

# Or ARM64 only
chmod +x build-arm.sh
./build-arm.sh
```

**Option 3: Manual Docker Build**
```bash
# Build for ARM64 (Apple Silicon)
docker build --platform linux/arm64 -t leetcode-telegram-bot .

# Run
docker run -d --name leetcode-bot \
  -e TELEGRAM_BOT_TOKEN="your_token" \
  -e TELEGRAM_GROUP_ID="your_group_id" \
  -v $(pwd)/data:/data \
  leetcode-telegram-bot
```

### Step 4: Run directly

```bash
# Download dependencies
go mod download

# Build
go build -o main .

# Run
./main
```

## Project Structure

```
leetcode-telegram-bot/
├── main.go                    # Entry point
├── internal/
│   ├── bot/                   # Telegram bot logic
│   │   └── bot.go
│   ├── config/                # Configuration management
│   │   └── config.go
│   ├── database/              # Database operations
│   │   └── database.go
│   ├── models/                # Data models
│   │   └── models.go
│   └── scheduler/             # Cron job scheduler
│       └── scheduler.go
├── problem_deduplicated.yaml  # LeetCode problems data
├── Dockerfile                 # Docker configuration
├── docker-compose.yml         # Docker Compose configuration
└── README.md                  # This file
```

## How to get Telegram Bot Token and Group ID

### Bot Token

1. Create a new bot with [@BotFather](https://t.me/BotFather)
2. Send `/newbot` and follow the instructions
3. Save the provided token

### Group ID

1. Add the bot to your group
2. Send any message in the group
3. Visit: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. Look for `"chat":{"id":` in the response - that's your Group ID

## Database Schema

The bot uses SQLite with the following tables:

- `problems`: Stores LeetCode problems
- `users`: Telegram user information
- `submissions`: User submissions
- `daily_challenges`: Daily challenges with day counter
- `challenge_counter`: Stores the current day number (starting from 9)

## Cron Jobs

- **07:00 (Mon-Fri)**: Post daily challenge (starting from Day 9)
- **15:00 (Mon-Fri)**: Afternoon reminder
- **22:00 (Mon-Fri)**: Evening reminder
- **Weekend**: No challenges posted

## Development

### Adding new problems

Edit the `problem_deduplicated.yaml` file following the existing format:

```yaml
# Add your problems here following the existing structure
```