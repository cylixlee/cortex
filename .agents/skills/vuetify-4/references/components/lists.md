# List - VList

`v-list` is a list component used to display a vertically arranged collection of items.

## Basic Usage

```vue
<v-list>
  <v-list-item>Item 1</v-list-item>
  <v-list-item>Item 2</v-list-item>
  <v-list-item>Item 3</v-list-item>
</v-list>
```

## Common Props

| Prop      | Description                           | Default |
| --------- | ------------------------------------- | ------- |
| `lines`   | Lines (one/two/three)                 | -       |
| `nav`     | Navigation mode                       | false   |
| `density` | Density (default/compact/comfortable) | default |
| `variant` | Variant (flat/outlined)               | flat    |

## Basic Lists

### Simple List

```vue
<v-list>
  <v-list-item title="Item 1" subtitle="Subtitle 1"></v-list-item>
  <v-list-item title="Item 2" subtitle="Subtitle 2"></v-list-item>
  <v-list-item title="Item 3" subtitle="Subtitle 3"></v-list-item>
</v-list>
```

### Lines

```vue
<v-list lines="one">
  <v-list-item title="Single line"></v-list-item>
</v-list>

<v-list lines="two">
  <v-list-item title="Title" subtitle="Subtitle"></v-list-item>
</v-list>

<v-list lines="three">
  <v-list-item title="Title" subtitle="Subtitle" text="Third line text"></v-list-item>
</v-list>
```

## Navigation List

```vue
<v-list nav>
  <v-list-item prepend-icon="mdi-home" title="Home" value="home"></v-list-item>
  <v-list-item prepend-icon="mdi-account" title="Profile" value="profile"></v-list-item>
  <v-list-item prepend-icon="mdi-cog" title="Settings" value="settings"></v-list-item>
</v-list>
```

### With Links

```vue
<v-list nav>
  <v-list-item
    prepend-icon="mdi-home"
    title="Home"
    value="home"
    href="#home"
  ></v-list-item>
  <v-list-item
    prepend-icon="mdi-account"
    title="Profile"
    value="profile"
    to="/profile"
  ></v-list-item>
</v-list>
```

## List with Icons

```vue
<v-list>
  <v-list-item
    prepend-icon="mdi-inbox"
    title="Inbox"
    subtitle="5 messages"
  ></v-list-item>
  
  <v-list-item
    prepend-icon="mdi-send"
    title="Sent"
    subtitle="10 messages"
  ></v-list-item>
  
  <v-list-item
    prepend-icon="mdi-delete"
    title="Trash"
    subtitle="0 messages"
  ></v-list-item>
</v-list>
```

## Avatar List

```vue
<v-list>
  <v-list-item
    prepend-avatar="https://randomuser.me/api/portraits/men/1.jpg"
    title="John Doe"
    subtitle="john@example.com"
  ></v-list-item>
  
  <v-list-item
    prepend-avatar="https://randomuser.me/api/portraits/women/2.jpg"
    title="Jane Smith"
    subtitle="jane@example.com"
  ></v-list-item>
</v-list>
```

## Nested List

```vue
<v-list>
  <v-list-group value="Users">
    <template v-slot:activator="{ props }">
      <v-list-item v-bind="props" prepend-icon="mdi-account-group" title="Users"></v-list-item>
    </template>
    
    <v-list-item title="Admin"></v-list-item>
    <v-list-item title="Editor"></v-list-item>
    <v-list-item title="Viewer"></v-list-item>
  </v-list-group>
  
  <v-list-group value="Settings">
    <template v-slot:activator="{ props }">
      <v-list-item v-bind="props" prepend-icon="mdi-cog" title="Settings"></v-list-item>
    </template>
    
    <v-list-item title="General"></v-list-item>
    <v-list-item title="Security"></v-list-item>
  </v-list-group>
</v-list>
```

## Selectable List

```vue
<template>
  <v-list selectable>
    <v-list-item
      v-for="item in items"
      :key="item.id"
      :title="item.title"
      :subtitle="item.subtitle"
      :value="item.id"
      :active="selected === item.id"
      @click="selected = item.id"
    ></v-list-item>
  </v-list>
</template>

<script setup>
import { ref } from 'vue'

const selected = ref('1')

const items = [
  { id: '1', title: 'Option 1', subtitle: 'Description 1' },
  { id: '2', title: 'Option 2', subtitle: 'Description 2' },
  { id: '3', title: 'Option 3', subtitle: 'Description 3' },
]
</script>
```

## Multi-Select List

```vue
<template>
  <v-list selectable multiple>
    <v-list-item
      v-for="item in items"
      :key="item.id"
      :title="item.title"
      :value="item.id"
      :active="selected.includes(item.id)"
      @click="toggle(item.id)"
    ></v-list-item>
  </v-list>
  
  <div>Selected: {{ selected }}</div>
</template>

<script setup>
import { ref } from 'vue'

const selected = ref(['1'])

const items = [
  { id: '1', title: 'Option 1' },
  { id: '2', title: 'Option 2' },
  { id: '3', title: 'Option 3' },
]

const toggle = (id) => {
  const index = selected.value.indexOf(id)
  if (index === -1) {
    selected.value.push(id)
  } else {
    selected.value.splice(index, 1)
  }
}
</script>
```

## Compact List

```vue
<v-list density="compact">
  <v-list-item title="Compact Item 1"></v-list-item>
  <v-list-item title="Compact Item 2"></v-list-item>
</v-list>
```

## Variants

### Outlined

```vue
<v-list variant="outlined">
  <v-list-item title="Item 1"></v-list-item>
  <v-list-item title="Item 2"></v-list-item>
</v-list>
```

### Nav Style

```vue
<v-list nav density="compact" rounded>
  <v-list-item prepend-icon="mdi-home" title="Home" value="home" active></v-list-item>
  <v-list-item prepend-icon="mdi-folder" title="Files" value="files"></v-list-item>
  <v-list-item prepend-icon="mdi-star" title="Favorites" value="favorites"></v-list-item>
</v-list>
```

## Complete Example

```vue
<template>
  <v-card width="300">
    <v-list>
      <v-list-subheader>REPORTS</v-list-subheader>
      
      <v-list-item
        v-for="item in reports"
        :key="item.title"
        :prepend-icon="item.icon"
        :title="item.title"
        :subtitle="item.subtitle"
        :value="item.value"
      ></v-list-item>
      
      <v-divider class="my-2"></v-divider>
      
      <v-list-subheader>ARCHIVE</v-list-subheader>
      
      <v-list-item
        v-for="item in archive"
        :key="item.title"
        :prepend-icon="item.icon"
        :title="item.title"
        :value="item.value"
      ></v-list-item>
    </v-list>
  </v-card>
</template>

<script setup>
const reports = [
  { title: 'Dashboard', subtitle: 'Overview', icon: 'mdi-view-dashboard', value: 'dashboard' },
  { title: 'Analytics', subtitle: 'Statistics', icon: 'mdi-chart-line', value: 'analytics' },
  { title: 'Reports', subtitle: 'Generated', icon: 'mdi-file-document', value: 'reports' },
]

const archive = [
  { title: 'Archived', icon: 'mdi-archive', value: 'archived' },
  { title: 'Trash', icon: 'mdi-delete', value: 'trash' },
]
</script>
```

## Grouped List

```vue
<v-list>
  <v-list-group>
    <template v-slot:activator="{ props }">
      <v-list-item v-bind="props" title="Group 1"></v-list-item>
    </template>
    
    <v-list-item title="Item 1.1"></v-list-item>
    <v-list-item title="Item 1.2"></v-list-item>
  </v-list-group>
  
  <v-list-group>
    <template v-slot:activator="{ props }">
      <v-list-item v-bind="props" title="Group 2"></v-list-item>
    </template>
    
    <v-list-item title="Item 2.1"></v-list-item>
    <v-list-item title="Item 2.2"></v-list-item>
  </v-list-group>
</v-list>
```

## List with Actions

```vue
<v-list>
  <v-list-item
    title="Clickable Item"
    subtitle="Click to see actions"
    @click="handleClick"
  >
    <template v-slot:append="{ isActive }">
      <v-btn
        icon="mdi-chevron-right"
        variant="text"
      ></v-btn>
    </template>
  </v-list-item>
</v-list>
```
