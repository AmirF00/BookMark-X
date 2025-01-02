package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
)

type Tweet struct {
	SNum       int    `json:"SNum"`
	Handle     string `json:"Handle"`
	Name       string `json:"Name"`
	ProfilePic string `json:"Profile Pic"`
	TweetText  string `json:"Tweet Text"`
	TweetLink  string `json:"Tweet Link"`
}

type Summary struct {
	SNum    int    `json:"SNum"`
	Handle  string `json:"Handle"`
	Link    string `json:"Link"`
	Summary string `json:"Summary"`
	Draft   bool   `json:"Draft"`
}

type PageData struct {
	Title string
	Data  interface{}
}

var (
	summaries []Summary
	mu        sync.RWMutex
	templates *template.Template
)

func init() {
	var err error
	templates, err = template.New("").ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Template parsing error:", err)
	}
}

func fileServer() http.Handler {
	return http.FileServer(http.Dir("static"))
}

func loadTweets() ([]Tweet, error) {
	data, err := os.ReadFile("static/tweets.json")
	if err != nil {
		return nil, err
	}
	var tweets []Tweet
	if err := json.Unmarshal(data, &tweets); err != nil {
		return nil, err
	}
	return tweets, nil
}

func loadSummaries() error {
	data, err := os.ReadFile("static/summaries.json")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		summaries = []Summary{}
		return nil
	}
	return json.Unmarshal(data, &summaries)
}

func saveSummaries() error {
	data, err := json.MarshalIndent(summaries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("static/summaries.json", data, 0644)
}

func tweetsHandler(w http.ResponseWriter, r *http.Request) {
	tweets, err := loadTweets()
	if err != nil {
		log.Printf("Error loading tweets: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/tweets.html"))
	err = tmpl.Execute(w, PageData{
		Title: "Tweets",
		Data:  tweets,
	})
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func summaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		mu.Lock()
		defer mu.Unlock()

		var summary Summary
		if err := json.NewDecoder(r.Body).Decode(&summary); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		found := false
		for i := range summaries {
			if summaries[i].SNum == summary.SNum {
				summaries[i] = summary
				found = true
				break
			}
		}
		if !found {
			summaries = append(summaries, summary)
		}

		if err := saveSummaries(); err != nil {
			log.Printf("Error saving summaries: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/summary.html"))
	err := tmpl.Execute(w, PageData{
		Title: "Summaries",
		Data:  summaries,
	})
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	tmpl := template.Must(template.ParseFiles("templates/read.html"))
	err := tmpl.Execute(w, PageData{
		Title: "Read Summaries",
		Data:  summaries,
	})
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func tipsHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	tmpl := template.Must(template.ParseFiles("templates/tips.html"))
	err := tmpl.Execute(w, PageData{
		Title: "Security Tips",
		Data:  summaries,
	})
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	if err := loadSummaries(); err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", fileServer()))
	http.HandleFunc("/twitts", tweetsHandler)
	http.HandleFunc("/summary", summaryHandler)
	http.HandleFunc("/read", readHandler)
	http.HandleFunc("/tips", tipsHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
