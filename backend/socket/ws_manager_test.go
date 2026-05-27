package socket

import "testing"

// newTestClient builds a Client usable for in-memory manager tests (no real
// network connection is needed for room bookkeeping).
func newTestClient(id string, buffer int) *Client {
	return &Client{ID: id, Username: id, Send: make(chan []byte, buffer)}
}

func TestWSManager_RoomMembershipAndBroadcast(t *testing.T) {
	m := NewWSManager()
	c1 := newTestClient("u1", 10)
	c2 := newTestClient("u2", 10)

	m.RegisterClient(c1)
	m.RegisterClient(c2)
	m.JoinRoom(c1, "room1")
	m.JoinRoom(c2, "room1")

	members := m.GetRoomMembers("room1")
	if len(members) != 2 {
		t.Fatalf("expected 2 room members, got %d", len(members))
	}

	m.BroadcastToRoom("room1", []byte("hello"), "u1")

	select {
	case <-c2.Send:
	default:
		t.Fatal("c2 should have received the broadcast")
	}
	select {
	case <-c1.Send:
		t.Fatal("sender u1 must not receive its own broadcast")
	default:
	}

	// broadcasting to an unknown room is a no-op
	m.BroadcastToRoom("ghost-room", []byte("x"), "")

	m.LeaveRoom(c1, "room1")
	if got := len(m.GetRoomMembers("room1")); got != 1 {
		t.Fatalf("expected 1 member after leave, got %d", got)
	}

	// leaving the last member removes the room entirely
	m.LeaveRoom(c2, "room1")
	if got := len(m.GetRoomMembers("room1")); got != 0 {
		t.Fatalf("expected empty room after last leave, got %d", got)
	}
}

func TestWSManager_Unregister(t *testing.T) {
	m := NewWSManager()
	c := newTestClient("solo", 10)
	m.RegisterClient(c)
	m.JoinRoom(c, "r")

	m.UnregisterClient(c)

	// Send channel is closed by UnregisterClient; receiving must not block.
	if _, open := <-c.Send; open {
		t.Fatal("expected Send channel to be closed after unregister")
	}
}

func TestSafeSend_BufferFullAndClosed(t *testing.T) {
	// full buffer: message is dropped without blocking
	full := make(chan []byte, 1)
	full <- []byte("first")
	safeSend(full, []byte("dropped"))

	// closed channel: the panic is recovered, not propagated
	closed := make(chan []byte, 1)
	close(closed)
	safeSend(closed, []byte("after-close"))
}
