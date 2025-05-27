// /static/js/main.js

import { initHeader } from './modules/header.js';
import { convertTimestampsToLocal } from './modules/timestamps.js';
import { initSlider } from './modules/slider.js';
import { initBall, initLogoHistoryRestore } from './modules/logo.js';
import { initNetworkSphere } from "./modules/auxiliary.js";

document.addEventListener('DOMContentLoaded', () => {
    function setFullViewportHeight() {
        const vh = window.innerHeight * 0.01;
        document.documentElement.style.setProperty('--vh', `${vh}px`);
    }
    setFullViewportHeight();
    window.addEventListener('resize', setFullViewportHeight);

    initNetworkSphere();
    initHeader();
    initBall();
    initLogoHistoryRestore();
    convertTimestampsToLocal();
});

document.addEventListener('htmx:afterSwap', (event) => {
    initNetworkSphere();
    const target = event.detail.target;
    if (target && target.querySelector('#slider-container')) {
        initSlider();
    }
});

window.addEventListener('htmx:historyRestore', () => {
    initNetworkSphere();
    if (document.querySelector('#slider-container')) {
        initSlider();
    }
});