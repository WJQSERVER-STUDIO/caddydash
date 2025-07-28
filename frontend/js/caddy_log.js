import { initializePage } from './common.js';
import { api } from './api.js';
import { notification } from './notifications.js';
import { t } from './locale.js';

function pageInit() {
    const tabs = document.querySelectorAll('.tab-link');
    const tabContents = document.querySelectorAll('.tab-content');

    const logContents = {
        stdout: document.getElementById('stdout-log-content'),
        caddy: document.getElementById('caddy-log-content'),
        sites: document.getElementById('site-log-content'),
    };

    const siteLogSelector = document.getElementById('site-log-selector');
    let eventSources = {};
    let lineCounters = {
        stdout: 0,
        caddy: 0,
        sites: 0,
    };

    function clearLog(logType) {
        if (logContents[logType]) {
            logContents[logType].innerHTML = '';
            lineCounters[logType] = 0;
        }
    }

    function connectToLogSource(logType, url) {
        if (eventSources[logType]) {
            eventSources[logType].close();
        }
        clearLog(logType);
        eventSources[logType] = new EventSource(url);

        eventSources[logType].addEventListener('connected', function(e) {
            console.log(`SSE connection for ${logType} established.`);
            notification.toast(t('toasts.log_connection_established'), 'success');
        });

        eventSources[logType].onmessage = function(event) {
            if (logContents[logType]) {
                const lines = event.data.split('\\n');
                lines.forEach(line => {
                    if (line.trim() !== '') {
                        lineCounters[logType]++;
                        const logLine = document.createElement('div');
                        logLine.classList.add('log-line');

                        const lineNumberDiv = document.createElement('div');
                        lineNumberDiv.classList.add('log-line-number');
                        lineNumberDiv.textContent = lineCounters[logType];

                        const lineContentDiv = document.createElement('div');
                        lineContentDiv.classList.add('log-line-content');
                        lineContentDiv.textContent = line;

                        logLine.appendChild(lineNumberDiv);
                        logLine.appendChild(lineContentDiv);

                        logContents[logType].appendChild(logLine);
                    }
                });
                logContents[logType].scrollTop = logContents[logType].scrollHeight;
            }
        };

        eventSources[logType].onerror = function(err) {
            console.error(`EventSource for ${logType} failed:`, err);
            let errorMessage = t('toasts.log_connection_error');
            if (err && err.target && err.target.readyState === EventSource.CLOSED) {
                errorMessage = t('toasts.log_connection_closed');
            }
            notification.toast(`${t(`pages.caddy_logs.log_type_${logType}`)}: ${errorMessage}`, 'error');
            eventSources[logType].close();
        };
    }

    function loadSiteLogs() {
        api.get('/caddy/logs/sites')
            .then(sites => {
                siteLogSelector.innerHTML = '';
                if (sites && sites.length > 0) {
                    sites.forEach(site => {
                        const option = document.createElement('option');
                        option.value = site;
                        option.textContent = site;
                        siteLogSelector.appendChild(option);
                    });
                    siteLogSelector.dispatchEvent(new Event('change'));
                } else {
                    const option = document.createElement('option');
                    option.textContent = t('pages.caddy_logs.no_site_logs');
                    option.disabled = true;
                    siteLogSelector.appendChild(option);
                }
            })
            .catch(error => {
                console.error('Error fetching site logs:', error);
                notification.toast(t('toasts.load_site_logs_error'), 'error');
            });
    }

    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            tabs.forEach(t => t.classList.remove('active'));
            tab.classList.add('active');

            tabContents.forEach(c => c.classList.remove('active'));
            document.getElementById(`${tab.dataset.tab}-tab`).classList.add('active');

            const logType = tab.dataset.tab;
            if (logType === 'sites') {
                if (siteLogSelector.options.length === 0 || siteLogSelector.options[0].disabled) {
                    loadSiteLogs();
                } else {
                    siteLogSelector.dispatchEvent(new Event('change'));
                }
            } else {
                connectToLogSource(logType, `/v0/api/caddy/logs/${logType}`);
            }
        });
    });

    siteLogSelector.addEventListener('change', function() {
        const sitename = this.value;
        if (sitename) {
            connectToLogSource('sites', `/v0/api/caddy/logs/site/${sitename}`);
        }
    });

    // Initialize the first tab
    connectToLogSource('stdout', '/v0/api/caddy/logs/stdout');
}

initializePage({ pageId: 'caddy_logs', pageInit: pageInit });
