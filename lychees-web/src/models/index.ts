// 书签操作后的的响应
export interface UpdateResponse {
  updateAt?: number
  error?: string
}
// 拖拽事件的序号
export interface SwapItem {
  rowIndex: number
  colIndex: number
}
export interface Item {
  isSvg: boolean
  icon: string
  turn: boolean
  color?: string
  // backgroundColor?: string

  title: string
  url: string
  iconSize?: string
}

export interface AddItem extends UpdateResponse {
  data: Item
}
export interface IconfontLink {
  iconfontLink: string
}
export interface PersonalInfo {
  nickname: string
  birthday?: Date
}

export interface Bookmarks extends PersonalInfo, IconfontLink {
  updateAt: number
  arrayBookmarks: Item[][]
}

export interface User {
  id: number
  email: string

  nickname: string
  createdAt: number
}
// 登录后的响应
export interface LoginResponse {
  userInfo: User
  token: string
}
// 登录时的数据定义
export interface LoginUser {
  email: string
  encrypted: string
  nanoid?: string
}
export interface LoginDeviceType {
  expireUnix: number
  loginTime: number
  token: string
  userAgent: string
  clientIp: string
}
export interface LoginDevicesType {
  devices: LoginDeviceType[]
}
// 注册时的数据定义
export interface RegisterUser extends LoginUser {
  nickname: string
}
// 注册后的响应
export interface RegisterResponse {
  message: string
  userId: number
  bookmarksId: number
}
export interface TokenInfo {
  token: string
  expireUnix: number
}
export interface Payload {
  exp: number
  iat: number
  id: number
  iss: string
  sub: string
  tokenInfo: TokenInfo
}

export interface PublicKey {
  publicKey: string
  nanoid: string
}

// svg图标建议
export interface RestaurantItem {
  value: string
}
