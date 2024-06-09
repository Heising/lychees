<script setup lang="ts">
import SearchContainer from '@/components/SearchContainer.vue'
import NavContainer from '@/components/NavContainer.vue'
import { useFocusStore } from '@/stores/searchStore'
import { onMounted, onUnmounted, provide, ref } from 'vue'
import { getdbWallpaper } from '@/stores/db'
import { clipboardKey, contextMenuStyle, getWallpaperKey } from '@/models/key'
const focusStore = useFocusStore()
const displayImg = ref(false)
const imageUrl = ref('')
const isNull = ref(false)
import { useNavStore } from '@/stores/navStore'
let currentInputElement: HTMLInputElement
// 剪切板菜单
const contextClipboard = ref<HTMLUListElement | null>(null)
const navStore = useNavStore()
// let img = new Image()
// function imgOnload() {
//   var canvas = document.createElement('canvas')
//   var ctx = canvas.getContext('2d') as CanvasRenderingContext2D
//   canvas.width = img.width
//   canvas.height = img.height
//   ctx.drawImage(img, 0, 0)

//   var imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
//   var totalPixels = imageData.data.length / 4 // 每个像素有4个值（RGBA通道）
//   var whitePixels = 0

//   // 遍历每个像素，检查是否为白色
//   for (var i = 0; i < imageData.data.length; i += 4) {
//     var red = imageData.data[i]
//     var green = imageData.data[i + 1]
//     var blue = imageData.data[i + 2]

//     // 判断是否为白色像素（可以根据实际情况调整阈值）
//     if (red > 200 && green > 200 && blue > 200) {
//       whitePixels++
//     }
//   }

//   // 计算白色像素比例
//   var whiteRatio = whitePixels / totalPixels
//   console.log('白色部分比例: ' + whiteRatio)
// }

const getWallpaper = async () => {
  getdbWallpaper(imageUrl)
}

provide(getWallpaperKey, getWallpaper)

const clipboard = async (event: MouseEvent) => {
  const inputElement = event.currentTarget as HTMLInputElement
  // console.log(event.currentTarget)
  // console.log('inputElement.value是：', inputElement.value)

  currentInputElement = inputElement
  if (!inputElement.value) {
    isNull.value = true
  } else {
    isNull.value = false
  }
  contextClipboardHide()

  let x = event.clientX // 获取鼠标的水平坐标
  let y = event.clientY // 获取鼠标的垂直坐标
  let distanceX = window.innerWidth - x // 鼠标与页面右边的距离
  let distanceY = window.innerHeight - y // 鼠标与页面底部的距离

  if (contextClipboard.value !== null) {
    // 先让其可见才有高度宽度
    contextClipboard.value.style.display = 'block'
    // 返回元素内容的总高度，包括溢出部分 包括不可见部分

    if (distanceX < contextClipboard.value.offsetWidth) {
      x -= contextClipboard.value.offsetWidth - distanceX
    }
    if (distanceY < contextClipboard.value.offsetHeight) {
      y -= contextClipboard.value.offsetHeight - distanceY
    }

    // 将contextMenuStyle对象中的样式应用于DOM元素
    Object.assign(contextClipboard.value.style, contextMenuStyle)

    contextClipboard.value.style.top = y.toString() + 'px'
    contextClipboard.value.style.left = x.toString() + 'px'
  }
}
const handleSearch = () => {
  // console.log(currentInputElement)

  // console.log('是搜索框')
  if (currentInputElement) {
    if (currentInputElement.type == 'search') {
      // 创建一个键盘事件对象
      var event = new KeyboardEvent('keyup', {
        key: 'Enter',
        code: 'Enter',
        keyCode: 13,
        which: 13,
        charCode: 13,
        bubbles: true,
      })

      // 触发事件
      currentInputElement.dispatchEvent(event)
    }
  }

  contextClipboardHide()
}
// 创建一个input事件
const inputDispatchEvent = (currentInputElement: HTMLInputElement) => {
  const inputEvent = new Event('input', {
    bubbles: true,
    cancelable: true,
  })
  // 触发input事件
  currentInputElement.dispatchEvent(inputEvent)
}
const handleCut = () => {
  navigator.clipboard.writeText(currentInputElement.value).then(
    function () {
      /* clipboard successfully set */
      if (currentInputElement) {
        currentInputElement.value = ''
        inputDispatchEvent(currentInputElement)
      }
    },
    function () {
      /* clipboard write failed */
    }
  )

  contextClipboardHide()
}
const handleCopy = () => {
  navigator.clipboard.writeText(currentInputElement.value).then(
    function () {
      /* clipboard successfully set */
    },
    function () {
      /* clipboard write failed */
    }
  )
  contextClipboardHide()
}
const handlePaste = () => {
  navigator.clipboard.readText().then((clipText) => {
    if (currentInputElement) {
      currentInputElement.value += clipText
      currentInputElement.focus()
      inputDispatchEvent(currentInputElement)
    }
  })
  contextClipboardHide()
}
// 隐藏右键菜单
const contextClipboardHide = () => {
  if (contextClipboard.value) {
    contextClipboard.value.removeAttribute('style')
  }
}
provide(clipboardKey, clipboard)
onMounted(() => {
  getWallpaper()
})
onUnmounted(() => {
  document.removeEventListener('click', handleClick)
})
document.addEventListener('click', handleClick)
function handleClick(event: Event) {
  // 检查点击的元素是否是要隐藏的元素
  // 检查点击的元素是否是要隐藏的元素或其子元素

  if (
    // event.target !== contextMenuTrigger.value &&
    contextClipboard.value !== null &&
    !contextClipboard.value.contains(event.target as Node)
  ) {
    // 如果不是要隐藏的元素，则隐藏它
    // contextMenuTrigger.value.style.display = 'none'
    contextClipboard.value.removeAttribute('style')
  }
}
// 预加载图片
async function preloadImage() {
  // var img = new Image()
  // img.src = wallpaper
  // imgOnload()
  // await img.decode() // 等待图片加载完成
  // console.log('图片已预加载')
  displayImg.value = true
}

// onMounted(async () => {
//   await preloadImage()
// })
</script>

<template>
  <div class="container" @click="focusStore.handleBlur">
    <div :class="{ mask: navStore.isMask }" v-show="focusStore.isFog"></div>
    <img
      v-show="displayImg && imageUrl"
      @load="preloadImage"
      :src="imageUrl"
      :class="{ focus: focusStore.isFog }"
    />
  </div>

  <div :class="{ mask: navStore.isMask }" v-show="!focusStore.isFog"></div>
  <SearchContainer></SearchContainer>
  <NavContainer></NavContainer>

  <ul class="menu clipboardMenu" ref="contextClipboard">
    <li
      class="menuItem"
      :class="{ hide: currentInputElement?.type !== 'search' }"
      @click="handleSearch"
    >
      <svg class="clipboardIcon icon" aria-hidden="true">
        <use xlink:href="#Heising-magnify"></use>
      </svg>
      搜索
    </li>
    <li class="menuItem" :class="{ hide: isNull }" @click="handleCut">
      <svg class="clipboardIcon icon" aria-hidden="true">
        <use xlink:href="#Heising-scissors"></use>
      </svg>
      剪切
    </li>
    <li class="menuItem" :class="{ hide: isNull }" @click="handleCopy">
      <svg class="clipboardIcon icon" aria-hidden="true">
        <use xlink:href="#Heising-clipboard"></use>
      </svg>
      复制
    </li>
    <li class="menuItem" @click="handlePaste">
      <svg class="clipboardIcon icon" aria-hidden="true">
        <use xlink:href="#Heising-copy"></use>
      </svg>
      粘贴
    </li>
  </ul>
</template>

<style scoped lang="scss">
.container {
  z-index: -999;
  position: fixed;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  img {
    z-index: -999;
    position: fixed;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    // transition: all 0.25s;
    transition: all 0.6s ease-in-out;

    &.focus {
      filter: blur(10px);
      transform: scale(1.2);
    }
  }
}

.mask {
  width: 100vw;
  height: 100vh;
  // z-index: -1;
  background: rgba(0, 0, 0, 0.5);
  // transition: all 0.6s ease-in-out;

  // backdrop-filter: blur(5px);
}

.clipboardMenu {
  z-index: 2023;
  background-color: rgba(233, 233, 233, 0.6);
  backdrop-filter: blur(30px);
  .hide {
    display: none;
  }
  .clipboardIcon {
    margin-right: 10px;
  }
}
</style>
