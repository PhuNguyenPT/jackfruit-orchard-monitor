document.addEventListener('DOMContentLoaded', function () {
    const toggles = document.querySelectorAll('[data-target]');

    toggles.forEach(function (toggle) {
        toggle.addEventListener('click', function (e) {
            e.stopPropagation();
            const target = document.getElementById(toggle.dataset.target);
            if (!target) return;

            const isOpen = !target.classList.contains('hidden');

            // Close any other open menus before opening this one
            toggles.forEach(function (other) {
                if (other === toggle) return;
                const otherTarget = document.getElementById(
                    other.dataset.target
                );
                if (otherTarget) otherTarget.classList.add('hidden');
            });

            target.classList.toggle('hidden', isOpen);
        });
    });

    // Close any open menu when tapping/clicking outside it
    document.addEventListener('click', function (e) {
        toggles.forEach(function (toggle) {
            const target = document.getElementById(toggle.dataset.target);
            if (!target) return;
            if (!target.contains(e.target) && !toggle.contains(e.target)) {
                target.classList.add('hidden');
            }
        });
    });
});
