<template>
  <div class="max-w-md mx-auto mb:pt-20 sm:pt-12">
    <form class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
      <div class="mb-4">
        <label class="block text-gray-700 text-sm font-bold mb-2" for="email">邮箱</label>
        <input
          v-model="formData.email"
          placeholder="请输入登录邮箱"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          id="email"
          type="email"
        />
      </div>
      <div class="mb-6">
        <label class="block text-gray-700 text-sm font-bold mb-2" for="password">密码</label>
        <input
          v-model="formData.password"
          placeholder="请输入登录密码"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
          id="password"
          type="password"
        />
      </div>
      <div class="flex items-center justify-between">
        <button
          @click="handleSubmit"
          class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
          type="button"
        >登录</button>
        <router-link
          to="/forget"
          class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800"
          href="#"
        >忘记密码?</router-link>
      </div>
    </form>
    <p class="text-center text-gray-500 text-xs">&copy;2020 Memory Card. All rights reserved.</p>
  </div>
</template>

<script>
import helper from "../../libs/helper";
export default {
  name: "Login",
  data() {
    return {
      formData: {
        email: "",
        password: ""
      }
    };
  },
  methods: {
    handleSubmit() {
      console.log(this.$store.state);
      if (this.validateForm()) {
        this.$http.post("/login", this.formData).then(res => {
          if ("token" in res) {
            this.$notify.success({
              title: "成功",
              message: "登录成功~即将跳转到个人页面..."
            });
            helper.setUser(res);
            setTimeout(() => {
              this.$notify.closeAll();
              this.$router.push({
                path: "/index"
              });
            }, 2000);
          }
        });
      }
    },
    validateForm() {
      let data = this.formData;
      if (!/^.+@.+$/.test(data.email)) {
        this.$notify.error({
          title: "错误",
          message: "邮箱格式错误"
        });
        return false;
      }
      if (data.password.length > 15 || data.password.length < 6) {
        this.$notify.error({
          title: "错误",
          message: "密码长度应该为6-15个字符"
        });
        return false;
      }
      return true;
    }
  }
};
</script>