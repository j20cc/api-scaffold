<template>
  <header class="body-font text-gray-700">
    <div class="container mx-auto flex justify-between p-5 items-center">
      <router-link to="/" class="flex title-font font-medium items-center text-gray-900 mb-0">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          stroke="currentColor"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          class="w-10 h-10 text-white p-2 bg-indigo-500 rounded-full"
          viewBox="0 0 24 24"
        >
          <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
        </svg>
        <span class="ml-3 text-xl">{{title}}</span>
      </router-link>

      <nav v-if="user" class="md:ml-auto flex flex-wrap items-center text-base justify-center">
        <router-link class="mr-5 hover:text-gray-900" to="/">{{user.name}}</router-link>
        <span class="mr-5 hover:text-gray-900 cursor-pointer" @click="logout">退出</span>
      </nav>

      <router-link to="/login" v-if="showLoginButton">
        <button
          class="inline-flex items-center bg-gray-200 border-0 py-1 px-3 focus:outline-none hover:bg-gray-300 rounded text-base mt-0"
        >
          登录
          <svg
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            class="w-4 h-4 ml-1"
            viewBox="0 0 24 24"
          >
            <path d="M5 12h14M12 5l7 7-7 7" />
          </svg>
        </button>
      </router-link>
      <router-link to="/register" v-if="showRegisterButton">
        <button
          class="inline-flex items-center bg-gray-200 border-0 py-1 px-3 focus:outline-none hover:bg-gray-300 rounded text-base mt-0"
        >
          注册
          <svg
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            class="w-4 h-4 ml-1"
            viewBox="0 0 24 24"
          >
            <path d="M5 12h14M12 5l7 7-7 7" />
          </svg>
        </button>
      </router-link>
    </div>
  </header>
</template>

<script>
import { mapState } from "vuex";
import helper from "../libs/helper";

export default {
  name: "MyHeader",
  data() {
    return {
      title: process.env.MIX_APP_NAME
    };
  },
  computed: {
    showLoginButton() {
      return this.$route.path != "/login" && !this.user;
    },
    showRegisterButton() {
      return this.$route.path == "/login" && !this.user;
    },
    ...mapState(["user"])
  },
  methods: {
    logout() {
      this.$confirm("确认退出吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          helper.removeUser();
          this.$notify.success({
            title: "成功",
            message: "退出成功~即将跳到登录页"
          });
          setTimeout(() => {
            this.$notify.closeAll();
            this.$router.push({
              path: "/login"
            });
          }, 2000);
        })
        .catch(() => {
          console.log("cancle");
        });
    }
  }
};
</script>