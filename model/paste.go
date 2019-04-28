package model

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"time"
)

type Paste struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PasteID     string `gorm:"varchar(32) index:Hash"`
	FileName    string
	ContentType string
	Raw         []byte `gorm:"type:blob"`
}

func CreatePaste(c *gin.Context) {
	m, _ := c.MultipartForm()

	var res []string

	for _, context := range m.File {
		paste := ParsePaste(context)

		if paste != nil {
			Db.Create(paste)
			res = append(res, "http://"+c.Request.Host+c.Request.RequestURI+paste.PasteID)
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "result": res})
}

func FetchPaste(c *gin.Context) {
	paste := Paste{PasteID: c.Param("hash")}

	Db.Where(&paste).First(&paste)

	c.Header("Content-Type", paste.ContentType)

	switch paste.ContentType {
	case "application/octet-stream":
		c.Header("Content-Disposition", "attachment;filename="+paste.FileName)
		break
	case "image/jpeg":
		//c.Header("Content-Length", strconv.Itoa(len(paste.Raw)))
		break
	default:
		break
		//c.Header("Content-Disposition", "text/plain;charset=utf-8")
	}

	c.String(http.StatusOK, string(paste.Raw))
}

func UpdatePaste(c *gin.Context) {

	paste := Paste{PasteID: c.Param("hash")}

	m, _ := c.MultipartForm()

	Db.Where(&paste).First(&paste)

	if paste.Raw == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
		return
	}
	var res []string

	for _, context := range m.File {

		p := ParsePaste(context)

		if p != nil {
			paste.Raw = p.Raw
			Db.Model(&paste).Update(paste)
			res = append(res, "http://"+c.Request.Host+c.Request.RequestURI)
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "result": res})
}

func DeletePaste(c *gin.Context) {
	paste := Paste{PasteID: c.Param("hash")}

	Db.Where(&paste).First(&paste)

	if paste.Raw == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
		return
	}

	Db.Delete(&paste)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func ParsePaste(headers []*multipart.FileHeader) *Paste {

	for _, header := range headers {

		fileName := ""

		switch header.Header.Get("Content-Type") {
		case "application/octet-stream":
			fileName = header.Filename
		case "image/jpeg":
			fileName = header.Filename
		}

		content := make([]byte, 20971520)

		tmpFile, _ := header.Open()

		n, _ := tmpFile.Read(content)

		tGob := make([]byte, 8)

		binary.LittleEndian.PutUint32(tGob, uint32(time.Now().Unix()))

		hash := md5.Sum(bytes.Join([][]byte{tGob, content[:n]}, []byte("")))

		paste := Paste{FileName: fileName, ContentType: header.Header.Get("Content-Type"), Raw: content[:n], PasteID: hex.EncodeToString(hash[:])}

		return &paste
	}
	return nil
}
