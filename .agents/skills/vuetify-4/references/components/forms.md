# Form Components

## VTextField - Text Input

### Basic Usage

```vue
<v-text-field label="Username"></v-text-field>
<v-text-field label="Email" type="email"></v-text-field>
```

### Placeholder

```vue
<v-text-field
  label="Search"
  placeholder="Enter keywords"
></v-text-field>
```

### Variants

```vue
<v-text-field label="Default Variant"></v-text-field>
<v-text-field label="Outlined" variant="outlined"></v-text-field>
<v-text-field label="Underlined" variant="underlined"></v-text-field>
<v-text-field label="Filled" variant="filled"></v-text-field>
<v-text-field label="Plain" variant="plain"></v-text-field>
```

### Prepend/Append Icons

```vue
<v-text-field
  label="Search"
  prepend-inner-icon="mdi-magnify"
></v-text-field>

<v-text-field
  label="Password"
  :type="show ? 'text' : 'password'"
  :append-inner-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
  @click:append-inner="show = !show"
></v-text-field>
```

### v-model Binding

```vue
<template>
  <v-text-field
    v-model="username"
    label="Username"
  ></v-text-field>
  <p>Value: {{ username }}</p>
</template>

<script setup>
import { ref } from 'vue'
const username = ref('')
</script>
```

### Password Input

```vue
<template>
  <v-text-field
    v-model="password"
    label="Password"
    :type="showPassword ? 'text' : 'password'"
    :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
    @click:append-inner="showPassword = !showPassword"
  ></v-text-field>
</template>

<script setup>
import { ref } from 'vue'
const password = ref('')
const showPassword = ref(false)
</script>
```

### Validation Rules

```vue
<template>
  <v-text-field
    v-model="email"
    label="Email"
    :rules="[rules.required, rules.email]"
  ></v-text-field>
</template>

<script setup>
import { ref } from 'vue'
const email = ref('')

const rules = {
  required: v => !!v || 'Required',
  email: v => /.+@.+\..+/.test(v) || 'Invalid email'
}
</script>
```

## VSelect - Dropdown Select

### Basic Usage

```vue
<template>
  <v-select
    v-model="selected"
    label="Select"
    :items="['Option 1', 'Option 2', 'Option 3']"
  ></v-select>
</template>

<script setup>
import { ref } from 'vue'
const selected = ref(null)
</script>
```

### Object Array Options

```vue
<template>
  <v-select
    v-model="selected"
    label="Select Color"
    :items="colors"
    item-title="name"
    item-value="value"
  ></v-select>
</template>

<script setup>
import { ref } from 'vue'
const selected = ref(null)
const colors = [
  { name: 'Red', value: 'red' },
  { name: 'Green', value: 'green' },
  { name: 'Blue', value: 'blue' }
]
</script>
```

### Multiple Selection

```vue
<v-select
  v-model="selected"
  label="Multiple"
  :items="['A', 'B', 'C']"
  multiple
></v-select>
```

### Clearable

```vue
<v-select
  v-model="selected"
  label="Clearable"
  :items="['A', 'B', 'C']"
  clearable
></v-select>
```

## VSwitch - Switch

### Basic Usage

```vue
<template>
  <v-switch
    v-model="enabled"
    label="Enable"
  ></v-switch>
  <p>Status: {{ enabled }}</p>
</template>

<script setup>
import { ref } from 'vue'
const enabled = ref(false)
</script>
```

### Colors

```vue
<v-switch color="primary" label="Primary"></v-switch>
<v-switch color="error" label="Error"></v-switch>
```

### Disabled State

```vue
<v-switch disabled label="Disabled"></v-switch>
```

## VCheckbox - Checkbox

### Basic Usage

```vue
<template>
  <v-checkbox
    v-model="agree"
    label="I agree to the terms"
  ></v-checkbox>
</template>

<script setup>
import { ref } from 'vue'
const agree = ref(false)
</script>
```

### Multiple Selection

```vue
<template>
  <v-checkbox
    v-for="item in items"
    :key="item"
    v-model="selected"
    :label="item"
    :value="item"
  ></v-checkbox>
  <p>Selected: {{ selected }}</p>
</template>

<script setup>
import { ref } from 'vue'
const items = ['A', 'B', 'C']
const selected = ref([])
</script>
```

### Indeterminate State

```vue
<v-checkbox
  v-model="value"
  label="Indeterminate State"
  indeterminate
></v-checkbox>
```

## VForm - Form Container

### Basic Usage

```vue
<template>
  <v-form @submit.prevent="submit">
    <v-text-field
      v-model="form.name"
      label="Name"
      :rules="[v => !!v || 'Required']"
    ></v-text-field>
    
    <v-text-field
      v-model="form.email"
      label="Email"
      :rules="[rules.required, rules.email]"
    ></v-text-field>
    
    <v-btn type="submit" color="primary">Submit</v-btn>
  </v-form>
</template>

<script setup>
import { reactive, ref } from 'vue'

const form = reactive({
  name: '',
  email: ''
})

const rules = {
  required: v => !!v || 'Required',
  email: v => /.+@.+\..+/.test(v) || 'Invalid email'
}

const formRef = ref(null)

const submit = async () => {
  const { valid } = await formRef.value.validate()
  if (valid) {
    console.log('Form submitted')
  }
}
</script>
```

### Form Validation Timing

```vue
<v-form validate-on="blur">
  <v-text-field
    v-model="value"
    :rules="[v => !!v || 'Required']"
  ></v-text-field>
</v-form>
```

### Validation Trigger Options

| Value              | Description             |
| ------------------ | ----------------------- |
| `lazy`             | Validate on submit only |
| `validate-on-blur` | Validate on blur        |
| `input`            | Validate on input       |
| `submit`           | Validate on submit      |

### Reset Form

```vue
<template>
  <v-form ref="formRef">
    <v-text-field v-model="name"></v-text-field>
    <v-btn @click="reset">Reset</v-btn>
  </v-form>
</template>

<script setup>
const formRef = ref(null)

const reset = () => {
  formRef.value.reset()
}
</script>
```

## Form Example: Login Page

```vue
<template>
  <v-card width="400" class="mx-auto">
    <v-card-title>Login</v-card-title>
    <v-card-text>
      <v-form ref="formRef" v-model="valid">
        <v-text-field
          v-model="form.username"
          label="Username"
          prepend-inner-icon="mdi-account"
          :rules="[rules.required]"
        ></v-text-field>
        
        <v-text-field
          v-model="form.password"
          label="Password"
          :type="showPassword ? 'text' : 'password'"
          prepend-inner-icon="mdi-lock"
          :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
          :rules="[rules.required]"
          @click:append-inner="showPassword = !showPassword"
        ></v-text-field>
        
        <v-checkbox
          v-model="form.remember"
          label="Remember me"
        ></v-checkbox>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-btn
        color="primary"
        block
        :disabled="!valid"
        @click="login"
      >
        Login
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup>
import { reactive, ref } from 'vue'

const formRef = ref(null)
const valid = ref(false)
const showPassword = ref(false)

const form = reactive({
  username: '',
  password: '',
  remember: false
})

const rules = {
  required: v => !!v || 'Required'
}

const login = () => {
  console.log('Login:', form)
}
</script>
```
