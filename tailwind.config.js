/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/**/*.{templ,html}"],
  theme: {
    extend: {
      gridTemplateColumns: {
        'descriptions': 'repeat(2, auto 1fr)'
      }
    },
  },
  plugins: [],
}

