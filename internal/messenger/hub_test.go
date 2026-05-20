package messenger

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func testHub(t *testing.T) (*Hub, context.CancelFunc) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	h := NewHub()
	go h.Run(ctx)
	return h, cancel
}

func testClient(id int64) *Client {
	return &Client{id: id, send: make(chan []byte, 256), done: make(chan struct{})}
}

func registerClient(t *testing.T, h *Hub, c *Client) {
	t.Helper()
	c.hub = h
	select {
	case h.register <- c:
	case <-time.After(time.Second):
		t.Fatal("timeout registering client")
	}
}

func TestHub_Run_DirectMessage(t *testing.T) {
	h, cancel := testHub(t)
	defer cancel()

	c1 := testClient(1)
	c2 := testClient(2)
	registerClient(t, h, c1)
	registerClient(t, h, c2)

	msg := DirectMessage{From: 1, To: 2, Message: "hello"}
	select {
	case h.direct <- msg:
	case <-time.After(time.Second):
		t.Fatal("timeout sending direct message")
	}

	select {
	case gotBytes := <-c2.send:
		var gotMsg DirectMessage
		if err := json.Unmarshal(gotBytes, &gotMsg); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if gotMsg.Message != "hello" || gotMsg.From != 1 {
			t.Errorf("got %+v, want %+v", gotMsg, msg)
		}
	case <-time.After(time.Second):
		t.Fatal("message not delivered to client")
	}

	select {
	case h.unregister <- c1:
	case <-time.After(time.Second):
		t.Fatal("timeout unregistering c1")
	}
	select {
	case h.unregister <- c2:
	case <-time.After(time.Second):
		t.Fatal("timeout unregistering c2")
	}

	select {
	case _, ok := <-c1.send:
		if ok {
			t.Error("expected closed channel c1")
		}
	case <-time.After(time.Second):
		t.Fatal("channel c1 not closed")
	}
}

func TestHub_Register_ReplacesOldConnection(t *testing.T) {
	h, cancel := testHub(t)
	defer cancel()

	c1Old := testClient(1)
	registerClient(t, h, c1Old)

	c1New := testClient(1)
	registerClient(t, h, c1New)

	select {
	case _, ok := <-c1Old.send:
		if ok {
			t.Error("expected old connection channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("old connection channel not closed")
	}
}

func TestHub_DirectMessage_NonExistentClient(t *testing.T) {
	h, cancel := testHub(t)
	defer cancel()

	msg := DirectMessage{From: 1, To: 999, Message: "hello"}
	select {
	case h.direct <- msg:
	case <-time.After(time.Second):
		t.Fatal("timeout sending to non-existent client")
	}
}

func TestHub_Unregister_NotFound(t *testing.T) {
	h, cancel := testHub(t)
	defer cancel()

	c := testClient(1)
	c.hub = h
	select {
	case h.unregister <- c:
	case <-time.After(time.Second):
		t.Fatal("timeout unregistering unknown client")
	}
}

func TestHub_Shutdown(t *testing.T) {
	h, cancel := testHub(t)

	c := testClient(1)
	registerClient(t, h, c)

	cancel()

	select {
	case _, ok := <-c.send:
		if ok {
			t.Error("expected client channel to be closed on shutdown")
		}
	case <-time.After(time.Second):
		t.Fatal("client channel not closed after shutdown")
	}
}

func TestHub_Shutdown_ViaStop(t *testing.T) {
	h, cancel := testHub(t)
	defer cancel()

	c := testClient(1)
	registerClient(t, h, c)

	h.Stop()

	select {
	case _, ok := <-c.send:
		if ok {
			t.Error("expected client channel to be closed on stop")
		}
	case <-time.After(time.Second):
		t.Fatal("client channel not closed after stop")
	}
}
