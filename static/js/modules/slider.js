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

    // Prevent multiple intervals
    if (slideInterval !== null) {
        clearInterval(slideInterval);
    }

    // Attach event listeners to navigation buttons
    prevButton.addEventListener('click', showPrevSlide);
    nextButton.addEventListener('click', showNextSlide);

    // Initialize dots (in case any adjustments are needed)
    updateDots();

    // Start the automatic slide
    slideInterval = setInterval(showNextSlide, 15000);
}

function showNextSlide() {
    const slides = document.querySelectorAll('.highlight-slide');
    const tinyCardsContainer = document.getElementById('tinyCardsContainer');

    if (slides.length <= 1 || !tinyCardsContainer) return;

    const totalSlides = slides.length;
    const visibleTinyCardsCount = 4; // Always display 4 tiny cards
    const nextSlideIndex = (slideIndex + 1) % totalSlides; // Determine the next slide index

    // Animate slide transition
    slides[slideIndex].classList.add('translate-x-full', 'opacity-0');
    slides[nextSlideIndex].classList.remove('translate-x-full', 'opacity-0');
    slideIndex = nextSlideIndex;

    // Update the active dot indicator
    updateDots();

    // Animate the tiny card list
    updateTinyCards(slides, tinyCardsContainer, totalSlides, visibleTinyCardsCount, 'next');
}

function showPrevSlide() {
    const slides = document.querySelectorAll('.highlight-slide');
    const tinyCardsContainer = document.getElementById('tinyCardsContainer');

    if (slides.length <= 1 || !tinyCardsContainer) return;

    const totalSlides = slides.length;
    const visibleTinyCardsCount = 4; // Always display 4 tiny cards
    const prevSlideIndex = (slideIndex - 1 + totalSlides) % totalSlides; // Determine the previous slide index

    // Animate slide transition
    slides[slideIndex].classList.add('translate-x-full', 'opacity-0');
    slides[prevSlideIndex].classList.remove('translate-x-full', 'opacity-0');
    slideIndex = prevSlideIndex;

    // Update the active dot indicator
    updateDots();

    // Animate the tiny card list
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
    const tinyCards = Array.from(tinyCardsContainer.children); // Current tiny cards

    // Animate the existing tiny cards
    tinyCards.forEach((card) => {
        if (direction === 'next') {
            card.style.transform = `translateY(-100%)`; // Move card up
        } else if (direction === 'prev') {
            card.style.transform = `translateY(100%)`; // Move card down
        }
        card.style.transition = 'transform 0.5s ease-in-out';
    });

    // After the animation completes, update the tiny cards
    setTimeout(() => {
        // Clear the current cards
        tinyCardsContainer.innerHTML = '';

        // Calculate the startIndex for the new tiny card list
        let startIndex;
        if (direction === 'next') {
            startIndex = (slideIndex + 1) % totalSlides;
        } else if (direction === 'prev') {
            startIndex = (slideIndex + 1) % totalSlides;
        }

        // Add the new set of cards
        for (let i = 0; i < visibleTinyCardsCount; i++) {
            const tinyIndex = (startIndex + i) % totalSlides;
            const newCard = createTinyCard(slides[tinyIndex]);

            // Animate the new card's entry
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
    }, 500); // Match the duration of the animation
}

function createTinyCard(slide) {
    const imgSrc = slide.querySelector('img').src;
    const headline = slide.querySelector('h2').textContent;

    const card = document.createElement('div');
    card.className = 'relative w-full h-28 bg-black rounded-lg overflow-hidden shadow-md group';
    card.innerHTML = `
        <img src="${imgSrc}" alt="Tiny Card Image" class="w-full h-full object-cover opacity-70 group-hover:opacity-50 transition-opacity duration-300"/>
        <div class="absolute inset-0 flex items-center justify-center p-2">
            <p class="text-white text-sm font-bold text-center leading-tight">${headline}</p>
        </div>
    `;
    return card;
}


