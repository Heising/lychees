<script setup lang="ts">
import { onMounted } from 'vue'
import { egg } from './utils/tools'

// 取消全局右键菜单
document.addEventListener('contextmenu', function (event) {
  event.preventDefault()
})

const date = new Date()
const age = Math.floor((date.getTime() - new Date('2024-05-01').getTime()) / 864e5)
const csstext = `font-size:20px;color:#01aebc;font-family:'Microsoft Yahei',PingFangSC-Regular,Helvetica,sans-serif,'等线'; line-height: 20px;`
onMounted(() => {
  egg(
    `%c已经上线 %c${age}%c 天\n
不好，被你发现了！%cCopyright Redli.cn Page \xa9 2024-%s`,
    csstext,
    `${csstext}; line-height: 20px;color:#999999;`,
    csstext,
    'font-size:12px;color:#999999;vertical-align: text-top;line-height: 30px;',
    date.getFullYear()
  )
})
</script>

<template>
  <!-- 一级路由出口组件 -->
  <RouterView />
</template>
