document.addEventListener('DOMContentLoaded', function() {
    const header = document.querySelector('.shrinkable-header');
    const SCROLL_THRESHOLD = 100; // px threshold

    window.addEventListener('scroll', function() {
        if (window.scrollY > SCROLL_THRESHOLD) {
            header.classList.add('scrolled');
        } else {
            header.classList.remove('scrolled');
        }
    });
});

function toggleSources(button) {
    const container = button.closest('.outcomes-container') || button.closest('.sources-list');
    const siblingContainer = container.nextElementSibling || container.previousElementSibling;

    container.classList.toggle('hidden');
    siblingContainer.classList.toggle('hidden');
}

function toggleMenu() {
    const menuIcon = document.querySelector('.menu-icon');
    const navMenu = document.getElementById('nav-menu');

    if (!menuIcon || !navMenu) {
        console.error("Menu Icon or Nav Menu not found!");
        return;
    }

    // Toggle visibility
    navMenu.classList.toggle('hidden');
    menuIcon.classList.toggle('active');

    console.log("Menu visibility toggled");
}

function convertTimestampsToLocal() {
    document.querySelectorAll('.timestamp').forEach(element => {
        const timestampUTC = element.getAttribute('data-timestamp');
        if (timestampUTC) {
            const localDate = new Date(timestampUTC);
            element.textContent = localDate.toLocaleString(undefined, {
                year: 'numeric', month: 'short', day: '2-digit',
                hour: '2-digit', minute: '2-digit'
            });
        }
    });

}

document.addEventListener('htmx:afterSwap', convertTimestampsToLocal);
