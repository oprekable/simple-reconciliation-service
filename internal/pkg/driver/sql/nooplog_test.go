package sql

import (
	"context"
	sqldblogger "github.com/simukti/sqldb-logger"
	"reflect"
	"testing"
)

func TestNewNoopLog(t *testing.T) {
	tests := []struct {
		name string
		want sqldblogger.Logger
	}{
		{
			name: "Ok",
			want: &noopLogAdapter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNoopLog(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNoopLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoopLogAdapter_Log(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 sqldblogger.Level
		in2 string
		in3 map[string]interface{}
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Ok",
			args: args{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zl := &noopLogAdapter{}
			zl.Log(tt.args.in0, tt.args.in1, tt.args.in2, tt.args.in3)
		})
	}
}
