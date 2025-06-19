package pubsub

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type PubSub struct {
	ClientID string
	QoS      byte

	subscriptions map[string]func(ctx context.Context, payload []byte) error
	connManager   *autopaho.ConnectionManager
}

func New(ctx context.Context, host string, port uint16, clientID string, qos byte) (*PubSub, error) {
	var p PubSub
	var err error

	p.ClientID = clientID
	p.QoS = qos
	p.subscriptions = make(map[string]func(ctx context.Context, payload []byte) error)

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
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				p.handleMessage,
			},
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

func (p *PubSub) Close(ctx context.Context) error {
	err := p.connManager.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("error disconnecting: %v", err)
	}

	return nil
}

func (p *PubSub) Subscribe(ctx context.Context, topic string, handler func(ctx context.Context, payload []byte) error) error {
	p.subscriptions[topic] = handler

	_, err := p.connManager.Subscribe(ctx, &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{Topic: topic, QoS: p.QoS},
		},
	})
	if err != nil {
		return fmt.Errorf("error subscribing to topic: %v", err)
	}

	return nil
}

func (p *PubSub) handleMessage(message paho.PublishReceived) (bool, error) {
	for topic, handler := range p.subscriptions {
		if message.Packet.Topic == topic {
			err := handler(context.Background(), message.Packet.Payload)
			if err != nil {
				return true, fmt.Errorf("error handling message: %v", err)
			}
		}
	}

	return true, nil
}

func (p *PubSub) handleConnectError(err error) {
	log.Printf("error with pubsub connection: %s", err)
}
