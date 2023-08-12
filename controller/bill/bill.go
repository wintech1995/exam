package bill

import (
	"fmt"
	"net/http"
	orm "wintech1995/my-app/orm"

	"github.com/gin-gonic/gin"
)

type BillRecord struct {
	Id        int
	CreatedAt string
	InvoiceNo string
	Status    int
	Total     float64
	Vat       float64
	Nettotal  float64
}

type OrderRecord struct {
	Id       int
	Price    float64
	Quantity float64
	Amount   float64
}

func ProductData(c *gin.Context) {
	var bill BillRecord
	var order []OrderRecord

	id := c.Param("id")

	orm.Db.Raw(`
		SELECT id, invoice_no, status, total, vat, nettotal,
			DATE_FORMAT(created_at, '%d/%m/%Y') AS created_at
		FROM tb_bill
		WHERE id = ?
	`, id).Scan(&bill)

	orm.Db.Raw(`
		SELECT id, bill_id, price, quantity, amount
		FROM tb_billorder
		WHERE bill_id = ?
	`, id).Scan(&order)

	var total float64
	var vat float64
	var nettotal float64

	for _, item := range order {
		fmt.Println(item)
		item.Amount = item.Price * item.Quantity

		total += item.Amount

		orm.Db.Exec(`
			UPDATE tb_billorder SET 
				amount = ?
			WHERE
				id = ?
		`, item.Amount, item.Id)
	}

	vat = total * 0.07
	nettotal = total + vat

	orm.Db.Exec(`
		UPDATE tb_bill SET 
			total = ?,
			vat = ?,
			nettotal = ?
		WHERE
			id = ?
	`, total, vat, nettotal, id)

	orm.Db.Raw(`
		SELECT id, invoice_no, status, total, vat, nettotal,
			DATE_FORMAT(created_at, '%d/%m/%Y') AS created_at
		FROM tb_bill
		WHERE id = ?
	`, id).Scan(&bill)

	orm.Db.Raw(`
		SELECT id, bill_id, price, quantity, amount
		FROM tb_billorder
		WHERE bill_id = ?
	`, id).Scan(&order)

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "success",
		"bill":    bill,
		"order":   order,
	})
}
