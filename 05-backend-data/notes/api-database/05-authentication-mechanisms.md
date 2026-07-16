# Session、Cookie、JWT、OAuth 2、OIDC 与 API Key

## 是什么

session 通常用随机标识关联服务端状态；Cookie 是浏览器状态传输机制；JWT 是可签名 claims 容器；OAuth 2 是委托授权框架；OIDC 在 OAuth 2 上定义身份认证；API key 标识调用应用。

## 为什么需要

混淆认证、授权和令牌格式会产生可伪造身份、泄漏凭据和错误信任边界。

## 关键特性或规则

session ID 高熵、服务端可撤销并轮换；Cookie 设置 Secure/HttpOnly/SameSite；JWT 验证 alg、签名、iss、aud、exp；OAuth 浏览器/原生应用用授权码+PKCE；OIDC 验证 nonce/state。

## 实际怎么使用

```go
cookie:=&http.Cookie{Name:"session",Value:randomID,Path:"/",HttpOnly:true,Secure:true,SameSite:http.SameSiteLaxMode}
http.SetCookie(w,cookie)
// Authorization: Bearer <access-token>
```

## 常见错误与边界

JWT payload 不是加密；不要接受算法混淆；OAuth access token 不是登录身份；API key 不代表最终用户；localStorage token 暴露给 XSS。

## 补充知识

密钥和令牌需要生命周期、轮换、撤销、最小 scope 与审计；优先成熟身份提供方。

## 来源

- [一手资料 1](https://www.rfc-editor.org/rfc/rfc6265.html)（访问日期：2026-07-16）
- [一手资料 2](https://www.rfc-editor.org/rfc/rfc7519.html)（访问日期：2026-07-16）
- [一手资料 3](https://www.rfc-editor.org/rfc/rfc6749.html)（访问日期：2026-07-16）
- [一手资料 4](https://openid.net/specs/openid-connect-core-1_0.html)（访问日期：2026-07-16）
- [一手资料 5](https://www.rfc-editor.org/rfc/rfc7636.html)（访问日期：2026-07-16）
