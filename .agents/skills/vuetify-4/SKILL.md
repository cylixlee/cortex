---
name: vuetify-4
description: Build Vue 3 applications with Vuetify 4 component library. Use when creating new Vuetify 4 projects, adding UI components, layouts, forms, or theming. Assumes environment is already configured with Vuetify 4.
license: MIT
---

# Vuetify 4 Agent Skill

Vuetify is a Material Design component library based on Vue 3, providing rich pre-built components and utility classes.

## Basic Concepts

### Application Structure

Basic structure of a Vuetify application:

```vue
<template>
  <v-app>
    <v-app-bar title="My App"></v-app-bar>
    <v-navigation-drawer></v-navigation-drawer>
    <v-main>
      <v-container>
        <!-- Page content -->
      </v-container>
    </v-main>
  </v-app>
</template>
```

### Core Layout Components

| Component         | Usage                                                                   |
| ----------------- | ----------------------------------------------------------------------- |
| `v-app`           | Application root component, parent container for all Vuetify components |
| `v-main`          | Main content area                                                       |
| `v-container`     | Responsive container                                                    |
| `v-row` / `v-col` | Grid system                                                             |

See [references/layout.md](references/layout.md) for detailed usage.

## Utility Classes Quick Reference

Vuetify provides rich utility classes:

### Spacing (ma-, pa-, mx-, my-)
```vue
<div class="ma-4 pa-2 mx-auto my-4">Spacing</div>
```

### Alignment (justify-, align-)
```vue
<v-row justify="center" align="center">
  <v-col>Centered Content</v-col>
</v-row>
```

### Display Control (d-*)
```vue
<div class="d-none d-md-flex">Desktop Display</div>
```

See [references/utilities.md](references/utilities.md) for detailed utility classes.

## Component Quick Reference

### Buttons & Cards
- `v-btn` - Button component
- `v-card` - Card component

### Dialogs & Notifications
- `v-dialog` - Modal dialog
- `v-snackbar` - Toast notification
- `v-alert` - Alert banner

### Lists
- `v-list` - List component

### Tabs & Chips
- `v-tabs` / `v-tab` - Tab navigation
- `v-chip` - Chip/tag component
- `v-chip-group` - Chip group

### Form Components
- `v-text-field` - Text input
- `v-select` - Dropdown select
- `v-switch` / `v-checkbox` - Switch/Checkbox
- `v-form` - Form validation

### Navigation Components
- `v-app-bar` - Top navigation bar
- `v-navigation-drawer` - Side navigation
- `v-bottom-navigation` - Bottom navigation

### Data Display
- `v-data-table` - Data table

See [references/components/](references/components/) for detailed usage.

## Theming & Colors

### Color Classes

```vue
<!-- Background -->
<div class="bg-primary">Primary Background</div>
<div class="bg-error">Error Background</div>

<!-- Text color -->
<span class="text-secondary">Secondary Text</span>
<span class="text-success">Success Text</span>
```

### Dark Mode

```vue
<!-- Component level dark mode -->
<v-btn dark>Dark Button</v-btn>

<!-- Global dark mode class -->
<v-app class="dark">
```

See [references/theming.md](references/theming.md) for detailed theming usage.

## CSS Layers Architecture

Vuetify 4 uses CSS Layers for style organization, enabling better compatibility with third-party CSS like TailwindCSS.

See [references/css-layers.md](references/css-layers.md) for detailed usage.

## Quick Examples

### Minimum Runnable Example

```vue
<template>
  <v-app>
    <v-main>
      <v-container>
        <v-card title="Welcome" subtitle="Vuetify 4">
          <v-card-text>
            This is a simple card component
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary">Confirm</v-btn>
            <v-btn variant="text">Cancel</v-btn>
          </v-card-actions>
        </v-card>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref } from 'vue'
</script>
```

### Login Page Example

```vue
<template>
  <v-app>
    <v-main>
      <v-container class="fill-height" fluid>
        <v-row justify="center" align="center">
          <v-col cols="12" sm="8" md="4">
            <v-card>
              <v-card-title class="text-center">Login</v-card-title>
              <v-card-text>
                <v-text-field label="Username" variant="outlined"></v-text-field>
                <v-text-field label="Password" type="password" variant="outlined"></v-text-field>
              </v-card-text>
              <v-card-actions>
                <v-btn block color="primary">Login</v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>
```

See more patterns in [references/patterns/](references/patterns/)

## Common Tasks Quick Reference

| Task              | Keywords                                           |
| ----------------- | -------------------------------------------------- |
| Create layout     | `v-app`, `v-main`, `v-container`                   |
| Add button        | `v-btn`                                            |
| Create card       | `v-card`                                           |
| Form input        | `v-text-field`, `v-select`                         |
| Form validation   | `v-form`, rules                                    |
| Data table        | `v-data-table`                                     |
| Navigation        | `v-app-bar`, `v-navigation-drawer`                 |
| Spacing/Alignment | Utility classes `ma-`, `pa-`, `justify-`, `align-` |
| Color theming     | `bg-*`, `text-*`, `dark`                           |
