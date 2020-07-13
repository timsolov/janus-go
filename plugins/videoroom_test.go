package plugins

import "testing"

func TestVideoroomRoom_AsMap(t *testing.T) {
	r := VideoroomRoom{}
	m := r.AsMap()
	for _, k := range []string{"description", "secret", "pin", "audiocodec", "videocodec", "vp9_profile", "h264_profile", "rec_dir"} {
		if _, ok := m[k]; ok {
			t.Errorf("empty field [%s] should have been omitted", k)
		}
	}
}

func TestVideoroomRoomForEdit_AsMap(t *testing.T) {
	r := VideoroomRoomForEdit{}
	m := r.AsMap()
	for _, k := range []string{"new_description", "new_secret", "new_pin"} {
		if _, ok := m[k]; ok {
			t.Errorf("empty field [%s] should have been omitted", k)
		}
	}
}
