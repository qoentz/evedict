/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./**/*.go"],
  safelist: [
    'htmx-indicator',
    'htmx-request',
  ],
  theme: {
    extend: {
      width: {
        '7/10': '70%',
        '3/10': '30%',
      },
      fontFamily: {
        palatino: ['Palatino', 'Palatino Linotype', 'Book Antiqua', 'serif'],
        libre: ['"Libre Baskerville"', 'serif'],
        gillsans: ["Gill Sans", "Gill Sans MT", "Calibri", "Trebuchet MS", "sans-serif"],
      },
      animation: {
        'fade-slide-in': 'fadeSlideIn 1.5s ease-out forwards',
        'letter-slide-in': 'letterSlideIn 0.5s ease-out forwards',
        'spin-slow': 'rotation 15s linear infinite',
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
            transform: 'translateY(2px)',
            opacity: 0,
          },
          '100%': {
            transform: 'translateY(0)',
            opacity: 1,
          },
        },
        rotation: {
          from: { transform: 'rotate(0deg)' },
          to: { transform: 'rotate(359deg)' },
        },
      },
      backgroundImage: {
        // Custom diagonal fade gradient
        'fade-to-black': 'linear-gradient(to bottom right, rgba(0, 0, 0, 0) 50%, rgba(0, 0, 0, 0.7) 80%, rgba(0, 0, 0, 1) 100%)',
      },
    },
  },
  plugins: [],
};



