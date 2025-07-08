// js/ui.js - 管理所有与UI渲染和DOM操作相关的函数

// 模块级私有变量, 用于存储翻译函数
let t;

// 新增: 初始化函数, 用于接收翻译函数
export function initUI(translator) {
    t = translator;
}

export const DOMElements = {
    sidebar: document.getElementById('sidebar'),
    menuToggleBtn: document.getElementById('menu-toggle-btn'),
    mainContent: document.querySelector('.main-content'),
    configListPanel: document.getElementById('config-list-panel'),
    configFormPanel: document.getElementById('config-form-panel'),
    renderedOutputPanel: document.getElementById('rendered-output-panel'),
    configForm: document.getElementById('config-form'),
    formTitle: document.getElementById('form-title'),
    backToListBtn: document.getElementById('back-to-list-btn'),
    domainInput: document.getElementById('domain'),
    originalFilenameInput: document.getElementById('original-filename'),
    headersContainer: document.getElementById('headers-container'),
    addNewConfigBtn: document.getElementById('add-new-config-btn'),
    cancelEditBtn: document.getElementById('cancel-edit-btn'),
    configListContainer: document.getElementById('config-list'),
    renderedContentCode: document.getElementById('rendered-content'),
    toastContainer: document.getElementById('toast-container'),
    dialogContainer: document.getElementById('dialog-container'),
    modalContainer: document.getElementById('modal-container'),
    themeToggleInput: document.getElementById('theme-toggle-input'),
    caddyStatusIndicator: document.getElementById('caddy-status-indicator'),
    caddyActionButtonContainer: document.getElementById('caddy-action-button-container'),
    logoutBtn: document.getElementById('logout-btn'),
    serviceModeControl: document.getElementById('service-mode-control'),
    upstreamFieldset: document.getElementById('upstream-fieldset'),
    fileserverFieldset: document.getElementById('fileserver-fieldset'),
    upstreamHeadersContainer: document.getElementById('upstream-headers-container'),
    mutiUpstreamCheckbox: document.getElementById('muti_upstream'),
    singleUpstreamGroup: document.getElementById('single-upstream-group'),
    multiUpstreamGroup: document.getElementById('multi-upstream-group'),
    multiUpstreamContainer: document.getElementById('multi-upstream-container'),
    addMultiUpstreamBtn: document.getElementById('add-multi-upstream-btn'),
};

export function switchView(viewToShow) {
    [DOMElements.configListPanel, DOMElements.configFormPanel, DOMElements.renderedOutputPanel]
        .forEach(view => view.classList.add('hidden'));
    if (viewToShow) viewToShow.classList.remove('hidden');
    DOMElements.addNewConfigBtn.disabled = (viewToShow === DOMElements.configFormPanel);
}

export function renderConfigList(filenames) {
    DOMElements.configListContainer.innerHTML = '';
    if (!filenames || filenames.length === 0) {
        DOMElements.configListContainer.innerHTML = `<p>${t('configs.no_configs')}</p>`;
        return;
    }
    filenames.forEach(filename => {
        const item = document.createElement('li');
        item.className = 'config-item';
        item.dataset.filename = filename;
        item.innerHTML = `<span class="config-item-name">${filename}</span><div class="config-item-actions"><button class="btn-icon edit-btn" title="${t('common.edit')}"><i class="fa-solid fa-pen-to-square"></i></button><button class="btn-icon delete-btn" title="${t('common.delete')}"><i class="fa-solid fa-trash-can"></i></button></div>`;
        DOMElements.configListContainer.appendChild(item);
    });
}

export function addKeyValueInput(container, keyName, valueName, key = '', value = '') {
    const div = document.createElement('div');
    div.className = 'header-entry';
    div.innerHTML = `
        <input type="text" name="${keyName}" placeholder="${t('form.key_placeholder')}" value="${key}">
        <input type="text" name="${valueName}" placeholder="${t('form.value_placeholder')}" value="${value}">
        <button type="button" class="btn-icon" onclick="this.parentElement.remove()" title="${t('common.remove_item')}">
            <i class="fa-solid fa-xmark"></i>
        </button>`;
    container.appendChild(div);
}

export function addSingleInput(container, inputName, placeholderKey, value = '') {
    const div = document.createElement('div');
    div.className = 'header-entry';
    div.style.gridTemplateColumns = '1fr auto';
    div.innerHTML = `
        <input type="text" name="${inputName}" placeholder="${t(placeholderKey)}" value="${value}">
        <button type="button" class="btn-icon" onclick="this.parentElement.remove()" title="${t('common.remove_upstream')}">
            <i class="fa-solid fa-xmark"></i>
        </button>`;
    container.appendChild(div);
}

export function fillForm(config, originalFilename) {
    DOMElements.originalFilenameInput.value = originalFilename;
    DOMElements.domainInput.value = config.domain_config?.domain || '';
    const enableUpstream = config.upstream_config?.enable_upstream || false;
    const enableFileServer = config.file_server_config?.enable_file_server || false;
    let mode = 'none';
    if (enableUpstream) mode = 'reverse_proxy';
    else if (enableFileServer) mode = 'file_server';
    const activeButton = DOMElements.serviceModeControl.querySelector(`[data-mode="${mode}"]`);
    if (activeButton) updateSegmentedControl(activeButton);
    const upstreamConfig = config.upstream_config || {};
    DOMElements.mutiUpstreamCheckbox.checked = upstreamConfig.muti_upstream || false;
    document.getElementById('upstream').value = upstreamConfig.upstream || '';
    DOMElements.multiUpstreamContainer.innerHTML = '';
    if (upstreamConfig.muti_upstream && upstreamConfig.upstream_servers) {
        upstreamConfig.upstream_servers.forEach(server => {
            addSingleInput(DOMElements.multiUpstreamContainer, 'upstream_servers', 'form.upstream_server_placeholder', server);
        });
    }
    DOMElements.upstreamHeadersContainer.innerHTML = '';
    if (upstreamConfig.upstream_headers) {
        Object.entries(upstreamConfig.upstream_headers).forEach(([k, v]) => v.forEach(val => addKeyValueInput(DOMElements.upstreamHeadersContainer, 'upstream_header_key', 'upstream_header_value', k, val)));
    }
    document.getElementById('file_dir_path').value = config.file_server_config?.file_dir_path || '';
    document.getElementById('enable_browser').checked = config.file_server_config?.enable_browser || false;
    DOMElements.headersContainer.innerHTML = '';
    if (config.headers) Object.entries(config.headers).forEach(([k, v]) => v.forEach(val => addKeyValueInput(DOMElements.headersContainer, 'header_key', 'header_value', k, val)));
    document.getElementById('enable_log').checked = config.log_config?.enable_log || false;
    document.getElementById('enable_error_page').checked = config.error_page_config?.enable_error_page || false;
    document.getElementById('enable_encode').checked = config.encode_config?.enable_encode || false;
    updateMultiUpstreamView(DOMElements.mutiUpstreamCheckbox.checked);
}

export function showRenderedConfig(configs, filename) {
    const targetConfig = configs.find(c => c.filename === filename);
    if (targetConfig && targetConfig.rendered_content) {
        DOMElements.renderedContentCode.textContent = atob(targetConfig.rendered_content);
        DOMElements.renderedOutputPanel.classList.remove('hidden');
    } else {
        DOMElements.renderedOutputPanel.classList.add('hidden');
    }
}

function createButton(text, className, onClick) {
    const button = document.createElement('button');
    button.className = `btn ${className}`;
    button.innerHTML = `<span>${text}</span>`;
    button.addEventListener('click', onClick);
    return button;
}

export function updateCaddyStatusView(status, handlers) {
    const { handleReloadCaddy, handleStopCaddy, handleStartCaddy } = handlers;
    const dot = DOMElements.caddyStatusIndicator.querySelector('.status-dot');
    const text = DOMElements.caddyStatusIndicator.querySelector('.status-text');
    const buttonContainer = DOMElements.caddyActionButtonContainer;
    if (!dot || !text || !buttonContainer) return;
    dot.className = 'status-dot';
    buttonContainer.innerHTML = '';
    let statusText, dotClass;
    switch (status) {
        case 'running':
            statusText = t('status.running'); dotClass = 'running';
            buttonContainer.appendChild(createButton(t('caddy.reload_btn'), 'btn-warning', handleReloadCaddy));
            buttonContainer.appendChild(createButton(t('caddy.stop_btn'), 'btn-danger', handleStopCaddy));
            break;
        case 'stopped':
            statusText = t('status.stopped'); dotClass = 'stopped';
            buttonContainer.appendChild(createButton(t('caddy.start_btn'), 'btn-success', handleStartCaddy));
            break;
        case 'checking': statusText = t('status.checking'); dotClass = 'checking'; break;
        default: statusText = t('status.unknown'); dotClass = 'error'; break;
    }
    text.textContent = statusText;
    dot.classList.add(dotClass);
}

export function updateServiceModeView(mode) {
    DOMElements.upstreamFieldset.classList.toggle('hidden', mode !== 'reverse_proxy');
    DOMElements.fileserverFieldset.classList.toggle('hidden', mode !== 'file_server');
}

export function updateMultiUpstreamView(isMulti) {
    DOMElements.singleUpstreamGroup.classList.toggle('hidden', isMulti);
    DOMElements.multiUpstreamGroup.classList.toggle('hidden', !isMulti);
}

export function updateSegmentedControl(activeButton) {
    const slider = document.getElementById('segmented-control-slider');
    const control = DOMElements.serviceModeControl;
    if (!activeButton || !slider || !control) return;
    control.querySelectorAll('button').forEach(btn => btn.classList.remove('active'));
    activeButton.classList.add('active');
    slider.style.width = `${activeButton.offsetWidth}px`;
    slider.style.transform = `translateX(${activeButton.offsetLeft}px)`;
}

export function createPresetSelectionModal(presets) {
    return new Promise(resolve => {
        const modalContainer = DOMElements.modalContainer;
        if (!modalContainer) return resolve(null);

       const presetItems = presets.map(p => {
            // 根据规则构建 i18n 键
            const nameKey = `presets.${p.id}.name`;
            const descKey = `presets.${p.id}.description`;

            // 使用 t() 函数进行翻译.
            //    - 如果找到了翻译, t() 会返回翻译后的文本.
            //    - 如果找不到, t() 会原样返回 nameKey.
            const displayName = t(nameKey);
            const displayDesc = t(descKey);

            // 设计回退逻辑:
            //    - 如果翻译结果和 key 相同, 说明没找到翻译, 此时使用后端返回的 p.name.
            //    - 否则, 使用翻译结果.
            const finalName = (displayName === nameKey) ? p.name : displayName;
            const finalDesc = (displayDesc === descKey) ? p.description : displayDesc;

            return `
                <li data-preset-id="${p.id}">
                    <strong>${finalName}</strong>
                    <p>${finalDesc}</p>
                </li>
            `;
        }).join('');

        /*
        const presetItems = presets.map(p => `
            <li data-preset-id="${p.id}">
                <strong>${t(p.name_key) || p.name}</strong>
                <p>${t(p.desc_key) || p.description}</p>
            </li>
        `).join('');
        */

        const modalHTML = `
            <div class="modal-overlay"></div>
            <div class="modal-box">
                <header class="modal-header">
                    <h3>${t('form.fill_from_preset')}</h3>
                    <button class="btn-icon" data-modal-close><i class="fa-solid fa-xmark"></i></button>
                </header>
                <div class="modal-content">
                    <ul class="preset-list">${presetItems}</ul>
                </div>
            </div>
        `;
        modalContainer.innerHTML = modalHTML;
        requestAnimationFrame(() => modalContainer.classList.add('active'));
        const cleanupAndResolve = (value) => {
            modalContainer.removeEventListener('click', eventHandler);
            modalContainer.classList.remove('active');
            setTimeout(() => { modalContainer.innerHTML = ''; resolve(value); }, 300);
        };
        const eventHandler = (e) => {
            if (e.target.classList.contains('modal-overlay') || e.target.closest('[data-modal-close]')) {
                cleanupAndResolve(null);
            }
            const listItem = e.target.closest('li[data-preset-id]');
            if (listItem) {
                cleanupAndResolve(listItem.dataset.presetId);
            }
        };
        modalContainer.addEventListener('click', eventHandler);
    });
}

export function createCustomSelect(containerId, options, onSelect) {
    const container = document.getElementById(containerId);
    if (!container) return;
    
    const inputName = container.id.replace('select-', '').replace(/-/g, '_');
    container.innerHTML = `<div class="select-selected"></div><div class="select-items"></div><input type="hidden" name="${inputName}">`;
    
    const selectedDiv = container.querySelector('.select-selected');
    const itemsDiv = container.querySelector('.select-items');
    const hiddenInput = container.querySelector('input[type="hidden"]');
    itemsDiv.innerHTML = '';

    if (!options || options.length === 0) {
        selectedDiv.textContent = t('common.no_options');
        return;
    }
    
    options.forEach((option, index) => {
        const item = document.createElement('div');
        const optionText = typeof option === 'object' ? option.name : option;
        const optionValue = typeof option === 'object' ? option.value : option;
        
        item.textContent = optionText;
        item.dataset.value = optionValue;

        if (index === 0) {
            selectedDiv.textContent = optionText;
            hiddenInput.value = optionValue;
        }
        item.addEventListener('click', function(e) {
            selectedDiv.textContent = this.textContent;
            hiddenInput.value = this.dataset.value;
            itemsDiv.classList.remove('select-show');
            selectedDiv.classList.remove('select-arrow-active');
            onSelect && onSelect(this.dataset.value);
            e.stopPropagation();
        });
        itemsDiv.appendChild(item);
    });

    selectedDiv.addEventListener('click', (e) => {
        e.stopPropagation();
        document.querySelectorAll('.select-items.select-show').forEach(openSelect => {
            if (openSelect !== itemsDiv) {
                openSelect.classList.remove('select-show');
                openSelect.previousElementSibling.classList.remove('select-arrow-active');
            }
        });
        itemsDiv.classList.toggle('select-show');
        selectedDiv.classList.toggle('select-arrow-active');
    });

    document.addEventListener('click', () => {
        itemsDiv.classList.remove('select-show');
        if(selectedDiv) selectedDiv.classList.remove('select-arrow-active');
    });
}