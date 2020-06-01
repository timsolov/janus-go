package janus

import (
	"testing"
)

func TestAdminTokens(t *testing.T) {
	api, err := NewAdminAPI("http://localhost:7088/admin", "janus-go")
	noError(t, err)

	resp, err := api.AddToken("test-token", []string{"janus.plugin.videoroom"})
	noError(t, err)
	if resp == nil {
		t.Errorf("resp is nil")
	}

	st := findToken(t, api, "test-token")
	if st == nil {
		t.Errorf("token not in list response")
	} else {
		if len(st.Plugins) != 1 {
			t.Errorf("expecting 1 plugin got %d", len(st.Plugins))
			return
		}
		if st.Plugins[0] != "janus.plugin.videoroom" {
			t.Errorf("expecting plugin %s != %s", "janus.plugin.videoroom", st.Plugins[0])
		}
	}

	resp, err = api.AllowToken("test-token", []string{"janus.plugin.echotest"})
	noError(t, err)

	st = findToken(t, api, "test-token")
	if st == nil {
		t.Errorf("token not in list response")
	} else {
		if len(st.Plugins) != 2 {
			t.Errorf("expecting 2 plugin got %d", len(st.Plugins))
		}
	}

	resp, err = api.DisallowToken("test-token", []string{"janus.plugin.videoroom"})
	noError(t, err)

	st = findToken(t, api, "test-token")
	if st == nil {
		t.Errorf("token not in list response")
	} else {
		if len(st.Plugins) != 1 {
			t.Errorf("expecting 1 plugin got %d", len(st.Plugins))
			return
		}
		if st.Plugins[0] != "janus.plugin.echotest" {
			t.Errorf("expecting plugin %s != %s", "janus.plugin.echotest", st.Plugins[0])
		}
	}

	resp, err = api.RemoveToken("test-token")
	noError(t, err)

	st = findToken(t, api, "test-token")
	if st != nil {
		t.Errorf("token in list response after removal")
	}

}

func findToken(t *testing.T, api AdminAPI, token string) *StoredToken {
	resp, err := api.ListTokens()
	noError(t, err)

	tokens, ok := resp.(*ListTokensResponse)
	if !ok {
		t.Errorf("wrong type: ListTokensResponse != %v", resp)
		return nil
	}

	for _, x := range tokens.Data["tokens"] {
		if x.Token == token {
			return x
		}
	}

	return nil
}

func TestAdminAPIImpl_ListSessions(t *testing.T) {
	client, err := Connect("ws://localhost:8188/")
	noError(t, err)
	defer client.Close()

	api, err := NewAdminAPI("http://localhost:7088/admin", "janus-go")
	noError(t, err)

	_, err = api.AddToken("test-token", []string{})
	noError(t, err)
	defer api.RemoveToken("test-token")
	client.Token = "test-token"

	session, err := client.Create()
	noError(t, err)
	defer session.Destroy()

	resp, err := api.ListSessions()
	noError(t, err)

	tResp, ok := resp.(*ListSessionsResponse)
	if !ok {
		t.Errorf("wrong type: ListSessionsResponse != %v", resp)
		return
	}
	if len(tResp.Sessions) != 1 {
		t.Errorf("expecting exactly 1 session, found %d", len(tResp.Sessions))
		return
	}

	if session.Id != tResp.Sessions[0] {
		t.Errorf("sessionID mismatch, expected %d got %d", session.Id, tResp.Sessions[0])
		return
	}
}

func TestAdminAPIImpl_MessagePlugin(t *testing.T) {
	client, err := Connect("ws://localhost:8188/")
	noError(t, err)
	defer client.Close()

	api, err := NewAdminAPI("http://localhost:7088/admin", "janus-go")
	noError(t, err)

	_, err = api.AddToken("test-token", []string{})
	noError(t, err)
	defer api.RemoveToken("test-token")
	client.Token = "test-token"

	resp, err := api.MessagePlugin("janus.plugin.videoroom", map[string]interface{}{"request": "list"})
	noError(t, err)

	tResp, ok := resp.(*MessagePluginResponse)
	if !ok {
		t.Errorf("wrong type: MessagePluginResponse != %v", resp)
		return
	}
	if len(tResp.Response) == 0 {
		t.Error("tResp.Response is empty")
		return
	}
}

func TestAdminAPIImpl_ListHandles(t *testing.T) {
	client, err := Connect("ws://localhost:8188/")
	noError(t, err)
	defer client.Close()

	api, err := NewAdminAPI("http://localhost:7088/admin", "janus-go")
	noError(t, err)

	_, err = api.AddToken("test-token", []string{})
	noError(t, err)
	defer api.RemoveToken("test-token")
	client.Token = "test-token"

	session, err := client.Create()
	noError(t, err)
	defer session.Destroy()

	handle, err := session.Attach("janus.plugin.videoroom")
	noError(t, err)
	defer handle.Detach()

	resp, err := api.ListHandles(session.Id)
	noError(t, err)
	if resp == nil {
		t.Error("resp is nil")
		return
	}
	tResp, ok := resp.(*ListHandlesResponse)
	if !ok {
		t.Errorf("wrong type: ListHandlesResponse != %v", resp)
		return
	}
	if len(tResp.Handles) != 1 {
		t.Errorf("expecting exactly 1 handle, found %d", len(tResp.Handles))
		return
	}

	if handle.Id != tResp.Handles[0] {
		t.Errorf("handleID mismatch, expected %d got %d", handle.Id, tResp.Handles[0])
		return
	}
}

func TestAdminAPIImpl_HandleInfo(t *testing.T) {
	client, err := Connect("ws://localhost:8188/")
	noError(t, err)
	defer client.Close()

	api, err := NewAdminAPI("http://localhost:7088/admin", "janus-go")
	noError(t, err)

	_, err = api.AddToken("test-token", []string{})
	noError(t, err)
	defer api.RemoveToken("test-token")
	client.Token = "test-token"

	session, err := client.Create()
	noError(t, err)
	defer session.Destroy()

	handle, err := session.Attach("janus.plugin.videoroom")
	noError(t, err)
	defer handle.Detach()

	resp, err := api.HandleInfo(session.Id+1, handle.Id+1)
	noError(t, err)
	if resp == nil {
		t.Error("resp is nil")
		return
	}
	tResp, ok := resp.(*HandleInfoResponse)
	if !ok {
		t.Errorf("wrong type: HandleInfoResponse != %v", resp)
		return
	}
	if tResp.Info == nil {
		t.Error("expected info to not be nil")
		return
	}
}

func noError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
