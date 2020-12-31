package login

import (
	"fmt"
	"github.com/tealeg/xlsx"
)


func ReadExcel(excelFilenameing string) []string{
	var strs []string
	xlFile, err := xlsx.OpenFile(excelFilenameing)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
	}
	for _, sheet := range xlFile.Sheets {
		fmt.Printf("Sheet Name: %s\n", sheet.Name)
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text := cell.String()
				strs = append(strs,text)
			}
		}
	}
	return strs
}




