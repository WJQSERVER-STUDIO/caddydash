// js/app.js - 主应用入口

import { state } from './state.js';
import { api } from './api.js';
import { initializePage } from './common.js';
import { t } from './locale.js';
import { notification } from './notifications.js';
import { DOMElements, switchView, renderConfigList, addKeyValueInput, addSingleInput, fillForm, showRenderedConfig, updateCaddyStatusView, updateSegmentedControl, updateServiceModeView, updateMultiUpstreamView, createPresetSelectionModal } from './ui.js';

// --- 事件处理与逻辑流 (所有 handle* 函数保持不变) ---

function getFormStateAsString() {
    const formData = new FormData(DOMElements.configForm);
    const data = {};
    const activeModeButton = DOMElements.serviceModeControl.querySelector('button.active');
    data.active_mode = activeModeButton ? activeModeButton.dataset.mode : 'none';
    for (const [key, value] of formData.entries()) {
        const el = DOMElements.configForm.querySelector(`[name="${key}"]`);
        if (el?.type === 'checkbox') data[key] = el.checked;
        else if (data[key]) {
            if (!Array.isArray(data[key])) data[key] = [data[key]];
            data[key].push(value);
        } else data[key] = value;
    }
    return JSON.stringify(data);
}

async function attemptExitForm() {
    if (await getFormStateAsString() !== state.initialFormState) {
        if (await notification.confirm(t('dialogs.unsaved_changes_msg'), t('dialogs.unsaved_changes_title'))) {
            switchView(DOMElements.configListPanel);
        }
    } else {
        switchView(DOMElements.configListPanel);
    }
}

async function handleLogout() {
    if (!await notification.confirm(t('dialogs.logout_msg'))) return;
    notification.toast(t('toasts.logout_processing'), 'info');
    setTimeout(() => { window.location.href = `/v0/api/auth/logout`; }, 500);
}

async function loadAllConfigs() {
    try {
        const filenames = await api.get('/config/filenames');
        renderConfigList(filenames);
    } catch (error) { if (error.message) notification.toast(t('toasts.load_configs_error', { error: error.message }), 'error'); }
}

async function handleEditConfig(originalFilename) {
    try {
        const [config, rendered] = await Promise.all([
            api.get(`/config/file/${originalFilename}`),
            api.get('/config/files/rendered')
        ]);
        state.isEditing = true;
        switchView(DOMElements.configFormPanel);
        DOMElements.formTitle.textContent = t('pages.configs.form_title_edit');
        fillForm(config, originalFilename);
        showRenderedConfig(rendered, originalFilename);
        const mode = config.upstream_config?.enable_upstream ? 'reverse_proxy' : (config.file_server_config?.enable_file_server ? 'file_server' : 'none');
        updateServiceModeView(mode);

        // Manually trigger change event for cache checkbox to set initial state
        const enableCacheCheckbox = DOMElements.configForm.querySelector('#enable_cache');
        enableCacheCheckbox.dispatchEvent(new Event('change'));

        state.initialFormState = await getFormStateAsString();
    } catch (error) { notification.toast(t('toasts.load_config_detail_error', { error: error.message }), 'error'); }
}

async function handleDeleteConfig(filename) {
    if (!await notification.confirm(t('dialogs.delete_config_msg', { filename }), t('dialogs.delete_config_title'))) return;
    try {
        await api.delete(`/config/file/${filename}`);
        notification.toast(t('toasts.delete_success'), 'success');
        loadAllConfigs();
    } catch (error) { notification.toast(t('toasts.delete_error', { error: error.message }), 'error'); }
}

async function handleSaveConfig(e) {
    e.preventDefault();
    const formData = new FormData(DOMElements.configForm);
    const domain = formData.get('domain');
    if (!domain) {
        notification.toast(t('toasts.error_domain_empty'), 'error');
        return;
    }
    const getHeadersMap = (keyName, valueName) => {
        const headers = {};
        formData.getAll(keyName).forEach((key, i) => {
            if (key) {
                if (!headers[key]) headers[key] = [];
                headers[key].push(formData.getAll(valueName)[i]);
            }
        });
        return headers;
    };
    const activeModeButton = DOMElements.serviceModeControl.querySelector('button.active');
    const activeMode = activeModeButton ? activeModeButton.dataset.mode : 'none';
    const isMutiUpstream = DOMElements.mutiUpstreamCheckbox.checked;

    const configData = {
        mode: "uni",
        domain_config: { domain: domain, muti_domains: false, domains: [] },
        upstream_config: {
            enable_upstream: activeMode === 'reverse_proxy',
            upstream: formData.get('upstream'),
            muti_upstream: isMutiUpstream,
            upstream_servers: formData.getAll('upstream_servers').filter(s => s),
            upstream_headers: getHeadersMap('upstream_header_key', 'upstream_header_value'),
        },
        file_server_config: {
            enable_file_server: activeMode === 'file_server',
            file_dir_path: formData.get('file_dir_path'),
            enable_browser: DOMElements.configForm.querySelector('#enable_browser').checked
        },
        headers: getHeadersMap('header_key', 'header_value'),
        log_config: { enable_log: DOMElements.configForm.querySelector('#enable_log').checked, log_domain: domain },
        error_page_config: { enable_error_page: DOMElements.configForm.querySelector('#enable_error_page').checked },
        encode_config: { enable_encode: DOMElements.configForm.querySelector('#enable_encode').checked },
        tls_snippet_config: { enable_site_tls_snippet: DOMElements.configForm.querySelector('#enable_tls_snippet').checked },
    };
    try {
        const result = await api.put(`/config/file/${domain}`, configData);
        state.isEditing = false;
        notification.toast(result.message || t('toasts.save_success'), 'success');
        setTimeout(() => {
            switchView(DOMElements.configListPanel);
            loadAllConfigs();
        }, 500);
    } catch (error) { notification.toast(t('toasts.save_error', { error: error.message }), 'error'); }
}

async function openPresetModal(targetType) {
    const applicablePresets = state.headerPresets.filter(p => p.target === targetType || p.target === 'any');
    if (applicablePresets.length === 0) {
        notification.toast(t('toasts.no_presets_available'), 'info');
        return;
    }

    const presetId = await createPresetSelectionModal(applicablePresets, t);
    if (!presetId) return;

    let targetContainer, keyName, valueName;
    if (targetType === 'global') {
        targetContainer = DOMElements.headersContainer;
        keyName = 'header_key';
        valueName = 'header_value';
    } else if (targetType === 'upstream') {
        targetContainer = DOMElements.upstreamHeadersContainer;
        keyName = 'upstream_header_key';
        valueName = 'upstream_header_value';
    } else {
        return;
    }

    try {
        notification.toast(t('toasts.loading_preset'), 'info', 1000);
        const preset = await api.get(`/config/headers-presets/${presetId}`);
        if (!preset.headers) {
            notification.toast(t('toasts.preset_no_data'), 'info');
            return;
        }

        const choice = await notification.confirm(t('dialogs.preset_fill_msg'), t('dialogs.preset_fill_title'), { confirmText: t('common.append'), cancelText: t('common.replace') });

        if (choice === false) {
            targetContainer.innerHTML = '';
        }

        Object.entries(preset.headers).forEach(([key, values]) => {
            values.forEach(value => {
                addKeyValueInput(targetContainer, keyName, valueName, key, value);
            });
        });
        const presetNameKey = `presets.${preset.id}.name`;
        notification.toast(t('toasts.preset_fill_success', { presetName: t(presetNameKey) }), 'success');

    } catch (error) {
        notification.toast(t('toasts.load_preset_error', { error: error.message }), 'error');
    }
}

// --- 初始化与事件绑定 ---
function pageInit() {
    // 这个函数包含所有特定于 app.js 的初始化逻辑
    loadAllConfigs();

    api.get('/config/headers-presets')
        .then(presets => {
            state.headerPresets = presets || [];
        })
        .catch(err => {
            if (err.message) notification.toast(t('toasts.load_presets_error', { error: err.message }), 'error');
        });

    DOMElements.addNewConfigBtn.addEventListener('click', async () => {
        state.isEditing = false;
        switchView(DOMElements.configFormPanel);
        DOMElements.formTitle.textContent = t('pages.configs.form_title_create');
        DOMElements.configForm.reset();

        const noneButton = DOMElements.serviceModeControl.querySelector('[data-mode="none"]');
        if (noneButton) updateSegmentedControl(noneButton);
        updateServiceModeView('none');
        updateMultiUpstreamView(false);

        state.initialFormState = await getFormStateAsString();
        DOMElements.headersContainer.innerHTML = '';
        DOMElements.upstreamHeadersContainer.innerHTML = '';
        DOMElements.multiUpstreamContainer.innerHTML = '';
        DOMElements.originalFilenameInput.value = '';
    });

    DOMElements.backToListBtn.addEventListener('click', attemptExitForm);
    DOMElements.cancelEditBtn.addEventListener('click', attemptExitForm);

    DOMElements.configForm.addEventListener('submit', handleSaveConfig);
    DOMElements.logoutBtn.addEventListener('click', handleLogout);

    DOMElements.configListContainer.addEventListener('click', e => {
        const button = e.target.closest('button');
        if (!button) return;
        const filename = button.closest('.config-item').dataset.filename;
        if (button.classList.contains('edit-btn')) handleEditConfig(filename);
        if (button.classList.contains('delete-btn')) handleDeleteConfig(filename);
    });

    DOMElements.configForm.querySelector('button[data-add-target="global"]').addEventListener('click', (e) => {
        e.preventDefault();
        addKeyValueInput(DOMElements.headersContainer, 'header_key', 'header_value');
    });

    DOMElements.configForm.querySelector('button[data-add-target="upstream"]').addEventListener('click', (e) => {
        e.preventDefault();
        addKeyValueInput(DOMElements.upstreamHeadersContainer, 'upstream_header_key', 'upstream_header_value');
    });

    DOMElements.configForm.querySelector('button[data-preset-target="global"]').addEventListener('click', (e) => {
        e.preventDefault();
        openPresetModal('global');
    });

    DOMElements.configForm.querySelector('button[data-preset-target="upstream"]').addEventListener('click', (e) => {
        e.preventDefault();
        openPresetModal('upstream');
    });

    DOMElements.addMultiUpstreamBtn.addEventListener('click', (e) => {
        e.preventDefault();
        addSingleInput(DOMElements.multiUpstreamContainer, 'upstream_servers', t('form.upstream_server_placeholder'));
    });

    DOMElements.serviceModeControl.addEventListener('click', (e) => {
        const button = e.target.closest('button');
        if (button) {
            e.preventDefault();
            const mode = button.dataset.mode;
            updateSegmentedControl(button);
            updateServiceModeView(mode);
        }
    });

    DOMElements.mutiUpstreamCheckbox.addEventListener('change', (e) => {
        updateMultiUpstreamView(e.target.checked);
    });

    // --- Re-implemented Accordion Logic ---

    const accordionHeaders = DOMElements.configForm.querySelectorAll('.accordion-header');

    // Function to update the max-height of an accordion content area
    const updateAccordionHeight = (content) => {
        // Only update if the accordion is currently open
        if (content.style.maxHeight !== '0px' && content.style.maxHeight) {
            // Set max-height to its own scroll height, allowing it to expand or shrink
            content.style.maxHeight = content.scrollHeight + 'px';
        }
    };

    accordionHeaders.forEach(header => {
        header.addEventListener('click', () => {
            const content = header.nextElementSibling;
            const parentAccordion = header.closest('.accordion');

            // Toggle the active class on the header
            header.classList.toggle('active');

            // Toggle the max-height to show or hide the content
            if (content.style.maxHeight) {
                content.style.maxHeight = null;
            } else {
                content.style.maxHeight = content.scrollHeight + 'px';
            }

            // After toggling, we may need to update the height of parent accordions.
            // This is a complex task because of animations. We will listen for the transition to end.
            const updateParentAccordions = () => {
                let parent = header.parentElement.closest('.accordion-content');
                while (parent) {
                    // Update the height of the parent accordion content
                    parent.style.maxHeight = parent.scrollHeight + 'px';
                    parent = parent.parentElement.closest('.accordion-content');
                }
            };

            // When our animation finishes, we'll call the update function.
            content.addEventListener('transitionend', updateParentAccordions, { once: true });
        });
    });

    const enableCacheCheckbox = DOMElements.configForm.querySelector('#enable_cache');
    const cacheSubAccordion = DOMElements.configForm.querySelector('#cache-sub-accordion');
    const mainAccordionContent = cacheSubAccordion.closest('.accordion-content');

    // Function to handle changes and update parent accordion
    const handleCacheToggle = (isChecked) => {
        cacheSubAccordion.classList.toggle('hidden', !isChecked);
        // We must update the parent accordion's height *after* the sub-accordion's visibility has changed
        if (mainAccordionContent) {
            // Use a short timeout to allow the browser to render the change in visibility,
            // then update the height. This ensures scrollHeight is calculated correctly.
            setTimeout(() => updateAccordionHeight(mainAccordionContent), 50);
        }
    };

    enableCacheCheckbox.addEventListener('change', (e) => {
        handleCacheToggle(e.target.checked);
    });

    // Also, observe the sub-accordion for any changes that might affect its height
    // This is robust for cases where content is dynamically added/removed from the sub-accordion
    const observer = new MutationObserver(() => {
        if (mainAccordionContent) {
            updateAccordionHeight(mainAccordionContent);
        }
    });

    // Start observing the sub-accordion's content area for changes
    const subAccordionContent = cacheSubAccordion.querySelector('.accordion-content-inner');
    if (subAccordionContent) {
        observer.observe(subAccordionContent, {
            childList: true, // watch for direct children changes
            subtree: true,   // watch for all descendants
            attributes: true, // watch for attribute changes
            characterData: true // watch for text changes
        });
    }
}

// 使用通用初始化函数启动页面
initializePage({ pageId: 'configs', pageInit: pageInit });