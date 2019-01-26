package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type apiSuite struct {
	suite.Suite
	router *mux.Router
}

func (suite *apiSuite) SetupSuite() {
	suite.router = NewRouter()
}

func (suite *apiSuite) HappyPath() {
	const message = "foo"

	postBody := strings.NewReader(fmt.Sprintf(`{"message":"%s"}`, message))
	req, err := http.NewRequest("POST", "/messages", postBody)
	suite.Require().NoError(err, "error building POST request")

	recorder := httptest.NewRecorder()
	suite.router.ServeHTTP(recorder, req)

	// verify a 200 response with no error message
	suite.Require().Equal(http.StatusOK, recorder.Code)
	decoder := json.NewDecoder(recorder.Body)
	var body putResponse
	suite.Require().NoError(decoder.Decode(&body), "malformed HTTP response to putting message %s", "cat")
	suite.Require().Empty(body.Error)

	// and verify the SHA digest in the response body is correct
	digest := body.Digest
	digestBytes := sha256.Sum256([]byte(message))
	expectedDigest := fmt.Sprintf("%x", digestBytes)
	suite.Require().Equal(expectedDigest, digest)

	// now do a GET with the digest we got, and verify we get our original message back
	req, err = http.NewRequest("GET", fmt.Sprintf("/messages/%s", digest), nil)
	suite.Require().NoError(err, "error building POST request")

	suite.router.ServeHTTP(recorder, req)
	suite.Require().Equal(http.StatusOK, recorder.Code)

	decoder = json.NewDecoder(recorder.Body)
	var getBody getResponse
	suite.Require().NoError(decoder.Decode(&getBody), "malformed HTTP response to getting message for digest %s", digest)

	suite.Require().Equal(message, getBody.Message)
}

func (suite *apiSuite) PutMalformed() {
	postBody := strings.NewReader(`{"purr":"meow"}`)
	req, err := http.NewRequest("POST", "/messages", postBody)
	suite.Require().NoError(err, "error building POST request")

	recorder := httptest.NewRecorder()
	suite.router.ServeHTTP(recorder, req)

	// verify a 400 response when our request JSON doesn't contain the required "message" field
	suite.Require().Equal(http.StatusBadRequest, recorder.Code)
}

func (suite *apiSuite) PutEmpty() {
	postBody := strings.NewReader(`{"message":""}`)
	req, err := http.NewRequest("POST", "/messages", postBody)
	suite.Require().NoError(err, "error building POST request")

	recorder := httptest.NewRecorder()
	suite.router.ServeHTTP(recorder, req)

	// verify a 400 response since our JSON contains no message text
	suite.Require().Equal(http.StatusBadRequest, recorder.Code)

	decoder := json.NewDecoder(recorder.Body)
	var body putResponse
	suite.Require().NoError(decoder.Decode(&body), "malformed HTTP response to putting empty message")
	suite.Require().NotEmpty(body.Error)
}

func (suite *apiSuite) GetEmpty() {
	// expect a 400 if we give no digest at all

	req, err := http.NewRequest("GET", "/messages", nil)
	suite.Require().NoError(err, "error building GET request")

	recorder := httptest.NewRecorder()
	suite.router.ServeHTTP(recorder, req)

	suite.Require().Equal(http.StatusBadRequest, recorder.Code)

	req, err = http.NewRequest("GET", "/messages/", nil)
	suite.Require().NoError(err, "building GET request")

	suite.Require().Equal(http.StatusBadRequest, recorder.Code)
}

func (suite *apiSuite) GetNotFound() {
	// expect a 404 if we give an unrecognized digest

	req, err := http.NewRequest("GET", "/messages/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil)
	suite.Require().NoError(err, "error building GET request")

	recorder := httptest.NewRecorder()
	suite.router.ServeHTTP(recorder, req)

	suite.Require().Equal(http.StatusNotFound, recorder.Code)
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(apiSuite))
}
