package auth

import (
	"strings"
	"testing"
)

func TestBrowserPKCEDefaultsRemainMatched(t *testing.T) {
	for _, key := range []string{
		"M365_CLIENT_ID",
		"M365_AUTHORITY",
		"M365_REDIRECT_URI",
		"M365_SCOPE",
	} {
		t.Setenv(key, "")
	}

	if got, want := ClientID(), "efcea265-005c-4f0a-97c2-b3ab369c8484"; got != want {
		t.Fatalf("ClientID() = %q, want %q", got, want)
	}
	if got, want := Authority(), "https://login.microsoftonline.com/common"; got != want {
		t.Fatalf("Authority() = %q, want %q", got, want)
	}
	if got, want := RedirectURI(), "http://127.0.0.1:4141/api/auth/callback"; got != want {
		t.Fatalf("RedirectURI() = %q, want %q", got, want)
	}
	if got, want := Scope(), "openid profile offline_access https://substrate.office.com/sydney/.default"; got != want {
		t.Fatalf("Scope() = %q, want %q", got, want)
	}
}

func TestDefaultAuthorityIsMultitenant(t *testing.T) {
	t.Setenv("M365_AUTHORITY", "")
	if got := Authority(); got != "https://login.microsoftonline.com/common" {
		t.Fatalf("Authority() = %q", got)
	}
	if !strings.Contains(AuthorizeEndpoint(), "/common/") {
		t.Fatal("default authorize endpoint must be multitenant")
	}
}

func TestAuthorityOverride(t *testing.T) {
	const custom = "https://login.microsoftonline.com/organizations"
	t.Setenv("M365_AUTHORITY", custom)
	if got := Authority(); got != custom {
		t.Fatalf("Authority() = %q", got)
	}
}
