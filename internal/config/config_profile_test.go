package config

import (
	"testing"

	"github.com/99designs/keyring"

	"github.com/mahmoud/igpostercli/internal/secrets"
)

func withTestKeyring(t *testing.T) func() {
	t.Helper()

	dir := t.TempDir()

	return secrets.SetOpenKeyringForTests(func() (keyring.Keyring, error) {
		return keyring.Open(keyring.Config{
			ServiceName:      "igpostercli-test",
			AllowedBackends:  []keyring.BackendType{keyring.FileBackend},
			FileDir:          dir,
			FilePasswordFunc: func(string) (string, error) { return "test-pass", nil },
		})
	})
}

func TestLoadWithProfileOverridesEnv(t *testing.T) {
	restore := withTestKeyring(t)
	defer restore()

	base := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", base)
	t.Setenv("HOME", base)

	t.Setenv("IG_APP_ID", "app")
	t.Setenv("IG_APP_SECRET", "secret")
	t.Setenv("IG_ACCESS_TOKEN", "env-token")
	t.Setenv("IG_USER_ID", "env-user")
	t.Setenv("IG_PAGE_ID", "env-page")
	t.Setenv("IG_BUSINESS_ID", "env-biz")

	profiles := ProfilesFile{
		Profiles: map[string]Profile{
			"agent": {
				IGUserID:   "profile-user",
				PageID:     "profile-page",
				BusinessID: "profile-biz",
			},
		},
	}

	if err := WriteProfiles(profiles); err != nil {
		t.Fatalf("write profiles: %v", err)
	}

	if err := secrets.SetAccessToken("agent", "keychain-token"); err != nil {
		t.Fatalf("set access token: %v", err)
	}

	cfg, err := LoadWithProfile("Agent")
	if err != nil {
		t.Fatalf("load with profile: %v", err)
	}

	if cfg.AppID != "app" {
		t.Fatalf("unexpected app id: %s", cfg.AppID)
	}

	if cfg.AppSecret != "secret" {
		t.Fatalf("unexpected app secret: %s", cfg.AppSecret)
	}

	if cfg.AccessToken != "keychain-token" {
		t.Fatalf("unexpected access token: %s", cfg.AccessToken)
	}

	if cfg.IGUserID != "profile-user" {
		t.Fatalf("unexpected user id: %s", cfg.IGUserID)
	}

	if cfg.PageID != "profile-page" {
		t.Fatalf("unexpected page id: %s", cfg.PageID)
	}

	if cfg.BusinessID != "profile-biz" {
		t.Fatalf("unexpected business id: %s", cfg.BusinessID)
	}
}

func TestLoadWithProfileFallsBackToEnv(t *testing.T) {
	restore := withTestKeyring(t)
	defer restore()

	base := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", base)
	t.Setenv("HOME", base)

	t.Setenv("IG_ACCESS_TOKEN", "env-token")
	t.Setenv("IG_USER_ID", "env-user")
	t.Setenv("IG_PAGE_ID", "env-page")
	t.Setenv("IG_BUSINESS_ID", "env-biz")

	cfg, err := LoadWithProfile("missing")
	if err != nil {
		t.Fatalf("load with profile: %v", err)
	}

	if cfg.AccessToken != "env-token" {
		t.Fatalf("unexpected access token: %s", cfg.AccessToken)
	}

	if cfg.IGUserID != "env-user" {
		t.Fatalf("unexpected user id: %s", cfg.IGUserID)
	}

	if cfg.PageID != "env-page" {
		t.Fatalf("unexpected page id: %s", cfg.PageID)
	}

	if cfg.BusinessID != "env-biz" {
		t.Fatalf("unexpected business id: %s", cfg.BusinessID)
	}
}
