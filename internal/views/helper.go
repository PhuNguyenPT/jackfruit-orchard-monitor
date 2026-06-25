package views

import (
	"GoApp/internal/database"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var vietnamTZ = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)

func formatPrice(price string, currency string) string {
	f, err := strconv.ParseFloat(strings.TrimSpace(price), 64)
	if err != nil {
		return price
	}

	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "VND", "₫":
		s := strconv.FormatInt(int64(f), 10)
		result := []byte{}
		for i, ch := range s {
			if i > 0 && (len(s)-i)%3 == 0 {
				result = append(result, ',')
			}
			result = append(result, byte(ch))
		}
		return fmt.Sprintf("%s ₫", string(result))
	case "USD":
		return fmt.Sprintf("$%.2f", f)
	case "EUR":
		return fmt.Sprintf("€%.2f", f)
	case "GBP":
		return fmt.Sprintf("£%.2f", f)
	default:
		// fallback: use currency code as prefix with 2 decimal places
		return fmt.Sprintf("%s %.2f", strings.ToUpper(currency), f)
	}
}

func FormatMonthYear(t time.Time, lang string) string {
	if lang == "vi" {
		months := []string{
			"Tháng 1", "Tháng 2", "Tháng 3", "Tháng 4",
			"Tháng 5", "Tháng 6", "Tháng 7", "Tháng 8",
			"Tháng 9", "Tháng 10", "Tháng 11", "Tháng 12",
		}
		return months[t.Month()-1] + " năm " + strconv.Itoa(t.Year())
	}
	return t.Format("January 2006")
}

func FormatDateTime(t time.Time, lang string) string {
	t = t.In(vietnamTZ)
	if lang == "vi" {
		months := []string{
			"Tháng 1", "Tháng 2", "Tháng 3", "Tháng 4",
			"Tháng 5", "Tháng 6", "Tháng 7", "Tháng 8",
			"Tháng 9", "Tháng 10", "Tháng 11", "Tháng 12",
		}
		return strconv.Itoa(t.Day()) + " " + months[t.Month()-1] + ", " + strconv.Itoa(t.Year()) + " " + t.Format("15:04:05")
	}
	return t.Format("Jan 2, 2006 15:04:05")
}

func paginationPages(current, total int) []int {
	if total <= 1 {
		return nil
	}

	set := map[int]bool{}
	pages := []int{}

	add := func(p int) {
		if p >= 1 && p <= total && !set[p] {
			set[p] = true
		}
	}

	add(1)
	add(total)
	add(current)
	for i := -2; i <= 2; i++ {
		add(current + i)
	}

	// build sorted list
	sorted := []int{}
	for p := 1; p <= total; p++ {
		if set[p] {
			sorted = append(sorted, p)
		}
	}

	// insert 0 as ellipsis where gaps exist
	for i, p := range sorted {
		if i == 0 {
			pages = append(pages, p)
			continue
		}
		if p-pages[len(pages)-1] > 1 {
			pages = append(pages, 0) // ellipsis
		}
		pages = append(pages, p)
	}

	return pages
}

func calculateSoilPercentage(raw int16, dry, wet int) float32 {
	dryVal := int16(dry)
	wetVal := int16(wet)
	if raw >= dryVal {
		return 0.0
	}
	if raw <= wetVal {
		return 100.0
	}
	return float32(dryVal-raw) / float32(dryVal-wetVal) * 100.0
}

func marshalSHT40Rows(rows []database.GetAirTempHumidReadingsByAddrRow, lang string) string {
	type point struct {
		T     string  `json:"t"`
		Temp  float32 `json:"temp"`
		Humid float32 `json:"humid"`
	}
	format := "02-01 15:04:05" // DD-MM for VI
	if lang != "vi" {
		format = "01-02 15:04:05" // MM-DD for EN
	}
	pts := make([]point, len(rows))
	for i, r := range rows {
		pts[len(rows)-1-i] = point{
			T:     r.CreatedAt.In(vietnamTZ).Format(format),
			Temp:  float32(r.Temperature) / 10,
			Humid: float32(r.Humidity) / 10,
		}
	}
	b, err := json.Marshal(pts)
	if err != nil {
		log.Printf("[views] marshalSHT40Rows: %v", err)
		return "[]"
	}
	return string(b)
}

func marshalSoilRows(rows []database.GetSoilMoistureReadingsBySensorIdxRow, soilDry, soilWet int, lang string) string {
	type point struct {
		T   string  `json:"t"`
		Pct float32 `json:"pct"`
		Raw int16   `json:"raw"`
	}
	format := "02-01 15:04:05"
	if lang != "vi" {
		format = "01-02 15:04:05"
	}
	pts := make([]point, len(rows))
	for i, r := range rows {
		pts[len(rows)-1-i] = point{
			T:   r.CreatedAt.In(vietnamTZ).Format(format),
			Pct: calculateSoilPercentage(r.Raw, soilDry, soilWet),
			Raw: r.Raw,
		}
	}
	b, err := json.Marshal(pts)
	if err != nil {
		log.Printf("[views] marshalSoilRows: %v", err)
		return "[]"
	}
	return string(b)
}
