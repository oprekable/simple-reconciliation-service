package process

import (
	"reflect"
	"simple-reconciliation-service/internal/app/component"
	"simple-reconciliation-service/internal/app/component/cconfig"
	"simple-reconciliation-service/internal/app/component/cerror"
	"simple-reconciliation-service/internal/app/component/cfs"
	"simple-reconciliation-service/internal/app/component/clogger"
	"simple-reconciliation-service/internal/app/component/csqlite"
	"simple-reconciliation-service/internal/app/repository"
	mockprocess "simple-reconciliation-service/internal/app/repository/process/_mock"
	mocksample "simple-reconciliation-service/internal/app/repository/sample/_mock"
	"testing"
)

func TestProviderSvc(t *testing.T) {
	type args struct {
		comp *component.Components
		repo *repository.Repositories
	}

	tests := []struct {
		args args
		want *Svc
		name string
	}{
		{
			name: "Ok",
			args: args{
				comp: component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				),
				repo: repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			},
			want: ProviderSvc(
				component.NewComponents(
					&cconfig.Config{},
					&clogger.Logger{},
					&cerror.Error{},
					&csqlite.DBSqlite{},
					&cfs.Fs{},
				), repository.NewRepositories(
					mocksample.NewRepository(t),
					mockprocess.NewRepository(t),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProviderSvc(tt.args.comp, tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProviderSvc() = %v, want %v", got, tt.want)
			}
		})
	}
}
