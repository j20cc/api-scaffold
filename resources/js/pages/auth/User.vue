<template>
  <div class="max-w-md mx-auto mb:pt-20 sm:pt-12">
    <form class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4" v-if="user">
      <div class="mb-6">
        <label class="block text-gray-700 text-sm font-bold mb-2">
          头像
        </label>
        <img :src="'https://picsum.photos/200/200'" class="w-16 h-16 rounded-full" alt="">
      </div>
      <div class="mb-6">
        <label class="block text-gray-700 text-sm font-bold mb-2">
          邮箱
        </label>
        <input
          v-model="user.email"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          disabled
          type="email"
        />
      </div>
      <div class="mb-6">
        <label class="block text-gray-700 text-sm font-bold mb-2">
          注册时间
        </label>
        <input
          v-model="user.created_at"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          disabled
          type="text"
        />
      </div>

      <button
        @click="sendVerifyEmail"
        v-if="!user.email_verified_at"
        class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        type="button"
      >发送激活邮件</button>
    </form>

    <CopyRight></CopyRight>
  </div>
</template>

<script>
import CopyRight from "../../components/CopyRight";

export default {
  name: "Login",
  data() {
    return {
      user: null
    };
  },
  mounted() {
    this.getUser()
  },
  methods: {
    sendVerifyEmail() {
      this.$http.post("/verification/email", { email: this.user.email }).then(res => {
        if (res.message == "success") {
          this.$notify.success({
            title: "成功",
            message: "发送激活邮件成功~请检查收件箱..."
          });
        }
      });
    },
    getUser() {
      this.$http.get("/profile").then(res => {
        this.user = res
      });
    }
  },
  components: {
    CopyRight
  }
};
</script>