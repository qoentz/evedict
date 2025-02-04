// /static/js/main.js

import { initHeader } from './modules/header.js';
import { initMenu } from './modules/menu.js';
import { convertTimestampsToLocal } from './modules/timestamps.js';
import { toggleSources } from './modules/sources.js';
import { initSlider } from './modules/slider.js';
import { initBall} from "./modules/logo.js";

document.addEventListener('DOMContentLoaded', () => {
    initHeader();
    convertTimestampsToLocal();
    initBall()
});

document.addEventListener('htmx:afterSwap', (event) => {
    initSlider();
    initMenu()
});

