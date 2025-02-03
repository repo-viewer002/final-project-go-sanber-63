package books

import (
	"time"
)

type Book struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Authors      string    `json:"authors"`
	Publisher    string    `json:"publisher"`
	Publish_Year uint      `json:"publish_year"`
	Stock        uint      `json:"stock"`
	Borrowed     uint      `json:"borrowed"`
	Genres       []string  `json:"genres"`
	Created_At   time.Time `json:"created_at"`
	Created_By   string    `json:"created_by"`
	Modified_At  time.Time `json:"modified_at"`
	Modified_By  string    `json:"modified_by"`
}

type SearchBook struct {
	Name              string `json:"name"`
	Authors           string `json:"authors"`
	Publisher         string `json:"publisher"`
	Publish_Year      string `json:"publish_Year"`
	Genre_Search_Type string `json:"genre_search_type"`
	Genres            []string `json:"genres"`
}
