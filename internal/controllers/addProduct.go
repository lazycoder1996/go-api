package controllers

import (
	"database/sql"
	"fmt"
	"go_dev/internal"
	"net/http"
	"time"

	"go_dev/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addProduct struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Description string  `json:"description" binding:"omitempty,max=250"`
}

func AddProduct(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload addProduct
		if e := ctx.ShouldBindJSON(&payload); e != nil {
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			ctx.IndentedJSON(http.StatusBadRequest, res)
			return
		}
		var guid = uuid.New().String()
		var createdAt = time.Now().Format(time.RFC3339)
		if _, e := db.Exec(`insert into products (guid, name, price, description, createdAt) values ($1, $2, $3, $4, $5)`, guid,
			payload.Name,
			payload.Price,
			payload.Description,
			createdAt,
		); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			ctx.IndentedJSON(http.StatusInternalServerError, res)
			return
		}
		var product models.Product
		var row = db.QueryRow(`select guid, name, price, description, createdAt from products where guid=$1`, guid)
		if e := row.Scan(&product.GUID, &product.Name, &product.Price, &product.Description, &product.CreatedAt); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			ctx.IndentedJSON(http.StatusInternalServerError, res)
			return
		}
		var res = internal.NewHTTPResponse(http.StatusCreated, product)
		ctx.Writer.Header().Add("Location", fmt.Sprintf("/products/%s", guid))
		ctx.IndentedJSON(http.StatusOK, res)
	}
}
