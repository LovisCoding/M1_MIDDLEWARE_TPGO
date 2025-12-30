package nats

import (
	"github.com/nats-io/nats.go"
)

type Client struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewClient(url string) (*Client, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &Client{
		nc: nc,
		js: js,
	}, nil
}

func (c *Client) InitStream(streamName string, subjects []string) error {
	_, err := c.js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: subjects,
	})
	return err
}

func (c *Client) Publish(subject string, data []byte) error {
	pubAckFuture, err := c.js.PublishAsync(subject, data)
	if err != nil {
		return err
	}

	select {
	case <-pubAckFuture.Ok():
		return nil
	case err := <-pubAckFuture.Err():
		return err
	}
}

func (c *Client) Close() {
	if c.nc != nil {
		c.nc.Close()
	}
}
