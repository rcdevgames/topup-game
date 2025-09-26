/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'primary-dark': '#35374B',
        'primary-medium': '#344955',
        'primary-light': '#50727B',
        'success': '#78A083',
        'text-primary': '#1F2937',
        'text-secondary': '#6B7280',
      },
    },
  },
  plugins: [],
}