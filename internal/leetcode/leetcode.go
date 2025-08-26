package leetcode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RecentACEntry struct {
	Title     string `json:"title"`
	TitleSlug string `json:"titleSlug"`
	Timestamp string `json:"timestamp"`
}

type RecentACData struct {
	RecentACSubmissionList []RecentACEntry `json:"recentAcSubmissionList"`
}

type RecentACListResponse struct {
	Data RecentACData `json:"data"`
}

type RecentAC struct {
	Title     string
	TitleSlug string
	Timestamp time.Time
}

func queryRecentACList(username string) ([]RecentACEntry, error) {
	// Call leetcode graphql api
	query := fmt.Sprintf(`{"query":"\n    query recentAcSubmissions($username: String!, $limit: Int!) {\n  recentAcSubmissionList(username: $username, limit: $limit) {\n    id\n    title\n    titleSlug\n    timestamp\n  }\n}\n    ","variables":{"username":"%s","limit":15},"operationName":"recentAcSubmissions"}`, username)
	resp, err := http.Post("https://leetcode.com/graphql", "application/json", strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent AC submissions: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch recent AC submissions: %s", resp.Status)
	}
	var result RecentACListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Data.RecentACSubmissionList, nil
}

func GetRecentACByUsername(username string) ([]RecentAC, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	recentACList, err := queryRecentACList(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent AC submissions for user %s: %w", username, err)
	}
	var result []RecentAC
	for _, ac := range recentACList {
		i, err := strconv.ParseInt(ac.Timestamp, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp %s: %w", ac.Timestamp, err)
		}
		timestamp := time.Unix(i, 0)
		result = append(result, RecentAC{
			Title:     ac.Title,
			TitleSlug: ac.TitleSlug,
			Timestamp: timestamp,
		})
	}
	return result, nil
}

