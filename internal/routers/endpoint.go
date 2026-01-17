package routers

import (
	"context"

	"github.com/gin-gonic/gin"
)

func handleResult[T any](handler func(c *gin.Context) (T, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := handler(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.JSON(200, result)
	}
}

type EndpointFunc[P, R any] func(ctx context.Context, param P) (R, error)

type BinderFunc[P any] func(ctx *gin.Context) (P, error)

func Endpoint[P, R any](
	endpoint EndpointFunc[P, R],
	binder BinderFunc[P],
) gin.HandlerFunc {
	return handleResult(func(c *gin.Context) (R, error) {
		param, err := binder(c)
		if err != nil {
			var zero R
			return zero, err
		}
		ctx := c.Request.Context()
		return endpoint(ctx, param)
	})
}

func JSONBinder[P any]() BinderFunc[P] {
	return func(c *gin.Context) (P, error) {
		var param P
		if err := c.ShouldBindJSON(&param); err != nil {
			return param, err
		}
		return param, nil
	}
}

func QueryBinder[P any]() BinderFunc[P] {
	return func(c *gin.Context) (P, error) {
		var param P
		if err := c.ShouldBindQuery(&param); err != nil {
			return param, err
		}
		return param, nil
	}
}

func PathBinder[P any]() BinderFunc[P] {
	return func(c *gin.Context) (P, error) {
		var param P
		if err := c.ShouldBindUri(&param); err != nil {
			return param, err
		}
		return param, nil
	}
}

func EmptyBinder[P any]() BinderFunc[P] {
	return func(c *gin.Context) (P, error) {
		var param P
		return param, nil
	}
}
