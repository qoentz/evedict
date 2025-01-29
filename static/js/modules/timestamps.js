export function convertTimestampsToLocal() {
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
