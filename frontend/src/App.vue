<template>
  <div class="container">
    <header class="header">
      <div class="header-row">
        <!-- 检测后才显示返回icon（无按钮背景，纯图标，透明） -->
        <img
            v-if="((currentPage === 'soft-detect' && hasCheckedSoft) || (currentPage === 'ip-detect' && hasCheckedIp))"
            src="../statics/return.png"
            alt="返回"
            class="return-icon-link"
            @click="goHome"
        />

        <h1 class="title" v-if="currentPage === 'home'">银狐-自动化排障工具</h1>
        <h1 class="title" v-else-if="currentPage === 'soft-detect'">银狐-第三方商业软件识别</h1>
        <h1 class="title" v-else-if="currentPage === 'ip-detect'">银狐-恶意外联IP自动化检测</h1>

        <!-- 检测按钮检测后放标题右侧，格式不变 -->
        <button
            v-if="currentPage === 'soft-detect' && hasCheckedSoft"
            class="btn-check"
            :disabled="isLoading"
            @click="checkProcesses"
        >
          {{ isLoading ? "正在加载..." : "检测" }}
        </button>
        <button
            v-if="currentPage === 'ip-detect' && hasCheckedIp"
            class="btn-check"
            :disabled="isLoadingIp"
            @click="checkIPConnections"
        >
          {{ isLoadingIp ? "正在加载..." : "检测" }}
        </button>
      </div>
      <link rel="icon" href="../statics/icon.ico" type="image/x-icon" />
    </header>

    <main>
      <!-- 首页选择 -->
      <div v-if="currentPage === 'home'" class="home-options">
        <div class="option-row">
          <span class="option-title">银狐-第三方商业软件识别</span>
          <button class="btn-check" @click="currentPage = 'soft-detect'">检测</button>
        </div>
        <div class="option-row">
          <span class="option-title">银狐-恶意外联IP自动化检测</span>
          <button class="btn-check" @click="currentPage = 'ip-detect'">检测</button>
        </div>
      </div>

      <!-- 第三方商业软件识别页面 -->
      <div v-else-if="currentPage === 'soft-detect'">
        <!-- 页面下方按钮，检测后隐藏 -->
        <div v-if="!hasCheckedSoft" class="main-btn-row">
          <button class="btn-check" :disabled="isLoading" @click="checkProcesses">
            {{ isLoading ? "正在加载..." : "检测" }}
          </button>
          <button class="btn-back" @click="goHome">返回</button>
        </div>
        <Loader v-if="isLoading" />
        <ErrorMessage v-if="errorMessage" :message="errorMessage" />
        <div v-if="processes.length > 0" class="results">
          <ProcessItem
              v-for="(process, index) in processes"
              :key="index"
              :data="process"
          />
        </div>
      </div>

      <!-- 恶意外联IP自动化检测页面 -->
      <div v-else-if="currentPage === 'ip-detect'">
        <!-- 页面下方按钮，检测后隐藏 -->
        <div v-if="!hasCheckedIp" class="main-btn-row">
          <button class="btn-check" :disabled="isLoadingIp" @click="checkIPConnections">
            {{ isLoadingIp ? "正在加载..." : "检测" }}
          </button>
          <button class="btn-back" @click="goHome">返回</button>
        </div>
        <Loader v-if="isLoadingIp" />
        <ErrorMessage v-if="ipErrorMessage" :message="ipErrorMessage" />
        <div v-if="ipConnections.length > 0" class="results">
          <table class="ip-table">
            <thead>
            <tr>
              <th>进程名</th>
              <th>PID</th>
              <th>目的IP:端口</th>
              <th>状态</th>
              <th>地理位置</th>
              <th>威胁情报</th>
              <th>威胁标签</th>
            </tr>
            </thead>
            <tbody>
            <tr
                v-for="(conn, idx) in ipConnections"
                :key="idx"
                :class="getTIClass(conn.tiResult)"
            >
              <td>{{ conn.process }}</td>
              <td>{{ conn.pid }}</td>
              <td>{{ conn.remote }}</td>
              <td>{{ conn.state }}</td>
              <td>{{ conn.geo || '-' }}</td>
              <td>{{ conn.Tags.length > 0 ? conn.Tags.join(', ') : '-' }}</td>
              <td>
                <span
                    v-if="getTIClass(conn.tiResult) === 'ti-safe'"
                    class="ti-result ti-safe"
                >
                  <span v-if="formatTI(conn.tiResult) === '安全'" class="safe-word">安全</span>
                  <span v-else>{{ formatTI(conn.tiResult) }}</span>
                </span>
                              <span
                                  v-else-if="getTIClass(conn.tiResult) === 'ti-malicious'"
                                  class="ti-malicious-text"
                              >
                  {{ formatTI(conn.tiResult) }}
                </span>
                              <span
                                  v-else
                                  class="ti-result ti-unknown"
                              >
                  {{ formatTI(conn.tiResult) }}
                </span>
              </td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>

    <footer class="footer">
      <p><strong>© 2025 SHAUN👑</strong></p>
    </footer>
  </div>
</template>

<script>
import Loader from "./components/ProcessLoader.vue";
import ProcessItem from "./components/ProcessItem.vue";
import ErrorMessage from "./components/ProcessErrorMessage.vue";

const mockData = `[
  {
    "processName": ["wrdlv4.exe", "winrdlv3.exe"],
    "describe": "ipguard",
    "isExist": "No",
    "pid": "",
    "connections": null
  },
  {
    "processName": ["NSecRTS.exe"],
    "describe": "Nsec(ping32)",
    "isExist": "No",
    "pid": "",
    "connections": null
  },
  {
    "processName": ["poda64.exe"],
    "describe": "固信",
    "isExist": "No",
    "pid": "",
    "connections": null
  },
  {
    "processName": ["ClashX"],
    "describe": "mac下测试使用",
    "isExist": "True",
    "pid": "902",
    "connections": [
      "223.5.5.5:443:ESTABLISHED [China-Hangzhou]",
      "113.240.72.99:443:ESTABLISHED [China-Qingyuan]"
    ]
  }
]`;

const mockIpData = `[
  {
    "process": "ClashX",
    "pid": "902",
    "remote": "223.5.5.5:443",
    "state": "ESTABLISHED",
    "geo": "China-Hangzhou",
    "tiResult": "未知",
    "Tags": "-"
  },
  {
    "process": "chrome.exe",
    "pid": "1103",
    "remote": "104.2.40.14:53",
    "state": "ESTABLISHED",
    "geo": "USA-Mountain View",
    "tiResult": "恶意",
    "Tags": "SilverFox"
  },
  {
    "process": "chrome.exe",
    "pid": "1103",
    "remote": "8.8.8.8:53",
    "state": "ESTABLISHED",
    "geo": "USA-Mountain View",
    "tiResult": "安全",
    "Tags": "-"
  }
]`;

export default {
  components: {
    Loader,
    ProcessItem,
    ErrorMessage,
  },
  data() {
    return {
      currentPage: "home", // home, soft-detect, ip-detect
      isLoading: false,
      errorMessage: "",
      processes: [],
      isLoadingIp: false,
      ipErrorMessage: "",
      ipConnections: [],
      hasCheckedSoft: false,
      hasCheckedIp: false,
    };
  },
  methods: {
    goHome() {
      this.currentPage = "home";
      this.errorMessage = "";
      this.processes = [];
      this.isLoading = false;
      this.ipErrorMessage = "";
      this.ipConnections = [];
      this.isLoadingIp = false;
      this.hasCheckedSoft = false;
      this.hasCheckedIp = false;
    },
    async checkProcesses() {
      this.isLoading = true;
      this.errorMessage = "";
      try {
        if (process.env.NODE_ENV === "development") {
          this.processes = JSON.parse(mockData);
        } else {
          const result = await window.checkProcesses();
          this.processes = JSON.parse(result);
        }
        this.hasCheckedSoft = true;
      } catch (error) {
        this.errorMessage = "加载数据的时候发生错误。。";
      } finally {
        this.isLoading = false;
      }
    },
    async checkIPConnections() {
      this.isLoadingIp = true;
      this.ipErrorMessage = "";
      try {
        if (process.env.NODE_ENV === "development") {
          this.ipConnections = JSON.parse(mockIpData);
        } else {
          console.log("开始进行外联IP测试")
          const result = await window.analyzeExternalConnections();
          console.log(result)
          this.ipConnections = JSON.parse(result);
        }
        this.hasCheckedIp = true;
      } catch (error) {
        this.ipErrorMessage = "加载外联IP数据时发生错误。。";
        console.log(error)
      } finally {
        this.isLoadingIp = false;
      }
    },
    getTIClass(ti) {
      if (!ti || ti === "未知") return "ti-unknown";
      if (ti === "安全") return "ti-safe";
      if (ti === "恶意") return "ti-malicious";
      return "ti-unknown";
    },
    formatTI(ti) {
      if (!ti || ti === "未知") return "未知";
      if (ti === "安全") return "安全";
      if (ti === "恶意") return "恶意";
      return ti;
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
  position: relative;
}

/* 加flex让按钮能和标题同行 */
.header-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 32px;
  position: relative;
}

/* 返回icon纯图标样式（无背景），可随背景色融合，具备hover放大高亮 */
.return-icon-link {
  width: 32px;
  height: 32px;
  cursor: pointer;
  margin-right: 10px;
  margin-left: 0;
  border-radius: 50%;
  background: transparent;
  transition: box-shadow 0.18s, transform 0.18s, filter 0.18s;
  box-shadow: none;
  filter: brightness(0.95);
}

.return-icon-link:hover {
  box-shadow: 0 0 10px 2px #e1cdb5cc;
  filter: brightness(1.1) drop-shadow(0 2px 6px #decba4bb);
  transform: scale(1.16);
}

/* 首页选项样式 */
.home-options {
  display: flex;
  flex-direction: column;
  gap: 35px;
  align-items: center;
  justify-content: center;
  margin-top: 80px;
}

.option-row {
  display: flex;
  align-items: center;
  gap: 30px;
  background: #fff9f3;
  border: 1px solid #decba4;
  border-radius: 18px;
  padding: 26px 50px;
  box-shadow: 0 8px 18px rgba(0, 0, 0, 0.11);
}

.option-title {
  font-size: 1.6rem;
  color: #8a6d3b;
  font-weight: bold;
  margin-right: 10px;
  width: 330px;
  display: inline-block;
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

.btn-back {
  background: #ecd9c6;
  color: #8a6d3b;
  border: none;
  border-radius: 25px;
  padding: 15px 40px;
  font-size: 1.5rem;
  font-weight: bold;
  cursor: pointer;
  margin-left: 12px;
  margin-right: 12px;
  box-shadow: 0 10px 15px rgba(0,0,0,0.2);
  transition: background 0.25s, transform 0.3s;
}
.btn-back:hover {
  background: #decba4;
  transform: scale(1.05);
}

/* 主体页面按钮组（居中一行） */
.main-btn-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 24px;
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

/* ip-detect 页面表格样式 */
.ip-table {
  width: 100%;
  border-collapse: collapse;
  background: #fff;
  font-size: 1.1rem;
}
.ip-table th, .ip-table td {
  border: 1px solid #decba4;
  padding: 10px 15px;
  text-align: center;
}
.ip-table thead {
  background: #f2e6d6;
}
.ip-table tbody tr:nth-child(odd) {
  background: #fdf8f3;
}

/* ====== 威胁情报美化 ====== */
/* ====== 威胁情报美化 ====== */
/* “安全”只有“安全”两个字是绿色，其它内容和未知一样灰色和字号 */

.ti-safe {
  /* 整体灰色，字号一致 */
  color: #b0b0b0 !important;
  font-weight: normal;
  background: none !important;
  border: none !important;
  box-shadow: none !important;
  padding: 0 !important;
  border-radius: 0 !important;
  font-size: 0.80em;
  vertical-align: middle;
}

/* “安全”两个字专用绿色高亮 */
.ti-safe .safe-word {
  color: #1bb669 !important;
  font-weight: bold;
  font-size: 1em;
}

/* “恶意”红色，字体略大，其他同普通 */
.ti-malicious,
.ti-malicious-text {
  color: #e74d3d !important;
  font-weight: bold;
  background: none !important;
  border: none !important;
  box-shadow: none !important;
  padding: 0 !important;
  border-radius: 0 !important;
  font-size: 0.88em;
  vertical-align: middle;
}

/* “未知”灰色，字号略小 */
.ti-unknown {
  color: #b0b0b0 !important;
  font-weight: normal;
  background: none !important;
  border: none !important;
  box-shadow: none !important;
  padding: 0 !important;
  border-radius: 0 !important;
  font-size: 0.80em;
  vertical-align: middle;
}

/* 页脚 */
.footer {
  margin-top: auto;
  padding: 20px;
  text-align: center;
  color: #8a6d3b;
  font-size: 1.5rem;
  font-weight: bold;
  font-style: italic;
  letter-spacing: 2px;
  text-shadow: 1px 1px 6px #decba4;
}
</style>