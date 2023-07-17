package domain

import (
	"keyvaluedb/storage"
	"reflect"
	"testing"
)

func TestKeyValueDBExecute(t *testing.T) {
	tests := []struct {
		name     string
		commands []Command
		expected []interface{}
	}{
		{
			name: "Set",
			commands: []Command{
				NewCommand(SET, "foo", "bar"),
			},
			expected: []interface{}{"OK"},
		},
		{
			name: "Set with invalid argument",
			commands: []Command{
				NewCommand(SET, "foo"),
			},
			expected: []interface{}{"(error) ERR wrong number of arguments for 'set' command"},
		},
		{
			name: "Get for nonexisting key",
			commands: []Command{
				NewCommand(GET, "nonexisting"),
			},
			expected: []interface{}{nil},
		},
		{
			name: "SetAndGet",
			commands: []Command{
				NewCommand(SET, "foo", "bar"),
				NewCommand(GET, "foo"),
			},
			expected: []interface{}{"OK", "bar"},
		},
		{
			name: "Get with invalid argument",
			commands: []Command{
				NewCommand(GET, ""),
			},
			expected: []interface{}{"(error) ERR wrong number of arguments for 'get' command"},
		},
		{
			name: "Delete for nonexisting key",
			commands: []Command{
				NewCommand(DEL, "nonexisting"),
			},
			expected: []interface{}{0},
		},
		{
			name: "SetAndDelete",
			commands: []Command{
				NewCommand(SET, "foo", "bar"),
				NewCommand(DEL, "foo"),
				NewCommand(GET, "foo"),
			},
			expected: []interface{}{"OK", 1, nil},
		},
		{
			name: "Increment for nonexisting key",
			commands: []Command{
				NewCommand(INCR, "counter"),
				NewCommand(GET, "counter"),
			},
			expected: []interface{}{"1", "1"},
		},
		{
			name: "Increment",
			commands: []Command{
				NewCommand(SET, "counter", "3"),
				NewCommand(INCR, "counter"),
				NewCommand(GET, "counter"),
			},
			expected: []interface{}{"OK", "4", "4"},
		},
		{
			name: "Increment non-integer value",
			commands: []Command{
				NewCommand(SET, "counter", "non-integer"),
				NewCommand(INCR, "counter"),
			},
			expected: []interface{}{"OK", "(error) ERR value is not an integer or out of range"},
		},
		{
			name: "IncrementBy for nonexisting key",
			commands: []Command{
				NewCommand(INCRBY, "counter", "10"),
				NewCommand(GET, "counter"),
			},
			expected: []interface{}{"10", "10"},
		},
		{
			name: "IncrementBy with invalid value",
			commands: []Command{
				NewCommand(INCRBY, "counter", nil),
				NewCommand(GET, "counter"),
			},
			expected: []interface{}{"(error) ERR wrong number of arguments for 'incrby' command", nil},
		},
		{
			name: "IncrementBy",
			commands: []Command{
				NewCommand(SET, "counter", "3"),
				NewCommand(INCRBY, "counter", "10"),
				NewCommand(GET, "counter"),
			},
			expected: []interface{}{"OK", "13", "13"},
		},
		{
			name: "IncrementBy integer value for exiting non-integer value",
			commands: []Command{
				NewCommand(SET, "counter", "non-integer"),
				NewCommand(INCRBY, "counter", "10"),
			},
			expected: []interface{}{"OK", "(error) ERR value is not an integer or out of range"},
		},
		{
			name: "IncrementBy non-integer value for exiting integer value",
			commands: []Command{
				NewCommand(SET, "counter", "10"),
				NewCommand(INCRBY, "counter", "non-integer"),
			},
			expected: []interface{}{"OK", "(error) ERR value is not an integer or out of range"},
		},
		{
			name: "SetAndGetMultipleKeys",
			commands: []Command{
				NewCommand(SET, "foo", "bar"),
				NewCommand(SET, "baz", "qux"),
				NewCommand(GET, "foo"),
				NewCommand(GET, "baz"),
			},
			expected: []interface{}{"OK", "OK", "bar", "qux"},
		},
		{
			name: "MultiBlock",
			commands: []Command{
				NewCommand(MULTI),
				NewCommand(SET, "foo", "bar"),
				NewCommand(GET, "foo"),
				NewCommand(EXEC),
			},
			expected: []interface{}{"OK", "QUEUED", "QUEUED", []interface{}{"OK", "bar"}},
		},
		{
			name: "MultiBlock with error in one of the command in the transaction",
			commands: []Command{
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
			commands: []Command{
				NewCommand(MULTI),
				NewCommand(SET, "foo", "bar"),
				NewCommand(DISCARD),
				NewCommand(GET, "foo"),
			},
			expected: []interface{}{"OK", "QUEUED", "OK", nil},
		},
		{
			name: "Compact",
			commands: []Command{
				NewCommand(SET, "foo", "bar"),
				NewCommand(SET, "baz", "qux"),
				NewCommand(COMPACT),
			},
			expected: []interface{}{"OK", "OK", []interface{}{"SET foo bar", "SET baz qux"}},
		},
		{
			name: "Select with wrong number of arguments",
			commands: []Command{
				NewCommand(SELECT),
			},
			expected: []interface{}{"(error) ERR wrong number of arguments for 'select' command"},
		},
		{
			name: "Select with invalid database index",
			commands: []Command{
				NewCommand(SELECT, "invalid"),
			},
			expected: []interface{}{"(error) ERR value is not an integer or out of range"},
		},
		{
			name: "Select with out of range database index",
			commands: []Command{
				NewCommand(SELECT, "40"),
			},
			expected: []interface{}{"(error) ERR DB index is out of range"},
		},
		{
			name: "Select valid database",
			commands: []Command{
				NewCommand(SELECT, "1"),
			},
			expected: []interface{}{"OK"},
		},
		{
			name: "Invalid command",
			commands: []Command{
				NewCommand("INVALID", "foo", "bar"),
			},
			expected: []interface{}{"(error) ERR unknown command `INVALID`, with args beginning with: `foo`, `bar`,"},
		},
	}

	for _, test := range tests {
		kvdb := NewKeyValueDB(storage.NewInMemory("2"))
		t.Run(test.name, func(t *testing.T) {
			for idx, cmd := range test.commands {
				_, got := kvdb.Execute(0, cmd)
				want := test.expected[idx]
				if !reflect.DeepEqual(got, want) {
					t.Errorf("command %v returned %q, expected %q", cmd, got, want)
				}
			}
		})
	}
}
