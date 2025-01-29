// /static/js/modules/menu.js

export function initMenu() {
    const menuIcon = document.querySelector('.menu-icon');
    const mobileNav = document.getElementById('mobile-nav');
    const overlay = document.getElementById('overlay');
    const body = document.body;

    // Select the Close Button within Mobile Navigation
    const closeButton = document.querySelector('.close-menu-button');

    if (!menuIcon || !mobileNav || !overlay || !closeButton) {
        console.error("Menu Icon, Mobile Nav Menu, Overlay, or Close Button not found!");
        return;
    }

    const toggleMenu = () => {
        mobileNav.classList.toggle('-translate-x-full'); // Slide out
        mobileNav.classList.toggle('translate-x-0'); // Slide in
        overlay.classList.toggle('hidden');
        menuIcon.classList.toggle('active');
        body.classList.toggle('overflow-hidden');

        const isExpanded = menuIcon.classList.contains('active');
        menuIcon.setAttribute('aria-expanded', isExpanded);
    };

    // Attach event listeners to Menu Icon, Close Button, and Overlay
    menuIcon.addEventListener('click', toggleMenu);
    closeButton.addEventListener('click', toggleMenu);
    overlay.addEventListener('click', toggleMenu);

    // Close the menu when any navigation link is clicked
    const mobileNavLinks = document.getElementById('mobile-nav-links');
    if (mobileNavLinks) {
        mobileNavLinks.addEventListener('click', (event) => {
            const target = event.target;
            if (target.tagName === 'A') {
                toggleMenu();
            }
        });
    }
}
