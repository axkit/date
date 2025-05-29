package date_test

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/axkit/date"
)

func TestNew(t *testing.T) {
	tests := []struct {
		year   int
		month  time.Month
		day    int
		expect date.Date
	}{
		{2023, time.January, 1, date.New(2023, time.January, 1)},
		{0, 0, 0, date.Null()},
	}

	for _, tt := range tests {
		if result := date.New(tt.year, tt.month, tt.day); result != tt.expect {
			t.Errorf("New(%d, %d, %d) = %v; want %v", tt.year, tt.month, tt.day, result, tt.expect)
		}
	}
}

func TestToday(t *testing.T) {
	today := date.Today()
	now := time.Now()
	expected := date.New(now.Year(), now.Month(), now.Day())

	if today != expected {
		t.Errorf("Today() = %v; want %v", today, expected)
	}
}

func TestDate_String(t *testing.T) {
	tests := []struct {
		date   date.Date
		expect string
	}{
		{date.New(2023, time.January, 1), "2023-01-01"},
		{date.Null(), ""},
	}

	for _, tt := range tests {
		if result := tt.date.String(); result != tt.expect {
			fmt.Println("result:", result)
			t.Errorf("Date.String() = %v; want %v", result, tt.expect)
		}
	}
}

func TestDate_Add(t *testing.T) {
	tests := []struct {
		date   date.Date
		years  int
		months int
		days   int
		expect date.Date
	}{
		{date.New(2023, time.January, 1), 1, 0, 0, date.New(2024, time.January, 1)},
		{date.New(2023, time.January, 1), 0, 1, 0, date.New(2023, time.February, 1)},
		{date.New(2023, time.January, 1), 0, 0, 1, date.New(2023, time.January, 2)},
	}

	for _, tt := range tests {
		if result := tt.date.Add(tt.years, tt.months, tt.days); result != tt.expect {
			t.Errorf("Date.Add(%d, %d, %d) = %v; want %v", tt.years, tt.months, tt.days, result, tt.expect)
		}
	}
}

func TestDate_Parse(t *testing.T) {
	tests := []struct {
		input  string
		expect date.Date
		err    bool
	}{
		{"2023-01-01", date.New(2023, time.January, 1), false},
		{"0000-00-00", date.Null(), true},
	}

	for _, tt := range tests {
		var d date.Date
		err := d.Parse(tt.input)
		if tt.err {
			if err == nil {
				t.Errorf("Date.Parse(%q) expected error; got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("Date.Parse(%q) unexpected error: %v", tt.input, err)
			}
			if d != tt.expect {
				t.Errorf("Date.Parse(%q) = %v; want %v", tt.input, d, tt.expect)
			}
		}
	}
}

func TestDate_Year(t *testing.T) {
	date := date.New(2023, time.January, 1)
	if year := date.Year(); year != 2023 {
		t.Errorf("Date.Year() = %d; want %d", year, 2023)
	}
}

func TestDate_Month(t *testing.T) {
	date := date.New(2023, time.January, 1)
	if month := date.Month(); month != time.January {
		t.Errorf("Date.Month() = %v; want %v", month, time.January)
	}
}

func TestDate_Day(t *testing.T) {
	date := date.New(2023, time.January, 1)
	if day := date.Day(); day != 1 {
		t.Errorf("Date.Day() = %d; want %d", day, 1)
	}
}

func TestDate_Value(t *testing.T) {
	tests := []struct {
		date   date.Date
		expect driver.Value
	}{
		{date.New(2023, time.January, 1), []byte("2023-01-01")},
		{date.Null(), nil},
	}

	for _, tt := range tests {
		value, _ := tt.date.Value()

		if !reflect.DeepEqual(value, tt.expect) {
			t.Errorf("Date.Value() = %v; want %v", value, tt.expect)
		}
	}
}

func TestDate_Valid(t *testing.T) {
	tests := []struct {
		date   date.Date
		expect bool
	}{
		{date.New(2023, time.January, 1), true},
		{date.Null(), false},
	}

	for _, tt := range tests {
		if valid := tt.date.Valid(); valid != tt.expect {
			t.Errorf("Date.Valid() = %v; want %v", valid, tt.expect)
		}
	}
}

func TestDate_Scan(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect date.Date
		err    bool
	}{
		{time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local), date.New(2023, time.January, 1), false},
		{nil, date.Null(), false},
		{"Invalid", date.Null(), true},
	}

	for _, tt := range tests {
		var d date.Date
		err := d.Scan(tt.input)
		if tt.err {
			if err == nil {
				t.Errorf("Date.Scan(%v) expected error; got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("Date.Scan(%v) unexpected error: %v", tt.input, err)
			}
			if d != tt.expect {
				t.Errorf("Date.Scan(%v) = %v; want %v", tt.input, d, tt.expect)
			}
		}
	}
}

func TestInitPreformattedValues(t *testing.T) {
	from := date.New(2023, time.January, 1)
	to := date.New(2023, time.January, 5)
	date.InitPreformattedValues(from, to)
	if from.String() != "2023-01-01" {
		t.Errorf("Date.String() = %v; want %v", from.String(), "2023-01-01")
	}
	if to.String() != "2023-01-05" {
		t.Errorf("Date.String() = %v; want %v", to.String(), "2023-01-05")
	}
}
func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		date   date.Date
		expect []byte
		err    bool
	}{
		{date.New(2023, time.January, 1), []byte(`"2023-01-01"`), false},
		{date.Null(), []byte("null"), false},
	}

	for _, tt := range tests {
		result, err := tt.date.MarshalJSON()
		if tt.err {
			if err == nil {
				t.Errorf("Date.MarshalJSON() expected error; got nil")
			}
		} else {
			if err != nil {
				t.Errorf("Date.MarshalJSON() unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expect) {
				t.Errorf("Date.MarshalJSON() = %v; want %v", string(result), string(tt.expect))
			}
		}
	}
}
func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input  []byte
		expect date.Date
		err    bool
	}{
		{[]byte(`"2023-01-01"`), date.New(2023, time.January, 1), false},
		{[]byte("null"), date.Null(), false},
		{[]byte("invalid"), date.Null(), true},
		{[]byte(`""`), date.Null(), false},
	}

	for _, tt := range tests {
		var d date.Date
		err := d.UnmarshalJSON(tt.input)
		if tt.err {
			if err == nil {
				t.Errorf("Date.UnmarshalJSON(%v) expected error; got nil", string(tt.input))
			}
		} else {
			if err != nil {
				t.Errorf("Date.UnmarshalJSON(%v) unexpected error: %v", string(tt.input), err)
			}
			if d != tt.expect {
				t.Errorf("Date.UnmarshalJSON(%v) = %v; want %v", string(tt.input), d, tt.expect)
			}
		}
	}
}

// Generate benchmarks for date.Date methods
func BenchmarkDate_New(b *testing.B) {
	for b.Loop() {
		_ = date.New(2023, time.January, 1)
	}
}

func BenchmarkDate_Today(b *testing.B) {
	for b.Loop() {
		_ = date.Today()
	}
}

func TestMain(m *testing.M) {
	date.InitPreformattedValues(date.New(2022, time.January, 1), date.New(2024, time.January, 5))
	m.Run()
}

func BenchmarkDate_String(b *testing.B) {
	d := date.New(2023, time.January, 1)
	for b.Loop() {
		_ = d.String()
	}
}

func BenchmarkDate_Parse(b *testing.B) {
	for b.Loop() {
		var d date.Date
		_ = d.Parse("2023-01-01")
	}
}

func BenchmarkDate_Year(b *testing.B) {
	d := date.New(2023, time.January, 1)
	for b.Loop() {
		_ = d.Year()
	}
}

func BenchmarkDate_Month(b *testing.B) {
	d := date.New(2023, time.January, 1)
	for b.Loop() {
		_ = d.Month()
	}
}

func BenchmarkDate_Day(b *testing.B) {
	d := date.New(2023, time.January, 1)
	for b.Loop() {
		_ = d.Day()
	}
}

func BenchmarkDate_Value(b *testing.B) {
	d := date.New(2023, time.January, 1)
	for b.Loop() {
		_, _ = d.Value()
	}
}

func BenchmarkDate_Valid(b *testing.B) {
	d := date.New(2023, time.January, 1)
	for b.Loop() {
		_ = d.Valid()
	}
}

func BenchmarkDate_Scan(b *testing.B) {
	for b.Loop() {
		var d date.Date
		_ = d.Scan(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local))
	}
}

func BenchmarkDate_MarshalJSON(b *testing.B) {
	d := date.New(2023, time.January, 1)
	for b.Loop() {
		_, _ = d.MarshalJSON()
	}
}

func BenchmarkDate_UnmarshalJSON(b *testing.B) {
	input := []byte(`"2023-01-01"`)
	for b.Loop() {
		var d date.Date
		_ = d.UnmarshalJSON(input)
	}
}
