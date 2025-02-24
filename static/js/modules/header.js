export function initHeader() {
    updateHeader();
    window.addEventListener('scroll', updateHeader);
}

export function updateHeader() {
    const header = document.querySelector('.shrinkable-header');
    const SCROLL_THRESHOLD = 100;
    if (!header) return;
    if (window.scrollY > SCROLL_THRESHOLD) {
        header.classList.add('scrolled');
    } else {
        header.classList.remove('scrolled');
    }
}
