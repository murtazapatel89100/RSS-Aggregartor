package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/murtazapatel89100/RSS-Aggregartor/internal/database"
)

func (config ApiConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	user, ok := GetUserFromContext(r)
	if !ok {
		RespondWithError(w, 403, "User not found in context")
		return
	}

	feed, err := config.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   params.Name,
		Url:    params.Url,
		UserID: user.ID,
	})
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Failed to create feed: %v", err))
		return
	}

	RespondWithJSON(w, 201, feed)
}

func (config ApiConfig) HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feed, err := config.DB.GetFeeds(r.Context())
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Failed to get feeds: %v", err))
		return
	}

	RespondWithJSON(w, 200, feed)

}

func (config ApiConfig) HandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	user, ok := GetUserFromContext(r)
	if !ok {
		RespondWithError(w, 403, "User not found in context")
		return
	}

	feedFollow, err := config.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Failed to create feed follow: %v", err))
		return
	}

	RespondWithJSON(w, 201, feedFollow)
}
