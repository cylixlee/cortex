# Theming & Colors

## Color Classes

### Background Color (bg-*)

```vue
<div class="bg-primary">Primary Background</div>
<div class="bg-secondary">Secondary Background</div>
<div class="bg-error">Error Background</div>
<div class="bg-warning">Warning Background</div>
<div class="bg-info">Info Background</div>
<div class="bg-success">Success Background</div>

<!-- Variants -->
<div class="bg-primary-lighten-1">Lighten</div>
<div class="bg-primary-darken-1">Darken</div>
```

### Text Color (text-*)

```vue
<p class="text-primary">Primary Text</p>
<p class="text-secondary">Secondary Text</p>
<p class="text-error">Error Text</p>
<p class="text-warning">Warning Text</p>
<p class="text-info">Info Text</p>
<p class="text-success">Success Text</p>

<!-- White/Black -->
<p class="text-white">White Text</p>
<p class="text-black">Black Text</p>
```

### Color Variants

Use `lighten-1` to `lighten-5` and `darken-1` to `darken-5`:

```vue
<!-- Lighten variants -->
<div class="bg-primary-lighten-2">Lighter</div>

<!-- Darken variants -->
<div class="bg-primary-darken-2">Darker</div>
```

## Default Color Palette

Vuetify comes with the following built-in colors:

| Color Name  | Usage           |
| ----------- | --------------- |
| `primary`   | Primary color   |
| `secondary` | Secondary color |
| `error`     | Error state     |
| `warning`   | Warning state   |
| `info`      | Info state      |
| `success`   | Success state   |

## Dark Mode

### Component Level Dark Mode

Use `dark` prop on components:

```vue
<v-btn dark>Dark Button</v-btn>
<v-card dark>
  <v-card-title>Dark Card</v-card-title>
</v-card>
<v-text-field dark label="Dark Input"></v-text-field>
```

### Global Dark Mode

Add `dark` class to `v-app`:

```vue
<v-app class="dark">
  <!-- Entire app uses dark theme -->
</v-app>
```

### Responsive Dark Mode

```vue
<v-app :class="{ 'dark': isDark }">
  <v-btn @click="isDark = !isDark">Toggle Mode</v-btn>
</v-app>
```

### Override Colors in Dark Mode

```vue
<!-- Custom background in dark mode -->
<v-card class="bg-grey-darken-4">
  Custom Background
</v-card>
```

## Combining Colors with Components

### Button Colors

```vue
<v-btn color="primary">Primary Button</v-btn>
<v-btn color="error">Error Button</v-btn>
<v-btn color="success" variant="outlined">Outlined</v-btn>
```

### Card Colors

```vue
<v-card color="primary" dark>
  <v-card-title>Title</v-card-title>
  <v-card-text>Content</v-card-text>
</v-card>
```

### Text Field Colors

```vue
<v-text-field label="Primary" color="primary"></v-text-field>
<v-text-field label="Error" color="error"></v-text-field>
```

### Background Container Example

```vue
<template>
  <v-app>
    <v-main>
      <v-container>
        <!-- Colored background section -->
        <v-sheet class="bg-primary pa-4 mb-4" rounded>
          <h2 class="text-white">Primary Section</h2>
        </v-sheet>
        
        <v-sheet class="bg-success-lighten-4 pa-4" rounded>
          <h2 class="text-success">Light Green Section</h2>
        </v-sheet>
      </v-container>
    </v-main>
  </v-app>
</template>
```

## Common Color Combinations

```vue
<!-- Success state -->
<div class="bg-success pa-4 rounded">
  <span class="text-white">✓ Success Message</span>
</div>

<!-- Error state -->
<div class="bg-error pa-4 rounded">
  <span class="text-white">✗ Error Message</span>
</div>

<!-- Info -->
<div class="bg-info pa-4 rounded">
  <span class="text-white">ℹ Info</span>
</div>

<!-- Warning state -->
<div class="bg-warning pa-4 rounded">
  <span class="text-black">⚠ Warning</span>
</div>
```
