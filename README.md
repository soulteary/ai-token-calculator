# AI 模型 Token 计算器

一个可以验证和计算文本消耗 Token 的小工具。

网页版本，汉化自 OpenAI Tokenizer。

## 使用方式

使用方式有很多种。

TBD

### Nginx

你可以使用 Nginx 来快速使用这个项目，如果你使用 Docker，可以使用类似 [./docker-compose.nginx.yml](./docker-compose.nginx.yml) 中的方式：

```yaml
version: "3"

services:
  web:
    image: nginx:1.25.3-alpine
    volumes:
      - ./public:/usr/share/nginx/html
    ports:
      - "8080:80"
```

