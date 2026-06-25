package views

import (
	"testing"
	"time"
)

func TestFormatPrice(t *testing.T) {
	tests := []struct {
		price    string
		currency string
		expected string
	}{
		{"33590000.0", "VND", "33,590,000 ₫"},
		{"1000.0", "VND", "1,000 ₫"},
		{"0.0", "VND", "0 ₫"},
		{"invalid", "VND", "invalid"},
		{"  18690000.0  ", "VND", "18,690,000 ₫"},
		{"99.99", "USD", "$99.99"},
		{"49.99", "EUR", "€49.99"},
		{"29.99", "GBP", "£29.99"},
		{"100.00", "JPY", "JPY 100.00"},
		{"100.00", "", " 100.00"},
	}
	for _, tt := range tests {
		got := formatPrice(tt.price, tt.currency)
		if got != tt.expected {
			t.Errorf("formatPrice(%q, %q) = %q, want %q", tt.price, tt.currency, got, tt.expected)
		}
	}
}

func TestFormatMonthYear(t *testing.T) {
	ts := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	if got := FormatMonthYear(ts, "en"); got != "March 2024" {
		t.Errorf("expected 'March 2024', got %q", got)
	}
	if got := FormatMonthYear(ts, "vi"); got != "Tháng 3 năm 2024" {
		t.Errorf("expected 'Tháng 3 năm 2024', got %q", got)
	}
}

func TestFormatDateTime(t *testing.T) {
	// 14:30 UTC = 21:30 UTC+7
	ts := time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC)
	if got := FormatDateTime(ts, "en"); got != "Mar 15, 2024 21:30:00" {
		t.Errorf("expected 'Mar 15, 2024 21:30:00', got %q", got)
	}
	if got := FormatDateTime(ts, "vi"); got != "15 Tháng 3, 2024 21:30:00" {
		t.Errorf("expected '15 Tháng 3, 2024 21:30:00', got %q", got)
	}
}

func TestFormatDateTime_UTC7Conversion(t *testing.T) {
	// 07:37 UTC = 14:37 UTC+7 — mirrors the actual sensor log
	utc := time.Date(2026, 6, 21, 7, 37, 0, 0, time.UTC)
	if got := FormatDateTime(utc, "en"); got != "Jun 21, 2026 14:37:00" {
		t.Errorf("got %q, want %q", got, "Jun 21, 2026 14:37:00")
	}
}

func TestFormatDateTime_MidnightUTC(t *testing.T) {
	// 00:00 UTC = 07:00 UTC+7, same day
	utc := time.Date(2026, 6, 21, 0, 0, 0, 0, time.UTC)
	if got := FormatDateTime(utc, "en"); got != "Jun 21, 2026 07:00:00" {
		t.Errorf("got %q, want %q", got, "Jun 21, 2026 07:00:00")
	}
}

func TestFormatDateTime_DayRollover(t *testing.T) {
	// 18:00 UTC = 01:00 UTC+7 next day
	utc := time.Date(2026, 6, 21, 18, 0, 0, 0, time.UTC)
	if got := FormatDateTime(utc, "en"); got != "Jun 22, 2026 01:00:00" {
		t.Errorf("got %q, want %q", got, "Jun 22, 2026 01:00:00")
	}
}

func TestPaginationPages(t *testing.T) {
	tests := []struct {
		current  int
		total    int
		expected []int
	}{
		{1, 1, nil},
		{1, 5, []int{1, 2, 3, 0, 5}},
		{1, 20, []int{1, 2, 3, 0, 20}},
		{8, 20, []int{1, 0, 6, 7, 8, 9, 10, 0, 20}},
		{19, 20, []int{1, 0, 17, 18, 19, 20}},
		{20, 20, []int{1, 0, 18, 19, 20}},
		{3, 5, []int{1, 2, 3, 4, 5}},
	}
	for _, tt := range tests {
		got := paginationPages(tt.current, tt.total)
		if !equalSlices(got, tt.expected) {
			t.Errorf("paginationPages(%d, %d) = %v, want %v", tt.current, tt.total, got, tt.expected)
		}
	}
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCalculateSoilPercentage(t *testing.T) {
	tests := []struct {
		raw  int16
		dry  int
		wet  int
		want float32
	}{
		{3500, 3500, 1760, 0.0},   // at dry threshold
		{3600, 3500, 1760, 0.0},   // above dry — clamp to 0
		{1760, 3500, 1760, 100.0}, // at wet threshold
		{1500, 3500, 1760, 100.0}, // below wet — clamp to 100
		{2630, 3500, 1760, 50.0},  // midpoint
	}
	for _, tt := range tests {
		got := calculateSoilPercentage(tt.raw, tt.dry, tt.wet)
		if got != tt.want {
			t.Errorf("calculateSoilPercentage(%d, %d, %d) = %.1f, want %.1f",
				tt.raw, tt.dry, tt.wet, got, tt.want)
		}
	}
}
