package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func getTopStories() ([]int, error) {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Trims '[' and ']' from response
	body = body[1 : len(body)-1]
	ids_s := strings.Split(string(body), ",")
	ids := make([]int, len(ids_s))
	for i, s := range ids_s {
		ids[i], err = strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
	}

	return ids, err
}

type story struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

func fetchStory(id int) (story, error) {
	story := story{}

	resp, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id))
	if err != nil {
		return story, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return story, err
	}

	if err := json.Unmarshal(body, &story); err != nil {
		return story, err
	}

	return story, nil
}

func GetStories(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()

	if r.URL.Path != "/" {
		http.Error(w, "NOT FOUND", 404)
		return
	}

	// Gets top stories
	topStories, err := getTopStories()
	if err != nil {
		http.Error(w, "Failed to fetch the top stories", 500)
		return
	}

	// Get the stories
	stories := make([]story, 30)
	var storiesWg sync.WaitGroup

	for i := 0; i < 30; i++ {
		storiesWg.Add(1)

		go func() {
			defer storiesWg.Done()

			story, err := fetchStory(topStories[i])
			if err != nil {
				http.Error(w, "Couldn't fetch a story", http.StatusNotFound)
				return
			}

			stories[i] = story
		}()
	}

	storiesWg.Wait()

	// Sends html response
	tmpl := template.Must(template.ParseFiles("index.tmpl"))
	if err := tmpl.Execute(w, stories); err != nil {
		http.Error(w, "Couldn't parse the template", http.StatusInternalServerError)
		return
	}

	// Time logging
	fmt.Printf("It took %v to procces a request\n", time.Now().Sub(timeStart))
}
