package date_test

import (
	"encoding/json"
	"testing"
	"time"

	"net/url"

	"github.com/axkit/date"
	"github.com/gorilla/schema"
)

func TestDate_Today(t *testing.T) {
	d := date.Today()
	t.Log(d)
}

func TestNewYMD(t *testing.T) {
	d := date.New(2018, 1, 31)
	if d.String() != "2018-01-31" {
		t.Error("NewYMD failed ", d.String())
	}
}

func TestAddDate(t *testing.T) {
	d := date.New(2018, 1, 31)
	d = d.Add(0, 0, 1)

	if d != date.New(2018, 2, 1) {
		t.Error("AddDate() failed", d.String())
	}
}

func TestDay(t *testing.T) {
	d := date.New(2018, 1, 31)
	if d.Day() != 31 {
		t.Error("Day() failed. Expected 31, got: ", d.Day())
	}
}

func TestMonth(t *testing.T) {
	d := date.New(2018, 1, 31)
	if d.Month() != 1 {
		t.Error("Day() failed. Expected 1, got: ", d.Day())
	}
}

func TestYear(t *testing.T) {
	d := date.New(2018, 1, 31)
	if d.Year() != 2018 {
		t.Error("Year() failed. Expected 2018, got: ", d.Year())
	}
}

func TestString(t *testing.T) {
	d := date.Date(uint32(0x20180131))

	if s := d.String(); s != "2018-01-31" {
		t.Error("String() failed", s)
	}
}

func TestParse(t *testing.T) {

	m := map[string]struct {
		d  date.Date
		ok bool
	}{
		"2017-06-02": {date.New(2017, time.Month(6), 2), true},
		"2017-02-28": {date.New(2017, time.Month(2), 28), true},
		"2017-22-28": {date.New(2017, time.Month(22), 28), false},
	}

	for k, v := range m {
		d, err := date.Parse(k)
		if err != nil && v.ok == false {
			continue
		}

		if err == nil && v.ok == false {
			t.Errorf("Test case %s failed!", k)
			continue
		}

		if v.d != d {
			t.Errorf(k, v)
		}
	}
}

func TestDateMarshal(t *testing.T) {

	d := date.New(2019, 3, 1)
	buf, err := json.Marshal(d)

	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(buf))
	}

	if string(buf) != `"2019-03-01"` {
		t.Error("marshal failed. got:", string(buf))
	}

}

func TestDateUnmarshal(t *testing.T) {

	var d struct {
		D date.Date `json:"d"`
		N date.Date `json:"n"`
		E date.Date `json:"e"`
	}

	d.N = date.Today()

	j := []byte(`{"d" : "2018-01-31", "n": null}`)

	err := json.Unmarshal(j, &d)

	if err != nil {
		t.Error(err)
	}

	if d.N.Valid() {
		t.Fatal("ашибка")
	}

	t.Log(d.D.String())
}

func BenchmarkDateStringPreCached(b *testing.B) {

	date.InitPreformattedValues(date.FiveYearBefore, date.FiveYearAfter)
	b.ResetTimer()

	d := date.Today()
	for n := 0; n < b.N; n++ {
		s := d.String()
		_ = s
	}
}

func BenchmarkDateString(b *testing.B) {

	date.InitPreformattedValues(date.FiveYearBefore, date.FiveYearAfter)
	b.ResetTimer()

	d := date.Today().Add(10, 0, 0)
	for n := 0; n < b.N; n++ {
		s := d.String()
		_ = s
	}
}

func BenchmarkTimeString(b *testing.B) {
	t := time.Now()
	for n := 0; n < b.N; n++ {
		_ = t.String()
	}
}

func BenchmarkDateMarshalPreCached(b *testing.B) {

	date.InitPreformattedValues(date.FiveYearBefore, date.FiveYearAfter)
	b.ResetTimer()

	d := date.Today()
	for n := 0; n < b.N; n++ {
		_, _ = json.Marshal(&d)
	}
}

func BenchmarkDateMarshal(b *testing.B) {

	date.InitPreformattedValues(date.FiveYearBefore, date.FiveYearAfter)
	b.ResetTimer()

	d := date.Today().Add(10, 0, 0)
	for n := 0; n < b.N; n++ {
		_, _ = json.Marshal(&d)
	}
}

func BenchmarkTimeMarshal(b *testing.B) {

	date.InitPreformattedValues(date.FiveYearBefore, date.FiveYearAfter)
	b.ResetTimer()

	d := time.Now()
	for n := 0; n < b.N; n++ {
		_, _ = json.Marshal(d)
	}

}

func BenchmarkDatePreCachedUnmarshal(b *testing.B) {

	b.ResetTimer()

	var d date.Date
	buf := []byte(`"2018-01-31"`)
	for n := 0; n < b.N; n++ {
		if err := json.Unmarshal(buf, &d); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkDateUnmarshal(b *testing.B) {

	b.ResetTimer()

	var d date.Date
	buf := []byte(`"2011-01-31"`)
	for n := 0; n < b.N; n++ {
		if err := json.Unmarshal(buf, &d); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkTimeUnmarshal(b *testing.B) {

	b.ResetTimer()

	var t time.Time
	buf := []byte(`"2018-01-31T01:02:03Z"`)
	for n := 0; n < b.N; n++ {
		if err := json.Unmarshal(buf, &t); err != nil {
			b.Error(err)
		}
	}
}

func TestSchema(t *testing.T) {

	type filter struct {
		FromDt date.Date `schema:"from_dt"`
		ToDt   date.Date `schema:"to_dt"`
	}
	var f filter

	var data url.Values = map[string][]string{
		"from_dt": {"2019-03-01"},
		"to_dt":   {"2019-03-11"},
		"page":    {"0"},
		"rows":    {"10"},
	}

	decoder := schema.NewDecoder()
	//decoder.RegisterConverter(filter.FromDt, converter)
	//decoder := schema.NewDecoder()
	decoder.RegisterConverter(f.FromDt, date.Converter)
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(&f, data); err != nil {
		t.Error(err)
	}

	t.Logf("%x:%x", f.FromDt, f.ToDt)
}
~ ` `