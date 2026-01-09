// Strix Halo Installer - Web Frontend

document.addEventListener('DOMContentLoaded', init);

async function init() {
    await loadDevice();
    await loadStages();
    setupEventHandlers();
}

async function loadDevice() {
    try {
        const res = await fetch('/api/device');
        const data = await res.json();
        document.getElementById('device-name').textContent = data.name;
    } catch (err) {
        document.getElementById('device-name').textContent = 'Detection failed';
    }
}

async function loadStages() {
    try {
        const res = await fetch('/api/stages');
        const stages = await res.json();

        const list = document.getElementById('stage-list');
        list.innerHTML = '';

        stages.forEach((stage, i) => {
            const li = document.createElement('li');
            li.id = `stage-${stage.id}`;
            li.innerHTML = `
                <span class="stage-status" id="status-${stage.id}"></span>
                <span class="stage-name">${stage.name}</span>
                ${stage.optional ? '<span class="optional">(optional)</span>' : ''}
            `;
            list.appendChild(li);
        });
    } catch (err) {
        console.error('Failed to load stages:', err);
    }
}

function setupEventHandlers() {
    document.getElementById('start-btn').addEventListener('click', startInstall);
    document.getElementById('tui-btn').addEventListener('click', () => {
        alert('Close this window and run: strix-install --tui');
    });
}

async function startInstall() {
    const btn = document.getElementById('start-btn');
    btn.disabled = true;
    btn.textContent = '⏳ Installing...';

    document.getElementById('log').classList.remove('hidden');

    try {
        const res = await fetch('/api/run', { method: 'POST' });
        const data = await res.json();

        if (data.status === 'started') {
            log('Installation started...');
            // In a real implementation, we'd poll for progress or use WebSocket
            // For now, just show that it started
        }
    } catch (err) {
        log('Error: ' + err.message);
        btn.disabled = false;
        btn.textContent = '▶ Retry Installation';
    }
}

function log(message) {
    const output = document.getElementById('log-output');
    const time = new Date().toLocaleTimeString();
    output.textContent += `[${time}] ${message}\n`;
    output.scrollTop = output.scrollHeight;
}

function setStageStatus(stageId, status) {
    const el = document.getElementById(`status-${stageId}`);
    if (el) {
        el.className = 'stage-status ' + status;
        if (status === 'done') {
            el.textContent = '✓';
        } else if (status === 'running') {
            el.textContent = '●';
        }
    }
}
