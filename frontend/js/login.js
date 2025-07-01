// js/login.js - 登录页面的独立逻辑

document.addEventListener('DOMContentLoaded', () => {

    const DOMElements = {
        loginForm: document.getElementById('login-form'),
        toastContainer: document.getElementById('toast-container'),
        langSwitcherBtn: document.getElementById('lang-switcher-btn'),
        langOptionsList: document.getElementById('lang-options-list'), // 从第一个片段引入
    };
    const loginButton = DOMElements.loginForm.querySelector('button[type="submit"]');
    const LOGIN_API_URL = '/v0/api/auth/login';

    const i18n = {
        currentLocale: {},
        currentLang: 'en',
        // 从第一个片段引入, 使用对象更方便显示语言名称
        supportedLangs: { 'en': 'English', 'zh-CN': '简体中文' }, 
        t: function(key, replacements = {}) {
            const translation = key.split('.').reduce((obj, k) => obj && obj[k], this.currentLocale) || key;
            let result = translation;
            if (typeof result === 'string') {
                for (const placeholder in replacements) {
                    result = result.replace(`{${placeholder}}`, replacements[placeholder]);
                }
            }
            return result;
        },
        applyTranslations: function() {
            // 优化后的翻译应用逻辑, 优先更新span, 其次更新非空文本节点, 最后直接更新元素文本
            document.querySelectorAll('[data-i18n]').forEach(el => {
                const key = el.dataset.i18n;
                const translation = this.t(key);
                if (translation !== key) { // 仅当找到翻译时才应用
                    const spanChild = el.querySelector('span');
                    if (spanChild) {
                        spanChild.textContent = translation;
                    } else {
                        // 查找直接的、非空文本节点进行替换
                        const textNode = Array.from(el.childNodes).find(node => node.nodeType === Node.TEXT_NODE && node.textContent.trim().length > 0);
                        if (textNode) {
                            textNode.textContent = translation;
                        } else {
                            // 备用方案: 直接设置元素的textContent
                            el.textContent = translation;
                        }
                    }
                }
            });
            // 从第一个片段引入, 处理data-i18n-title属性
            document.querySelectorAll('[data-i18n-title]').forEach(el => {
                el.title = this.t(el.dataset.i18nTitle);
            });
            document.title = this.t('pages.login.page_title');
        },
        loadLocale: async function(lang) {
            try {
                const response = await fetch(`/locales/${lang}.json`);
                if (!response.ok) throw new Error('File not found');
                this.currentLocale = await response.json();
                this.currentLang = lang;
                document.documentElement.lang = lang; // 设置HTML语言属性
                localStorage.setItem('appLanguage', lang); // 从第一个片段引入, 保存到localStorage
            } catch (e) {
                console.error(`Could not load locale for ${lang}, using fallback.`, e);
                this.currentLocale = {}; 
            }
        },
        init: async function() {
            // 从第一个片段引入, 优先使用保存的语言, 其次使用浏览器语言
            const savedLang = localStorage.getItem('appLanguage');
            const browserLang = navigator.language.startsWith('zh') ? 'zh-CN' : 'en';
            const langToLoad = savedLang || browserLang;
            await this.loadLocale(langToLoad);
            this.applyTranslations();
            this.populateLangOptions(); // 从第一个片段引入, 初始化语言选项列表
        },
        // 从第一个片段引入, 用于动态生成语言选项列表
        populateLangOptions: function() {
            // 清空现有选项
            DOMElements.langOptionsList.innerHTML = ''; 
            for (const [code, name] of Object.entries(this.supportedLangs)) {
                const li = document.createElement('li');
                li.dataset.lang = code;
                li.textContent = name;
                if (code === this.currentLang) {
                    li.classList.add('active'); // 标记当前选中语言
                }
                DOMElements.langOptionsList.appendChild(li);
            }
        }
        // 移除 i18n.toggleLanguage, 因为有新的语言选择机制
    };
    
    // 从第二个片段完整引入toast对象
    const toast = {
        show: function(message, type = 'info', duration = 3000) {
            if (!DOMElements.toastContainer) return;
            const icons = { success: 'fa-check-circle', error: 'fa-times-circle', info: 'fa-info-circle' };
            const toastElement = document.createElement('div');
            toastElement.className = `toast ${type}`;
            toastElement.innerHTML = `<i class="toast-icon fa-solid ${icons[type]}"></i><p class="toast-message">${message}</p><button class="toast-close" data-toast-close>×</button>`;
            DOMElements.toastContainer.appendChild(toastElement);
            requestAnimationFrame(() => toastElement.classList.add('show'));
            const timeoutId = setTimeout(() => this._hide(toastElement), duration);
            toastElement.querySelector('[data-toast-close]').addEventListener('click', () => {
                clearTimeout(timeoutId);
                this._hide(toastElement);
            });
        },
        _hide: function(toastElement) {
            if (!toastElement) return;
            toastElement.classList.remove('show');
            toastElement.addEventListener('transitionend', () => toastElement.remove(), { once: true });
        }
    };

    // 从第二个片段完整引入handleLogin函数
    async function handleLogin(e) {
        e.preventDefault();
        const username = DOMElements.loginForm.username.value.trim();
        const password = DOMElements.loginForm.password.value.trim();

        if (username === '') {
            toast.show(i18n.t('toasts.error_username_empty'), 'error');
            DOMElements.loginForm.username.focus();
            return;
        }
        if (password === '') {
            toast.show(i18n.t('toasts.error_password_empty'), 'error');
            DOMElements.loginForm.password.focus();
            return;
        }

        loginButton.disabled = true;
        loginButton.querySelector('span').textContent = i18n.t('pages.login.logging_in_btn');

        try {
            const response = await fetch(LOGIN_API_URL, {
                method: 'POST',
                body: new URLSearchParams(new FormData(DOMElements.loginForm))
            });
            const result = await response.json();
            if (response.ok) {
                toast.show(i18n.t('toasts.login_success'), 'success');
                setTimeout(() => { window.location.href = '/'; }, 500);
            } else {
                throw new Error(result.error || i18n.t('toasts.login_error_generic'));
            }
        } catch (error) {
            toast.show(error.message, 'error');
            loginButton.disabled = false;
            loginButton.querySelector('span').textContent = i18n.t('pages.login.login_btn');
        }
    }

    async function initApp() {
        // 主题设置逻辑
        const storedTheme = localStorage.getItem('theme');
        const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        document.documentElement.dataset.theme = storedTheme || (systemPrefersDark ? 'dark' : 'light');
        
        // 初始化国际化
        await i18n.init();
        
        // 登录表单事件监听
        if (DOMElements.loginForm) {
            DOMElements.loginForm.addEventListener('submit', handleLogin);
        }
        // 语言切换按钮事件监听 (从第一个片段引入)
        if (DOMElements.langSwitcherBtn) {
            DOMElements.langSwitcherBtn.addEventListener('click', (e) => {
                e.stopPropagation(); // 阻止事件冒泡, 防止立即触发document的点击事件
                DOMElements.langOptionsList.classList.toggle('hidden');
            });
        }
        // 语言选项列表事件监听 (从第一个片段引入)
        if (DOMElements.langOptionsList) {
            DOMElements.langOptionsList.addEventListener('click', async (e) => {
                const target = e.target.closest('li[data-lang]'); // 查找最近的语言li元素
                if (target) {
                    await i18n.loadLocale(target.dataset.lang); // 加载新语言
                    i18n.applyTranslations(); // 应用翻译
                    i18n.populateLangOptions(); // 更新语言选项列表的激活状态
                    DOMElements.langOptionsList.classList.add('hidden'); // 隐藏列表
                }
            });
        }
        // 文档点击事件, 用于点击外部时隐藏语言选项列表 (从第一个片段引入)
        document.addEventListener('click', () => {
            if (DOMElements.langOptionsList && !DOMElements.langOptionsList.classList.contains('hidden')) {
                DOMElements.langOptionsList.classList.add('hidden');
            }
        });
    }

    initApp();
});