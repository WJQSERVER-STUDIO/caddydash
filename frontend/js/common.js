// js/common.js - 存放共享模块

import { initI18n, t } from './locale.js';
import { notification } from './notifications.js';
import { initCaddyStatus } from './caddy.js';
import { initUI } from './ui.js';

const theme = {
    init: (toggleElement) => {
        const storedTheme = localStorage.getItem('theme');
        const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        const currentTheme = storedTheme || (systemPrefersDark ? 'dark' : 'light');
        theme.apply(currentTheme);
        if (toggleElement) {
            toggleElement.addEventListener('change', (e) => theme.apply(e.target.checked ? 'dark' : 'light'));
        }
    },
    apply: (themeName) => {
        document.documentElement.dataset.theme = themeName;
        localStorage.setItem('theme', themeName);
        const themeToggleInput = document.getElementById('theme-toggle-input');
        if (themeToggleInput) {
            themeToggleInput.checked = themeName === 'dark';
        }
    }
};

function hideToast(toastElement) {
    if (!toastElement) return;
    toastElement.classList.remove('show');
    toastElement.addEventListener('transitionend', () => toastElement.remove(), { once: true });
}

const toast = {
    show: (message, type = 'info', duration = 3000) => {
        const toastContainer = document.getElementById('toast-container');
        if (!toastContainer) return;
        const icons = { success: 'fa-check-circle', error: 'fa-times-circle', info: 'fa-info-circle' };
        const iconClass = icons[type] || 'fa-info-circle';
        const toastElement = document.createElement('div');
        toastElement.className = `toast ${type}`;
        toastElement.innerHTML = `<i class="toast-icon fa-solid ${iconClass}"></i><p class="toast-message">${message}</p><button class="toast-close" data-toast-close>×</button>`;
        toastContainer.appendChild(toastElement);
        requestAnimationFrame(() => toastElement.classList.add('show'));
        const timeoutId = setTimeout(() => hideToast(toastElement), duration);
        toastElement.querySelector('[data-toast-close]').addEventListener('click', () => {
            clearTimeout(timeoutId);
            hideToast(toastElement);
        });
    }
};

function activateNav(pageId) {
    const navLinks = document.querySelectorAll('.sidebar-nav a');
    navLinks.forEach(link => {
        link.classList.remove('active');
        if (link.dataset.navId === pageId) {
            link.classList.add('active');
        }
    });
}

/**
 * 初始化移动端侧边栏的开关逻辑
 */
function initSidebar() {
    const sidebar = document.getElementById('sidebar');
    const menuToggleBtn = document.getElementById('menu-toggle-btn');
    
    // 动态创建并管理遮罩层
    let overlay = document.querySelector('.sidebar-overlay');
    if (!overlay) {
        overlay = document.createElement('div');
        overlay.className = 'sidebar-overlay';
        document.body.appendChild(overlay);
    }

    const openSidebar = () => {
        if (!sidebar || !overlay) return;
        sidebar.classList.add('is-open');
        overlay.classList.add('is-visible');
    };

    const closeSidebar = () => {
        if (!sidebar || !overlay) return;
        sidebar.classList.remove('is-open');
        overlay.classList.remove('is-visible');
    };

    if (menuToggleBtn) {
        menuToggleBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            openSidebar();
        });
    }

    overlay.addEventListener('click', closeSidebar);
}

/**
 * 通用的页面初始化函数
 * @param {object} options - 初始化选项
 * @param {string} options.pageId - 当前页面的ID, 用于导航高亮
 * @param {function} [options.pageInit=null] - 特定于该页面的额外初始化逻辑
 */
export async function initializePage(options) {
    // 1. 初始化国际化 (最高优先级)
    await initI18n();

    initUI(t);

    // 2. 初始化UI模块和通用功能
    theme.init(document.getElementById('theme-toggle-input'));
    notification.init(
        document.getElementById('toast-container'), 
        document.getElementById('dialog-container'),
        document.getElementById('modal-container')
    );
    activateNav(options.pageId);
    initSidebar();
    initCaddyStatus(t);

    // 3. 如果有特定页面的初始化逻辑, 则执行它
    if (options.pageInit && typeof options.pageInit === 'function') {
        options.pageInit();
    }
}

// 导出模块
export { theme, toast, activateNav };