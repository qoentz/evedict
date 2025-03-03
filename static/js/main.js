// /static/js/main.js

import { initHeader } from './modules/header.js';
import { convertTimestampsToLocal } from './modules/timestamps.js';
import { initSlider } from './modules/slider.js';
import { initBall, initLogoHistoryRestore } from './modules/logo.js';

document.addEventListener('DOMContentLoaded', () => {
    initHeader();
    initBall();
    initLogoHistoryRestore();
    convertTimestampsToLocal();
});

document.addEventListener('htmx:afterSwap', (event) => {
    const target = event.detail.target;
    if (target && target.querySelector('#slider-container')) {
        initSlider();
    }
});

window.addEventListener('htmx:historyRestore', () => {
    if (document.querySelector('#slider-container')) {
        initSlider();
    }
});