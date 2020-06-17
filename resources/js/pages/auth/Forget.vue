<template>
  <div class="max-w-md mx-auto mb:pt-20 sm:pt-12">
    <form class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
      <div class="mb-6">
        <label class="block text-gray-700 text-sm font-bold mb-2" for="email">
          邮箱
          <i class="text-red-500 mx-1">*</i>
        </label>
        <input
          v-model="email"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          id="email"
          type="email"
        />
      </div>
      <button
        @click="sendResetEmail"
        class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        type="button"
      >发送重置邮件</button>
    </form>
    <p class="text-center text-gray-500 text-xs">&copy;2020 Memory Card. All rights reserved.</p>
  </div>
</template>

<script>
export default {
  name: "Login",
  data() {
    return {
      email: ""
    };
  },
  methods: {
    sendResetEmail() {
      if (!/^.+@.+$/.test(this.email)) {
        this.$notify.error({
          title: "错误",
          message: "邮箱格式错误"
        });
        return false;
      }
      this.$http.post("/password/email", { email: this.email }).then(res => {
        if (res.message == "success") {
          this.$notify.success({
            title: "成功",
            message: "发送重置邮件成功~请检查收件箱..."
          });
        }
      });
    }
  }
};
</script>