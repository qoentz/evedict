export function toggleSources(button) {
    const container = button.closest('.outcomes-container') || button.closest('.sources-list');
    const siblingContainer = container.nextElementSibling || container.previousElementSibling;

    container.classList.toggle('hidden');
    siblingContainer.classList.toggle('hidden');
}
