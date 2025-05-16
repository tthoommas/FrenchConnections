package endpoints_test

import (
	"FrenchConnections/internal/db"
	"FrenchConnections/internal/endpoints"
	"FrenchConnections/internal/models"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/v3/assert"
)

func TestListGamesHandler(t *testing.T) {
	r := gin.Default()
	r.GET("/game/list", endpoints.List)
	testGame := &models.Game{
		CreatedBy: "sami",
		GameCategories: []models.GameCategory{
			{
				CategoryTitle: "Category A",
				Words:         []string{"A", "B", "C", "D"},
			},
			{
				CategoryTitle: "Category B",
				Words:         []string{"E", "F", "G", "H"},
			},
			{
				CategoryTitle: "Category C",
				Words:         []string{"I", "J", "K", "L"},
			},
			{
				CategoryTitle: "Category D",
				Words:         []string{"M", "N", "O", "P"},
			},
		}}

	db.GetDBClient().Migrator().DropTable(&models.Game{}, &models.GameCategory{})
	db.GetDBClient().AutoMigrate(&models.Game{}, &models.GameCategory{})
	db.GetDBClient().Create(&testGame)
	testGame.ID = 0
	db.GetDBClient().Create(&testGame)
	testGame.ID = 0
	db.GetDBClient().Create(&testGame)

	req, _ := http.NewRequest("GET", "/game/list", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	expectedResponse := `\[{"id":1,"createdAt":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}[+-]\d{2}:\d{2}","updatedAt":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}[+-]\d{2}:\d{2}","createdBy":"sami"},{"id":2,"createdAt":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}[+-]\d{2}:\d{2}","updatedAt":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}[+-]\d{2}:\d{2}","createdBy":"sami"},{"id":3,"createdAt":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}[+-]\d{2}:\d{2}","updatedAt":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}[+-]\d{2}:\d{2}","createdBy":"sami"}]`
	if matched, err := regexp.MatchString(expectedResponse, string(responseData)); !matched || err != nil {
		assert.Error(t, fmt.Errorf("response body did not match expected"), expectedResponse, string(responseData))
	}
}

func TestCreateGameHandler(t *testing.T) {

	testCases := []struct {
		Name                 string
		RequestBody          string
		ExpectedStatus       int
		ExpectedResponseBody string
	}{
		{
			Name:                 "Invalid request body",
			RequestBody:          `{"random": "invalid body"}`,
			ExpectedStatus:       http.StatusBadRequest,
			ExpectedResponseBody: `{"error":"invalid or unknown fields in game data"}`,
		},
		{
			Name:                 "Wrong number of game categories",
			RequestBody:          `{"createdBy":"sami","gameCategories":[{"categoryTitle":"WASHING MACHINE CYCLES\/SETTINGS","words":["BULKY","COTTON","DELICATE","SPIN"]},{"categoryTitle":"WORDS SAID FREQUENTLY IN THE \u201cBILL AND TED\u201d MOVIES","words":["BOGUS","DUDE","EXCELLENT","TOTALLY"]},{"categoryTitle":"___BOX","words":["CHATTER","JUKE","SHADOW","SOAP"]}]}`,
			ExpectedStatus:       http.StatusBadRequest,
			ExpectedResponseBody: `{"error":"Invalid number of game categories provided, must be 4 got 3"}`,
		},
		{
			Name:                 "Wrong creator",
			RequestBody:          `{"createdBy":"","gameCategories":[{"categoryTitle":"FAUX","words":["ARTIFICIAL","FAKE","IMITATION","MOCK"]},{"categoryTitle":"WASHING MACHINE CYCLES\/SETTINGS","words":["BULKY","COTTON","DELICATE","SPIN"]},{"categoryTitle":"WORDS SAID FREQUENTLY IN THE \u201cBILL AND TED\u201d MOVIES","words":["BOGUS","DUDE","EXCELLENT","TOTALLY"]},{"categoryTitle":"___BOX","words":["CHATTER","JUKE","SHADOW","SOAP"]}]}`,
			ExpectedStatus:       http.StatusBadRequest,
			ExpectedResponseBody: `{"error":"Invalid creator name."}`,
		},
		{
			Name:                 "Empty category title",
			RequestBody:          `{"createdBy":"sami","gameCategories":[{"categoryTitle":"","words":["ARTIFICIAL","FAKE","IMITATION","MOCK"]},{"categoryTitle":"WASHING MACHINE CYCLES\/SETTINGS","words":["BULKY","COTTON","DELICATE","SPIN"]},{"categoryTitle":"WORDS SAID FREQUENTLY IN THE \u201cBILL AND TED\u201d MOVIES","words":["BOGUS","DUDE","EXCELLENT","TOTALLY"]},{"categoryTitle":"___BOX","words":["CHATTER","JUKE","SHADOW","SOAP"]}]}`,
			ExpectedStatus:       http.StatusBadRequest,
			ExpectedResponseBody: `{"error":"invalid game provided, category titles must be at least one character, the category at index 0 has an empty title"}`,
		},
		{
			Name:                 "Empty word in the game",
			RequestBody:          `{"createdBy":"sami","gameCategories":[{"categoryTitle":"FAUX","words":["","FAKE","IMITATION","MOCK"]},{"categoryTitle":"WASHING MACHINE CYCLES\/SETTINGS","words":["BULKY","FAKE","DELICATE","SPIN"]},{"categoryTitle":"WORDS SAID FREQUENTLY IN THE \u201cBILL AND TED\u201d MOVIES","words":["BOGUS","DUDE","EXCELLENT","TOTALLY"]},{"categoryTitle":"___BOX","words":["CHATTER","JUKE","SHADOW","SOAP"]}]}`,
			ExpectedStatus:       http.StatusBadRequest,
			ExpectedResponseBody: `{"error":"invalid game provided, words must be at least one character, the category \[FAUX] contains an empty word"}`,
		},
		{
			Name:                 "Duplicated word 'FAKE' in the game",
			RequestBody:          `{"createdBy":"sami","gameCategories":[{"categoryTitle":"FAUX","words":["ARTIFICIAL","FAKE","IMITATION","MOCK"]},{"categoryTitle":"WASHING MACHINE CYCLES\/SETTINGS","words":["BULKY","FAKE","DELICATE","SPIN"]},{"categoryTitle":"WORDS SAID FREQUENTLY IN THE \u201cBILL AND TED\u201d MOVIES","words":["BOGUS","DUDE","EXCELLENT","TOTALLY"]},{"categoryTitle":"___BOX","words":["CHATTER","JUKE","SHADOW","SOAP"]}]}`,
			ExpectedStatus:       http.StatusBadRequest,
			ExpectedResponseBody: `{"error":"invalid game provided, words must be unique but \[FAKE] is present several times"}`,
		},
		{
			Name:                 "Valid game creation",
			RequestBody:          `{"createdBy":"sami","gameCategories":[{"categoryTitle":"FAUX","words":["ARTIFICIAL","FAKE","IMITATION","MOCK"]},{"categoryTitle":"WASHING MACHINE CYCLES\/SETTINGS","words":["BULKY","COTTON","DELICATE","SPIN"]},{"categoryTitle":"WORDS SAID FREQUENTLY IN THE \u201cBILL AND TED\u201d MOVIES","words":["BOGUS","DUDE","EXCELLENT","TOTALLY"]},{"categoryTitle":"___BOX","words":["CHATTER","JUKE","SHADOW","SOAP"]}]}`,
			ExpectedStatus:       http.StatusCreated,
			ExpectedResponseBody: `{"gameId":\d}`,
		},
	}
	r := gin.Default()
	r.POST("/game", endpoints.Create)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/game", strings.NewReader(testCase.RequestBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			responseData, _ := io.ReadAll(w.Body)
			assert.Equal(t, testCase.ExpectedStatus, w.Code)
			if matched, err := regexp.MatchString(testCase.ExpectedResponseBody, string(responseData)); !matched || err != nil {
				assert.Error(t, fmt.Errorf("response body did not match expected"), testCase.ExpectedResponseBody, string(responseData))
			}
		})
	}
}

func TestRetrieveGameHandler(t *testing.T) {
	testCases := []struct {
		Name                 string
		GameId               string
		ExpectedStatus       int
		ExpectedResponseBody string
	}{
		{
			Name:                 "Fetch game success",
			GameId:               "1",
			ExpectedStatus:       http.StatusOK,
			ExpectedResponseBody: `{"createdBy":"sami","creationDate":"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}[+-]\d{2}:\d{2}","game":\["N","G","F","A","C","D","K","L","P","O","B","E","I","J","H","M"]}`,
		},
		{
			Name:                 "Fetch unexisting game",
			GameId:               "20",
			ExpectedStatus:       http.StatusNotFound,
			ExpectedResponseBody: ``,
		},
		{
			Name:                 "Fetch game wrong ID format",
			GameId:               "toto",
			ExpectedStatus:       http.StatusBadRequest,
			ExpectedResponseBody: ``,
		},
	}

	r := gin.Default()
	r.GET("/game/:gameId", endpoints.Retrieve)

	// Create a mock games for tests
	testGame := &models.Game{
		CreatedBy: "sami",
		GameCategories: []models.GameCategory{
			{
				CategoryTitle: "Category A",
				Words:         []string{"A", "B", "C", "D"},
			},
			{
				CategoryTitle: "Category B",
				Words:         []string{"E", "F", "G", "H"},
			},
			{
				CategoryTitle: "Category C",
				Words:         []string{"I", "J", "K", "L"},
			},
			{
				CategoryTitle: "Category D",
				Words:         []string{"M", "N", "O", "P"},
			},
		}}
	// Clear DB before running tests in case we run all tests at the same time
	db.GetDBClient().Migrator().DropTable(&models.Game{}, &models.GameCategory{})
	db.GetDBClient().AutoMigrate(&models.Game{}, &models.GameCategory{})
	db.GetDBClient().Create(&testGame)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/game/%s", testCase.GameId), nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			responseData, _ := io.ReadAll(w.Body)
			assert.Equal(t, testCase.ExpectedStatus, w.Code)
			if matched, err := regexp.MatchString(testCase.ExpectedResponseBody, string(responseData)); !matched || err != nil {
				assert.Error(t, fmt.Errorf("response body did not match expected"), testCase.ExpectedResponseBody, string(responseData))
			}
		})
	}
}

func TestGuessHandler(t *testing.T) {
	testCases := []struct {
		Name                 string
		GameId               string
		Guess                string
		ExpectedStatus       int
		ExpectedResponseBody string
	}{
		{
			Name:                 "No guess",
			GameId:               "1",
			Guess:                ``,
			ExpectedStatus:       http.StatusBadRequest,
			ExpectedResponseBody: ``,
		},
		{
			Name:                 "Invalid game ID",
			GameId:               "43",
			Guess:                ``,
			ExpectedStatus:       http.StatusNotFound,
			ExpectedResponseBody: ``,
		},
		{
			Name:                 "No game ID",
			GameId:               "",
			Guess:                ``,
			ExpectedStatus:       http.StatusBadRequest,
			ExpectedResponseBody: ``,
		},
		{
			Name:                 "Invalid guess",
			GameId:               "1",
			Guess:                `{"proposition": ["A", "E", "I", "M"]}`,
			ExpectedStatus:       http.StatusOK,
			ExpectedResponseBody: `{"success":false,"isOneAway":false,"categoryTitle":""}`,
		},
		{
			Name:                 "Invalid guess 2",
			GameId:               "1",
			Guess:                `{"proposition": ["A", "E", "I"]}`,
			ExpectedStatus:       http.StatusBadRequest,
			ExpectedResponseBody: ``,
		},
		{
			Name:                 "One away",
			GameId:               "1",
			Guess:                `{"proposition": ["A", "B", "C", "E"]}`,
			ExpectedStatus:       http.StatusOK,
			ExpectedResponseBody: `{"success":false,"isOneAway":true,"categoryTitle":""}`,
		},
		{
			Name:                 "Good guess",
			GameId:               "1",
			Guess:                `{"proposition": ["E", "F", "G", "H"]}`,
			ExpectedStatus:       http.StatusOK,
			ExpectedResponseBody: `{"success":true,"isOneAway":false,"categoryTitle":"Category B"}`,
		},
	}

	r := gin.Default()
	r.POST("/game/:gameId/guess", endpoints.Guess)

	// Create a mock games for tests
	testGame := &models.Game{
		CreatedBy: "sami",
		GameCategories: []models.GameCategory{
			{
				CategoryTitle: "Category A",
				Words:         []string{"A", "B", "C", "D"},
			},
			{
				CategoryTitle: "Category B",
				Words:         []string{"E", "F", "G", "H"},
			},
			{
				CategoryTitle: "Category C",
				Words:         []string{"I", "J", "K", "L"},
			},
			{
				CategoryTitle: "Category D",
				Words:         []string{"M", "N", "O", "P"},
			},
		}}
	// Clear DB before running tests in case we run all tests at the same time
	db.GetDBClient().Migrator().DropTable(&models.Game{}, &models.GameCategory{})
	db.GetDBClient().AutoMigrate(&models.Game{}, &models.GameCategory{})
	db.GetDBClient().Create(&testGame)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", fmt.Sprintf("/game/%s/guess", testCase.GameId), strings.NewReader(testCase.Guess))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			responseData, _ := io.ReadAll(w.Body)
			assert.Equal(t, testCase.ExpectedStatus, w.Code)
			if matched, err := regexp.MatchString(testCase.ExpectedResponseBody, string(responseData)); !matched || err != nil {
				assert.Error(t, fmt.Errorf("response body did not match expected"), testCase.ExpectedResponseBody, string(responseData))
			}
		})
	}
}
