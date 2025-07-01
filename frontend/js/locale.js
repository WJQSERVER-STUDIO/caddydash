// js/locale.js - 国际化 (i18n) 核心模块

let currentLocale = {};
let currentLang = 'en'; // 默认语言
const supportedLangs = ['en', 'zh-CN']; // 应用支持的语言列表

/**
 * 加载指定的语言文件 (JSON)
 * @param {string} lang - 语言代码 (e.g., 'en', 'zh-CN')
 */
async function loadLocale(lang) {
    try {
        const response = await fetch(`/locales/${lang}.json`);
        if (!response.ok) {
            throw new Error(`Language file for ${lang} not found (status: ${response.status}).`);
        }
        currentLocale = await response.json();
        currentLang = lang;
        document.documentElement.lang = lang;
    } catch (error) {
        console.error("i18n Error:", error);
        // 如果加载目标语言失败, 安全回退到默认的英语
        if (lang !== 'en') {
            console.warn(`Falling back to default language 'en'.`);
            await loadLocale('en');
        }
    }
}

/**
 * 将加载的翻译应用到所有带有 data-i18n 属性的DOM元素上
 */
function applyTranslationsToDOM() {
    document.querySelectorAll('[data-i18n]').forEach(el => {
        const key = el.dataset.i18n;
        const translation = t(key);
        if (translation !== key) {
            // 优先替换元素的第一个文本节点, 避免覆盖内部的 <i> 等元素
            const textNode = Array.from(el.childNodes).find(node => node.nodeType === Node.TEXT_NODE && node.textContent.trim());
            if (textNode) {
                // 在图标和文本之间保留一个空格
                textNode.textContent = el.querySelector('i') ? ` ${translation}` : translation;
            } else {
                el.textContent = translation;
            }
        }
    });
    // 特殊处理 title 属性
    document.querySelectorAll('[data-i18n-title]').forEach(el => {
        const key = el.dataset.i18nTitle;
        const translation = t(key);
        if (translation !== key) el.title = translation;
    });
    document.querySelectorAll('[data-i18n-placeholder]').forEach(el => {
        const key = el.dataset.i18nPlaceholder;
        const translation = t(key);
        if (translation !== key && el.placeholder !== undefined) { // 确保元素有placeholder属性
            el.placeholder = translation;
        }
    });
}

/**
 * 获取翻译文本, 支持点分隔的路径和占位符替换
 * @param {string} key - 翻译键 (e.g., 'pages.login.welcome')
 * @param {object} [replacements={}] - 用于替换占位符的键值对
 * @returns {string} - 翻译后的字符串
 */
export function t(key, replacements = {}) {
    // 通过路径 'a.b.c' 在嵌套对象中查找值: currentLocale['a']['b']['c']
    const translation = key.split('.').reduce((obj, k) => obj && obj[k], currentLocale);

    let result = translation || key; // 如果找不到, 返回原始key作为回退

    // 处理占位符替换, e.g., {filename: 'example.com'}
    if (typeof result === 'string') {
        for (const placeholder in replacements) {
            result = result.replace(`{${placeholder}}`, replacements[placeholder]);
        }
    }

    return result;
}

/**
 * 初始化 i18n 系统: 检测语言, 加载语言包, 并应用翻译
 */
export async function initI18n() {
    const urlParams = new URLSearchParams(window.location.search);
    const langFromUrl = urlParams.get('lang');
    const langFromStorage = localStorage.getItem('appLanguage');
    const browserLang = navigator.language.startsWith('zh') ? 'zh-CN' : 'en';
    let langToLoad = 'en';

    if (langFromUrl && supportedLangs.includes(langFromUrl)) {
        langToLoad = langFromUrl;
    } else if (langFromStorage && supportedLangs.includes(langFromStorage)) {
        langToLoad = langFromStorage;
    } else if (supportedLangs.includes(browserLang)) {
        langToLoad = browserLang;
    }

    await loadLocale(langToLoad);
    applyTranslationsToDOM();
}

/**
 * 切换应用语言
 * @param {string} lang - 目标语言代码
 */
export async function setLanguage(lang) {
    if (supportedLangs.includes(lang) && lang !== currentLang) {
        localStorage.setItem('appLanguage', lang);
        window.location.reload(); // 刷新页面以应用所有翻译是最简单可靠的方式
    }
}

export function getCurrentLanguage() {
    return currentLang;
}