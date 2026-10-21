package useragent

import (
	"os"
	"strings"
	"testing"
)

// TestDocs_DeviceIdentityNote protects the guidance that a User-Agent is not a
// unique device identifier and callers must not treat it as one.
func TestDocs_DeviceIdentityNote(t *testing.T) {
	data, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	readme := strings.ToLower(string(data))
	for _, phrase := range []string{
		"not a unique device identifier",
		"device.vendor",
		"device.model",
		"application-level session",
	} {
		if !strings.Contains(readme, phrase) {
			t.Errorf("README.md is missing %q", phrase)
		}
	}
}
