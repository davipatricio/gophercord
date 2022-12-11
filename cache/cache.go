package cache

import "github.com/davipatricio/gophercord/structs"

type CachingManager struct {
	// A map of cached channels, mapped by their ID
	Channels map[string]*structs.Channel
	// A map of cached guilds, mapped by their ID
	Guilds map[string]*structs.Guild
	// A map of cached users, mapped by their ID
	Users map[string]*structs.User
}

type CacheConfig struct {
	// Maximum number of channels to cache
	ChannelLimit uint
	// Maximum number of guilds to cache
	GuildLimit uint
	// Maximum number of users to cache
	UserLimit uint
}

func NewCachingManager(config CacheConfig) *CachingManager {
	return &CachingManager{
		Channels: make(map[string]*structs.Channel, config.ChannelLimit),
		Guilds: make(map[string]*structs.Guild, config.GuildLimit),
		Users: make(map[string]*structs.User, config.UserLimit),
	}
}