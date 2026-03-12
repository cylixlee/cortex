# Alert - VAlert

`v-alert` is an alert/banner component used to display important information to users.

## Basic Usage

```vue
<v-alert type="info">This is an info alert</v-alert>
<v-alert type="success">This is a success alert</v-alert>
<v-alert type="warning">This is a warning alert</v-alert>
<v-alert type="error">This is an error alert</v-alert>
```

## Common Props

| Prop        | Description       | Default            |
| ----------- | ----------------- | ------------------ |
| `type`      | Alert type        | info               |
| `title`     | Title             | -                  |
| `text`      | Content text      | -                  |
| `variant`   | Style variant     | default            |
| `icon`      | Icon              | Auto based on type |
| `closable`  | Closable          | false              |
| `prominent` | Prominent display | false              |

## Types

```vue
<!-- Info -->
<v-alert type="info" title="Info Title">Info message</v-alert>

<!-- Success -->
<v-alert type="success" title="Success Title">Success message</v-alert>

<!-- Warning -->
<v-alert type="warning" title="Warning Title">Warning message</v-alert>

<!-- Error -->
<v-alert type="error" title="Error Title">Error message</v-alert>
```

## Variants

### Default (elevated)

```vue
<v-alert type="info">Elevated (default)</v-alert>
```

### Flat

```vue
<v-alert type="info" variant="flat">Flat variant</v-alert>
<v-alert type="success" variant="flat">Success flat</v-alert>
<v-alert type="warning" variant="flat">Warning flat</v-alert>
<v-alert type="error" variant="flat">Error flat</v-alert>
```

### Tonal

```vue
<v-alert type="info" variant="tonal">Tonal variant</v-alert>
<v-alert type="success" variant="tonal">Success tonal</v-alert>
<v-alert type="warning" variant="tonal">Warning tonal</v-alert>
<v-alert type="error" variant="tonal">Error tonal</v-alert>
```

### Outlined

```vue
<v-alert type="info" variant="outlined">Outlined variant</v-alert>
<v-alert type="success" variant="outlined">Success outlined</v-alert>
<v-alert type="warning" variant="outlined">Warning outlined</v-alert>
<v-alert type="error" variant="outlined">Error outlined</v-alert>
```

### Text

```vue
<v-alert type="info" variant="text">Text variant</v-alert>
<v-alert type="success" variant="text">Success text</v-alert>
<v-alert type="warning" variant="text">Warning text</v-alert>
<v-alert type="error" variant="text">Error text</v-alert>
```

## Icons

### Auto Icon

```vue
<v-alert type="info">Auto icon</v-alert>
<v-alert type="success">Auto icon</v-alert>
<v-alert type="warning">Auto icon</v-alert>
<v-alert type="error">Auto icon</v-alert>
```

### Custom Icon

```vue
<v-alert type="info" icon="mdi-information">Custom icon</v-alert>
<v-alert type="success" icon="mdi-check-circle">Custom icon</v-alert>
<v-alert type="warning" icon="mdi-alert">Custom icon</v-alert>
<v-alert type="error" icon="mdi-close-circle">Custom icon</v-alert>
```

### Remove Icon

```vue
<v-alert type="info" :icon="false">No icon</v-alert>
```

### Prominent Icon

```vue
<v-alert type="info" prominent icon="mdi-information" variant="tonal">
  Prominent icon
</v-alert>
```

## Closable

### Basic Usage

```vue
<v-alert closable>Closable alert</v-alert>
```

### Custom Close Icon

```vue
<v-alert closable close-icon="mdi-close-circle">
  Custom close icon
</v-alert>
```

### Listen to Close Event

```vue
<template>
  <v-alert
    v-model="alert"
    closable
    title="Closable Alert"
    @click:close="handleClose"
  >
    This alert can be closed
  </v-alert>
</template>

<script setup>
import { ref } from 'vue'

const alert = ref(true)

const handleClose = () => {
  console.log('Alert closed!')
}
</script>
```

### Manual Control

```vue
<template>
  <v-btn v-if="!alert" @click="alert = true">Show Alert</v-btn>
  
  <v-alert
    v-model="alert"
    closable
    type="success"
    title="Success"
  >
    Operation completed successfully
  </v-alert>
</template>

<script setup>
import { ref } from 'vue'
const alert = ref(true)
</script>
```

## Title and Text

### Using title Prop

```vue
<v-alert type="info" title="Info Title">
  This is the alert content
</v-alert>
```

### Using text Prop

```vue
<v-alert type="success" text="Simple text content"></v-alert>
```

### Complete Example

```vue
<v-alert type="warning" title="Warning Title" variant="tonal">
  <template v-slot:text>
    This is the detailed message content.
    It can span multiple lines.
  </template>
</v-alert>
```

## Transitions

```vue
<v-alert type="info" transition="scale-transition">Scale transition</v-alert>
<v-alert type="success" transition="slide-x-transition">Slide transition</v-alert>
<v-alert type="error" transition="fade-transition">Fade transition</v-alert>
```

## Combined Examples

```vue
<template>
  <v-container>
    <v-alert
      type="success"
      variant="tonal"
      closable
      title="Changes Saved"
      class="mb-4"
    >
      Your changes have been saved successfully.
    </v-alert>

    <v-alert
      type="error"
      variant="outlined"
      icon="mdi-alert-circle"
      class="mb-4"
    >
      <strong>Error:</strong> Unable to connect to server.
    </v-alert>

    <v-alert
      type="info"
      prominent
      title="System Update"
      class="mb-4"
    >
      A new version is available. Please refresh the page.
    </v-alert>
  </v-container>
</template>
```

## Borders

```vue
<v-alert type="info" border="start">Start border</v-alert>
<v-alert type="success" border="end">End border</v-alert>
<v-alert type="warning" border="top">Top border</v-alert>
<v-alert type="error" border="bottom">Bottom border</v-alert>
```

## Color Overrides

```vue
<!-- Custom color -->
<v-alert color="purple" variant="flat">Custom purple color</v-alert>

<!-- Combined -->
<v-alert type="warning" color="deep-orange" variant="tonal">
  Warning with custom color
</v-alert>
```
