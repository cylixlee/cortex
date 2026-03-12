# Chips - VChip

`v-chip` is a compact component used to display information, filter content, or trigger actions.

## Basic Usage

```vue
<v-chip>Default Chip</v-chip>
<v-chip label>Label Chip</v-chip>
<v-chip closable>Closable Chip</v-chip>
```

## Common Props

| Prop       | Description         | Default |
| ---------- | ------------------- | ------- |
| `label`    | Rounded label style | false   |
| `closable` | Show close button   | false   |
| `color`    | Chip color          | default |
| `variant`  | Style variant       | flat    |
| `size`     | Chip size           | default |
| `prepend`  | Icon before text    | -       |
| `append`   | Icon after text     | -       |

## Types / Colors

```vue
<!-- Colors -->
<v-chip color="primary">Primary</v-chip>
<v-chip color="secondary">Secondary</v-chip>
<v-chip color="error">Error</v-chip>
<v-chip color="warning">Warning</v-chip>
<v-chip color="info">Info</v-chip>
<v-chip color="success">Success</v-chip>
```

## Variants

```vue
<!-- Default (flat) -->
<v-chip color="primary">Flat</v-chip>

<!-- Outlined -->
<v-chip color="primary" variant="outlined">Outlined</v-chip>

<!-- Tonal -->
<v-chip color="primary" variant="tonal">Tonal</v-chip>

<!-- Text -->
<v-chip color="primary" variant="text">Text</v-chip>
```

## Icons

### Prepend Icon

```vue
<v-chip prepend="mdi-account">User</v-chip>
<v-chip prepend="mdi-check" color="success">Verified</v-chip>
```

### Append Icon

```vue
<v-chip append="mdi-close-circle">Filter</v-chip>
```

### Icon Only

```vue
<v-chip icon="mdi-magnify" variant="text"></v-chip>
```

## Clickable Chips

```vue
<!-- Clickable -->
<v-chip @click="handleClick">Clickable</v-chip>

<!-- Filter chip -->
<v-chip
  :model-value="selected"
  @click:close="selected = false"
  filter
>
  Filter Selected
</v-chip>
```

## Chip Group - VChipGroup

Used to group related chips with selection behavior.

```vue
<template>
  <v-chip-group v-model="selected" mandatory>
    <v-chip value="all">All</v-chip>
    <v-chip value="active">Active</v-chip>
    <v-chip value="inactive">Inactive</v-chip>
  </v-chip-group>
</template>

<script setup>
import { ref } from 'vue'
const selected = ref('all')
</script>
```

### With Column

```vue
<v-chip-group v-model="selected" column>
  <v-chip value="s" color="primary">Small</v-chip>
  <v-chip value="m" color="secondary">Medium</v-chip>
  <v-chip value="l" color="error">Large</v-chip>
</v-chip-group>
```

### With Icons

```vue
<v-chip-group v-model="selected" mandatory>
  <v-chip value="phone" prepend-icon="mdi-phone">Phone</v-chip>
  <v-chip value="tablet" prepend-icon="mdi-tablet">Tablet</v-chip>
  <v-chip value="laptop" prepend-icon="mdi-laptop">Laptop</v-chip>
</v-chip-group>
```

## Combined Examples

```vue
<template>
  <v-container>
    <div class="mb-4">
      <h3 class="mb-2">Status Chips</h3>
      <v-chip color="success" prepend="mdi-check-circle">Active</v-chip>
      <v-chip color="error" prepend="mdi-close-circle" class="ml-2">Inactive</v-chip>
      <v-chip color="warning" prepend="mdi-clock" class="ml-2">Pending</v-chip>
    </div>

    <div class="mb-4">
      <h3 class="mb-2">Filter Chips</h3>
      <v-chip-group v-model="filters" multiple>
        <v-chip filter>Marketing</v-chip>
        <v-chip filter>Development</v-chip>
        <v-chip filter>Design</v-chip>
      </v-chip-group>
    </div>

    <div class="mb-4">
      <h3 class="mb-2">Closable Chips</h3>
      <v-chip closable class="mr-2 mb-2">Tag 1</v-chip>
      <v-chip closable color="primary" class="mr-2 mb-2">Tag 2</v-chip>
      <v-chip closable color="error" variant="outlined" class="mr-2 mb-2">Tag 3</v-chip>
    </div>

    <div>
      <h3 class="mb-2">Sizes</h3>
      <v-chip size="x-small">Extra Small</v-chip>
      <v-chip size="small" class="ml-2">Small</v-chip>
      <v-chip class="ml-2">Default</v-chip>
      <v-chip size="large" class="ml-2">Large</v-chip>
      <v-chip size="x-large" class="ml-2">Extra Large</v-chip>
    </div>
  </v-container>
</template>

<script setup>
import { ref } from 'vue'
const filters = ref([])
</script>
```

## Usage with Data Table

```vue
<template>
  <v-data-table :headers="headers" :items="items">
    <template v-slot:item.tags="{ item }">
      <v-chip
        v-for="tag in item.tags"
        :key="tag"
        size="small"
        class="mr-1"
      >
        {{ tag }}
      </v-chip>
    </template>
    
    <template v-slot:item.status="{ item }">
      <v-chip
        :color="item.status === 'published' ? 'success' : 'warning'"
        size="small"
      >
        {{ item.status }}
      </v-chip>
    </template>
  </v-data-table>
</template>

<script setup>
import { ref } from 'vue'

const headers = [
  { title: 'Title', key: 'title' },
  { title: 'Tags', key: 'tags' },
  { title: 'Status', key: 'status' }
]

const items = ref([
  { title: 'Article 1', tags: ['news', 'tech'], status: 'published' },
  { title: 'Article 2', tags: ['blog'], status: 'draft' }
])
</script>
```
