package http_result

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
	"math"
	"strconv"
)

type Meta struct {
	TotalPage   int   `json:"total_page"`
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
}

type Links struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
}

type PageResult struct {
	Data  any   `json:"data"` // List of data
	Total int64 `json:"total"`
	Links Links `json:"links"`
	Meta  Meta  `json:"meta"`
}

// SearchByParams
// ?name=xxx&pageSize=1&currentPage=1&sort=xxx&order=xxx
func (h *HttpResult) SearchByParams(params map[string]string, excepts ...string) *HttpResult {
	for _, except := range excepts {
		delete(params, except)
	}
	query := facades.Orm().Query()
	h.Query = func(q orm.Query) orm.Query {
		for key, value := range params {
			if value == "" || key == "pageSize" || key == "total" || key == "currentPage" || key == "sort" || key == "order" {
				continue
			} else {
				q = q.Where(key+" like ?", "%"+value+"%")
			}
		}
		return q
	}(query)
	return h
}

func (r *HttpResult) ResultPagination(ctx http.Context, dest any) (http.Response, error) {
	r.Context = ctx
	request := ctx.Request()
	pageSize := request.Query("pageSize", "10")
	pageSizeInt := cast.ToInt(pageSize)
	currentPage := request.Query("currentPage", "1")
	currentPageInt := cast.ToInt(currentPage)
	total := int64(0)
	r.Query.Model(dest).Paginate(currentPageInt, pageSizeInt, dest, &total)

	URL_PATH := ctx.Request().Origin().URL.Path
	var proto = "http://"
	if request.Origin().TLS != nil {
		proto = "https://"
	}
	// Corrected links generation
	links := Links{
		First: proto + request.Origin().Host + URL_PATH + "?pageSize=" + pageSize + "&currentPage=1",
		Last:  proto + request.Origin().Host + URL_PATH + pageSize + "&currentPage=" + strconv.Itoa(int(total)/pageSizeInt),
		Prev:  proto + request.Origin().Host + URL_PATH + pageSize + "&currentPage=" + strconv.Itoa(currentPageInt-1),
		Next:  proto + request.Origin().Host + URL_PATH + pageSize + "&currentPage=" + strconv.Itoa(currentPageInt+1),
	}

	// Corrected total page calculation
	totalPage := int(math.Ceil(float64(total) / float64(pageSizeInt)))

	meta := Meta{
		TotalPage:   totalPage,
		CurrentPage: currentPageInt,
		PerPage:     pageSizeInt,
		Total:       total,
	}

	pageResult := PageResult{
		Data:  dest,
		Total: total,
		Links: links,
		Meta:  meta,
	}

	// 返回构建好的分页结果
	return r.Context.Response().Success().Json(pageResult), nil
}
