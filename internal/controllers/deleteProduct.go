package controllers

import (
	"database/sql"
	"go_dev/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var binding guiBinding
		var c = ctx.Request.Context()
		if e := ctx.ShouldBindUri(&binding); e != nil {
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			ctx.IndentedJSON(http.StatusBadRequest, res)
			return
		}
		var result sql.Result
		var e error
		if result, e = db.ExecContext(c, `delete from products where guid=$1`, binding.GUID); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			ctx.IndentedJSON(http.StatusInternalServerError, res)
			return
		}
		if nProducts, _ := result.RowsAffected(); nProducts == 0 {
			var res = internal.NewHTTPResponse(http.StatusNotFound, sql.ErrNoRows)
			ctx.IndentedJSON(http.StatusNotFound, res)
			return

		}

		ctx.IndentedJSON(http.StatusNoContent, nil)
	}
}
