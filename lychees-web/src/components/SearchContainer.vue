<script setup lang="ts">
import { useFocusStore } from '@/stores/searchStore'
import { useSearchHistory } from '@/stores/historyStore'
import { db } from '@/stores/db'

import { ref, onMounted, onUnmounted, reactive, inject } from 'vue'
import { unix } from '@/utils/tools'
import { clipboardKey } from '@/models/key'

const focusStore = useFocusStore()
const searchHistory = useSearchHistory()
const isHistory = ref<boolean>(false)
const word = ref<string>('')
let date = new Date()
const suggestArray = ref<string[]>([])
const clipboardEvent = inject(clipboardKey, (enevt: MouseEvent) => {
  console.log('上层组件找不到clipboardKey！')
})
const time = reactive({
  hours: '',
  minutes: '',
})

const getTime = () => {
  date = new Date()
  time.hours = date.getHours().toString().padStart(2, '0')

  time.minutes = date.getMinutes().toString().padStart(2, '0')
}
// 搜索框聚焦
const handleFocus = () => {
  focusStore.isFog = true
  focusStore.isFocus = true
  focusStore.isDisplayEngine = true
  if (!word.value) {
    getHistory()
  }
}
// 获取搜索历史的列表
const getHistory = () => {
  isHistory.value = true
  suggestArray.value = searchHistory.searchHistoryArray
}
// 搜索框失焦
const handleBlur = () => {
  focusStore.isFocus = false
}
// 推荐列表引擎
const suggestion = 'https://suggestion.baidu.com/su?wd='

// 时间标签悬停过渡动画
let animationTimeout: number
const timeStyles = ref({
  transform: '',
})
// 开始悬停
const startHover = () => {
  clearTimeout(Number(animationTimeout))

  timeStyles.value.transform = 'scale(1.15)'
  animationTimeout = setTimeout(() => {
    timeStyles.value.transform = 'scale(1.1)'
  }, 250)
}
// 结束悬停
const endHover = (enevt: Event) => {
  const timeElement = enevt.currentTarget as HTMLTimeElement
  clearTimeout(animationTimeout)
  timeStyles.value.transform = 'scale(0.95)'
  animationTimeout = setTimeout(() => {
    timeStyles.value.transform = ''
    timeElement.removeAttribute('style')
    clearTimeout(animationTimeout)
  }, 250)
}

let getTimeInterval: number
// 挂载后
onMounted(() => {
  getTime()
  getTimeInterval = setInterval(() => {
    getTime()
  }, 1000)
})

// 卸载后，移除定时器,防止内存溢出
onUnmounted(() => {
  getTimeInterval && clearInterval(getTimeInterval)
})

let engineIndex = ref(0)

const start = () => {
  // const wordTrim = word.value.trim()
  if (word.value) {
    searchHistory.add(word.value)
  }

  window.open(`${startEngineList[engineIndex.value].url}${encodeURIComponent(word.value)}`)
  hoverSuggestIndex.value = -1
  word.value = ''
}

const startEngineList = [
  {
    url: 'https://www.baidu.com/s?ie=utf-8&word=',
    name: 'baidu',
  },
  {
    url: 'https://www.bing.com/search?q=',
    name: 'bing',
  },
  {
    url: 'https://www.google.com/search?q=',
    name: 'google',
  },
]

const changeEngine = (index: number) => {
  // inputSearch.value.focus()
  // focusStore.onfocus()

  if (engineIndex.value !== index) {
    engineIndex.value = index
    // engineIndex.value = startEngineList.findIndex((element) =>element.name == item.name)
  }

  if (word.value != '') {
    start()
  }
}
function handleAnimationend(event: Event) {
  const timeElement = event.currentTarget as HTMLTimeElement
  timeElement.classList.remove('rubberBand')
}
const rubberBand = (event: Event) => {
  focusStore.change()
  const timeElement = event.currentTarget as HTMLTimeElement
  if (timeElement.classList.contains('rubberBand')) {
    timeElement.classList.remove('rubberBand')
  }
  // 我们的回调在下一次重绘页面之前被调用，避免在重绘之前被调用，而此时样式还没有被真正重新计算。
  requestAnimationFrame(() => {
    timeElement.classList.add('rubberBand')
  })
}
// 是否在合成中
let isComposite = false

// 如果已经有baidu则不添加属性
// if (!window.hasOwnProperty('baidu')) {

const baidu = {
  sug: ({ q, p, s }: { q: string; p: boolean; s: string[] }) => {
    db.searchSsuggestList.add({
      keyword: q,
      suggestList: s,
      time: unix(),
    })
    isHistory.value = false
    suggestArray.value = s
  },
}
// 允许重新定义，防止拿不到vue组件的实例
Object.defineProperty(window, 'baidu', {
  value: baidu,
  writable: true,
})

const hoverSuggestIndex = ref(-1)
function handleArrowUp(event: Event) {
  // event.preventDefault() // 阻止默认行为
  if (hoverSuggestIndex.value == -1 || hoverSuggestIndex.value == 0) {
    hoverSuggestIndex.value = suggestArray.value.length - 1
  } else {
    hoverSuggestIndex.value--
  }
  word.value = suggestArray.value[hoverSuggestIndex.value]

  // console.log('按下方向键上')
  // 执行向上操作
}
function handleArrowDown(event: Event) {
  // event.preventDefault() // 阻止默认行为
  if (hoverSuggestIndex.value == -1 || hoverSuggestIndex.value >= suggestArray.value.length - 1) {
    hoverSuggestIndex.value = 0
  } else {
    hoverSuggestIndex.value++
  }
  word.value = suggestArray.value[hoverSuggestIndex.value]
  // console.log('按下方向键下')
  // 执行向下操作
}
function handlePreventDefault(event: Event) {
  event.preventDefault() // 阻止默认行为
}
let valveTimer: number
// 获取推荐词列表
const handleInput = () => {
  if (isComposite) {
    // console.log('在合成中')
    return
  }
  valveTimer && clearTimeout(valveTimer)
  deBounce()
}
const handleCompositionStart = () => {
  // console.log('合成开始')

  isComposite = true
}
const handleCompositionEnd = () => {
  // console.log('合成结束')

  isComposite = false
  handleInput()
}
// 根据索引删除本地搜索记录
const handDeleteHistory = (index: number) => {
  searchHistory.cleanIndex(index)
}
const handleSuggestion = (keyword: string) => {
  word.value = keyword
  // console.log(keyword)

  start()
}
// 函数防抖
const deBounce = () => {
  if (hoverSuggestIndex.value != -1) {
    hoverSuggestIndex.value = -1
  }

  // 内容为空 啥也不管
  if (!word.value) {
    getHistory()
    return
  }

  // 清理时间久远的数据
  db.searchSsuggestList
    .where('time')
    .below(unix() - 604800)
    .toArray()
    .then(function (result) {
      // 处理查询结果
      if (result.length > 0) {
        const bulk: number[] = result.reduce((array: number[], element) => {
          if (element.id) {
            array.push(element.id)
          }
          return array
        }, [])
        // console.log('即将的数据有', bulk)
        db.searchSsuggestList.bulkDelete(bulk)
      } else {
        // console.log('没有需要清理的数据')
      }
    })
  // 查询存在的推荐列表
  db.searchSsuggestList
    .where('keyword')
    .equals(word.value)
    .first()
    .then(function (result) {
      if (result) {
        // 如果查询到
        console.log('记录存在')
        isHistory.value = false
        suggestArray.value = result.suggestList
      } else {
        // 没有查询到则发起函数防抖
        valveTimer = setTimeout(() => {
          // 使用示例
          const script = document.createElement('script')
          // 监听脚本加载完成后的回调
          script.onload = () => {
            // 此处递归调用，慎用警告！！！
            // deBounce()
            // console.log('网络搜索完成')

            // 在脚本加载完成后执行你的逻辑
            script.remove()
          }
          script.src = suggestion + encodeURIComponent(word.value) + '&t=' + Date.now()
          document.body.appendChild(script)
        }, 500)
      }
    })
}
</script>

<template>
  <!-- 时间容器 -->
  <div class="time-container">
    <!-- 时间 -->
    <time
      class="animated"
      @click="rubberBand"
      :style="timeStyles"
      @mouseover="startHover"
      @mouseout="endHover"
      @animationend="handleAnimationend"
    >
      {{ time.hours }}
      <span>:</span>
      {{ time.minutes }}
    </time>
  </div>

  <!-- 搜索框 -->
  <label
    class="search-container"
    :class="{ focus: focusStore.isFocus || focusStore.isDisplayEngine }"
  >
    <input
      v-model.trim="word"
      type="search"
      @focus="handleFocus"
      @blur="handleBlur"
      @keyup.enter="start"
      @input="handleInput"
      @compositionstart="handleCompositionStart"
      @compositionend="handleCompositionEnd"
      autocomplete="off"
      accesskey="s"
      @contextmenu.prevent="clipboardEvent"
      @keyup.up="handleArrowUp"
      @keyup.down="handleArrowDown"
      @keydown.up="handlePreventDefault"
      @keydown.down="handlePreventDefault"
    />
    <svg
      class="clipboardIcon icon cross"
      aria-hidden="true"
      v-show="word && focusStore.isDisplayEngine"
      @click="word = ''"
    >
      <use xlink:href="#Heising-cross"></use>
    </svg>
  </label>
  <!-- 切换引擎 -->
  <div class="engine-container" v-show="focusStore.isDisplayEngine">
    <span
      v-for="(item, index) in startEngineList"
      :key="index"
      @click="changeEngine(index)"
      :class="{ selected: index == engineIndex }"
    >
      <svg class="heisingfont icon heising-alipay" aria-hidden="true">
        <use :xlink:href="'#Heising-' + item.name"></use>
      </svg>
    </span>
  </div>
  <!-- 推荐列表 -->

  <ul class="suggestion" v-show="focusStore.isDisplayEngine && suggestArray.length != 0">
    <li
      v-for="(v, i) in suggestArray"
      @click.self="handleSuggestion(v)"
      :class="{ hover: hoverSuggestIndex == i }"
      :key="v"
    >
      <svg class="clipboardIcon icon" aria-hidden="true">
        <use :xlink:href="`#Heising-${isHistory ? 'counter-clockwise' : 'magnify'}`"></use>
      </svg>
      {{ v }}
      <svg
        class="clipboardIcon icon cross"
        aria-hidden="true"
        v-if="isHistory"
        @click="handDeleteHistory(i)"
      >
        <use xlink:href="#Heising-cross"></use>
      </svg>
    </li>
  </ul>
</template>

<style lang="scss">
.animated {
  -webkit-animation-duration: 1s;
  animation-duration: 1s;
  -webkit-animation-duration: 1s;
  animation-duration: 1s;
  -webkit-animation-fill-mode: both;
  animation-fill-mode: both;
}
@-webkit-keyframes rubberBand-animate {
  0% {
    -webkit-transform: scaleX(1);
    transform: scaleX(1);
  }

  30% {
    -webkit-transform: scale3d(1.25, 0.75, 1);
    transform: scale3d(1.25, 0.75, 1);
  }

  40% {
    -webkit-transform: scale3d(0.75, 1.25, 1);
    transform: scale3d(0.75, 1.25, 1);
  }

  50% {
    -webkit-transform: scale3d(1.15, 0.85, 1);
    transform: scale3d(1.15, 0.85, 1);
  }

  65% {
    -webkit-transform: scale3d(0.95, 1.05, 1);
    transform: scale3d(0.95, 1.05, 1);
  }

  75% {
    -webkit-transform: scale3d(1.05, 0.95, 1);
    transform: scale3d(1.05, 0.95, 1);
  }

  to {
    -webkit-transform: scaleX(1);
    transform: scaleX(1);
  }
}

@keyframes rubberBand-animate {
  0% {
    -webkit-transform: scaleX(1);
    transform: scaleX(1);
  }

  30% {
    -webkit-transform: scale3d(1.25, 0.75, 1);
    transform: scale3d(1.25, 0.75, 1);
  }

  40% {
    -webkit-transform: scale3d(0.75, 1.25, 1);
    transform: scale3d(0.75, 1.25, 1);
  }

  50% {
    -webkit-transform: scale3d(1.15, 0.85, 1);
    transform: scale3d(1.15, 0.85, 1);
  }

  65% {
    -webkit-transform: scale3d(0.95, 1.05, 1);
    transform: scale3d(0.95, 1.05, 1);
  }

  75% {
    -webkit-transform: scale3d(1.05, 0.95, 1);
    transform: scale3d(1.05, 0.95, 1);
  }

  to {
    -webkit-transform: scaleX(1);
    transform: scaleX(1);
  }
}

.rubberBand {
  -webkit-animation-name: rubberBand-animate;
  animation-name: rubberBand-animate;
}

.time-container {
  position: fixed;
  top: 100px;
  left: 50%;
  transform: translateX(-50%);
  cursor: pointer;

  time {
    display: block;
    // max-width: 300px;
    color: #fff;
    font-weight: lighter;
    font-size: 40px;
    // white-space: nowrap;
    text-overflow: ellipsis;
    transition: all 0.25s;
    // transition: transform 0.3s;
  }
}
.search-container {
  position: fixed;
  top: 180px;
  left: 50%;
  border-radius: 20px;
  transform: translateX(-50%);
  width: 255px;
  max-width: 80%;
  transition:
    width 0.25s,
    backdrop-filter 0.25s;
  &.focus,
  &:hover {
    width: 530px;
    backdrop-filter: blur(15px);
    input[type='search'] {
      color: #333;
    }
  }

  input[type='search'] {
    width: 100%;
    height: 45px;
    padding: 0 50px;
    border-radius: 20px;
    color: transparent;
    background-color: rgba(255, 255, 255, 0.5);
    box-shadow: 0 3px 12px rgba(10, 10, 10, 0.3);
    outline: 0;
    border: none;
    text-align: center;
    opacity: 0.6;
    // backdrop-filter: blur(10px);
    transition: all 0.25s;

    &:hover,
    &:focus {
      background-color: rgba(255, 255, 255, 0.75);
      // width: 530px;
      color: #333;
    }

    &:focus {
      background-color: rgba(255, 255, 255, 0.9);
      color: black;
    }

    &::-webkit-search-cancel-button {
      -webkit-appearance: none;
    }

    &::selection {
      // color: red;
      background: var(--theme-color-heising);
    }

    &::placeholder {
      color: #fff;
      /* text-shadow: 0 0 10px rgba(0,0,0,.3); */
      // -webkit-text-stroke: thin rgba(0, 0, 0, 0.3);
      transition: 0.25s;
    }

    &:hover::placeholder {
      color: rgba(255, 255, 255, 0.8);
      // -webkit-text-stroke-width: 0;
    }

    &:focus::placeholder {
      // -webkit-text-stroke-width: 0;
      color: transparent;
    }
  }
}
.engine-container {
  display: flex;
  justify-content: space-between;
  position: absolute;
  top: 235px;
  left: 50%;
  transform: translateX(-50%);
  width: 255px;
  // transition: 0.25s;

  span {
    position: relative;
    width: 50px;
    height: 28px;
    line-height: 28px;
    text-align: center;
    background-color: rgba(255, 255, 255, 0.5);

    padding: 5px 30px;

    border-radius: 18px;
    cursor: pointer;
    transition: 0.3s;

    // 给中间加空隙
    &:nth-child(2) {
      margin-left: 10px;
      margin-right: 10px;
    }

    &:hover {
      background-color: rgba(255, 255, 255, 1);

      .heisingfont {
        transform: translate(-50%, -50%) scale(1.2);
        transition: all 0.25s;
      }
    }

    &.selected {
      background-color: rgba(255, 255, 255, 0.9);
    }
  }
}

.suggestion {
  background-color: rgba(233, 233, 233, 0.3);
  width: 525px;
  max-width: 80%;
  overflow-y: hidden;
  border-radius: 15px;
  z-index: 40;
  position: absolute;
  top: 283px;
  left: 50%;
  backdrop-filter: blur(20px);
  transform: translateX(-50%);

  // transition: 0.3s;

  li {
    position: relative;
    height: 30px;
    // border-radius: 10px;
    padding-right: 15px;
    text-indent: 15px;
    line-height: 30px;
    // border: #fff;
    // color: var(--heising-background-color);
    color: rgba(255, 255, 255, 0.6);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    transition: 0.15s cubic-bezier(0.88, 0.12, 0.39, 1);
    // transition: all 0.4s cubic-bezier(0.31, -0.1, 0.43, 1.59);
    .cross {
      display: none;
    }
    &:hover,
    &.hover {
      backdrop-filter: blur(10px);
      color: var(--heising-color);
      text-indent: 20px;
      background-color: rgba(233, 233, 233, 0.3);
    }
    &:hover {
      .cross {
        display: block;
      }
    }
  }
}
.cross {
  cursor: pointer;
  position: absolute;
  right: 15px;
  top: 50%;
  color: rgba(255, 255, 255, 0.6);
  transform: translateY(-50%);
  &:hover {
    color: rgba(255, 255, 255, 1);
  }
}
@media (max-width: 1024px) {
  .time-container {
    top: 50px;
  }

  .search-container {
    top: 130px;
  }

  .engine-container {
    top: 185px;
  }

  .suggestion {
    top: 233px;
  }
}
</style>
