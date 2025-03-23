package service

import (
	"reflect"
	"simple-reconciliation-service/internal/app/service/process"
	mockprocess "simple-reconciliation-service/internal/app/service/process/_mock"
	"simple-reconciliation-service/internal/app/service/sample"
	mocksample "simple-reconciliation-service/internal/app/service/sample/_mock"
	"testing"
)

func TestNewServices(t *testing.T) {
	type args struct {
		svcSample  sample.Service
		svcProcess process.Service
	}

	tests := []struct {
		args args
		want *Services
		name string
	}{
		{
			name: "Ok",
			args: args{
				svcSample:  mocksample.NewService(t),
				svcProcess: mockprocess.NewService(t),
			},
			want: &Services{
				SvcSample:  mocksample.NewService(t),
				SvcProcess: mockprocess.NewService(t),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServices(tt.args.svcSample, tt.args.svcProcess); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServices() = %v, want %v", got, tt.want)
			}
		})
	}
}
