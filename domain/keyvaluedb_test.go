package domain

import (
	"reflect"
	"testing"
)

func TestKeyValueDBExecute(t *testing.T) {
	tests := []struct {
		name     string
		commands []command
		expected []interface{}
	}{
		{
			name: "Set",
			commands: []command{
				NewCommand(SET, "foo", "bar"),
			},
			expected: []interface{}{"OK"},
		},
		{
			name: "Set with invalid argument",
			commands: []command{
				NewCommand(SET, "foo"),
			},
			expected: []interface{}{"(error) ERR wrong number of arguments for 'set' command"},
		},
		{
			name: "Get for nonexisting key",
			commands: []command{
				NewCommand(GET, "nonexisting"),
			},
			expected: []interface{}{nil},
		},
		{
			name: "SetAndGet",
			commands: []command{
				NewCommand(SET, "foo", "bar"),
				NewCommand(GET, "foo"),
			},
			expected: []interface{}{"OK", "bar"},
		},
		{
			name: "Get with invalid argument",
			commands: []command{
				NewCommand(GET, ""),
			},
			expected: []interface{}{"(error) ERR wrong number of arguments for 'get' command"},
		},
		{
			name: "Delete for nonexisting key",
			commands: []command{
				NewCommand(DEL, "nonexisting"),
			},
			expected: []interface{}{0},
		},
		{
			name: "SetAndDelete",
			commands: []command{
				NewCommand(SET, "foo", "bar"),
				NewCommand(DEL, "foo"),
				NewCommand(GET, "foo"),
			},
			expected: []interface{}{"OK", 1, nil},
		},
		{
			name: "Increment for nonexisting key",
			commands: []command{
				NewCommand(INCR, "counter"),
				NewCommand(GET, "counter"),
			},
			expected: []interface{}{"1", "1"},
		},
		{
			name: "Increment",
			commands: []command{
				NewCommand(SET, "counter", "3"),
				NewCommand(INCR, "counter"),
				NewCommand(GET, "counter"),
			},
			expected: []interface{}{"OK", "4", "4"},
		},
		{
			name: "Increment non-integer value",
			commands: []command{
				NewCommand(SET, "counter", "non-integer"),
				NewCommand(INCR, "counter"),
			},
			expected: []interface{}{"OK", "(error) ERR value is not an integer or out of range"},
		},
		{
			name: "IncrementBy for nonexisting key",
			commands: []command{
				NewCommand(INCRBY, "counter", "10"),
				NewCommand(GET, "counter"),
			},
			expected: []interface{}{"10", "10"},
		},
		{
			name: "IncrementBy with invalid value",
			commands: []command{
				NewCommand(INCRBY, "counter", nil),
				NewCommand(GET, "counter"),
			},
			expected: []interface{}{"(error) ERR wrong number of arguments for 'incrby' command", nil},
		},
		{
			name: "IncrementBy",
			commands: []command{
				NewCommand(SET, "counter", "3"),
				NewCommand(INCRBY, "counter", "10"),
				NewCommand(GET, "counter"),
			},
			expected: []interface{}{"OK", "13", "13"},
		},
		{
			name: "IncrementBy integer value for exiting non-integer value",
			commands: []command{
				NewCommand(SET, "counter", "non-integer"),
				NewCommand(INCRBY, "counter", "10"),
			},
			expected: []interface{}{"OK", "(error) ERR value is not an integer or out of range"},
		},
		{
			name: "IncrementBy non-integer value for exiting integer value",
			commands: []command{
				NewCommand(SET, "counter", "10"),
				NewCommand(INCRBY, "counter", "non-integer"),
			},
			expected: []interface{}{"OK", "(error) ERR value is not an integer or out of range"},
		},
		{
			name: "SetAndGetMultipleKeys",
			commands: []command{
				NewCommand(SET, "foo", "bar"),
				NewCommand(SET, "baz", "qux"),
				NewCommand(GET, "foo"),
				NewCommand(GET, "baz"),
			},
			expected: []interface{}{"OK", "OK", "bar", "qux"},
		},
		{
			name: "MultiBlock",
			commands: []command{
				NewCommand(MULTI),
				NewCommand(SET, "foo", "bar"),
				NewCommand(GET, "foo"),
				NewCommand(EXEC),
			},
			expected: []interface{}{"OK", "QUEUED", "QUEUED", []interface{}{"OK", "bar"}},
		},
		{
			name: "MultiBlock with error in one of the command in the transaction",
			commands: []command{
				NewCommand(MULTI),
				NewCommand(SET, "foo", "bar"),
				NewCommand(GET, "foo"),
				NewCommand(INCR, "foo"),
				NewCommand(EXEC),
			},
			expected: []interface{}{"OK", "QUEUED", "QUEUED", "QUEUED", []interface{}{"OK", "bar", "(error) ERR value is not an integer or out of range"}},
		},
		{
			name: "DiscardMultiBlock",
			commands: []command{
				NewCommand(MULTI),
				NewCommand(SET, "foo", "bar"),
				NewCommand(DISCARD),
				NewCommand(GET, "foo"),
			},
			expected: []interface{}{"OK", "QUEUED", "OK", nil},
		},
		{
			name: "Compact",
			commands: []command{
				NewCommand(SET, "foo", "bar"),
				NewCommand(SET, "baz", "qux"),
				NewCommand(COMPACT),
			},
			expected: []interface{}{"OK", "OK", []interface{}{"SET foo bar", "SET baz qux"}},
		},
		{
			name: "Invalid command",
			commands: []command{
				NewCommand("INVALID", "foo", "bar"),
			},
			expected: []interface{}{"(error) ERR unknown command `INVALID`, with args beginning with: 'foo', 'bar',"},
		},
	}

	for _, test := range tests {
		kvdb := NewKeyValueDB()
		t.Run(test.name, func(t *testing.T) {
			for idx, cmd := range test.commands {
				got := kvdb.Execute(cmd)
				want := test.expected[idx]
				if !reflect.DeepEqual(got, want) {
					t.Errorf("command %v returned %q, expected %q", cmd, got, want)
				}
			}
		})
	}
}
