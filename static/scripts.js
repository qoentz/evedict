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
    const mobileNav = document.getElementById('mobile-nav');
    const overlay = document.getElementById('overlay');
    const body = document.body;

    if (!menuIcon || !mobileNav || !overlay) {
        console.error("Menu Icon, Mobile Nav Menu, or Overlay not found!");
        return;
    }

    // Toggle the mobile navigation menu
    mobileNav.classList.toggle('-translate-x-full'); // Slide out
    mobileNav.classList.toggle('translate-x-0'); // Slide in

    // Toggle the overlay
    overlay.classList.toggle('hidden');

    // Toggle the 'active' state on the menu icon
    menuIcon.classList.toggle('active');

    // Toggle body overflow to prevent background scrolling
    body.classList.toggle('overflow-hidden');

    // Update aria-expanded for accessibility
    const isExpanded = menuIcon.classList.contains('active');
    menuIcon.setAttribute('aria-expanded', isExpanded);

    console.log("Mobile menu visibility toggled");
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

document.addEventListener('DOMContentLoaded', function() {
    const overlay = document.getElementById('overlay');
    const mobileNavLinks = document.getElementById('mobile-nav-links');

    if (overlay) {
        overlay.addEventListener('click', toggleMenu);
    }

    if (mobileNavLinks) {
        mobileNavLinks.addEventListener('click', function(event) {
            const target = event.target;
            if (target.tagName === 'A') {
                // Optionally close the menu when a link is clicked
                toggleMenu();
            }
        });
    }
});

