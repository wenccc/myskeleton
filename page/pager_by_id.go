package page

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"strconv"
)

type PagerByID[T any] struct {
	Page      uint
	TotalPage uint
	Limit     uint
	Total     uint
	Data      []T
	StartId   uint64
}

// PaginateById 通过Id来实现分页，而不是offset,但是如果起始页太小，也继续用
func PaginateById[T any](ctx *gin.Context, builder *gorm.DB) (*PagerByID[T], error) {
	var count int64
	builder.Count(&count)
	p := PagerByID[T]{
		Page:      getPage(ctx),
		TotalPage: 0,
		Limit:     getLimit(ctx),
		Total:     uint(count),
		Data:      make([]T, 0),
		StartId:   0,
	}
	p.TotalPage = uint(math.Ceil(float64(p.Total) / float64(p.Limit)))

	startId := getStartId(ctx)
	p.StartId = uint64(startId)

	var res *gorm.DB
	if startId > 0 {
		res = builder.Where("id > ?", startId).Limit(int(p.Limit)).Find(&p.Data)
	} else {
		res = builder.Offset(int(p.Page*p.Limit - p.Limit)).Limit(int(p.Limit)).Find(&p.Data)
	}

	if res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	return &p, nil
}

func getStartId(ctx *gin.Context) uint {

	id := ctx.DefaultQuery("start_id", "0")
	l, err := strconv.Atoi(id)
	if err != nil {
		return Limit
	}
	if l < 1 {
		return Limit
	}
	return uint(l)
}
