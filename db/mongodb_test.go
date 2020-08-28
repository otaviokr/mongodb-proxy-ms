package db_test

import (
	"strconv"
	"testing"

	"github.com/otaviokr/mongodb-proxy-ms/db"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	inputs   map[string]string
	expected string
	hasError bool
}

func TestNewConnection(t *testing.T) {

	testCases := []TestCase{
		{
			inputs:   map[string]string{"h": "cool_host", "p": "123456", "u": "cool_user", "pw": "supers3cr3t"},
			expected: "mongodb://cool_user:supers3cr3t@cool_host:123456",
			hasError: false,
		},
		{
			inputs:   map[string]string{"h": "cool_host", "p": "123456", "u": "", "pw": ""},
			expected: "mongodb://cool_host:123456",
			hasError: false,
		},
		{
			inputs:   map[string]string{"h": "cool_host", "p": "123456", "u": "", "pw": "supers3cr3t"},
			expected: "mongodb://cool_host:123456",
			hasError: false,
		},
		{
			inputs:   map[string]string{"h": "cool_host", "p": "123456", "u": "cool_user", "pw": ""},
			expected: "mongodb://cool_host:123456",
			hasError: false,
		},
		{
			inputs:   map[string]string{"h": "cool_host", "p": "-1", "u": "", "pw": ""},
			expected: "mongodb://cool_host",
			hasError: false,
		},
		{
			inputs:   map[string]string{"h": "cool_host", "p": "0", "u": "", "pw": ""},
			expected: "mongodb://cool_host",
			hasError: false,
		},
		{
			inputs:   map[string]string{"h": "cool_host", "p": "1", "u": "", "pw": ""},
			expected: "mongodb://cool_host:1",
			hasError: false,
		},
		{
			inputs:   map[string]string{"h": "cool_host", "p": "0", "u": "spec!@lCh@r%", "pw": "S:per/s3cr3t"},
			expected: "mongodb://spec!%40lCh%40r%25:S%3Aper%2Fs3cr3t@cool_host",
			hasError: false,
		},
	}

	for _, tc := range testCases {
		port, err := strconv.Atoi(tc.inputs["p"])
		if err != nil {
			t.FailNow()
		}

		actual, err := db.NewConnection(tc.inputs["h"], port, tc.inputs["u"], tc.inputs["pw"])
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, tc.expected, actual.URI, "URI is not the same")
	}
}
