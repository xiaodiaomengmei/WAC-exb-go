package files

import "os"

func Path() string {
	pwd,_:=os.Getwd()
	targetPath := pwd+"\\files\\excelfile.xlsx"
	return targetPath
}
