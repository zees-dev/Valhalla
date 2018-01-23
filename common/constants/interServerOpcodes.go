package constants

const (
	LOGIN_REQUEST_WORLD_SUMMARY  = 0x00
	LOGIN_REQUEST_WORLD_STATUS   = 0x01
	LOGIN_REQUEST_MIGRATION_INFO = 0x02

	WORLD_REQUEST_ID     = 0x03
	WORLD_UPDATE_WORLD   = 0x04
	WORLD_UPDATE_CHANNEL = 0x05
	WORLD_DROPPED        = 0x06
	WORLD_NEW_CHANNEL    = 0x07
	WORLD_DELETE_CHANNEL = 0x08
	WORLD_SEND_CHANNELS  = 0x09

	CHANNEL_REGISTER         = 0x10
	CHANNEL_UPDATE           = 0x11
	CHANNEL_DROPPED          = 0x12
	CHANNEL_REQUEST_ID       = 0x13
	CHANNEL_USE_SAVED_IDs    = 0x14
	CHANNEL_GET_INTERNAL_IDS = 0x15
)
