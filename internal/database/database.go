package database

import (
	"database/sql"
	"fmt"
	"log"

	"leetcode-telegram-bot/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the database connection
type DB struct {
	conn *sql.DB
}

// New creates a new database connection and initializes tables
func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// createTables creates all necessary database tables
func (db *DB) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS problems (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL,
			category TEXT NOT NULL,
			used BOOLEAN DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			username TEXT,
			first_name TEXT,
			last_name TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS submissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			problem_id INTEGER NOT NULL,
			submitted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			date TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id),
			FOREIGN KEY (problem_id) REFERENCES problems (id),
			UNIQUE(user_id, problem_id, date)
		)`,
		`CREATE TABLE IF NOT EXISTS daily_challenges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			problem_id INTEGER NOT NULL,
			date TEXT NOT NULL UNIQUE,
			posted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			day_number INTEGER NOT NULL,
			FOREIGN KEY (problem_id) REFERENCES problems (id)
		)`,
		`CREATE TABLE IF NOT EXISTS user_leetcode_profiles (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    user_id INTEGER NOT NULL,
		    leetcode_username TEXT NOT NULL UNIQUE,
		    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (user_id) REFERENCES users (id)
		)`,
		`CREATE TABLE IF NOT EXISTS challenge_counter (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			current_day INTEGER NOT NULL DEFAULT 9,
			last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := db.conn.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %q: %w", query, err)
		}
	}

	// Initialize challenge counter if it doesn't exist
	_, err := db.conn.Exec(`INSERT OR IGNORE INTO challenge_counter (id, current_day) VALUES (1, 9)`)
	if err != nil {
		return fmt.Errorf("failed to initialize challenge counter: %w", err)
	}

	return nil
}

// AddProblem adds a new problem to the database
func (db *DB) AddProblem(problem *models.Problem) error {
	query := `INSERT OR IGNORE INTO problems (title, url, category) VALUES (?, ?, ?)`
	_, err := db.conn.Exec(query, problem.Title, problem.URL, problem.Category)
	return err
}

// GetRandomUnusedProblem gets a random unused problem
func (db *DB) GetRandomUnusedProblem() (*models.Problem, error) {
	query := `SELECT id, title, url, category FROM problems WHERE used = FALSE ORDER BY RANDOM() LIMIT 1`
	row := db.conn.QueryRow(query)

	var problem models.Problem
	err := row.Scan(&problem.ID, &problem.Title, &problem.URL, &problem.Category)
	if err != nil {
		return nil, err
	}

	return &problem, nil
}

// MarkProblemAsUsed marks a problem as used
func (db *DB) MarkProblemAsUsed(problemID int) error {
	query := `UPDATE problems SET used = TRUE WHERE id = ?`
	_, err := db.conn.Exec(query, problemID)
	return err
}

// AddUser adds or updates a user in the database
func (db *DB) AddUser(user *models.User) error {
	query := `INSERT OR REPLACE INTO users (id, username, first_name, last_name, created_at) 
			  VALUES (?, ?, ?, ?, COALESCE((SELECT created_at FROM users WHERE id = ?), CURRENT_TIMESTAMP))`
	_, err := db.conn.Exec(query, user.ID, user.Username, user.FirstName, user.LastName, user.ID)
	return err
}

// AddSubmission adds a new submission
func (db *DB) AddSubmission(submission *models.Submission) error {
	query := `INSERT OR IGNORE INTO submissions (user_id, problem_id, date) VALUES (?, ?, ?)`
	_, err := db.conn.Exec(query, submission.UserID, submission.ProblemID, submission.Date)
	return err
}

// HasUserSubmittedToday checks if a user has submitted today
func (db *DB) HasUserSubmittedToday(userID int64, date string) (bool, error) {
	query := `SELECT COUNT(*) FROM submissions WHERE user_id = ? AND date = ?`
	row := db.conn.QueryRow(query, userID, date)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// AddDailyChallenge adds a new daily challenge
func (db *DB) AddDailyChallenge(challenge *models.DailyChallenge) error {
	query := `INSERT OR IGNORE INTO daily_challenges (problem_id, date, day_number) VALUES (?, ?, ?)`
	_, err := db.conn.Exec(query, challenge.ProblemID, challenge.Date, challenge.DayNumber)
	return err
}

// GetCurrentDayNumber gets the current day number
func (db *DB) GetCurrentDayNumber() (int, error) {
	query := `SELECT current_day FROM challenge_counter WHERE id = 1`
	row := db.conn.QueryRow(query)

	var dayNumber int
	err := row.Scan(&dayNumber)
	if err != nil {
		return 9, err // Default to day 9 if error
	}

	return dayNumber, nil
}

// IncrementDayNumber increments the day number and returns the new value
func (db *DB) IncrementDayNumber() (int, error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Get current day number
	var currentDay int
	err = tx.QueryRow(`SELECT current_day FROM challenge_counter WHERE id = 1`).Scan(&currentDay)
	if err != nil {
		return 0, err
	}

	// Increment and update
	newDay := currentDay + 1
	_, err = tx.Exec(`UPDATE challenge_counter SET current_day = ?, last_updated = CURRENT_TIMESTAMP WHERE id = 1`, newDay)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return newDay, nil
}

// ResetDayNumber resets the day number back to 8 (so next challenge will be Day 9)
func (db *DB) ResetDayNumber() error {
	query := `UPDATE challenge_counter SET current_day = 8, last_updated = CURRENT_TIMESTAMP WHERE id = 1`
	_, err := db.conn.Exec(query)
	return err
}

// GetTodaysChallenge gets today's challenge
func (db *DB) GetTodaysChallenge(date string) (*models.Problem, error) {
	query := `SELECT p.id, p.title, p.url, p.category 
			  FROM problems p 
			  JOIN daily_challenges dc ON p.id = dc.problem_id 
			  WHERE dc.date = ?`
	row := db.conn.QueryRow(query, date)

	var problem models.Problem
	err := row.Scan(&problem.ID, &problem.Title, &problem.URL, &problem.Category)
	if err != nil {
		return nil, err
	}

	return &problem, nil
}

// GetTodaysChallengeWithDay gets today's challenge with day number
func (db *DB) GetTodaysChallengeWithDay(date string) (*models.Problem, int, error) {
	query := `SELECT p.id, p.title, p.url, p.category, dc.day_number
			  FROM problems p 
			  JOIN daily_challenges dc ON p.id = dc.problem_id 
			  WHERE dc.date = ?`
	row := db.conn.QueryRow(query, date)

	var problem models.Problem
	var dayNumber int
	err := row.Scan(&problem.ID, &problem.Title, &problem.URL, &problem.Category, &dayNumber)
	if err != nil {
		return nil, 0, err
	}

	return &problem, dayNumber, nil
}

// GetLeaderboard gets the leaderboard with user statistics
func (db *DB) GetLeaderboard(limit int) ([]models.LeaderboardEntry, error) {
	query := `SELECT u.id, u.username, u.first_name, u.last_name, COUNT(s.id) as total_solved
			  FROM users u
			  LEFT JOIN submissions s ON u.id = s.user_id
			  GROUP BY u.id, u.username, u.first_name, u.last_name
			  ORDER BY total_solved DESC, u.first_name ASC
			  LIMIT ?`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []models.LeaderboardEntry
	for rows.Next() {
		var entry models.LeaderboardEntry
		err := rows.Scan(&entry.UserID, &entry.Username, &entry.FirstName, &entry.LastName, &entry.TotalSolved)
		if err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, entry)
	}

	return leaderboard, nil
}

// GetUsersWhoDidntSubmitToday gets users who haven't submitted today
func (db *DB) GetUsersWhoDidntSubmitToday(date string) ([]models.User, error) {
	query := `SELECT u.id, u.username, u.first_name, u.last_name
			  FROM users u
			  WHERE u.id NOT IN (
				  SELECT DISTINCT user_id FROM submissions WHERE date = ?
			  )`

	rows, err := db.conn.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// LoadProblemsFromYAML loads problems from YAML file into database
func (db *DB) LoadProblemsFromYAML(problems models.ProblemsData) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for category, problemList := range problems {
		for _, problem := range problemList {
			_, err := tx.Exec(
				`INSERT OR IGNORE INTO problems (title, url, category) VALUES (?, ?, ?)`,
				problem.Title, problem.URL, category,
			)
			if err != nil {
				log.Printf("Error inserting problem %s: %v", problem.Title, err)
				continue
			}
		}
	}

	return tx.Commit()
}

// GetLeetcodeProfile gets leetcode profile of user by id
func (db *DB) GetLeetcodeProfile(id int64) (*models.User, error) {
	query := `SELECT u.id, u.leetcode_username
			  FROM user_leetcode_profiles u
			  WHERE u.user_id = ?`

	rows, err := db.conn.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, fmt.Errorf("no leetcode profile found for user with id %d", id)
}

// RegisterLeetcodeProfile registers a leetcode profile for a user
func (db *DB) RegisterLeetcodeProfile(userID int64, leetcodeUsername string) error {
	query := `INSERT OR REPLACE INTO user_leetcode_profiles (user_id, leetcode_username, created_at) 
			  VALUES (?, ?, COALESCE((SELECT created_at FROM user_leetcode_profiles WHERE user_id = ?), CURRENT_TIMESTAMP))`
	_, err := db.conn.Exec(query, userID, leetcodeUsername, userID)
	return err
}
