# Tabs - VTabs

`v-tabs` is used to organize content into separate views where only one view is visible at a time.

## Basic Usage

```vue
<template>
  <v-card>
    <v-tabs v-model="tab" color="primary">
      <v-tab value="one">Tab 1</v-tab>
      <v-tab value="two">Tab 2</v-tab>
      <v-tab value="three">Tab 3</v-tab>
    </v-tabs>

    <v-card-text>
      <v-tabs-window v-model="tab">
        <v-tabs-window-item value="one">
          Content of Tab 1
        </v-tabs-window-item>
        <v-tabs-window-item value="two">
          Content of Tab 2
        </v-tabs-window-item>
        <v-tabs-window-item value="three">
          Content of Tab 3
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { ref } from 'vue'
const tab = ref('one')
</script>
```

## Common Props

| Prop         | Description             | Default |
| ------------ | ----------------------- | ------- |
| `color`      | Tab indicator color     | primary |
| `align`      | Tab alignment           | start   |
| `centered`   | Center the tabs         | false   |
| `grow`       | Grow tabs to fill width | false   |
| `fixed-tabs` | Fixed width tabs        | false   |
| `height`     | Tab height              | -       |

## Fixed Tabs

```vue
<v-tabs fixed-tabs color="primary">
  <v-tab value="one">Tab 1</v-tab>
  <v-tab value="two">Tab 2</v-tab>
  <v-tab value="three">Tab 3</v-tab>
</v-tabs>
```

## Alignment

### Start (Default)

```vue
<v-tabs align="start">
  <v-tab value="one">Tab 1</v-tab>
  <v-tab value="two">Tab 2</v-tab>
</v-tabs>
```

### Center

```vue
<v-tabs align="center" centered>
  <v-tab value="one">Tab 1</v-tab>
  <v-tab value="two">Tab 2</v-tab>
</v-tabs>
```

### End

```vue
<v-tabs align="end">
  <v-tab value="one">Tab 1</v-tab>
  <v-tab value="two">Tab 2</v-tab>
</v-tabs>
```

## Grow Tabs

```vue
<v-tabs grow color="primary">
  <v-tab value="one">Tab 1</v-tab>
  <v-tab value="two">Tab 2</v-tab>
  <v-tab value="three">Tab 3</v-tab>
</v-tabs>
```

## With Icons

```vue
<template>
  <v-tabs v-model="tab" color="primary">
    <v-tab value="home">
      <v-icon start>mdi-home</v-icon>
      Home
    </v-tab>
    <v-tab value="settings">
      <v-icon start>mdi-cog</v-icon>
      Settings
    </v-tab>
    <v-tab value="about">
      <v-icon start>mdi-information</v-icon>
      About
    </v-tab>
  </v-tabs>
</template>

<script setup>
import { ref } from 'vue'
const tab = ref('home')
</script>
```

### Icons Only

```vue
<v-tabs v-model="tab" color="primary" align="center">
  <v-tab value="home" icon="mdi-home"></v-tab>
  <v-tab value="settings" icon="mdi-cog"></v-tab>
  <v-tab value="about" icon="mdi-information"></v-tab>
</v-tabs>
```

## Dynamic Tabs

```vue
<template>
  <v-card>
    <v-card-title>Dynamic Tabs</v-card-title>
    
    <v-tabs v-model="activeTab" color="primary">
      <v-tab
        v-for="(tab, index) in tabs"
        :key="index"
        :value="tab.value"
      >
        {{ tab.title }}
      </v-tab>
    </v-tabs>

    <v-card-text>
      <v-tabs-window v-model="activeTab">
        <v-tabs-window-item
          v-for="(tab, index) in tabs"
          :key="index"
          :value="tab.value"
        >
          {{ tab.content }}
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { ref } from 'vue'

const activeTab = ref('users')

const tabs = ref([
  { title: 'Users', value: 'users', content: 'User management content...' },
  { title: 'Roles', value: 'roles', content: 'Role management content...' },
  { title: 'Permissions', value: 'permissions', content: 'Permissions content...' }
])
</script>
```

## Using VWindow (Legacy)

In Vuetify 4, `v-tabs-window` replaces `v-window`:

```vue
<v-tabs v-model="tab" color="primary">
  <v-tab value="one">Tab 1</v-tab>
  <v-tab value="two">Tab 2</v-tab>
</v-tabs>

<v-tabs-window v-model="tab">
  <v-tabs-window-item value="one">
    Tab 1 content
  </v-tabs-window-item>
  <v-tabs-window-item value="two">
    Tab 2 content
  </v-tabs-window-item>
</v-tabs-window>
```

## Combined Examples

```vue
<template>
  <v-card>
    <v-tabs
      v-model="tab"
      color="primary"
      align="centered"
      bg-color="grey-lighten-4"
    >
      <v-tab value="overview">
        <v-icon start>mdi-view-dashboard</v-icon>
        Overview
      </v-tab>
      <v-tab value="analytics">
        <v-icon start>mdi-chart-line</v-icon>
        Analytics
      </v-tab>
      <v-tab value="reports">
        <v-icon start>mdi-file-document</v-icon>
        Reports
      </v-tab>
      <v-tab value="settings">
        <v-icon start>mdi-cog</v-icon>
        Settings
      </v-tab>
    </v-tabs>

    <v-divider></v-divider>

    <v-tabs-window v-model="tab">
      <v-tabs-window-item value="overview">
        <v-card-text>
          <h3>Overview Dashboard</h3>
          <p>Welcome to the overview section.</p>
        </v-card-text>
      </v-tabs-window-item>

      <v-tabs-window-item value="analytics">
        <v-card-text>
          <h3>Analytics</h3>
          <p>View your analytics data here.</p>
        </v-card-text>
      </v-tabs-window-item>

      <v-tabs-window-item value="reports">
        <v-card-text>
          <h3>Reports</h3>
          <p>Generate and view reports.</p>
        </v-card-text>
      </v-tabs-window-item>

      <v-tabs-window-item value="settings">
        <v-card-text>
          <h3>Settings</h3>
          <p>Configure your preferences.</p>
        </v-card-text>
      </v-tabs-window-item>
    </v-tabs-window>
  </v-card>
</template>

<script setup>
import { ref } from 'vue'
const tab = ref('overview')
</script>
```

## With Pagination

```vue
<template>
  <v-card>
    <v-tabs
      v-model="tab"
      color="primary"
      show-arrows
    >
      <v-tab v-for="i in 15" :key="i" :value="'tab-' + i">
        Tab {{ i }}
      </v-tab>
    </v-tabs>

    <v-tabs-window v-model="tab">
      <v-tabs-window-item
        v-for="i in 15"
        :key="i"
        :value="'tab-' + i"
      >
        <v-card-text>
          Content for Tab {{ i }}
        </v-card-text>
      </v-tabs-window-item>
    </v-tabs-window>
  </v-card>
</template>

<script setup>
import { ref } from 'vue'
const tab = ref('tab-1')
</script>
```

## Vertical Tabs

```vue
<template>
  <v-row>
    <v-col cols="3">
      <v-tabs v-model="tab" direction="vertical" color="primary">
        <v-tab value="one">Tab 1</v-tab>
        <v-tab value="two">Tab 2</v-tab>
        <v-tab value="three">Tab 3</v-tab>
      </v-tabs>
    </v-col>
    <v-col cols="9">
      <v-tabs-window v-model="tab">
        <v-tabs-window-item value="one">
          Content 1
        </v-tabs-window-item>
        <v-tabs-window-item value="two">
          Content 2
        </v-tabs-window-item>
        <v-tabs-window-item value="three">
          Content 3
        </v-tabs-window-item>
      </v-tabs-window>
    </v-col>
  </v-row>
</template>

<script setup>
import { ref } from 'vue'
const tab = ref('one')
</script>
```
