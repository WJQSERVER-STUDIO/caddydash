// js/caddy.js - Caddy 实例状态管理与控制

import { api } from './api.js';
import { notification } from './notifications.js';

let caddyStatusInterval;
const POLLING_INTERVAL = 5000;

// 将 t 函数保存在模块作用域内
let translate; 

const DOMElements = {
    caddyStatusIndicator: document.getElementById('caddy-status-indicator'),
    caddyActionButtonContainer: document.getElementById('caddy-action-button-container'),
};

function createButton(text, className, onClick) {
    const button = document.createElement('button');
    button.className = `btn ${className}`;
    button.innerHTML = `<span>${text}</span>`;
    button.addEventListener('click', onClick);
    return button;
}

function updateCaddyStatusView(status) {
    const dot = DOMElements.caddyStatusIndicator.querySelector('.status-dot');
    const text = DOMElements.caddyStatusIndicator.querySelector('.status-text');
    const buttonContainer = DOMElements.caddyActionButtonContainer;
    
    if(!dot || !text || !buttonContainer) return;

    dot.className = 'status-dot';
    buttonContainer.innerHTML = '';
    let statusText, dotClass;
    switch (status) {
        case 'running':
            statusText = translate('status.running'); 
            dotClass = 'running';
            buttonContainer.appendChild(createButton(translate('caddy.reload_btn'), 'btn-warning', handleReloadCaddy));
            buttonContainer.appendChild(createButton(translate('caddy.stop_btn'), 'btn-danger', handleStopCaddy));
            break;
        case 'stopped':
            statusText = translate('status.stopped'); 
            dotClass = 'stopped';
            buttonContainer.appendChild(createButton(translate('caddy.start_btn'), 'btn-success', handleStartCaddy));
            break;
        case 'checking': 
            statusText = translate('status.checking'); 
            dotClass = 'checking'; 
            break;
        default: 
            statusText = translate('status.unknown'); 
            dotClass = 'error'; 
            break;
    }
    text.textContent = statusText;
    dot.classList.add(dotClass);
}

async function checkCaddyStatus() {
    try {
        const response = await api.get('/caddy/status');
        updateCaddyStatusView(response.message === 'Caddy is running' ? 'running' : 'stopped');
    } catch (error) { 
        console.error('Error checking Caddy status:', error); 
        updateCaddyStatusView('error');
    }
}

async function handleStartCaddy() {
    try {
        const result = await api.post('/caddy/run');
        notification.toast(result.message || translate('toasts.start_cmd_sent'), 'success');
        setTimeout(checkCaddyStatus, 500);
    } catch (error) { notification.toast(translate('toasts.start_error', { error: error.message }), 'error'); }
}

async function handleStopCaddy() {
    if (!await notification.confirm(translate('dialogs.stop_caddy_msg'))) return;
    try {
        const result = await api.post('/caddy/stop');
        notification.toast(result.message || translate('toasts.stop_cmd_sent'), 'info');
        setTimeout(checkCaddyStatus, 500);
    } catch(error) { notification.toast(translate('toasts.action_error', { error: error.message }), 'error'); }
}

async function handleReloadCaddy() {
    if (!await notification.confirm(translate('dialogs.reload_caddy_msg'))) return;
    try {
        const result = await api.post('/caddy/restart');
        notification.toast(result.message || translate('toasts.reload_sent'), 'success');
        setTimeout(checkCaddyStatus, 500);
    } catch(error) { notification.toast(translate('toasts.reload_error', { error: error.message }), 'error'); }
}

// initCaddyStatus 现在接收 t 函数作为参数
export function initCaddyStatus(translator) {
    // 保存翻译函数以供模块内其他函数使用
    translate = translator;

    const dialogContainer = document.getElementById('dialog-container');
    const toastContainer = document.getElementById('toast-container');
    if (dialogContainer && toastContainer) {
        notification.init(toastContainer, dialogContainer, null, translate); // 将 t 函数传递给通知模块
    }
    
    checkCaddyStatus();
    if (caddyStatusInterval) {
        clearInterval(caddyStatusInterval);
    }
    caddyStatusInterval = setInterval(checkCaddyStatus, POLLING_INTERVAL);
}