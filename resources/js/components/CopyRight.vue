<template>
  <p class="text-center text-gray-500 text-xs">&copy;2020 {{name}}. All rights reserved.</p>
</template>

<script>
export default {
  name: 'CopyRight',
  data() {
    return {
      name: process.env.MIX_APP_NAME
    }
  }
}
</script>