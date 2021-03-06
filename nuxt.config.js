
export default {
  mode: 'spa',
  srcDir: 'src/',
  generate: {
    dir: 'nuxt-dist',
  },
  env: {
    serverPort: process.env.NODE_ENV === 'dev' ? '10098' : null
  },
  head: {
    title: process.env.npm_package_name || '',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: process.env.npm_package_description || '' }
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
    ],
    bodyAttrs: {
      style: 'background-color: #000000;',
    }
  },
  /*
  ** Customize the progress-bar color
  */
  loading: { color: '#fff' },
  /*
  ** Global CSS
  */
  css: [
  ],
  /*
  ** Plugins to load before mounting the App
  */
  plugins: [
  ],
  /*
  ** Nuxt.js dev-modules
  */
  buildModules: [
    '@nuxt/typescript-build'
  ],
  /*
  ** Nuxt.js modules
  */
  modules: [
  ],
  /*
  ** Build configuration
  */
 build: {
  /*
  ** You can extend webpack config here
  */
  extend(config, ctx) {
    if (ctx.isClient) {
      config.devtool = '#source-map'
    }


  }
},
extensions: ['ts'],
}
