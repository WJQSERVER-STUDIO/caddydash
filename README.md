# CaddyDash - 现代化 Caddy Web 管理面板

CaddyDash 是一个现代化、用户友好的 Web 界面，用于管理和配置您的 Caddy 服务器。它提供了一个直观的图形界面，让您能够轻松地管理站点配置、调整全局设置，并监控 Caddy 实例的状态。

[English](README_EN.md) | 简体中文

## ✨ 特性

*   **直观的站点配置**: 通过表单轻松创建、编辑和删除站点配置，支持反向代理和文件服务模式。
*   **全局 Caddyfile 管理**: 直接在 Web 界面中修改和保存 Caddy 的全局配置文件。
*   **Caddy 实例控制**: 一键启动、停止和重载 Caddy 服务。
*   **多语言支持**: 内置国际化 (i18n) 功能，支持中英文切换，未来可扩展更多语言。
*   **响应式设计**: 界面适配桌面和移动设备，提供一致的用户体验。
*   **主题切换**: 支持明亮与暗色主题，满足不同用户偏好。
*   **用户认证**: 提供安全的登录和初始化流程，保障面板访问安全。
*   **预设管理**: 支持从预设填充常用的请求头配置，提高配置效率。

## 🚀 技术栈

**前端:**

*   **纯原生 HTML5/CSS3/JavaScript (ESM)**: 完全基于浏览器原生技术构建，响应式设计, 移动端友好

**后端:**

*   **Go 语言**: 高性能、并发友好的后端服务
*   **Touka 框架**: [Touka](https://github.com/infinite-iroha/touka)高效多层次的HTTP框架
*   **SQLite**: 轻量级嵌入式数据库，用于存储用户和配置数据
*   **CaddyServer**: 作为核心组件

## 💡 架构概览

CaddyDash 前端采用**多页面应用 (MPA)** 架构，每个主要功能模块都对应一个独立的 HTML 页面，并由其专属的 JavaScript 入口文件驱动。

*   **高度模块化**: 所有 JavaScript 代码都以 **ESM (ECMAScript Modules)** 形式组织，通过 `import/export` 机制实现代码复用和职责分离。
*   **共享组件**: `js/common.js`、`js/locale.js`、`js/notifications.js`、`js/ui.js`、`js/api.js` 等模块封装了跨页面共享的功能，如页面初始化、国际化、通知、UI操作和后端 API 调用。
*   **独立页面逻辑**: `js/app.js` (站点配置), `js/global.js` (全局配置), `js/settings.js` (面板设置), `js/login.js` (登录), `js/init.js` (初始化) 分别处理各自页面的特定业务逻辑。

## 部署

推荐使用docker compose部署

```yml
version: '3.8'

services:
  caddydash:
    image: wjqserver/caddydash:latest
    container_name: caddydash
    network_mode: host # 使用宿主机网络模式

    # 映射宿主机目录到容器内部，用于持久化数据和配置
    volumes:
      # CaddyDash 的配置目录，包含 config.toml
      - ./config:/data/caddy/config
      # Caddy 的额外配置文件片段目录
      - ./config.d:/data/caddy/config.d
      # 日志目录
      - ./log:/data/caddy/log
      # CaddyDash 的数据库文件
      - ./db:/data/caddy/db
      # Caddy 自身的内部数据目录（如证书存储），Caddy Dash 可能会管理此处
      - ./caddy_internal_data:/root/.local/share/caddy # Caddy 默认数据目录之一
    
    restart: unless-stopped
```

## 🌐 国际化 (i18n)

CaddyDash 前端支持多语言显示。

*   **语言包**: 翻译文本存储在 `locales/en.json` (英文) 和 `locales/zh-CN.json` (简体中文) 文件中。
*   **动态翻译**: `js/locale.js` 模块负责加载正确的语言包，并动态地将 HTML 元素中带有 `data-i18n` (内容)、`data-i18n-title` (标题属性) 和 `data-i18n-placeholder` (输入框提示) 的文本进行翻译。
*   **切换语言**: 在登录页、初始化页或面板设置页，您可以通过界面上的语言切换选项来改变语言。

**如何添加新的翻译条目:**

1.  在 `locales/en.json` 和 `locales/zh-CN.json` 中添加新的键值对。请遵循现有的点分隔命名约定（例如 `pages.feature.new_text`）。
2.  在 HTML 中使用 `data-i18n="your.new.key"`、`data-i18n-title="your.new.key"` 或 `data-i18n-placeholder="your.new.key"` 属性。
3.  在 JavaScript 代码中，使用 `t('your.new.key', { replacements })` 函数来获取翻译文本。

## 相关项目

### 核心

[Caddy官方仓库](https://github.com/caddyserver/caddy)

### 后端

[Touka框架](https://github.com/infinite-iroha/touka)
[SQLite驱动](https://modernc.org/sqlite)
[Touka中间件...](https://github.com/fenthope)

## 关于

*   [![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/WJQSERVER-STUDIO/caddydash) 可供参考, AI生成存在幻觉, 不完全可靠, 请注意辨别

## 🤝 贡献

我们欢迎并鼓励任何形式的贡献！如果您有任何功能建议、bug 报告或代码改进，请随时通过 Issues 或 Pull Requests 提交。

## 📜 许可证

Copyright © 2025 WJQSERVER

本项目 CaddyDash 在 **Mozilla Public License 2.0 (MPL 2.0)** 许可证下授权