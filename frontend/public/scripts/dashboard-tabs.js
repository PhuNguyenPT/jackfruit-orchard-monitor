document.addEventListener('DOMContentLoaded', function () {
    const tabButtons = document.querySelectorAll('[data-dashboard-tab]');
    if (!tabButtons.length) return;

    const activeClasses = ['text-blue-600', 'border-blue-600'];
    const inactiveClasses = [
        'text-gray-500',
        'border-transparent',
        'hover:text-gray-700',
    ];

    function activate(tabName) {
        tabButtons.forEach(function (btn) {
            const isActive = btn.dataset.dashboardTab === tabName;
            btn.classList.remove(...activeClasses, ...inactiveClasses);
            btn.classList.add(...(isActive ? activeClasses : inactiveClasses));
        });

        document
            .querySelectorAll('[data-dashboard-panel]')
            .forEach(function (panel) {
                panel.classList.toggle(
                    'hidden',
                    panel.dataset.dashboardPanel !== tabName
                );
            });
    }

    tabButtons.forEach(function (btn) {
        btn.addEventListener('click', function () {
            const tab = btn.dataset.dashboardTab;
            activate(tab);

            const url = new URL(window.location.href);
            if (tab === 'account') {
                url.searchParams.delete('tab');
            } else {
                url.searchParams.set('tab', tab);
            }
            history.replaceState(null, '', url.pathname + url.search);
        });
    });
});
