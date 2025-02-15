// /static/js/slider.js

let slideIndex = 0;
let slideInterval = null;

export function initSlider() {
    const prevButton = document.getElementById('prevSlide');
    const nextButton = document.getElementById('nextSlide');
    const tinyCardsContainer = document.getElementById('tinyCardsContainer');
    const slides = document.querySelectorAll('.highlight-slide');

    if (!prevButton || !nextButton || !tinyCardsContainer || slides.length === 0) {
        console.error("Slider elements not found or no slides available!");
        return;
    }

    // Clear any previous interval so it doesn't continue running
    if (slideInterval !== null) {
        clearInterval(slideInterval);
    }

    // Reset the slide index
    slideIndex = 0;

    // Reset all slides to the initial state:
    // Show only the first slide, hide the rest.
    slides.forEach((slide, index) => {
        if (index === 0) {
            slide.classList.remove('translate-x-full', 'opacity-0');
        } else {
            slide.classList.add('translate-x-full', 'opacity-0');
        }
    });

    // Use the onclick property to ensure only one event listener is bound
    prevButton.onclick = showPrevSlide;
    nextButton.onclick = showNextSlide;

    // Initialize the dots indicator based on the new state
    updateDots();

    // Start the automatic slide interval
    slideInterval = setInterval(showNextSlide, 15000);
}

function showNextSlide() {
    const slides = document.querySelectorAll('.highlight-slide');
    const tinyCardsContainer = document.getElementById('tinyCardsContainer');

    if (slides.length <= 1 || !tinyCardsContainer) return;

    const totalSlides = slides.length;
    const visibleTinyCardsCount = 4;
    const nextSlideIndex = (slideIndex + 1) % totalSlides;

    // Transition current slide out and next slide in
    slides[slideIndex].classList.add('translate-x-full', 'opacity-0');
    slides[nextSlideIndex].classList.remove('translate-x-full', 'opacity-0');
    slideIndex = nextSlideIndex;

    updateSlidePointerEvents(slides, slideIndex);

    updateDots();
    updateTinyCards(slides, tinyCardsContainer, totalSlides, visibleTinyCardsCount, 'next');
}

function showPrevSlide() {
    const slides = document.querySelectorAll('.highlight-slide');
    const tinyCardsContainer = document.getElementById('tinyCardsContainer');

    if (slides.length <= 1 || !tinyCardsContainer) return;

    const totalSlides = slides.length;
    const visibleTinyCardsCount = 4;
    const prevSlideIndex = (slideIndex - 1 + totalSlides) % totalSlides;

    slides[slideIndex].classList.add('translate-x-full', 'opacity-0');
    slides[prevSlideIndex].classList.remove('translate-x-full', 'opacity-0');
    slideIndex = prevSlideIndex;

    updateSlidePointerEvents(slides, slideIndex);

    updateDots();
    updateTinyCards(slides, tinyCardsContainer, totalSlides, visibleTinyCardsCount, 'prev');
}

function updateDots() {
    const dotContainer = document.getElementById('slideDots');
    if (!dotContainer) return;
    const dots = dotContainer.children;
    for (let i = 0; i < dots.length; i++) {
        if (i === slideIndex) {
            dots[i].classList.remove('bg-gray-400');
            dots[i].classList.add('bg-white');
        } else {
            dots[i].classList.remove('bg-white');
            dots[i].classList.add('bg-gray-400');
        }
    }
}

function updateTinyCards(slides, tinyCardsContainer, totalSlides, visibleTinyCardsCount, direction) {
    const tinyCards = Array.from(tinyCardsContainer.children);
    tinyCards.forEach((card) => {
        if (direction === 'next') {
            card.style.transform = 'translateY(-100%)';
        } else if (direction === 'prev') {
            card.style.transform = 'translateY(100%)';
        }
        card.style.transition = 'transform 0.5s ease-in-out';
    });

    setTimeout(() => {
        tinyCardsContainer.innerHTML = '';

        let startIndex = (slideIndex + 1) % totalSlides;

        for (let i = 0; i < visibleTinyCardsCount; i++) {
            const tinyIndex = (startIndex + i) % totalSlides;
            const newCard = createTinyCard(slides[tinyIndex]);

            if (direction === 'prev' && i === 0) {
                newCard.style.transform = 'translateY(-100%)';
                newCard.style.opacity = '0';
                newCard.style.transition = 'transform 0.4s ease-in-out, opacity 0.5s ease-in-out';
                setTimeout(() => {
                    newCard.style.transform = 'translateY(0)';
                    newCard.style.opacity = '1';
                }, 50);
            } else if (direction === 'next' && i === visibleTinyCardsCount - 1) {
                newCard.style.transform = 'translateY(100%)';
                newCard.style.opacity = '0';
                newCard.style.transition = 'transform 0.4s ease-in-out, opacity 0.5s ease-in-out';
                setTimeout(() => {
                    newCard.style.transform = 'translateY(0)';
                    newCard.style.opacity = '1';
                }, 50);
            }

            tinyCardsContainer.appendChild(newCard);
        }
    }, 500); // Match the animation duration
}

function createTinyCard(slide) {
    // Get the same URL from the original clickable element in the slide.
    const clickableEl = slide.querySelector('[hx-get]');
    const hxGet = clickableEl ? clickableEl.getAttribute('hx-get') : '';
    const imgSrc = slide.querySelector('img').src;
    const headline = slide.querySelector('h2').textContent;

    const card = document.createElement('div');
    card.className =
        'cursor-pointer relative w-full h-28 bg-gray-700 rounded-lg overflow-hidden shadow-md group snap-start';

    // Set HTMX attributes so the card behaves like the original.
    card.setAttribute('hx-get', hxGet);
    card.setAttribute('hx-trigger', 'click');
    card.setAttribute('hx-target', '#forecast-feed');
    card.setAttribute('hx-swap', 'innerHTML');

    card.innerHTML = `
        <img src="${imgSrc}" alt="Tiny Card Image" class="w-full h-full object-cover opacity-70 group-hover:opacity-50 transition-opacity duration-300"/>
        <div class="absolute inset-0 flex items-center justify-center p-2">
            <p class="text-white text-sm font-bold text-center leading-tight">${headline}</p>
        </div>
    `;

    // Explicitly reset any inline transforms.
    card.style.transform = 'translateY(0)';

    // Re-process the card with HTMX so that its hx-* attributes are bound.
    if (window.htmx) {
        htmx.process(card);
    }

    return card;
}

function updateSlidePointerEvents(slides, activeIndex) {
    slides.forEach((slide, index) => {
        if (index === activeIndex) {
            slide.classList.remove('pointer-events-none');
        } else {
            slide.classList.add('pointer-events-none');
        }
    });
}



