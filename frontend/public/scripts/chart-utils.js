/**
 * Shared Chart.js helpers for sensor history pages.
 */

/**
 * Returns a consistent Chart.js options object.
 *
 * @param {object} opts
 * @param {string}  opts.suffix        - Tooltip/tick suffix, e.g. '°C' or '%'
 * @param {number}  [opts.min]         - Hard y-axis minimum (use for bounded ranges like 0–100)
 * @param {number}  [opts.max]         - Hard y-axis maximum
 * @param {number}  [opts.suggestedMin] - Soft minimum (Chart.js will not go below this unless data forces it)
 * @param {number}  [opts.suggestedMax] - Soft maximum
 */
export function sensorChartOptions({
    suffix = '',
    min,
    max,
    suggestedMin,
    suggestedMax,
} = {}) {
    return {
        responsive: true,
        maintainAspectRatio: false,
        interaction: { mode: 'index', intersect: false },
        plugins: {
            legend: { display: true },
            tooltip: {
                callbacks: {
                    label(ctx) {
                        return ctx.dataset.label + ': ' + ctx.parsed.y + suffix;
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
                ...(min !== undefined && { min }),
                ...(max !== undefined && { max }),
                ...(suggestedMin !== undefined && { suggestedMin }),
                ...(suggestedMax !== undefined && { suggestedMax }),
                ticks: { font: { size: 11 }, color: '#6b7280' },
                grid: { color: 'rgba(0,0,0,0.04)' },
            },
        },
    };
}

/**
 * Computes a padded y-axis range from an array of values.
 * Returns { suggestedMin, suggestedMax } suitable for spreading into sensorChartOptions.
 *
 * @param {number[]} values
 * @param {number}   pad  - Padding added to each side (default 2)
 */
export function paddedRange(values, pad = 2) {
    const min = Math.min(...values);
    const max = Math.max(...values);
    return {
        suggestedMin: Math.floor(min - pad),
        suggestedMax: Math.ceil(max + pad),
    };
}
