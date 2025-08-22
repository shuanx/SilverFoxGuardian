<template>
  <div class="loader"></div>
</template>

<style>
.loader {
  border: 6px solid rgba(0, 0, 0, 0.1);
  border-top: 6px solid #8a6d3b;
  border-radius: 50%;
  width: 50px;
  height: 50px;
  animation: spin 1s linear infinite;
  margin: 20px auto;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
