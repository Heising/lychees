<script setup lang="ts">
import { useFocusStore } from '@/stores/searchStore'
import { useNavStore } from '@/stores/navStore'
import {
  bookmarksDeleteApi,
  bookmarksMoveApi,
  bookmarksSwapApi,
  bookmarksUpdateApi,
} from '@/apis/bookmark'

import type { AddItem, Item, SwapItem } from '@/models'

import { inject, nextTick, onMounted, onUnmounted, ref } from 'vue'
import CornerContainer from '@/components/CornerContainer.vue'
import { ElNotification, type FormInstance } from 'element-plus'
import { openUrl } from '@/utils/tools'
import {
  extractWebsiteInfo,
  querySearch,
  removeSymbolElements,
  bookmarksRules,
  reValidate,
} from '@/utils/tools'
import { clipboardKey, contextMenuStyle } from '@/models/key'
import './menu.scss'
const focusStore = useFocusStore()
const navStore = useNavStore()
const color = 'var(--heising-color)'
const scrollIndex = ref<number>(0)
const clipboardEvent = inject(clipboardKey, (enevt: MouseEvent) => {
  console.log('上层组件找不到clipboardKey！')
})
let flag = true
// const ulPage = ref<HTMLUListElement>()
const ulPage = ref()

// const scrollSnapType = ref<string>('y proximity')

const asidePage = ref()
let sectionsArray: HTMLElement[]
const bookmarksRuleFormRef = ref<FormInstance | null>(null)
async function getBookmarks() {
  // 保留旧的第三方link
  const iconfontLink = navStore.bookmarks.iconfontLink

  await navStore.getBookmarks().catch((error) => {
    console.error(error.code, error.config.url)
  })

  // 添加第三方symbol 如果不一样，则重新渲染注入svg
  if (iconfontLink !== navStore.bookmarks.iconfontLink) {
    rendering()
  }
}
function rendering() {
  // console.log(navStore.bookmarks.iconfontLink)
  removeSymbolElements()
  if (navStore.bookmarks.iconfontLink) {
    const script: HTMLScriptElement = document.createElement('script')
    script.onload = () => {
      // console.log('注入svg完成')
      // 在脚本加载完成后执行你的逻辑
      script.remove()
    }
    script.src = navStore.bookmarks.iconfontLink

    document.body.appendChild(script)
  }
}
// 节流阀
let bookmarksTimeout: number
let throttle = false
function handleVisible() {
  if (throttle) return
  throttle = true
  // ulPage.value?.onscroll(new Event('onscroll'))
  if (document.visibilityState === 'visible') {
    // console.log('回到页面')
    getBookmarks()
  }
  // 节流十秒
  bookmarksTimeout = setTimeout(() => {
    throttle = false
    clearTimeout(bookmarksTimeout)
  }, 10000)
}

function handleScroll(e: Event) {
  // sectionsArray = Array.from(ulPage.value?.children) as HTMLUListElement[]
  sectionsArray = Array.from(ulPage.value?.children) as HTMLUListElement[]

  sectionsArray.forEach(function (item, index) {
    // console.log(`元素是${index}`, ulPage.value.scrollTop - item.offsetTop)
    if (
      ulPage.value.scrollTop - item.offsetTop < 10 &&
      ulPage.value.scrollTop - item.offsetTop > -10
    ) {
      scrollIndex.value = index
      return
    }
    // 上一级着色是 如果是0则不操作
    if (
      index &&
      ulPage.value.scrollTop - item.offsetTop > -50 &&
      ulPage.value.scrollTop - item.offsetTop < -10
    ) {
      scrollIndex.value = index - 1
    }
  })
}

onUnmounted(() => {
  document.removeEventListener('visibilitychange', handleVisible)
  document.removeEventListener('click', handleClick)
})
onMounted(async () => {
  rendering()
  await getBookmarks()

  // 切换回到页面都更新一下
  document.addEventListener('visibilitychange', handleVisible)
})

const openUrlSelf = (url: string) => {
  window.open(url, '_self')
  contextMenuTriggerHide()
}
const changeScroll = async (index: number) => {
  if (!flag) {
    return
  }

  flag = false
  nextTick(() => {
    ulPage.value.scroll({
      top: ulPage.value.children[index].offsetTop,
      behavior: 'smooth',
    })
  })

  flag = true
}
// 右键菜单逻辑
const contextMenuTrigger = ref<HTMLUListElement | null>(null)

let currentEditingRowIndex = -1
let currentEditingColIndex = -1

// 当前编辑的元素
const currentEditingItem = ref<Item>({
  title: '',
  url: '',
  icon: '',
  isSvg: false,
  turn: false,
})

// 删除某个元素
// const deleteItem = async () => {}

const confirmDeleteItem = () => {
  navStore.bookmarks.arrayBookmarks[currentEditingRowIndex].splice(currentEditingColIndex, 1)
  bookmarksDeleteApi(currentEditingRowIndex, currentEditingColIndex)
    .then((response) => {
      response.updateAt && (navStore.bookmarks.updateAt = response.updateAt)
      ElNotification({
        title: '删除成功',
        type: 'success',
      })
    })
    .catch(() => {
      ElNotification({
        title: '删除失败',
        type: 'error',
      })

      getBookmarks()
    })

  contextMenuTriggerHide()
}
const cancelEvent = () => {
  contextMenuTriggerHide()
}

// 隐藏右键菜单
const contextMenuTriggerHide = () => {
  if (contextMenuTrigger.value) {
    contextMenuTrigger.value.removeAttribute('style')
  }
}
// 打开右键菜单
const showContextMenu = (event: MouseEvent, rowIndex: number, colIndex: number, item: Item) => {
  contextMenuTriggerHide()

  // 给当前要编辑的元素索引赋值
  currentEditingRowIndex = rowIndex
  currentEditingColIndex = colIndex

  currentEditingItem.value = { ...item }
  let x = event.clientX // 获取鼠标的水平坐标
  let y = event.clientY // 获取鼠标的垂直坐标
  let distanceX = window.innerWidth - x // 鼠标与页面右边的距离
  let distanceY = window.innerHeight - y // 鼠标与页面底部的距离

  if (contextMenuTrigger.value !== null) {
    // 先让其可见才有高度宽度
    contextMenuTrigger.value.style.display = 'block'
    // 返回元素内容的总高度，包括溢出部分 包括不可见部分
   

    // 如果屏幕距离小于菜单，则补偿
    if (distanceX < contextMenuTrigger.value.offsetWidth) {
      x -= contextMenuTrigger.value.offsetWidth - distanceX
    }
    if (distanceY < contextMenuTrigger.value.offsetHeight) {
      y -= contextMenuTrigger.value.offsetHeight - distanceY
    }

    // 将contextMenuStyle对象中的样式应用于DOM元素
    Object.assign(contextMenuTrigger.value.style, contextMenuStyle)

    contextMenuTrigger.value.style.top = y.toString() + 'px'
    contextMenuTrigger.value.style.left = x.toString() + 'px'
  }
}
// 添加点击事件监听器到document
document.addEventListener('click', handleClick)
function handleClick(event: Event) {
  // 检查点击的元素是否是要隐藏的元素
  // 检查点击的元素是否是要隐藏的元素或其子元素

  if (
    // event.target !== contextMenuTrigger.value &&
    contextMenuTrigger.value !== null &&
    !contextMenuTrigger.value.contains(event.target as Node)
  ) {
    // 如果不是要隐藏的元素，则隐藏它
    // contextMenuTrigger.value.style.display = 'none'
    contextMenuTrigger.value.removeAttribute('style')
  }
}

// 移动端右键菜单逻辑
let timer: number

const handleTouchstart = (event: TouchEvent) => {
  // 创建一个鼠标事件
  const ev = new MouseEvent('contextmenu', {})
  timer = setTimeout(function () {
    event.target?.dispatchEvent(ev)
  }, 700)
}

const handleTouchend = (event: TouchEvent) => {
  clearTimeout(timer)
}

// let dragInit: SwapItem = {
//   rowIndex: -1,
//   colIndex: -1,
// }
let dragFrom: SwapItem = {
  rowIndex: -1,
  colIndex: -1,
}

const handleDragstart = (e: DragEvent, rowIndex: number, colIndex: number) => {
  dragFrom = {
    rowIndex: -1,
    colIndex: -1,
  }

  console.log('开始拖拽', e.target, rowIndex, colIndex)
  dragFrom = {
    rowIndex,
    colIndex,
  }
}
const handleDragend = (e: DragEvent) => {
  clearDropStyle()
}
const handleDragover = (e: DragEvent) => {
  e.preventDefault()
}
const handleDragenter = (e: DragEvent) => {
  clearDropStyle()

  const dropNode = e.currentTarget as HTMLLIElement
  dropNode.classList.add('drop-over')
}
// 被放置到有效的放置目标上时触发
const handleDrop = (e: DragEvent, rowIndex: number, colIndex: number) => {
  clearDropStyle()
  if (dragFrom.colIndex == colIndex) {
    // 拖拽元素是同一个
    // console.log('拖拽元素是同一个')

    return
  }
  // 先备份拖拽大小
  const dragItemBackupIconSize =
    navStore.bookmarks.arrayBookmarks[dragFrom.rowIndex][dragFrom.colIndex].iconSize
  // 备份被放置的大小
  const dropItemBackupIconSize = navStore.bookmarks.arrayBookmarks[rowIndex][colIndex].iconSize

  // 备份被放置
  const dropItem = navStore.bookmarks.arrayBookmarks[rowIndex][colIndex]
  // 开始放置play
  navStore.bookmarks.arrayBookmarks[rowIndex][colIndex] =
    navStore.bookmarks.arrayBookmarks[dragFrom.rowIndex][dragFrom.colIndex]
  navStore.bookmarks.arrayBookmarks[rowIndex][colIndex].iconSize = dropItemBackupIconSize

  navStore.bookmarks.arrayBookmarks[dragFrom.rowIndex][dragFrom.colIndex] = dropItem
  navStore.bookmarks.arrayBookmarks[dragFrom.rowIndex][dragFrom.colIndex].iconSize =
    dragItemBackupIconSize
  // console.log('拖拽到哪里', e.currentTarget, rowIndex, colIndex)
  // 已经交换
  bookmarksSwapApi(dragFrom, { rowIndex, colIndex })
    .then((response) => {
      response.updateAt && (navStore.bookmarks.updateAt = response.updateAt)
      ElNotification({
        title: '交换成功',
        type: 'success',
      })
    })
    .catch(() => {
      ElNotification({
        title: '交换失败',
        type: 'error',
      })

      getBookmarks()
    })
}

// 清除样式
function clearDropStyle() {
  document.querySelectorAll('.drop-over').forEach((node) => {
    node.classList.remove('drop-over')
  })
}

const props = {
  expandTrigger: 'hover' as const,
}
// 编辑元素
const dialogFormVisible = ref(false)

interface Options {
  value: number
  label: string
  children?: Options[]
}
let positionOptions: Options[] = []

const position = ref([-1, -1])
let isMove = false
const handleChange = (value: any) => {
  isMove = true
  // console.log(value)
}
const handleEditingItem = () => {
  isMove = false
  position.value = [currentEditingRowIndex, currentEditingColIndex]

  contextMenuTriggerHide()
  bookmarksRuleFormRef.value?.clearValidate()
  dialogFormVisible.value = true

  // 这段卡
  positionOptions = navStore.bookmarks.arrayBookmarks.map((item, rowIndex) => {
    const children = item.map((_, colIndex) => {
      return {
        value: colIndex,
        label: '位置 ' + (colIndex + 1),
      }
    })
    children.push({
      value: item.length,
      label: '最后 ' + (item.length + 1),
    })

    return {
      value: rowIndex,
      label: '页面 ' + (rowIndex + 1),
      children: children,
    }
  })

  positionOptions.push({
    value: navStore.bookmarks.arrayBookmarks.length,
    label: '新增 ' + (navStore.bookmarks.arrayBookmarks.length + 1),
    children: [{ value: 6, label: '位置 ' + 1 }],
  })
}

const bookmarksUpdate = async () => {
  if (!bookmarksRuleFormRef.value) {
    return
  }
  let validate = false
  await bookmarksRuleFormRef.value.validate((valid, fields) => {
    if (valid) {
      validate = true
    } else {
      validate = false
    }
  })
  if (!validate) {
    return
  }
  dialogFormVisible.value = false

  extractWebsiteInfo(currentEditingItem)

  // console.log('添加的索引是', currentEditingRowIndex.value)
  let res: Promise<AddItem>
  let from = { rowIndex: currentEditingRowIndex, colIndex: currentEditingColIndex }
  if (isMove) {
    if (navStore.bookmarks.arrayBookmarks.length <= position.value[0]) {
      // 判断需要新增页面

      navStore.bookmarks.arrayBookmarks[currentEditingRowIndex].splice(currentEditingColIndex, 1)
      navStore.bookmarks.arrayBookmarks.push([
        {
          ...currentEditingItem.value,
        },
      ])
      position.value[1] = 0
      // 坐标修改为数组最后
      position.value[0] = navStore.bookmarks.arrayBookmarks.length - 1
    } else if (navStore.bookmarks.arrayBookmarks[position.value[0]].length <= position.value[1]) {
      // 不需要新增页面 但是位置比长度大或者等于直接放最后
      navStore.bookmarks.arrayBookmarks[currentEditingRowIndex].splice(currentEditingColIndex, 1)

      navStore.bookmarks.arrayBookmarks[position.value[0]].push({ ...currentEditingItem.value })
    } else {
      // 同一个页面插入 先删除后添加
      navStore.bookmarks.arrayBookmarks[currentEditingRowIndex].splice(currentEditingColIndex, 1)
      navStore.bookmarks.arrayBookmarks[position.value[0]].splice(position.value[1], 0, {
        ...currentEditingItem.value,
      })
    }
    // console.log(position.value);
    changeScroll(position.value[0])

    res = bookmarksMoveApi(
      from,
      { rowIndex: position.value[0], colIndex: position.value[1] },
      currentEditingItem.value
    )
    from = { rowIndex: position.value[0], colIndex: position.value[1] }
  } else {
    navStore.bookmarks.arrayBookmarks[currentEditingRowIndex][currentEditingColIndex] =
      currentEditingItem.value
    res = bookmarksUpdateApi(from, currentEditingItem.value)
  }

  res
    .then((response) => {
      response.updateAt && (navStore.bookmarks.updateAt = response.updateAt)
      ElNotification({
        title: '修改成功',
        type: 'success',
      })
      if (response.updateAt) {
        navStore.bookmarks.updateAt = response.updateAt
      }

      navStore.bookmarks.arrayBookmarks[from.rowIndex][from.colIndex] = response.data
    })
    .catch(() => {
      ElNotification({
        title: '修改失败',
        type: 'error',
      })
      getBookmarks()
    })
}
// 转码html字符
function decodeHtml(html: string) {
  let txt = document.createElement('textarea')
  txt.innerHTML = html
  return txt.value
}
</script>
<template>
  <CornerContainer
    @changeScroll="changeScroll"
    @getBookmarks="getBookmarks"
    v-show="!focusStore.isFog && !focusStore.isDisplayEngine"
  ></CornerContainer>
  <transition-group name="fade">
    <div
      class="nav-container"
      ref="ulPage"
      @scroll="handleScroll"
      v-show="!focusStore.isFog && !focusStore.isDisplayEngine"
      key="container"
    >
      <ul
        class="page"
        v-for="(v, rowIndex) in navStore.bookmarks.arrayBookmarks"
        :key="rowIndex"
        @dragover="handleDragover"
        @dragend="handleDragend"
        @touchstart="handleTouchstart"
        @touchend="handleTouchend"
      >
        <li
          @dragenter="handleDragenter"
          @drop="handleDrop($event, rowIndex, colIndex)"
          class="custom-nav"
          :class="item.iconSize ? 'icon-size-' + item.iconSize : ''"
          :data-content-after="decodeHtml(item.title)"
          v-for="(item, colIndex) in v"
          :key="colIndex + item.url"
          :style="{
            '--heising-background-color': item.turn ? item.color : color,
            '--item-color': item.turn ? color : item.color,
          }"
          @contextmenu.prevent="showContextMenu($event, rowIndex, colIndex, item)"
        >
          <a
            @dragstart="handleDragstart($event, rowIndex, colIndex)"
            @click.prevent="openUrl(item.url)"
            class="btn"
            :class="item.isSvg ? 'svg' : ''"
            draggable="true"
          >
            <!-- <div
              class="shade"
              :style="{
                backgroundColor: item.turn ? color : item.color,
                '--shade': item.turn ? color : item.color,
              }"
            ></div> -->
            <!-- 生成svg -->
            <svg class="heisingfont icon" aria-hidden="true" v-if="item.isSvg">
              <use :xlink:href="'#' + item.icon"></use>
            </svg>
            <img class="heisingfont icon" :src="item.icon" alt="" srcset="" loading="lazy" v-else />
          </a>
        </li>
      </ul>
    </div>
    <!-- 滚动条 -->
    <aside
      class="aside"
      key="aside"
      v-if="!focusStore.isFog && navStore.bookmarks.arrayBookmarks?.length > 1"
    >
      <ul class="aside-list" id="aside" ref="asidePage">
        <li
          v-for="(_, index) in navStore.bookmarks.arrayBookmarks"
          :class="{ active: index == scrollIndex }"
          :key="index"
          @click="changeScroll(index)"
        ></li>
      </ul>
    </aside>
  </transition-group>

  <ul class="menu" ref="contextMenuTrigger">
    <li class="menuItem" @click="openUrlSelf(currentEditingItem!.url)">当前页面打开</li>
    <li class="menuItem" @click="handleEditingItem">编辑这个书签</li>
    <el-popconfirm
      confirm-button-text="确认"
      cancel-button-text="取消"
      title="确认删除这个书签吗？"
      @confirm="confirmDeleteItem"
      @cancel="cancelEvent"
      :hide-icon="true"
      confirm-button-type="warning"
      :hide-after="0"
    >
      <template #reference>
        <li class="menuItem warn">删除这个书签</li>
      </template>
    </el-popconfirm>
  </ul>

  <el-dialog v-model="dialogFormVisible" title="编辑该书签">
    <el-form :model="currentEditingItem" :rules="bookmarksRules" ref="bookmarksRuleFormRef">
      <el-form-item label="请编辑链接" label-width="100px" prop="url">
        <el-input
          v-model.trim="currentEditingItem.url"
          autocomplete="off"
          clearable
          @contextmenu.prevent="clipboardEvent"
          @input="reValidate($event, bookmarksRuleFormRef)"
        />
      </el-form-item>

      <el-form-item label="请编辑标题" label-width="100px">
        <el-input
          v-model.trim="currentEditingItem.title"
          autocomplete="off"
          placeholder="为空则自动识别"
          clearable
          @contextmenu.prevent="clipboardEvent"
        />
      </el-form-item>

      <el-form-item label="请选择样式" label-width="100px">
        <el-radio-group v-model="currentEditingItem.isSvg" class="ml-4">
          <el-radio :value="true" size="large">svg</el-radio>
          <el-radio :value="false" size="large">url</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="请更改颜色" label-width="100px">
        <el-color-picker v-model="currentEditingItem.color" :disabled="!currentEditingItem.isSvg" />
      </el-form-item>

      <el-form-item label="请选择翻转" label-width="100px">
        <el-checkbox
          v-model="currentEditingItem.turn"
          label="翻转颜色"
          size="large"
          :disabled="!currentEditingItem.isSvg"
        />
      </el-form-item>

      <el-form-item label="请选择大小" label-width="100px">
        <el-radio-group v-model="currentEditingItem.iconSize" class="ml-4">
          <el-radio :value="undefined" size="large">1x1</el-radio>
          <el-radio :value="'1x2'" size="large">1x2</el-radio>
          <el-radio :value="'2x1'" size="large">2x1</el-radio>
          <el-radio :value="'2x3'" size="large">2x3</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="请输入图标" label-width="100px">
        <el-autocomplete
          v-model.trim="currentEditingItem.icon"
          :fetch-suggestions="currentEditingItem.isSvg ? querySearch : []"
          :debounce="0"
          clearable
          placeholder="选择SVG或者输入图标的URL"
          @contextmenu.prevent="clipboardEvent"
        />
      </el-form-item>

      <el-form-item label="移动书签到" label-width="100px">
        <el-cascader
          v-model="position"
          :options="positionOptions"
          :props="props"
          @change="handleChange"
          placeholder="移动书签到"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取消</el-button>
        <el-button type="primary" @click="bookmarksUpdate">确认</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss">
/* 选择前四个元素 */
// &:nth-child(-n + 4) {
//   /* 添加样式 */
//   width: 24px;
//   margin-right: 8px;
//   display: inline-block;
//   text-align: center;
//   &.active {
//     color: white;
//     background-color: rgba(0, 0, 0, 0.3);
//   }
// }

.aside {
  position: fixed;
  top: 50%;
  right: 2.6042vw;
  z-index: 1;
  transform: translateY(-50%);

  li {
    width: 8px;
    /* height: 8px; */
    height: 8px;

    background-color: rgba(204, 204, 204, 0.5);

    border-radius: 4px;
    margin-bottom: 8px;
    transition: all 0.3s;
    cursor: pointer;
    box-shadow: 0 3px 12px rgba(10, 10, 10, 0.3);

    &:first-child.active {
      background-color: var(--heising-background-color);
    }

    &.active {
      height: 24px;

      background: var(--theme-color-heising);
    }
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s;
}

.fade-enter,
.fade-leave-to

/* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}

// .v-enter-active,
// .v-leave-active {
//   transition: all 0.25s ease;
// }
// .v-enter-from,
// .v-leave-to {
//   opacity: 0;
//   transform: translateY(10px);
// }
.nav-container {
  position: absolute;
  top: 273px;
  left: 50%;
  height: calc(100% - 283px);
  transform: translateX(-50%);
  overflow: auto;
  scroll-snap-type: y mandatory;
  scrollbar-width: none;

  &::-webkit-scrollbar {
    display: none;
  }
}

.page {
  box-sizing: border-box;
  scroll-snap-align: start;
  padding-top: 10px;
  padding-bottom: 30px;

  min-height: 100%;
  --has-columns: 13;
  --has-columns-length: 90px;
  --has-rows: 90px;
  --has-column-gap: 20px;
  --has-row-gap: 30px;
  --icon-radius: 20px;
  display: grid;
  grid-template-columns: repeat(var(--has-columns), var(--has-columns-length));
  grid-auto-rows: var(--has-rows);
  gap: var(--has-column-gap);
  row-gap: var(--has-row-gap);
  grid-auto-flow: dense;
  animation: shadeAnimation linear both;
  // 触发动画时机
  animation-timeline: view();
  // 防止太快触发动画
  animation-range: entry 35% cover 0%;

  .custom-nav {
    display: flex;
    position: relative;
    border-radius: var(--icon-radius);
    background-color: var(--heising-background-color);
    transition: all 0.2s linear;
;
    &.drop-over {
      opacity: 0.3;
    }

    &:hover {
      transform: translateY(-5px);

      .btn {
        .heisingfont {
          color: var(--heising-background-color);
          transform: translate(-50%, -50%) scale(1.2);
        }

        &.svg::before {
          top: -15%;
          left: -15%;
        }
      }
    }

    &::after {
      content: attr(data-content-after);
      font-family: 'Microsoft YaHei';
      position: absolute;
      font-size: 12px;
      bottom: -20px;
      right: 50%;
      transform: translateX(50%);
      color: var(--heising-color);
      width: 100%;
      text-align: center;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .btn {
      overflow: hidden;
      position: relative;
      width: 100%;
      height: 100%;
      text-align: center;
      border-radius: var(--icon-radius);
      transition:
        transform 0.2s linear,
        background-color 0.4s linear;

      .heisingfont {
        color: var(--item-color);
        font-size: 38px;
        transition: all 0.4s cubic-bezier(0.31, -0.1, 0.43, 1.59);
      }
    }
  }

  @media (min-width: 1024px) {
    .btn.svg::before {
      background-color: var(--item-color);
      content: '';
      position: absolute;
      top: 80%;
      left: -110%;
      width: 130%;
      height: 130%;
      transform: rotate(45deg);
      transition: all 0.4s cubic-bezier(0.31, -0.1, 0.43, 1.59);
    }
  }

  .icon-size-2x3 {
    grid-row: span 2;
    grid-column: span 3;

    &:hover {
      .btn {
        .heisingfont {
          /* color: #F7931A; */
          /* color: #00758f; */
          transform: translate(-50%, -50%) scale(2.2);
        }

        &.svg::before {
          top: -34%;
          left: -13%;
        }
      }
    }

    .btn {
      .heisingfont {
        transform: translate(-50%, -50%) scale(1.76);
        line-height: calc(2 * var(--has-rows) + var(--has-row-gap));
        transition: all 0.4s cubic-bezier(0.31, -0.1, 0.43, 1.59);
      }

      &.svg::before {
        width: 130%;
        height: 174%;
      }
    }
  }

  .icon-size-1x2 {
    grid-row: span 1;
    grid-column: span 2;

    &:hover .btn.svg::before {
      top: -57%;
      left: -15%;
    }

    .btn.svg::before {
      position: absolute;
      height: 220%;
      transform: rotate(45deg);
    }
  }

  .icon-size-2x1 {
    grid-row: span 2;

    &:hover .btn.svg::before {
      top: 0%;
      left: -65%;
    }

    .btn {
      &.svg::before {
        top: 83%;
        left: -230%;
        width: 230%;
        height: 100%;
      }

      .heisingfont {
        line-height: calc(2 * var(--has-rows) + var(--has-row-gap));
      }
    }
  }

}

@media (max-width: 1800px) {
  .page {
    --has-columns: 10;
  }

  @media (max-width: 1300px) {
    .page {
      --has-columns: 9;
    }

    @media (max-width: 1200px) {
      .page {
        --has-columns: 8;
      }
    }
  }
}

@media (max-width: 1024px) {
  .nav-container {
    top: 223px;
  }

  .page {
    --has-column-gap: 15px;
    --has-columns: 6;
    /* --has-row-gap: 30px; */
  }

  .page .custom-nav:hover {
    transform: none;

    .btn {
      background-color: var(--item-color);
    }
  }
}

@media (max-width: 750px) {
  .page {
    --has-columns: 4;
  }

  @media (max-width: 600px) {
    .page {
      --has-columns-length: 90px;
      --has-rows: 90px;
      --has-row-gap: 25px;
    }

    
  }

  @media (max-width: 460px) {
    .page {
      /* --has-columns: 3; */
      /* --has-row-gap: 25px; */
      --has-columns-length: 70px;
      --has-rows: 70px;
    }

    .page .custom-nav::after {
      font-size: 10px;
    }

  }

  @media (max-width: 400px) {
    .page {
      --has-columns-length: 65px;
      --has-rows: 65px;
    }
  }

  @media (max-width: 370px) {
    .page {
      --has-columns: 3;
    }
  }
}

@keyframes shadeAnimation {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}
</style>
