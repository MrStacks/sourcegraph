package shared

import (
	"github.com/sourcegraph/sourcegraph/internal/env"
)

type Config struct {
	env.BaseConfig

	Address string

	PubSub struct {
		ProjectID string
		TopicID   string
	}
}

func (c *Config) Load() {
	c.Address = c.Get("PING_SERVICE_ADDR", ":10086", "Address to serve Ping service on.")
	c.PubSub.ProjectID = c.Get("PING_SERVICE_PUBSUB_PROJECT_ID", "", "The project ID for the Pub/Sub.")
	c.PubSub.TopicID = c.Get("PING_SERVICE_PUBSUB_TOPIC_ID", "", "The topic ID for the Pub/Sub.")
}
