<template>
  <div class="materials-container">
    <h1>练习素材</h1>
    
    <div class="materials-panel">
      <div class="materials-list">
        <h2>内置素材</h2>
        <div 
          v-for="(material, index) in builtInMaterials" 
          :key="'builtin-' + index"
          class="material-item"
          @click="selectMaterial(material)"
        >
          <div class="material-preview">{{ material.title }}</div>
          <div class="material-desc">{{ truncateText(material.text, 100) }}</div>
        </div>
        
        <h2>自定义素材</h2>
        <div 
          v-for="material in customMaterials" 
          :key="material.id"
          class="material-item"
        >
          <div class="material-preview">{{ material.title || '自定义素材' }}</div>
          <div class="material-desc">{{ truncateText(material.text, 100) }}</div>
          <div class="material-actions">
            <button class="action-btn select-btn" @click="selectMaterial(material)">选择</button>
            <button class="action-btn delete-btn" @click="deleteMaterial(material.id)">删除</button>
          </div>
        </div>
      </div>
      
      <div class="add-material-panel">
        <h2>添加新素材</h2>
        <div class="form-group">
          <label for="material-title">标题:</label>
          <input 
            type="text" 
            id="material-title" 
            v-model="newMaterial.title" 
            placeholder="输入素材标题..."
          >
        </div>
        <div class="form-group">
          <label for="material-text">文本内容:</label>
          <textarea 
            id="material-text" 
            v-model="newMaterial.text" 
            placeholder="输入练习文本内容..."
            rows="8"
          ></textarea>
        </div>
        <button class="save-btn" @click="saveMaterial">保存素材</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'MaterialsView',
  data() {
    return {
      builtInMaterials: [
        {
          id: 'builtin-1',
          title: '编程关键词',
          text: 'function return const let var class interface extends implements static void public private protected import export default switch case break continue for while do if else try catch finally throw new typeof instanceof'
        },
        {
          id: 'builtin-2',
          title: '常用编程变量',
          text: 'data value item element node index count sum total result temp key name id title text content message error status code flag isValid isReady isLoading isActive current prev next start end min max size length width height'
        },
        {
          id: 'builtin-3',
          title: 'Git命令',
          text: 'git init clone add commit push pull fetch merge rebase checkout branch status log diff reset stash cherry-pick revert bisect remote tag show reflog config help'
        },
        {
          id: 'builtin-4',
          title: 'HTML标签',
          text: 'html head body div span p h1 h2 h3 h4 h5 h6 ul ol li a img input button form label select option table tr td th script style link meta title'
        },
        {
          id: 'builtin-5',
          title: 'CSS属性',
          text: 'color background margin padding border width height display position top left right bottom font size weight flex grid align justify transform transition animation opacity visibility overflow float z-index content box-sizing'
        }
      ],
      customMaterials: [],
      userId: 'user_default',
      newMaterial: {
        title: '',
        text: ''
      }
    }
  },
  methods: {
    truncateText(text, maxLength) {
      if (text.length <= maxLength) return text;
      return text.substring(0, maxLength) + '...';
    },
    loadUserId() {
      const savedUserId = localStorage.getItem('userId')
      if (savedUserId) {
        this.userId = savedUserId
      } else {
        this.userId = 'user_' + Math.random().toString(36).substring(2, 9)
        localStorage.setItem('userId', this.userId)
      }
    },
    async loadCustomMaterials() {
      try {
        const response = await axios.get(`/api/texts/${this.userId}`)
        this.customMaterials = response.data
        
        // 添加标题属性，如果API中没有提供
        this.customMaterials = this.customMaterials.map(material => ({
          ...material,
          title: material.title || '自定义素材'
        }))
      } catch (error) {
        console.error('加载自定义素材失败:', error)
      }
    },
    selectMaterial(material) {
      localStorage.setItem('selectedMaterial', JSON.stringify(material))
      this.$router.push('/')
    },
    async saveMaterial() {
      if (!this.newMaterial.text.trim()) {
        alert('请输入文本内容')
        return
      }
      
      try {
        const response = await axios.post('/api/texts', {
          user_id: this.userId,
          text: this.newMaterial.text,
          title: this.newMaterial.title || '自定义素材'
        })
        
        // 清空表单
        this.newMaterial = {
          title: '',
          text: ''
        }
        
        // 重新加载自定义素材
        this.loadCustomMaterials()
      } catch (error) {
        // console.error('保存素材失败:', error)
      }
    },
    async deleteMaterial(id) {
      try {
        await axios.delete(`/api/texts/${id}`)
        // 重新加载自定义素材
        this.loadCustomMaterials()
      } catch (error) {
        // console.error('删除素材失败:', error)
      }
    }
  },
  mounted() {
    this.loadUserId()
    this.loadCustomMaterials()
  }
}
</script>

<style scoped>
.materials-container {
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
  margin-bottom: 15px;
  margin-top: 20px;
}

.materials-panel {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.materials-list, .add-material-panel {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.material-item {
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 4px;
  margin-bottom: 15px;
  cursor: pointer;
  transition: all 0.3s;
}

.material-item:hover {
  background-color: #f9f9f9;
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.material-preview {
  font-weight: bold;
  margin-bottom: 5px;
  color: #2c3e50;
}

.material-desc {
  color: #666;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.material-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px;
}

.action-btn {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  margin-left: 10px;
}

.select-btn {
  background-color: #42b983;
  color: white;
}

.delete-btn {
  background-color: #ff5757;
  color: white;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  color: #666;
  font-weight: bold;
}

input[type="text"], textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  font-family: monospace;
}

textarea {
  resize: vertical;
}

.save-btn {
  background-color: #42b983;
  color: white;
  border: none;
  padding: 10px 15px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  width: 100%;
}

@media (max-width: 768px) {
  .materials-panel {
    grid-template-columns: 1fr;
  }
}
</style> 