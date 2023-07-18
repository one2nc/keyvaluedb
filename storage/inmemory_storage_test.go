package storage

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewInMemory(t *testing.T) {
	tests := []struct {
		name     string
		dbCntStr string
		want     Storage
	}{
		{
			name:     "Empty dbCountStr should default to 16",
			dbCntStr: "",
			want:     NewInMemory("16"),
		},
		{
			name:     "Valid dbCountStr",
			dbCntStr: "10",
			want:     NewInMemory("10"),
		},
		{
			name:     "Invalid dbCountStr should default to 16",
			dbCntStr: "abc",
			want:     NewInMemory("16"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewInMemory(tt.dbCntStr)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInMemory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemorySelect(t *testing.T) {
	tests := []struct {
		name       string
		dbCntStr   string
		dbIndexStr string
		want       int
		wantErr    error
	}{
		{
			name:       "Valid dbIndexStr within range",
			dbCntStr:   "16",
			dbIndexStr: "0",
			want:       0,
			wantErr:    nil,
		},
		{
			name:       "Valid dbIndexStr within range",
			dbCntStr:   "16",
			dbIndexStr: "10",
			want:       10,
			wantErr:    nil,
		},
		{
			name:       "Invalid dbIndexStr (out of range)",
			dbCntStr:   "16",
			dbIndexStr: "-1",
			want:       0,
			wantErr:    fmt.Errorf("(error) ERR DB index is out of range"),
		},
		{
			name:       "Invalid dbIndexStr (not an integer)",
			dbCntStr:   "16",
			dbIndexStr: "abc",
			want:       0,
			wantErr:    fmt.Errorf("(error) ERR value is not an integer or out of range"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := NewInMemory(tt.dbCntStr)

			got, err := in.Select(tt.dbIndexStr)
			if err != nil {
				if tt.wantErr == nil {
					t.Errorf("inMemory.Select() error = %v, wantErr <nil>", err)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("inMemory.Select() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if tt.wantErr != nil {
				t.Errorf("inMemory.Select() error = <nil>, wantErr %v", tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("inMemory.Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemorySetGet(t *testing.T) {
	type args struct {
		dbIndex int
		key     string
		value   interface{}
	}
	tests := []struct {
		name     string
		dbCntStr string
		key      string
		setArgs  args
		want     interface{}
	}{
		{
			name: "Set a new key-value pair",
			setArgs: args{
				dbIndex: 0,
				key:     "key1",
				value:   "value1",
			},
			key:  "key1",
			want: "value1",
		},
		{
			name: "Update the value of an existing key",
			setArgs: args{
				dbIndex: 0,
				key:     "key1",
				value:   "value2",
			},
			key:  "key1",
			want: "value2",
		},
		{
			name: "Get a value for nonexisting key",
			key:  "nonexisting",
			want: nil,
		},
		{
			name: "Get the value of an existing key",
			setArgs: args{
				dbIndex: 0,
				key:     "key1",
				value:   "value2",
			},
			key:  "key1",
			want: "value2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := NewInMemory(tt.dbCntStr)

			in.Set(tt.setArgs.dbIndex, tt.setArgs.key, tt.setArgs.value)

			got := in.Get(tt.setArgs.dbIndex, tt.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inMemory.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemorySetDel(t *testing.T) {
	type args struct {
		dbIndex int
		key     string
		value   interface{}
	}
	tests := []struct {
		name     string
		dbCntStr string
		key      string
		setArgs  args
		want     interface{}
	}{
		{
			name: "Del a value for nonexisting key",
			key:  "nonexisting",
			want: 0,
		},
		{
			name: "Del the value of an existing key",
			setArgs: args{
				dbIndex: 0,
				key:     "key1",
				value:   "value2",
			},
			key:  "key1",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := NewInMemory(tt.dbCntStr)

			in.Set(tt.setArgs.dbIndex, tt.setArgs.key, tt.setArgs.value)

			got := in.Del(tt.setArgs.dbIndex, tt.key)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inMemory.Set() = %v, want %v", got, tt.want)
			}

			if in.Get(tt.setArgs.dbIndex, tt.key) != nil {
				t.Errorf("Del(%d, %s) did not delete the key properly", tt.setArgs.dbIndex, tt.key)
			}
		})
	}
}

func TestInMemoryGetAll(t *testing.T) {
	type fields struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name     string
		dbCntStr string
		dbIndex  int
		fields   []fields
		wantAll  []string
	}{
		{
			name:     "Get all key-value pairs",
			dbCntStr: "1",
			dbIndex:  0,
			fields: []fields{
				{
					key:   "key1",
					value: "value1",
				},
				{
					key:   "key2",
					value: "value2",
				},
				{
					key:   "key3",
					value: "value3",
				},
			},
			wantAll: []string{"key1 value1", "key2 value2", "key3 value3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := NewInMemory(tt.dbCntStr)

			for _, f := range tt.fields {
				in.Set(tt.dbIndex, f.key, f.value)
			}

			allChan := in.GetAll(tt.dbIndex)

			var gotAll []string
			for s := range allChan {
				gotAll = append(gotAll, s)
			}
			if len(gotAll) != len(tt.wantAll) {
				t.Errorf("GetAll(%d) returned %d items, want %d items", tt.dbIndex, len(gotAll), len(tt.wantAll))
			}
			for i, got := range gotAll {
				if got != tt.wantAll[i] {
					t.Errorf("GetAll(%d) returned %s, want %s", tt.dbIndex, got, tt.wantAll[i])
				}
			}
		})
	}
}
