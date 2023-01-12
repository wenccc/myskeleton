package page

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"strconv"
)

const (
	Limit = 10
)

type Pager[T any] struct {
	Page      uint
	TotalPage uint
	Limit     uint
	Total     uint
	Data      []T
}

func Paginate[T any](ctx *gin.Context, builder *gorm.DB) (*Pager[T], error) {
	var count int64
	builder.Count(&count)
	p := Pager[T]{
		Page:      getPage(ctx),
		TotalPage: 0,
		Limit:     getLimit(ctx),
		Total:     uint(count),
		Data:      make([]T, 0),
	}

	p.TotalPage = uint(math.Ceil(float64(p.Total) / float64(p.Limit)))
	res := builder.Offset(int(p.Page*p.Limit - p.Limit)).Limit(int(p.Limit)).Find(&p.Data)
	if res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	return &p, nil
}

func getPage(ctx *gin.Context) uint {
	page := ctx.DefaultQuery("page", "1")
	p, err := strconv.Atoi(page)
	if err != nil {
		return 1
	}
	if p < 1 {
		return 1
	}
	return uint(p)
}
func getLimit(ctx *gin.Context) uint {

	limit := ctx.DefaultQuery("limit", "10")
	l, err := strconv.Atoi(limit)
	if err != nil {
		return Limit
	}
	if l < 1 {
		return Limit
	}
	return uint(l)
}
