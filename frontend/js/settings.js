// js/settings.js - 设置页面的逻辑

import { initializePage } from './common.js';
import { api } from './api.js';
import { notification } from './notifications.js';
import { t, setLanguage, getCurrentLanguage } from './locale.js';
import { createCustomSelect } from './ui.js';

const DOMElements = {
    resetForm: document.getElementById('reset-password-form'),
    logoutBtn: document.getElementById('logout-btn'),
};
const resetButton = DOMElements.resetForm.querySelector('button[type="submit"]');

async function handleResetPassword(e) {
    e.preventDefault();
    const newPassword = DOMElements.resetForm.new_password.value;
    const confirmPassword = DOMElements.resetForm.confirm_new_password.value;
    const currentPassword = DOMElements.resetForm.old_password.value;
    const username = DOMElements.resetForm.username.value;

    if (!username || !currentPassword || !newPassword || !confirmPassword) {
        notification.toast(t('toasts.error_all_fields_required'), 'error');
        return;
    }
    if (newPassword !== confirmPassword) {
        notification.toast(t('toasts.init_error_mismatch'), 'error');
        return;
    }
    if (newPassword.length < 8) {
        notification.toast(t('toasts.init_error_short'), 'error');
        return;
    }

    resetButton.disabled = true;
    resetButton.querySelector('span').textContent = t('pages.settings.resetting_password_btn');

    try {
        const result = await api.post('/auth/resetpwd', new URLSearchParams(new FormData(DOMElements.resetForm)));
        notification.toast(t('toasts.pwd_reset_success'), 'success');
        setTimeout(() => { window.location.href = '/v0/api/auth/logout'; }, 1500);
    } catch (error) {
        notification.toast(`${t('common.error_prefix')}: ${error.message}`, 'error');
        resetButton.disabled = false;
        resetButton.querySelector('span').textContent = t('pages.settings.reset_password_btn');
    }
}

async function handleLogout() {
    if (await notification.confirm(t('dialogs.logout_msg'))) {
        notification.toast(t('toasts.logout_processing'), 'info');
        setTimeout(() => { window.location.href = '/v0/api/auth/logout'; }, 500);
    }
}

// 页面特有的初始化逻辑
function pageInit() {
    const langOptions = { 'en': 'English', 'zh-CN': '简体中文' };
    const langSelectOptions = Object.keys(langOptions).map(key => ({ name: langOptions[key], value: key }));

    createCustomSelect('select-language', langSelectOptions, (selectedValue) => {
        setLanguage(selectedValue);
    });

    const langSelect = document.getElementById('select-language');
    if (langSelect) {
        const currentLangName = langOptions[getCurrentLanguage()];
        const selectedDiv = langSelect.querySelector('.select-selected');
        if (selectedDiv) selectedDiv.textContent = currentLangName;
    }

    DOMElements.resetForm.addEventListener('submit', handleResetPassword);
    DOMElements.logoutBtn.addEventListener('click', handleLogout);
}

// 使用通用初始化函数
initializePage({ pageId: 'settings', pageInit: pageInit });