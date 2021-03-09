package network

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ProxyDialerTestSuite struct {
	suite.Suite

	u *url.URL
}

func (suite *ProxyDialerTestSuite) SetupSuite() {
	u, _ := url.Parse("socks5://hello:world@10.0.0.10:3128")
	suite.u = u
}

func (suite *ProxyDialerTestSuite) TestSetupDefaults() {
	d := newProxyDialer(&DialerMock{}, suite.u).(*circuitBreakerDialer)

	suite.EqualValues(ProxyDialerOpenThreshold, d.openThreshold)
	suite.EqualValues(ProxyDialerHalfOpenTimeout, d.halfOpenTimeout)
	suite.EqualValues(ProxyDialerResetFailuresTimeout, d.resetFailuresTimeout)
}

func (suite *ProxyDialerTestSuite) TestSetupValuesAllOk() {
	query := url.Values{}

	query.Set("open_threshold", "30")
	query.Set("reset_failures_timeout", "1s")
	query.Set("half_open_timeout", "2s")

	suite.u.RawQuery = query.Encode()
	d := newProxyDialer(&DialerMock{}, suite.u).(*circuitBreakerDialer)

	suite.EqualValues(30, d.openThreshold)
	suite.EqualValues(2*time.Second, d.halfOpenTimeout)
	suite.EqualValues(time.Second, d.resetFailuresTimeout)
}

func (suite *ProxyDialerTestSuite) TestOpenThreshold() {
	query := url.Values{}
	params := []string{"-30", "aaa", "1.0", "-1.0"}

	for _, v := range params {
		suite.T().Run(fmt.Sprintf("param=%s", v), func(t *testing.T) {
			query.Set("open_threshold", v)

			suite.u.RawQuery = query.Encode()
			d := newProxyDialer(&DialerMock{}, suite.u).(*circuitBreakerDialer)

			suite.EqualValues(ProxyDialerOpenThreshold, d.openThreshold)
		})
	}
}

func (suite *ProxyDialerTestSuite) TestHalfOpenTimeout() {
	query := url.Values{}
	params := []string{"-30", "30", "aaa", "-3.0", "3.0"}

	for _, v := range params {
		suite.T().Run(fmt.Sprintf("param=%s", v), func(t *testing.T) {
			query.Set("half_open_timeout", v)

			suite.u.RawQuery = query.Encode()
			d := newProxyDialer(&DialerMock{}, suite.u).(*circuitBreakerDialer)

			suite.EqualValues(ProxyDialerHalfOpenTimeout, d.halfOpenTimeout)
		})
	}
}

func (suite *ProxyDialerTestSuite) TestResetFailuresTimeout() {
	query := url.Values{}
	params := []string{"-30", "30", "aaa", "-3.0", "3.0"}

	for _, v := range params {
		suite.T().Run(fmt.Sprintf("param=%s", v), func(t *testing.T) {
			query.Set("reset_failures_timeout", v)

			suite.u.RawQuery = query.Encode()
			d := newProxyDialer(&DialerMock{}, suite.u).(*circuitBreakerDialer)

			suite.EqualValues(ProxyDialerHalfOpenTimeout, d.halfOpenTimeout)
		})
	}
}

func TestProxyDialer(t *testing.T) {
	suite.Run(t, &ProxyDialerTestSuite{})
}
