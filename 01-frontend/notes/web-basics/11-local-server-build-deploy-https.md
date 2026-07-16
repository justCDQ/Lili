# 本地开发服务器、构建产物、部署与 HTTPS

## 是什么与为什么需要

本地开发服务器通过 HTTP 提供项目文件，常带热更新和开发错误信息。构建把源文件转换、打包或优化为部署产物。部署把产物和服务配置发布到可访问环境。HTTPS 是 HTTP 经 TLS 提供传输加密、完整性与服务器身份验证。

直接双击 `file:` 页面缺少真实 HTTP 环境，模块、请求和路径行为可能不同；本地服务器更接近生产访问方式。

## 实际使用

```sh
npx serve .
# 或项目脚本
npm run dev
npm run build
```

开发时确认终端输出的 host/port；构建后检查 `dist/`，用静态服务器预览而非直接打开。部署时只上传约定产物，配置正确入口、404 回退、缓存和 HTTPS，最后用公开 URL 检查 Network 与证书。

## 关键规则

- 开发服务器通常不适合公网生产：性能、安全和错误输出配置不同。
- 构建产物是生成文件，应由同一源码和锁文件重建；是否提交由项目约定。
- HTTPS 不保证站点业务可信，只保护客户端到证书对应服务器之间的传输。
- HTTPS 页面加载 HTTP 子资源会形成 mixed content，并可能被浏览器阻止。

## 常见错误与边界

部署子路径时绝对 `/assets/...` 可能指向域名根。单页应用直接刷新深层 URL 需要服务器回退配置。环境变量嵌入前端构建后对用户可见。开发环境缓存和生产缓存策略不同。

## 补充知识

TLS 证书需覆盖域名并在有效期内；自动续期仍需监控。静态托管、CDN、容器和服务器都是部署方式，选择取决于是否需要动态执行环境。

## 来源

- [MDN：Publishing your website](https://developer.mozilla.org/en-US/docs/Learn_web_development/Getting_started/Your_first_website/Publishing_your_website)
- [MDN：HTTPS](https://developer.mozilla.org/en-US/docs/Glossary/HTTPS)
- [MDN：Mixed content](https://developer.mozilla.org/en-US/docs/Web/Security/Mixed_content)

访问日期：2026-07-16。
