# CSS Layers

CSS Layers is a core new feature introduced in Vuetify 4, providing a modern way to organize and manage CSS.

## What are CSS Layers

CSS Layers (CSS @layer) is a CSS specification that allows developers to organize styles into named layers and control the priority between these layers.

```css
@layer base {
  button { ... }
}

@layer components {
  .card { ... }
}
```

Vuetify 4 organizes all styles into 5 top-level layers:

## 5 Top-Level Layers

| Layer                | Description      | Usage                         |
| -------------------- | ---------------- | ----------------------------- |
| `vuetify-core`       | Core base styles | CSS variables, reset styles   |
| `vuetify-components` | Component styles | All Vuetify component styles  |
| `vuetify-overrides`  | Override layer   | Default style override points |
| `vuetify-utilities`  | Utility classes  | Utility class styles          |
| `vuetify-final`      | Final layer      | Highest priority              |

## Coexisting with Third-Party CSS

CSS Layers allow Vuetify to safely coexist with third-party CSS like TailwindCSS.

### Integrating with TailwindCSS

```css
/* main.css */
@layer tailwind-base {
  @tailwind base;
}

@tailwind components;

@layer tailwind-utilities {
  @tailwind utilities;
}
```

### Adjusting Layer Order

If you need to adjust priority, you can use `@layer`:

```css
/* Make Vuetify utilities higher priority than custom utilities */
@layer vuetify-utilities, custom-utilities;
```

## Overriding Vuetify Styles

### Using @layer

```css
@layer vuetify-components {
  .v-btn {
    text-transform: none;
    letter-spacing: normal;
  }
}
```

### Using More Specific Selectors

```css
/* Use more specific selectors in components */
.v-application .v-btn {
  font-weight: 600;
}
```

### Using CSS Variable Overrides

```css
/* Override button height globally */
.v-btn {
  --v-btn-height: 40px;
}
```

## Common Override Examples

### Removing Border Radius

```css
@layer vuetify-components {
  .v-card {
    border-radius: 0;
  }
}
```

### Custom Theme Color

```css
:root {
  --v-theme-primary: #6200ee;
}
```

### Overriding Specific Component Variants

```css
@layer vuetify-components {
  .v-btn--variant-flat {
    background-color: rgb(var(--v-theme-surface));
  }
}
```

## Debugging Layer Order

View current layer order:

```css
/* View in browser console */
@layer;
```

## Differences from SCSS

Vuetify 4 migrated from SCSS to pure CSS, using CSS Layers instead of SCSS's style organization:

| Old (Vuetify 3)   | New (Vuetify 4)            |
| ----------------- | -------------------------- |
| SCSS variables    | CSS variables              |
| SCSS mixins       | CSS @layer                 |
| Dark mode classes | CSS variables + dark class |
