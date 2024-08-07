package constants

const (
	ProductionMode  = "production"
	DevelopmentMode = "development"
	TestMode        = "test"

	DateTimeLayout    = "2006-01-02 15:04:05"
	DateTimeUtcLayout = "2006-01-02T15:04:05.000Z"
	DateTimeWithMS    = "2006-01-02 15:04:05.999999"
	DateTimeWithT     = "2006-01-02T15:04:05.999999"
	DatetimeMicro     = "060102150405.999999"
	DatetimeFormat1   = "20060102150405" // YYYYMMddHHmmss
	DateFormat1       = "20060102"       // YYYYMMdd
	DateLayout        = "2006-01-02"
	MonthLayout       = "2006-01"
	DateSimpleLayout  = "060102"
	DateMonthLayout   = "0102" // 月日

	RunEnv      = "RUN_ENV"      // 运行环境
	RunPlatform = "RUN_PLATFORM" // 运行平台
	RunApp      = "RUN_APP"      // 运行应用 sh bj

)
