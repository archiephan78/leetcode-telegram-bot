package models

import (
	"time"
)

// Problem represents a LeetCode problem
type Problem struct {
	ID       int    `json:"id" db:"id"`
	Title    string `json:"title" db:"title"`
	URL      string `json:"url" db:"url"`
	Category string `json:"category" db:"category"`
	Used     bool   `json:"used" db:"used"`
}

// User represents a Telegram user
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserLeetcodeProfile represents a user's LeetCode profile
type UserLeetcodeProfile struct {
	ID          int64     `json:"id" db:"id"`
	UserId		 int64     `json:"user_id" db:"user_id"`
	LeetCodeUsername string    `json:"leetcode_username" db:"leetcode_username"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Submission represents a user's submission for a daily challenge
type Submission struct {
	ID          int       `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	ProblemID   int       `json:"problem_id" db:"problem_id"`
	SubmittedAt time.Time `json:"submitted_at" db:"submitted_at"`
	Date        string    `json:"date" db:"date"` // Format: YYYY-MM-DD
}

// DailyChallenge represents the daily challenge posted
type DailyChallenge struct {
	ID        int       `json:"id" db:"id"`
	ProblemID int       `json:"problem_id" db:"problem_id"`
	Date      string    `json:"date" db:"date"` // Format: YYYY-MM-DD
	PostedAt  time.Time `json:"posted_at" db:"posted_at"`
	DayNumber int       `json:"day_number" db:"day_number"` // Day counter (starting from 9)
}

// LeaderboardEntry represents a user's statistics for leaderboard
type LeaderboardEntry struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	TotalSolved int    `json:"total_solved"`
}

// ProblemsData represents the structure of the YAML file
type ProblemsData map[string][]struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

// ChallengeCounter represents the global challenge counter
type ChallengeCounter struct {
	ID          int       `json:"id" db:"id"`
	CurrentDay  int       `json:"current_day" db:"current_day"`
	LastUpdated time.Time `json:"last_updated" db:"last_updated"`
}
