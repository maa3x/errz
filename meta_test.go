package errz

import (
	"strings"
	"testing"
)

func TestMetadata(t *testing.T) {
	t.Run("nil metadata", func(t *testing.T) {
		var m Metadata
		if m.String() != "" {
			t.Error("nil metadata should return empty string")
		}
		if m.Has("key") {
			t.Error("nil metadata should not have any keys")
		}
		if values := m.Get("key"); values != nil {
			t.Error("nil metadata should return nil values")
		}
		if m = m.Add("key", "value"); m != nil {
			t.Error("nil metadata should remain nil after Add")
		}
	})

	t.Run("add and get", func(t *testing.T) {
		m := make(Metadata, 0)
		m = m.Add("key1", "value1")
		m = m.Add("key2", 42)
		m = m.Add("key1", "value3")

		if !m.Has("key1") {
			t.Error("metadata should have key1")
		}
		if !m.Has("key2") {
			t.Error("metadata should have key2")
		}
		if m.Has("key3") {
			t.Error("metadata should not have key3")
		}

		values := m.Get("key1")
		if len(values) != 2 {
			t.Errorf("expected 2 values for key1, got %d", len(values))
		}
		if values[0] != "value1" || values[1] != "value3" {
			t.Error("unexpected values for key1")
		}

		values = m.Get("key2")
		if len(values) != 1 || values[0] != 42 {
			t.Error("unexpected value for key2")
		}
	})

	t.Run("string representation", func(t *testing.T) {
		m := make(Metadata, 0)
		m = m.Add("key1", "value1")
		m = m.Add("key2", 42)

		str := m.String()
		if !strings.Contains(str, "key1") || !strings.Contains(str, "value1") {
			t.Error("string representation should contain key1 and value1")
		}
		if !strings.Contains(str, "key2") || !strings.Contains(str, "42") {
			t.Error("string representation should contain key2 and 42")
		}
	})

	t.Run("detail string", func(t *testing.T) {
		d := detail{Key: "key", Value: "value"}
		if d.String() != "key\tvalue" {
			t.Errorf("unexpected detail string: %s", d.String())
		}
	})
}
