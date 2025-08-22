<template>
  <div class="process-item">
    <div class="process-header">
      <h2>{{ data.describe }}</h2>
      <span class="status" :class="{ running: isRunning, 'not-running': !isRunning }">
        {{ isRunning ? "运行中" : "未运行" }}
      </span>
    </div>
    <div class="details">
      <table>
        <thead>
        <tr>
          <th>字段</th>
          <th>值</th>
        </tr>
        </thead>
        <tbody>
        <tr>
          <td>进程名字</td>
          <td>{{ data.processName.join(", ") }}</td>
        </tr>
        <tr>
          <td>PID</td>
          <td>{{ data.pid || "无" }}</td>
        </tr>
        </tbody>
      </table>
      <h5 class="connections-title">外联信息</h5>
      <div v-if="!data.connections || data.connections.length === 0" class="no-connections">
        暂无相关外联信息。
      </div>
      <table v-else>
        <thead>
        <tr>
          <th>IP:PORT</th>
          <th>状态</th>
          <th>来源国家</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(conn, index) in parsedConnections" :key="index">
          <td>{{ conn.ip }}</td>
          <td>{{ conn.status }}</td>
          <td>{{ conn.country }}</td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    data: {
      type: Object,
      required: true,
    },
  },
  computed: {
    isRunning() {
      return this.data.isExist === "True";
    },
    parsedConnections() {
      return this.data.connections.map((conn) => {
        const [ip, port, status, country] = conn.match(/(.*?):(.*?):(.*?)\s\[(.*?)\]/).slice(1);
        return { ip: `${ip}:${port}`, status, country };
      });
    },
  },
};
</script>

<style>
.process-item {
  border: 1px solid #decba4;
  border-radius: 15px;
  padding: 20px;
  margin-bottom: 30px;
  background: #fffaf2;
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.1);
}

.process-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.process-header h2 {
  font-size: 1.8rem;
  font-family: "Georgia", serif;
  color: #8a6d3b;
  margin: 0;
}

.status {
  padding: 5px 15px;
  border-radius: 20px;
  font-size: 1rem;
  font-family: "Georgia", serif;
}

.status.running {
  background: #8a6d3b;
  color: white;
}

.status.not-running {
  background: #e4d5c4;
  color: #8a6d3b;
}

.details table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

.details table th,
.details table td {
  border: 1px solid #decba4;
  padding: 10px;
  color: #333;
  text-align: left;
}

.details table th {
  background: #f9f4ef;
  font-weight: bold;
  font-family: "Georgia", serif;
}
</style>
