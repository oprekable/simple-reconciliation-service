package csvhelper

import (
	"context"
	"path/filepath"

	"github.com/aaronjan/hunch"

	"github.com/jszwec/csvutil"

	"github.com/spf13/afero"
)

func StructToCSVFile(ctx context.Context, fs afero.Fs, filePath string, structData interface{}, isDeleteDirectory bool) error {
	_, err := hunch.Waterfall(
		ctx,
		func(c context.Context, i interface{}) (interface{}, error) {
			filePath = filepath.Clean(filePath)
			basePath := filepath.Dir(filePath)

			if isDeleteDirectory {
				return afero.Glob(fs, basePath+"/*")
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if isDeleteDirectory {
				contents := i.([]string)
				for _, item := range contents {
					e := fs.RemoveAll(item)
					if e != nil {
						return nil, e
					}
				}
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return nil, fs.MkdirAll(filepath.Dir(filePath), 0755)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return csvutil.Marshal(structData)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			if i != nil {
				marshal := i.([]byte)
				return nil, afero.WriteFile(fs, filePath, marshal, 0644)
			}

			return nil, nil
		},
	)

	return err
}

//
//func CSVFileToStruct[T any](ctx context.Context, fs afero.Fs, bank string, filePath string) (returnData T, err error) {
//	_, err := hunch.Waterfall(
//		ctx,
//		func(c context.Context, i interface{}) (interface{}, error) {
//			return csvutil.Marshal(structData)
//		},
//		func(c context.Context, i interface{}) (interface{}, error) {
//			if i != nil {
//				marshal := i.([]byte)
//				return nil, afero.WriteFile(fs, filePath, marshal, 0644)
//			}
//
//			return nil, nil
//		},
//	)
//
//	return err
//}
