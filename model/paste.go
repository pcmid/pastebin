package model

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"time"
)

type PasteMetaData struct {
	ID          uint32 `gorm:"PRIMARY_KEY"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	MD5Hash     string `gorm:"varchar(32) index:Hash"`
	FileName    string
	ContentType string
}

type Paste struct {
	//ID          uint32 `gorm:"PRIMARY_KEY"`
	//CreatedAt   time.Time
	//UpdatedAt   time.Time
	//MD5Hash     string `gorm:"varchar(32) index:Hash"`
	//FileName    string
	//ContentType string
	PasteMetaData
	Raw []byte `gorm:"type:blob"`
}

func CreatePaste(c *gin.Context) {
	m, _ := c.MultipartForm()

	var res []string

	for _, context := range m.File {
		paste := ParsePaste(context)

		if paste != nil {
			Db.Create(paste).First(&paste)
			res = append(res, "http://"+c.Request.Host+c.Request.RequestURI+GenSortUrl(paste.ID))
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "result": res})
}

func FetchPaste(c *gin.Context) {
	paste := Paste{PasteMetaData: PasteMetaData{ID: UrlToID(c.Param("url"))}}

	Db.Where(&paste).First(&paste)

	if paste.Raw == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
		return
	}

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

func FetchMeta(c *gin.Context) {
	paste := Paste{PasteMetaData: PasteMetaData{ID: UrlToID(c.Param("url"))}}

	Db.Where(&paste).First(&paste)

	if paste.Raw == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "result": paste.PasteMetaData})

}

func UpdatePaste(c *gin.Context) {

	paste := Paste{PasteMetaData: PasteMetaData{ID: UrlToID(c.Param("url"))}}

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
			res = append(res, "http://"+c.Request.Host+c.Request.RequestURI+GenSortUrl(paste.ID))
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "result": res})
}

func DeletePaste(c *gin.Context) {
	paste := Paste{PasteMetaData: PasteMetaData{ID: UrlToID(c.Param("url"))}}

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

		hash := md5.Sum(content[:n])

		paste := Paste{PasteMetaData: PasteMetaData{FileName: fileName, ContentType: header.Header.Get("Content-Type"), MD5Hash: hex.EncodeToString(hash[:])}, Raw: content[:n]}

		return &paste
	}
	return nil
}
