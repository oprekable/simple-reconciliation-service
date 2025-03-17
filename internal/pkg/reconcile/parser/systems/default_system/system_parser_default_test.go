package default_system

import (
	"bytes"
	"context"
	"encoding/csv"
	"reflect"
	"simple-reconciliation-service/internal/pkg/reconcile/parser/systems"
	"testing"
	"time"
)

func TestCSVSystemTrxDataGetAmount(t *testing.T) {
	type fields struct {
		TrxID           string
		TransactionTime string
		Type            string
		Amount          float64
	}

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVSystemTrxData{
				TrxID:           tt.fields.TrxID,
				TransactionTime: tt.fields.TransactionTime,
				Type:            tt.fields.Type,
				Amount:          tt.fields.Amount,
			}
			if got := u.GetAmount(); got != tt.want {
				t.Errorf("GetAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVSystemTrxDataGetTransactionTime(t *testing.T) {
	type fields struct {
		TrxID           string
		TransactionTime string
		Type            string
		Amount          float64
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVSystemTrxData{
				TrxID:           tt.fields.TrxID,
				TransactionTime: tt.fields.TransactionTime,
				Type:            tt.fields.Type,
				Amount:          tt.fields.Amount,
			}
			if got := u.GetTransactionTime(); got != tt.want {
				t.Errorf("GetTransactionTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVSystemTrxDataGetTrxID(t *testing.T) {
	type fields struct {
		TrxID           string
		TransactionTime string
		Type            string
		Amount          float64
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVSystemTrxData{
				TrxID:           tt.fields.TrxID,
				TransactionTime: tt.fields.TransactionTime,
				Type:            tt.fields.Type,
				Amount:          tt.fields.Amount,
			}
			if got := u.GetTrxID(); got != tt.want {
				t.Errorf("GetTrxID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVSystemTrxDataGetType(t *testing.T) {
	type fields struct {
		TrxID           string
		TransactionTime string
		Type            string
		Amount          float64
	}

	tests := []struct {
		name   string
		fields fields
		want   systems.TrxType
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVSystemTrxData{
				TrxID:           tt.fields.TrxID,
				TransactionTime: tt.fields.TransactionTime,
				Type:            tt.fields.Type,
				Amount:          tt.fields.Amount,
			}
			if got := u.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVSystemTrxDataToSystemTrxData(t *testing.T) {
	type fields struct {
		TrxID           string
		TransactionTime string
		Type            string
		Amount          float64
	}

	tests := []struct {
		name           string
		fields         fields
		wantReturnData *systems.SystemTrxData
		wantErr        bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CSVSystemTrxData{
				TrxID:           tt.fields.TrxID,
				TransactionTime: tt.fields.TransactionTime,
				Type:            tt.fields.Type,
				Amount:          tt.fields.Amount,
			}
			gotReturnData, err := u.ToSystemTrxData()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToSystemTrxData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("ToSystemTrxData() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}

func TestNewSystemParser(t *testing.T) {
	type args struct {
		dataStruct   systems.SystemTrxDataInterface
		csvReader    *csv.Reader
		isHaveHeader bool
	}

	tests := []struct {
		name    string
		args    args
		want    *SystemParser
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				dataStruct:   &CSVSystemTrxData{},
				csvReader:    csv.NewReader(nil),
				isHaveHeader: true,
			},
			want: &SystemParser{
				dataStruct:   &CSVSystemTrxData{},
				csvReader:    csv.NewReader(nil),
				parser:       systems.DefaultSystemParser,
				isHaveHeader: true,
			},
			wantErr: false,
		},
		{
			name: "Error csvReader is nil",
			args: args{
				dataStruct:   &CSVSystemTrxData{},
				csvReader:    nil,
				isHaveHeader: true,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSystemParser(tt.args.dataStruct, tt.args.csvReader, tt.args.isHaveHeader)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewSystemParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSystemParser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemParserToSystemTrxData(t *testing.T) {
	layoutTime := "2006-01-02 15:04:05"

	type fields struct {
		dataStruct   systems.SystemTrxDataInterface
		csvReader    *csv.Reader
		parser       systems.SystemParserType
		isHaveHeader bool
	}

	type args struct {
		ctx      context.Context
		filePath string
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantReturnData []*systems.SystemTrxData
		wantErr        bool
	}{
		{
			name: "Ok with header",
			fields: fields{
				dataStruct: &CSVSystemTrxData{},
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						`TrxID,TransactionTime,Type,Amount
0012d068c53eb0971fc8563343c5d81f,2025-03-15 10:51:52,CREDIT,20500
005dcbc9e27365a072be5393ea8d0f37,2025-03-14 18:29:01,CREDIT,42100`,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: true,
			},
			args: args{
				ctx:      context.Background(),
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: []*systems.SystemTrxData{
				{
					TrxID: "0012d068c53eb0971fc8563343c5d81f",
					TransactionTime: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-15 10:51:52")
						return t
					}(),
					Type:     "CREDIT",
					FilePath: "/tmp/foo.csv",
					Amount:   20500,
				},
				{
					TrxID: "005dcbc9e27365a072be5393ea8d0f37",
					TransactionTime: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-14 18:29:01")
						return t
					}(),
					Type:     "CREDIT",
					FilePath: "/tmp/foo.csv",
					Amount:   42100,
				},
			},
			wantErr: false,
		},
		{
			name: "Error decode with header",
			fields: fields{
				dataStruct: &CSVSystemTrxData{},
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						``,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: true,
			},
			args: args{
				ctx:      context.Background(),
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: nil,
			wantErr:        true,
		},
		{
			name: "Error parse TransactionTime",
			fields: fields{
				dataStruct: &CSVSystemTrxData{},
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						`0012d068c53eb0971fc8563343c5d81f,random string,CREDIT,20500`,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: false,
			},
			args: args{
				ctx:      context.Background(),
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: nil,
			wantErr:        false,
		},
		{
			name: "Ok without header",
			fields: fields{
				dataStruct: &CSVSystemTrxData{},
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						`0012d068c53eb0971fc8563343c5d81f,2025-03-15 10:51:52,CREDIT,20500
005dcbc9e27365a072be5393ea8d0f37,2025-03-14 18:29:01,CREDIT,42100`,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: false,
			},
			args: args{
				ctx:      context.Background(),
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: []*systems.SystemTrxData{
				{
					TrxID: "0012d068c53eb0971fc8563343c5d81f",
					TransactionTime: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-15 10:51:52")
						return t
					}(),
					Type:     "CREDIT",
					FilePath: "/tmp/foo.csv",
					Amount:   20500,
				},
				{
					TrxID: "005dcbc9e27365a072be5393ea8d0f37",
					TransactionTime: func() time.Time {
						t, _ := time.Parse(layoutTime, "2025-03-14 18:29:01")
						return t
					}(),
					Type:     "CREDIT",
					FilePath: "/tmp/foo.csv",
					Amount:   42100,
				},
			},
			wantErr: false,
		},
		{
			name: "Error nil csvReader without header",
			fields: fields{
				dataStruct:   &CSVSystemTrxData{},
				csvReader:    nil,
				parser:       "",
				isHaveHeader: false,
			},
			args: args{
				ctx:      context.Background(),
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: nil,
			wantErr:        false,
		},
		{
			name: "Error nil dataStruct without header",
			fields: fields{
				dataStruct: nil,
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						``,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: false,
			},
			args: args{
				ctx:      context.Background(),
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: nil,
			wantErr:        true,
		},
		{
			name: "Error decode without header",
			fields: fields{
				dataStruct: &CSVSystemTrxData{},
				csvReader: func() *csv.Reader {
					f := bytes.NewBufferString(
						``,
					)
					return csv.NewReader(f)
				}(),
				parser:       "",
				isHaveHeader: false,
			},
			args: args{
				ctx:      context.Background(),
				filePath: "/tmp/foo.csv",
			},
			wantReturnData: nil,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &SystemParser{
				dataStruct:   tt.fields.dataStruct,
				csvReader:    tt.fields.csvReader,
				parser:       tt.fields.parser,
				isHaveHeader: tt.fields.isHaveHeader,
			}

			gotReturnData, err := d.ToSystemTrxData(tt.args.ctx, tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToSystemTrxData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotReturnData, tt.wantReturnData) {
				t.Errorf("ToSystemTrxData() gotReturnData = %v, want %v", gotReturnData, tt.wantReturnData)
			}
		})
	}
}
