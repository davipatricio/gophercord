package client

import (
	"time"

	"github.com/davipatricio/gophercord/cache"
	"nhooyr.io/websocket"
)

type Client struct {
	ApplicationID   string                `json:"application_id"`
	CachingManager  *cache.CachingManager `json:"caching_manager"`
	Token           string
	Ready           bool
	sequence        int
	ws              *websocket.Conn
	heartbeatTicker *time.Ticker
}

// NewClient creates a new bot client
//
//	bot := client.NewClient("token")
func NewClient(token string) (client *Client) {
	client = &Client{}
	client.Token = token
	client.SetCachingManager(cache.NewCachingManager(cache.CacheConfig{ChannelLimit: 100, GuildLimit: 100, UserLimit: 100}))

	return
}

// SetCachingManager sets the caching manager for the client
//
//	client.SetCachingManager(cache.NewCachingManager(cache.CacheConfig{100, 100, 100}))
func (client *Client) SetCachingManager(manager *cache.CachingManager) {
	client.CachingManager = manager
}
