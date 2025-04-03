package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var amount int = 320

type updateAmountPayload struct {
	Amount int `json:"amount"`
}

func (app *Config) UpdateAmount(c *gin.Context) {
	var payload updateAmountPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed while updating water bill",
		})
		return
	}

	amount += payload.Amount
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Water bill updated successfully",
	})
}

func (app *Config) Water(c *gin.Context) {
	key := c.Param("api_key")

	if key == "c68cdd9a-e1ff-4d24-8ca9-b5fe18a649c8" {
		waterBills := []map[string]string{
			{
				"amount":   fmt.Sprintf("%d birr", amount),
				"due_date": "20-05-2017",
				"status":   "paid",
			},
			{
				"amount":   "270 birr",
				"due_date": "18-06-2017",
				"status":   "unpaid",
			},
			{
				"amount":   "680 birr",
				"due_date": "15-07-2017",
				"status":   "overdue",
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"bills": waterBills,
		})

	} else {
		waterBills := []map[string]string{
			{
				"amount":   "120 birr",
				"due_date": "10-03-2017",
				"status":   "paid",
			},
			{
				"amount":   "370 birr",
				"due_date": "13-04-2017",
				"status":   "paid",
			},
			{
				"amount":   "280 birr",
				"due_date": "19-07-2017",
				"status":   "unpaid",
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"bills": waterBills,
		})
	}
}

func (app *Config) Electric(c *gin.Context) {
	key := c.Param("api_key")

	if key == "c68cdd9a-e1ff-4d24-8ca9-b5fe18a649c8" {
		electricBills := []map[string]string{
			{
				"amount":   "400 birr",
				"due_date": "21-12-2016",
				"status":   "paid",
			},
			{
				"amount":   "370 birr",
				"due_date": "09-02-2017",
				"status":   "paid",
			},
			{
				"amount":   "777 birr",
				"due_date": "24-04-2017",
				"status":   "unpaid",
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"bills": electricBills,
		})

	} else {
		electricBills := []map[string]string{
			{
				"amount":   "550 birr",
				"due_date": "11-09-2016",
				"status":   "paid",
			},
			{
				"amount":   "5100 birr",
				"due_date": "19-05-2017",
				"status":   "unpaid",
			},
			{
				"amount":   "800 birr",
				"due_date": "07-07-2017",
				"status":   "overdue",
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"bills": electricBills,
		})
	}
}

func (app *Config) Internet(c *gin.Context) {
	key := c.Param("api_key")

	if key == "c68cdd9a-e1ff-4d24-8ca9-b5fe18a649c8" {
		internetBills := []map[string]string{
			{
				"amount":   "1100 birr",
				"due_date": "30-11-2016",
				"status":   "paid",
			},
			{
				"amount":   "1100 birr",
				"due_date": "30-03-2017",
				"status":   "unpaid",
			},
			{
				"amount":   "1100 birr",
				"due_date": "30-04-2017",
				"status":   "unpaid",
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"bills": internetBills,
		})

	} else {
		internetBills := []map[string]string{
			{
				"amount":   "3100 birr",
				"due_date": "05-12-2016",
				"status":   "paid",
			},
			{
				"amount":   "3100 birr",
				"due_date": "05-05-2017",
				"status":   "paid",
			},
			{
				"amount":   "3100 birr",
				"due_date": "05-07-2017",
				"status":   "unpaid",
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"bills": internetBills,
		})
	}
}
