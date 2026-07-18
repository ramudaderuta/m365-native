package auth

import "testing"

func TestBrowserPKCEDefaultsRemainMatched(t *testing.T) {
	for _, key := range []string{
		"M365_CLIENT_ID",
		"M365_AUTHORITY",
		"M365_REDIRECT_URI",
	} {
		t.Setenv(key, "")
	}

	if got, want := ClientID(), "c0ab8ce9-e9a0-42e7-b064-33d422df41f1"; got != want {
		t.Fatalf("ClientID() = %q, want %q", got, want)
	}
	if got, want := Authority(), "https://login.microsoftonline.com/common"; got != want {
		t.Fatalf("Authority() = %q, want %q", got, want)
	}
	if got, want := RedirectURI(), "https://login.microsoftonline.com/common/oauth2/nativeclient"; got != want {
		t.Fatalf("RedirectURI() = %q, want %q", got, want)
	}
}
