import Dexie, { type Table } from 'dexie'
import type { Ref } from 'vue'
import { useUserStore } from '@/stores/userStore'

export interface searchSsuggest {
  id?: number
  keyword: string
  suggestList: string[]
  time: number
}

export class MySearchSsuggests extends Dexie {
  // 'friends' is added by dexie when declaring the stores()
  // We just tell the typing system this is the case
  searchSsuggestList!: Table<searchSsuggest>

  constructor() {
    super('searchSsuggests')
    this.version(1).stores({
      searchSsuggestList: '++id, keyword, suggestList, time', // Primary key and indexed props
    })
  }
}

export const db = new MySearchSsuggests()

// 使用.clear()方法清空数据库中的所有对象
// async function clearDatabase() {
//   try {
//     await db.searchSsuggestList.clear();
//     console.log("数据库中的所有对象已成功清空");
//   } catch (error) {
//     console.error("清空数据库对象时出错：", error);
//   }
// }

// 使用.delete()方法删除整个数据库
export async function cleanData() {
  try {
    // await dbWallpaper.dbWallpaperList.clear()
    await db.searchSsuggestList.clear()
    console.log('数据库成功清空')
  } catch (error) {
    console.error('清空数据库时出错：', error)
  }
}

export async function cleanWallpaperData() {
  try {
    await dbWallpaper.dbWallpaperList.clear()
    console.log('壁纸成功清空')
  } catch (error) {
    console.error('清空壁纸出错：', error)
  }
}

// export async function deleteDatabase() {
//   try {
//     await Dexie.delete('searchSsuggests')
//     await Dexie.delete('wallpaper')
//     console.log('数据库成功删除')
//   } catch (error) {
//     console.error('删除数据库时出错：', error)
//   }
// }
export interface Wallpaper {
  id?: number
  wallpaper: Blob
}

export class MyWallpaper extends Dexie {
  // 'friends' is added by dexie when declaring the stores()
  // We just tell the typing system this is the case
  dbWallpaperList!: Table<Wallpaper>

  constructor() {
    super('wallpaper')
    this.version(1).stores({
      dbWallpaperList: '++id, wallpaper', // Primary key and indexed props
    })
  }
}

export const dbWallpaper = new MyWallpaper()

// 使用.clear()方法清空数据库中的所有对象
// export async function clearDatabase() {
//   try {
//     await dbWallpaper.dbWallpaperList.clear()
//     console.log('壁纸成功删除')
//   } catch (error) {
//     console.error('删除壁纸出错：', error)
//   }

//   console.log('壁纸成功删除')
// }
// 查找壁纸
export const getdbWallpaper = async (imageUrl: Ref<string>) => {
  const userStore = useUserStore()
  await dbWallpaper.dbWallpaperList
    .get(userStore.userInfo.userInfo.id, function (image) {
      if (image) {
        // 如果找到图片数据
        imageUrl.value = URL.createObjectURL(image.wallpaper) // 创建一个URL对象
      } else {
        console.log('Image not found in the database.')
        imageUrl.value = ''
      }
    })
    .catch(function (error) {
      console.log('Error opening database: ' + error)
      imageUrl.value = ''
    })
}
