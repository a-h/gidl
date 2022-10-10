package model_test

import (
	"encoding/json"
	"testing"

	"github.com/a-h/gidl/model"
	"github.com/a-h/gidl/model/tests/anonymous"
	"github.com/a-h/gidl/model/tests/chans"
	"github.com/a-h/gidl/model/tests/docs"
	"github.com/a-h/gidl/model/tests/enum"
	"github.com/a-h/gidl/model/tests/functions"
	"github.com/a-h/gidl/model/tests/functiontypes"
	"github.com/a-h/gidl/model/tests/pointers"
	"github.com/a-h/gidl/model/tests/privatetypes"
	"github.com/a-h/gidl/model/tests/publictypes"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		pkg      string
		expected string
	}{
		{
			name:     "private structs are ignored",
			pkg:      "github.com/a-h/gidl/model/tests/privatetypes",
			expected: privatetypes.Expected,
		},
		{
			name:     "public structs are included",
			pkg:      "github.com/a-h/gidl/model/tests/publictypes",
			expected: publictypes.Expected,
		},
		{
			name:     "string and integer enums are supported",
			pkg:      "github.com/a-h/gidl/model/tests/enum",
			expected: enum.Expected,
		},
		{
			name:     "pointers to pointers become a single pointer",
			pkg:      "github.com/a-h/gidl/model/tests/pointers",
			expected: pointers.Expected,
		},
		{
			name:     "functions and method receivers are ignored",
			pkg:      "github.com/a-h/gidl/model/tests/functions",
			expected: functions.Expected,
		},
		{
			name:     "fields of type channel are ignored",
			pkg:      "github.com/a-h/gidl/model/tests/chans",
			expected: chans.Expected,
		},
		{
			name:     "anonymous structs are ignored",
			pkg:      "github.com/a-h/gidl/model/tests/anonymous",
			expected: anonymous.Expected,
		},
		{
			name:     "function fields and function types are ignored",
			pkg:      "github.com/a-h/gidl/model/tests/functiontypes",
			expected: functiontypes.Expected,
		},
		{
			name:     "stuct, field and constant comments are extracted",
			pkg:      "github.com/a-h/gidl/model/tests/docs",
			expected: docs.Expected,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			m, err := model.Get(test.pkg)
			if err != nil {
				t.Fatalf("failed to get model %q: %v", test.pkg, err)
			}

			var expected *model.Model
			err = json.Unmarshal([]byte(test.expected), &expected)
			if err != nil {
				t.Fatalf("snapshot load failed: %v", err)
			}

			if diff := cmp.Diff(expected, m, cmpopts.IgnoreUnexported(*m)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
