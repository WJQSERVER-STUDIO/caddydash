# CaddyDash - Modern Caddy Web Administration Panel

CaddyDash is a modern, user-friendly web interface for managing and configuring your Caddy server. It provides an intuitive graphical interface, allowing you to easily manage site configurations, adjust global settings, and monitor the status of your Caddy instances.

[简体中文](README.md) | English

## ✨ Features

*   **Intuitive Site Configuration**: Easily create, edit, and delete site configurations through forms, supporting reverse proxy and file serving modes.
*   **Global Caddyfile Management**: Directly modify and save Caddy's global configuration file within the web interface.
*   **Caddy Instance Control**: One-click start, stop, and reload of the Caddy service.
*   **Multi-language Support**: Built-in internationalization (i18n) functionality, supporting English and Chinese switching, with future expandability for more languages.
*   **Responsive Design**: The interface adapts to desktop and mobile devices, providing a consistent user experience.
*   **Theme Switching**: Supports light and dark themes to suit different user preferences.
*   **User Authentication**: Provides secure login and initialization processes to ensure panel access security.
*   **Preset Management**: Supports populating common request header configurations from presets, improving configuration efficiency.

## 🚀 Tech Stack

**Frontend:**

*   **Pure Native HTML5/CSS3/JavaScript (ESM)**: Built entirely on native browser technologies, responsive design, and mobile-friendly.

**Backend:**

*   **Go Language**: High-performance, concurrency-friendly backend service.
*   **Touka Framework**: HTTP framework built on Go, used for handling web requests.
*   **SQLite**: Lightweight embedded database for storing user and configuration data.
*   **CaddyServer**: The core component.

## 💡 Architecture Overview

CaddyDash frontend adopts a **Multi-Page Application (MPA)** architecture, where each main functional module corresponds to an independent HTML page, driven by its dedicated JavaScript entry file.

*   **Highly Modular**: All JavaScript code is organized in **ESM (ECMAScript Modules)** form, enabling code reuse and separation of concerns through `import/export` mechanisms.
*   **Shared Components**: Modules like `js/common.js`, `js/locale.js`, `js/notifications.js`, `js/ui.js`, `js/api.js` encapsulate functionalities shared across pages, such as page initialization, internationalization, notifications, UI operations, and backend API calls.
*   **Independent Page Logic**: `js/app.js` (site configuration), `js/global.js` (global configuration), `js/settings.js` (panel settings), `js/login.js` (login), `js/init.js` (initialization) handle their respective page-specific business logic.

## 🌐 Internationalization (i18n)

CaddyDash frontend supports multi-language display.

*   **Language Packs**: Translation texts are stored in `locales/en.json` (English) and `locales/zh-CN.json` (Simplified Chinese) files.
*   **Dynamic Translation**: The `js/locale.js` module is responsible for loading the correct language pack and dynamically translating texts within HTML elements that have `data-i18n` (content), `data-i18n-title` (title attribute), and `data-i18n-placeholder` (input placeholder) attributes.
*   **Switching Languages**: On the login page, initialization page, or panel settings page, you can change the language through the language toggle options in the interface.

**How to add new translation entries:**

1.  Add new key-value pairs to `locales/en.json` and `locales/zh-CN.json`. Please follow the existing dot-separated naming convention (e.g., `pages.feature.new_text`).
2.  In HTML, use `data-i18n="your.new.key"`, `data-i18n-title="your.new.key"`, or `data-i18n-placeholder="your.new.key"` attributes.
3.  In JavaScript code, use the `t('your.new.key', { replacements })` function to retrieve translated text.

## 🤝 Contributing

We welcome and encourage contributions of all forms! If you have any feature suggestions, bug reports, or code improvements, please feel free to submit them via Issues or Pull Requests.

## 📜 License

Copyright © 2025 WJQSERVER

This project CaddyDash is licensed under the **Mozilla Public License 2.0 (MPL 2.0)**.