# 文件上传 API 文档

## 概述

本系统支持多种文件类型的上传，包括图片、PDF、Word 文档、视频等。通过配置文件可以灵活控制每种文件类型的大小限制和存储路径。

## 核心特性

- ✅ **自动识别文件类型**：根据文件扩展名自动识别和验证
- ✅ **灵活配置**：每种文件扩展名都可以单独配置大小限制和存储路径
- ✅ **安全上传**：使用预签名 URL，不暴露凭证
- ✅ **支持多种文件**：图片、PDF、Word、Excel、PowerPoint、视频等
- ✅ **统一扩展名**：自动将扩展名转为小写（.PNG → .png）
- ✅ **路径自动管理**：文件路径和名称完全由后端控制，保证安全性

## 支持的文件类型

### 图片类型
- `.jpg` / `.jpeg` - 最大 10MB
- `.png` - 最大 10MB
- `.gif` - 最大 5MB
- `.webp` - 最大 10MB

### 文档类型
- `.pdf` - 最大 50MB
- `.doc` / `.docx` - 最大 20MB
- `.xls` / `.xlsx` - 最大 20MB
- `.ppt` / `.pptx` - 最大 50MB
- `.txt` - 最大 1MB

### 视频类型
- `.mp4` - 最大 100MB
- `.mov` - 最大 100MB
- `.avi` - 最大 100MB

> **注意**：可以在 `conf/config.yaml` 中的 `r2.allowed_extensions` 配置新的文件类型。

## API 接口

### 1. 生成上传 Token

**接口地址**：`POST /api/v1/upload/token`

**请求头**：
```
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json
```

**请求参数**：

```json
{
  "file_extension": ".pdf",      // 必需：文件扩展名（会自动转为小写）
  "file_size": 1024000           // 可选：文件大小（字节）
}
```

**参数说明**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `file_extension` | string | 是 | 文件扩展名，如 `.jpg`, `.png`, `.pdf`，会自动转为小写 |
| `file_size` | int64 | 否 | 文件大小（字节），用于验证 |

**路径管理**：

文件路径和名称完全由后端控制，保证安全性和一致性：
- 扩展名自动转为小写（`.PNG` → `.png`）
- 根据扩展名自动选择存储目录（配置文件中的 `default_path`）
- 文件名格式：`{时间戳}_{随机数}.{扩展名}`
- 按年月自动分目录：`{目录}/{年月}/{文件名}`
- 示例：`images/2025/01/1737849600_abc123.png`

**响应示例**：

```json
{
  "upload_url": "https://xxx.r2.cloudflarestorage.com/orbia/documents/2025/01/1737849600_abc123.pdf?X-Amz-...",
  "public_url": "https://pub-xxx.r2.dev/documents/2025/01/1737849600_abc123.pdf",
  "expires_in": 1800,
  "headers": {
    "Content-Type": "application/pdf",
    "Content-Length": "1024000"
  },
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

### 2. 验证文件 URL

**接口地址**：`POST /api/v1/upload/validate`

**请求头**：
```
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json
```

**请求参数**：

```json
{
  "file_url": "https://pub-xxx.r2.dev/documents/2025/01/xxx.pdf"
}
```

**响应示例**：

```json
{
  "is_valid": true,
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

## 使用示例

### 示例 1：上传图片

```javascript
// 1. 获取上传 token
const tokenResponse = await fetch('/api/v1/upload/token', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    file_extension: '.png',  // 即使传 .PNG 也会自动转为 .png
    file_size: file.size
  })
});

const uploadToken = await tokenResponse.json();

// 2. 上传文件
await fetch(uploadToken.upload_url, {
  method: 'PUT',
  headers: uploadToken.headers,
  body: file
});

// 3. 使用 public_url
console.log('文件地址：', uploadToken.public_url);
// 输出示例：https://pub-xxx.r2.dev/images/2025/01/1737849600_abc123.png
```

### 示例 2：上传 PDF 文档

```javascript
const tokenResponse = await fetch('/api/v1/upload/token', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    file_extension: '.pdf',
    file_size: file.size
  })
});
// 文件会自动存储到 documents 目录（配置文件中的 default_path）
```

### 示例 3：上传视频

```javascript
const tokenResponse = await fetch('/api/v1/upload/token', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    file_extension: '.mp4',
    file_size: file.size
  })
});
// 文件会自动存储到 videos 目录
```

### 示例 4：React 完整示例

```typescript
import { useState } from 'react';

function FileUploader() {
  const [uploading, setUploading] = useState(false);
  const [fileUrl, setFileUrl] = useState('');

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      setUploading(true);

      // 获取文件扩展名
      const extension = '.' + file.name.split('.').pop();

      // 1. 获取上传 token
      const tokenResponse = await fetch('/api/v1/upload/token', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          file_extension: extension,
          file_size: file.size
        })
      });

      const tokenData = await tokenResponse.json();
      if (tokenData.base_resp.code !== 0) {
        throw new Error(tokenData.base_resp.message);
      }

      // 2. 上传到 R2
      await fetch(tokenData.upload_url, {
        method: 'PUT',
        headers: tokenData.headers,
        body: file
      });

      // 3. 保存文件 URL
      setFileUrl(tokenData.public_url);
      alert('上传成功！');

    } catch (error) {
      console.error('上传失败：', error);
      alert(`上传失败：${error.message}`);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div>
      <input
        type="file"
        onChange={handleUpload}
        disabled={uploading}
      />
      {uploading && <p>上传中...</p>}
      {fileUrl && (
        <div>
          <p>文件地址：{fileUrl}</p>
          <a href={fileUrl} target="_blank" rel="noopener noreferrer">
            查看文件
          </a>
        </div>
      )}
    </div>
  );
}
```

## 配置说明

在 `conf/config.yaml` 中配置文件类型：

```yaml
r2:
  # ... 其他配置 ...
  max_file_size: 10485760  # 默认最大文件大小（10MB）
  
  # 按文件扩展名配置
  allowed_extensions:
    .jpg:
      max_size: 10485760  # 10MB
      default_path: "images"
    .pdf:
      max_size: 52428800  # 50MB
      default_path: "documents"
    .mp4:
      max_size: 104857600  # 100MB
      default_path: "videos"
```

**添加新的文件类型**：

只需在配置文件中添加新的扩展名配置，无需修改代码：

```yaml
allowed_extensions:
  .zip:
    max_size: 104857600  # 100MB
    default_path: "attachments"
  .rar:
    max_size: 104857600  # 100MB
    default_path: "attachments"
```

## 错误码

| Code | Message | 说明 |
|------|---------|------|
| 400 | unsupported file format: .xxx | 不支持的文件格式 |
| 400 | file size exceeds limit | 文件大小超过限制 |
| 401 | unauthorized | 未登录或 token 无效 |
| 500 | failed to generate upload token | 生成上传 token 失败 |

## 常见问题

### Q1: 如何添加新的文件类型？

在 `conf/config.yaml` 中的 `r2.allowed_extensions` 添加新的扩展名配置即可，无需修改代码。

### Q2: 文件路径的组织结构是什么？

文件路径格式为：`{upload_path}/{年月}/{时间戳}_{随机数}.{扩展名}`

例如：
- `avatars/2025/01/1737849600_abc123.jpg`
- `documents/2025/01/1737849600_def456.pdf`

### Q3: 为什么要统一扩展名为小写？

1. **一致性**：避免同一个文件类型有多种表示（.jpg vs .JPG vs .Jpg）
2. **配置简单**：只需在配置文件中维护小写的扩展名
3. **URL 美观**：生成的文件 URL 更统一美观

### Q4: 如何限制特定用户上传特定类型的文件？

可以在业务逻辑层（service 层）添加额外的权限检查。目前接口只验证文件扩展名和大小。

### Q5: 上传失败怎么办？

检查以下几点：
1. 文件扩展名是否在配置的 `allowed_extensions` 中
2. 文件大小是否超过该扩展名的 `max_size` 限制
3. 用户是否已登录（需要有效的 JWT token）
4. 预签名 URL 是否已过期（30分钟有效期）

## 最佳实践

1. **前端验证**：在调用接口前，先在前端验证文件类型和大小，提升用户体验
2. **显示进度**：使用 XMLHttpRequest 或 Axios 显示上传进度
3. **错误处理**：妥善处理各种错误情况，给用户明确的提示
4. **文件名处理**：服务器会生成唯一的文件名，前端不需要处理文件名冲突
5. **CDN 加速**：`public_url` 已经是 CDN 地址，可以直接使用

## 安全性

- ✅ 使用预签名 URL，不暴露 R2 凭证
- ✅ 限制文件类型和大小
- ✅ 需要用户登录认证
- ✅ 预签名 URL 有 30 分钟过期时间
- ✅ 文件路径由服务器生成，防止路径遍历攻击

