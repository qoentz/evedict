export function convertTimestampsToLocal() {
    document.querySelectorAll('.timestamp').forEach(element => {
        const timestampUTC = element.getAttribute('data-timestamp');
        if (timestampUTC) {
            const localDate = new Date(timestampUTC);

            const datePart = localDate.toLocaleDateString('en-GB', {
                year: 'numeric',
                month: 'short',
                day: 'numeric'
            }).replace(/,/g, ''); // removes commas from en-GB format

            const timePart = localDate.toLocaleTimeString('sv-SE', {
                hour: '2-digit',
                minute: '2-digit',
                hour12: false
            });

            element.textContent = `${datePart}, ${timePart}`;
        }
    });
}