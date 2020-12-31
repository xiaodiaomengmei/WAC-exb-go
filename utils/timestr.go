package utils

import (
	"fmt"
	"time"
	"strconv"
	//"strings"
)

//Str2Time:=Str2Time("2017-09-12 12:03:40")
//fmt.Println(Str2Time)
//2017-09-12 12:03:40 +0800 CST

//Str2Stamp:=Str2Stamp("2017-09-12 12:03:40")
//fmt.Println(Str2Stamp)
//1505189020000

//Time2Str:=Time2Str()
//fmt.Println(Time2Str)
//2018-11-11 17:50:50

//GetStamp:=Time2Stamp()
//fmt.Println(GetStamp)
//1541497850321

//Stamp2Str:=Stamp2Str(1505189020000)
//fmt.Println(Stamp2Str)
//2017-09-12 12:03:40

//Stamp2Time:=Stamp2Time(1505188820000)
//fmt.Println(Stamp2Time)
//2017-09-12 12:00:20 +0800 CST

/**字符串->时间对象*/
func Str2Time(formatTimeStr string) time.Time{
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型

	return theTime

}
/**字符串->时间戳*/
func Str2Stamp(formatTimeStr string) int64 {
	timeStruct:=Str2Time(formatTimeStr)
	millisecond:=timeStruct.UnixNano()/1e6
	return  millisecond
}

/**时间对象->字符串*/
func Time2Str() string {
	const shortForm = "2006-01-02 15:04:05"
	t := time.Now()
	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	return str
}


/*时间对象->时间戳*/
func Time2Stamp()int64{
	t:=time.Now()
	millisecond:=t.UnixNano()/1e6
	return  millisecond
}
/*时间戳->字符串*/
func Stamp2Str(stamp int64) string{
	timeLayout := "2006-01-02 15:04:05"
	str:=time.Unix(stamp/1000,0).Format(timeLayout)
	return str
}
/*时间戳->时间对象*/
func Stamp2Time(stamp int64)time.Time{
	stampStr:=Stamp2Str(stamp)
	timer:=Str2Time(stampStr)
	return timer
}
/**字符串2-字符串1->时间段字符串*/
func Str2Duration(formatT1Str string,formatT2Str string) string{
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	t1, _ := time.ParseInLocation(timeLayout, formatT1Str, loc) //使用模板在对应时区转化为time.time类型
	t2, _ := time.ParseInLocation(timeLayout, formatT2Str, loc) //使用模板在对应时区转化为time.time类型
	timesub := t2.Sub(t1)
	DurationStr := fmt.Sprintf("%s",timesub)
	return DurationStr
}


var weekday = [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
// get week date
//0--current week,1--before week
func GetWeekDayStr(WeekNum int) []string {
	var ArrayweekdayStr []string
	now := time.Now()
	//if weeknum==0
	WeekFirstDay :=  now.AddDate(0, 0, 0-int(now.Weekday()))
	RecycleMaxNum := int(now.Weekday())
	//else week num != 0
	if 0 != WeekNum {
		WeekFirstDay = now.AddDate(0, 0, -7*WeekNum-int(now.Weekday())) 
		RecycleMaxNum = 6
	}
	
	for i:=0; i<=RecycleMaxNum ; i++ {
		WeekDay := WeekFirstDay.AddDate(0, 0, i)
		WeekDayStr := fmt.Sprintf("%s_%s",weekday[int(WeekDay.Weekday())],WeekDay.Format("2006-01-02"))
		ArrayweekdayStr = append(ArrayweekdayStr, WeekDayStr)
	}
	return ArrayweekdayStr
}
func GetWeekDate(WeekNum int) []time.Time {
	var ArrayweekDate []time.Time
	now := time.Now()
	//if weeknum==0
	WeekFirstDay :=  now.AddDate(0, 0, 0-int(now.Weekday()))
	RecycleMaxNum := int(now.Weekday())
	//else week num != 0
	if 0 != WeekNum {
		WeekFirstDay = now.AddDate(0, 0, -7*WeekNum-int(now.Weekday())) 
		RecycleMaxNum = 6
	}
	for i:=0; i<=RecycleMaxNum ; i++ {
		WeekDay := WeekFirstDay.AddDate(0, 0, i)
		ArrayweekDate = append(ArrayweekDate, WeekDay)
	}
	return ArrayweekDate
}

func GetMonthDateStr(myYear int,myMonth string)  []string {
	var ArrayMonthdayStr []string
    // 数字月份必须前置补零
	if len(myMonth)==1 {
        myMonth = "0"+myMonth
	}
	yearStr:= strconv.Itoa(myYear)

    timeLayout := "2006-01-02 15:04:05"
    loc, _ := time.LoadLocation("Local")
    theTime, _ := time.ParseInLocation(timeLayout, yearStr+"-"+myMonth+"-01 00:00:00", loc)
    newMonth := theTime.Month()

    monthStartDate := time.Date(myYear,newMonth, 1, 0, 0, 0, 0, time.Local)
	//monthEndDate := time.Date(myYear,newMonth+1, 0, 0, 0, 0, 0, time.Local)
	monthdate:=monthStartDate
	for i:=1;i<=31;i++{
		ArrayMonthdayStr = append(ArrayMonthdayStr, fmt.Sprintf("%s_%s",weekday[int(monthdate.Weekday())],monthdate.Format("2006-01-02")))
		monthdate=monthdate.AddDate(0,0,1)
	}
    return ArrayMonthdayStr
}

func GetMonthDate(myYear int,myMonth string) []time.Time {
	var ArrayMonthday []time.Time
    // 数字月份必须前置补零
	if len(myMonth)==1 {
        myMonth = "0"+myMonth
	}
	yearStr:= strconv.Itoa(myYear)

    timeLayout := "2006-01-02 15:04:05"
    loc, _ := time.LoadLocation("Local")
    theTime, _ := time.ParseInLocation(timeLayout, yearStr+"-"+myMonth+"-01 00:00:00", loc)
    newMonth := theTime.Month()

    monthStartDate := time.Date(myYear,newMonth, 1, 0, 0, 0, 0, time.Local)
	//monthEndDate := time.Date(myYear,newMonth+1, 0, 0, 0, 0, 0, time.Local)
	monthdate:=monthStartDate
	for i:=1;i<=31;i++{
		ArrayMonthday = append(ArrayMonthday, monthdate)
		monthdate=monthdate.AddDate(0,0,1)
	}
    return ArrayMonthday
}