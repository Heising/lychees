import type { InjectionKey } from 'vue'

export const getWallpaperKey = Symbol() as InjectionKey<() => Promise<void>>
export const clipboardKey = Symbol() as InjectionKey<(event: MouseEvent) => Promise<void>>
export const contextMenuStyle = {
  opacity: 1,
  transform: 'scale(1)',
}
