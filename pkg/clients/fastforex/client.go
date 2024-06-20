package fastforex

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/romutchio/crypto-calculator/internal/config"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

type Client struct {
	config *config.FastForex
	ln     *fasthttputil.InmemoryListener
}

var (
	ErrEmptyConfig           = errors.New("field of config can not be empty")
	ErrBadResponseStatusCode = errors.New("bad response status code")
)

// New создаёт новый клиент fastforex, проверяет что конфиг поля не пустые.
func New(conf *config.FastForex, ln *fasthttputil.InmemoryListener) (*Client, error) {
	switch {
	case conf.URL == "":
		return nil, fmt.Errorf("%w: url is empty", ErrEmptyConfig)
	case conf.Token == "":
		return nil, fmt.Errorf("%w: token is empty", ErrEmptyConfig)
	}

	c := &Client{
		config: conf,
		ln:     ln,
	}
	return c, nil
}

func (c *Client) request(url, method string, body []byte) ([]byte, error) {
	client := &fasthttp.Client{}
	if c.config.IsLocal && c.ln != nil {
		client = &fasthttp.Client{
			Dial: func(addr string) (net.Conn, error) {
				conn, err := c.ln.Dial()
				if err != nil {
					return nil, fmt.Errorf("client dial failed %w", err)
				}
				return conn, nil
			},
		}
	}

	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	req.Header.SetMethod(method)
	req.Header.SetContentType("application/json")
	req.SetRequestURI(c.config.URL + url)
	if body != nil {
		req.SetBody(body)
	}

	if err := client.DoTimeout(req, res, c.config.Timeout); err != nil {
		return nil, fmt.Errorf("get from fastforex: %w", err)
	}
	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%w: url '%s' answer with code %d: '%s'", ErrBadResponseStatusCode, url, res.StatusCode(), body)
	}

	return res.Body(), nil
}
