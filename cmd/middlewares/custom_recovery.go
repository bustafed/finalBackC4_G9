package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RecoveryWithLog(ctx *gin.Context, err any) {

	fmt.Printf("error is %+v\n", err)

	fmt.Printf("method is %+v\n", ctx.Request.Method)

	fmt.Printf("path is %+v\n", ctx.Request.URL.Path)

	fmt.Printf("the time is %+v\n", time.Now())

	ctx.AbortWithError(http.StatusInternalServerError, errors.New("internal server error"))
}
