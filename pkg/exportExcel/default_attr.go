package exportExcel

import "github.com/xuri/excelize/v2"

// GetDefaultAttr 获取文件默认属性
func GetDefaultAttr(f *excelize.File) CommonAttr {
	return CommonAttr{
		FirstRowOpts:            FirstRowOpts,
		ContentRowOpts:          ContentRowOpts,
		FirstRowStyleId:         FirstRowVerDistributedStyleId(f),
		HorizontalLeftStyleId:   DefaultHorizontalLeftStyleId(f),
		HorizontalCenterStyleId: DefaultHorizontalCenterStyleId(f),
		QtyRightStyleId:         DefaultQtyRightStyleId(f),
		AmtRightStyleId:         DefaultAmtRightStyleId(f),
	}
}

// DefaultFirstRowStyleId 默认首行样式
func DefaultFirstRowStyleId(f *excelize.File) (firstRowStyleId int) {
	firstRowStyleId, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: "宋体", Color: "000000", Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    BlackBorder(),
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#e7f1de"}, Pattern: 1},
	})
	return firstRowStyleId
}

// FirstRowVerDistributedStyleId 首行样式 - 垂直方向分布式
func FirstRowVerDistributedStyleId(f *excelize.File) (firstRowStyleId int) {
	firstRowStyleId, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: "宋体", Color: "000000", Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "distributed"},
		Border:    BlackBorder(),
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#e7f1de"}, Pattern: 1},
	})
	return firstRowStyleId
}

// DefaultHorizontalLeftStyleId 默认水平左对齐样式
func DefaultHorizontalLeftStyleId(f *excelize.File) (horizontalLeftStyleId int) {
	horizontalLeftStyleId, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: "宋体", Color: "000000", Bold: false, Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border:    BlackBorder(),
	})
	return horizontalLeftStyleId
}

// DefaultHorizontalCenterStyleId 默认水平居中对齐样式
func DefaultHorizontalCenterStyleId(f *excelize.File) (horizontalLeftStyleId int) {
	horizontalLeftStyleId, _ = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: "宋体", Color: "000000", Bold: false, Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    BlackBorder(),
	})
	return horizontalLeftStyleId
}

// DefaultQtyRightStyleId 默认数量右对齐样式
func DefaultQtyRightStyleId(f *excelize.File) (qtyRightStyleId int) {
	qtyRightStyleId, _ = f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Family: "宋体", Color: "000000", Bold: false, Size: 11},
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		Border:       BlackBorder(),
		CustomNumFmt: &CustomNumFmt3,
	})
	return qtyRightStyleId
}

// DefaultAmtRightStyleId 默认金额右对齐样式
func DefaultAmtRightStyleId(f *excelize.File) (amtRightStyleId int) {
	amtRightStyleId, _ = f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Family: "宋体", Color: "000000", Bold: false, Size: 11},
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		Border:       BlackBorder(),
		CustomNumFmt: &CustomNumFmt4,
	})
	return amtRightStyleId
}
