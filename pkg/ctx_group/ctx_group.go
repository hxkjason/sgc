package ctx_group

import (
	"context"
	"errors"
	"github.com/hxkjason/sgc/utils"
	"golang.org/x/sync/errgroup"
	"math"
)

func GetCtxCancelFuncGroup() (context.Context, context.CancelFunc, *errgroup.Group) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	group, errCtx := errgroup.WithContext(ctx)
	return errCtx, cancelFunc, group
}

// CheckGoroutineHasErr 检测goroutine是否有错误
func CheckGoroutineHasErr(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// CalcGoroutineNumAndOneDealNum 计算需要多少协程和每个协程处理的数量
func CalcGoroutineNumAndOneDealNum(rowsLen, planGoroutineNum, maxGoroutineNum int) (goroutineNum, oneGoroutineDealNum int, e error) {

	if rowsLen < 1 || planGoroutineNum < 1 {
		return 0, 0, errors.New("calc params is invalid")
	}
	// 总行数小于等于计划协程数
	if rowsLen <= planGoroutineNum {
		return 1, rowsLen, nil
	}
	// 限制最大协程数
	if planGoroutineNum > maxGoroutineNum {
		planGoroutineNum = maxGoroutineNum
	}
	// 计算每个协程的处理数量
	oneGoroutineDealNum = int(math.Ceil(utils.DecimalFloat(float64(rowsLen)/float64(planGoroutineNum), 6)))
	// 实际需要协程数量
	needGoroutineNum := int(math.Ceil(utils.DecimalFloat(float64(rowsLen)/float64(oneGoroutineDealNum), 6)))

	return needGoroutineNum, oneGoroutineDealNum, nil
}
