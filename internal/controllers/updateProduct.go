package controllers

import (
	"database/sql"
	"go_dev/internal"
	"go_dev/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type putProduct struct {
	Name        string  `json:"name" binding:"required_without_all=Price Description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Description string  `json:"description" binding:"omitempty,max=250"`
}

func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var binding guiBinding
		var payload putProduct

		if e := ctx.ShouldBindUri(&binding); e != nil {
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			ctx.IndentedJSON(http.StatusBadRequest, res)
			return
		}
		if e := ctx.ShouldBindJSON(&payload); e != nil {
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			ctx.IndentedJSON(http.StatusBadRequest, res)
			return
		}
		var row = db.QueryRow(`select name, price, description from products where guid = $1`, binding.GUID)

		var currentProduct models.Product
		if e := row.Scan(&currentProduct.Name, &currentProduct.Price, &currentProduct.Description); e != nil {
			if e == sql.ErrNoRows {
				var res = internal.NewHTTPResponse(http.StatusNotFound, sql.ErrNoRows)
				ctx.IndentedJSON(http.StatusNotFound, res)
				return

			}
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			ctx.IndentedJSON(http.StatusInternalServerError, res)
			return

		}
		var option = copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
		}
		if e := copier.CopyWithOption(&currentProduct, &payload, option); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			ctx.IndentedJSON(http.StatusInternalServerError, res)
			return

		}
		if _, e := db.Exec(`update products set name=$2, price=$3, description=$4 where guid=$1`,
			binding.GUID, currentProduct.Name, currentProduct.Price, currentProduct.Description); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			ctx.IndentedJSON(http.StatusInternalServerError, res)
			return

		}

		var updateRow = db.QueryRow(`select guid, name, price, description, createdAt from products where guid = $1`, binding.GUID)

		var product models.Product
		if e := updateRow.Scan(&product.GUID, &product.Name, &product.Price, &product.Description, &product.CreatedAt); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			ctx.IndentedJSON(http.StatusInternalServerError, res)
			return

		}
		var res = internal.NewHTTPResponse(http.StatusOK, product)
		ctx.IndentedJSON(http.StatusOK, res)
	}

}
