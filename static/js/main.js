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

    initNetworkSphere();
    initAmbientBackground();

    initHeader();
    initBall();
    initLogoHistoryRestore();
    convertTimestampsToLocal();
});

function handlePageChange() {
    setTimeout(() => {
        initNetworkSphere();
        initAmbientBackground();
        convertTimestampsToLocal();

        if (document.querySelector('#slider-container')) {
            initSlider();
        }
    }, 50);
}

document.addEventListener('htmx:afterSwap', (event) => {
    handlePageChange();
});

window.addEventListener('htmx:historyRestore', () => {
    handlePageChange();
});