// js/theme.js - 处理深色/浅色主题切换

// 这个模块在初始化时需要知道切换按钮的DOM元素
let themeToggleInput;

function applyTheme(themeName) {
    document.documentElement.dataset.theme = themeName;
    if (themeToggleInput) {
        themeToggleInput.checked = themeName === 'dark';
    }
    localStorage.setItem('theme', themeName);
}

function handleToggle(e) {
    const newTheme = e.target.checked ? 'dark' : 'light';
    applyTheme(newTheme);
}

export function initTheme(toggleElement) {
    themeToggleInput = toggleElement;
    const storedTheme = localStorage.getItem('theme');
    const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const currentTheme = storedTheme || (systemPrefersDark ? 'dark' : 'light');
    applyTheme(currentTheme);
    themeToggleInput.addEventListener('change', handleToggle);
}