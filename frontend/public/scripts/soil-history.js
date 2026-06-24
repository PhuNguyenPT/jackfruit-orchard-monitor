import Chart from 'chart.js/auto';

const dataEl = document.getElementById('chart-data');
const rows = JSON.parse(dataEl.dataset.points);
const labels = rows.map(function (r) {
    return r.t;
});

new Chart(document.getElementById('soil-chart'), {
    type: 'line',
    data: {
        labels,
        datasets: [
            {
                label: 'Soil Moisture (%)',
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
    options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: { mode: 'index', intersect: false },
        plugins: {
            legend: { display: true },
            tooltip: {
                callbacks: {
                    label: function (ctx) {
                        return ctx.dataset.label + ': ' + ctx.parsed.y + '%';
                    },
                },
            },
        },
        scales: {
            x: {
                ticks: {
                    maxTicksLimit: 8,
                    font: { size: 11 },
                    color: '#9ca3af',
                },
                grid: { display: false },
            },
            y: {
                min: 0,
                max: 100,
                ticks: { font: { size: 11 }, color: '#6b7280' },
                grid: { color: 'rgba(0,0,0,0.04)' },
            },
        },
    },
});
