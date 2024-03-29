package domain

import (
	"reflect"
	"testing"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name string
		c    string
		args []interface{}
		want Command
	}{
		{
			name: "Empty command",
			c:    "",
			args: []interface{}{},
			want: Command{},
		},
		{
			name: "Command with signle argument",
			c:    GET,
			args: []interface{}{"foo"},
			want: Command{
				Name: GET,
				Key:  "foo",
			},
		},
		{
			name: "Command with two arguments",
			c:    SET,
			args: []interface{}{"foo", "bar"},
			want: Command{
				Name:  SET,
				Key:   "foo",
				Value: "bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommand(tt.c, tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
