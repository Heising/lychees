import type { Item, RestaurantItem } from '@/models'
import axios from 'axios'
import type { AutocompleteFetchSuggestionsCallback, FormInstance } from 'element-plus'
import { reactive, type Ref } from 'vue'
// 获取unix秒级的时间戳
export function unix() {
  return Math.trunc(Date.now() / 1000)
}

export function getFilterSymbol(): RestaurantItem[] {
  const symbolElement = document.querySelectorAll('symbol')

  const filterSymbol = Array.from(symbolElement)
    .filter((item) => !item.id.startsWith('Heising-'))
    .map((item) => {
      return { value: item.id }
    })

  return filterSymbol
}
// 移除以前导入的symbol
export function removeSymbolElements() {
  const symbolElement = document.querySelectorAll('symbol')

  Array.from(symbolElement)
    .filter((item) => !item.id.startsWith('Heising-'))
    .forEach((item) => {
      // 移除父级元素
      if (item.parentNode) {
        ;(item.parentNode as SVGSVGElement).remove()
      }

      // console.log(item.parentNode)

      // 移除子元素
      item.remove()
    })
}

let restaurants: RestaurantItem[] = []

export const querySearch = (queryString: string, cb: AutocompleteFetchSuggestionsCallback) => {
  restaurants = getFilterSymbol()

  console.log('queryString', queryString)

  const results = queryString ? restaurants.filter(createFilter(queryString)) : restaurants
  // call callback function to return suggestions
  console.log('results', results)
  cb(results)
}

const createFilter = (queryString: string) => {
  return (restaurant: RestaurantItem) => {
    return restaurant.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0
  }
}
function isURL(str: string) {
  try {
    new URL(str)
    return true
  } catch (error) {
    return false
  }
}
const httpInstance = axios.create()
export const extractWebsiteInfo = async (currentEditingItem: Ref<Item>) => {

  if (!isURL(currentEditingItem.value.url)) return
  if (
    currentEditingItem.value.title &&
    !(!currentEditingItem.value.isSvg && !currentEditingItem.value.icon)
  ) {
    return
  }

  try {
    await httpInstance({
      url: currentEditingItem.value.url,
      method: 'GET',
    }).then((response) => {
      console.log('date是', response)

      // 拿标题
      if (!currentEditingItem.value.title) {
        // 创建正则表达式来匹配title属性值
        const regex = /<title[^>]*>([^<]+)<\/title>/

        // 使用match方法匹配正则表达式并提取href属性值
        const titleMatch = response.data.match(regex)

        // 输出匹配到的title属性值
        if (titleMatch) {
          currentEditingItem.value.title = titleMatch[1]
        } else {
          console.log('未找到标题')
        }
      }
      // 拿图标
      if (!currentEditingItem.value.isSvg && !currentEditingItem.value.icon) {
        // 创建正则表达式来匹配href属性值

        const hrefRegex = /<link rel="(?:shortcut |)icon".*?href="(.*?)"/

        // 使用match方法匹配正则表达式并提取href属性值
        const hrefMatch = response.data.match(hrefRegex)

        // 输出匹配到的href属性值
        if (hrefMatch) {
          let hrefValue
          if (hrefMatch[1].startsWith('//')) {
            hrefValue = hrefMatch[1]
          } else if (hrefMatch[1].startsWith('/')) {
            const url = new URL(currentEditingItem.value.url)
            const protocol = url.protocol // 获取协议（例如：http:）
            const hostname = url.hostname // 获取域名（例如：www.example.com）

            console.log(protocol) // 输出协议
            console.log(hostname) // 输出域名
            hrefValue = protocol + '//' + hostname + hrefMatch[1]
          } else {
            hrefValue = hrefMatch[1]
          }

          currentEditingItem.value.icon = hrefValue
        } else {
          console.log('未找到href属性值')
        }
      }
    })
  } catch (error) {
    console.log('可能同源策略拿不到', error)
  }
}

// interface RuleForm {
//   url: string
// }
export const bookmarksRules = reactive({
  url: [{ required: true, message: '请输入链接', trigger: 'blur' }],
})

// 重新验证表单
export const reValidate = (event: InputEvent, bookmarksRuleFormRef: FormInstance | null) => {
  bookmarksRuleFormRef?.validate()
}
// 给vue template拿不到window使用
export const openUrl = (url: string) => {
  window.open(url)
}
// 彩蛋打印
export const egg = function (...restParams: any[]) {
  queueMicrotask(console.log.bind(console, ...restParams))
}
