package entity

import "testing"

func TestMakeChatID(t *testing.T) {
	tests := []struct {
		name   string
		userA  int64
		userB  int64
		expect string
	}{
		{"AscOrder", 100, 200, "100:200"},
		{"DescOrder", 200, 100, "100:200"},
		{"Symmetric", 42, 99, "42:99"},
		{"SymmetricReverse", 99, 42, "42:99"},
		{"SameUser", 5, 5, "5:5"},
		{"LargeIDs", 100000000001, 100000000002, "100000000001:100000000002"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeChatID(tt.userA, tt.userB)
			if got != tt.expect {
				t.Errorf("MakeChatID(%d, %d) = %q, want %q", tt.userA, tt.userB, got, tt.expect)
			}
		})
	}
}

func TestMakeChatID_Deterministic(t *testing.T) {
	id1 := MakeChatID(123, 456)
	id2 := MakeChatID(456, 123)
	if id1 != id2 {
		t.Errorf("MakeChatID not deterministic: %q != %q", id1, id2)
	}
}
