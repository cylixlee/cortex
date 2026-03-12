# Data Table - VDataTable

`v-data-table` is used to display structured data.

## Basic Usage

```vue
<template>
  <v-data-table
    :headers="headers"
    :items="items"
  ></v-data-table>
</template>

<script setup>
import { ref } from 'vue'

const headers = [
  { title: 'Name', key: 'name' },
  { title: 'Age', key: 'age' },
  { title: 'City', key: 'city' }
]

const items = ref([
  { name: 'John', age: 25, city: 'New York' },
  { name: 'Jane', age: 30, city: 'Los Angeles' },
  { name: 'Bob', age: 28, city: 'Chicago' }
])
</script>
```

## Common Props

| Prop             | Type    | Description                                   |
| ---------------- | ------- | --------------------------------------------- |
| `headers`        | array   | Header configuration                          |
| `items`          | array   | Data items                                    |
| `items-length`   | number  | Total items (used for server-side pagination) |
| `items-per-page` | number  | Items per page                                |
| `page`           | number  | Current page                                  |
| `loading`        | boolean | Loading state                                 |
| `hover`          | boolean | Hover effect                                  |

## Header Configuration

```vue
const headers = [
  { title: 'Name', key: 'name', align: 'start' },
  { title: 'Age', key: 'age', align: 'center' },
  { title: 'City', key: 'city', align: 'end' },
  { title: 'Actions', key: 'actions', sortable: false }
]
```

| Config     | Description                       |
| ---------- | --------------------------------- |
| `title`    | Display title                     |
| `key`      | Corresponding data field          |
| `align`    | Alignment: 'start'/'center'/'end' |
| `sortable` | Whether sortable                  |
| `width`    | Column width                      |

## Pagination

### Client-side Pagination

```vue
<v-data-table
  :headers="headers"
  :items="items"
  :items-per-page="5"
></v-data-table>
```

### Server-side Pagination

```vue
<template>
  <v-data-table
    :headers="headers"
    :items="items"
    :items-length="totalItems"
    :loading="loading"
    @update:options="loadItems"
  ></v-data-table>
</template>

<script setup>
import { ref } from 'vue'

const items = ref([])
const totalItems = ref(0)
const loading = ref(false)

const headers = [
  { title: 'Name', key: 'name' },
  { title: 'Age', key: 'age' }
]

const loadItems = async ({ page, itemsPerPage }) => {
  loading.value = true
  // Simulate API request
  const result = await fetchData(page, itemsPerPage)
  items.value = result.items
  totalItems.value = result.total
  loading.value = false
}
</script>
```

## Sorting

```vue
<v-data-table
  :headers="headers"
  :items="items"
  :sort-by="[{ key: 'age', order: 'asc' }]"
></v-data-table>
```

## Row Actions

```vue
<template>
  <v-data-table
    :headers="headers"
    :items="items"
  >
    <template v-slot:item.actions="{ item }">
      <v-btn
        icon="mdi-pencil"
        variant="text"
        size="small"
        @click="edit(item)"
      ></v-btn>
      <v-btn
        icon="mdi-delete"
        variant="text"
        size="small"
        color="error"
        @click="remove(item)"
      ></v-btn>
    </template>
  </v-data-table>
</template>

<script setup>
const headers = [
  { title: 'Name', key: 'name' },
  { title: 'Actions', key: 'actions', sortable: false }
]

const items = ref([
  { name: 'John' },
  { name: 'Jane' }
])

const edit = (item) => console.log('Edit', item)
const remove = (item) => console.log('Delete', item)
</script>
```

## Search

```vue
<template>
  <v-card>
    <v-card-title>
      User List
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        label="Search"
        prepend-inner-icon="mdi-magnify"
        single-line
        hide-details
      ></v-text-field>
    </v-card-title>
    
    <v-data-table
      :headers="headers"
      :items="items"
      :search="search"
    ></v-data-table>
  </v-card>
</template>

<script setup>
import { ref } from 'vue'

const search = ref('')
const headers = [
  { title: 'Name', key: 'name' },
  { title: 'Age', key: 'age' }
]
const items = ref([
  { name: 'John', age: 25 },
  { name: 'Jane', age: 30 }
])
</script>
```

## Row Selection

```vue
<template>
  <v-data-table
    v-model="selected"
    :headers="headers"
    :items="items"
    show-select
  ></v-data-table>
  <p>Selected: {{ selected }}</p>
</template>

<script setup>
import { ref } from 'vue'

const selected = ref([])
const headers = [{ title: 'Name', key: 'name' }]
const items = ref([
  { name: 'John' },
  { name: 'Jane' }
])
</script>
```

## Loading State

```vue
<v-data-table
  :headers="headers"
  :items="items"
  loading
></v-data-table>

<v-data-table
  :headers="headers"
  :items="items"
  loading-text="Loading..."
></v-data-table>
```

## Hide Pagination

```vue
<v-data-table
  :headers="headers"
  :items="items"
  :hide-default-footer="true"
></v-data-table>
```

## Complete Example: Data Table with CRUD

```vue
<template>
  <v-card>
    <v-card-title class="d-flex align-center">
      User Management
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        label="Search"
        prepend-inner-icon="mdi-magnify"
        single-line
        hide-details
        density="compact"
        class="mr-4"
        style="max-width: 300px"
      ></v-text-field>
      <v-btn color="primary" prepend-icon="mdi-plus">
        Add New
      </v-btn>
    </v-card-title>
    
    <v-data-table
      :headers="headers"
      :items="items"
      :search="search"
      :loading="loading"
      hover
    >
      <template v-slot:item.status="{ item }">
        <v-chip
          :color="item.status === 'active' ? 'success' : 'grey'"
          size="small"
        >
          {{ item.status }}
        </v-chip>
      </template>
      
      <template v-slot:item.actions="{ item }">
        <v-btn
          icon="mdi-pencil"
          variant="text"
          size="small"
          @click="editItem(item)"
        ></v-btn>
        <v-btn
          icon="mdi-delete"
          variant="text"
          size="small"
          color="error"
          @click="deleteItem(item)"
        ></v-btn>
      </template>
    </v-data-table>
  </v-card>
</template>

<script setup>
import { ref } from 'vue'

const search = ref('')
const loading = ref(false)

const headers = [
  { title: 'ID', key: 'id', width: '80px' },
  { title: 'Name', key: 'name' },
  { title: 'Email', key: 'email' },
  { title: 'Status', key: 'status' },
  { title: 'Actions', key: 'actions', sortable: false, align: 'center' }
]

const items = ref([
  { id: 1, name: 'John Doe', email: 'john@example.com', status: 'active' },
  { id: 2, name: 'Jane Smith', email: 'jane@example.com', status: 'inactive' },
  { id: 3, name: 'Bob Wilson', email: 'bob@example.com', status: 'active' }
])

const editItem = (item) => console.log('Edit', item)
const deleteItem = (item) => console.log('Delete', item)
</script>
```
