import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import VueI18n from 'vue-i18n'
Vue.use(VueI18n)

// import { vDebounceThrottle } from 'directives'
// Vue.use(vDebounceThrottle)

const i18n = new VueI18n({
  locale: "en", // 定义默认语言为中文
  messages: {
    zh: require("./locales/zh.json"),//json结构见下方
    en: require("./locales/en.json")
  }
});

Vue.config.productionTip = false
import { Row, Col, Drawer, Form, FormItem, Input, Button, Select, Option, Dialog, Upload, Message, Tabs, TabPane, Table, TableColumn, Loading} from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css';
Vue.use(Row)
Vue.use(Col)
Vue.use(Drawer)
Vue.use(Form)
Vue.use(FormItem)
Vue.use(Input)
Vue.use(Button)
Vue.use(Select)
Vue.use(Option)
Vue.use(Dialog)
Vue.use(Upload)
Vue.use(Tabs)
Vue.use(TabPane)
Vue.use(Table)
Vue.use(TableColumn)
Vue.use(Loading)
Vue.prototype.$message = Message;
new Vue({
  router,
    store,
    i18n,
  render: h => h(App)
}).$mount('#app')
