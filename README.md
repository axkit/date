# Date Package

[![Build Status](https://github.com/axkit/date/actions/workflows/go.yml/badge.svg)](https://github.com/axkit/date/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axkit/date)](https://goreportcard.com/report/github.com/axkit/date)
[![GoDoc](https://pkg.go.dev/badge/github.com/axkit/date)](https://pkg.go.dev/github.com/axkit/date)
[![Coverage Status](https://coveralls.io/repos/github/axkit/date/badge.svg?branch=main)](https://coveralls.io/github/axkit/date?branch=main)


The `date` package provides a custom date type in Go, designed specifically for use cases where only the date (year, month, day) is needed, without any time or timezone information. This type replaces the standard `time.Time` type in scenarios where time and timezone components are unnecessary.

## Features

- Nullable date type, comparable and sortable.
- Stores date as a hex integer (0xYYYYMMDD), making it efficient in memory and database storage.
- Converts easily between `date.Date` and `time.Time`.
- Marshals and unmarshals to and from JSON.
- Compatible with SQL database operations.

## Installation

To install the package, use:

```bash
go get github.com/axkit/date
```

## Usage

Here’s how to use the `date` package:

### Import the Package

```go
import "github.com/axkit/date"
```

### Creating a New Date

You can create a new date using the `New` function, which accepts the year, month, and day as arguments:

```go
d := date.New(2023, time.January, 1)
fmt.Println(d.String()) // Output: 2023-01-01
```

### Working with Today’s Date

To get today’s date:

```go
today := date.Today()
fmt.Println(today.String()) // Output: current date in YYYY-MM-DD format
```

### Converting to time.Time

To convert `date.Date` to `time.Time`:

```go
t := d.Time() // Local time zone
utc := d.UTC() // UTC time zone
```

### Adding Years, Months, and Days

You can add time to a `date.Date`:

```go
d := date.New(2023, time.January, 1)
newDate := d.Add(1, 0, 0) // Adds 1 year
fmt.Println(newDate.String()) // Output: 2024-01-01
```

### Parsing a Date from String

To parse a date from a string:

```go
var d date.Date
err := d.Parse("2023-01-01")
if err != nil {
    fmt.Println("Error parsing date:", err)
}
fmt.Println(d.String()) // Output: 2023-01-01
```

### Checking Validity

You can check if a `date.Date` is null or valid:

```go
if d.Valid() {
    fmt.Println("Date is valid")
} else {
    fmt.Println("Date is null")
}
```

### Database Operations

The `date.Date` type is compatible with SQL databases, implementing both the `Scanner` and `Valuer` interfaces.

#### Storing to Database

```go
stmt, _ := db.Prepare("INSERT INTO dates (date) VALUES (?)")
_, err := stmt.Exec(d)
if err != nil {
    fmt.Println("Error inserting date:", err)
}
```

#### Retrieving from Database

```go
var d date.Date
err := db.QueryRow("SELECT date FROM dates WHERE id = ?", id).Scan(&d)
if err != nil {
    fmt.Println("Error scanning date:", err)
}
fmt.Println(d.String())
```

## JSON Marshaling and Unmarshaling

The `date.Date` type supports JSON encoding and decoding:

```go
type Example struct {
    Date date.Date `json:"date"`
}

e := Example{Date: date.New(2023, time.January, 1)}
jsonData, _ := json.Marshal(e)
fmt.Println(string(jsonData)) // Output: {"date":"2023-01-01"}
```

## Running Tests

To run tests for the `date` package:

```bash
go test github.com/axkit/date
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```

This README provides an overview of the package's features, how to install it, and how to use its primary functions.