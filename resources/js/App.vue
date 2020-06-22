<template>
  <div id="app">
    <MyHeader v-if="showHeader"></MyHeader>

    <transition name="fade" mode="out-in">
      <router-view></router-view>
    </transition>
  </div>
</template>

<script>
import MyHeader from "./components/MyHeader.vue";
import helper from "./libs/helper";
import "./assets/main.css";

export default {
  name: "App",
  beforeCreate() {
    let user = helper.getUser();
    if (user != null) {
      helper.setUser(user);
    }
  },
  computed: {
    showHeader() {
      return this.$route.name != "card";
    }
  },
  components: {
    MyHeader
  }
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
