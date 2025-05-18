<template>
  <div class="app-container">
    <header>
      <h1>打字练习</h1>
      <nav>
        <router-link to="/" class="nav-link">练习</router-link>
        <router-link to="/materials" class="nav-link">素材库</router-link>
        <router-link to="/history" class="nav-link">历史记录</router-link>
      </nav>
    </header>
    <main>
      <router-view />
    </main>
    <footer>
      <p>&copy; {{ new Date().getFullYear() }} 打字练习应用</p>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Arial', sans-serif;
  line-height: 1.6;
  color: #333;
  background-color: #f5f5f5;
}

.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

header {
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1rem 2rem;
  text-align: center;
}

header h1 {
  color: #2c3e50;
  margin-bottom: 1rem;
}

nav {
  display: flex;
  justify-content: center;
  gap: 1.5rem;
}

.nav-link {
  color: #2c3e50;
  text-decoration: none;
  font-weight: bold;
  padding: 0.5rem 0;
  position: relative;
}

.nav-link:after {
  content: '';
  position: absolute;
  width: 0;
  height: 2px;
  bottom: 0;
  left: 0;
  background-color: #42b983;
  transition: width 0.3s ease;
}

.nav-link:hover:after,
.router-link-active:after {
  width: 100%;
}

.router-link-active {
  color: #42b983;
}

main {
  flex: 1;
}

footer {
  text-align: center;
  padding: 1rem;
  background-color: #fff;
  box-shadow: 0 -2px 4px rgba(0, 0, 0, 0.05);
  margin-top: 2rem;
}
</style> 