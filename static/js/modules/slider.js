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
    prevButton.onclick = () => {
        showPrevSlide();
        resetAutoSlideTimer(); // Reset timer on manual navigation
    };
    nextButton.onclick = () => {
        showNextSlide();
        resetAutoSlideTimer(); // Reset timer on manual navigation
    };

    // Initialize the dots indicator based on the new state
    updateDots();

    // Start the automatic slide interval
    startAutoSlideTimer();
}

function startAutoSlideTimer() {
    slideInterval = setInterval(showNextSlide, 15000);
}

function resetAutoSlideTimer() {
    // Clear existing interval
    if (slideInterval !== null) {
        clearInterval(slideInterval);
    }
    // Start a new interval
    startAutoSlideTimer();
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

    tinyCards.forEach((cardWrapper) => {
        const reflection = cardWrapper.querySelector('.absolute.-bottom-16');
        if (reflection) {
            reflection.style.opacity = '0';
        }
    });

    tinyCards.forEach((card, index) => {
        if (direction === 'next') {
            if (index === 0) {
                // Top card exits faster
                card.style.transition = 'transform 0.7s ease-in-out';
                card.style.transform = 'translateY(-145%)';
            } else {
                // Middle cards keep the smooth speed
                card.style.transition = 'transform 0.7s ease-in-out';
                card.style.transform = 'translateY(-115%)';
            }
        } else if (direction === 'prev') {
            if (index === tinyCards.length - 1) {
                // Bottom card exits faster
                card.style.transition = 'transform 0.7s ease-in-out';
                card.style.transform = 'translateY(145%)';
            } else {
                // Middle cards keep the smooth speed
                card.style.transition = 'transform 0.7s ease-in-out';
                card.style.transform = 'translateY(115%)';
            }
        }
    });

    setTimeout(() => {
        tinyCardsContainer.innerHTML = '';

        let startIndex = (slideIndex + 1) % totalSlides;

        for (let i = 0; i < visibleTinyCardsCount; i++) {
            const tinyIndex = (startIndex + i) % totalSlides;

            const isSlidingCard =
                (direction === 'prev' && i === 0) ||
                (direction === 'next' && i === visibleTinyCardsCount - 1);

            // Always add reflection to the bottom card (index 3)
            const isBottomCard = i === 3;
            const newCard = createTinyCard(slides[tinyIndex], isBottomCard);

            // Hide reflection initially for bottom card when going previous
            if (direction === 'prev' && isBottomCard) {
                const reflection = newCard.querySelector('.absolute.-bottom-16');
                if (reflection) {
                    const reflectionDiv = reflection.querySelector('div');
                    if (reflectionDiv) {
                        reflectionDiv.style.setProperty('opacity', '0', 'important');
                    }
                }
            }

            if (isSlidingCard) {
                newCard.style.transform = direction === 'prev' ? 'translateY(-100%)' : 'translateY(100%)';
                newCard.style.opacity = '0';
                newCard.style.transition = 'transform 0.4s ease-in-out, opacity 0.5s ease-in-out';
                setTimeout(() => {
                    newCard.style.transform = 'translateY(0)';
                    newCard.style.opacity = '1';
                }, 50);
            }

            tinyCardsContainer.appendChild(newCard);
        }

        // Delayed reflection for bottom card when going previous
        if (direction === 'prev') {
            setTimeout(() => {
                const bottomCard = tinyCardsContainer.children[3]; // Bottom card is at index 3
                if (bottomCard) {
                    const reflection = bottomCard.querySelector('.absolute.-bottom-16');
                    if (reflection) {
                        const reflectionDiv = reflection.querySelector('div');
                        if (reflectionDiv) {
                            reflectionDiv.style.transition = 'opacity 0.4s ease-in-out';
                            reflectionDiv.style.setProperty('opacity', '0.8', 'important');
                        }
                    }
                }
            }, 150); // Delay the reflection appearance
        }
    }, 800);
}

function createTinyCard(slide, addReflection = false) {
    const clickableEl = slide.querySelector('[hx-get]');
    const hxGet = clickableEl ? clickableEl.getAttribute('hx-get') : '';
    const imgSrc = slide.querySelector('img').src;
    const headline = slide.querySelector('h2').textContent;

    const cardWrapper = document.createElement('div');
    cardWrapper.className = 'relative w-full';

    const card = document.createElement('div');
    card.className =
        'cursor-pointer relative w-full h-28 bg-gray-700 rounded-lg overflow-hidden shadow-md group snap-start z-10';

    // Set HTMX attributes
    card.setAttribute('hx-get', hxGet);
    card.setAttribute('hx-trigger', 'click');
    card.setAttribute('hx-target', '#forecast-feed');
    card.setAttribute('hx-swap', 'innerHTML');

    card.innerHTML = `
        <img src="${imgSrc}" alt="" class="w-full h-full object-cover opacity-70 group-hover:opacity-50 transition-opacity duration-300"/>
        <div class="absolute inset-0 flex items-center justify-center p-2">
            <p class="text-white text-sm font-bold text-center leading-tight">${headline}</p>
        </div>
    `;

    // Add bottom reflection if requested
    if (addReflection) {
        const reflection = document.createElement('div');
        reflection.className = 'absolute -bottom-16 -left-4 w-[calc(100%+2rem)] h-16 overflow-visible pointer-events-none z-0';
        const reflectionDiv = document.createElement('div');
        reflectionDiv.className = 'relative w-full h-full scale-y-[-1] opacity-8';

        // Set mask property using setProperty to handle vendor prefixes properly
        reflectionDiv.style.setProperty('mask', 'radial-gradient(ellipse 150% 120% at center top, rgba(0,0,0,0.9) 0%, rgba(0,0,0,0.6) 20%, rgba(0,0,0,0.3) 40%, rgba(0,0,0,0.1) 60%, rgba(0,0,0,0) 80%)');

        reflectionDiv.innerHTML = `<img src="${imgSrc}" alt="" class="w-full h-full object-cover blur-md scale-110"/>`;

        reflection.appendChild(reflectionDiv);
        cardWrapper.appendChild(reflection);
    }

    cardWrapper.appendChild(card);

    // Re-process HTMX
    if (window.htmx) {
        htmx.process(card);
    }

    return cardWrapper;
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
