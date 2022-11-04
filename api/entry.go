package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mrityunjaygr8/simplebank/db/sqlc"
)

type listEntriesParams struct {
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
}

func (server *Server) listEntries(ctx *gin.Context) {
	var req listEntriesParams

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListEntriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	entries, err := server.store.ListEntries(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

type listEntriesForAccountQueryParams struct {
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
}

type listEntriesForAccountURI struct {
	AccountID int64 `uri:"account_id" binding:"required,min=1"`
}

func (server *Server) listEntriesForAccount(ctx *gin.Context) {
	var queryParams listEntriesForAccountQueryParams
	var req listEntriesForAccountURI

	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListEntriesForAccountParams{
		Limit:     queryParams.PageSize,
		Offset:    (queryParams.PageID - 1) * queryParams.PageSize,
		AccountID: req.AccountID,
	}

	entries, err := server.store.ListEntriesForAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}
