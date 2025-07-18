package middlewares

import (
	"log"
	"net/http"
	"product-service/config"
	"product-service/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 验证JWT令牌的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// 检查Bearer格式
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Expected 'Bearer <token>'",
			})
			return
		}

		// 提取令牌
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 加载配置获取JWT密钥
		cfg := config.LoadConfig()
		jwtSecret := cfg.JWTSecret

		// 验证令牌
		userID, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
			return
		}

		// 设置用户ID到上下文
		c.Set("userID", userID)
		c.Next()
	}
}

// AdminMiddleware 验证管理员权限的中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行基础认证
		userIDValue, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		userID := userIDValue.(int)

		// 在实际应用中，这里应该查询数据库验证用户角色
		// 为简化演示，我们假设用户ID为1的是管理员
		if userID != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Admin privileges required",
			})
			return
		}

		c.Next()
	}
}

// CORSMiddleware 处理跨域请求的中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware 简单的速率限制中间件
func RateLimitMiddleware() gin.HandlerFunc {
	// 在实际应用中应该使用更健壮的解决方案如redis
	// 这里使用简单的内存计数器作为示例
	requestCount := make(map[string]int)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 每分钟限制60个请求
		if count, exists := requestCount[clientIP]; exists && count >= 60 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}

		requestCount[clientIP]++

		// 每分钟重置计数器（在实际应用中应该使用定时器）
		// 这里只是示例，生产环境应该使用更可靠的方案
		go func(ip string) {
			<-time.After(time.Minute)
			if _, exists := requestCount[ip]; exists {
				requestCount[ip] = 0
			}
		}(clientIP)

		c.Next()
	}
}

// LoggingMiddleware 请求日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录请求完成后的信息
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if errorMessage != "" {
			log.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s?%s\n%s",
				end.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
				query,
				errorMessage,
			)
		} else {
			log.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s?%s",
				end.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
				query,
			)
		}
	}
}
