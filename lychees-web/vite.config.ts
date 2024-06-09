import { fileURLToPath, URL } from 'node:url'
import { resolve } from 'node:path'

import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { compression } from 'vite-plugin-compression2'

import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import { envParse } from 'vite-plugin-env-parse'
import { createHtmlPlugin } from 'vite-plugin-html'
import { VitePWA } from 'vite-plugin-pwa'
// https://vitejs.dev/config/
export default defineConfig({
  envDir: './env',
  build: {
    rollupOptions: {
      // 配置rollup的一些构建策略
      output: {
        // output: {
        chunkFileNames: 'static/js/[name].[hash].js',
        entryFileNames: 'static/js/[name].[hash].js',
        assetFileNames: 'static/[ext]/[name].[hash].[ext]',
        manualChunks(id: string) {
          // console.log('id是', id)
          if (id.includes('svg')) {
            // 分离svg.js
            return 'svg'
          }
          // svg: ['svg'],
          // 分包
          if (id.includes('element-plus')) {
            return 'element-plus'
          }
          if (id.includes('node_modules')) {
            return 'node_modules'
          }
        },
      },
    },
    // 小于指定阈值的静态资源以base64数据URL的形式内联到生成的文件中
    assetsInlineLimit: 0, // 4kb
    // outDir: "dists", // 配置输出目录
    // assetsDir: "static", // 配置输出目录中的静态资源目录
    // emptyOutDir: true, // 清除输出目录中的所有文件

    // 启用/禁用 CSS 代码拆分。当启用时，在异步 chunk 中导入的 CSS 将内联到异步 chunk 本身，并在其被加载时一并获取。
    // 如果禁用，整个项目中的所有 CSS 将被提取到一个 CSS 文件中。
    // 如果指定了 build.lib，build.cssCodeSplit 会默认为 false。
    cssCodeSplit: true,
  },
  plugins: [
    vue(),
    createSvgIconsPlugin({
      // 指定需要缓存的图标文件夹
      iconDirs: [resolve(process.cwd(), 'src/assets/svgs')],
      // 指定symbolId格式
      symbolId: 'Heising-[name]',

      /**
       * custom dom id
       * @default: __svg__icons__dom__
       */
      customDomId: 'svg_dom',
    }),
    createHtmlPlugin({
      minify: true,
    }),
    AutoImport({
      // Auto import functions from Element Plus, e.g. ElMessage, ElMessageBox... (with style)
      // 自动导入 Element Plus 相关函数，如：ElMessage, ElMessageBox... (带样式)
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [
        // Auto register Element Plus components
        // 自动导入 Element Plus 组件
        ElementPlusResolver(),
      ],
    }),
    envParse({}),
    VitePWA({
      manifest: {
        theme_color: '#01aebc',
      },
      registerType: 'autoUpdate',
    }),
    // 新增gzip压缩
    compression({ algorithm: 'brotliCompress', exclude: [/\.(br)$/, /\.(gz)$/] }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
