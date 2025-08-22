<template>
  <div class="container">
    <header class="header">
<!--      <img src="/path-to-logo.png" alt="National Logo" class="logo" />-->
      <h1 class="title">第三方白利用识别</h1>
    </header>

    <main>
      <button class="btn-check" :disabled="isLoading" @click="checkProcesses">
        {{ isLoading ? "正在加载..." : "检测" }}
      </button>

      <Loader v-if="isLoading" />
      <ErrorMessage v-if="errorMessage" :message="errorMessage" />
      <div v-if="processes.length > 0" class="results">
        <ProcessItem
            v-for="(process, index) in processes"
            :key="index"
            :data="process"
        />
      </div>
    </main>

    <footer class="footer">
      <p>© 2025 Shaun</p>
    </footer>
  </div>
</template>

<script>
import Loader from "./components/ProcessLoader.vue";
import ProcessItem from "./components/ProcessItem.vue";
import ErrorMessage from "./components/ProcessErrorMessage.vue";

export default {
  components: {
    Loader,
    ProcessItem,
    ErrorMessage,
  },
  data() {
    return {
      isLoading: false,
      errorMessage: "",
      processes: [],
    };
  },
  methods: {
    async checkProcesses() {
      this.isLoading = true;
      this.errorMessage = "";
      try {
        const result = await window.checkProcesses(); // 调用 Webview 绑定的 Go 函数
        this.processes = JSON.parse(result);
      } catch (error) {
        this.errorMessage = "加载数据的时候发生错误。。";
      } finally {
        this.isLoading = false;
      }
    },
  },
};
</script>

<style>
/* 主容器 */
.container {
  font-family: "Georgia", serif;
  color: #333;
  background: linear-gradient(135deg, #fdfcfb, #e2d1c3);
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
  box-sizing: border-box;
}

/* 头部设计 */
.header {
  width: 100%;
  max-width: 1200px;
  text-align: center;
  margin-bottom: 30px;
  padding: 20px;
  border-bottom: 2px solid #decba4;
}

.logo {
  width: 100px;
  height: auto;
  margin-bottom: 15px;
}

.title {
  font-size: 3rem;
  font-weight: bold;
  color: #8a6d3b;
  font-family: "Times New Roman", Times, serif;
}

/* 按钮样式 */
.btn-check {
  background: linear-gradient(90deg, #c4a77d, #8a6d3b);
  color: white;
  border: none;
  border-radius: 25px;
  padding: 15px 40px;
  font-size: 1.5rem;
  font-weight: bold;
  cursor: pointer;
  box-shadow: 0 10px 15px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease-in-out;
}

.btn-check:hover {
  background: linear-gradient(90deg, #8a6d3b, #c4a77d);
  transform: scale(1.05);
}

.btn-check:disabled {
  background: #e4d5c4;
  color: #8a6d3b;
  cursor: not-allowed;
  box-shadow: none;
}

/* 结果容器 */
.results {
  margin-top: 30px;
  width: 100%;
  max-width: 800px;
  background: #fff9f3;
  border: 1px solid #decba4;
  border-radius: 15px;
  padding: 30px;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

/* 页脚 */
.footer {
  margin-top: auto;
  padding: 20px;
  text-align: center;
  color: #8a6d3b;
  font-size: 1rem;
  font-style: italic;
}
</style>
