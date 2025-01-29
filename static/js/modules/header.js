export function initHeader() {
    const header = document.querySelector('.shrinkable-header');
    const SCROLL_THRESHOLD = 100; // px threshold

    window.addEventListener('scroll', () => {
        if (window.scrollY > SCROLL_THRESHOLD) {
            header.classList.add('scrolled');
        } else {
            header.classList.remove('scrolled');
        }
    });
}
