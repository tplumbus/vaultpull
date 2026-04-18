package vault

import "testing"

func TestKVVersionFromString(t *testing.T) {
	tests := []struct {
		input    string
		want     KVVersion
		wantErr  bool
	}{
		{"1", KVVersion1, false},
		{"2", KVVersion2, false},
		{"", KVVersion2, false},
		{"v1", 0, true},
		{"3", 0, true},
		{"0", 0, true},
	}

	for _, tc := range tests {
		t.Run("input="+tc.input, func(t *testing.T) {
			got, err := KVVersionFromString(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("KVVersionFromString(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestKVVersion_String(t *testing.T) {
	if KVVersion1.String() != "1" {
		t.Errorf("expected '1', got %q", KVVersion1.String())
	}
	if KVVersion2.String() != "2" {
		t.Errorf("expected '2', got %q", KVVersion2.String())
	}
}
