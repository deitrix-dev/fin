/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.go"],
  theme: {
    extend: {
      gridTemplateColumns: {
        'descriptions': 'repeat(2, auto 1fr)'
      }
    },
  },
  plugins: [],
}

