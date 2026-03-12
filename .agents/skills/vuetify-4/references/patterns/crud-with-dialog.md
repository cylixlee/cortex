# CRUD with Dialog

Complete example of a list page with add/edit dialog and delete confirmation.

```vue
<template>
  <v-app>
    <v-app-bar elevation="1">
      <v-app-bar-title>User Management</v-app-bar-title>
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        label="Search"
        prepend-inner-icon="mdi-magnify"
        hide-details
        single-line
        density="compact"
        class="mr-4"
        style="max-width: 300px"
      ></v-text-field>
      <v-btn color="primary" prepend-icon="mdi-plus" @click="openAddDialog">
        Add User
      </v-btn>
    </v-app-bar>
    
    <v-main>
      <v-container fluid>
        <v-card>
          <v-data-table
            :headers="headers"
            :items="items"
            :search="search"
            :loading="loading"
            hover
          >
            <template v-slot:item.name="{ item }">
              <div class="d-flex align-center py-2">
                <v-avatar color="primary" size="32" class="mr-2">
                  <span class="text-white">{{ item.name[0] }}</span>
                </v-avatar>
                {{ item.name }}
              </div>
            </template>
            
            <template v-slot:item.status="{ item }">
              <v-chip
                :color="item.status === 'active' ? 'success' : 'grey'"
                size="small"
              >
                {{ item.status === 'active' ? 'Active' : 'Inactive' }}
              </v-chip>
            </template>
            
            <template v-slot:item.actions="{ item }">
              <v-btn
                icon="mdi-pencil"
                variant="text"
                size="small"
                @click="openEditDialog(item)"
              ></v-btn>
              <v-btn
                icon="mdi-delete"
                variant="text"
                size="small"
                color="error"
                @click="openDeleteDialog(item)"
              ></v-btn>
            </template>
          </v-data-table>
        </v-card>
      </v-container>
    </v-main>

    <!-- Add/Edit Dialog -->
    <v-dialog v-model="userDialog" max-width="600" persistent>
      <v-card>
        <v-card-title>
          {{ isEditing ? 'Edit User' : 'Add User' }}
        </v-card-title>
        
        <v-card-text>
          <v-form ref="form" v-model="valid">
            <v-text-field
              v-model="userForm.name"
              label="Name"
              :rules="[rules.required]"
              variant="outlined"
              class="mb-2"
            ></v-text-field>
            
            <v-text-field
              v-model="userForm.email"
              label="Email"
              type="email"
              :rules="[rules.required, rules.email]"
              variant="outlined"
              class="mb-2"
            ></v-text-field>
            
            <v-select
              v-model="userForm.role"
              label="Role"
              :items="['Admin', 'User', 'Guest']"
              :rules="[rules.required]"
              variant="outlined"
              class="mb-2"
            ></v-select>
            
            <v-select
              v-model="userForm.status"
              label="Status"
              :items="['active', 'inactive']"
              :rules="[rules.required]"
              variant="outlined"
            ></v-select>
          </v-form>
        </v-card-text>
        
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="closeUserDialog">Cancel</v-btn>
          <v-btn
            color="primary"
            :disabled="!valid"
            :loading="saving"
            @click="saveUser"
          >
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title>Confirm Delete</v-card-title>
        <v-card-text>
          Are you sure you want to delete <strong>{{ userToDelete?.name }}</strong>?
          This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn
            color="error"
            :loading="deleting"
            @click="confirmDelete"
          >
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Success Snackbar -->
    <v-snackbar v-model="snackbar" :color="snackbarColor" timeout="3000">
      {{ snackbarText }}
      <template v-slot:actions>
        <v-btn variant="text" @click="snackbar = false">Close</v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<script setup>
import { ref, reactive } from 'vue'

const search = ref('')
const loading = ref(false)
const valid = ref(false)
const saving = ref(false)
const deleting = ref(false)
const userDialog = ref(false)
const deleteDialog = ref(false)
const isEditing = ref(false)
const snackbar = ref(false)
const snackbarText = ref('')
const snackbarColor = ref('success')
const userToDelete = ref(null)

const userForm = reactive({
  id: null,
  name: '',
  email: '',
  role: 'User',
  status: 'active'
})

const emptyForm = {
  id: null,
  name: '',
  email: '',
  role: 'User',
  status: 'active'
}

const rules = {
  required: v => !!v || 'This field is required',
  email: v => /.+@.+\..+/.test(v) || 'Invalid email address'
}

const headers = [
  { title: 'Name', key: 'name' },
  { title: 'Email', key: 'email' },
  { title: 'Role', key: 'role' },
  { title: 'Status', key: 'status' },
  { title: 'Actions', key: 'actions', sortable: false, align: 'end' }
]

const items = ref([
  { id: 1, name: 'John Doe', email: 'john@example.com', role: 'Admin', status: 'active' },
  { id: 2, name: 'Jane Smith', email: 'jane@example.com', role: 'User', status: 'active' },
  { id: 3, name: 'Bob Wilson', email: 'bob@example.com', role: 'User', status: 'inactive' }
])

const showSnackbar = (text, color = 'success') => {
  snackbarText.value = text
  snackbarColor.value = color
  snackbar.value = true
}

const openAddDialog = () => {
  isEditing.value = false
  Object.assign(userForm, emptyForm)
  userDialog.value = true
}

const openEditDialog = (item) => {
  isEditing.value = true
  Object.assign(userForm, item)
  userDialog.value = true
}

const closeUserDialog = () => {
  userDialog.value = false
  Object.assign(userForm, emptyForm)
}

const openDeleteDialog = (item) => {
  userToDelete.value = item
  deleteDialog.value = true
}

const saveUser = async () => {
  saving.value = true
  
  try {
    await new Promise(resolve => setTimeout(resolve, 500))
    
    if (isEditing.value) {
      const index = items.value.findIndex(i => i.id === userForm.id)
      if (index !== -1) {
        items.value[index] = { ...userForm }
      }
      showSnackbar('User updated successfully')
    } else {
      const newUser = {
        ...userForm,
        id: Math.max(...items.value.map(i => i.id)) + 1
      }
      items.value.push(newUser)
      showSnackbar('User added successfully')
    }
    
    closeUserDialog()
  } catch (error) {
    showSnackbar('Failed to save user', 'error')
  } finally {
    saving.value = false
  }
}

const confirmDelete = async () => {
  deleting.value = true
  
  try {
    await new Promise(resolve => setTimeout(resolve, 500))
    
    items.value = items.value.filter(i => i.id !== userToDelete.value.id)
    showSnackbar('User deleted successfully')
    deleteDialog.value = false
    userToDelete.value = null
  } catch (error) {
    showSnackbar('Failed to delete user', 'error')
  } finally {
    deleting.value = false
  }
}
</script>
```

## Key Features

- **Search**: Filter table data in real-time
- **Add Dialog**: Form with validation for creating new users
- **Edit Dialog**: Pre-populated form for editing existing users
- **Delete Confirmation**: Confirmation dialog before deletion
- **Form Validation**: Required fields and email validation
- **Loading States**: Visual feedback during save/delete operations
- **Snackbar Notifications**: Success/error messages

## Form Validation Rules

```javascript
const rules = {
  required: v => !!v || 'This field is required',
  email: v => /.+@.+\..+/.test(v) || 'Invalid email address',
  minLength: v => (v && v.length >= 6) || 'Minimum 6 characters'
}
```

## Dialog Best Practices

- Use `persistent` prop for add/edit dialogs to prevent accidental close
- Use `max-width` to control dialog size
- Always show loading state during async operations
- Clear form data when dialog closes
- Use snackackbar for operation feedback
