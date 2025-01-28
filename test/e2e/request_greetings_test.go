package e2e

import (
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

func TestGetDateAndGreetingsSuite(t *testing.T) {
	suite.Run(t, new(RequestGreetingsTestSuite))
}

func (suite *RequestGreetingsTestSuite) TestGivenAnGreetingsRequestWhenClientDontPassApiKeyRateLimitMustBeBasedOnIpAddress() {
	for i := 0; i < suite.RateLimiterControl.IpMaxRequests; i++ {
		resp, err := suite.server.Client().Get(suite.server.URL + "/time/greetings")
		suite.NoError(err)
		suite.Equal(http.StatusOK, resp.StatusCode)
	}

	resp, err := suite.server.Client().Get(suite.server.URL + "/time/greetings")
	suite.NoError(err)
	suite.Equal(http.StatusTooManyRequests, resp.StatusCode)

	// A second before rate limit window ends
	seconds := time.Duration(suite.RateLimiterControl.IpWindowTime-1) * time.Second
	suite.redisMock.FastForward(seconds)

	// Still rate limited
	resp, err = suite.server.Client().Get(suite.server.URL + "/time/greetings")
	suite.NoError(err)
	suite.Equal(http.StatusTooManyRequests, resp.StatusCode)

	// A second after rate limit window ends
	seconds = time.Duration(suite.RateLimiterControl.IpWindowTime-1) * time.Second
	suite.redisMock.FastForward(seconds)

	// Not rate limited anymore
	resp, err = suite.server.Client().Get(suite.server.URL + "/time/greetings")
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *RequestGreetingsTestSuite) TestGivenAnGreetingsRequestWhenClientPassApiKeyRateLimitMustBeBasedOnApiKey() {
	req, err := http.NewRequest("GET", suite.server.URL+"/time/greetings", nil)
	suite.NoError(err)
	req.Header.Add("API_KEY", "123abc")

	for i := 0; i < suite.RateLimiterControl.ApiKeyMaxRequests; i++ {
		resp, err := suite.server.Client().Do(req)
		suite.NoError(err)
		suite.Equal(http.StatusOK, resp.StatusCode)
	}

	resp, err := suite.server.Client().Do(req)
	suite.NoError(err)
	suite.Equal(http.StatusTooManyRequests, resp.StatusCode)

	// A second before rate limit window ends
	seconds := time.Duration(suite.RateLimiterControl.IpWindowTime-1) * time.Second
	suite.redisMock.FastForward(seconds)

	resp, err = suite.server.Client().Do(req)
	suite.NoError(err)
	suite.Equal(http.StatusTooManyRequests, resp.StatusCode)

	// A second after rate limit window ends
	seconds = time.Duration(suite.RateLimiterControl.IpWindowTime-1) * time.Second
	suite.redisMock.FastForward(seconds)

	resp, err = suite.server.Client().Do(req)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
}
