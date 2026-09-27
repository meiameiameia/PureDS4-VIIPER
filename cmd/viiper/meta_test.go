package main

import (
	"strings"
	"testing"
)

func TestDescriptionIdentifiesPureDS4Fork(t *testing.T) {
	description := Description()
	if !strings.Contains(description, "Source:  https://github.com/meiameiameia/PureDS4-VIIPER") {
		t.Fatalf("Description() does not identify the PureDS4 fork: %q", description)
	}
}
