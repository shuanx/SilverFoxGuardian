<template>
  <div class="error">{{ message }}</div>
</template>

<script>
export default {
  props: {
    message: {
      type: String,
      required: true,
    },
  },
};
</script>

<style>
.error {
  color: #d9534f;
  font-weight: bold;
  font-family: "Georgia", serif;
  margin-top: 20px;
  text-align: center;
  font-size: 1.2rem;
}
</style>
