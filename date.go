// Package date provides a custom date type that omits time and timezone components,
// serving as a replacement for the standard time.Time type where only the date part is needed.
package date

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// Date represents a nullable date type that excludes time and timezone information.
// Internally, it stores the year, month, and day as a combined hex integer (0xYYYYMMDD),
// making it comparable and sortable as an integer. In a database, it is stored as DATE.
// A null (empty) value is represented in memory as 0.
type Date uint32

// Predefined variables
var (
	Separator         byte = '-'
	DatabaseSeparator byte = '-'
	tmpl                   = [10]byte{'0', '0', '0', '0', Separator, '0', '0', Separator, '0', '0'}
	tmplDB                 = [10]byte{'0', '0', '0', '0', DatabaseSeparator, '0', '0', DatabaseSeparator, '0', '0'}
	tmplJS                 = [12]byte{'"', '0', '0', '0', '0', Separator, '0', '0', Separator, '0', '0', '"'}

	pfm map[Date]string
	pjm map[string]Date

	null  = []byte("null")
	cache string

	// FiveYearBefore represents the date five years before today.
	FiveYearBefore = Today().Add(-5, 0, 0)
	// FiveYearAfter represents the date five years from today.
	FiveYearAfter = Today().Add(5, 0, 0)
)

// InitPreformattedValues initializes preformatted date values between a specified range.
// This function populates two maps, pfm and pjm, for date-to-string and string-to-date
// conversions, respectively, to optimize date marshaling and unmarshaling.
func InitPreformattedValues(from, to Date) {
	pfm = make(map[Date]string, (to.Time().Sub(from.Time()))/(24*time.Hour)+1)
	pjm = make(map[string]Date, (to.Time().Sub(from.Time()))/(24*time.Hour)+1)
	d := from
	for d < to {
		d = d.Add(0, 0, 1)
		pfm[d] = d.String()
		pjm[d.String()] = d
	}
}

// Time converts a Date to a time.Time instance in the local timezone.
func (d Date) Time() time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
}

// UTC converts a Date to a time.Time instance in UTC.
func (d Date) UTC() time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

// In converts a Date to a time.Time instance in the specified location.
func (d Date) In(loc *time.Location) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
}

// String returns the date as a string formatted as "YYYY-MM-DD".
// If the date is null, an empty string is returned.
func (d Date) String() string {
	if s, ok := pfm[d]; ok {
		return s
	}
	if !d.Valid() {
		return ""
	}
	return d.string()
}

// string is a private helper that returns the Date as a formatted string "YYYY-MM-DD".
func (d Date) string() string {
	var buf [10]byte
	d.byteArr(&buf)
	return string(buf[:])
}

func (d Date) byteArr(res *[10]byte) {

	*res = tmpl
	i := uint32(d)
	zero := byte('0')

	res[0] = byte((i>>28)&0x0000000F) + zero
	res[1] = byte((i>>24)&0x0000000F) + zero
	res[2] = byte((i>>20)&0x0000000F) + zero
	res[3] = byte((i>>16)&0x0000000F) + zero

	res[5] = byte((i>>12)&0x0000000F) + zero
	res[6] = byte((i>>8)&0x0000000F) + zero

	res[8] = byte((i>>4)&0x0000000F) + zero
	res[9] = byte(i&0x0000000F) + zero
}

// Add returns a new Date that is the result of adding the specified number
// of years, months, and days to the original Date.
func (d Date) Add(years, months, days int) Date {
	t := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	t = t.AddDate(years, months, days)
	return New(t.Year(), t.Month(), t.Day())
}

// Parse sets the Date from a string in the format "YYYY-MM-DD".
func (d *Date) Parse(s string) error {
	dt, err := Parse(s)
	if err != nil {
		return err
	}
	*d = dt
	return nil
}

// Year returns the year component of the Date.
func (d Date) Year() int {
	i := uint32(d)
	return int((i>>28)*1000 + ((i>>24)&0x000F)*100 + ((i>>20)&0x000F)*10 + (i>>16)&0x000F)
}

// Month returns the month component of the Date.
func (d Date) Month() time.Month {
	i := (uint32(d) >> 8) & 0x000000FF
	return time.Month((i>>4)*10 + i&0x0F)
}

// Day returns the day component of the Date.
func (d Date) Day() int {
	i := uint32(d) & 0x000000FF
	return int((i>>4)*10 + i&0x0F)
}

// Parse converts a string in "YYYY-MM-DD" format into a Date.
// If the input string is enclosed in quotes, they are removed before parsing.
func Parse(s string) (Date, error) {
	if s[0] == '"' {
		s = s[1:]
	}
	t, err := parseYYYYMMDD([]byte(s))
	if err != nil {
		return 0, err
	}
	return newDate(t.Year(), t.Month(), t.Day()), nil
}

// Null returns a null Date (equivalent to zero).
func Null() Date {
	return 0
}

// Value returns the Date as a driver.Value for database storage.
func (d Date) Value() (driver.Value, error) {
	if !d.Valid() {
		return nil, nil
	}
	return []byte(d.String()), nil
}

// Valid returns false if the Date is null.
func (d Date) Valid() bool {
	return d > 0
}

// Scan sets the Date from a database value.
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = 0
		return nil
	}
	v, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("Date.Scan: expected date.Date, got %T (%q)", value, value)
	}
	*d = New(v.Year(), v.Month(), v.Day())
	return nil
}

// Today returns the current date as a Date instance.
func Today() Date {
	t := time.Now()
	return newDate(t.Year(), t.Month(), t.Day())
}

// New creates a Date from the specified year, month, and day.
func New(y int, m time.Month, d int) Date {
	return newDate(y, m, d)
}

// newDate creates a Date by encoding the year, month, and day into a hex integer.
func newDate(y int, m time.Month, d int) Date {
	return Date(dec2hexy(uint32(y))<<16 | (dec2hex(uint32(m)) << 8) | dec2hex(uint32(d)))
}

// dec2hex converts a month or day (0-31) to hex, e.g., 31 to 0x31.
func dec2hex(i uint32) uint32 {
	return (i/10)*16 + i%10
}

// dec2hexy converts a year (e.g., 2018) to hex, e.g., 2018 to 0x2018.
func dec2hexy(i uint32) uint32 {
	return (i/1000)*16*16*16 + ((i%1000)/100)*16*16 + (i%100)/10*16 + i%10
}

// parseYYYYMMDD decodes a byte array in the format "YYYY-MM-DD" to a Date.
func parseYYYYMMDD(b []byte) (Date, error) {
	if len(b) < 10 {
		return Null(), errors.New("input date length less than 10 bytes")
	}
	y := int(b[0]-'0')*1000 + int(b[1]-'0')*100 + int(b[2]-'0')*10 + int(b[3]-'0')
	m := int(b[5]-'0')*10 + int(b[6]-'0')
	d := int(b[8]-'0')*10 + int(b[9]-'0')
	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	if t.Year() != y || t.Month() != time.Month(m) || t.Day() != d {
		return Null(), errors.New("date parse error")
	}
	return newDate(y, time.Month(m), d), nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if !d.Valid() {
		return null, nil
	}

	if s, ok := pfm[d]; ok {
		return []byte(s), nil
	}

	var buf [10]byte
	d.byteArr(&buf)

	return buf[:], nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	if len(b) == 0 ||
		(len(b) == 4 && b[0] == 'n' && b[1] == 'u' && b[2] == 'l' && b[3] == 'l') ||
		(len(b) == 2 && b[0] == '"' && b[1] == '"') {
		*d = 0
		return nil
	}

	pd, err := parseYYYYMMDD(b)
	if err != nil {
		return err
	}
	*d = pd
	return nil
}
