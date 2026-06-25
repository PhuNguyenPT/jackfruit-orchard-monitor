import Chart from 'chart.js/auto';
import { sensorChartOptions, paddedRange } from './chart-utils.js';

const dataEl = document.getElementById('chart-data');
const rows = JSON.parse(dataEl.dataset.points);
const nonce = dataEl.dataset.nonce;
const tempLabel = dataEl.dataset.labelTemp;
const humidLabel = dataEl.dataset.labelHumid;
const wsUrl = dataEl.dataset.wsUrl;
const labels = rows.map((r) => r.t);
let lastTs = rows.length > 0 ? rows[rows.length - 1].t : null;

const TEMP_MIN = -40;
const TEMP_MAX = 100;
const MAX_POINTS = 100;

// ── Live indicator ────────────────────────────────────────────────────────
const indicatorWrap = document.createElement('div');
indicatorWrap.className = 'flex items-center gap-2 mb-2';
indicatorWrap.innerHTML =
    '<span id="live-dot" class="w-2 h-2 rounded-full bg-gray-300 inline-block transition-colors"></span>' +
    '<span id="live-text" class="text-xs text-gray-400">Connecting…</span>';

const tempCanvasEl = document.getElementById('temp-chart');
tempCanvasEl.parentElement.insertBefore(indicatorWrap, tempCanvasEl);

// Second indicator — same state, mirrors the temp one
const indicatorWrap2 = document.createElement('div');
indicatorWrap2.className = 'flex items-center gap-2 mb-2';
indicatorWrap2.innerHTML =
    '<span id="live-dot2" class="w-2 h-2 rounded-full bg-gray-300 inline-block transition-colors"></span>' +
    '<span id="live-text2" class="text-xs text-gray-400">Connecting…</span>';

const humidCanvasEl = document.getElementById('humid-chart');
humidCanvasEl.parentElement.insertBefore(indicatorWrap2, humidCanvasEl);

function setLive(on) {
    const dot = document.getElementById('live-dot');
    const txt = document.getElementById('live-text');
    const dot2 = document.getElementById('live-dot2');
    const txt2 = document.getElementById('live-text2');
    if (on) {
        const dotClass =
            'w-2 h-2 rounded-full bg-green-500 animate-pulse inline-block transition-colors';
        const txtClass = 'text-xs text-green-600';
        dot.className = dotClass;
        txt.textContent = 'Live';
        txt.className = txtClass;
        dot2.className = dotClass;
        txt2.textContent = 'Live';
        txt2.className = txtClass;
    } else {
        const dotClass =
            'w-2 h-2 rounded-full bg-gray-300 inline-block transition-colors';
        const txtClass = 'text-xs text-gray-400';
        dot.className = dotClass;
        txt.textContent = 'Reconnecting…';
        txt.className = txtClass;
        dot2.className = dotClass;
        txt2.textContent = 'Reconnecting…';
        txt2.className = txtClass;
    }
}

// ── Marker ────────────────────────────────────────────────────────────────
function updateTempMarker(latestTemp) {
    const marker = document.getElementById('temp-marker');
    const label = document.getElementById('temp-current-label');
    if (!marker) return;
    const clamped = Math.min(TEMP_MAX, Math.max(TEMP_MIN, latestTemp));
    const pct = ((clamped - TEMP_MIN) / (TEMP_MAX - TEMP_MIN)) * 100;
    const halfWidthPx = marker.offsetWidth / 2;
    let styleTag = document.getElementById('temp-marker-style');
    if (!styleTag) {
        styleTag = document.createElement('style');
        styleTag.id = 'temp-marker-style';
        styleTag.nonce = nonce;
        document.head.appendChild(styleTag);
    }
    styleTag.textContent = `#temp-marker { left: calc(${pct}% - ${halfWidthPx}px) !important; }`;
    if (label) label.textContent = latestTemp.toFixed(1) + '°C';
}

if (rows.length > 0) updateTempMarker(rows[rows.length - 1].temp);

const temps = rows.map((r) => r.temp);

const tempChart = new Chart(document.getElementById('temp-chart'), {
    type: 'line',
    data: {
        labels: [...labels],
        datasets: [
            {
                label: tempLabel,
                data: temps,
                borderColor: '#f97316',
                backgroundColor: 'rgba(249,115,22,0.08)',
                borderWidth: 2,
                pointRadius: 2,
                tension: 0.3,
                fill: true,
            },
        ],
    },
    options: sensorChartOptions({ suffix: '°C', ...paddedRange(temps, 2) }),
});

const humidChart = new Chart(document.getElementById('humid-chart'), {
    type: 'line',
    data: {
        labels: [...labels],
        datasets: [
            {
                label: humidLabel,
                data: rows.map((r) => r.humid),
                borderColor: '#3b82f6',
                backgroundColor: 'rgba(59,130,246,0.08)',
                borderWidth: 2,
                pointRadius: 2,
                tension: 0.3,
                fill: true,
            },
        ],
    },
    options: sensorChartOptions({ suffix: '%', min: 0, max: 100 }),
});

// ── Resize on tab focus (fixes white-space after hidden tab) ──────────────
document.addEventListener('visibilitychange', () => {
    if (!document.hidden) {
        tempChart.resize();
        humidChart.resize();
    }
});

// ── Real-time updates ─────────────────────────────────────────────────────
const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
let retryDelay = 1000;
function applyPoint(t, temp, humid) {
    tempChart.data.labels.push(t);
    tempChart.data.datasets[0].data.push(temp);
    humidChart.data.labels.push(t);
    humidChart.data.datasets[0].data.push(humid);
    if (tempChart.data.labels.length > MAX_POINTS) {
        tempChart.data.labels.shift();
        tempChart.data.datasets[0].data.shift();
        humidChart.data.labels.shift();
        humidChart.data.datasets[0].data.shift();
    }
    lastTs = t;
}

function onMessage(event) {
    const msg = JSON.parse(event.data);
    if (msg.batch) {
        for (const pt of msg.points) applyPoint(pt.t, pt.temp, pt.humid);
        tempChart.update('none');
        humidChart.update('none');
        if (msg.points.length > 0) {
            updateTempMarker(msg.points[msg.points.length - 1].temp);
        }
        return;
    }
    applyPoint(msg.t, msg.temp, msg.humid);
    tempChart.update('none');
    humidChart.update('none');
    updateTempMarker(msg.temp);
}

function connect() {
    const sinceParam = lastTs ? '?since=' + encodeURIComponent(lastTs) : '';
    const ws = new WebSocket(proto + '//' + location.host + wsUrl + sinceParam);
    ws.addEventListener('open', () => {
        setLive(true);
        retryDelay = 1000;
    });
    ws.addEventListener('message', onMessage);
    ws.addEventListener('close', () => {
        setLive(false);
        setTimeout(connect, retryDelay);
        retryDelay = Math.min(retryDelay * 2, 30000);
    });
    ws.addEventListener('error', () => ws.close());
}

connect();
