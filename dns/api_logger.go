package dns

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// APILogger API 响应日志器，负责将原始响应存入文件并保留最近100条
type APILogger struct {
	logPath  string // 日志文件路径，位于配置目录下的API.log
	maxLines int    // 最大保留日志行数
	mu       sync.Mutex // 保证日志写入线程安全
}

// NewAPILogger 初始化API日志器
// configDir: 配置文件所在目录
// maxLines: 最大保留日志条数
func NewAPILogger(configDir string, maxLines int) *APILogger {
	logPath := filepath.Join(configDir, "API.log")
	return &APILogger{
		logPath:  logPath,
		maxLines: maxLines,
	}
}

// WriteLog 写入API原始响应日志，自动截断超出最大行数的旧日志
func (l *APILogger) WriteLog(rawResp string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 1. 读取现有日志内容
	var lines []string
	file, err := os.Open(l.logPath)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		_ = file.Close()
	}

	// 2. 追加新日志（带时间戳）
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	newLogLine := fmt.Sprintf("[%s] %s", timestamp, strings.TrimSpace(rawResp))
	lines = append(lines, newLogLine)

	// 3. 截断日志，只保留最近maxLines条
	if len(lines) > l.maxLines {
		lines = lines[len(lines)-l.maxLines:]
	}

	// 4. 写入文件
	file, err = os.Create(l.logPath)
	if err != nil {
		return fmt.Errorf("创建API.log失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("写入日志行失败: %w", err)
		}
	}
	return writer.Flush()
}