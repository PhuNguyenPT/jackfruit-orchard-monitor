package views

import (
	"context"
	"strings"
	"testing"
)

func TestHeader_LanguageSelectorUsesClickToggle(t *testing.T) {
	var sb strings.Builder
	err := Header("home", "", "en").Render(context.Background(), &sb)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	html := sb.String()

	// Must NOT rely on CSS-only hover (the original mobile bug)
	if strings.Contains(html, "group-hover:block") {
		t.Error("language selector still uses group-hover; this doesn't work on touch devices")
	}

	// Must have the click-toggle wiring
	if !strings.Contains(html, `id="lang-menu-btn"`) {
		t.Error("missing lang-menu-btn id")
	}
	if !strings.Contains(html, `data-target="lang-menu"`) {
		t.Error("missing data-target=lang-menu on the toggle button")
	}
	if !strings.Contains(html, `id="lang-menu"`) {
		t.Error("missing lang-menu target element")
	}

	// Dropdown must start hidden so it doesn't flash open before JS runs
	if !strings.Contains(html, `id="lang-menu" class="hidden`) {
		t.Error("lang-menu should start with the hidden class")
	}
}

func TestHeader_ShowsCorrectLanguageFlag(t *testing.T) {
	cases := []struct {
		lang string
		want string
	}{
		{"en", "🇬🇧 EN"},
		{"vi", "🇻🇳 VI"},
	}
	for _, tc := range cases {
		var sb strings.Builder
		if err := Header("home", "", tc.lang).Render(context.Background(), &sb); err != nil {
			t.Fatalf("render failed: %v", err)
		}
		if !strings.Contains(sb.String(), tc.want) {
			t.Errorf("lang=%s: expected %q in output", tc.lang, tc.want)
		}
	}
}
