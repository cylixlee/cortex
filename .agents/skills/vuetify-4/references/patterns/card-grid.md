# Card Grid Layout

Responsive card grid layout for product listings.

```vue
<template>
  <v-app>
    <v-app-bar>
      <v-app-bar-title>Product List</v-app-bar-title>
    </v-app-bar>
    
    <v-main>
      <v-container>
        <v-row>
          <v-col
            v-for="product in products"
            :key="product.id"
            cols="12"
            sm="6"
            md="4"
            lg="3"
          >
            <v-card hover>
              <v-img
                :src="product.image"
                height="200"
                cover
              ></v-img>
              
              <v-card-item>
                <v-card-title>{{ product.name }}</v-card-title>
                <v-card-subtitle>{{ product.category }}</v-card-subtitle>
              </v-card-item>
              
              <v-card-text>
                <div class="text-h6 text-primary mb-2">
                  ${{ product.price }}
                </div>
                <div class="text-body-2 text-medium-emphasis">
                  {{ product.description }}
                </div>
              </v-card-text>
              
              <v-divider></v-divider>
              
              <v-card-actions>
                <v-btn variant="text" color="primary">
                  Details
                </v-btn>
                <v-spacer></v-spacer>
                <v-btn
                  icon="mdi-cart"
                  variant="text"
                  color="primary"
                ></v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref } from 'vue'

const products = ref([
  { id: 1, name: 'Product A', category: 'Electronics', price: 299, description: 'High quality electronics', image: 'https://picsum.photos/400/300?random=1' },
  { id: 2, name: 'Product B', category: 'Clothing', price: 199, description: 'Fashion clothing', image: 'https://picsum.photos/400/300?random=2' },
  { id: 3, name: 'Product C', category: 'Home', price: 399, description: 'Premium home goods', image: 'https://picsum.photos/400/300?random=3' },
  { id: 4, name: 'Product D', category: 'Sports', price: 299, description: 'Sports equipment', image: 'https://picsum.photos/400/300?random=4' },
  { id: 5, name: 'Product E', category: 'Books', price: 99, description: 'Bestselling books', image: 'https://picsum.photos/400/300?random=5' },
  { id: 6, name: 'Product F', category: 'Food', price: 59, description: 'Specialty foods', image: 'https://picsum.photos/400/300?random=6' }
])
</script>
```
