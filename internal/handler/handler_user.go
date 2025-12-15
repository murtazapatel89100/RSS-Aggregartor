package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/murtazapatel89100/RSS-Aggregartor/internal/database"
)

func (config ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode((&params))
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}

	user, err := config.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:       uuid.New(),
		Username: params.Name,
	})
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	RespondWithJSON(w, 201, user)
}

func (config ApiConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		RespondWithError(w, 403, "User not found in context")
		return
	}

	RespondWithJSON(w, 200, user)
}

func (config ApiConfig) HandlerGetUserFeeds(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		RespondWithError(w, 403, "User not found in context")
		return
	}

	posts, err := config.DB.GetPostForUser(r.Context(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(10),
		Offset: int32(0),
	})
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Failed to get user posts: %v", err))
		return
	}

	// Return empty array instead of null if no posts found
	if posts == nil {
		posts = []database.Post{}
	}

	RespondWithJSON(w, 200, posts)
}
