/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./**/*.go"],
  theme: {
    extend: {
      width: {
        '7/10': '70%',
        '3/10': '30%',
      },
      fontFamily: {
        'palatino': ['Palatino', 'Palatino Linotype', 'Book Antiqua', 'serif'],
      },
      animation: {
        'fade-slide-in': 'fadeSlideIn 1.5s ease-out forwards',
        'letter-slide-in': 'letterSlideIn 0.5s ease-out forwards',
        'spin-slow': 'rotation 15s linear infinite', // Re-add the spinning animation
      },
      keyframes: {
        fadeSlideIn: {
          '0%': {
            opacity: 0,
            transform: 'translateX(-20px)',
          },
          '100%': {
            opacity: 1,
            transform: 'translateX(0)',
          },
        },
        letterSlideIn: {
          '0%': {
            transform: 'translateY(2px)', // Smaller vertical movement
            opacity: 0, // Optional: Start faded if desired
          },
          '100%': {
            transform: 'translateY(0)', // End at the original position
            opacity: 1, // Ensure full opacity at the end
          },
        },
        rotation: {
          from: { transform: 'rotate(0deg)' },
          to: { transform: 'rotate(359deg)' },
        },
      },
    },
  },
  plugins: [],
};


