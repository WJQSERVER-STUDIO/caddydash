// js/global.js - 全局配置页面的逻辑

import { initializePage } from './common.js'; // 导入通用初始化函数
import { api } from './api.js';
import { notification } from './notifications.js';
import { createCustomSelect } from './ui.js';

const DOMElements = {
    globalForm: document.getElementById('global-caddy-form'),
    logoutBtn: document.getElementById('logout-btn'),
    enableDnsChallengeCheckbox: document.getElementById('enable_dns_challenge'),
    globalTlsConfigGroup: document.getElementById('global-tls-config-group'),
    enableEchCheckbox: document.getElementById('enable_ech'),
    echConfigGroup: document.getElementById('ech-config-group'),
};
// submitButton 在 pageInit 中获取, 确保DOM已加载
let submitButton;

function getGlobalConfigFromForm() {
    const formData = new FormData(DOMElements.globalForm);
    const enableEch = DOMElements.enableEchCheckbox.checked;

    return {
        debug: DOMElements.globalForm.querySelector('[name="debug"]').checked,
        ports_config: {
            admin_port: formData.get('admin_port'),
            http_port: parseInt(formData.get('http_port'), 10) || 80,
            https_port: parseInt(formData.get('https_port'), 10) || 443,
        },
        metrics: DOMElements.globalForm.querySelector('[name="metrics"]').checked,
        log_config: {
            level: DOMElements.globalForm.querySelector('input[name="log_level"]').value,
            rotate_size: formData.get('log_rotate_size'),
            rotate_keep: formData.get('log_rotate_keep'),
            rotate_keep_for_time: formData.get('log_rotate_keep_for_time'),
        },
        tls_config: {
            enable_dns_challenge: DOMElements.enableDnsChallengeCheckbox.checked,
            provider: DOMElements.globalForm.querySelector('input[name="tls_provider"]').value,
            token: formData.get('tls_token'),
            echouter_sni: enableEch ? formData.get('tls_ech_sni') : "",
            email: formData.get('tls_email'),
        },
        tls_snippet_config: {},
    };
}

function fillGlobalConfigForm(config) {
    if (!config) return;

    DOMElements.globalForm.querySelector('[name="debug"]').checked = config.debug || false;
    DOMElements.globalForm.querySelector('[name="metrics"]').checked = config.metrics || false;

    const ports = config.ports_config || {};
    DOMElements.globalForm.querySelector('[name="admin_port"]').value = ports.admin_port || ':2019';
    DOMElements.globalForm.querySelector('[name="http_port"]').value = ports.http_port || 80;
    DOMElements.globalForm.querySelector('[name="https_port"]').value = ports.https_port || 443;

    const log = config.log_config || {};
    const logLevelSelect = document.getElementById('select-log-level');
    const logLevel = log.level || 'INFO';
    if (logLevelSelect && logLevelSelect.querySelector('.select-selected')) {
        logLevelSelect.querySelector('.select-selected').textContent = logLevel;
        const hiddenInput = logLevelSelect.querySelector('input[name="log_level"]');
        if (hiddenInput) hiddenInput.value = logLevel;
    }
    DOMElements.globalForm.querySelector('[name="log_rotate_size"]').value = log.rotate_size || '10MB';
    DOMElements.globalForm.querySelector('[name="log_rotate_keep"]').value = log.rotate_keep || '10';
    DOMElements.globalForm.querySelector('[name="log_rotate_keep_for_time"]').value = log.rotate_keep_for_time || '24h';

    const tls = config.tls_config || {};
    DOMElements.enableDnsChallengeCheckbox.checked = tls.enable_dns_challenge || false;
    DOMElements.globalTlsConfigGroup.classList.toggle('hidden', !DOMElements.enableDnsChallengeCheckbox.checked);

    const tlsProviderSelect = document.getElementById('select-tls-provider');
    const provider = tls.provider || '';
    if (tlsProviderSelect && provider && tlsProviderSelect.querySelector('.select-selected')) {
        tlsProviderSelect.querySelector('.select-selected').textContent = provider;
        const hiddenProviderInput = tlsProviderSelect.querySelector('input[name="tls_provider"]');
        if (hiddenProviderInput) hiddenProviderInput.value = provider;
    }
    DOMElements.globalForm.querySelector('[name="tls_token"]').value = tls.token || '';
    DOMElements.globalForm.querySelector('[name="tls_email"]').value = tls.email || '';

    const echOuterSni = tls.echouter_sni || '';
    DOMElements.enableEchCheckbox.checked = !!echOuterSni;
    DOMElements.echConfigGroup.classList.toggle('hidden', !DOMElements.enableEchCheckbox.checked);
    DOMElements.globalForm.querySelector('[name="tls_ech_sni"]').value = echOuterSni;
}

async function handleSaveGlobalConfig(e) {
    e.preventDefault();
    const configData = getGlobalConfigFromForm();
    submitButton.disabled = true;
    submitButton.querySelector('span').textContent = "保存中...";

    try {
        const result = await api.put('/global/config', configData);
        notification.toast(result.message || '全局配置已成功保存，Caddy正在重载...', 'success');
    } catch (error) {
        notification.toast(`保存失败: ${error.message}`, 'error');
    } finally {
        submitButton.disabled = false;
        submitButton.querySelector('span').textContent = "保存并重载";
    }
}

async function handleLogout() {
    if (await notification.confirm('您确定要退出登录吗?')) {
        notification.toast('正在退出...', 'info');
        setTimeout(() => { window.location.href = '/v0/api/auth/logout'; }, 500);
    }
}

// 页面特有的初始化逻辑
function pageInit() {
    // 在这里获取 submitButton, 确保 DOM 已加载
    submitButton = DOMElements.globalForm.querySelector('button[type="submit"]');

    api.get('/global/log/levels')
        .then(levels => createCustomSelect('select-log-level', Object.keys(levels)))
        .catch(err => notification.toast(`加载日志级别失败: ${err.message}`, 'error'));

    api.get('/global/tls/providers')
        .then(providers => createCustomSelect('select-tls-provider', Object.keys(providers)))
        .catch(err => notification.toast(`加载TLS提供商失败: ${err.message}`, 'error'));

    api.get('/global/config')
        .then(config => fillGlobalConfigForm(config))
        .catch(err => notification.toast(`加载全局配置失败: ${err.message}`, 'error'));

    DOMElements.globalForm.addEventListener('submit', handleSaveGlobalConfig);
    DOMElements.logoutBtn.addEventListener('click', handleLogout);

    DOMElements.enableDnsChallengeCheckbox.addEventListener('change', (e) => {
        DOMElements.globalTlsConfigGroup.classList.toggle('hidden', !e.target.checked);
    });
    DOMElements.enableEchCheckbox.addEventListener('change', (e) => {
        DOMElements.echConfigGroup.classList.toggle('hidden', !e.target.checked);
    });
}

// 使用通用初始化函数启动页面
initializePage({ pageId: 'global', pageInit: pageInit });