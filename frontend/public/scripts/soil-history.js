import Chart from 'chart.js/auto';
import { sensorChartOptions } from './chart-utils.js';

const dataEl = document.getElementById('chart-data');
const rows = JSON.parse(dataEl.dataset.points);
const soilLabel = dataEl.dataset.labelSoil;
const wsUrl = dataEl.dataset.wsUrl;
const labels = rows.map((r) => r.t);
let lastTs = rows.length > 0 ? rows[rows.length - 1].t : null;

const MAX_POINTS = 100;

// ── Live indicator ────────────────────────────────────────────────────────
const indicatorWrap = document.createElement('div');
indicatorWrap.className = 'flex items-center gap-2 mb-2';
indicatorWrap.innerHTML =
    '<span id="live-dot" class="w-2 h-2 rounded-full bg-gray-300 inline-block transition-colors"></span>' +
    '<span id="live-text" class="text-xs text-gray-400">Connecting…</span>';

const soilCanvasEl = document.getElementById('soil-chart');
soilCanvasEl.parentElement.insertBefore(indicatorWrap, soilCanvasEl);

function setLive(on) {
    const dot = document.getElementById('live-dot');
    const txt = document.getElementById('live-text');
    if (on) {
        dot.className =
            'w-2 h-2 rounded-full bg-green-500 animate-pulse inline-block transition-colors';
        txt.textContent = 'Live';
        txt.className = 'text-xs text-green-600';
    } else {
        dot.className =
            'w-2 h-2 rounded-full bg-gray-300 inline-block transition-colors';
        txt.textContent = 'Reconnecting…';
        txt.className = 'text-xs text-gray-400';
    }
}

// ── Marker ────────────────────────────────────────────────────────────────
const marker = document.getElementById('soil-marker');
const label = document.getElementById('soil-current-label');
const bar = document.querySelector('.soil-gradient-bar');

function updateMarker(pct) {
    const clamped = Math.min(100, Math.max(0, pct));
    const barWidth = bar.getBoundingClientRect().width;
    marker.style.left = `${(clamped / 100) * barWidth - 6}px`;
    label.textContent = `${clamped.toFixed(1)}%`;
}

const latest = rows[rows.length - 1];
if (latest) updateMarker(latest.pct);

const soilChart = new Chart(document.getElementById('soil-chart'), {
    type: 'line',
    data: {
        labels,
        datasets: [
            {
                label: soilLabel,
                data: rows.map((r) => r.pct),
                borderColor: '#10b981',
                backgroundColor: 'rgba(16,185,129,0.08)',
                borderWidth: 2,
                pointRadius: 2,
                tension: 0.3,
                fill: true,
            },
        ],
    },
    options: sensorChartOptions({ suffix: '%', min: 0, max: 100 }),
});

document.addEventListener('visibilitychange', () => {
    if (!document.hidden) {
        soilChart.resize();
    }
});

// ── Real-time updates ─────────────────────────────────────────────────────
const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
let retryDelay = 1000;

function applyPoint(t, pct) {
    soilChart.data.labels.push(t);
    soilChart.data.datasets[0].data.push(pct);
    if (soilChart.data.labels.length > MAX_POINTS) {
        soilChart.data.labels.shift();
        soilChart.data.datasets[0].data.shift();
    }
    lastTs = t;
}

function onMessage(event) {
    const msg = JSON.parse(event.data);
    if (msg.batch) {
        for (const pt of msg.points) applyPoint(pt.t, pt.pct);
        soilChart.update('none');
        if (msg.points.length > 0) {
            updateMarker(msg.points[msg.points.length - 1].pct);
        }
        return;
    }
    applyPoint(msg.t, msg.pct);
    soilChart.update('none');
    updateMarker(msg.pct);
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
