package endpoints

import (
	"FrenchConnections/internal/db"
	"FrenchConnections/internal/models"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/utils"
)

func Create(c *gin.Context) {
	var gameToCreate models.Game

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&gameToCreate); err != nil {
		slog.Debug("error decoding create game data", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid or unknown fields in game data"})
		return
	}

	if len(gameToCreate.GameCategories) != 4 {
		errMsg := fmt.Sprintf("Invalid number of game categories provided, must be 4 got %d", len(gameToCreate.GameCategories))
		slog.Debug(errMsg)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	if len(gameToCreate.CreatedBy) <= 2 {
		errMsg := fmt.Sprintf("Invalid creator name. Must be at least 2 characters was %d", len(gameToCreate.CreatedBy))
		slog.Debug(errMsg)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	// Check that games are valid:
	// - Words are present only one time in the game
	// - Words are not empty (len > 0)
	// - Category titles are valid (len > 0)
	seenWords := map[string]bool{}
	for i, cat := range gameToCreate.GameCategories {
		if len(cat.CategoryTitle) == 0 {
			errMsg := fmt.Sprintf("invalid game provided, category titles must be at least one character, the category at index %d has an empty title", i)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
		for _, word := range cat.Words {
			if len(word) == 0 {
				errMsg := fmt.Sprintf("invalid game provided, words must be at least one character, the category [%s] contains an empty word", cat.CategoryTitle)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errMsg})
				return
			}
			_, exists := seenWords[word]
			if exists {
				errMsg := fmt.Sprintf("invalid game provided, words must be unique but [%s] is present several times", word)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errMsg})
				return
			} else {
				seenWords[word] = true
			}
		}
	}

	gameToCreate.ID = 0
	gameToCreate.CreatedAt = time.Now()

	// Save the new game in DB
	result := db.GetDBClient().Create(&gameToCreate)
	if result.Error != nil || result.RowsAffected != 1 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "couldn't insert new game into database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"gameId": gameToCreate.ID})
}

func Retrieve(c *gin.Context) {
	gameIdInt, err := getGameId(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var game models.Game
	result := db.GetDBClient().Preload("GameCategories").First(&game, gameIdInt)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	// Collect all words of the game
	wordsList := make([]string, 0, 4*4)
	for _, gc := range game.GameCategories {
		wordsList = append(wordsList, gc.Words...)
	}

	// Shuffle the words in a reproducible manner to always get same result for the same game
	r := rand.New(rand.NewPCG(1234, 5678))
	r.Shuffle(len(wordsList), func(i, j int) {
		wordsList[i], wordsList[j] = wordsList[j], wordsList[i]
	})

	shuffledGame := models.ShuffledGame{CreatedBy: game.CreatedBy, CreationDate: game.CreatedAt, ShuffledWords: wordsList}
	c.JSON(http.StatusOK, shuffledGame)
}

func Guess(c *gin.Context) {
	gameId, err := getGameId(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var game models.Game
	result := db.GetDBClient().Preload("GameCategories").First(&game, gameId)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	var proposition models.Guess

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&proposition); err != nil {
		slog.Debug("error decoding guess data", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid or unknown fields in guess data"})
		return
	}

	if len(proposition.Proposition) != 4 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	guessResponse := &models.GuessResponse{Success: false, IsOneAway: false}
	for _, category := range game.GameCategories {
		i := 0
		for _, word := range category.Words {
			if !utils.Contains(proposition.Proposition, word) {
				break
			} else {
				i += 1
			}
		}
		if i == 3 {
			guessResponse.IsOneAway = true
			break
		} else if i == 4 {
			guessResponse.Success = true
			guessResponse.CategoryTitle = category.CategoryTitle
			break
		}
	}

	c.JSON(http.StatusOK, guessResponse)
}

func getGameId(c *gin.Context) (*int, error) {
	gameId := c.Param("gameId")
	if len(gameId) == 0 {
		return nil, fmt.Errorf("invalid game id")
	}
	gameIdInt, err := strconv.Atoi(gameId)
	return &gameIdInt, err
}
