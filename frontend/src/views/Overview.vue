<template>
  <div>
    <div class="toolbar">
      <div>
        <h2 class="page-title">概览</h2>
        <p class="page-sub">工地、探方与出土文物登记统计</p>
      </div>
      <button class="btn secondary" @click="load">刷新</button>
    </div>

    <div v-if="error" class="error card">{{ error }}</div>

    <div class="stats">
      <div class="stat card">
        <div class="label">发掘工地</div>
        <div class="value">{{ data.siteCount ?? '-' }}</div>
      </div>
      <div class="stat card">
        <div class="label">探方/发掘单位</div>
        <div class="value">{{ data.unitCount ?? '-' }}</div>
      </div>
      <div class="stat card">
        <div class="label">出土文物总数</div>
        <div class="value">{{ data.findCount ?? '-' }}</div>
      </div>
      <div class="stat card">
        <div class="label">近7日新增文物</div>
        <div class="value">{{ data.last7DaysNewFinds ?? '-' }}</div>
        <div class="hint">按东八区自然日计</div>
      </div>
    </div>

    <div class="card">
      <h3>按工地对比</h3>
      <table class="table" v-if="data.bySite?.length">
        <thead>
          <tr>
            <th>工地名称</th>
            <th>探方数</th>
            <th>文物数</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="row in data.bySite"
            :key="row.siteId"
            class="link-row"
            title="查看该工地的探方"
            @click="goSiteUnits(row)"
          >
            <td>{{ row.siteName }}</td>
            <td>{{ row.unitCount }}</td>
            <td>{{ row.findCount }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else class="page-sub">暂无工地数据</p>
      <p class="page-sub hint">点击行跳转到该工地的探方筛选</p>
    </div>

    <div class="card">
      <h3>按器物类型统计</h3>
      <table class="table" v-if="data.byType?.length">
        <thead>
          <tr>
            <th>器物类型</th>
            <th>数量</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in data.byType" :key="row.artifactType">
            <td><span class="tag">{{ row.artifactType }}</span></td>
            <td>{{ row.count }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else class="page-sub">暂无统计数据</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api/http'

const router = useRouter()

const data = reactive({
  siteCount: 0,
  unitCount: 0,
  findCount: 0,
  last7DaysNewFinds: 0,
  byType: [],
  bySite: []
})
const error = ref('')

async function load() {
  error.value = ''
  try {
    const { data: res } = await api.get('/overview')
    Object.assign(data, res)
  } catch (e) {
    error.value = e.response?.data?.error || '加载失败'
  }
}

function goSiteUnits(row) {
  router.push({ name: 'units', query: { siteId: String(row.siteId) } })
}

onMounted(load)
</script>

<style scoped>
.stats {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

.stat .label {
  color: var(--muted);
  font-size: 0.9rem;
}

.stat .value {
  font-size: 2rem;
  font-weight: 700;
  margin-top: 0.35rem;
  color: var(--accent);
}

.stat .hint {
  margin-top: 0.25rem;
  color: var(--muted);
  font-size: 0.78rem;
}

.card {
  margin-bottom: 1rem;
}

.hint {
  margin-top: 0.5rem;
  font-size: 0.8rem;
}

.link-row {
  cursor: pointer;
}

.link-row:hover {
  background: rgba(139, 90, 43, 0.08);
}

h3 {
  margin: 0 0 0.75rem;
}

@media (max-width: 1000px) {
  .stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 800px) {
  .stats {
    grid-template-columns: 1fr;
  }
}
</style>
