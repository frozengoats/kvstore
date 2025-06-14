package kvstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFromMappingGood(t *testing.T) {
	m := map[string]any{
		"hello": "world",
		"super": map[string]any{
			"awesome": 10,
			"happy":   nil,
		},
	}
	_, err := FromMapping(m)
	assert.NoError(t, err)
}

func TestCreateFromMappingBadMapping(t *testing.T) {
	m := map[string]any{
		"hello": "world",
		"super": map[int]any{
			100: 10,
			200: nil,
		},
	}
	_, err := FromMapping(m)
	assert.Error(t, err)
}

func TestCreateFromMappingBadArray(t *testing.T) {
	m := map[string]any{
		"hello": "world",
		"super": map[string]any{
			"10":  10,
			"200": nil,
		},
		"stuff": []int{1, 2, 3},
	}
	_, err := FromMapping(m)
	assert.Error(t, err)
}

func TestStringArray(t *testing.T) {
	store := NewStore()
	err := store.Set([]string{"a", "b", "c"}, "arr")
	assert.NoError(t, err)
	s := store.GetStringArray("arr")
	assert.Equal(t, []string{"a", "b", "c"}, s)
}

func TestIntArray(t *testing.T) {
	store := NewStore()
	err := store.Set([]int{1, 2, 3}, "arr")
	assert.NoError(t, err)
	s := store.GetIntArray("arr")
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestFloatArray(t *testing.T) {
	store := NewStore()
	err := store.Set([]float64{1.5, 2.5, 3.5}, "arr")
	assert.NoError(t, err)
	s := store.GetFloatArray("arr")
	assert.Equal(t, []float64{1.5, 2.5, 3.5}, s)
}

func TestMappingArray(t *testing.T) {
	store := NewStore()
	err := store.Set([]map[string]any{
		{
			"a": 1,
			"b": 2,
		},
	}, "arr")
	assert.NoError(t, err)
	s := store.GetMappingArray("arr")
	assert.Equal(t, []map[string]any{
		{
			"a": 1,
			"b": 2,
		},
	}, s)
}

func TestOverlay(t *testing.T) {
	base := NewStore()
	err := base.Set("abc", "first", "second", "third")
	assert.NoError(t, err)
	err = base.Set(map[string]any{
		"a": "howdy",
		"b": 123,
		"c": 10.5,
	}, "first", "second", "third-b")
	assert.NoError(t, err)

	ovl := NewStore()
	err = ovl.Set("def", "first", "second-two")
	assert.NoError(t, err)
	err = ovl.Set([]any{"hello", "world"}, "first", "second", "third", "fourth")
	assert.NoError(t, err)

	final, err := base.Overylay(ovl)
	assert.NoError(t, err)
	thirdMapping := final.GetMapping("first", "second", "third")
	assert.Equal(t, len(thirdMapping), 1)
	assert.Equal(t, "howdy", final.GetString("first", "second", "third-b", "a"))
	assert.Equal(t, 123, final.GetInt("first", "second", "third-b", "b"))
	assert.Equal(t, "def", final.GetString("first", "second-two"))
	contents := final.GetStringArray("first", "second", "third", "fourth")
	assert.Len(t, contents, 2)
}

func TestNoNsSet(t *testing.T) {
	s := NewStore()
	err := s.Set("hello")
	assert.Error(t, err)
}

func TestIndexIdentifier(t *testing.T) {
	key, index, ok := parseArrayKey("test[100]")
	assert.Equal(t, "test", key)
	assert.Equal(t, 100, index)
	assert.True(t, ok)
}

func TestAccessNotation(t *testing.T) {
	s := NewStore()
	err := s.Set([]any{map[string]any{
		"x": 1,
		"y": 2,
		"z": 3,
	},
		map[string]any{
			"x": 10,
			"y": 20,
			"z": 30,
		}}, "first", "second")

	assert.NoError(t, err)
	value := s.GetInt(s.ParseNamespaceString("first.second[1].x")...)
	assert.Equal(t, 10, value)

	err = s.Set(22, s.ParseNamespaceString("first.second[0].y")...)
	assert.NoError(t, err)
	value = s.GetInt(s.ParseNamespaceString("first.second[0].y")...)
	assert.Equal(t, 22, value)

}
