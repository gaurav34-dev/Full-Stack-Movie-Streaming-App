package controllers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gaurav34-dev/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/gaurav34-dev/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const testDBName = "magicstream_test"

var testClient *mongo.Client

// TestMain connects to MongoDB once, runs all tests, then drops the test database.
func TestMain(m *testing.M) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	os.Setenv("DATABASE_NAME", testDBName)

	var err error
	testClient, err = mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic("mongodb connect: " + err.Error())
	}

	code := m.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = testClient.Database(testDBName).Drop(ctx)
	_ = testClient.Disconnect(ctx)

	os.Exit(code)
}

// seedMovies inserts movies into the test collection and registers cleanup.
func seedMovies(t *testing.T, movies []models.Movie) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := testClient.Database(testDBName).Collection("movies")
	docs := make([]interface{}, len(movies))
	for i, mv := range movies {
		docs[i] = mv
	}
	if _, err := col.InsertMany(ctx, docs); err != nil {
		t.Fatalf("seed movies: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		col.DeleteMany(ctx, bson.D{})
	})
}

func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/movies", controllers.GetMovies(testClient))
	return r
}

// TestGetMovies_ReturnsSeededMovies checks that all inserted movies are returned.
func TestGetMovies_ReturnsSeededMovies(t *testing.T) {
	seedMovies(t, []models.Movie{
		{
			ImdbID:     "tt0000001",
			Title:      "Test Movie One",
			PosterPath: "https://example.com/poster1.jpg",
			YouTubeID:  "abc123",
			Genre:      []models.Genre{{GenreID: 1, GenreName: "Action"}},
			Ranking:    models.Ranking{RankingValue: 1, RankingName: "Excellent"},
		},
		{
			ImdbID:     "tt0000002",
			Title:      "Test Movie Two",
			PosterPath: "https://example.com/poster2.jpg",
			YouTubeID:  "def456",
			Genre:      []models.Genre{{GenreID: 2, GenreName: "Drama"}},
			Ranking:    models.Ranking{RankingValue: 2, RankingName: "Good"},
		},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	newRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d — body: %s", w.Code, w.Body.String())
	}

	var got []models.Movie
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 movies, got %d", len(got))
	}

	byID := map[string]models.Movie{}
	for _, mv := range got {
		byID[mv.ImdbID] = mv
	}
	if _, ok := byID["tt0000001"]; !ok {
		t.Error("tt0000001 missing from response")
	}
	if _, ok := byID["tt0000002"]; !ok {
		t.Error("tt0000002 missing from response")
	}
}

// TestGetMovies_EmptyCollection checks that an empty collection returns 200 with an empty body.
func TestGetMovies_EmptyCollection_ReturnsEmpty(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	testClient.Database(testDBName).Collection("movies").Drop(ctx)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	newRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	// MongoDB driver may encode a nil slice as null — both null and [] are acceptable.
	body := w.Body.String()
	if body != "null" && body != "[]" {
		var got []models.Movie
		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected 0 movies, got %d", len(got))
		}
	}
}

// TestGetMovies_ResponseShape verifies the JSON keys match the model json tags.
func TestGetMovies_ResponseShape(t *testing.T) {
	seedMovies(t, []models.Movie{
		{
			ImdbID:     "tt9999999",
			Title:      "Shape Test Movie",
			PosterPath: "https://example.com/shape.jpg",
			YouTubeID:  "zzz999",
			Genre:      []models.Genre{{GenreID: 5, GenreName: "Comedy"}},
			Ranking:    models.Ranking{RankingValue: 3, RankingName: "Average"},
		},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	newRouter().ServeHTTP(w, req)

	var raw []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(raw) == 0 {
		t.Fatal("expected at least one movie in response")
	}

	movie := raw[0]
	for _, field := range []string{"imdb_id", "title", "poster_path", "youtube_id", "genre", "ranking"} {
		if _, ok := movie[field]; !ok {
			t.Errorf("response body missing field %q", field)
		}
	}
}
