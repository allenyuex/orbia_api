package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/smtp"
	"strings"

	"orbia_api/biz/infra/config"
)

// 验证码邮件模板
const (
	// EmailVerificationCodeTemplate 验证码邮件模板（包含邮件头部）
	EmailVerificationCodeTemplate = `Subject: Orbia App - Your Verification Code
From: {{.From}}
To: {{.To}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verification Code</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 10px;
            padding: 40px;
            text-align: center;
        }
        .content {
            background: white;
            border-radius: 8px;
            padding: 40px;
            margin-top: 20px;
        }
        .logo {
            font-size: 32px;
            font-weight: bold;
            color: white;
            margin-bottom: 10px;
        }
        .subtitle {
            color: rgba(255, 255, 255, 0.9);
            font-size: 16px;
            margin-bottom: 20px;
        }
        h1 {
            color: #333;
            font-size: 24px;
            margin-bottom: 20px;
        }
        .code-container {
            background: #f7fafc;
            border: 2px dashed #cbd5e0;
            border-radius: 8px;
            padding: 20px;
            margin: 30px 0;
        }
        .code {
            font-size: 36px;
            font-weight: bold;
            color: #667eea;
            letter-spacing: 8px;
            font-family: 'Courier New', monospace;
        }
        .info {
            color: #718096;
            font-size: 14px;
            margin: 20px 0;
        }
        .warning {
            background: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 12px;
            margin-top: 20px;
            text-align: left;
            font-size: 14px;
            color: #856404;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #e2e8f0;
            color: #a0aec0;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">🚀 Orbia App</div>
        <div class="subtitle">Your Gateway to Web3 Influence</div>
        
        <div class="content">
            <h1>Email Verification Code</h1>
            <p class="info">You are attempting to sign in to your Orbia account. Please use the verification code below:</p>
            
            <div class="code-container">
                <div class="code">{{.Code}}</div>
            </div>
            
            <p class="info">This verification code will expire in <strong>{{.ExpireMinutes}} minutes</strong>.</p>
            
            <div class="warning">
                <strong>⚠️ Security Notice:</strong><br>
                If you did not request this code, please ignore this email. Do not share this code with anyone.
            </div>
            
            <div class="footer">
                <p>This is an automated message from Orbia App.</p>
                <p>© 2025 Orbia. All rights reserved.</p>
            </div>
        </div>
    </div>
</body>
</html>
`
)

// GenerateVerificationCode 生成指定长度的数字验证码
func GenerateVerificationCode(length int) string {
	if length <= 0 {
		length = 6 // 默认6位
	}

	code := ""
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += fmt.Sprintf("%d", n)
	}
	return code
}

// SendVerificationEmail 发送验证码邮件
func SendVerificationEmail(to, code string, expireMinutes int) error {
	cfg := config.GlobalConfig.SMTP

	log.Printf("[Email Debug] Starting to send verification email")
	log.Printf("[Email Debug] SMTP Config - Server: %s, Port: %s, Username: %s, Email: %s, FromName: %s",
		cfg.Server, cfg.Port, cfg.Username, cfg.Email, cfg.FromName)
	log.Printf("[Email Debug] Recipient: %s, Code: %s, ExpireMinutes: %d", to, code, expireMinutes)

	if cfg.Server == "" || cfg.Port == "" {
		return fmt.Errorf("SMTP configuration is not set")
	}

	// SMTP认证
	log.Printf("[Email Debug] Creating SMTP auth with username: %s, server: %s", cfg.Username, cfg.Server)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Server)

	// 准备模板数据
	data := struct {
		From          string
		To            string
		Code          string
		ExpireMinutes int
	}{
		From:          cfg.Email,
		To:            to,
		Code:          code,
		ExpireMinutes: expireMinutes,
	}

	// Parse template
	log.Printf("[Email Debug] Parsing email template")
	tmpl, err := template.New("email").Parse(EmailVerificationCodeTemplate)
	if err != nil {
		log.Printf("[Email Debug] Failed to parse template: %v", err)
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Execute template
	log.Printf("[Email Debug] Executing template with data")
	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
		log.Printf("[Email Debug] Failed to execute template: %v", err)
		return fmt.Errorf("error executing template: %w", err)
	}

	log.Printf("[Email Debug] Email body size: %d bytes", emailBody.Len())
	log.Printf("[Email Debug] First 300 chars of email: %s", string(emailBody.Bytes()[:min(300, emailBody.Len())]))

	// 发送邮件（完全按照另一个项目的方式）
	addr := fmt.Sprintf("%s:%s", cfg.Server, cfg.Port)
	log.Printf("[Email Debug] Sending email to SMTP server: %s", addr)
	log.Printf("[Email Debug] From: %s, To: %s", cfg.Email, to)

	err = smtp.SendMail(
		addr,
		auth,
		cfg.Email,
		[]string{to},
		emailBody.Bytes(),
	)
	if err != nil {
		log.Printf("[Email Debug] Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("[Email Debug] Email sent successfully!")
	return nil
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	// 简单的邮箱格式验证
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}
