package plugins

// VideoroomRequestFactory factory to make requests to VideoRoom plugin
type VideoroomRequestFactory struct {
	PluginRequestFactory
}

// NewVideoroomRequestFactory creates new instance of factory
func NewVideoroomRequestFactory(adminKey string) *VideoroomRequestFactory {
	return &VideoroomRequestFactory{
		PluginRequestFactory: *NewPluginRequestFactory("janus.plugin.videoroom", adminKey),
	}
}

//
func (f *VideoroomRequestFactory) ListRequest() *BasePluginRequest {
	request := f.make("list")
	return &request
}

func (f *VideoroomRequestFactory) CreateRequest(room *VideoroomRoom, permanent bool, allowed []string) *VideoroomCreateRequest {
	return &VideoroomCreateRequest{
		BasePluginRequest: f.make("create"),
		Room:              room,
		Permanent:         permanent,
		Allowed:           allowed,
	}
}

func (f *VideoroomRequestFactory) EditRequest(room *VideoroomRoomEdit, permanent bool, secret string) *VideoroomEditRequest {
	return &VideoroomEditRequest{
		BasePluginRequest: f.make("edit"),
		Room:              room,
		Permanent:         permanent,
		Secret:            secret,
	}
}

func (f *VideoroomRequestFactory) DestroyRequest(roomID int, permanent bool, secret string) *VideoroomDestroyRequest {
	return &VideoroomDestroyRequest{
		BasePluginRequest: f.make("destroy"),
		RoomID:            roomID,
		Permanent:         permanent,
		Secret:            secret,
	}
}
