package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNormalizeProfileName(t *testing.T) {
	got, err := NormalizeProfileName("Brand_A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "brand_a" {
		t.Fatalf("unexpected name: %s", got)
	}

	_, err = NormalizeProfileName("bad name")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestProfilesReadWrite(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", base)
	t.Setenv("HOME", base)

	cfg := ProfilesFile{
		Profiles: map[string]Profile{
			"default": {
				IGUserID:   "ig1",
				PageID:     "page1",
				BusinessID: "biz1",
			},
		},
	}

	if err := WriteProfiles(cfg); err != nil {
		t.Fatalf("write profiles: %v", err)
	}

	path, err := ProfilesPath()
	if err != nil {
		t.Fatalf("profiles path: %v", err)
	}

	if _, statErr := os.Stat(path); statErr != nil {
		t.Fatalf("expected profiles file to exist: %v", statErr)
	}

	read, err := ReadProfiles()
	if err != nil {
		t.Fatalf("read profiles: %v", err)
	}

	if !reflect.DeepEqual(cfg, read) {
		t.Fatalf("unexpected profiles: %#v", read)
	}
}

func TestListProfiles(t *testing.T) {
	cfg := ProfilesFile{
		Profiles: map[string]Profile{
			"b": {},
			"a": {},
		},
	}

	got := ListProfiles(cfg)
	expected := []string{"a", "b"}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("unexpected list: %#v", got)
	}
}

func TestProfilesPathUsesConfigDir(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", base)
	t.Setenv("HOME", base)

	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("user config dir: %v", err)
	}

	path, err := ProfilesPath()
	if err != nil {
		t.Fatalf("profiles path: %v", err)
	}

	expectedDir := filepath.Join(configDir, AppName)
	if filepath.Dir(path) != expectedDir {
		t.Fatalf("unexpected profiles dir: %s", filepath.Dir(path))
	}
}
