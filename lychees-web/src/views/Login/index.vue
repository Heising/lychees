<template>
  <svg aria-hidden="true" style="position: absolute; width: 0px; height: 0px; overflow: hidden">
    <symbol id="email" viewBox="0 0 15 15">
      <path
        d="M1 2C0.447715 2 0 2.44772 0 3V12C0 12.5523 0.447715 13 1 13H14C14.5523 13 15 12.5523 15 12V3C15 2.44772 14.5523 2 14 2H1ZM1 3L14 3V3.92494C13.9174 3.92486 13.8338 3.94751 13.7589 3.99505L7.5 7.96703L1.24112 3.99505C1.16621 3.94751 1.0826 3.92486 1 3.92494V3ZM1 4.90797V12H14V4.90797L7.74112 8.87995C7.59394 8.97335 7.40606 8.97335 7.25888 8.87995L1 4.90797Z"
        fill-rule="evenodd"
        clip-rule="evenodd"
      ></path>
    </symbol>

    <symbol id="clear" viewBox="64 64 896 896">
      <path
        d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm165.4 618.2l-66-.3L512 563.4l-99.3 118.4-66.1.3c-4.4 0-8-3.5-8-8 0-1.9.7-3.7 1.9-5.2l130.1-155L340.5 359a8.32 8.32 0 01-1.9-5.2c0-4.4 3.6-8 8-8l66.1.3L512 464.6l99.3-118.4 66-.3c4.4 0 8 3.5 8 8 0 1.9-.7 3.7-1.9 5.2L553.5 514l130 155c1.2 1.5 1.9 3.3 1.9 5.2 0 4.4-3.6 8-8 8z"
      ></path>
    </symbol>

    <symbol id="lock" viewBox="0 0 15 15">
      <path
        d="M5 4.63601C5 3.76031 5.24219 3.1054 5.64323 2.67357C6.03934 2.24705 6.64582 1.9783 7.5014 1.9783C8.35745 1.9783 8.96306 2.24652 9.35823 2.67208C9.75838 3.10299 10 3.75708 10 4.63325V5.99999H5V4.63601ZM4 5.99999V4.63601C4 3.58148 4.29339 2.65754 4.91049 1.99307C5.53252 1.32329 6.42675 0.978302 7.5014 0.978302C8.57583 0.978302 9.46952 1.32233 10.091 1.99162C10.7076 2.65557 11 3.57896 11 4.63325V5.99999H12C12.5523 5.99999 13 6.44771 13 6.99999V13C13 13.5523 12.5523 14 12 14H3C2.44772 14 2 13.5523 2 13V6.99999C2 6.44771 2.44772 5.99999 3 5.99999H4ZM3 6.99999H12V13H3V6.99999Z"
        fill-rule="evenodd"
        clip-rule="evenodd"
      ></path>
    </symbol>

    <symbol id="eye" viewBox="64 64 896 896">
      <path
        d="M942.2 486.2C847.4 286.5 704.1 186 512 186c-192.2 0-335.4 100.5-430.2 300.3a60.3 60.3 0 000 51.5C176.6 737.5 319.9 838 512 838c192.2 0 335.4-100.5 430.2-300.3 7.7-16.2 7.7-35 0-51.5zM512 766c-161.3 0-279.4-81.8-362.7-254C232.6 339.8 350.7 258 512 258s279.4 81.8 362.7 254C791.5 684.2 673.4 766 512 766z"
      ></path>
      <path
        d="M508 336c-97.2 0-176 78.8-176 176s78.8 176 176 176 176-78.8 176-176-78.8-176-176-176zm0 288c-61.9 0-112-50.1-112-112s50.1-112 112-112 112 50.1 112 112-50.1 112-112 112z"
      ></path>
    </symbol>

    <symbol id="eye-invisible" viewBox="64 64 896 896">
      <path
        d="M942.2 486.2Q889.47 375.11 816.7 305l-50.88 50.88C807.31 395.53 843.45 447.4 874.7 512 791.5 684.2 673.4 766 512 766q-72.67 0-133.87-22.38L323 798.75Q408 838 512 838q288.3 0 430.2-300.3a60.29 60.29 0 000-51.5zm-63.57-320.64L836 122.88a8 8 0 00-11.32 0L715.31 232.2Q624.86 186 512 186q-288.3 0-430.2 300.3a60.3 60.3 0 000 51.5q56.69 119.4 136.5 191.41L112.48 835a8 8 0 000 11.31L155.17 889a8 8 0 0011.31 0l712.15-712.12a8 8 0 000-11.32zM149.3 512C232.6 339.8 350.7 258 512 258c54.54 0 104.13 9.36 149.12 28.39l-70.3 70.3a176 176 0 00-238.13 238.13l-83.42 83.42C223.1 637.49 183.3 582.28 149.3 512zm246.7 0a112.11 112.11 0 01146.2-106.69L401.31 546.2A112 112 0 01396 512z"
      ></path>
      <path
        d="M508 624c-3.46 0-6.87-.16-10.25-.47l-52.82 52.82a176.09 176.09 0 00227.42-227.42l-52.82 52.82c.31 3.38.47 6.79.47 10.25a111.94 111.94 0 01-112 112z"
      ></path>
    </symbol>
  </svg>
  <div class="bg"></div>
  <div class="form-container">
    <h1>进入奇妙世界</h1>
    <!-- .prevent阻止默认事件 相当于event.preventDefault-->
    <form class="form" @submit.prevent="login">
      <label>
        <svg class="icon left" aria-hidden="true">
          <use xlink:href="#email"></use>
        </svg>
        <input
          type="email"
          inputmode="email"
          placeholder="请输入邮箱号"
          autocomplete="email"
          name="email"
          v-model.trim="account.email"
          @input="checkEmail"
        />
        <svg class="icon right" aria-hidden="true" @click="clear">
          <use xlink:href="#clear"></use>
        </svg>
      </label>
      <h4 data-value="account"></h4>

      <label>
        <svg class="icon left" aria-hidden="true">
          <use xlink:href="#lock"></use>
        </svg>
        <input
          type="password"
          placeholder="请输入密码"
          autocomplete="current-password"
          maxlength="16"
          name="encrypted"
          v-model.trim="account.password"
          @input="checkEncrypted"
        />
        <svg class="icon right" aria-hidden="true" @click="eye">
          <use xlink:href="#eye"></use>
        </svg>
      </label>
      <h4 data-value="password"></h4>
      <a class="btn-login" href="https://account.redli.cn/register">注册</a>
      <button class="btn-login">登录</button>
    </form>
  </div>
  <ul class="bg-squares"></ul>
</template>
<script lang="ts" setup>
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/userStore'
import { getKeyApi } from '@/apis/user'
import JSEncrypt from 'jsencrypt'
const router = useRouter()
const userStore = useUserStore()
const account = reactive({ email: '', password: '' })
let isInvalid = ref<boolean>(false)
let verifyNumber = ref<number>(0b0)
watch(verifyNumber, () => {
  // console.log('有变动')
  if (verifyNumber.value === 0b11) {
    isInvalid.value = true
    // console.log('合法数据')
  } else {
    isInvalid.value = false
    // console.log('非法数据')
  }
})

onMounted(() => {
  const spanElemnt = document.querySelector('.footer_hide') as HTMLSpanElement
  // console.log('挂载', spanElemnt)
  // spanElemnt.style.display = 'inline'
  spanElemnt.classList.remove('footer_hide')
  if (userStore.userInfo.token) {
    ElMessage({
      message: '你已经登录，即将跳转',
    })
    const timeID = setTimeout(() => {
      router.replace({ path: '/' })
      clearTimeout(timeID)
    }, 3000)
  }
})
onUnmounted(() => {
  const spanElemnt = (document.querySelector('.footer') as HTMLElement).querySelector('span')
  // console.log('卸载', spanElemnt)
  spanElemnt?.classList.add('footer_hide')
})
// 清空输入框
const clear = (event: Event) => {
  // 隐藏自身
  const svgElement = event.currentTarget as SVGSVGElement
  svgElement.removeAttribute('style')
  const inputElements = document.querySelectorAll('input')

  // 将伪的数组转为真数组，并且遍历数组，将main的子元素的高度放到数组中
  const sectionsArrayArray = Array.from(inputElements)

  // Array.from(sections)
  sectionsArrayArray.forEach((item) => {
    item.removeAttribute('style')
  })

  // 清空输入框的值
  account.email = ''
  account.password = ''
}

// 显示密码
const eye = (event: Event) => {
  const svgElement = event.currentTarget as SVGSVGElement
  const inputElement = svgElement.parentNode?.querySelector('input')
  const useElement = svgElement.querySelector('use')

  if (inputElement && useElement) {
    if (inputElement.type == 'text') {
      useElement.setAttribute('xlink:href', '#eye')

      inputElement.type = 'password'
    } else {
      useElement.setAttribute('xlink:href', '#eye-invisible')
      inputElement.type = 'text'
    }
  }
}
// 检测邮箱
const checkEmail = (event: Event) => {
  // 获取触发事件的DOM元素
  const inputElement = event.currentTarget as HTMLInputElement
  const svgElement = (inputElement.parentNode as HTMLLabelElement).querySelector(
    'svg.right'
  ) as SVGSVGElement
  if (svgElement) {
    if (inputElement.value) {
      svgElement.style.display = 'block'
    } else {
      svgElement.removeAttribute('style')
    }
  }

  // 改变DOM的样式
  inputElement.style.setProperty('background-color', 'rgba(238, 119, 82, 0.5)')
  // if (inputElement.value.length < 6) {
  //   console.log('密码长度过短')
  //   return
  // }
  // // 去除空格
  // account.email = account.email.trim()
  // 全部小写
  account.email = account.email.toLowerCase()
  const regx = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  let mask = 0b10 // 要修改的位的掩码
  if (regx.test(account.email)) {
    // 修改二进制数字的特定位
    verifyNumber.value = verifyNumber.value | mask // 使用按位或操作符来修改二进制数字的特定位
    inputElement.style.backgroundColor = 'var(--theme-color-heising)'
  } else {
    inputElement.style.backgroundColor = 'rgba(255, 0, 0, 0.5)'

    verifyNumber.value = verifyNumber.value & ~mask // 使用按位或操作符来修改二进制数字的特定位
  }
}

// 检测密码
const checkEncrypted = (event: Event) => {
  // 获取触发事件的DOM元素
  const inputElement = event.currentTarget as HTMLInputElement
  const svgElement = (inputElement.parentNode as HTMLLabelElement).querySelector(
    'svg.right'
  ) as SVGSVGElement
  if (svgElement) {
    if (inputElement.value) {
      svgElement.style.display = 'block'
    } else {
      svgElement.removeAttribute('style')
    }
  }

  // 改变DOM的样式
  // inputElement.style.setProperty('background-color', 'rgba(238, 119, 82, 0.5)')

  // 去除空格
  // account.password = account.password.trim()
  const regx = new RegExp(
    '^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[/!@#$%^&*(),.?"\\[\\];:{}|<>\\\\])([a-zA-Z0-9/!@#$%^&*(),.?"\\[\\];:{}|<>\\\\]){6,16}$'
  )
  let mask = 0b1 // 要修改的位的掩码
  if (regx.test(account.password)) {
    // 修改二进制数字的特定位
    inputElement.style.backgroundColor = 'var(--theme-color-heising)'
    verifyNumber.value = verifyNumber.value | mask // 使用按位或操作符来修改二进制数字的特定位
  } else {
    inputElement.style.backgroundColor = 'rgba(255, 0, 0, 0.5)'
    verifyNumber.value = verifyNumber.value & ~mask // 使用按位或操作符来修改二进制数字的特定位
  }
}

const login = async () => {
  if (isInvalid.value) {
    const email = account.email
    const result = await getKeyApi()
    const nanoid = result.data.nanoid

    const encryptor = new JSEncrypt() // 新建JSEncrypt对象
    encryptor.setPublicKey(result.data.publicKey) // 设置公钥
    const encrypted = encryptor.encrypt(result.data.nanoid + account.password)

    if (encrypted) {
      // 让按钮失焦
      const buttonElement = document.activeElement as HTMLButtonElement
      buttonElement.blur()
      await userStore.getUserInfo({ email, encrypted, nanoid }).then(() => {
        router.replace({ path: '/' })
      })
    }

    // 2. 跳转首页
  } else {
    // 错误提示
    ElMessage({
      type: 'warning',
      message: '数据格式不对，请检查',
      grouping: true,
    })
  }
}
</script>
<style scoped lang="scss">
.bg {
  position: fixed;
  height: 100vh;
  width: 100vw;
  background-color: var(--theme-color-heising);
  z-index: -100;
}

.form-container {
  height: 100vh;
  width: 100vw;
  /* 溢出隐藏 */
  overflow: hidden;
  background-size: cover;

  text-align: center;
  color: #fff;
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
  place-items: center;
  place-content: center;
  // transition: all 0.6s;
  transition: background-color 0.6s ease-in-out;

  &:focus-within {
    /* backdrop-filter: blur(10px);*/
    background-color: rgba(0, 0, 0, 0.5);
    // backdrop-filter: blur(10px);
    /* z-index: -1;
    opacity: 0; */
    /* position: fixed; */
    /* left: 0;
    top: 0;
    width: 100%;
    height: 100%; */
    /* background: radial-gradient(rgba(0,0,0,0) 0,rgba(0,0,0,0.5) 100%),radial-gradient(rgba(0,0,0,0) 33%,rgba(0,0,0,.3) 166%); */
    /* transition: .25s */
  }
  h1 {
    position: relative;
    top: -40px;
    font-size: 30px;
    font-weight: 100;
    letter-spacing: 2px;
    /* margin-bottom: 30px; */
    /* 过渡效果 */
    /* transition: translateY 1s ease-in-out; */
  }
}

.form {
  display: flex;
  /* flex-direction: column; */
  /* align-items: center; */
  /* justify-items: center; */
  justify-content: space-around;
  /* position: relative;
    z-index: 2; */
  gap: normal 20px;
  opacity: 1;
  width: 282px;
  flex-wrap: wrap;
  /* 不透明度改变时的过渡效果 */
  /* transition: opacity 0.5s; */
  // transition: 1s ease-in-out;
  input {
    outline: none;
    /* border: 1px solid rgba(255,255,255,0.4); */
    border: 1px solid rgba(255, 255, 255, 0.5);
    border-right: 1px solid rgba(255, 255, 255, 0.2);
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
    background-color: rgba(255, 255, 255, 0.2);
    width: 250px;
    padding: 10px 15px;
    border-radius: 3px;
    margin: 0 auto 0 auto;
    /* text-align: center; */
    color: #fff;
    font-size: 14px;
    transition: 0.25s;
    &::placeholder {
      color: rgba(255, 255, 255, 0.8);
      font-size: 14px;
      font-weight: 300;
      transition: 0.25s;
      &:hover {
        background-color: rgba(255, 255, 255, 0.4);
        border-color: var(--theme-color-heising);
        /* background-color: #fff; */
        /* width: 300px; */
      }
    }
  }
  h4 {
    width: 250px;
    padding: 0 15px;
    height: 18px;
    margin: 0 auto 0 auto;
    text-align: right;
    color: var(--theme-color-heising);
    font-size: 10px;
    font-weight: lighter;
    transition: 0.25s;
  }
}

.form input:hover::placeholder {
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
  font-weight: 300;
}

.form input:focus::placeholder {
  color: transparent;
}

.form input:focus {
  background-color: rgba(255, 255, 255, 0.3);
  box-shadow: 0px 0px 2px 1px var(--theme-color-heising);
}
.form input:autofill {
  -webkit-transition: all 5000s ease-in-out 0s !important;
  transition: all 5000s ease-in-out 0s !important;
}
.form input:-webkit-autofill {
  -webkit-transition: all 5000s ease-in-out 0s !important;
  transition: all 5000s ease-in-out 0s !important;
}

.btn-login {
  outline: none;
  background-color: #fff;
  color: var(--theme-color-heising);
  border: none;
  /* width: 45%; */
  flex: 1;
  /* flex-basis:auto; */
  padding: 10px 15px;
  border-radius: 3px;
  font-size: 15px;
  cursor: pointer;
  transition: 0.25s;
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-right: 1px solid rgba(255, 255, 255, 0.2);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  /* margin: 0 10px; */
  &:active {
    // background-color: #f5f7f9;
    /* width: 332px; */
    background-color: var(--theme-color-heising);
    color: #fff;
  }
  &:disabled {
    border-right: 1px solid rgba(255, 255, 255, 0.5);
    border-bottom: 1px solid rgba(255, 255, 255, 0.5);
    background-color: dimgrey;
    color: linen;
    opacity: 0.8;
  }
}

label {
  position: relative;
}

.icon {
  fill: rgba(255, 255, 255, 0.5);

  cursor: pointer;
  position: absolute;
  font-size: 18px;

  top: 50%;
  transform: translateY(-50%);

  transition: all 0.3s;

  &.left {
    right: auto;
    left: -26px;
  }
  &.right {
    display: none;
    right: 8px;
  }
  &:active {
    fill: var(--theme-color-heising);
  }
}

@media (any-hover: hover) {
  .btn-login {
    &:hover {
      background-color: rgba(255, 255, 255, 0.5);
    }
  }
  .icon:hover {
    fill: rgba(255, 255, 255, 1);
  }
}
</style>
