// /static/js/main.js

import { initHeader } from './modules/header.js';
import { convertTimestampsToLocal } from './modules/timestamps.js';
import { initSlider } from './modules/slider.js';
import { initBall, initLogoHistoryRestore } from './modules/logo.js';
import { initNetworkSphere, initAmbientBackground } from "./modules/auxiliary.js";

document.addEventListener('DOMContentLoaded', () => {
    function setFullViewportHeight() {
        const vh = window.innerHeight * 0.01;
        document.documentElement.style.setProperty('--vh', `${vh}px`);
    }
    setFullViewportHeight();
    window.addEventListener('resize', setFullViewportHeight);

    // Initialize both scenes
    initNetworkSphere();      // For pages with globe
    initAmbientBackground();  // For carousel/pages without globe

    initHeader();
    initBall();
    initLogoHistoryRestore();
    convertTimestampsToLocal();

    // Handle HTMX page swaps to reinitialize scenes
    document.body.addEventListener('htmx:afterSwap', (e) => {
        initNetworkSphere();
        initAmbientBackground();
    });
});

// Handle window resize for both scenes
window.addEventListener('resize', () => {
    // The individual fit() functions are called automatically by each scene's resize listeners
});

document.addEventListener('htmx:afterSwap', (event) => {
    initNetworkSphere();
    initAmbientBackground();
    convertTimestampsToLocal();
    const target = event.detail.target;
    if (target && target.querySelector('#slider-container')) {
        initSlider();
    }
});

window.addEventListener('htmx:historyRestore', () => {
    initNetworkSphere();
    initAmbientBackground();
    if (document.querySelector('#slider-container')) {
        initSlider();
    }
});