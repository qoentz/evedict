// /static/js/main.js

import { initHeader } from './modules/header.js';
import { initMenu } from './modules/menu.js';
import { convertTimestampsToLocal } from './modules/timestamps.js';
import { toggleSources } from './modules/sources.js';
import { initSlider } from './modules/slider.js';

document.addEventListener('DOMContentLoaded', () => {
    initHeader();
    convertTimestampsToLocal();
});

document.addEventListener('htmx:afterSwap', (event) => {
    initSlider();
    initMenu()
});
