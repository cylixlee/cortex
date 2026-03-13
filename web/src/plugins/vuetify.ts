import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import { useTheme } from 'vuetify'

import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

export const statusColors = {
  pending: { bg: 'warning-lighten-5', text: 'warning-darken-3' },
  processing: { bg: 'info-lighten-5', text: 'info-darken-3' },
  completed: { bg: 'success-lighten-5', text: 'success-darken-3' },
  failed: { bg: 'error-lighten-5', text: 'error-darken-3' },
  default: { bg: 'grey-lighten-4', text: 'grey-darken-3' },
  uploading: { bg: 'info-lighten-5', text: 'info-darken-3' },
  extracting: { bg: 'warning-lighten-5', text: 'warning-darken-3' },
  analyzing: { bg: 'purple-lighten-5', text: 'purple-darken-3' },
  generating: { bg: 'pink-lighten-5', text: 'pink-darken-3' },
  unknown: { bg: 'grey-lighten-4', text: 'grey-darken-3' },
}

export const spacing = {
  xs: '4px',
  sm: '8px',
  md: '16px',
  lg: '24px',
  xl: '32px',
  chatMessageGap: '48px',
}

export const layout = {
  pageMaxWidth: '1200px',
  chatMaxWidth: '900px',
  cardPadding: '24px',
  pagePadding: '24px',
}

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
          success: '#2E7D32',
          warning: '#F9A825',
          background: '#FFFFFF',
          surface: '#FFFFFF',
          'on-background': '#1D1B20',
          'on-surface': '#1D1B20',
          sidebar: '#F5F5F5',
          chat: '#FFFFFF',
          muted: '#757575',
          'surface-variant': '#F5F5F5',
          'warning-lighten-5': '#FEF3C7',
          'warning-darken-3': '#92400E',
          'info-lighten-5': '#DBEAFE',
          'info-darken-3': '#1E40AF',
          'success-lighten-5': '#D1FAE5',
          'success-darken-3': '#065F46',
          'error-lighten-5': '#FEE2E2',
          'error-darken-3': '#991B1B',
          'purple-lighten-5': '#E0E7FF',
          'purple-darken-3': '#4338CA',
          'pink-lighten-5': '#FCE7F3',
          'pink-darken-3': '#9D174D',
          'grey-lighten-4': '#F3F4F6',
          'grey-darken-3': '#374151',
          'button-secondary': '#F3F4F6',
          'button-secondary-text': '#374151',
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
