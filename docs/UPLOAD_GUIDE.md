# Cloudflare R2 上传指南

本指南说明如何使用 Orbia API 的上传服务将文件上传到 Cloudflare R2 存储。

## 目录
- [概述](#概述)
- [工作原理](#工作原理)
- [接口说明](#接口说明)
- [前端使用示例](#前端使用示例)
- [安全性说明](#安全性说明)

## 概述

我们使用 **预签名 URL（Presigned URL）** 方式实现安全的文件上传。这是 Cloudflare R2 和 AWS S3 的最佳实践，具有以下优势：

- ✅ **安全性高**：不在前端暴露 R2 访问凭证
- ✅ **性能好**：前端直接上传到 R2，不经过后端服务器
- ✅ **灵活性强**：支持设置过期时间、文件大小限制等
- ✅ **简单易用**：前端只需要使用标准的 HTTP PUT 请求

## 工作原理

```
┌─────────┐                          ┌─────────────┐                      ┌──────────────┐
│         │  1. 请求上传 token        │             │                      │              │
│  前端   │ ────────────────────────> │  后端 API   │                      │ Cloudflare   │
│         │                          │             │                      │      R2      │
│         │  2. 返回预签名 URL        │             │                      │              │
│         │ <──────────────────────── │             │                      │              │
│         │                          └─────────────┘                      │              │
│         │                                                                │              │
│         │  3. 使用预签名 URL 直接上传文件                                 │              │
│         │ ──────────────────────────────────────────────────────────────> │              │
│         │                                                                │              │
│         │  4. 返回上传结果                                                │              │
│         │ <────────────────────────────────────────────────────────────── │              │
└─────────┘                                                                └──────────────┘
```

## 接口说明

### 生成上传 Token

**接口地址：** `POST /api/v1/upload/token`

**请求参数：**

```json
{
  "image_type": 1,           // 图片类型：1=头像(AVATAR), 2=团队图标(TEAM_ICON)
  "file_extension": ".jpg",  // 文件扩展名（包含点号）
  "file_size": 1024000       // 文件大小（字节，可选）
}
```

**响应参数：**

```json
{
  "upload_url": "https://xxx.r2.cloudflarestorage.com/bucket/path?签名参数...",
  "public_url": "https://pub-xxx.r2.dev/avatars/xxx.jpg",
  "expires_in": 1800,  // 过期时间（秒）
  "headers": {
    "Content-Type": "image/jpeg",
    "Content-Length": "1024000"
  },
  "base_resp": {
    "code": 0,
    "message": "success"
  }
}
```

**字段说明：**

- `upload_url`: 预签名的上传 URL，使用 PUT 方法上传
- `public_url`: 上传成功后的公开访问 URL
- `expires_in`: URL 有效期（秒），默认 30 分钟
- `headers`: 上传时必需的 HTTP 请求头

## 前端使用示例

### JavaScript (Fetch API)

```javascript
// 1. 获取上传 token
async function getUploadToken(file) {
  const response = await fetch('/api/v1/upload/token', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer YOUR_JWT_TOKEN'
    },
    body: JSON.stringify({
      image_type: 1, // AVATAR
      file_extension: `.${file.name.split('.').pop()}`,
      file_size: file.size
    })
  });
  
  const data = await response.json();
  if (data.base_resp.code !== 0) {
    throw new Error(data.base_resp.message);
  }
  
  return data;
}

// 2. 上传文件到 R2
async function uploadFile(file, uploadToken) {
  const response = await fetch(uploadToken.upload_url, {
    method: 'PUT',
    headers: uploadToken.headers,
    body: file
  });
  
  if (!response.ok) {
    throw new Error(`Upload failed: ${response.statusText}`);
  }
  
  return uploadToken.public_url;
}

// 3. 完整的上传流程
async function uploadImage(file) {
  try {
    // 验证文件
    if (!file.type.startsWith('image/')) {
      throw new Error('只支持图片文件');
    }
    
    if (file.size > 10 * 1024 * 1024) { // 10MB
      throw new Error('文件大小不能超过 10MB');
    }
    
    // 获取上传凭证
    console.log('正在获取上传凭证...');
    const uploadToken = await getUploadToken(file);
    
    // 上传文件
    console.log('正在上传文件...');
    const publicUrl = await uploadFile(file, uploadToken);
    
    console.log('上传成功！图片地址：', publicUrl);
    return publicUrl;
    
  } catch (error) {
    console.error('上传失败：', error.message);
    throw error;
  }
}

// 使用示例
document.getElementById('fileInput').addEventListener('change', async (e) => {
  const file = e.target.files[0];
  if (file) {
    try {
      const imageUrl = await uploadImage(file);
      console.log('图片已上传，URL:', imageUrl);
      // 可以将 imageUrl 保存到表单或发送到后端
    } catch (error) {
      alert('上传失败：' + error.message);
    }
  }
});
```

### React 示例

```typescript
import { useState } from 'react';

interface UploadToken {
  upload_url: string;
  public_url: string;
  expires_in: number;
  headers: Record<string, string>;
}

function ImageUploader() {
  const [uploading, setUploading] = useState(false);
  const [imageUrl, setImageUrl] = useState<string>('');

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      setUploading(true);

      // 1. 获取上传 token
      const tokenResponse = await fetch('/api/v1/upload/token', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          image_type: 1,
          file_extension: `.${file.name.split('.').pop()}`,
          file_size: file.size
        })
      });

      const tokenData = await tokenResponse.json();
      if (tokenData.base_resp.code !== 0) {
        throw new Error(tokenData.base_resp.message);
      }

      const uploadToken: UploadToken = tokenData;

      // 2. 上传到 R2
      const uploadResponse = await fetch(uploadToken.upload_url, {
        method: 'PUT',
        headers: uploadToken.headers,
        body: file
      });

      if (!uploadResponse.ok) {
        throw new Error('Upload failed');
      }

      // 3. 设置图片 URL
      setImageUrl(uploadToken.public_url);
      alert('上传成功！');

    } catch (error) {
      console.error('Upload error:', error);
      alert(`上传失败：${error.message}`);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div>
      <input
        type="file"
        accept="image/*"
        onChange={handleUpload}
        disabled={uploading}
      />
      {uploading && <p>上传中...</p>}
      {imageUrl && (
        <div>
          <p>上传成功！</p>
          <img src={imageUrl} alt="Uploaded" style={{ maxWidth: '300px' }} />
          <p>图片地址：{imageUrl}</p>
        </div>
      )}
    </div>
  );
}

export default ImageUploader;
```

### Axios 示例

```javascript
import axios from 'axios';

async function uploadImageWithAxios(file) {
  try {
    // 1. 获取上传 token
    const tokenResponse = await axios.post('/api/v1/upload/token', {
      image_type: 1,
      file_extension: `.${file.name.split('.').pop()}`,
      file_size: file.size
    }, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    const { upload_url, public_url, headers } = tokenResponse.data;

    // 2. 上传文件到 R2
    await axios.put(upload_url, file, {
      headers: headers
    });

    console.log('上传成功！图片地址：', public_url);
    return public_url;

  } catch (error) {
    console.error('上传失败：', error);
    throw error;
  }
}
```

## 安全性说明

### 我们实施的安全措施

1. **不暴露凭证**
   - 接口返回的是预签名 URL，不包含 R2 的 Access Key 和 Secret Key
   - 前端无法获取到任何敏感凭证信息

2. **限制上传范围**
   - 预签名 URL 只能上传到指定的路径
   - 文件路径由后端生成，前端无法自定义

3. **时间限制**
   - 预签名 URL 有 30 分钟有效期
   - 过期后自动失效，需要重新获取

4. **文件类型和大小限制**
   - 后端验证文件扩展名（只支持图片格式）
   - 限制文件大小不超过 10MB
   - 预签名 URL 中包含 Content-Type 和 Content-Length 限制

5. **用户认证**
   - 必须登录才能获取上传 token
   - 防止匿名上传

### 前端注意事项

1. **不要缓存上传 token**
   - 每次上传都重新获取 token
   - 不要存储或重复使用 upload_url

2. **验证文件**
   ```javascript
   // 验证文件类型
   if (!file.type.startsWith('image/')) {
     throw new Error('只支持图片文件');
   }
   
   // 验证文件大小
   if (file.size > 10 * 1024 * 1024) {
     throw new Error('文件大小不能超过 10MB');
   }
   ```

3. **错误处理**
   ```javascript
   try {
     await uploadFile(file, token);
   } catch (error) {
     if (error.response?.status === 403) {
       console.error('上传 URL 已过期，请重新获取');
     } else if (error.response?.status === 413) {
       console.error('文件太大');
     } else {
       console.error('上传失败：', error.message);
     }
   }
   ```

## 支持的文件格式

目前支持以下图片格式：
- `.jpg` / `.jpeg`
- `.png`
- `.gif`
- `.webp`

## 常见问题

### Q: 为什么使用 PUT 而不是 POST？
A: 预签名 URL 生成时使用的是 PutObject 操作，对应 HTTP PUT 方法。这是 S3 兼容 API 的标准做法。

### Q: 可以批量上传吗？
A: 可以。为每个文件分别获取 token，然后并发或串行上传。

### Q: 上传失败怎么办？
A: 
- 检查 token 是否过期（30分钟有效期）
- 检查请求头是否正确设置
- 检查文件大小是否超限
- 查看浏览器控制台的错误信息

### Q: 可以上传到 CDN 吗？
A: `public_url` 已经是 CDN 地址，可以直接使用。Cloudflare R2 自带 CDN 加速。

### Q: 如何显示上传进度？
A: 使用 XMLHttpRequest 或 Axios 的进度回调：

```javascript
const xhr = new XMLHttpRequest();
xhr.upload.addEventListener('progress', (e) => {
  if (e.lengthComputable) {
    const percent = (e.loaded / e.total) * 100;
    console.log(`上传进度: ${percent.toFixed(2)}%`);
  }
});

xhr.open('PUT', uploadToken.upload_url);
Object.entries(uploadToken.headers).forEach(([key, value]) => {
  xhr.setRequestHeader(key, value);
});
xhr.send(file);
```

## 参考资料

- [Cloudflare R2 文档](https://developers.cloudflare.com/r2/)
- [AWS S3 预签名 URL 文档](https://docs.aws.amazon.com/AmazonS3/latest/userguide/PresignedUrlUploadObject.html)

