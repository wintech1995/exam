package user

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	orm "wintech1995/my-app/orm"

	"github.com/gin-gonic/gin"
)

type ParamRecord struct {
	FileData string `json:"filedata" binding:"required"`
}

func ReadAll(c *gin.Context) {
	var users []orm.User

	orm.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "success",
		"users":   users,
	})
}

func Profile(c *gin.Context) {
	var user orm.User

	userId := c.MustGet("userId").(float64)
	orm.Db.Raw(`
		SELECT * 
		FROM users
		WHERE id = ?
	`, userId).Scan(&user)
	// orm.Db.First(&user, userId)
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "success",
		"user":    user,
	})
}

func UploadAvatar(c *gin.Context) {
	var json ParamRecord
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := strconv.Itoa(int(time.Now().Unix()))
	filepath := "img/avatar/" + r + ".jpg"
	base64toJpg(json.FileData, filepath)

	userId := c.MustGet("userId").(float64)

	orm.Db.Exec(`
		UPDATE users SET 
			avatar = ?
		WHERE
			id = ?
	`, filepath, userId)

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "upload success",
		"img":     filepath,
	})
}

// Take an existing jpg srcFileName and decode/encode it
func createJpg() {

	srcFileName := "flower.jpg"
	dstFileName := "newFlower.jpg"
	// Decode the JPEG data. If reading from file, create a reader with
	reader, err := os.Open(srcFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	//Decode from reader to image format
	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Got format String", formatString)
	fmt.Println(m.Bounds())

	//Encode from image format to writer
	f, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = jpeg.Encode(f, m, &jpeg.Options{Quality: 75})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Jpg file", dstFileName, "created")

}

// Take an existing png srcFileName and decode/encode it
func createPng() {
	srcFileName := "mouse.png"
	dstFileName := "newMouse.png"
	reader, err := os.Open(srcFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	//Decode from reader to image format
	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Got format String", formatString)
	fmt.Println(m.Bounds())

	//Encode from image format to writer
	f, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = png.Encode(f, m)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Png file", dstFileName, "created")

}

// Converts pre-existing base64 data (found in example of https://golang.org/pkg/image/#Decode) to test.png
func base64toPng(data string) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	fmt.Println(bounds, formatString)

	//Encode from image format to writer
	pngFilename := "test.png"
	f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = png.Encode(f, m)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Png file", pngFilename, "created")

}

// Given a base64 string of a JPEG, encodes it into an JPEG image test.jpg
func base64toJpg(data string, filepath string) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	fmt.Println("base64toJpg", bounds, formatString)

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = jpeg.Encode(f, m, &jpeg.Options{Quality: 75})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Jpg file", filepath, "created")
}

func getJPEGbase64(fileName string) string {
	imgFile, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	//fmt.Println("Base64 string is:", imgBase64Str)
	return imgBase64Str
}
