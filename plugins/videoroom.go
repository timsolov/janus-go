package plugins

import (
	"github.com/timsolov/janus-go"
)

// VideoroomResponse base response
type VideoroomResponse struct {
	Videoroom string `json:"videoroom"`
}

// VideoroomErrorResponse error response
type VideoroomErrorResponse struct {
	VideoroomResponse
	PluginError
}

func (err *VideoroomErrorResponse) Error() string {
	return err.PluginError.Error()
}

// VideoroomListResponse list of rooms
type VideoroomListResponse struct {
	VideoroomResponse
	Rooms []*VideoroomRoomListEntry `json:"list"`
}

// VideoroomCreateRequest create room
type VideoroomCreateRequest struct {
	BasePluginRequest
	Room      *VideoroomRoom
	Permanent bool
	Allowed   []string
}

// Payload ...
func (r *VideoroomCreateRequest) Payload() map[string]interface{} {
	payload := r.BasePluginRequest.Payload()
	payload["permanent"] = r.Permanent
	if len(r.Allowed) > 0 {
		payload["allowed"] = r.Allowed
	}
	mergeMap(payload, r.Room.AsMap())
	return payload
}

// VideoroomCreateResponse success response on create room
type VideoroomCreateResponse struct {
	VideoroomResponse
	RoomID    int  `json:"room"`
	Permanent bool `json:"permanent"`
}

// VideoroomEditRequest edit room
type VideoroomEditRequest struct {
	BasePluginRequest
	Room      *VideoroomRoomEdit
	Secret    string
	Permanent bool
}

// Payload ...
func (r *VideoroomEditRequest) Payload() map[string]interface{} {
	payload := r.BasePluginRequest.Payload()
	payload["permanent"] = r.Permanent
	if r.Secret != "" {
		payload["secret"] = r.Secret
	}
	mergeMap(payload, r.Room.AsMap())
	return payload
}

// VideoroomEditResponse success response on edit room
type VideoroomEditResponse struct {
	VideoroomResponse
	RoomID int `json:"room"`
}

// VideoroomDestroyRequest destroy room
type VideoroomDestroyRequest struct {
	BasePluginRequest
	RoomID    int
	Secret    string
	Permanent bool
}

// Payload ...
func (r *VideoroomDestroyRequest) Payload() map[string]interface{} {
	payload := r.BasePluginRequest.Payload()
	payload["room"] = r.RoomID
	payload["permanent"] = r.Permanent
	if r.Secret != "" {
		payload["secret"] = r.Secret
	}
	return payload
}

// VideoroomDestroyResponse success response on destroy room
type VideoroomDestroyResponse struct {
	VideoroomResponse
	RoomID int `json:"room"`
}

// VideoroomRoom describes room settings
type VideoroomRoom struct {
	Room               int    `json:"room,omitempty"`
	Description        string `json:"description,omitempty"`
	IsPrivate          bool   `json:"is_private"`
	Secret             string `json:"secret,omitempty"`
	Pin                string `json:"pin,omitempty"`
	RequirePvtID       bool   `json:"require_pvtid"`
	RequireE2ee        bool   `json:"require_e2ee"`
	Publishers         int    `json:"publishers"`
	Bitrate            int    `json:"bitrate"`
	FirFreq            int    `json:"fir_freq"`
	AudioCodec         string `json:"audiocodec,omitempty"`
	VideoCodec         string `json:"videocodec,omitempty"`
	Vp9Profile         string `json:"vp9_profile,omitempty"`
	H264Profile        string `json:"h264_profile,omitempty"`
	OpusFec            bool   `json:"opus_fec"`
	VideoSvc           bool   `json:"video_svc"`
	AudioLevelExt      bool   `json:"audiolevel_ext"`
	AudioLevelEvent    bool   `json:"audiolevel_event"`
	AudioActivePackets int    `json:"audio_active_packets,omitempty"`
	AudioLevelAverage  int    `json:"audio_level_average,omitempty"`
	VideoOrientExt     bool   `json:"videoorient_ext"`
	PlayoutDelayExt    bool   `json:"playoutdelay_ext"`
	TransportWideCCExt bool   `json:"transport_wide_cc_ext"`
	Record             bool   `json:"record"`
	RecDir             string `json:"rec_dir,omitempty"`
	LockRecord         bool   `json:"lock_record"`
	NotifyJoining      bool   `json:"notify_joining"`
}

// AsMap convert struct to map
func (r *VideoroomRoom) AsMap() map[string]interface{} {
	m, _ := janus.StructToMap(r)
	return m
}

// VideoroomRoomListEntry each record from room list
type VideoroomRoomListEntry struct {
	VideoroomRoom
	PinRequired     bool `json:"pin_required"`
	MaxPublishers   int  `json:"max_publishers"`
	BitrateCap      bool `json:"bitrate_cap"`
	NumParticipants int  `json:"num_participants"`
}

// VideoroomRoomEdit edit room
type VideoroomRoomEdit struct {
	Room         int    `json:"room"`
	Description  string `json:"new_description,omitempty"`
	IsPrivate    bool   `json:"new_is_private"`
	Secret       string `json:"new_secret,omitempty"`
	Pin          string `json:"new_pin,omitempty"`
	RequirePvtID bool   `json:"new_require_pvtid"`
	Publishers   int    `json:"new_publishers"`
	Bitrate      int    `json:"new_bitrate"`
	FirFreq      int    `json:"new_fir_freq"`
	LockRecord   bool   `json:"new_lock_record"`
}

// AsMap convert struct to map
func (r *VideoroomRoomEdit) AsMap() map[string]interface{} {
	m, _ := janus.StructToMap(r)
	return m
}
