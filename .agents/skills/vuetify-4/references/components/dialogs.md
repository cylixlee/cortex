# Dialog - VDialog

`v-dialog` is a modal dialog component used to display temporary content or user interactions.

## Basic Usage

```vue
<template>
  <v-btn @click="dialog = true">Open Dialog</v-btn>

  <v-dialog v-model="dialog" width="500">
    <v-card title="Dialog Title">
      <v-card-text>
        This is the dialog content.
      </v-card-text>
      <v-card-actions>
        <v-btn @click="dialog = false">Close</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref } from 'vue'
const dialog = ref(false)
</script>
```

## Common Props

| Prop         | Description                    | Default |
| ------------ | ------------------------------ | ------- |
| `v-model`    | Controls dialog visibility     | false   |
| `width`      | Dialog width                   | -       |
| `max-width`  | Maximum width                  | -       |
| `title`      | Dialog title                   | -       |
| `persistent` | Prevent close on outside click | false   |
| `fullscreen` | Fullscreen display             | false   |
| `scrollable` | Scrollable content             | false   |

## Complete Example

```vue
<template>
  <v-btn color="primary" @click="dialog = true">
    Open Dialog
  </v-btn>

  <v-dialog v-model="dialog" max-width="500">
    <v-card>
      <v-card-title>
        <span class="text-h5">Edit Profile</span>
      </v-card-title>
      
      <v-card-text>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-text-field
                label="Email"
                variant="outlined"
              ></v-text-field>
            </v-col>
            <v-col cols="12">
              <v-text-field
                label="Password"
                type="password"
                variant="outlined"
              ></v-text-field>
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>
      
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn @click="dialog = false">Cancel</v-btn>
        <v-btn color="primary" @click="dialog = false">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref } from 'vue'
const dialog = ref(false)
</script>
```

## Fullscreen Dialog

```vue
<v-dialog v-model="dialog" fullscreen>
  <v-card>
    <v-toolbar color="primary">
      <v-btn icon @click="dialog = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
      <v-toolbar-title>Settings</v-toolbar-title>
    </v-toolbar>
    
    <v-card-text>
      Fullscreen dialog content
    </v-card-text>
  </v-card>
</v-dialog>
```

## Persistent Dialog

When `persistent` is set, clicking the overlay won't close the dialog:

```vue
<v-dialog v-model="dialog" persistent width="500">
  <v-card>
    <v-card-title>Persistent Dialog</v-card-title>
    <v-card-text>
      Click outside to try closing - it won't work!
    </v-card-text>
    <v-card-actions>
      <v-btn @click="dialog = false">Close</v-btn>
    </v-card-actions>
  </v-card>
</v-dialog>
```

## Form Dialog

```vue
<template>
  <v-dialog v-model="dialog" max-width="500" persistent>
    <v-form @submit.prevent="submit">
      <v-card>
        <v-card-title>Login</v-card-title>
        
        <v-card-text>
          <v-text-field
            v-model="form.email"
            label="Email"
            :rules="[rules.required, rules.email]"
            variant="outlined"
          ></v-text-field>
          
          <v-text-field
            v-model="form.password"
            label="Password"
            type="password"
            variant="outlined"
          ></v-text-field>
        </v-card-text>
        
        <v-card-actions>
          <v-btn @click="dialog = false">Cancel</v-btn>
          <v-btn type="submit" color="primary">Login</v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>

<script setup>
import { ref, reactive } from 'vue'

const dialog = ref(false)
const form = reactive({
  email: '',
  password: ''
})

const rules = {
  required: v => !!v || 'Required',
  email: v => /.+@.+\..+/.test(v) || 'Invalid email'
}

const submit = () => {
  console.log(form)
  dialog.value = false
}
</script>
```

## Scroll Behavior

### Scrollable Dialog

```vue
<v-dialog v-model="dialog" scrollable width="500">
  <v-card title="Scrollable Dialog">
    <v-divider></v-divider>
    <v-card-text style="height: 300px;">
      <!-- Long content here -->
    </v-card-text>
  </v-card>
</v-dialog>
```

### Internal Scrolling

```vue
<v-dialog v-model="dialog" width="500">
  <v-card>
    <v-card-title>Title</v-card-title>
    <v-divider></v-divider>
    <v-card-text class="overflow-y-auto" style="max-height: 400px;">
      <!-- Scrollable content -->
    </v-card-text>
    <v-card-actions>
      <v-btn @click="dialog = false">Close</v-btn>
    </v-card-actions>
  </v-card>
</v-dialog>
```

## Dialog Sizes

```vue
<!-- Default -->
<v-dialog v-model="dialog">...</v-dialog>

<!-- Small -->
<v-dialog v-model="dialog" width="300">...</v-dialog>

<!-- Large -->
<v-dialog v-model="dialog" width="800">...</v-dialog>

<!-- No max width -->
<v-dialog v-model="dialog" width="auto">...</v-dialog>
```

## Transition

```vue
<v-dialog v-model="dialog" transition="dialog-bottom-transition">
  ...
</v-dialog>

<v-dialog v-model="dialog" transition="scale-transition">
  ...
</v-dialog>
```

## Using v-card Props

```vue
<v-dialog v-model="dialog" width="500">
  <v-card title="Dialog Title" subtitle="Optional subtitle">
    <v-card-text>
      Your content here
    </v-card-text>
    <v-card-actions>
      <v-btn @click="dialog = false">Close</v-btn>
    </v-card-actions>
  </v-card>
</v-dialog>
```
