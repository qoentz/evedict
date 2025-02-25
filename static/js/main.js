// /static/js/main.js

import { initHeader } from './modules/header.js';
import { convertTimestampsToLocal } from './modules/timestamps.js';
import { initSlider } from './modules/slider.js';
import { initBall} from "./modules/logo.js";

document.addEventListener('DOMContentLoaded', () => {
    initHeader()
    initBall()
    convertTimestampsToLocal();
});

document.addEventListener('htmx:afterSwap', (event) => {
    const target = event.detail.target;
    if (target && target.querySelector('#slider-container')) {
        initSlider();
    }
});

document.addEventListener("htmx:historyRestore", (event) => {
    console.log("htmx:historyRestore fired");

    // Remove letter animation classes only from the logo text
    document.querySelectorAll(".logo-text-wrapper span.animate-letter-slide-in")
        .forEach(el => el.classList.remove("animate-letter-slide-in"));

    // Remove the fade-slide animation from the container (but not the ball image)
    document.querySelectorAll(".logo-wrapper.animate-fade-slide-in")
        .forEach(el => el.classList.remove("animate-fade-slide-in"));
});
















