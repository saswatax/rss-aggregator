package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/saswatax/rss-aggregator/internal/database"
	"github.com/saswatax/rss-aggregator/internal/server/services"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	payload, err := decode[request](r)
	if err != nil {
		log.Println(err)
		encode(w, http.StatusBadRequest, "invalid payload")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), s.Config.Env.BCRYPT_COST)
	if err != nil {
		log.Println(err)
		encode(w, http.StatusInternalServerError, "failed to generate password hash")
		return
	}

	user, err := s.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  string(passwordHash),
	})
	if err != nil {
		log.Println(err)
		encode(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	encode(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
	})
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	payload, err := decode[request](r)
	if err != nil {
		log.Println(err)
		encode(w, http.StatusBadRequest, "invalid payload")
		return
	}

	user, err := s.DB.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		log.Println(err)
		encode(w, http.StatusNotFound, fmt.Sprintf("no user found with %v", payload.Email))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		log.Println(err)
		encode(w, http.StatusUnauthorized, "wrong password")
		return
	}

	token, err := services.GenerateToken(user.ID, s.Config.Env.JWT_SECRET)
	if err != nil {
		log.Println(err)
		encode(w, http.StatusInternalServerError, "failed to generated authentication token")
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	encode(w, http.StatusOK, response{
		Token: token,
	})
}

func (s *Server) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ctxUserID{}).(uuid.UUID)
	if !ok {
		log.Println("userID not found")
		encode(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	user, err := s.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Println(err)
		encode(w, http.StatusNotFound, "no user found")
		return
	}

	encode(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
	})
}

func (s *Server) CreateFeed(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	payload, err := decode[request](r)
	if err != nil {
		log.Println(err)
		encode(w, http.StatusBadRequest, "invalid payload")
		return
	}

	userID, ok := r.Context().Value(ctxUserID{}).(uuid.UUID)
	if !ok {
		log.Println("userID not found")
		encode(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	feed, err := s.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      payload.Name,
		Url:       payload.URL,
		UserID:    userID,
	})
	if err != nil {
		log.Println(err)
		encode(w, http.StatusInternalServerError, "failed to create feed")
		return
	}

	encode(w, http.StatusCreated, Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		URL:       feed.Url,
		UserID:    feed.UserID,
	})
}

func (s *Server) GetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := s.DB.GetFeeds(r.Context())
	if err != nil {
		log.Println(err)
		encode(w, http.StatusInternalServerError, "failed to fetch feeds")
		return
	}

	response := make([]Feed, len(feeds))
	for _, feed := range feeds {
		response = append(response, Feed{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			URL:       feed.Url,
			UserID:    feed.UserID,
		})
	}

	encode(w, http.StatusOK, response)
}

func (s *Server) CreateFeedFollow(w http.ResponseWriter, r *http.Request) {
	type request struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	payload, err := decode[request](r)
	if err != nil {
		log.Println(err)
		encode(w, http.StatusBadRequest, "invalid payload")
		return
	}

	userID, ok := r.Context().Value(ctxUserID{}).(uuid.UUID)
	if !ok {
		log.Println("userID not found")
		encode(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	feedFollow, err := s.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userID,
		FeedID:    payload.FeedID,
	})
	if err != nil {
		log.Println(err)
		encode(w, http.StatusInternalServerError, "failed to create feed follow")
		return
	}

	encode(w, http.StatusCreated, FeedFollow{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserID:    feedFollow.UserID,
		FeedID:    feedFollow.FeedID,
	})
}

func (s *Server) GetFeedFollows(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ctxUserID{}).(uuid.UUID)
	if !ok {
		log.Println("userID not found")
		encode(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	feedFollows, err := s.DB.GetFeedFollows(r.Context(), userID)
	if err != nil {
		log.Println(err)
		encode(w, http.StatusInternalServerError, "failed to fetch feeds")
		return
	}

	response := make([]FeedFollow, len(feedFollows))
	for _, feed := range feedFollows {
		response = append(response, FeedFollow{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			UserID:    feed.UserID,
			FeedID:    feed.FeedID,
		})
	}

	encode(w, http.StatusOK, response)
}

func (s *Server) DeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedIDStr := chi.URLParam(r, "id")

	feedID, err := uuid.Parse(feedIDStr)
	if err != nil {
		log.Println(err)
		encode(w, http.StatusBadRequest, "invalid feed id")
		return
	}

	userID, ok := r.Context().Value(ctxUserID{}).(uuid.UUID)
	if !ok {
		log.Println("userID not found")
		encode(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	err = s.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedID,
		UserID: userID,
	})
	if err != nil {
		log.Println(err)
		encode(w, http.StatusInternalServerError, "failed to fetch feeds")
		return
	}

	encode(w, http.StatusOK, "deleted successfully")
}

func (s *Server) GetPosts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ctxUserID{}).(uuid.UUID)
	if !ok {
		log.Println("userID not found")
		encode(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	posts, err := s.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: userID,
		Limit:  10,
	})
	if err != nil {
		log.Println(err)
		encode(w, http.StatusInternalServerError, "failed to fetch posts")
		return
	}

	response := make([]Post, len(posts))
	for _, post := range posts {
		response = append(response, Post{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Url:         post.Url,
			Title:       post.Title,
			Description: post.Description,
			PublishedAt: post.PublishedAt,
			FeedID:      post.FeedID,
		})
	}

	encode(w, http.StatusOK, response)
}
