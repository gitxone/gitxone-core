<template>
<div>
<sidebar-menu :menu="menu" :collapsed="true" @item-click="handleMenuClick" :showChild="true">
  <span slot="toggle-icon"><i class="fas fa-angle-right" /></span>
  <span slot="dropdown-icon"><i class="fas fa-caret-right" /></span>
</sidebar-menu>
<div class="container">
  <div class="selector">
    <input v-model="path" @keyup.enter="go" />
    <button @click="go">Go</button>
  </div>
  <table class="history">
    <thead>
      <tr>
        <th>Path</th>
        <th width="170">Last access</th>
        <th width="100"></th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="[p, info] in repoItems" :key="p">
        <td class="path"><nuxt-link :to="`/repo/${p}`">{{ p }}</nuxt-link></td>
        <td class="last-access">{{ info.lastAccess }}</td>
        <td class="operation">
          <button @click="handleDelRepo(p)">Delete history</button>
        </td>
      </tr>
    </tbody>
  </table>
</div>
</div>
</template>

<script lang="ts">
import '@fortawesome/fontawesome-free/js/all.js'

import { mapMutations } from 'vuex'
import { Prop, Vue } from 'vue-property-decorator'
import Component from 'vue-class-component'

import { SidebarMenu } from 'vue-sidebar-menu'
import 'vue-sidebar-menu/dist/vue-sidebar-menu.css'

import {
  DEL_REPO,
  SAVE_REPOS,
} from '@/store/mutationTypes'

import {
  RepoType,
  ReposType,
} from '@/store/types'

@Component({
  components: {
    'sidebar-menu': SidebarMenu,
  }
})
export default class Page extends Vue {
  menu = [
    {
      operation: 'url',
      url: 'https://github.com/gitxone/gitxone-core',
      title: 'Github',
      icon: 'fab fa-github',
    },
  ]
  path = '~/'

  get repos (this: any): ReposType {
    return this.$store.state.repos
  } 

  get repoItems (this:any): any {
    const items = [... Object.entries(this.repos || {})]
    const sorter = (a: any, b: any) => a[1].lastAccess > b[1].lastAccess ? -1 : 1
    return items.sort(sorter)
  }

  public go (this: any) {
    const path = `/repo/${this.path.replace(/^\//g, '')}`
    this.$router.history.push(path)
  }
  public handleDelRepo (path: string) {
    this.$store.commit(DEL_REPO, {path})
    this.$store.commit(SAVE_REPOS)
  }
  public handleMenuClick (event: Event, item: any) {
    switch(item.operation) {
      case 'url':
        window.open(item.url)
        break
    }
  }
}
</script>

<style lang="stylus">
.container
  position relative
  margin-left 50px

.selector
  display flex
  margin 30px

  input
    font-size 30px
    padding 10px
    background-color #222222
    border solid 1px #333333
    color #ffffff
    flex 1

  button
    width 100px
    font-size 30px
    cursor pointer
    background-color #666666
    border solid 1px #333333
    color #ffffff
    font-family serif

.history
  width 100%
  color #ffffff
  border-collapse collapse

  .path
    font-size 20px

    a
      color #aaa
      text-decoration none
      &:hover
        color #fff

  td
    border-top solid 1px #888

  th, td
    padding 10px

  button
    font-size 12px
    background-color #450000
    border solid 1px #990000
    color #ffffff
    cursor pointer

</style>
