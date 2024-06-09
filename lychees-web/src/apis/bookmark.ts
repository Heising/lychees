// 封装所有和用户相关的接口函数
import request from '@/utils/http'
import type {
  AddItem,
  Bookmarks,
  IconfontLink,
  Item,
  PersonalInfo,
  SwapItem,
  UpdateResponse,
} from '@/models'
import { useNavStore } from '@/stores/navStore'

// 查找用户的书签信息
export const bookmarksFindApi = (updateAt: number) => {
  return request<Bookmarks>({
    url: '/auth/bookmarks/find/' + updateAt,
  })
}
// 删除出一个书签
export const bookmarksDeleteApi = (rowIndex: number, colIndex: number) => {
  return request<any, UpdateResponse>({
    url: `/auth/bookmarks/${rowIndex}/${colIndex}`,
    method: 'DELETE',
  })
}
// 增加一个标签
export const bookmarksPushApi = (rowIndex: number, item: Item) => {
  return request<any, AddItem>({
    url: `/auth/bookmarks/push/${rowIndex}`,
    method: 'POST',
    data: item,
  })
}

// 删除页面
export const DeletePageApi = (rowIndexs: number[]) => {
  return request<any, UpdateResponse>({
    url: `/auth/bookmarks/page`,
    method: 'DELETE',
    data: rowIndexs,
  })
}

// 交换元素
export const bookmarksSwapApi = (dragFrom: SwapItem, dragTo: SwapItem) => {
  const navStore = useNavStore()

  return request<any, UpdateResponse>({
    url: `/auth/bookmarks/swap/${dragFrom.rowIndex}`,
    method: 'PUT',
    data: {
      [dragFrom.colIndex]: navStore.bookmarks.arrayBookmarks[dragFrom.rowIndex][dragFrom.colIndex],
      [dragTo.colIndex]: navStore.bookmarks.arrayBookmarks[dragTo.rowIndex][dragTo.colIndex],
    },
  })
}

// 移动元素
export const bookmarksMoveApi = (from: SwapItem, to: SwapItem, item: Item) => {
  return request<any, AddItem>({
    url: `/auth/bookmarks/move/${from.rowIndex}/${from.colIndex}`,
    method: 'PUT',
    data: {
      newRowIndex: to.rowIndex,
      newColIndex: to.colIndex,
      item,
    },
  })
}

export const bookmarksUpdateApi = (from: SwapItem, item: Item) => {
  return request<any, AddItem>({
    url: `/auth/bookmarks/update/${from.rowIndex}/${from.colIndex}`,
    method: 'PUT',
    data: item,
  })
}
// 更新第三方图标链接
export const updateIconfontApi = (iconfontlink: IconfontLink) => {
  return request<any, UpdateResponse>({
    url: `/auth/iconfont`,
    method: 'PUT',
    data: iconfontlink,
  })
}
export const updatePersonInfoApi = (personalInfo: PersonalInfo) => {
  return request<any, UpdateResponse>({
    url: `/auth/personalinfo`,
    method: 'PUT',
    data: personalInfo,
  })
}
