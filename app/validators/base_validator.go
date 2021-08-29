package validators

type BaseListPage struct {
	Page     int `form:"page" zh:"页码" en:"page" binding:"omitempty"`
	PageSize int `form:"page_size" zh:"页面条数" en:"Page size" binding:"omitempty"`
}



