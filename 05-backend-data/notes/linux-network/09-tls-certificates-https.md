# TLS、证书与 HTTPS

## 是什么

TLS 在传输层上提供机密性、完整性和对端认证；证书把公钥与域名等身份绑定并由信任链验证。HTTPS 是 HTTP over TLS。

## 为什么需要

没有 TLS，凭据和数据可被路径上的参与者读取或修改，也无法可靠确认服务端身份。

## 关键特性或规则

客户端验证证书链、有效期、主机名和用途；SNI/ALPN 在握手协商域名和协议；私钥最小权限并轮换；自动续期仍需告警。

## 实际怎么使用

```sh
openssl s_client -connect example.com:443 -servername example.com -showcerts
curl -v https://example.com/health
```

## 常见错误与边界

加密不代表应用安全；忽略证书验证会失去身份认证；证书更新后旧连接仍可能使用旧会话；时间错误会导致验证失败。

## 补充知识

TLS 1.3 减少握手往返并移除旧算法；mTLS 可认证客户端但增加证书生命周期复杂度。

## 来源

- [一手资料 1](https://www.rfc-editor.org/rfc/rfc8446.html)（访问日期：2026-07-16）
- [一手资料 2](https://www.rfc-editor.org/rfc/rfc6125.html)（访问日期：2026-07-16）
- [一手资料 3](https://developer.mozilla.org/en-US/docs/Glossary/HTTPS)（访问日期：2026-07-16）
- [一手资料 4](https://docs.openssl.org/master/man1/openssl-s_client/)（访问日期：2026-07-16）
