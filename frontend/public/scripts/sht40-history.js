import Chart from 'chart.js/auto';
import { sensorChartOptions, paddedRange } from './chart-utils.js';

const dataEl = document.getElementById('chart-data');
const rows = JSON.parse(dataEl.dataset.points);
const nonce = dataEl.dataset.nonce;

const tempLabel = dataEl.dataset.labelTemp;
const humidLabel = dataEl.dataset.labelHumid;
const labels = rows.map(function (r) {
    return r.t;
});

const TEMP_MIN = -40;
const TEMP_MAX = 100;

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

    if (label) {
        label.textContent = latestTemp.toFixed(1) + '°C';
    }
}

if (rows.length > 0) {
    updateTempMarker(rows[rows.length - 1].temp);
}

const temps = rows.map(function (r) {
    return r.temp;
});

new Chart(document.getElementById('temp-chart'), {
    type: 'line',
    data: {
        labels,
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
    options: sensorChartOptions({
        suffix: '°C',
        // Soft range: padded ±2 °C around actual data so small fluctuations
        // don't look like dramatic spikes, while still showing the real scale.
        ...paddedRange(temps, 2),
    }),
});

new Chart(document.getElementById('humid-chart'), {
    type: 'line',
    data: {
        labels,
        datasets: [
            {
                label: humidLabel,
                data: rows.map(function (r) {
                    return r.humid;
                }),
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
