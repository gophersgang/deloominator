package db_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/deluminator/pkg/db"
)

func TestNewDSN(t *testing.T) {
	tests := []struct {
		name, source string
		expected     *db.DSN
	}{
		{
			"postgres full",
			"postgresql://grace:hopper@hal9000:3000/test?foo=bar",
			&db.DSN{
				Driver:   "postgresql",
				Username: "grace",
				Pass:     "hopper",
				Host:     "hal9000",
				Port:     3000,
				DBName:   "test",
				Options:  "foo=bar",
			},
		},
		{
			"postgres no port",
			"postgresql://grace:hopper@hal9000/test?foo=bar",
			&db.DSN{
				Driver:   "postgresql",
				Username: "grace",
				Pass:     "hopper",
				Host:     "hal9000",
				DBName:   "test",
				Options:  "foo=bar",
			},
		},
		{
			"postgres no credentials",
			"postgresql://hal9000:3000/test?foo=bar",
			&db.DSN{
				Driver:  "postgresql",
				Host:    "hal9000",
				Port:    3000,
				DBName:  "test",
				Options: "foo=bar",
			},
		},
		{
			"postgres no credentials no database",
			"postgresql://hal9000:3000/?foo=bar",
			&db.DSN{
				Driver:  "postgresql",
				Host:    "hal9000",
				Port:    3000,
				Options: "foo=bar",
			},
		},
		{
			"mysql full",
			"mysql://grace:hopper@hal9000:3000/test?foo=bar",
			&db.DSN{
				Driver:   "mysql",
				Username: "grace",
				Pass:     "hopper",
				Host:     "hal9000",
				Port:     3000,
				DBName:   "test",
				Options:  "foo=bar",
			},
		},
		{
			"mysql no port",
			"mysql://grace:hopper@hal9000/test?foo=bar",
			&db.DSN{
				Driver:   "mysql",
				Username: "grace",
				Pass:     "hopper",
				Host:     "hal9000",
				DBName:   "test",
				Options:  "foo=bar",
			},
		},
		{
			"mysql no credentials",
			"mysql://hal9000:3000/test?foo=bar",
			&db.DSN{
				Driver:  "mysql",
				Host:    "hal9000",
				Port:    3000,
				DBName:  "test",
				Options: "foo=bar",
			},
		},
		{
			"mysql no credentials no database",
			"mysql://hal9000:3000/?foo=bar",
			&db.DSN{
				Driver:  "mysql",
				Host:    "hal9000",
				Port:    3000,
				Options: "foo=bar",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := db.NewDSN(test.source)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("Expected %v to equal %v", actual, test.expected)
			}
		})
	}
}

func TestNewDSNWrongInput(t *testing.T) {
	tests := []struct {
		dsn      string
		expected string
	}{
		{dsn: "the answer is forthy two", expected: "the answer is forthy two is not a valid Data Source Name"},
		{dsn: "postgres://localhost:fortytwo/anwswer", expected: "postgres://localhost:fortytwo/anwswer is not a valid Data Source Name. fortytwo is not a valid port"},
	}

	for _, test := range tests {
		t.Run(test.dsn, func(t *testing.T) {
			_, err := db.NewDSN(test.dsn)
			if err == nil {
				t.Fatalf("Expected error but got none")
			}

			if !reflect.DeepEqual(test.expected, err.Error()) {
				t.Fatalf("Expected %s, got %s", test.expected, err.Error())
			}
		})
	}
}

func TestDSNFormat(t *testing.T) {
	tests := []struct {
		name     string
		source   *db.DSN
		expected string
	}{
		{
			"postgres full",
			&db.DSN{
				Driver:   "postgres",
				Username: "grace",
				Pass:     "hopper",
				Host:     "hal9000",
				Port:     3000,
				DBName:   "test",
				Options:  "foo=bar",
			},
			"postgres://grace:hopper@hal9000:3000/test?foo=bar",
		},
		{
			"postgres no port",
			&db.DSN{
				Driver:   "postgres",
				Username: "grace",
				Pass:     "hopper",
				Host:     "hal9000",
				DBName:   "test",
				Options:  "foo=bar",
			},
			"postgres://grace:hopper@hal9000/test?foo=bar",
		},
		{
			"postgres no credentials",
			&db.DSN{
				Driver:  "postgres",
				Host:    "hal9000",
				Port:    3000,
				DBName:  "test",
				Options: "foo=bar",
			},
			"postgres://hal9000:3000/test?foo=bar",
		},
		{
			"postgres no options",
			&db.DSN{
				Driver:   "postgres",
				Username: "grace",
				Pass:     "hopper",
				Host:     "hal9000",
				Port:     3000,
				DBName:   "test",
			},
			"postgres://grace:hopper@hal9000:3000/test",
		},
		{
			"mysql full",
			&db.DSN{
				Driver:   "mysql",
				Username: "grace",
				Pass:     "hopper",
				Host:     "hal9000",
				Port:     3000,
				DBName:   "test",
				Options:  "foo=bar",
			},
			"grace:hopper@hal9000:3000/test?foo=bar",
		},
		{
			"mysql no pass",
			&db.DSN{
				Driver:   "mysql",
				Username: "grace",
				Host:     "hal9000",
				Port:     3000,
				DBName:   "test",
				Options:  "foo=bar",
			},
			"grace@hal9000:3000/test?foo=bar",
		},
		{
			"mysql no database",
			&db.DSN{
				Driver:   "mysql",
				Username: "grace",
				Host:     "hal9000",
				Port:     3000,
				Options:  "foo=bar",
			},
			"grace@hal9000:3000/?foo=bar",
		},
		{
			"mysql no address",
			&db.DSN{
				Driver:   "mysql",
				Username: "grace",
				Pass:     "hopper",
				DBName:   "test",
				Options:  "foo=bar",
			},
			"grace:hopper@/test?foo=bar",
		},
		{
			"mysql user pass and database name",
			&db.DSN{
				Driver:   "mysql",
				Username: "grace",
				Pass:     "hopper",
				DBName:   "test_123123123",
			},
			"grace:hopper@/test_123123123",
		},
		{
			"mysql nothing",
			&db.DSN{
				Driver:   "mysql",
				Username: "grace",
			},
			"grace@/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.source.Format()

			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("Expected %v to equal %v", actual, test.expected)
			}
		})
	}
}