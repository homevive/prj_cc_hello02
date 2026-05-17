package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	githubAPI = "https://api.github.com"
	repoOwner = "larksuite"
	repoName  = "cli"
)

func getToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

// GitHub API 返回的 issue/PR 结构
type ghIssue struct {
	ID        int       `json:"id"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	User      *ghUser   `json:"user"`
	Labels    []ghLabel `json:"labels"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PullReq   *struct{} `json:"pull_request"`
}

type ghUser struct {
	Login string `json:"login"`
}

type ghLabel struct {
	Name string `json:"name"`
}

func fetchGH(path string) ([]ghIssue, error) {
	url := fmt.Sprintf("%s%s", githubAPI, path)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+getToken())
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "prj_cc_hello02")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("fetch %s: HTTP %d", path, resp.StatusCode)
	}

	var items []ghIssue
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return items, nil
}

func syncIssues() (int, error) {
	items, err := fetchGH(fmt.Sprintf("/repos/%s/%s/issues?state=all&per_page=100&filter=all", repoOwner, repoName))
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range items {
		if item.PullReq != nil {
			continue // 跳过 PR（GitHub Issues API 也返回 PR）
		}
		labels := ""
		for i, l := range item.Labels {
			if i > 0 {
				labels += ","
			}
			labels += l.Name
		}
		author := ""
		if item.User != nil {
			author = item.User.Login
		}
		if err := upsertItem(GitHubItem{
			ID:        int64(item.ID),
			Number:    item.Number,
			Title:     item.Title,
			State:     item.State,
			ItemType:  "issue",
			Author:    author,
			Labels:    labels,
			URL:       item.HTMLURL,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}); err != nil {
			log.Printf("[INFO] upsert issue #%d failed: %v\n", item.Number, err)
		} else {
			count++
		}
	}
	return count, nil
}

func syncPRs() (int, error) {
	items, err := fetchGH(fmt.Sprintf("/repos/%s/%s/pulls?state=all&per_page=100", repoOwner, repoName))
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range items {
		labels := ""
		for i, l := range item.Labels {
			if i > 0 {
				labels += ","
			}
			labels += l.Name
		}
		author := ""
		if item.User != nil {
			author = item.User.Login
		}

		state := item.State
		if item.State == "closed" {
			// 进一步检查是否已合并
			if merged, _ := checkPRMerged(item.Number); merged {
				state = "merged"
			}
		}

		if err := upsertItem(GitHubItem{
			ID:        int64(item.ID),
			Number:    item.Number,
			Title:     item.Title,
			State:     state,
			ItemType:  "pr",
			Author:    author,
			Labels:    labels,
			URL:       item.HTMLURL,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}); err != nil {
			log.Printf("[INFO] upsert PR #%d failed: %v\n", item.Number, err)
		} else {
			count++
		}
	}
	return count, nil
}

func checkPRMerged(number int) (bool, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/merge", githubAPI, repoOwner, repoName, number)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+getToken())
	req.Header.Set("User-Agent", "prj_cc_hello02")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == 204, nil
}

func syncAll() {
	log.Printf("[INFO] Syncing GitHub data from %s/%s...\n", repoOwner, repoName)
	appendLog(fmt.Sprintf("Syncing from %s/%s...", repoOwner, repoName), 0)

	issueCount, err := syncIssues()
	if err != nil {
		log.Printf("[INFO] Sync issues error: %v\n", err)
		appendLog(fmt.Sprintf("Sync issues failed: %v", err), 0)
	} else {
		log.Printf("[INFO] Synced %d issues\n", issueCount)
		appendLog(fmt.Sprintf("Synced %d issues", issueCount), issueCount)
	}

	prCount, err := syncPRs()
	if err != nil {
		log.Printf("[INFO] Sync PRs error: %v\n", err)
		appendLog(fmt.Sprintf("Sync PRs failed: %v", err), 0)
	} else {
		log.Printf("[INFO] Synced %d PRs\n", prCount)
		appendLog(fmt.Sprintf("Synced %d PRs", prCount), prCount)
	}

	logStats()
}
