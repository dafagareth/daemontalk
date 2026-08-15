/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/templates/**/*.templ",
    "./web/templates/**/*.go",
  ],
  theme: {
    extend: {
      fontFamily: {
        mono: ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'monospace'],
      },
    },
  },
  plugins: [],
}
