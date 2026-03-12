package cache

import "testing"

func TestNewCache(t *testing.T) {
	c := NewCache[string]()

	if c == nil {
		t.Errorf("expected non-nil cache, got nil")
	}

	obj, ok := c.Get("nonexistent")
	if ok {
		t.Errorf("expected Get to return false for nonexistent key, got true")
	}

	if obj != "" {
		t.Errorf("expected Get to return zero value for nonexistent key, got %v", obj)
	}
}

func TestSetCache(t *testing.T) {
	cases := []struct {
		name      string
		sets      []struct{ key, value string }
		getKey    string
		expected  string
		wantFound bool
	}{
		{
			name:      "store and retrieve",
			sets:      []struct{ key, value string }{{"key1", "value1"}},
			getKey:    "key1",
			expected:  "value1",
			wantFound: true,
		},
		{
			name: "overwrite existing key",
			sets: []struct{ key, value string }{
				{"key1", "value1"},
				{"key1", "value2"},
			},
			getKey:    "key1",
			expected:  "value2",
			wantFound: true,
		},
		{
			name: "set one key does not affect another",
			sets: []struct{ key, value string }{
				{"key1", "value1"},
			},
			getKey:    "key2",
			expected:  "",
			wantFound: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCache[string]()
			for _, s := range tc.sets {
				c.Set(s.key, s.value)
			}

			result, ok := c.Get(tc.getKey)

			if ok != tc.wantFound {
				t.Errorf("expected ok=%v for key %q, got %v", tc.wantFound, tc.getKey, ok)
			}
			if result != tc.expected {
				t.Errorf("expected Get(%q)=%q, got %q", tc.getKey, tc.expected, result)
			}
		})
	}
}

func TestGetCache(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		expected  string
		wantFound bool
	}{
		{
			name:      "existing key",
			input:     "key1",
			expected:  "value1",
			wantFound: true,
		},
		{
			name:      "nonexistent key",
			input:     "nonexistent",
			expected:  "",
			wantFound: false,
		},
	}

	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	c := Cache[string]{testData}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, ok := c.Get(tc.input)

			if ok != tc.wantFound {
				t.Errorf("expected ok=%v for key %q, got %v", tc.wantFound, tc.input, ok)
			}
			if result != tc.expected {
				t.Errorf("expected Get(%q)=%q, got %q", tc.input, tc.expected, result)
			}
		})
	}
}
