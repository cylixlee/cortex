# Buttons - VBtn

`v-btn` is the most commonly used interactive component in Vuetify.

## Basic Usage

```vue
<v-btn>Default Button</v-btn>
<v-btn color="primary">Primary Button</v-btn>
<v-btn color="error">Error Button</v-btn>
```

## Variants (variant)

```vue
<!-- Default: elevated -->
<v-btn variant="elevated">Elevated</v-btn>

<!-- Flat -->
<v-btn variant="flat">Flat</v-btn>

<!-- Outlined -->
<v-btn variant="outlined">Outlined</v-btn>

<!-- Text (no background) -->
<v-btn variant="text">Text</v-btn>

<!-- Tonal (no shadow) -->
<v-btn variant="tonal">Tonal</v-btn>
```

## Colors

Use `color` prop to set color:

```vue
<v-btn color="primary">Primary</v-btn>
<v-btn color="secondary">Secondary</v-btn>
<v-btn color="error">Error</v-btn>
<v-btn color="warning">Warning</v-btn>
<v-btn color="info">Info</v-btn>
<v-btn color="success">Success</v-btn>
```

Combine with variants:

```vue
<v-btn color="primary" variant="outlined">Primary Outlined</v-btn>
<v-btn color="error" variant="text">Error Text</v-btn>
```

## Sizes

```vue
<v-btn size="x-small">Extra Small</v-btn>
<v-btn size="small">Small</v-btn>
<v-btn>Default</v-btn>
<v-btn size="large">Large</v-btn>
<v-btn size="x-large">Extra Large</v-btn>
```

## Shapes

```vue
<!-- Rounded -->
<v-btn rounded="lg">Rounded</v-btn>

<!-- Circle (suitable for icons) -->
<v-btn icon="mdi-heart" circle></v-btn>

<!-- Square (no rounding) -->
<v-btn rounded="0">Square</v-btn>
```

## Block Button

```vue
<v-btn block>Block Button</v-btn>
```

## Loading State

```vue
<template>
  <v-btn
    :loading="loading"
    @click="loading = true"
  >
    Click to Load
  </v-btn>
</template>

<script setup>
import { ref } from 'vue'
const loading = ref(false)
</script>
```

## Icon Buttons

### Using icon prop

```vue
<v-btn icon="mdi-heart"></v-btn>
<v-btn icon="mdi-heart" color="error"></v-btn>
```

### Text button with icon

```vue
<v-btn>
  <v-icon start>mdi-plus</v-icon>
  Add
</v-btn>

<v-btn>
  Delete
  <v-icon end>mdi-delete</v-icon>
</v-btn>
```

### Circular icon buttons

```vue
<v-btn icon="mdi-magnify" variant="text"></v-btn>
<v-btn icon="mdi-heart" color="error" variant="flat"></v-btn>
```

## Disabled State

```vue
<v-btn disabled>Disabled</v-btn>
<v-btn disabled variant="outlined">Disabled Outline</v-btn>
```

## Combined Examples

```vue
<template>
  <v-container>
    <v-row class="mb-4">
      <v-col>
        <v-btn color="primary">Primary Action</v-btn>
        <v-btn variant="text">Secondary Action</v-btn>
      </v-col>
    </v-row>
    
    <v-row class="mb-4">
      <v-col>
        <v-btn variant="outlined" prepend-icon="mdi-arrow-left">Back</v-btn>
        <v-btn variant="tonal" append-icon="mdi-arrow-right">Next</v-btn>
      </v-col>
    </v-row>
    
    <v-row>
      <v-col>
        <v-btn icon="mdi-magnify" variant="text"></v-btn>
        <v-btn icon="mdi-pencil" variant="flat" color="primary"></v-btn>
        <v-btn icon="mdi-delete" variant="outlined" color="error"></v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>
```
