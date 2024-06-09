<template>
  <div>
    <div class="corner-container" ref="cornerContainer">
      <svg class="icon person" aria-hidden="true" @click="showMenu(personInfoMenu)">
        <use :xlink:href="'#Heising-person'"></use>
      </svg>
      <svg class="icon settings" aria-hidden="true" @click="showMenu(settingsMenu)">
        <use :xlink:href="'#Heising-gear'"></use>
      </svg>
    </div>
    <ul class="menu" ref="personInfoMenu">
      <li class="menuItem" @click="openEditPersonalInfo">编辑用户信息</li>
      <li class="menuItem" @click="openDevices">查看登录设备</li>
      <el-popconfirm
        confirm-button-text="确认"
        cancel-button-text="取消"
        title="确认退出登录？"
        @confirm="confirmLogOut"
        :hide-icon="true"
        confirm-button-type="warning"
        :hide-after="0"
      >
        <template #reference>
          <li class="menuItem warn">退出登录状态</li>
        </template>
      </el-popconfirm>
    </ul>

    <ul class="menu" ref="settingsMenu">
      <li class="menuItem" @click="openAdd">添加一个书签</li>
      <el-divider />
      <li class="menuItem" @click="openEditWallpaper">更改本地壁纸</li>
      <li class="menuItem" @click="openupdateIconfontLink">更新图标资源</li>
      <el-divider />
      <li class="menuItem warn" @click="removeEmptyPage">移除空的页面</li>
    </ul>

    <el-dialog v-model="devicesVisible" title="登录设备信息">
      <ul class="device-list">
        <li v-for="item in loginDevices" :key="item.token" class="device-info">
          <div class="device-header">
            {{ item.clientIp ? item.clientIp : '未知IP' }}

            <span>{{ unixTimestampConvert(item.loginTime) }}</span>
            <span v-if="isCurrentDevice(item.token)">当前设备</span>

            <el-popconfirm
              confirm-button-text="确认"
              cancel-button-text="取消"
              title="确认移除该设备？"
              @confirm="expelDevice(item.token)"
              :hide-icon="true"
              confirm-button-type="warning"
              :hide-after="0"
              v-else
            >
              <template #reference>
                <svg class="icon" aria-hidden="true">
                  <use xlink:href="#Heising-cross"></use>
                </svg>
              </template>
            </el-popconfirm>
          </div>
          <div>
            <svg class="icon" aria-hidden="true">
              <use :xlink:href="`#Heising-${isMobile(item.userAgent)}`"></use>
            </svg>
            {{ item.userAgent ? item.userAgent : '未知设备' }}
          </div>
        </li>
      </ul>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="devicesVisible = false" round>取消</el-button>
          <el-popconfirm
            confirm-button-text="确认"
            cancel-button-text="取消"
            title="确认驱逐全部登录设备？"
            @confirm="expelDevices()"
            :hide-icon="true"
            confirm-button-type="warning"
            :hide-after="0"
          >
            <template #reference>
              <el-button type="danger" round>驱逐全部</el-button>
            </template>
          </el-popconfirm>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="updateIconfontLinkVisible" title="更新第三方图标链接">
      <el-form :model="currentEditingIconfontLink">
        <el-form-item label="请输入链接" label-width="100px">
          <el-input
            v-model.trim="currentEditingIconfontLink.iconfontLink"
            autocomplete="off"
            clearable
            @contextmenu.prevent="clipboardEvent"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="updateIconfontLinkVisible = false">取消</el-button>
          <el-button type="primary" @click="updateIconfontLink">确认</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="updatePersonVisible" title="更改用户信息">
      <el-form :model="currentEditingPersonalInfo">
        <el-form-item label="昵称" label-width="100px">
          <el-input
            v-model.trim="currentEditingPersonalInfo.nickname"
            autocomplete="off"
            clearable
          />
        </el-form-item>
      </el-form>
      <el-form :model="currentEditingPersonalInfo">
        <el-form-item label="生日" label-width="100px">
          <el-date-picker
            v-model="currentEditingPersonalInfo.birthday"
            type="date"
            placeholder="选择你的生日"
            :size="size"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="openUrl('https://account.redli.cn/')">
            账号高级管理
          </el-button>

          <el-button @click="updatePersonVisible = false">取消</el-button>
          <el-button type="primary" @click="updatePersonalInfo">确认</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog
      class="wallpaper-dialog"
      v-model="updateWallpaperVisible"
      title="更改本地壁纸"
      center
    >
      <div class="preview-container">
        <div :style="{ backgroundImage: `url(${imageUrl})` }" class="preview-image"></div>
        <div class="wallpaper-btn-container">
          <el-button type="danger" @click="cleanWallpaper" plain>移除壁纸</el-button>

          <el-upload
            ref="uploadWallpaper"
            :limit="1"
            :on-exceed="handleExceed"
            action="#"
            :show-file-list="false"
            :auto-upload="false"
            :on-change="handleChange"
            accept="image/*"
          >
            <el-button type="success" plain>上传壁纸</el-button>
          </el-upload>
          <div>
            <el-button type="warning" @click="toggleCover" plain>
              {{ navStore.isMask ? '关闭' : '启用' }}遮罩
            </el-button>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="updateWallpaperVisible = false">取消操作</el-button>
          <el-button type="primary" @click="updateWallpaper">确认更换</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="addBookmarkFormVisible" title="添加一个书签">
      <el-form :model="currentEditingItem" :rules="bookmarksRules" ref="bookmarksRuleFormRef">
        <el-form-item label="请输入链接" label-width="100px" prop="url">
          <el-input
            v-model.trim="currentEditingItem.url"
            autocomplete="off"
            clearable
            @contextmenu.prevent="clipboardEvent"
            @input="reValidate($event, bookmarksRuleFormRef)"
          />
        </el-form-item>

        <el-form-item label="请输入标题" label-width="100px">
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

        <el-form-item label="请输入颜色" label-width="100px">
          <el-color-picker
            v-model="currentEditingItem.color"
            :disabled="!currentEditingItem.isSvg"
          />
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

        <el-form-item label="添加到页面" label-width="100px">
          <el-select v-model="currentRowIndex" placeholder="选择你的页面">
            <el-option
              v-for="(_, index) in navStore.bookmarks.arrayBookmarks"
              :label="index + 1"
              :value="index"
              :key="index"
            />
            <el-option
              :label="navStore.bookmarks.arrayBookmarks.length + 1 + ': 新增一页'"
              :value="navStore.bookmarks.arrayBookmarks.length"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="addBookmarkFormVisible = false">取消</el-button>
          <el-button type="primary" @click="bookmarksPush">确认</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { inject, onMounted, onUnmounted, ref } from 'vue'
import { useUserStore } from '@/stores/userStore'
import { useNavStore } from '@/stores/navStore'
import type { FormInstance, UploadInstance, UploadProps, UploadRawFile } from 'element-plus'
import { ElMessage, genFileId } from 'element-plus'
import './menu.scss'
import type { PersonalInfo, Item, IconfontLink, LoginDeviceType, Payload } from '@/models'
import { extractWebsiteInfo, bookmarksRules, reValidate } from '@/utils/tools'
import {
  DeletePageApi,
  bookmarksPushApi,
  updateIconfontApi,
  updatePersonInfoApi,
} from '@/apis/bookmark'
import { ElNotification } from 'element-plus'
import { querySearch } from '@/utils/tools'
import { expelDeviceApi, expelDevicesApi, getDevicesApi, logOutApi } from '@/apis/user'
import router from '@/router'
import { dbWallpaper, getdbWallpaper, cleanWallpaperData } from '@/stores/db'
import { clipboardKey, contextMenuStyle, getWallpaperKey } from '@/models/key'
import { openUrl } from '@/utils/tools'
const bookmarksRuleFormRef = ref<FormInstance | null>(null)
const userStore = useUserStore()
const emit = defineEmits<{
  changeScroll: [index: number] // 具名元组语法
  getBookmarks: []
}>()
const clipboardEvent = inject(clipboardKey, (enevt: MouseEvent) => {
  console.log('上层组件找不到clipboardKey！')
})
const addBookmarkFormVisible = ref(false)
const updateIconfontLinkVisible = ref(false)
const updatePersonVisible = ref(false)
const devicesVisible = ref(false)
const updateWallpaperVisible = ref(false)
const loginDevices = ref<LoginDeviceType[]>([])
// const isEditing = ref<boolean>(false)
const currentEditingItem = ref<Item>({
  isSvg: false,
  icon: '',
  title: '',
  url: '',
  turn: false,
})
const currentEditingIconfontLink = ref<IconfontLink>({
  iconfontLink: '',
})
const size = ref<'default' | 'large' | 'small'>('default')
const currentEditingPersonalInfo = ref<PersonalInfo>({
  nickname: '',
})
const currentRowIndex = ref(0)

const navStore = useNavStore()
// 右键菜单逻辑
const settingsMenu = ref<HTMLUListElement | null>(null)
const personInfoMenu = ref<HTMLUListElement | null>(null)

const openEditPersonalInfo = () => {
  currentEditingPersonalInfo.value.nickname = navStore.bookmarks.nickname
  currentEditingPersonalInfo.value.birthday = navStore.bookmarks.birthday
  personInfoMenuHide()
  updatePersonVisible.value = true
}
const isMobile = (userAgent: string) => {
  if (userAgent.includes('Mobile')) {
    return 'mobile'
  } else {
    return 'desktop'
  }
}

const unixTimestampConvert = (unixTimestamp: number) => {
  // 创建一个新的JavaScript Date对象，基于Unix时间戳（乘以1000转换为毫秒）
  const date = new Date(unixTimestamp * 1000)

  // 生成日期字符串
  // console.log(date.toLocaleDateString('zh-CN')) // 输出: 2022/5/6

  // 生成时间字符串
  // console.log(date.toLocaleTimeString('zh-CN')) // 输出: 13:10:34

  return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN')
}

const openDevices = () => {
  personInfoMenuHide()
  getDevices().then(() => {
    devicesVisible.value = true
  })
}
const getDevices = async () => {
  const result = await getDevicesApi()
  loginDevices.value = result.data.devices
}
const expelDevice = async (token: string) => {
  await expelDeviceApi(token).then(() => {
    getDevices()
  })
}

const expelDevices = async () => {
  await expelDevicesApi().finally(() => {
    devicesVisible.value = false

    userStore.clearUserInfo()
  })
  router.push('/login')
}
const isCurrentDevice = (token: string) => {
  const currentJwt = userStore.userInfo.token
  if (currentJwt) {
    const jwtArray = currentJwt.split('.')
    // 拿到负载段，并且base64转json，再json解码
    const payload: Payload = JSON.parse(atob(jwtArray[1]))

    // 判断是否相等
    if (payload.tokenInfo.token === token) {
      return true
    }
    return false
  }
  return false
}
const uploadWallpaper = ref<UploadInstance>()
// 限制一个
const handleExceed: UploadProps['onExceed'] = (files) => {
  uploadWallpaper.value!.clearFiles()
  const file = files[0] as UploadRawFile
  file.uid = genFileId()
  uploadWallpaper.value!.handleStart(file)
}

const imageUrl = ref('')
const getWallpaper = inject(getWallpaperKey, () => {
  console.log('上层组件找不到getWallpaperKey！')
})
let imageRaw: Blob
// 清空本地壁纸
async function cleanWallpaper() {
  await cleanWallpaperData()
  updateWallpaperVisible.value = false
  navStore.isMask = true
  getWallpaper && getWallpaper()
}
// 打开编辑本地壁纸
const openEditWallpaper = async () => {
  await getdbWallpaper(imageUrl)
  srcList[0] = imageUrl.value
  settingsMenuHide()
  updateWallpaperVisible.value = true
}

// 更新壁纸
const updateWallpaper = () => {
  updateWallpaperVisible.value = false
  if (imageRaw) {
    dbWallpaper.dbWallpaperList
      .put({ id: userStore.userInfo.userInfo.id, wallpaper: imageRaw })
      .then(function () {
        console.log('Image saved successfully!')
      })
      .catch(function (error) {
        console.error('Error saving image: ' + error)
      })
    // 调用爷爷级组件
    getWallpaper && getWallpaper()
  }
}

const toggleCover = () => {
  navStore.isMask = !navStore.isMask
}
const srcList = ['']
const handleChange: UploadProps['onChange'] = (uploadFile, uploadFiles) => {
  // console.log(uploadFile.raw)
  if (uploadFile.raw && !uploadFile.raw.type.startsWith('image/')) {
    ElMessage.error('不是图片格式!')
    return
  } else if (uploadFile.raw?.size! / 20480 / 1024 > 2) {
    ElMessage.error('图片大小不能超过20MiB!')
    return
  }
  const imgUrl = URL.createObjectURL(uploadFile.raw!)
  srcList[0] = imgUrl
  imageUrl.value = imgUrl
  imageRaw = uploadFile.raw!
}

// 打开添加书签
const openAdd = () => {
  currentEditingItem.value = {
    isSvg: false,
    icon: '',
    title: '',
    url: '',
    turn: false,
  }
  settingsMenuHide()
  bookmarksRuleFormRef.value?.clearValidate()
  addBookmarkFormVisible.value = true
}
const openupdateIconfontLink = () => {
  currentEditingIconfontLink.value.iconfontLink = navStore.bookmarks.iconfontLink
  settingsMenuHide()
  updateIconfontLinkVisible.value = true
}



const cornerContainer = ref<HTMLDivElement | null>(null)
// 隐藏用户菜单
const personInfoMenuHide = () => {
  if (personInfoMenu.value) {
    personInfoMenu.value.removeAttribute('style')
  }
}

// 隐藏设置菜单
const settingsMenuHide = () => {
  if (settingsMenu.value) {
    settingsMenu.value.removeAttribute('style')
  }
}
onUnmounted(() => {
  document.removeEventListener('click', handleClick)
})
onMounted(() => {
  document.addEventListener('click', handleClick)
})

// 添加点击事件监听器到document
function handleClick(event: Event) {
  // 检查点击的元素是否是要隐藏的元素
  // 检查点击的元素是否是要隐藏的元素或其子元素

  if (
    // event.target !== clickMenu.value &&
    cornerContainer.value !== null &&
    !cornerContainer.value.parentNode?.contains(event.target as Node)
  ) {
    // 如果不是要隐藏的元素，则隐藏它
    // contextMenuTrigger.value.style.display = 'none'
    settingsMenu.value?.removeAttribute('style')
    personInfoMenu.value?.removeAttribute('style')
  }
}



// 打开菜单
const showMenu = (ulElement: HTMLUListElement | null) => {
  personInfoMenuHide()
  settingsMenuHide()

  if (cornerContainer.value !== null) {
    let x = cornerContainer.value.getBoundingClientRect().right // 获取元素的right
    let y = cornerContainer.value.getBoundingClientRect().bottom // 获取元素的bottom最底部

    if (ulElement !== null) {
      // 先让其可见才有高度宽度
      ulElement.style.display = 'block'

      x -= ulElement.offsetWidth
      y += 10
      Object.assign(ulElement.style, contextMenuStyle)

      ulElement.style.top = y.toString() + 'px'
      ulElement.style.left = x.toString() + 'px'
    }
  }
}
const updatePersonalInfo = async () => {
  updatePersonVisible.value = false

  updatePersonInfoApi(currentEditingPersonalInfo.value)
    .then((response) => {
      response.updateAt && (navStore.bookmarks.updateAt = response.updateAt)
      ElNotification({
        title: '更新用户信息成功',
        type: 'success',
      })
      navStore.bookmarks.nickname = currentEditingPersonalInfo.value.nickname
      navStore.bookmarks.birthday = currentEditingPersonalInfo.value.birthday
    })
    .catch(() => {
      ElNotification({
        title: '更新用户信息失败',
        type: 'error',
      })
    })
}

const updateIconfontLink = async () => {
  updateIconfontLinkVisible.value = false
  if (currentEditingIconfontLink.value.iconfontLink == navStore.bookmarks.iconfontLink) {
    return
  }

  updateIconfontApi(currentEditingIconfontLink.value)
    .then((response) => {
      ElNotification({
        title: '更新第三方图标链接成功',
        type: 'success',
      })
      emit('getBookmarks')
    })
    .catch(() => {
      ElNotification({
        title: '更新第三方图标链接失败',
        type: 'error',
      })
    })
}
const bookmarksPush = async () => {
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
  addBookmarkFormVisible.value = false

  extractWebsiteInfo(currentEditingItem)
  
  bookmarksPushApi(currentRowIndex.value, currentEditingItem.value)
    .then((response) => {
      response.updateAt && (navStore.bookmarks.updateAt = response.updateAt)
      ElNotification({
        title: '添加成功',
        type: 'success',
      })
      if (response.updateAt) {
        navStore.bookmarks.updateAt = response.updateAt
      }

      if (navStore.bookmarks.arrayBookmarks.length <= currentRowIndex.value) {
        navStore.bookmarks.arrayBookmarks.push([
          {
            ...response.data,
          },
        ])
      } else {
        navStore.bookmarks.arrayBookmarks[currentRowIndex.value].push({
          ...response.data,
        })
      }

      // console.log(response)
      emit('changeScroll', currentRowIndex.value)
    })
    .catch(() => {
      ElNotification({
        title: '添加失败',
        type: 'error',
      })
    })
}

const removeEmptyPage = () => {
  settingsMenuHide()
  const emptypage: number[] = []
  navStore.bookmarks.arrayBookmarks.forEach((item, index) => {
    if (item === null || item.length === 0) {
      emptypage.push(index)
    }
  })

  if (emptypage.length > 0) {
    DeletePageApi(emptypage)
      .then((response) => {
        response.updateAt && (navStore.bookmarks.updateAt = response.updateAt)
        ElNotification({
          title: '删除成功',
          type: 'success',
        })
        if (response.updateAt) {
          navStore.bookmarks.updateAt = response.updateAt
        }
        navStore.bookmarks.arrayBookmarks = navStore.bookmarks.arrayBookmarks.filter((item) => {
          if (item === null || item.length === 0) {
            return false
          }
          return item.length > 0
        })
      })
      .catch(() => {
        ElNotification({
          title: '删除失败',
          type: 'error',
        })
      })
  }
  emit('changeScroll', 0)
}

const confirmLogOut = async () => {
  await logOutApi().finally(() => {
    userStore.clearUserInfo()
  })
  personInfoMenuHide()
  router.push('/login')
}
</script>
<style lang="scss">
.device-list {
  .device-info {
    border-bottom: thin solid var(--theme-color-grey);
    margin-bottom: 5px;
    padding-bottom: 5px;
    &:last-child {
      border-bottom-style: none;
    }
    .icon {
      vertical-align: baseline;
    }
    .device-header {
      position: relative;
      font-weight: bold;
      span:first-child {
        position: absolute;
        right: 80px;
      }
      span:last-child {
        float: right;
        .icon {
          position: absolute;
          top: 50%;
          right: 5px;
          transform: translateY(-50%);
        }
      }
    }
  }
}
.corner-container {
  position: fixed;
  top: 30px;
  right: 2.6042vw;

  .person,
  .settings {
    color: var(--heising-background-color);
    cursor: pointer;
    transition: all 0.25s;
    &:hover {
      color: #fff;
      transform: scale(1.2);
    }
  }
  .person {
    // right: calc(2.6042vw + 30px);
    margin-right: 14px;
  }
  .settings {
    right: 2.6042vw;
    &:active {
      transform: rotate(80deg) scale(1.2);
    }
  }
}

.wallpaper-dialog {
  .preview-container {
    display: flex;
    justify-content: space-evenly;
    .preview-image {
      background-position: center;
      background-color: #8c939d;
      width: 288px;
      height: 162px;
      background-size: cover;
    }

    .wallpaper-btn-container {
      display: flex;
      // gap: 25px;
      justify-content: space-between;
      flex-direction: column;
    }
  }
}
@media (max-width: 750px) {
  .wallpaper-dialog {
    .preview-container {
      flex-direction: column;
      align-items: center;
      .preview-image {
        // width: 150px;

        // height: 266.8px;
        width: 200px;

        height: 355.7px;
        margin-bottom: 20px;
      }
      .wallpaper-btn-container {
        width: 100%;

        flex-direction: row;
      }
    }
  }
}
</style>