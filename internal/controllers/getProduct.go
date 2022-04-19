package controllers

import (
	"database/sql"
	"go_dev/internal"
	"go_dev/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type guiBinding struct {
	GUID string `uri:"guid" binding:"required,uuid4"`
}

func GetProduct(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var binding guiBinding
		var c = ctx.Request.Context()
		if e := ctx.ShouldBindUri(&binding); e != nil {
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			ctx.IndentedJSON(http.StatusBadRequest, res)
			return
		}

		var row = db.QueryRowContext(c, `select guid, name, price, description, createdAt from products where guid=$1`, binding.GUID)
		var product models.Product
		if e := row.Scan(&product.GUID, &product.Name, &product.Price, &product.Description, &product.CreatedAt); e != nil {
			if e == sql.ErrNoRows {
				var res = internal.NewHTTPResponse(http.StatusNotFound, e)
				ctx.IndentedJSON(http.StatusNotFound, res)
				return
			}
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			ctx.IndentedJSON(http.StatusInternalServerError, res)
			return
		}
		var res = internal.NewHTTPResponse(http.StatusOK, product)
		ctx.IndentedJSON(http.StatusOK, res)
	}
}
