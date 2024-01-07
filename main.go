package main

import (
	todo_txt "github.com/1set/todotxt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const default_local_html_path = "./html"
const layout = "2006-01-02"

func parseTodoTmpFile(tmp_file *os.File) (todo_txt.TaskList, error) {
	parsed_tasks, err := todo_txt.LoadFromFile(tmp_file)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return parsed_tasks, nil
}

func uploadTempFile(c *gin.Context) {
	// check hash
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		c.Abort()
	}
	log.Println(file.Filename)
	// Write to tempfile
	f, err := os.CreateTemp("", "todo_tmp_*.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		c.Abort()
	}
	defer f.Close()
	c.SaveUploadedFile(file, f.Name())
	log.Println(f.Name())
	// Send URI to client
	c.JSON(http.StatusOK, gin.H{
		"filename": f.Name(),
	})
}

func getUserAgent(c *gin.Context) string {
	ua := c.Request.Header.Get("User-Agent")
	return ua
}

func getTodo(c *gin.Context) {
	query := c.Request.URL.Query()
	log.Println(query)
	ua := getUserAgent(c)
	log.Println(ua)
	tmp_file_name := c.Query("filename")
	temp_file_name_decoded, _ := url.QueryUnescape(tmp_file_name)
	tmp_file, err := os.Open(temp_file_name_decoded)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		log.Fatal(err)
		c.Abort()
	}
	defer tmp_file.Close()
	parsed_tasks, err := parseTodoTmpFile(tmp_file)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
	var matched_tasks []todo_txt.Task
	for _, task := range parsed_tasks {

		if ShouldSkip(query, "completed_only", strconv.FormatBool(task.Completed)) ||
			ShouldSkip(query, "uncompleted_only", strconv.FormatBool(!task.Completed)) ||
			ShouldSkip(query, "due_date", task.DueDate.Format(layout)) ||
			ShouldSkip(query, "priority", task.Priority) ||
			ShouldSkip(query, "created_date", task.CreatedDate.Format(layout)) {
			continue
		}

		matched_tasks = append(matched_tasks, task)
	}
	c.JSON(http.StatusOK, gin.H{
		"matched_tasks": matched_tasks,
		"DEBUG":         "DEBUG",
	})
	log.Println("TASK:", matched_tasks)
}

func indexHandler(c *gin.Context) {
	ua := getUserAgent(c)
	log.Println(ua)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Todo.txt",
	})
}

func deleteTempFile(tmp_file_name string) {
	os.Remove(tmp_file_name)
}

func deletTempFileHandler(c *gin.Context) {
	tmp_file_name := c.Query("filename")
	temp_file_name_decoded, err := url.QueryUnescape(tmp_file_name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
	deleteTempFile(temp_file_name_decoded)
	c.JSON(http.StatusOK, gin.H{
		"filename": temp_file_name_decoded,
	})
}

func main() {
	log.Println("start server")
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/filter", getTodo)
		v1.POST("/upload", AuthToken(), uploadTempFile)
		v1.DELETE("/delete", AuthToken(), deletTempFileHandler)
	}

	router.Run(":8080")
	log.Println("end server")
}
