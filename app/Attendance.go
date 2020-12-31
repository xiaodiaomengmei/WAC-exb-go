//application package
package app
//attendance sub application, use user history trace table,
//generate attendance by user_name and date,then export xls, send email to manager
import (
	"fmt"
	//"log"
	"time"
	"strings"
	//"github.com/Luxurioust/excelize"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"wifidog-server/model"
	"wifidog-server/dao"
	"wifidog-server/utils"
	)
type AttendUser struct{
	AccountName  string 
	HunmanName   string 
	Date         string
	LoginTime    string 
	LogoutTime   string 
	Duration     float64 
}
var userDao=new(dao.UserDao)
var userhistorytraceDao=new(dao.UserHistoryTraceDao)
var usernowtraceDao=new(dao.UserNowTraceDao)
//获取总工时：下班时间-上班时间
func GetAttendanByUsernameAndDate(user_atten *AttendUser) {
	var userhistorytraces []model.UserHistoryTrace
	userhistorytraces = userhistorytraceDao.GetUserHistoryTraceByUserNameAndLoginDate(user_atten.AccountName ,user_atten.Date)

	var logintime string
	var logouttime string
	logintime = "Z"    //set max
	logouttime = "0"   //set min
	for _,v := range userhistorytraces{
		if strings.Compare(logintime, v.Logintime) > 0{
			logintime = v.Logintime
		}
		if strings.Compare(logouttime, v.Logouttime) < 0{
			logouttime=v.Logouttime
		}
	}
	t1 := utils.Str2Time(logintime)
	user_atten.LoginTime = t1.Format("15:04:05")
	t1 = utils.Str2Time(logouttime)
	user_atten.LogoutTime = t1.Format("15:04:05")
	
	time_part,_:= time.ParseDuration(utils.Str2Duration(logintime, logouttime))
	user_atten.Duration = time_part.Hours()
}
//每天午夜12点,将当前轨迹表导入到用户历史轨迹，然后删除所有的用户当前轨迹
//AP ipset list auto clear at 24'clock,too 
func DelAllUserNowTrace(){
	//update user history trace table
	var usernowtraces []model.UserNowTrace
	usernowtraces = usernowtraceDao.GetUserNowTraceDbAll()
	for _,v := range usernowtraces{
	   if v.State != "离线"{
		  userhistorytraceDao.AddUserHistoryTrace(v)
	   }
	}
	//clear user now trace table
	usernowtraceDao.ClearAllUserNowTrace()
 }
func EveryDatAttenSheet(xlsx *excelize.File, datestr string){
	// Create a new sheet.
	newsheet := datestr
    index := xlsx.NewSheet(newsheet)
    // Set value of a cell.
    xlsx.SetCellValue(newsheet, "A1", "姓名")
	xlsx.SetCellValue(newsheet, "B1", "到岗")
	xlsx.SetCellValue(newsheet, "C1", "离岗")
	xlsx.SetCellValue(newsheet, "D1", "工时")
    // Set active sheet of the workbook.
	xlsx.SetActiveSheet(index)
	
	UserCnt := 0
	useraccounts := userDao.GetUserAll()
	for _,v := range useraccounts{
		UserCnt++
		var user_atten AttendUser
		user_atten.AccountName=v.Name
		user_atten.HunmanName=v.Human
		user_atten.Date=datestr
		GetAttendanByUsernameAndDate(&user_atten)
		UserCntStr := fmt.Sprintf("%d", UserCnt+1)
		xlsx.SetCellValue(newsheet, "A"+UserCntStr, user_atten.HunmanName)
		xlsx.SetCellValue(newsheet, "B"+UserCntStr, user_atten.LoginTime)
		xlsx.SetCellValue(newsheet, "C"+UserCntStr, user_atten.LogoutTime)
		xlsx.SetCellValue(newsheet, "D"+UserCntStr, fmt.Sprintf("%.1f", user_atten.Duration))
		//fmt.Println(user_atten)
	}
}
//at 24'clock, execute the timer function
func AttendanceTmrDay(){
	//DelAllUserNowTrace()
	// Save xlsx file by the given path.
	xlsx := excelize.NewFile()
	datestr := time.Now().Format("2006-01-02")
	EveryDatAttenSheet(xlsx, datestr)

	xlsFileName := datestr+"日考勤表.xlsx"
    err := xlsx.SaveAs(xlsFileName)
    if err != nil {
        fmt.Println(err)
	}

	//send email
	emailSub := datestr+"日考勤数据"
	utils.SendMailForDay(emailSub,xlsFileName)
}
//GetWeekPosi =0,current week, =1,last week
func AttendanceTmrWeekend(){
	
	currentTime := time.Now()
	datestr := currentTime.Format("2006-01-02")
	emailSub := datestr+"周考勤数据"

	xlsFileName := datestr+"周考勤表.xlsx"
	xlsx := excelize.NewFile()

	GetWeekPosi := 1
	var cellname string
	var user_atten AttendUser

	//total sum attendance
	totalsheet := "周工时汇总"
	xlsx.NewSheet(totalsheet)
	//total get current days of week
	xlsx.SetColWidth(totalsheet, "A", "A", 20)
	WeekDates := utils.GetWeekDate(GetWeekPosi)
	WeekDaysStr := utils.GetWeekDayStr(GetWeekPosi)
	rowindex := 2
	for _,v := range WeekDaysStr{
		cellname,_=excelize.CoordinatesToCellName(1, rowindex)
		xlsx.SetCellValue(totalsheet, cellname, v)
		rowindex++
	}
	cellname,_=excelize.CoordinatesToCellName(1, rowindex)
	xlsx.SetCellValue(totalsheet,cellname, "汇总")
	rowindex++
	cellname,_=excelize.CoordinatesToCellName(1, rowindex)
	xlsx.SetCellValue(totalsheet,cellname, "平均")

	UserColumn := 2
	//dateRow := 2
	useraccounts := userDao.GetUserAll()
	for _,v := range useraccounts{
		var totalHours float64
		
		//set human name
		cellname,_=excelize.CoordinatesToCellName(UserColumn, 1)
		xlsx.SetCellValue(totalsheet, cellname, v.Human)
		user_atten.AccountName=v.Name
		user_atten.HunmanName=v.Human
		
		rowindex := 2
		for _,v := range WeekDates{
			user_atten.Date = v.Format("2006-01-02")
			GetAttendanByUsernameAndDate(&user_atten)
			cellname,_=excelize.CoordinatesToCellName(UserColumn, rowindex)
			xlsx.SetCellValue(totalsheet, cellname, fmt.Sprintf("%.1f", user_atten.Duration))
			
			totalHours = totalHours+user_atten.Duration
			rowindex++
		}
		cellname,_=excelize.CoordinatesToCellName(UserColumn, rowindex)
		rowindex++
		xlsx.SetCellValue(totalsheet, cellname, fmt.Sprintf("%.1f", totalHours))

		cellname,_=excelize.CoordinatesToCellName(UserColumn, rowindex)
		rowindex++
		xlsx.SetCellValue(totalsheet, cellname, fmt.Sprintf("%.1f", totalHours/5))
		UserColumn++
	}

	//every day attendance data
	for _,v := range WeekDates{
		datestr = v.Format("2006-01-02")
		EveryDatAttenSheet(xlsx, datestr)
	}
	xlsx.DeleteSheet("Sheet1")
	xlsx.SetActiveSheet(0)
	// Save xlsx file by the given path.
    xlsx.SaveAs(xlsFileName)
	//send email
	utils.SendMailForWeek(emailSub,xlsFileName)
}

//GetMonthPosi=0,current month, =1,last month
func AttendanceTmrMonth(){

	currentTime := time.Now()
	year := currentTime.Year()
	monthStr := fmt.Sprintf("%02d",currentTime.Month()-1)
	//if current month =1,then get last year December month
	if(1 == currentTime.Month()){
		year = year-1
		monthStr = fmt.Sprintf("12")
	}

	datestr := currentTime.Format("2006-01-02")
	emailSub := datestr+"月考勤数据"

	xlsFileName := datestr+"月考勤表.xlsx"
	xlsx := excelize.NewFile()

	var cellname string
	var user_atten AttendUser

	//total sum attendance
	totalsheet := "月工时汇总"
	xlsx.NewSheet(totalsheet)
	//total get current days of week
	xlsx.SetColWidth(totalsheet, "A", "A", 20)
	MonthDates := utils.GetMonthDate(year, monthStr)
	MonthDaysStr := utils.GetMonthDateStr(year, monthStr)
	rowindex := 2
	for _,v := range MonthDaysStr{
		cellname,_=excelize.CoordinatesToCellName(1, rowindex)
		xlsx.SetCellValue(totalsheet, cellname, v)
		rowindex++
	}
	cellname,_=excelize.CoordinatesToCellName(1, rowindex)
	xlsx.SetCellValue(totalsheet,cellname, "汇总")
	rowindex++
	cellname,_=excelize.CoordinatesToCellName(1, rowindex)
	xlsx.SetCellValue(totalsheet,cellname, "平均")

	UserColumn := 2
	//dateRow := 2
	useraccounts := userDao.GetUserAll()
	for _,v := range useraccounts{
		var totalHours float64
		
		//set human name
		cellname,_=excelize.CoordinatesToCellName(UserColumn, 1)
		xlsx.SetCellValue(totalsheet, cellname, v.Human)
		user_atten.AccountName=v.Name
		user_atten.HunmanName=v.Human
		
		rowindex := 2
		for _,v := range MonthDates{
			user_atten.Date = v.Format("2006-01-02")
			GetAttendanByUsernameAndDate(&user_atten)
			cellname,_=excelize.CoordinatesToCellName(UserColumn, rowindex)
			xlsx.SetCellValue(totalsheet, cellname, fmt.Sprintf("%.1f", user_atten.Duration))
			
			totalHours = totalHours+user_atten.Duration
			rowindex++
		}
		cellname,_=excelize.CoordinatesToCellName(UserColumn, rowindex)
		rowindex++
		xlsx.SetCellValue(totalsheet, cellname, fmt.Sprintf("%.1f", totalHours))

		cellname,_=excelize.CoordinatesToCellName(UserColumn, rowindex)
		rowindex++
		xlsx.SetCellValue(totalsheet, cellname, fmt.Sprintf("%.1f", totalHours/5))
		UserColumn++
	}

	//every day attendance data
	for _,v := range MonthDates{
		datestr = v.Format("2006-01-02")
		EveryDatAttenSheet(xlsx, datestr)
	}
	xlsx.DeleteSheet("Sheet1")
	xlsx.SetActiveSheet(0)
	// Save xlsx file by the given path.
    xlsx.SaveAs(xlsFileName)
	//send email
	utils.SendMailForMonth(emailSub,xlsFileName)

}

