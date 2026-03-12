# Snackbar - VSnackbar

`v-snackbar` is a snackbar component used to display brief notification messages.

## Basic Usage

```vue
<template>
  <v-btn @click="snackbar = true">Show Snackbar</v-btn>

  <v-snackbar v-model="snackbar">
    Message text here
  </v-snackbar>
</template>

<script setup>
import { ref } from 'vue'
const snackbar = ref(false)
</script>
```

## Common Props

| Prop         | Description           | Default      |
| ------------ | --------------------- | ------------ |
| `v-model`    | Controls visibility   | false        |
| `text`       | Message text          | -            |
| `timeout`    | Display duration (ms) | 5000         |
| `color`      | Background color      | -            |
| `location`   | Display position      | bottom right |
| `multi-line` | Multi-line mode       | false        |

## Timeout Configuration

```vue
<!-- Default 5 seconds -->
<v-snackbar v-model="snackbar" text="Default 5s"></v-snackbar>

<!-- 3 seconds -->
<v-snackbar v-model="snackbar" text="3 seconds" :timeout="3000"></v-snackbar>

<!-- Never auto-close (0) -->
<v-snackbar v-model="snackbar" text="Persistent" :timeout="-1"></v-snackbar>

<!-- Custom duration -->
<v-snackbar v-model="snackbar" text="10 seconds" :timeout="10000"></v-snackbar>
```

## Multi-line Message

```vue
<v-snackbar v-model="snackbar" multi-line>
  This is a multi-line message that can display longer content.
  This is the second line of content.
</v-snackbar>
```

## With Action Button

```vue
<v-snackbar v-model="snackbar" timeout="-1">
  <span>Message with action</span>
  <template v-slot:actions>
    <v-btn variant="text" @click="snackbar = false">Close</v-btn>
    <v-btn variant="text" @click="handleAction">Action</v-btn>
  </template>
</v-snackbar>
```

## Color Variants

```vue
<!-- Default (theme color) -->
<v-snackbar v-model="snackbar">
  Default color
</v-snackbar>

<!-- Custom color -->
<v-snackbar v-model="snackbar" color="success">
  Success message
</v-snackbar>

<v-snackbar v-model="snackbar" color="error">
  Error message
</v-snackbar>

<v-snackbar v-model="snackbar" color="warning">
  Warning message
</v-snackbar>

<v-snackbar v-model="snackbar" color="info">
  Info message
</v-snackbar>
```

## Using text Prop

```vue
<v-snackbar v-model="snackbar" text="Simple text message"></v-snackbar>
```

## Location

### Bottom Location

```vue
<!-- Bottom left -->
<v-snackbar v-model="snackbar" location="bottom left">...</v-snackbar>

<!-- Bottom center -->
<v-snackbar v-model="snackbar" location="bottom">...</v-snackbar>

<!-- Bottom right (default) -->
<v-snackbar v-model="snackbar" location="bottom right">...</v-snackbar>
```

### Top Location

```vue
<!-- Top left -->
<v-snackbar v-model="snackbar" location="top left">...</v-snackbar>

<!-- Top center -->
<v-snackbar v-model="snackbar" location="top">...</v-snackbar>

<!-- Top right -->
<v-snackbar v-model="snackbar" location="top right">...</v-snackbar>
```

## Complete Example

```vue
<template>
  <v-container>
    <v-row>
      <v-col cols="auto">
        <v-btn color="primary" @click="showSuccess">
          Success
        </v-btn>
      </v-col>
      <v-col cols="auto">
        <v-btn color="error" @click="showError">
          Error
        </v-btn>
      </v-col>
      <v-col cols="auto">
        <v-btn color="warning" @click="showWarning">
          Warning
        </v-btn>
      </v-col>
    </v-row>

    <v-snackbar
      v-model="snackbar"
      :color="color"
      :timeout="timeout"
      location="top"
    >
      {{ text }}
      <template v-slot:actions>
        <v-btn variant="text" @click="snackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref } from 'vue'

const snackbar = ref(false)
const text = ref('')
const color = ref('success')
const timeout = ref(3000)

const showSuccess = () => {
  text.value = 'Operation successful!'
  color.value = 'success'
  snackbar.value = true
}

const showError = () => {
  text.value = 'An error occurred!'
  color.value = 'error'
  snackbar.value = true
}

const showWarning = () => {
  text.value = 'Warning message'
  color.value = 'warning'
  snackbar.value = true
}
</script>
```

## Icon

```vue
<v-snackbar v-model="snackbar" color="success">
  <template v-slot:prepend>
    <v-icon color="white">mdi-check-circle</v-icon>
  </template>
  Success message with icon
</v-snackbar>
```

## Vertical Layout

```vue
<v-snackbar v-model="snackbar" vertical>
  <div class="text-h6">Title</div>
  <div>First line</div>
  <div>Second line</div>
  <template v-slot:actions>
    <v-btn variant="text" @click="snackbar = false">Close</v-btn>
  </template>
</v-snackbar>
```

## Variants

```vue
<!-- Text style (no background) -->
<v-snackbar v-model="snackbar" variant="text">
  Text variant
</v-snackbar>

<!-- Flat style -->
<v-snackbar v-model="snackbar" variant="flat" color="primary">
  Flat variant
</v-snackbar>

<!-- Outlined style -->
<v-snackbar v-model="snackbar" variant="outlined" color="primary">
  Outlined variant
</v-snackbar>
```

## Rounded

```vue
<v-snackbar v-model="snackbar" rounded="lg">
  Rounded snackbar
</v-snackbar>
```
