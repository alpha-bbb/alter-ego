/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  safelist: [
    "red-mark",
    "blue-mark",
    "yellow-mark",
    "what",
    "example",
    "solution",
    "how",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};
