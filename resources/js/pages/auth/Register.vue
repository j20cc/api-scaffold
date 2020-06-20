<template>
  <div class="max-w-md mx-auto mb:pt-20 sm:pt-12">
    <form class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
      <div class="mb-4">
        <label class="block text-gray-700 text-sm font-bold mb-2" for="username">
          用户名
          <i class="text-red-500 mx-1">*</i>
        </label>
        <input
          v-model="formData.name"
          placeholder="请输入6-15位用户名"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          id="username"
          type="text"
        />
      </div>
      <div class="mb-4">
        <label class="block text-gray-700 text-sm font-bold mb-2" for="email">
          邮箱
          <i class="text-red-500 mx-1">*</i>
        </label>
        <input
          v-model="formData.email"
          placeholder="请输入有效邮箱"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          id="email"
          type="email"
        />
      </div>
      <div class="mb-4">
        <label class="block text-gray-700 text-sm font-bold mb-2" for="password">
          密码
          <i class="text-red-500 mx-1">*</i>
        </label>
        <input
          v-model="formData.password"
          placeholder="请输入6-15位密码"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
          id="password"
          type="password"
        />
      </div>
      <div class="mb-6">
        <label class="block text-gray-700 text-sm font-bold mb-2" for="password">
          确认密码
          <i class="text-red-500 mx-1">*</i>
        </label>
        <input
          v-model="formData.repassword"
          placeholder="请确认密码"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
          id="repassword"
          type="password"
        />
      </div>

      <button
        @click="handleSubmit"
        class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        type="button"
      >注册</button>
    </form>

    <CopyRight></CopyRight>
  </div>
</template>

<script>
import CopyRight from "../../components/CopyRight";
import helper from '../../libs/helper';

export default {
  name: "Register",
  data() {
    return {
      formData: {
        name: "",
        email: "",
        password: "",
        repassword: ""
      }
    };
  },
  methods: {
    handleSubmit() {
      if (this.validateForm()) {
        this.$http.post("/register", this.formData).then(res => {
          this.$notify.success({
            title: "成功",
            message: "注册成功~即将跳转到个人页面..."
          });
          helper.setUser(res)
          setTimeout(() => {
            this.$notify.closeAll();
            this.$router.push({
              path: "/user"
            });
          }, 2000);
        });
      }
    },
    validateForm() {
      let data = this.formData;
      if (data.name.length > 15 || data.name.length < 6) {
        this.$notify.error({
          title: "错误",
          message: "用户名应该为6-15个字符"
        });
        return false;
      }
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
          message: "密码应该为6-15个字符"
        });
        return false;
      }
      if (data.repassword != data.password) {
        this.$notify.error({
          title: "错误",
          message: "确认密码与密码不同"
        });
        return false;
      }
      return true;
    }
  },
  components: {
    CopyRight
  }
};
</script>