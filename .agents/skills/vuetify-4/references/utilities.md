# Utility Classes

Vuetify provides rich utility classes for spacing, alignment, display control, and more.

## Spacing

### Margin (ma-)

| Class     | Value |
| --------- | ----- |
| `ma-0`    | 0     |
| `ma-1`    | 4px   |
| `ma-2`    | 8px   |
| `ma-3`    | 12px  |
| `ma-4`    | 16px  |
| `ma-5`    | 20px  |
| `ma-auto` | auto  |

```vue
<div class="ma-4">16px margin on all sides</div>
```

### Padding (pa-)

```vue
<div class="pa-4">16px padding on all sides</div>
```

### Horizontal/Vertical Spacing

- `mx-*` - Horizontal (left and right) spacing
- `my-*` - Vertical (top and bottom) spacing

```vue
<div class="mx-auto">Center horizontally (auto margin left/right)</div>
<div class="my-4">16px margin top/bottom</div>
```

### Single Side Spacing

- `mt-`, `mb-`, `ml-`, `mr-` - Top, bottom, left, right
- `ms-` (start), `me-` (end) - Start/End (supports RTL)

```vue
<div class="mt-4 mb-2">16px top, 8px bottom</div>
```

## Alignment

### Main Axis Alignment (justify-)

| Class                   | Effect                      |
| ----------------------- | --------------------------- |
| `justify-start`         | Left align                  |
| `justify-center`        | Center                      |
| `justify-end`           | Right align                 |
| `justify-space-between` | Justify between             |
| `justify-space-around`  | Space around                |
| `justify-space-evenly`  | Even space (including ends) |

### Cross Axis Alignment (align-)

| Class           | Effect       |
| --------------- | ------------ |
| `align-start`   | Top align    |
| `align-center`  | Center       |
| `align-end`     | Bottom align |
| `align-stretch` | Stretch      |

### Self Alignment

| Class                | Effect            |
| -------------------- | ----------------- |
| `align-self-start`   | Self top align    |
| `align-self-center`  | Self center       |
| `align-self-end`     | Self bottom align |
| `align-self-stretch` | Self stretch      |

```vue
<v-row>
  <v-col align-self="center">Vertically Centered</v-col>
</v-row>
```

## Display Control

### Display Property (d-)

```vue
<div class="d-none">Hidden</div>
<div class="d-block">Display as block</div>
<div class="d-inline">Inline display</div>
<div class="d-flex">Flex container</div>
<div class="d-inline-flex">Inline flex</div>
```

### Responsive Display

```vue
<!-- Hidden on phone, displayed on desktop -->
<div class="d-none d-md-flex">Desktop Display</div>

<!-- Displayed on phone, hidden on desktop -->
<div class="d-flex d-md-none">Phone Display</div>
```

### Responsive Breakpoint Suffixes

| Suffix | Min Width |
| ------ | --------- |
| `sm`   | 600px     |
| `md`   | 960px     |
| `lg`   | 1280px    |
| `xl`   | 1920px    |
| `xxl`  | 2560px    |

### Common Combination Examples

```vue
<!-- Centered Flex Container -->
<div class="d-flex justify-center align-center">Content</div>

<!-- Responsive Grid -->
<div class="d-flex flex-wrap">
  <div class="w-50 w-md-33 w-lg-25">Responsive Item</div>
</div>
```

## Sizing

### Width

| Class    | Value |
| -------- | ----- |
| `w-25`   | 25%   |
| `w-50`   | 50%   |
| `w-75`   | 75%   |
| `w-100`  | 100%  |
| `w-auto` | auto  |

### Height

| Class      | Value |
| ---------- | ----- |
| `h-25`     | 25%   |
| `h-50`     | 50%   |
| `h-75`     | 75%   |
| `h-100`    | 100%  |
| `h-auto`   | auto  |
| `h-screen` | 100vh |

```vue
<v-container class="fill-height">Full viewport height</v-container>
```

## Text

### Text Alignment

```vue
<p class="text-left">Left align</p>
<p class="text-center">Center</p>
<p class="text-right">Right align</p>
<p class="text-justify">Justify</p>
```

### Responsive Text Alignment

```vue
<p class="text-md-center">Center on desktop</p>
```

### Text Transform

```vue
<p class="text-lowercase">lowercase</p>
<p class="text-uppercase">uppercase</p>
<p class="text-capitalize">capitalize</p>
```

### Typography

```vue
<p class="text-h1">Heading 1</p>
<p class="text-h2">Heading 2</p>
<p class="text-body-1">Body</p>
<p class="text-caption">Caption</p>
```

## MD3 Typography Classes

Vuetify 4 introduces Material Design 3 specification typography classes:

### Display

| Class                 | Description           |
| --------------------- | --------------------- |
| `text-display-large`  | Display Large (57px)  |
| `text-display-medium` | Display Medium (45px) |
| `text-display-small`  | Display Small (36px)  |

### Headline

| Class                  | Description            |
| ---------------------- | ---------------------- |
| `text-headline-large`  | Headline Large (32px)  |
| `text-headline-medium` | Headline Medium (28px) |
| `text-headline-small`  | Headline Small (24px)  |

### Title

| Class               | Description                     |
| ------------------- | ------------------------------- |
| `text-title-large`  | Title Large (22px)              |
| `text-title-medium` | Title Medium (16px, 500 weight) |
| `text-title-small`  | Title Small (14px, 500 weight)  |

### Body

| Class              | Description        |
| ------------------ | ------------------ |
| `text-body-large`  | Body Large (16px)  |
| `text-body-medium` | Body Medium (14px) |
| `text-body-small`  | Body Small (12px)  |

### Label

| Class               | Description         |
| ------------------- | ------------------- |
| `text-label-large`  | Label Large (14px)  |
| `text-label-medium` | Label Medium (12px) |
| `text-label-small`  | Label Small (11px)  |

```vue
<p class="text-display-large">Display Large</p>
<p class="text-headline-medium">Headline Medium</p>
<p class="text-title-small">Title Small</p>
<p class="text-body-large">Body Large</p>
<p class="text-label-medium">Label Medium</p>
```

## Text Emphasis

### High Emphasis

```vue
<span class="text-high-emphasis">High emphasis text</span>
```

### Medium Emphasis

```vue
<span class="text-medium-emphasis">Medium emphasis text</span>
```

### Disabled

```vue
<span class="text-disabled">Disabled text</span>
```

## Text Decoration

```vue
<p class="text-decoration-underline">Underline text</p>
<p class="text-decoration-line-through">Line through text</p>
<p class="text-decoration-overline">Overline text</p>
<p class="text-decoration-none">No decoration (reset)</p>
```
