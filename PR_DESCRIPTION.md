# Issue #2630 - FormFileContent 功能实现

## 功能描述

新增 `Context.FormFileContent()` 方法，允许用户直接获取上传文件的二进制内容，便于计算哈希值、内容检查等操作。

## 背景

Issue #2630 提到用户需要在上传文件后获取文件内容并计算哈希值（如 SHA1）。原有 `FormFile()` 方法只返回文件头信息，需要手动打开文件读取内容。

## 新增方法

```go
// FormFileContent returns the first file content and its hash for the provided form key.
// It returns the file header, content bytes, and an error if any.
// The content is read into memory, so it should be used with caution for large files.
func (c *Context) FormFileContent(name string) (*multipart.FileHeader, []byte, error)
```

## 使用示例

### 示例 1：获取文件并计算 SHA1 哈希

```go
package main

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadFile(c *gin.Context) {
	fileHeader, content, err := c.FormFileContent("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 计算 SHA1 哈希
	hash := sha1.Sum(content)
	hashHex := hex.EncodeToString(hash[:])

	c.JSON(http.StatusOK, gin.H{
		"file_name":  fileHeader.Filename,
		"size":       len(content),
		"sha1_hash":  hashHex,
	})
}

func main() {
	r := gin.Default()
	r.POST("/upload", UploadFile)
	r.Run(":8080")
}
```

### 示例 2：结合保存文件到磁盘

```go
func UploadAndProcess(c *gin.Context) {
	fileHeader, content, err := c.FormFileContent("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 计算哈希
	hash := sha256.Sum256(content)
	hashHex := hex.EncodeToString(hash[:])

	// 2. 内容检查
	if len(content) > 10<<20 { // 10MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
		return
	}

	// 3. 保存到磁盘
	if err := c.SaveUploadedFile(fileHeader, "/uploads/"+fileHeader.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_name":  fileHeader.Filename,
		"hash":       hashHex,
		"saved":      true,
	})
}
```

## 变更文件

| 文件 | 变更 |
|------|------|
| `context.go` | 新增 `FormFileContent()` 方法 |
| `context_test.go` | 新增 `TestContextFormFileContent()` 和 `TestContextFormFileContentFailed()` 测试 |

## 测试

```bash
go test -v -run "TestContextFormFileContent" ./...
```

结果：通过

## 兼容性

- **向后兼容**：新增方法，不影响现有代码
- **注意**：文件内容会读入内存，大文件请谨慎使用（建议配合 `MaxMultipartMemory` 限制）

## 关联 Issue

- Closes #2630
