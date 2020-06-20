<template>
  <div></div>
</template>

<script>
import helper from "../../libs/helper";
export default {
  name: "VerifyEmail",
  created() {
    if ("sign" in this.$route.query) {
      this.$http
        .post("/verification", { sign: this.$route.query.sign })
        .then(res => {
          helper.notifyAndRedirect("success", "验证成功~", "/user", 1000);
        })
        .catch(err => {
          this.$router.push({ path: "/" });
        });
    } else {
      helper.notifyAndRedirect("error", "注册失败~", "/", 1000);
    }
  }
};
</script>