package janus

import (
	"testing"
)

func TestAdminTokens(t *testing.T) {
	api, err := NewAdminAPI("http://localhost:7088/admin", "janus-go")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	resp, err := api.AddToken("test-token", []string{"janus.plugin.videoroom"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if resp == nil {
		t.Errorf("resp is nil")
	}

	st := findToken(t, api, "test-token")
	if st == nil {
		t.Errorf("token not in list response")
	} else {
		if len(st.Plugins) != 1 {
			t.Errorf("expecting 1 plugin got %d", len(st.Plugins))
		}
		if st.Plugins[0] != "janus.plugin.videoroom" {
			t.Errorf("expecting plugin %s != %s", "janus.plugin.videoroom", st.Plugins[0])
		}
	}

	resp, err = api.AllowToken("test-token", []string{"janus.plugin.echotest"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	st = findToken(t, api, "test-token")
	if st == nil {
		t.Errorf("token not in list response")
	} else {
		if len(st.Plugins) != 2 {
			t.Errorf("expecting 2 plugin got %d", len(st.Plugins))
		}
	}

	resp, err = api.DisallowToken("test-token", []string{"janus.plugin.videoroom"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	st = findToken(t, api, "test-token")
	if st == nil {
		t.Errorf("token not in list response")
	} else {
		if len(st.Plugins) != 1 {
			t.Errorf("expecting 1 plugin got %d", len(st.Plugins))
		}
		if st.Plugins[0] != "janus.plugin.echotest" {
			t.Errorf("expecting plugin %s != %s", "janus.plugin.echotest", st.Plugins[0])
		}
	}

	resp, err = api.RemoveToken("test-token")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	st = findToken(t, api, "test-token")
	if st != nil {
		t.Errorf("token in list response after removal")
	}

}

func findToken(t *testing.T, api *AdminAPI, token string) *StoredToken {
	resp, err := api.ListTokens()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tokens, ok := resp.(*ListTokensResponse)
	if !ok {
		t.Errorf("wrong type: ListTokensResponse != %v", resp)
	}
	for _, x := range tokens.Data["tokens"] {
		if x.Token == token {
			return x
		}
	}

	return nil
}
