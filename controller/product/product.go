package product

import (
	"net/http"
	orm "wintech1995/my-app/orm"

	"github.com/gin-gonic/gin"
)

type ParamRecord struct {
	Status string `json:"status" binding:"required"`
}

type ProductRecord struct {
	Id        int
	CreatedAt string
	Name      string
	Price     string
	Status    int
	ShopName  string
}

func ProductList(c *gin.Context) {
	var product []orm.Product

	orm.Db.Raw(`
		SELECT * 
		FROM tb_product
		WHERE id <> 0
	`).Scan(&product)

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "success",
		"product": product,
	})
}

func ProductData(c *gin.Context) {
	var product orm.Product

	id := c.Param("id")
	orm.Db.Raw(`
		SELECT * 
		FROM tb_product
		WHERE id = ?
	`, id).Scan(&product)

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "success",
		"product": product,
	})
}

func UpdateStatus(c *gin.Context) {
	var product []ProductRecord
	var json ParamRecord
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orm.Db.Exec(`
		UPDATE tb_product SET
			status = ?
		WHERE
			id <> 0
	`, json.Status)

	orm.Db.Raw(`
		SELECT pd.id, pd.name, pd.status, sh.name AS shop_name, 
			FORMAT(pd.price, 2) AS price, 
			DATE_FORMAT(pd.created_at, '%d/%m/%Y') AS created_at
		FROM tb_product pd, tb_shop sh
		WHERE
			sh.id = pd.shop_id AND 
			pd.id <> 0
	`).Scan(&product)

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "update success",
		"product": product,
	})
}
