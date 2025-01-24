let slideIndex = 0;

function showNextSlide() {
    const slides = document.querySelectorAll('.highlight-slide');
    const tinyCardsContainer = document.querySelector('#tinyCardsContainer');

    if (slides.length <= 1 || !tinyCardsContainer) return;

    const totalSlides = slides.length;
    const visibleTinyCardsCount = 4; // Always display 4 tiny cards
    const nextSlideIndex = (slideIndex + 1) % totalSlides; // Determine the next highlighted slide index

    // Update the highlighted slide
    slides[slideIndex].classList.add('translate-x-full', 'opacity-0');
    slides[nextSlideIndex].classList.remove('translate-x-full', 'opacity-0');
    slideIndex = nextSlideIndex;

    // Animate the tiny card list
    updateTinyCards(slides, tinyCardsContainer, totalSlides, visibleTinyCardsCount, 'next');
}

function showPrevSlide() {
    const slides = document.querySelectorAll('.highlight-slide');
    const tinyCardsContainer = document.querySelector('#tinyCardsContainer');

    if (slides.length <= 1 || !tinyCardsContainer) return;

    const totalSlides = slides.length;
    const visibleTinyCardsCount = 4; // Always display 4 tiny cards
    const prevSlideIndex = (slideIndex - 1 + totalSlides) % totalSlides; // Determine the previous highlighted slide index

    // Update the highlighted slide
    slides[slideIndex].classList.add('translate-x-full', 'opacity-0');
    slides[prevSlideIndex].classList.remove('translate-x-full', 'opacity-0');
    slideIndex = prevSlideIndex;

    // Animate the tiny card list
    updateTinyCards(slides, tinyCardsContainer, totalSlides, visibleTinyCardsCount, 'prev');
}

function updateTinyCards(slides, tinyCardsContainer, totalSlides, visibleTinyCardsCount, direction) {
    const tinyCards = Array.from(tinyCardsContainer.children); // Current tiny cards

    // Add animation to the existing tiny cards
    tinyCards.forEach((card) => {
        if (direction === 'next') {
            card.style.transform = `translateY(-100%)`; // Move each card up
        } else if (direction === 'prev') {
            card.style.transform = `translateY(100%)`; // Move each card down
        }
        card.style.transition = 'transform 0.5s ease-in-out';
    });

    // After the animation completes, update the tiny cards
    setTimeout(() => {
        // Clear the current cards
        tinyCardsContainer.innerHTML = '';

        // Correctly calculate the `startIndex` for the new tiny card list
        let startIndex;
        if (direction === 'next') {
            // For "next", the list moves forward
            startIndex = (slideIndex + 1) % totalSlides;
        } else if (direction === 'prev') {
            // For "prev", shift backward but keep the tiny list aligned with "next" logic
            startIndex = (slideIndex + 1) % totalSlides;
        }

        // Add the new set of cards
        for (let i = 0; i < visibleTinyCardsCount; i++) {
            const tinyIndex = (startIndex + i) % totalSlides;

            const newCard = createTinyCard(slides[tinyIndex]);

            // Animate the new card's entry
            if (direction === 'prev' && i === 0) {
                newCard.style.transform = 'translateY(-100%)'; // Start new card from above
                newCard.style.opacity = '0';
                newCard.style.transition = 'transform 0.4s ease-in-out, opacity 0.5s ease-in-out';
                setTimeout(() => {
                    newCard.style.transform = 'translateY(0)';
                    newCard.style.opacity = '1';
                }, 50);
            } else if (direction === 'next' && i === visibleTinyCardsCount - 1) {
                newCard.style.transform = 'translateY(100%)'; // Start new card from below
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

// Attach event listeners to navigation buttons
document.addEventListener('click', (event) => {
    const prevButton = event.target.closest('#prevSlide');
    const nextButton = event.target.closest('#nextSlide');

    if (prevButton) showPrevSlide();
    if (nextButton) showNextSlide();
});







document.addEventListener('DOMContentLoaded', () => {
    setInterval(() => {
        console.log("Automatic slide triggered");
        showNextSlide();
    }, 15000);
});


