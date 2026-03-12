# Layout System

## VApp - Application Root Component

`v-app` is the root component for all Vuetify applications and must be used as the outermost container.

```vue
<template>
  <v-app>
    <!-- Other Vuetify components -->
  </v-app>
</template>
```

### Common Props

| Prop          | Type    | Description                   |
| ------------- | ------- | ----------------------------- |
| `theme`       | string  | Specifies theme name          |
| `full-height` | boolean | Whether to occupy full height |

## VMain - Main Content Area

`v-main` is the main content area of the application, **must be a direct child of `v-app`**.

```vue
<v-app>
  <v-main>
    <!-- Page content -->
  </v-main>
</v-app>
```

### Common Props

| Prop  | Type   | Description                       |
| ----- | ------ | --------------------------------- |
| `tag` | string | Rendered HTML tag, default 'main' |

## VContainer - Responsive Container

`v-container` provides a responsive layout container, supporting fluid and fixed widths.

```vue
<v-container>
  <!-- Content -->
</v-container>
```

### Common Props

| Prop    | Type    | Description                         |
| ------- | ------- | ----------------------------------- |
| `fluid` | boolean | Enable fluid width (100%)           |
| `class` | string  | Add custom class like `fill-height` |

### Example: Centered Vertically Layout

```vue
<v-container class="fill-height" fluid>
  <v-row justify="center" align="center">
    <v-col cols="12" sm="8" md="4">
      <!-- Centered Content -->
    </v-col>
  </v-row>
</v-container>
```

## Grid System - VRow and VCol

Vuetify 4's grid system is based on Flexbox.

### Basic Usage

```vue
<v-container>
  <v-row>
    <v-col cols="12">Full Width</v-col>
  </v-row>
  <v-row>
    <v-col cols="6">50% Width</v-col>
    <v-col cols="6">50% Width</v-col>
  </v-row>
  <v-row>
    <v-col cols="4">33%</v-col>
    <v-col cols="4">33%</v-col>
    <v-col cols="4">33%</v-col>
  </v-row>
</v-container>
```

### Responsive Breakpoints

| Breakpoint | Min Width | Typical Device |
| ---------- | --------- | -------------- |
| `xs`       | 0px       | Phone          |
| `sm`       | 600px     | Tablet         |
| `md`       | 960px     | Small Laptop   |
| `lg`       | 1280px    | Desktop        |
| `xl`       | 1920px    | Large Screen   |
| `xxl`      | 2560px    | Extra Large    |

### Responsive Column Width

```vue
<v-row>
  <!-- Phone: 12 columns (100%), Tablet: 6 columns (50%), Desktop: 4 columns (33%) -->
  <v-col cols="12" sm="6" md="4">
    Responsive Column
  </v-col>
</v-row>
```

### Using gap Property

Vuetify 4 recommends using `gap` property instead of the old `gutter`:

```vue
<v-row gap="4">
  <v-col cols="6">Column 1</v-col>
  <v-col cols="6">Column 2</v-col>
</v-row>

<!-- Responsive gap -->
<v-row gap-md="8">
  <v-col cols="6">Column 1</v-col>
  <v-col cols="6">Column 2</v-col>
</v-row>
```

### Alignment

Use utility classes (recommended) or `align` / `justify` props:

```vue
<!-- Vertical alignment: start, center, end -->
<v-row align="center">
  <v-col>Vertically Centered</v-col>
</v-row>

<!-- Horizontal alignment: start, center, end, space-between, space-around -->
<v-row justify="center">
  <v-col>Horizontally Centered</v-col>
</v-row>
```

### Auto Column Width

When `v-col` does not specify `cols`, it automatically adjusts width based on content:

```vue
<v-row>
  <v-col>Auto Width</v-col>
  <v-col>Auto Width</v-col>
</v-row>
```

### Offsets

```vue
<v-row>
  <v-col cols="6" offset="3">Centered (offset 3 columns)</v-col>
</v-row>
```
