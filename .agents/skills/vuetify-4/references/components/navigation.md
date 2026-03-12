# Navigation Components

## VAppBar - Top Navigation Bar

### Basic Usage

```vue
<v-app>
  <v-app-bar title="App Title"></v-app-bar>
  <v-main>Content</v-main>
</v-app>
```

### Common Props

| Prop        | Type   | Description                                |
| ----------- | ------ | ------------------------------------------ |
| `title`     | string | Title text                                 |
| `elevation` | number | Elevation level                            |
| `color`     | string | Background color                           |
| `density`   | string | Density: 'default'/'comfortable'/'compact' |

### Position

```vue
<!-- Top (default) -->
<v-app-bar title="Top Navigation"></v-app-bar>

<!-- Bottom -->
<v-app-bar location="bottom" title="Bottom Navigation"></v-app-bar>
```

### Scroll Behavior

```vue
<v-app-bar
  title="Scroll Hide"
  scroll-threshold="100"
  scroll-behavior="hide"
>
</v-app-bar>
```

| scroll-behavior value | Description          |
| --------------------- | -------------------- |
| `hide`                | Hide on scroll up    |
| `show`                | Show on scroll down  |
| `inverted`            | Inverted behavior    |
| `fade-image`          | Fade image on scroll |

### Complete Example

```vue
<v-app-bar color="primary" elevation="2">
  <v-app-bar-nav-icon></v-app-bar-nav-icon>
  <v-app-bar-title>My App</v-app-bar-title>
  <v-spacer></v-spacer>
  <v-btn icon="mdi-magnify"></v-btn>
  <v-btn icon="mdi-bell"></v-btn>
  <v-btn icon="mdi-account"></v-btn>
</v-app-bar>
```

## VNavigationDrawer - Side Navigation

### Basic Usage

```vue
<v-app>
  <v-navigation-drawer v-model="drawer">
    <v-list>
      <v-list-item title="Nav Item 1"></v-list-item>
      <v-list-item title="Nav Item 2"></v-list-item>
    </v-list>
  </v-navigation-drawer>
  
  <v-app-bar>
    <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
    <v-app-bar-title>Title</v-app-bar-title>
  </v-app-bar>
  
  <v-main>Content</v-main>
</v-app>
</template>

<script setup>
import { ref } from 'vue'
const drawer = ref(true)
</script>
```

### Rail Mode (Collapsed State)

```vue
<v-navigation-drawer
  v-model="drawer"
  rail
>
  <v-list>
    <v-list-item icon="mdi-home"></v-list-item>
    <v-list-item icon="mdi-account"></v-list-item>
  </v-list>
</v-navigation-drawer>
```

### Toggle Rail Mode

```vue
<template>
  <v-navigation-drawer
    v-model="drawer"
    :rail="rail"
  >
    <v-list>
      <v-list-item
        :prepend-icon="rail ? 'mdi-menu' : 'mdi-backburger'"
        @click="rail = !rail"
      ></v-list-item>
      <v-divider></v-divider>
      <v-list-item prepend-icon="mdi-home" title="Home"></v-list-item>
      <v-list-item prepend-icon="mdi-account" title="Users"></v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup>
import { ref } from 'vue'
const drawer = ref(true)
const rail = ref(false)
</script>
```

### Position

```vue
<!-- Left (default) -->
<v-navigation-drawer location="left"></v-navigation-drawer>

<!-- Right -->
<v-navigation-drawer location="right"></v-navigation-drawer>
```

### Width

```vue
<v-navigation-drawer width="300"></v-navigation-drawer>
<v-navigation-drawer rail-width="72"></v-navigation-drawer>
```

### Temporary/Permanent Mode

```vue
<!-- Temporary drawer (default on mobile) -->
<v-navigation-drawer temporary></v-navigation-drawer>

<!-- Permanent -->
<v-navigation-drawer permanent></v-navigation-drawer>
```

## VBottomNavigation - Bottom Navigation

### Basic Usage

```vue
<template>
  <v-bottom-navigation v-model="value">
    <v-btn value="home">
      <v-icon>mdi-home</v-icon>
      <span>Home</span>
    </v-btn>
    <v-btn value="search">
      <v-icon>mdi-magnify</v-icon>
      <span>Search</span>
    </v-btn>
    <v-btn value="account">
      <v-icon>mdi-account</v-icon>
      <span>Profile</span>
    </v-btn>
  </v-bottom-navigation>
</template>

<script setup>
import { ref } from 'vue'
const value = ref('home')
</script>
```

### Icons Only Mode

```vue
<v-bottom-navigation v-model="value" icon>
  <v-btn value="home">
    <v-icon>mdi-home</v-icon>
  </v-btn>
  <v-btn value="search">
    <v-icon>mdi-magnify</v-icon>
  </v-btn>
  <v-btn value="account">
    <v-icon>mdi-account</v-icon>
  </v-btn>
</v-bottom-navigation>
```

### Colors

```vue
<v-bottom-navigation v-model="value" color="primary">
  <!-- Navigation items -->
</v-bottom-navigation>
```

### Scroll Hide

```vue
<v-bottom-navigation
  v-model="value"
  grow
  scroll-threshold="100"
  scroll-behavior="hide"
>
  <!-- Navigation items -->
</v-bottom-navigation>
```

### Using v-btn-group

```vue
<v-bottom-navigation v-model="value" grow>
  <v-btn value="home">
    <v-icon>mdi-home</v-icon>
  </v-btn>
  <v-btn value="search">
    <v-icon>mdi-magnify</v-icon>
  </v-btn>
  <v-btn value="settings">
    <v-icon>mdi-cog</v-icon>
  </v-btn>
</v-bottom-navigation>
```

## Complete Example: App Layout

```vue
<template>
  <v-app>
    <!-- Top Navigation -->
    <v-app-bar color="primary">
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-app-bar-title>My App</v-app-bar-title>
      <v-spacer></v-spacer>
      <v-btn icon="mdi-magnify"></v-btn>
      <v-btn icon="mdi-bell"></v-btn>
    </v-app-bar>
    
    <!-- Side Navigation -->
    <v-navigation-drawer v-model="drawer">
      <v-list>
        <v-list-item
          :prepend-avatar="user.avatar"
          :title="user.name"
          :subtitle="user.email"
        ></v-list-item>
      </v-list>
      <v-divider></v-divider>
      <v-list density="compact" nav>
        <v-list-item
          v-for="item in navItems"
          :key="item.title"
          :prepend-icon="item.icon"
          :title="item.title"
          :value="item.value"
        ></v-list-item>
      </v-list>
    </v-navigation-drawer>
    
    <!-- Main Content -->
    <v-main>
      <v-container>
        <p>Page content...</p>
      </v-container>
    </v-main>
    
    <!-- Bottom Navigation -->
    <v-bottom-navigation v-model="bottomNav" grow>
      <v-btn value="home" icon="mdi-home"></v-btn>
      <v-btn value="search" icon="mdi-magnify"></v-btn>
      <v-btn value="settings" icon="mdi-cog"></v-btn>
    </v-bottom-navigation>
  </v-app>
</template>

<script setup>
import { ref } from 'vue'

const drawer = ref(true)
const bottomNav = ref('home')

const user = {
  name: 'John Doe',
  email: 'john@example.com',
  avatar: 'https://example.com/avatar.jpg'
}

const navItems = [
  { title: 'Home', icon: 'mdi-home', value: 'home' },
  { title: 'About', icon: 'mdi-information', value: 'about' },
  { title: 'Settings', icon: 'mdi-cog', value: 'settings' }
]
</script>
```
