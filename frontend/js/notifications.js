// js/notifications.js - 提供Toast和Dialog两种通知

let toastContainer;
let dialogContainer;
let modalContainer; // 虽然在此文件中不直接使用, 但 init 中保留以示完整
let t; // 模块级翻译函数变量

function hideToast(toastElement) {
    if (!toastElement) return;
    toastElement.classList.remove('show');
    toastElement.addEventListener('transitionend', () => toastElement.remove(), { once: true });
}

export const notification = {
    init: (toastEl, dialogEl, modalEl, translator) => {
        toastContainer = toastEl;
        dialogContainer = dialogEl;
        modalContainer = modalEl;
        t = translator; // 保存从外部传入的翻译函数
    },
    toast: (message, type = 'info', duration = 3000) => {
        if (!toastContainer) return;
        const icons = { success: 'fa-check-circle', error: 'fa-times-circle', info: 'fa-info-circle', warning: 'fa-exclamation-triangle' };
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
    },
    confirm: (message, title = '', options = {}) => {
        return new Promise(resolve => {
            if (!dialogContainer || !t) {
                // 如果模块未初始化, 提供一个浏览器默认的 confirm作为回退
                console.warn('Notification module not initialized. Falling back to native confirm.');
                resolve(window.confirm(message));
                return;
            }

            // 使用 t 函数翻译按钮文本, 如果 options 中提供了自定义键, 则优先使用
            const confirmText = options.confirmText || t('dialogs.confirm_btn');
            const cancelText = options.cancelText || t('dialogs.cancel_btn');

            const dialogHTML = `
                <div class="dialog-box">
                    ${title ? `<h3>${title}</h3>` : ''}
                    <p class="dialog-message">${message}</p>
                    <div class="dialog-actions">
                        <button class="btn btn-secondary" data-action="cancel">${cancelText}</button>
                        <button class="btn btn-primary" data-action="confirm">${confirmText}</button>
                    </div>
                </div>`;
            
            dialogContainer.innerHTML = dialogHTML;
            dialogContainer.classList.add('active');

            const eventHandler = (e) => {
                const actionButton = e.target.closest('[data-action]');
                if (!actionButton) return;
                closeDialog(actionButton.dataset.action === 'confirm');
            };

            const closeDialog = (result) => {
                dialogContainer.removeEventListener('click', eventHandler);
                dialogContainer.classList.remove('active');
                setTimeout(() => { dialogContainer.innerHTML = ''; resolve(result); }, 200);
            };

            dialogContainer.addEventListener('click', eventHandler);
        });
    }
};