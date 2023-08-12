package chat

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ParamRecord struct {
	Chat string `json:"chat" binding:"required"`
}

func SendChat(c *gin.Context) {
	var json ParamRecord
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url := "https://api.line.me/v2/bot/message/multicast"
	method := "POST"

	payload := strings.NewReader(`{
		"to": ["U683008a38cc42b43b2ebbc11fd4d84bb"],
		"messages": [{
			"type": "text",
			"text": "` + json.Chat + `"
		}]
  	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer HYBE1QJLmn5LKUOsUElxM+eNZ4Rd+68QIi8KmLMQklqXwg3gGR9wFR9vRlOp55Gjw5sfBGfy6kCC0yFsKP7WRXVYkEPD6oqj6NCtOG+n2ooX7H9AYFShFsu9MPu6ASfCwogR1vTky8Ka+NfxX3vdeQdB04t89/1O/w1cDnyilFU=")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(body))

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "success",
		"chat":    json.Chat,
		"result":  string(body),
	})
}
