import Chart from 'chart.js/auto';
import { sensorChartOptions } from './chart-utils.js';

const dataEl = document.getElementById('chart-data');
const rows = JSON.parse(dataEl.dataset.points);
const soilLabel = dataEl.dataset.labelSoil;
const labels = rows.map(function (r) {
    return r.t;
});

const marker = document.getElementById('soil-marker');
const label = document.getElementById('soil-current-label');
const bar = document.querySelector('.soil-gradient-bar');

function updateMarker(pct) {
    // clamp 0–100
    const clamped = Math.min(100, Math.max(0, pct));
    const barWidth = bar.getBoundingClientRect().width;
    // offset by half marker width (6px) to centre the triangle
    marker.style.left = `${(clamped / 100) * barWidth - 6}px`;
    label.textContent = `${clamped.toFixed(1)}%`;
}

const latest = rows[rows.length - 1];
if (latest) updateMarker(latest.pct);

new Chart(document.getElementById('soil-chart'), {
    type: 'line',
    data: {
        labels,
        datasets: [
            {
                label: soilLabel,
                data: rows.map(function (r) {
                    return r.pct;
                }),
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
