const mix = require('laravel-mix')
const tailwindcss = require('tailwindcss')

mix.setPublicPath('public')

mix.js('resources/js/app.js', 'js')
  .sass('resources/scss/app.scss', 'css')
  .options({
    processCssUrls: false,
    postCss: [tailwindcss('./tailwind.config.js')],
  })

mix.webpackConfig({
  output: {
    chunkFilename: 'js/[name].js',
  },
})