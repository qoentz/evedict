/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./**/*.go"],
  theme: {
    extend: {
      fontFamily: {
        'palatino': ['Palatino', 'Palatino Linotype', 'Book Antiqua', 'serif'],
      },
      animation: {
        'spin-slow': 'rotation 10s linear infinite',
      },
      keyframes: {
        rotation: {
          from: { transform: 'rotate(0deg)' },
          to: { transform: 'rotate(359deg)' },
        },
      },
    },
  },
  plugins: [],
}

