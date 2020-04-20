package core

const (
	SUCCESS                    int64 = 1
	DB_CONNECTTION_FAILED    int64 = 10000
	DB_LOADDATA_FAILED        int64 = 10001
	REST_PARAM_INVALID        int64 = 20000
	REST_METHOD_INVALID       int64 = 20001
	REST_ILLEGAL_DATAFORMAT   int64 = 20002
	)

var ErrMap = map[int64]string{
	SUCCESS:                  "success",
	DB_CONNECTTION_FAILED:  "connect db error",
	DB_LOADDATA_FAILED:      "load db data error",
	REST_PARAM_INVALID:      "invalid rest parameter",
	REST_METHOD_INVALID:     "invalid rest method",
	REST_ILLEGAL_DATAFORMAT: "rest illegal data format",
}
