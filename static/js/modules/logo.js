export function initBall() {
    const ballLink = document.querySelector('.hand-ball-container a');
    if (ballLink) {
        ballLink.addEventListener('click', restartLetterAnimation);
    }
}

function restartLetterAnimation() {
    // Select all the spans inside your header (adjust the selector if needed)
    const letterSpans = document.querySelectorAll('h1 span');
    letterSpans.forEach(el => {
        // Remove the animation class
        el.classList.remove('animate-letter-slide-in');
        // Force a reflow to reset the animation (this trick makes the browser recognize a change)
        void el.offsetWidth;
        // Re‑add the animation class so the animation plays again
        el.classList.add('animate-letter-slide-in');
    });
}