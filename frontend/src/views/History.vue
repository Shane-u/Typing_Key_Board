<template>
  <div class="history-container">
    <h1>练习历史记录</h1>
    
    <div class="history-panel" v-if="historyRecords.length > 0">
      <div class="history-table">
        <table>
          <thead>
            <tr>
              <th>日期</th>
              <th>素材</th>
              <th>错误数</th>
              <th>准确率</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="record in historyRecords" :key="record.id">
              <td>{{ formatDate(record.completed_at) }}</td>
              <td>{{ getMaterialTitle(record.text_id) }}</td>
              <td>{{ record.error_count }}</td>
              <td>{{ record.accuracy }}%</td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <div class="history-stats">
        <h2>你的统计数据</h2>
        <div class="stat-card">
          <div class="stat-value">{{ totalPractices }}</div>
          <div class="stat-label">总练习次数</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ averageAccuracy }}%</div>
          <div class="stat-label">平均准确率</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ totalErrors }}</div>
          <div class="stat-label">总错误数</div>
        </div>
      </div>
    </div>
    
    <div class="no-history" v-else>
      <p>你还没有完成任何练习。</p>
      <router-link to="/materials" class="material-btn">开始练习</router-link>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'HistoryView',
  data() {
    return {
      historyRecords: [],
      userId: '',
      materials: []
    }
  },
  computed: {
    totalPractices() {
      return this.historyRecords.length;
    },
    averageAccuracy() {
      if (this.historyRecords.length === 0) return 0;
      const sum = this.historyRecords.reduce((acc, record) => acc + record.accuracy, 0);
      return Math.round(sum / this.historyRecords.length);
    },
    totalErrors() {
      return this.historyRecords.reduce((acc, record) => acc + record.error_count, 0);
    }
  },
  methods: {
    loadUserId() {
      const savedUserId = localStorage.getItem('userId')
      if (savedUserId) {
        this.userId = savedUserId
      } else {
        this.userId = 'user_' + Math.random().toString(36).substring(2, 9)
        localStorage.setItem('userId', this.userId)
      }
    },
    async loadHistory() {
      try {
        const response = await axios.get(`/api/history/${this.userId}`)
        this.historyRecords = response.data
      } catch (error) {
        console.error('加载历史记录失败:', error)
        this.historyRecords = []
      }
    },
    async loadMaterials() {
      try {
        // 加载内置素材
        const builtInMaterials = [
          { id: 'builtin-1', title: '编程关键词' },
          { id: 'builtin-2', title: '常用编程变量' },
          { id: 'builtin-3', title: 'Git命令' },
          { id: 'builtin-4', title: 'HTML标签' },
          { id: 'builtin-5', title: 'CSS属性' }
        ];
        
        // 自定义素材
        const response = await axios.get(`/api/texts/${this.userId}`);
        const customMaterials = response.data;
        
        // 合并内置和自定义素材
        this.materials = [...builtInMaterials, ...customMaterials];
      } catch (error) {
        console.error('加载素材失败:', error);
        this.materials = [];
      }
    },
    getMaterialTitle(textId) {
      const material = this.materials.find(m => m.id === textId);
      return material ? material.title : '未知素材';
    },
    formatDate(dateString) {
      if (!dateString) return '未知';
      
      const date = new Date(dateString);
      return new Intl.DateTimeFormat('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      }).format(date);
    }
  },
  mounted() {
    this.loadUserId();
    this.loadMaterials();
    this.loadHistory();
  }
}
</script>

<style scoped>
.history-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

h1 {
  color: #2c3e50;
  text-align: center;
  margin-bottom: 30px;
}

h2 {
  color: #2c3e50;
  margin-bottom: 20px;
}

.history-panel {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
}

.history-table, .history-stats {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

th {
  font-weight: bold;
  color: #2c3e50;
  background-color: #f9f9f9;
}

tbody tr:hover {
  background-color: #f9f9f9;
}

.stat-card {
  background-color: #f9f9f9;
  border-radius: 6px;
  padding: 15px;
  margin-bottom: 15px;
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #42b983;
  margin-bottom: 5px;
}

.stat-label {
  color: #666;
}

.no-history {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 40px;
  text-align: center;
}

.material-btn {
  display: inline-block;
  margin-top: 15px;
  background-color: #42b983;
  color: white;
  text-decoration: none;
  padding: 10px 20px;
  border-radius: 4px;
  font-weight: bold;
}

@media (max-width: 768px) {
  .history-panel {
    grid-template-columns: 1fr;
  }
  
  th, td {
    padding: 8px;
    font-size: 14px;
  }
}
</style> 