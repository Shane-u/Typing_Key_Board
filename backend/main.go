package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type TextSetting struct {
	ID     int64  `json:"id"`
	UserID string `json:"user_id"`
	Text   string `json:"text"`
	Title  string `json:"title"`
}

type TypingHistory struct {
	ID          int64     `json:"id"`
	UserID      string    `json:"user_id"`
	TextID      string    `json:"text_id"`
	ErrorCount  int       `json:"error_count"`
	Accuracy    int       `json:"accuracy"`
	CompletedAt time.Time `json:"completed_at"`
}

var db *sql.DB

func main() {
	// data目录
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatalf("创建数据目录失败: %v", err)
	}

	// SQLite
	var err error
	dbPath := filepath.Join("./data", "typing.db")
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 建表
	if err := createTables(); err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	// 设置Gin
	r := gin.Default()

	// 配置CORS TODO需要修改！
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://{修改为自己的服务器IP}:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 打字素材API路由
	r.GET("/api/texts/:user_id", getTexts)
	r.POST("/api/texts", saveText)
	r.PUT("/api/texts/:id", updateText)
	r.DELETE("/api/texts/:id", deleteText)

	// 历史记录API路由
	r.GET("/api/history/:user_id", getHistory)
	r.POST("/api/history", saveHistory)
	r.DELETE("/api/history/:id", deleteHistory)

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

func createTables() error {
	// 打字素材表
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS text_settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			text TEXT NOT NULL,
			title TEXT
		);
	`)
	if err != nil {
		return err
	}

	// 历史记录表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS typing_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			text_id TEXT NOT NULL,
			error_count INTEGER NOT NULL,
			accuracy INTEGER NOT NULL,
			completed_at TIMESTAMP NOT NULL
		);
	`)
	return err
}

// getTexts
func getTexts(c *gin.Context) {
	userID := c.Param("user_id")
	rows, err := db.Query("SELECT id, user_id, text, title FROM text_settings WHERE user_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var texts []TextSetting
	for rows.Next() {
		var text TextSetting
		if err := rows.Scan(&text.ID, &text.UserID, &text.Text, &text.Title); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		texts = append(texts, text)
	}

	c.JSON(http.StatusOK, texts)
}

func saveText(c *gin.Context) {
	var text TextSetting
	if err := c.BindJSON(&text); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 插入打字素材
	result, err := db.Exec("INSERT INTO text_settings (user_id, text, title) VALUES (?, ?, ?)",
		text.UserID, text.Text, text.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	text.ID = id
	c.JSON(http.StatusCreated, text)
}

// 更新打字素材
func updateText(c *gin.Context) {
	id := c.Param("id")
	var text TextSetting
	if err := c.BindJSON(&text); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE text_settings SET text = ?, title = ? WHERE id = ?", text.Text, text.Title, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// 删除打字素材
func deleteText(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM text_settings WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 获取历史记录
func getHistory(c *gin.Context) {
	userID := c.Param("user_id")
	rows, err := db.Query(`
		SELECT id, user_id, text_id, error_count, accuracy, completed_at 
		FROM typing_history 
		WHERE user_id = ? 
		ORDER BY completed_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var histories []TypingHistory
	for rows.Next() {
		var history TypingHistory
		var completedAtStr string
		if err := rows.Scan(&history.ID, &history.UserID, &history.TextID, &history.ErrorCount, &history.Accuracy, &completedAtStr); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 解析时间字符串
		history.CompletedAt, _ = time.Parse(time.RFC3339, completedAtStr)
		histories = append(histories, history)
	}

	c.JSON(http.StatusOK, histories)
}

func saveHistory(c *gin.Context) {
	var history TypingHistory
	if err := c.BindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "无法解析请求数据"})
		return
	}

	// 验证user_id和text_id
	if history.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id不能为空"})
		return
	}
	if history.TextID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text_id不能为空"})
		return
	}

	// √ 完成时间
	if history.CompletedAt.IsZero() {
		history.CompletedAt = time.Now()
	}

	result, err := db.Exec(`
		INSERT INTO typing_history (user_id, text_id, error_count, accuracy, completed_at) 
		VALUES (?, ?, ?, ?, ?)
	`, history.UserID, history.TextID, history.ErrorCount, history.Accuracy, history.CompletedAt.Format(time.RFC3339))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	history.ID = id
	c.JSON(http.StatusCreated, history)
}

func deleteHistory(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM typing_history WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
