package test

import (
	"testing"

	"ft_transcendence/backend/internal/socket"
)

func TestWSManager_LeaveRoom(t *testing.T) {
	m := socket.NewWSManager()
	client := &socket.Client{ID: "user1", Username: "test", Send: make(chan []byte, 1)}

	m.RegisterClient(client)
	m.JoinRoom(client, "room1")

	members := m.GetRoomMembers("room1")
	if len(members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(members))
	}

	m.LeaveRoom(client, "room1")

	members = m.GetRoomMembers("room1")
	if len(members) != 0 {
		t.Fatalf("expected 0 members after leave, got %d", len(members))
	}
}

func TestWSManager_GetRoomMembers(t *testing.T) {
	m := socket.NewWSManager()
	client1 := &socket.Client{ID: "user1", Username: "test1", Send: make(chan []byte, 1)}
	client2 := &socket.Client{ID: "user2", Username: "test2", Send: make(chan []byte, 1)}

	m.RegisterClient(client1)
	m.RegisterClient(client2)
	m.JoinRoom(client1, "room1")
	m.JoinRoom(client2, "room1")

	members := m.GetRoomMembers("room1")
	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
}

func TestWSManager_GetRoomMembersEmpty(t *testing.T) {
	m := socket.NewWSManager()

	members := m.GetRoomMembers("nonexistent")
	if len(members) != 0 {
		t.Fatalf("expected 0 members for nonexistent room, got %d", len(members))
	}
}

func TestWSManager_IsOnline(t *testing.T) {
	m := socket.NewWSManager()
	client := &socket.Client{ID: "user1", Username: "test", Send: make(chan []byte, 1)}

	if m.IsOnline("user1") {
		t.Fatal("expected offline before registration")
	}

	m.RegisterClient(client)

	if !m.IsOnline("user1") {
		t.Fatal("expected online after registration")
	}

	m.UnregisterClient(client)

	if m.IsOnline("user1") {
		t.Fatal("expected offline after unregistration")
	}
}

func TestWSManager_IsOnlineNonExistent(t *testing.T) {
	m := socket.NewWSManager()

	if m.IsOnline("nonexistent") {
		t.Fatal("expected offline for nonexistent user")
	}
}

func TestWSManager_BroadcastToRoom(t *testing.T) {
	m := socket.NewWSManager()
	client1 := &socket.Client{ID: "user1", Username: "test1", Send: make(chan []byte, 1)}
	client2 := &socket.Client{ID: "user2", Username: "test2", Send: make(chan []byte, 1)}

	m.RegisterClient(client1)
	m.RegisterClient(client2)
	m.JoinRoom(client1, "room1")
	m.JoinRoom(client2, "room1")

	m.BroadcastToRoom("room1", []byte("hello"), "user1")

	msg := <-client2.Send
	if string(msg) != "hello" {
		t.Fatalf("expected 'hello', got %q", string(msg))
	}

	select {
	case <-client1.Send:
		t.Fatal("sender should not receive broadcast")
	default:
	}
}

func TestWSManager_BroadcastToEmptyRoom(_ *testing.T) {
	m := socket.NewWSManager()

	m.BroadcastToRoom("nonexistent", []byte("hello"), "user1")
}

func TestWSManager_UnregisterCleansUpRooms(t *testing.T) {
	m := socket.NewWSManager()
	client := &socket.Client{ID: "user1", Username: "test", Send: make(chan []byte, 1)}

	m.RegisterClient(client)
	m.JoinRoom(client, "room1")
	m.JoinRoom(client, "room2")

	m.UnregisterClient(client)

	members1 := m.GetRoomMembers("room1")
	members2 := m.GetRoomMembers("room2")

	if len(members1) != 0 {
		t.Fatalf("expected 0 members in room1 after unregister, got %d", len(members1))
	}
	if len(members2) != 0 {
		t.Fatalf("expected 0 members in room2 after unregister, got %d", len(members2))
	}
}

func TestWSManager_MultipleRoomsPerClient(t *testing.T) {
	m := socket.NewWSManager()
	client := &socket.Client{ID: "user1", Username: "test", Send: make(chan []byte, 1)}

	m.RegisterClient(client)
	m.JoinRoom(client, "room1")
	m.JoinRoom(client, "room2")
	m.JoinRoom(client, "room3")

	if !m.IsOnline("user1") {
		t.Fatal("expected online")
	}

	members1 := m.GetRoomMembers("room1")
	members2 := m.GetRoomMembers("room2")
	members3 := m.GetRoomMembers("room3")

	if len(members1) != 1 || len(members2) != 1 || len(members3) != 1 {
		t.Fatal("expected client in all 3 rooms")
	}
}
