import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'cortexTheme',
    themes: {
      cortexTheme: {
        dark: false,
        colors: {
          primary: '#6750A4',
          secondary: '#625B71',
          error: '#B3261E',
          info: '#3B82F6',
          background: '#FFFFFF',
          surface: '#FFFFFF',
          'on-background': '#1D1B20',
          'on-surface': '#1D1B20',
          sidebar: '#F5F5F5',
          chat: '#FFFFFF',
          muted: '#757575',
          'surface-variant': '#F5F5F5',
        },
      },
    },
  },
  defaults: {
    VBtn: {
      variant: 'flat',
      rounded: 'lg',
    },
    VCard: {
      rounded: 'lg',
      elevation: 0,
    },
    VTextField: {
      variant: 'outlined',
      density: 'comfortable',
      rounded: 'lg',
    },
    VTextarea: {
      variant: 'outlined',
      rounded: 'lg',
    },
  },
})
