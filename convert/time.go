package convert

import "time"

// DateTimeLayout 是常用日期时间布局：yyyy-MM-dd HH:mm:ss。
const DateTimeLayout = "2006-01-02 15:04:05"

// DateLayout 是常用日期布局：yyyy-MM-dd。
const DateLayout = "2006-01-02"

// TimeLayout 是常用时间布局：HH:mm:ss。
const TimeLayout = "15:04:05"

// TimeToString 按指定布局格式化时间。
func TimeToString(t time.Time, layout string) string {
	return t.Format(layout)
}

// TimeToStringInLocation 将时间转换到指定时区后，按指定布局格式化。
func TimeToStringInLocation(t time.Time, layout string, loc *time.Location) string {
	return t.In(locationOrLocal(loc)).Format(layout)
}

// TimeToDateTimeString 使用 DateTimeLayout 格式化时间。
func TimeToDateTimeString(t time.Time) string {
	return t.Format(DateTimeLayout)
}

// TimeToDateTimeStringInLocation 将时间转换到指定时区后，使用 DateTimeLayout 格式化。
func TimeToDateTimeStringInLocation(t time.Time, loc *time.Location) string {
	return TimeToStringInLocation(t, DateTimeLayout, loc)
}

// TimeToDateString 使用 DateLayout 格式化时间。
func TimeToDateString(t time.Time) string {
	return t.Format(DateLayout)
}

// TimeToDateStringInLocation 将时间转换到指定时区后，使用 DateLayout 格式化。
func TimeToDateStringInLocation(t time.Time, loc *time.Location) string {
	return TimeToStringInLocation(t, DateLayout, loc)
}

// TimeToUnix 将时间转换为秒级 Unix 时间戳。
func TimeToUnix(t time.Time) int64 {
	return t.Unix()
}

// TimeToUnixMilli 将时间转换为毫秒级 Unix 时间戳。
func TimeToUnixMilli(t time.Time) int64 {
	return t.UnixMilli()
}

// UnixToTime 将秒级 Unix 时间戳转换为时间。
func UnixToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}

// UnixToTimeInLocation 将秒级 Unix 时间戳转换为指定时区的时间。
func UnixToTimeInLocation(sec int64, loc *time.Location) time.Time {
	return time.Unix(sec, 0).In(locationOrLocal(loc))
}

// UnixMilliToTime 将毫秒级 Unix 时间戳转换为时间。
func UnixMilliToTime(msec int64) time.Time {
	return time.UnixMilli(msec)
}

// UnixMilliToTimeInLocation 将毫秒级 Unix 时间戳转换为指定时区的时间。
func UnixMilliToTimeInLocation(msec int64, loc *time.Location) time.Time {
	return time.UnixMilli(msec).In(locationOrLocal(loc))
}

func locationOrLocal(loc *time.Location) *time.Location {
	if loc == nil {
		return time.Local
	}
	return loc
}
