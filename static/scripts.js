function toggleSources(button) {
    const container = button.closest('.outcomes-container') || button.closest('.sources-list');
    const siblingContainer = container.nextElementSibling || container.previousElementSibling;

    container.classList.toggle('hidden');
    siblingContainer.classList.toggle('hidden');
}

function toggleMenu() {
    const menuIcon = document.querySelector('.menu-icon');
    const navMenu = document.querySelector('.nav-menu');
    menuIcon.classList.toggle('active');
    navMenu.classList.toggle('active');
}
