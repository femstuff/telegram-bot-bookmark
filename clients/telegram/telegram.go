package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"telegram-bot/lib/e"
)

const (
	errMsg       = "cant do request"
	methodGetUpd = "getUpdates"
	methodSend   = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(methodGetUpd, q)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return res.Result, err
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return body, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chatID", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(methodSend, q)
	if err != nil {
		return e.Wrap("cant send message", err)
	}

	return nil
}
