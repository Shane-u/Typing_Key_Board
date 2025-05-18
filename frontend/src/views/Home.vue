<template>
  <div class="home-container">
    <div class="typing-area">
      <div class="practice-panel" v-if="selectedMaterial">
        <div class="target-text">
          <template v-for="(char, index) in targetChars" :key="index">
            <span 
              :class="{
                'correct': typedChars[index] === char,
                'incorrect': typedChars[index] !== undefined && typedChars[index] !== char,
                'current': currentIndex === index
              }"
            >{{ char }}</span>
          </template>
        </div>
        <div class="typing-input-container">
          <textarea 
            ref="typingInput"
            v-model="typedText"
            class="typing-input" 
            :disabled="isCompleted"
            @input="handleInput"
            @keydown="handleKeyDown"
            placeholder="点击这里开始输入..."
          ></textarea>
        </div>
        <div class="stats-panel">
          <div class="stat-item">
            <span class="stat-label">进度:</span>
            <span class="stat-value">{{ progress }}%</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">错误:</span>
            <span class="stat-value">{{ errorCount }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">准确率:</span>
            <span class="stat-value">{{ accuracy }}%</span>
          </div>
          <div v-if="isCompleted" class="completion-message">
            练习完成！
            <button @click="resetTyping" class="reset-btn">重新开始</button>
          </div>
        </div>
      </div>
      <div v-else class="no-material-message">
        <p>请先选择一个练习素材。</p>
        <router-link to="/materials" class="material-btn">前往素材库</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'HomeView',
  data() {
    return {
      targetText: '',
      targetChars: [],
      typedText: '',
      typedChars: [],
      currentIndex: 0,
      errorCount: 0,
      isStarted: false,
      isCompleted: false,
      errors: {},
      userId: 'user_' + Math.random().toString(36).substring(2, 9), 
      selectedMaterial: null
    }
  },
  computed: {
    progress() {
      if (this.targetChars.length === 0) return 0
      return Math.floor((this.currentIndex / this.targetChars.length) * 100)
    },
    accuracy() {
      if (this.targetChars.length === 0) return 100
      // 使用错误计数来计算准确率
      const totalChars = this.targetChars.length
      return Math.floor(((totalChars - this.errorCount) / totalChars) * 100)
    }
  },
  methods: {
    loadSettings() {
      // 从本地存储加载用户ID
      const savedUserId = localStorage.getItem('userId')
      if (savedUserId) {
        this.userId = savedUserId
      } else {
        localStorage.setItem('userId', this.userId)
      }
      
      // 从本地存储加载选中的素材
      const selectedMaterial = localStorage.getItem('selectedMaterial')
      if (selectedMaterial) {
        this.selectedMaterial = JSON.parse(selectedMaterial)
        this.targetText = this.selectedMaterial.text
        this.startTyping()
      }
    },
    startTyping() {
      if (!this.targetText) return
      
      this.targetChars = Array.from(this.targetText)
      this.typedText = ''
      this.typedChars = []
      this.currentIndex = 0
      this.errorCount = 0
      this.isStarted = true
      this.isCompleted = false
      this.errors = {}
      
      this.$nextTick(() => {
        if (this.$refs.typingInput) {
          this.$refs.typingInput.focus()
        }
      })
    },
    resetTyping() {
      this.startTyping()
    },
    handleInput(event) {
      this.typedText = event.target.value
      this.typedChars = Array.from(this.typedText)
      
      // 更新当前位置
      this.currentIndex = this.typedChars.length
      
      // 检查完成状态
      if (this.currentIndex >= this.targetChars.length) {
        this.isCompleted = true
        
        // 计算最终准确率（目标字符数-错误字符数）/目标字符数 * 100%
        const finalAccuracy = Math.floor(((this.targetChars.length - this.errorCount) / this.targetChars.length) * 100)
        // 保存到历史记录
        this.saveHistory(finalAccuracy)
      }
    },
    handleKeyDown(event) {
      // 允许删除错误输入
      if (event.key === 'Backspace' || event.key === 'Delete') {
        return true
      }
      
      // 检查是否输入错误
      if (event.key.length === 1) {
        const expectedChar = this.targetChars[this.currentIndex]
        // 只有当该位置还没有被记录为错误时（ !this.errors[this.currentIndex]）才增加错误计数
        if (event.key !== expectedChar && !this.errors[this.currentIndex]) {
          this.errorCount++
          this.errors[this.currentIndex] = true
        }
      }
    },
    async saveHistory(finalAccuracy) {
	  // 异步保存为历史记录
      try {
        // 确保selectedMaterial存在且有id
        if (!this.selectedMaterial || !this.selectedMaterial.id) {
          return;
        }
        
        const historyData = {
          user_id: this.userId,
          text_id: String(this.selectedMaterial.id),
          error_count: this.errorCount,
          accuracy: finalAccuracy || this.accuracy,
          completed_at: new Date().toISOString()
        };
        
        const response = await axios.post('/api/history', historyData);
      } catch (error) {
        if (error.response) {
          console.error('错误信息:', error.response.data);
        }
      }
    }
  },
  mounted() {
    this.loadSettings()
  }
}
</script>

<style scoped>
.home-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.typing-area {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.no-material-message {
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

.practice-panel {
  padding: 20px;
}

.target-text {
  font-size: 18px;
  line-height: 1.8;
  margin-bottom: 20px;
  background-color: #f9f9f9;
  padding: 15px;
  border-radius: 4px;
  min-height: 100px;
  white-space: pre-wrap;
  font-family: monospace;
}

.target-text span {
  border-radius: 2px;
  padding: 2px 0;
}

.target-text span.correct {
  background-color: rgba(66, 185, 131, 0.2);
  color: #42b983;
}

.target-text span.incorrect {
  background-color: rgba(255, 87, 87, 0.2);
  color: #ff5757;
}

.target-text span.current {
  border-bottom: 2px solid #42b983;
}

.typing-input-container {
  margin-bottom: 20px;
}

.typing-input {
  width: 100%;
  height: 80px;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  resize: none;
  font-size: 16px;
  font-family: monospace;
}

.stats-panel {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.stat-item {
  margin-right: 20px;
}

.stat-label {
  color: #666;
  margin-right: 5px;
}

.stat-value {
  font-weight: bold;
  color: #2c3e50;
}

.completion-message {
  margin-top: 10px;
  padding: 10px;
  background-color: #f0f7f3;
  border-radius: 4px;
  color: #42b983;
  font-weight: bold;
  display: flex;
  align-items: center;
}

.reset-btn {
  background-color: #42b983;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
  margin-left: 10px;
  font-size: 14px;
}

@media (max-width: 768px) {
  .home-container {
    padding: 10px;
  }
  
  .target-text {
    font-size: 16px;
  }
  
  .stats-panel {
    flex-direction: column;
  }
  
  .stat-item {
    margin-bottom: 10px;
  }
}
</style> 