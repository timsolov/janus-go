package plugins

import "testing"

func TestTextroomRoom_AsMap(t *testing.T) {
	r := TextroomRoom{}
	m := r.AsMap()
	for _, k := range []string{"description", "secret", "pin", "post"} {
		if _, ok := m[k]; ok {
			t.Errorf("empty field [%s] should have been omitted", k)
		}
	}
}

func TestTextroomRoomForEdit_AsMap(t *testing.T) {
	r := TextroomRoomForEdit{}
	m := r.AsMap()
	for _, k := range []string{"new_description", "new_secret", "new_pin", "new_post"} {
		if _, ok := m[k]; ok {
			t.Errorf("empty field [%s] should have been omitted", k)
		}
	}
}
