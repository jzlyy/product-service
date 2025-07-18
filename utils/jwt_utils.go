package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

// GenerateToken 生成JWT令牌
func GenerateToken(userID int, jwtSecret string) (string, error) {
	// 设置令牌有效期 (72小时)
	expirationTime := time.Now().Add(72 * time.Hour)

	// 创建声明
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证JWT令牌
func ValidateToken(tokenString, jwtSecret string) (int, error) {
	// 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	// 验证令牌
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 提取用户ID
		userID, ok := claims["user_id"]
		if !ok {
			return 0, errors.New("user_id claim is missing")
		}

		// 将用户ID转换为整数
		var id int
		switch v := userID.(type) {
		case float64:
			id = int(v)
		case int:
			id = v
		case string:
			id, err = strconv.Atoi(v)
			if err != nil {
				return 0, errors.New("invalid user_id format")
			}
		default:
			return 0, errors.New("invalid user_id type")
		}

		return id, nil
	}

	return 0, errors.New("invalid token")
}

// GetUserIDFromToken 从令牌中提取用户ID（不验证令牌）
func GetUserIDFromToken(tokenString string) (int, error) {
	// 注意：这个方法不验证令牌，只用于提取信息
	// 在生产环境中应谨慎使用

	// 解析令牌但不验证
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// 提取用户ID
		userID, ok := claims["user_id"]
		if !ok {
			return 0, errors.New("user_id claim is missing")
		}

		// 将用户ID转换为整数
		var id int
		switch v := userID.(type) {
		case float64:
			id = int(v)
		case int:
			id = v
		case string:
			id, err = strconv.Atoi(v)
			if err != nil {
				return 0, errors.New("invalid user_id format")
			}
		default:
			return 0, errors.New("invalid user_id type")
		}

		return id, nil
	}

	return 0, errors.New("invalid token claims")
}
