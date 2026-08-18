document.addEventListener('DOMContentLoaded', function () {
    document.querySelectorAll('[data-toggle-password]').forEach(function (btn) {
        const input = document.getElementById(btn.dataset.togglePassword);
        if (!input) return;

        const showIcon = btn.querySelector('[data-icon="show"]');
        const hideIcon = btn.querySelector('[data-icon="hide"]');

        btn.addEventListener('click', function () {
            const isHidden = input.type === 'password';
            input.type = isHidden ? 'text' : 'password';
            if (showIcon) showIcon.classList.toggle('hidden', isHidden);
            if (hideIcon) hideIcon.classList.toggle('hidden', !isHidden);
            btn.setAttribute('aria-pressed', isHidden ? 'true' : 'false');
        });
    });
});
