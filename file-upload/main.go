package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Transer struct {
	ID       uint   `json:"id"`
	Filename string `json:"filename"`
	Filesize int64  `json:"filesize"`
}

var db *gorm.DB
var err error

// file handler for route "/upload"
func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading file...")
	r.ParseMultipartForm(10 << 20) // "10 << 20" means that the maximum file size is 10MB

	file, handler, err := r.FormFile("fileUploadForm")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Uploaded file: %+v\n", handler.Filename)
	fmt.Printf("File size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// make sure we close the file
	defer file.Close()

	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	logItem := &Transer{
		Filename: handler.Filename,
		Filesize: handler.Size,
	}

	db.Create(&logItem)

	fmt.Println("File upload successfull")
}

func setupRoutes() {
	// serve html
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// add route
	http.HandleFunc("/upload", uploadFile)
	fmt.Println("Server is up and running!")

	// host the server on port 8080
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Setting up file-upload server...")

	db, err = gorm.Open("sqlite3", "./transfers.db")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Connected to database!")

	defer db.Close()
	db.AutoMigrate(&Transer{})

	router := gin.Default()
	router.GET("/transers", getTransfers)
	router.DELETE("/transfer/:id", deleteTransferEntry)

	// init routes
	setupRoutes()
}

func getTransfers(c *gin.Context) {
	var transfers []Transer
	if err := db.Find(&transfers).Error; err != nil {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, transfers)
	}
}

func deleteTransferEntry(c *gin.Context) {
	id := c.Params.ByName("id")
	var transfer Transer

	d := db.Where("id = ?", id).Delete(&transfer)
	fmt.Println(d)

	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
