package exportExcel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// docs: https://xuri.me/excelize/zh-hans/stream.html

type (
	ExportFileInfo struct {
		file         *excelize.File
		StreamWriter *excelize.StreamWriter
		sheetName    string //可定义默认sheet名称
	}

	// HeaderAttr excel 表头信息
	HeaderAttr struct {
		Index int     `json:"index"` // 列索引, 从 0 开始
		Title string  `json:"title"` // 列标题
		Key   string  `json:"key"`   // 对应的内容属性key
		Width float64 `json:"width"` // 列宽度
	}

	CommonAttr struct {
		FirstRowOpts            excelize.RowOpts
		ContentRowOpts          excelize.RowOpts
		FirstRowStyleId         int
		HorizontalLeftStyleId   int
		HorizontalCenterStyleId int
		QtyRightStyleId         int
		AmtRightStyleId         int
	}
)

var (
	varFilenamePrefix = "export-excel"                         // 默认文件前缀
	varSheetName      = "Sheet1"                               // 默认Sheet名称
	TopLineHeight     = 26.0                                   // 首行高度
	defaultHeight     = 20.0                                   // 默认行高度
	AlignHorizontal   = "left"                                 // 默认水平对齐方式
	CustomNumFmt1     = "￥#,##0.00;￥-#,##0.00"                 // 千分位 人民币金额
	CustomNumFmt2     = "#,##0.00;-#,##0.00"                   // 千分位 数字
	CustomNumFmt3     = "#,##0;-#,##0;_(* \"-\"_);_(@_)"       // 千分位 数字 整数，为0时展示 -
	CustomNumFmt4     = "#,##0.00;-#,##0.00;_(* \"-\"_);_(@_)" // 千分位 数字 2位小数，为0时展示 -
	FirstRowOpts      = excelize.RowOpts{Height: 20.1}
	ContentRowOpts    = excelize.RowOpts{Height: 18}
)

func NewStreamWriterExcel(filenamePrefix, alignHorizontal, sheetName string) *ExportFileInfo {
	// 文件名称前缀
	if filenamePrefix != "" {
		varFilenamePrefix = filenamePrefix
	}
	if sheetName != "" {
		varSheetName = sheetName
	}
	// 水平对其方式
	switch alignHorizontal {
	case "center", "right", "fill", "justify", "centerContinuous", "distributed":
		AlignHorizontal = alignHorizontal
	default:
		AlignHorizontal = "left"
	}
	xkExcelExportObj := &ExportFileInfo{file: createFile(), sheetName: varSheetName}
	_ = FreezeFirstRow(xkExcelExportObj.file, varSheetName)
	// 设置流式写入器
	xkExcelExportObj.StreamWriter, _ = xkExcelExportObj.file.NewStreamWriter(xkExcelExportObj.sheetName)
	return xkExcelExportObj
}

// SaveToLocal 保存在本地
func (efi *ExportFileInfo) SaveToLocal(headerAttrs []HeaderAttr, contentDataList []map[string]interface{}, fileDst string) {
	efi.WriteData(headerAttrs, contentDataList)
	if err := efi.file.SaveAs(fileDst); err != nil {
		fmt.Println(err)
	}
}

func (efi *ExportFileInfo) WriteData(params []HeaderAttr, data []map[string]interface{}) {
	efi.writeFileTop(params)
	efi.writeFileContent(params, data)
}

//设置首行
func (efi *ExportFileInfo) writeFileTop(headerAttrs []HeaderAttr) {
	topStyle, _ := efi.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: "宋体", Color: "#000000", Bold: false, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: AlignHorizontal, Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#d0f1d3"}, Pattern: 1},
	})
	topRowOpts := excelize.RowOpts{Height: TopLineHeight}

	titles := make([]interface{}, 0, len(headerAttrs))
	for _, headerAttr := range headerAttrs {
		titles = append(titles, excelize.Cell{Value: headerAttr.Title, StyleID: topStyle})
		columnIndex := headerAttr.Index
		// 设置列宽
		_ = efi.StreamWriter.SetColWidth(columnIndex+1, columnIndex+1, headerAttr.Width)
	}

	//  写入首行
	if err := efi.StreamWriter.SetRow("A1", titles, topRowOpts); err != nil {
		return
	}
}

// 设置文件内容
func (efi *ExportFileInfo) writeFileContent(params []HeaderAttr, data []map[string]interface{}) {

	contentStyle, _ := efi.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: "Microsoft YaHei", Color: "#000000", Bold: false, Size: 10},
		Alignment: &excelize.Alignment{Horizontal: AlignHorizontal, Vertical: "bottom"},
	})
	contentRowOpts := excelize.RowOpts{Height: defaultHeight, StyleID: contentStyle}

	columnNum := len(params)
	startRowId := 2
	// 写入内容
	for rowID := startRowId; rowID <= len(data)+startRowId-1; rowID++ {
		row := make([]interface{}, columnNum)
		for _, headerAttr := range params {
			row[headerAttr.Index] = excelize.Cell{Value: data[rowID-startRowId][headerAttr.Key], StyleID: contentStyle}
		}
		// 索引转单元格函数
		cell, _ := excelize.CoordinatesToCellName(1, rowID)
		_ = efi.StreamWriter.SetRow(cell, row, contentRowOpts)
	}
	_ = efi.StreamWriter.Flush()
}

func createFile() *excelize.File {
	f := excelize.NewFile()
	// 创建一个默认工作表
	sheetName := varSheetName
	index := f.NewSheet(sheetName)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	return f
}

func createFileName() string {
	return fmt.Sprintf("%s-%d.xlsx", varFilenamePrefix, time.Now().UnixMicro())
}

func WriteToCsv(filename, titles string, contentDataList []string) {
	//数据写入到csv文件
	t1 := time.Now()

	var stringBuilder strings.Builder
	stringBuilder.WriteString(titles)

	dataLen := len(contentDataList)

	for _, lineContent := range contentDataList {
		stringBuilder.WriteString(lineContent)
	}

	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModeAppend|os.ModePerm)
	dataString := stringBuilder.String()
	if _, err := file.WriteString(dataString); err != nil {
		fmt.Println("writeToCsv写入文件出错：" + err.Error())
	}

	err := file.Close()
	if err != nil {
		fmt.Println("writeToCsv关闭文件出错：" + err.Error())
	}

	fmt.Printf("writeToCsv总共%d条数据，总耗时%s\n", dataLen, time.Now().Sub(t1))
}

func ReadOneSheetAllRowsData(dst string, index int, removeAfterRead bool) (int, string, [][]string) {

	f, err := excelize.OpenFile(dst)
	if err != nil {
		fmt.Println(err)
		return 1, err.Error(), nil
	}

	defer func() {
		// 关闭文件
		if err = f.Close(); err != nil {
			fmt.Println(err)
		}
		// Remove file
		if removeAfterRead {
			if err = os.Remove(dst); err != nil {
				fmt.Println("删除文件失败,dst:", dst)
			}
		}
	}()

	// Get all the rows in the index Sheet.
	rows, err := f.GetRows(f.GetSheetName(index))
	if err != nil {
		return 1, err.Error(), nil
	}

	return 0, "ok", rows
}

func ReadOneSheetAllRowsDataOpenReader(file io.Reader, index int) (int, string, [][]string) {

	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Println(err)
		return 1, err.Error(), nil
	}

	defer func() {
		// 关闭文件
		if err = f.Close(); err != nil {
			fmt.Println(err)
		}

	}()

	// Get all the rows in the index Sheet.
	rows, err := f.GetRows(f.GetSheetName(index))
	if err != nil {
		return 1, err.Error(), nil
	}

	return 0, "ok", rows
}

// GetSheetColLetterName 获取表格列名称
// i 第几列
func GetSheetColLetterName(i int) string {
	const letterLen = 26
	letters := [letterLen]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	// 如果列数不大于 26，则直接返回列名称
	if i <= letterLen {
		return letters[i-1]
	}

	// 列数不能大于 26 * 27 列
	maxColNum := letterLen * (letterLen + 1)
	if i > maxColNum {
		panic("列数不能大于" + strconv.Itoa(maxColNum))
	}

	times := i / letterLen
	remainder := i % letterLen

	if remainder == 0 {
		return letters[times-2] + letters[letterLen-1]
	}
	return letters[times-1] + letters[remainder-1]
}

// BlackBorder 黑色边框样式
func BlackBorder() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}
}

// FreezeFirstRow // 冻结首行
func FreezeFirstRow(file *excelize.File, sheetName string) error {
	err := file.SetPanes(sheetName, `{
		"freeze": true,
		"split": false,
		"x_split": 0,
		"y_split": 1,
		"top_left_cell": "A2",
		"active_pane": "bottomLeft",
		"panes": [
		{
			"pane": "bottomLeft"
		}]
	}`)
	return err
}

// SetSheetTitleAndColWidth 普通写入 设置标题和列宽
func SetSheetTitleAndColWidth(file *excelize.File, sheetName string, headerAttrs []HeaderAttr) error {
	for _, headerAttr := range headerAttrs {
		colName := GetSheetColLetterName(headerAttr.Index + 1)
		if err := file.SetCellValue(sheetName, colName+"1", headerAttr.Title); err != nil {
			return err
		}
		if err := file.SetColWidth(sheetName, colName, colName, headerAttr.Width); err != nil {
			return err
		}
	}
	return nil
}

// StreamWriterSetSheetTitleAndColWidth 流式写入 设置标题和列宽
func StreamWriterSetSheetTitleAndColWidth(streamWriter *excelize.StreamWriter, headerAttrs []HeaderAttr, firstRowStyleID int, firstRowOpts excelize.RowOpts) error {
	// 设置首行标题、列宽
	firstRow := make([]interface{}, 0, 26)
	for _, headerAttr := range headerAttrs {
		firstRow = append(firstRow, excelize.Cell{StyleID: firstRowStyleID, Value: headerAttr.Title})
		// 设置列宽
		if err := streamWriter.SetColWidth(headerAttr.Index+1, headerAttr.Index+1, headerAttr.Width); err != nil {
			return err
		}
	}
	return streamWriter.SetRow("A1", firstRow, firstRowOpts)
}
