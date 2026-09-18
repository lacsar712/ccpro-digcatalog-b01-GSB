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
      <div class="stat card highlight">
        <div class="label">近 7 日新增文物</div>
        <div class="value">{{ data.last7DaysNewFinds ?? '-' }}</div>
        <div class="hint">{{ rangeText }}（UTC+8）</div>
      </div>
    </div>

    <div class="grid-2">
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

      <div class="card">
        <h3>按工地对比</h3>
        <table class="table" v-if="data.bySite?.length">
          <thead>
            <tr>
              <th>工地</th>
              <th>文物数</th>
              <th>探方数</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in sortedBySite"
              :key="row.siteId"
              class="clickable-row"
              @click="goUnits(row.siteId)"
            >
              <td>{{ row.siteName }}</td>
              <td>{{ row.findCount }}</td>
              <td>{{ row.unitCount }}</td>
              <td class="row-actions" @click.stop>
                <button class="btn secondary small" @click="goUnits(row.siteId)">探方筛选</button>
                <button class="btn secondary small" @click="goSite(row.siteId)">查看工地</button>
              </td>
            </tr>
          </tbody>
        </table>
        <p v-else class="page-sub">暂无工地数据</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api/http'

const router = useRouter()

const data = reactive({
  siteCount: 0,
  unitCount: 0,
  findCount: 0,
  byType: [],
  last7DaysNewFinds: 0,
  bySite: [],
  last7DaysRange: null
})
const error = ref('')

const rangeText = computed(() => {
  const r = data.last7DaysRange
  if (!r || !r.start) return '近 7 个自然日'
  return `${r.start} ~ ${r.end || r.start}`
})

// 文物数多的工地排前面，便于横向对比。
const sortedBySite = computed(() =>
  [...(data.bySite || [])].sort((a, b) => {
    if (b.findCount !== a.findCount) return b.findCount - a.findCount
    return b.unitCount - a.unitCount
  })
)

async function load() {
  error.value = ''
  try {
    const { data: res } = await api.get('/overview')
    Object.assign(data, res)
  } catch (e) {
    error.value = e.response?.data?.error || '加载失败'
  }
}

function goUnits(siteId) {
  router.push({ name: 'units', query: { siteId: String(siteId) } })
}

function goSite(siteId) {
  router.push({ name: 'sites', query: { focusId: String(siteId) } })
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

.stat.highlight {
  border-left: 3px solid var(--accent);
}

.stat .hint {
  margin-top: 0.35rem;
  color: var(--muted);
  font-size: 0.8rem;
}

.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

h3 {
  margin: 0 0 0.75rem;
}

.clickable-row {
  cursor: pointer;
}

.clickable-row:hover {
  background: var(--surface-hover, rgba(0, 0, 0, 0.03));
}

.row-actions {
  display: flex;
  gap: 0.4rem;
}

@media (max-width: 900px) {
  .stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .grid-2 {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 560px) {
  .stats {
    grid-template-columns: 1fr;
  }
}
</style>
