package controllers

import (
	"database/sql"
	"go_dev/internal"
	"go_dev/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rows *sql.Rows
		var e error
		if rows, e = db.Query(`select guid, name, price, description, createdAt from products`); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.IndentedJSON(http.StatusInternalServerError, res)
			return
		}
		defer rows.Close()
		var products []models.Product
		for rows.Next() {
			var product models.Product
			if e := rows.Scan(&product.GUID, &product.Name, &product.Price, &product.Description, &product.CreatedAt); e != nil {
				var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
				c.IndentedJSON(http.StatusInternalServerError, res)
				return

			}
			products = append(products, product)
		}

		if len(products) == 0 {
			var res = internal.NewHTTPResponse(http.StatusNotFound, sql.ErrNoRows)
			c.IndentedJSON(http.StatusNotFound, res)
			return

		}
		var res = internal.NewHTTPResponse(http.StatusOK, products)
		c.IndentedJSON(http.StatusOK, res)
	}
}
