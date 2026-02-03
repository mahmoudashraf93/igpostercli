package secrets

import (
	"errors"
	"testing"

	"github.com/99designs/keyring"
)

func withTestKeyring(t *testing.T) func() {
	t.Helper()

	dir := t.TempDir()

	return SetOpenKeyringForTests(func() (keyring.Keyring, error) {
		return keyring.Open(keyring.Config{
			ServiceName:      "igpostercli-test",
			AllowedBackends:  []keyring.BackendType{keyring.FileBackend},
			FileDir:          dir,
			FilePasswordFunc: func(string) (string, error) { return "test-pass", nil },
		})
	})
}

func TestAccessTokenRoundTrip(t *testing.T) {
	restore := withTestKeyring(t)
	defer restore()

	const profile = "test"
	const token = "tok123"

	_, ok, err := GetAccessToken(profile)
	if err != nil {
		t.Fatalf("get access token: %v", err)
	}

	if ok {
		t.Fatalf("expected no token")
	}

	err = SetAccessToken(profile, token)
	if err != nil {
		t.Fatalf("set access token: %v", err)
	}

	got, ok, err := GetAccessToken(profile)
	if err != nil {
		t.Fatalf("get access token: %v", err)
	}

	if !ok {
		t.Fatalf("expected token to be present")
	}

	if got != token {
		t.Fatalf("unexpected token: %s", got)
	}

	deleted, err := DeleteAccessToken(profile)
	if err != nil {
		t.Fatalf("delete access token: %v", err)
	}

	if !deleted {
		t.Fatalf("expected token to be deleted")
	}

	_, ok, err = GetAccessToken(profile)
	if err != nil {
		t.Fatalf("get access token: %v", err)
	}

	if ok {
		t.Fatalf("expected token to be removed")
	}
}

func TestDeleteMissingAccessToken(t *testing.T) {
	restore := withTestKeyring(t)
	defer restore()

	deleted, err := DeleteAccessToken("missing")
	if err == nil {
		t.Fatalf("expected error")
	}

	if !errors.Is(err, ErrTokenNotFound()) {
		t.Fatalf("unexpected error: %v", err)
	}

	if deleted {
		t.Fatalf("expected deleted=false")
	}
}
