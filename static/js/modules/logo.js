export function initBall() {
    // Attach click listener to restart letter animations.
    const ballLink = document.querySelector('.hand-ball-container a');
    if (ballLink) {
        ballLink.addEventListener('click', restartLetterAnimation);
    }

    // Start the ball spin loop.
    startBallSpin();
}

// Exported function to initialize our historyRestore handler.
export function initLogoHistoryRestore() {
    document.addEventListener("htmx:historyRestore", (event) => {
        // Remove letter animation classes from logo text.
        document.querySelectorAll(".logo-text-wrapper span.animate-letter-slide-in")
            .forEach(el => el.classList.remove("animate-letter-slide-in"));

        // Remove fade-slide animation from the logo wrapper.
        document.querySelectorAll(".logo-wrapper.animate-fade-slide-in")
            .forEach(el => el.classList.remove("animate-fade-slide-in"));

        // Scroll to the top of the page.
        setTimeout(() => {
            window.scrollTo(0, 0);
        }, 0);
    });
}

// Restart the letter animations on click.
function restartLetterAnimation() {
    const letterSpans = document.querySelectorAll('h1 span');
    letterSpans.forEach(el => {
        el.classList.remove('animate-letter-slide-in');
        // Force reflow to reset animation
        void el.offsetWidth;
        el.classList.add('animate-letter-slide-in');
    });
}

// --- Ball Spin Logic --- //

// Global variable for the requestAnimationFrame ID.
let ballAnimationFrame = null;

// Retrieve (or set) a persistent start time so the spin remains consistent.
let ballAnimationStart = sessionStorage.getItem("ballAnimationStart");
if (!ballAnimationStart) {
    ballAnimationStart = Date.now();
    sessionStorage.setItem("ballAnimationStart", ballAnimationStart.toString());
} else {
    ballAnimationStart = parseInt(ballAnimationStart, 10);
}

// Function that updates the ball's rotation using requestAnimationFrame.
function updateBallSpin() {
    const ball = document.getElementById("ball");
    if (!ball) return;

    const duration = 15; // seconds for one full rotation
    const elapsed = (Date.now() - ballAnimationStart) / 1000;
    // Calculate angle (0 to 360)
    const angle = (elapsed * 360 / duration) % 360;
    ball.style.transform = `rotate(${angle}deg)`;

    ballAnimationFrame = requestAnimationFrame(updateBallSpin);
}

// Starts the spin loop if not already running.
function startBallSpin() {
    if (ballAnimationFrame === null) {
        updateBallSpin();
    }
}
