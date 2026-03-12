# Cards - VCard

`v-card` is used to display groups of content, providing title, subtitle, body, and action areas.

## Basic Structure

```vue
<v-card>
  <v-card-title>Title</v-card-title>
  <v-card-subtitle>Subtitle</v-card-subtitle>
  <v-card-text>
    Card body content
  </v-card-text>
  <v-card-actions>
    <v-btn>Action</v-btn>
  </v-card-actions>
</v-card>
```

## Complete Example

```vue
<v-card title="Card Title" subtitle="Subtitle">
  <v-card-text>
    This is the card body content, can contain any HTML.
    Vuetify provides good default styles.
  </v-card-text>
  <v-card-actions>
    <v-btn color="primary">Confirm</v-btn>
    <v-btn variant="text">Cancel</v-btn>
  </v-card-actions>
</v-card>
```

## Media Area

### Top Image

```vue
<v-card
  image="https://example.com/image.jpg"
  title="Image Title"
>
</v-card>
```

### Bottom Image

```vue
<v-card>
  <v-card-title>Title</v-card-title>
  <v-card-text>Body content</v-card-text>
  <v-img
    src="https://example.com/image.jpg"
    height="200"
  ></v-img>
</v-card>
```

### Using prepend/append Slots

```vue
<v-card>
  <template v-slot:prepend>
    <v-avatar color="primary">
      <v-icon icon="mdi-account"></v-icon>
    </v-avatar>
  </template>
  
  <v-card-title>User Info</v-card-title>
  <v-card-subtitle>Admin</v-card-subtitle>
  
  <v-card-text>
    User detail content
  </v-card-text>
</v-card>
```

## Variants & Styles

### Colors

```vue
<v-card color="primary" dark>
  <v-card-title>Primary Card</v-card-title>
  <v-card-text>Dark background</v-card-text>
</v-card>

<v-card color="surface-variant">
  <v-card-title>Surface Variant</v-card-title>
</v-card>
```

### Height & Elevation

```vue
<!-- Fixed height -->
<v-card height="200">
  Fixed Height Card
</v-card>

<!-- Elevation levels: 0-24 -->
<v-card elevation="4">High Elevation</v-card>
<v-card elevation="0">No Elevation</v-card>
```

### Rounded

```vue
<v-card rounded="lg">Large Rounded</v-card>
<v-card rounded="0">No Rounded</v-card>
```

## Interactive Cards

### Clickable

```vue
<v-card 
  hover 
  @click="handleClick"
>
  <v-card-title>Clickable Card</v-card-title>
  <v-card-text>Click event</v-card-text>
</v-card>
```

### With Link

```vue
<v-card 
  :to="/page"
  link
>
  <v-card-title>Navigation Card</v-card-title>
</v-card>
```

## Loading State

```vue
<v-card :loading="loading">
  <v-card-title>Title</v-card-title>
  <v-card-text>Content</v-card-text>
</v-card>
```

Custom loading indicator:

```vue
<v-card loading="saving">
  <v-card-title>Saving...</v-card-title>
</v-card>
```

## Card Grid Layout

```vue
<template>
  <v-container>
    <v-row>
      <v-col
        v-for="n in 3"
        :key="n"
        cols="12"
        sm="6"
        md="4"
      >
        <v-card>
          <v-img
            :src="`https://picsum.photos/500/300?random=${n}`"
            height="200"
            cover
          ></v-img>
          <v-card-title>Card {{ n }}</v-card-title>
          <v-card-text>
            This is the card content description
          </v-card-text>
          <v-card-actions>
            <v-btn variant="text">Details</v-btn>
            <v-btn color="primary">Action</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
```

## Complex Example: News Card

```vue
<v-card>
  <v-img
    src="https://example.com/news.jpg"
    height="200"
    cover
  ></v-img>
  
  <v-card-item>
    <v-card-title>News Title</v-card-title>
    <v-card-subtitle>
      <v-icon icon="mdi-clock" size="small" class="mr-1"></v-icon>
      2 hours ago
    </v-card-subtitle>
  </v-card-item>
  
  <v-card-text>
    This is the news summary content, showing how to use card components to present news articles...
  </v-card-text>
  
  <v-divider></v-divider>
  
  <v-card-actions>
    <v-btn variant="text" prepend-icon="mdi-thumb-up">
      Like
    </v-btn>
    <v-btn variant="text" prepend-icon="mdi-share">
      Share
    </v-btn>
    <v-spacer></v-spacer>
    <v-btn variant="text" color="primary">
      Read More
    </v-btn>
  </v-card-actions>
</v-card>
```
