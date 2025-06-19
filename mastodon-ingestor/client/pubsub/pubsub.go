package pubsub

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type Client struct {
	ClientID string
	QoS      byte

	connManager *autopaho.ConnectionManager
}

func New(ctx context.Context, host string, port uint16, clientID string, qos byte) (*Client, error) {
	var p Client
	var err error

	p.ClientID = clientID
	p.QoS = qos

	brokerURL, err := url.Parse(fmt.Sprintf("mqtt://%s:%d", host, port))
	if err != nil {
		return nil, fmt.Errorf("error parsing broker URL: %v", err)
	}

	cfg := autopaho.ClientConfig{
		ServerUrls:            []*url.URL{brokerURL},
		SessionExpiryInterval: 10 * 60, // 10 minutes for reconnection
		OnConnectError:        p.handleConnectError,
		ClientConfig: paho.ClientConfig{
			ClientID: p.ClientID,
		},
	}

	p.connManager, err = autopaho.NewConnection(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating pubsub connection: %v", err)
	}

	err = p.connManager.AwaitConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error awaiting pubsub connection: %v", err)
	}

	return &p, nil
}

func (p *Client) Close(ctx context.Context) error {
	err := p.connManager.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("error disconnecting: %v", err)
	}

	return nil
}

func (p *Client) Publish(ctx context.Context, topic string, payload []byte) error {
	_, err := p.connManager.Publish(ctx, &paho.Publish{
		Topic:   topic,
		Payload: payload,
		QoS:     p.QoS,
	})
	if err != nil {
		return fmt.Errorf("error publishing message: %v", err)
	}

	return nil
}

func (p *Client) handleConnectError(err error) {
	log.Printf("error with pubsub connection: %s", err)
}
