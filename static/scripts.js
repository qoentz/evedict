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
