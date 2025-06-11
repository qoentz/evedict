// /static/js/main.js

import { initHeader } from './modules/header.js';
import { convertTimestampsToLocal } from './modules/timestamps.js';
import { initSlider } from './modules/slider.js';
import { initBall, initLogoHistoryRestore } from './modules/logo.js';
import { initNetworkSphere, initAmbientBackground } from "./modules/auxiliary.js";

function isMainPage() {
    return window.location.pathname === '/' || window.location.pathname === '';
}

document.addEventListener('DOMContentLoaded', () => {
    function setFullViewportHeight() {
        const vh = window.innerHeight * 0.01;
        document.documentElement.style.setProperty('--vh', `${vh}px`);
    }
    setFullViewportHeight();
    window.addEventListener('resize', setFullViewportHeight);

    if (!isMainPage()) {
        initNetworkSphere();
        initAmbientBackground();
    }

    initHeader();
    initBall();
    initLogoHistoryRestore();
    convertTimestampsToLocal();
});

function handlePageChange() {
    initNetworkSphere();
    convertTimestampsToLocal();

    if (document.querySelector('#slider-container')) {
        initSlider();
    }
}

document.addEventListener('htmx:afterSwap', (event) => {
    initAmbientBackground();
    handlePageChange();
});

window.addEventListener('htmx:historyRestore', () => {
    handlePageChange();
    if (document.querySelector('#progress-bar-container')) {
        setTimeout(() => {
            initAmbientBackground();
        }, 50);
    } else {
        initAmbientBackground();
    }
});