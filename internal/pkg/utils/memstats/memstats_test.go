package memstats

import (
	"bytes"
	"strings"
	"testing"
)

func TestMemStats(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Ok",
			want: "us",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MemStats()
			h := humanizeNano(got.TotalAlloc)

			if !strings.Contains(h, tt.want) {
				t.Errorf("MemStats() = %v, want %v", h, tt.want)
			}
		})
	}
}

func TestPrintMemoryStats(t *testing.T) {
	tests := []struct {
		name  string
		wantW string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			PrintMemoryStats(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("PrintMemoryStats() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_humanizeNano(t *testing.T) {
	type args struct {
		n uint64
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Nano Second",
			args: args{
				n: 1,
			},
			want: "1ns",
		},
		{
			name: "Micro Second",
			args: args{
				n: 1_001,
			},
			want: "1us",
		},
		{
			name: "Milli Second",
			args: args{
				n: 1_000_001,
			},
			want: "1ms",
		},
		{
			name: "Second",
			args: args{
				n: 1_000_000_001,
			},
			want: "1s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := humanizeNano(tt.args.n); got != tt.want {
				t.Errorf("humanizeNano() = %v, want %v", got, tt.want)
			}
		})
	}
}
