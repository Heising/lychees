import path from 'path'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import { envParse } from 'vite-plugin-env-parse'
import { compression } from 'vite-plugin-compression2'
export default defineConfig({
  envDir: './env',
  build: {
    rollupOptions: {
      // 配置rollup的一些构建策略
      output: {
        // 控制输出
        //         // 在rollup里面, hash代表将你的文件名和文件内容进行组合计算得来的结果
        //         // assetFileNames: "[hash].[name].[ext]",
        //     // },
        // output: {
        chunkFileNames: 'static/js/[name].[hash].js',
        entryFileNames: 'static/js/[name].[hash].js',
        assetFileNames: 'static/[ext]/[name].[hash].[ext]',
        manualChunks(id: string) {
          // console.log('id是', id)
          // svg: ['svg'],
          // 分包
          // if (id.includes('radix-ui')) {
          //   return 'radix-ui'
          // }
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
    cssCodeSplit: true,
  },
  plugins: [
    react(),
    envParse({}),
    // 新增gzip压缩
    compression({ algorithm: 'brotliCompress', exclude: [/\.(br)$/, /\.(gz)$/] }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
})
