<template>
<div>
<sidebar-menu :menu="menu" :collapsed="true" @item-click="handleMenuClick" :showChild="true">
  <span slot="toggle-icon"><i class="fas fa-angle-right" /></span>
  <span slot="dropdown-icon"><i class="fas fa-caret-right" /></span>
</sidebar-menu>
<div class="container">
  <Pane 
    v-for="[id, pane] in paneItems"
    :id="id"
    :key="id"
    :x="pane.x >= 0 ? pane.x : 0"
    :y="pane.y >= 0 ? pane.y : 0"
    :z="pane.z || 0"
    :width="pane.width"
    :height="pane.height"
    :initialCommand="pane.command"
    :post="pane.post"
    :path="path"
  ></Pane>
</div>
</div>
</template>

<script lang="ts">
import moment from 'moment'
import '@fortawesome/fontawesome-free/js/all.js'

import { Component, Prop, Vue } from 'vue-property-decorator'
import { mapMutations } from 'vuex'
import VueDraggableResizable from 'vue-draggable-resizable'
import { SidebarMenu } from 'vue-sidebar-menu'
import 'vue-sidebar-menu/dist/vue-sidebar-menu.css'

import {
  ADD_PANE,
  INC_TOPZ,
  SET_REPO,
  SAVE_REPOS,
} from '@/store/mutationTypes'
import { RepoType, PanesType } from '@/store/types'
import defaultState from '@/defaultState.json'

import Pane from '@/components/Pane.vue'

interface MenuType {
  header?: boolean,
  hiddenOnCollapse?: boolean,
  title?: string,
  operation?: string,
  command: string,
  post?: string,
  icon?: string,
  href?: string,
  url?: string,
  child?: MenuType[],
}


@Component({
  components: {
    'vue-draggable-resizable': VueDraggableResizable,
    'sidebar-menu': SidebarMenu,
    Pane,
  }
})
export default class VueComponent extends Vue {
  menu = [
    {
      header: true,
      title: 'Gitxone',
      hiddenOnCollapse: true
    },
    {
      title: 'New command',
      operation: 'command',
      post: '',
      icon: 'fa fa-terminal',
      child: [
        {
          title: '&& clear',
          icon: 'far fa-window-restore',
          operation: 'command',
          post: 'clear',
          child: [
            {
              title: 'git commit -m ',
              operation: 'command',
              command: 'commit -m ',
              post: 'clear',
            },
          ],
        },
        {
          title: '&& exit (close)',
          icon: 'fas fa-times',
          operation: 'command',
          post: 'exit',
          child: [
            {
              title: 'git add *',
              operation: 'command',
              command: 'add *',
              post: 'exit',
            },
          ],
        },

      ]
    },
    {
      href: '/',
      title: 'Home',
      icon: 'fa fa-home'
    },
    {
      operation: 'url',
      url: 'https://github.com/gitxone/gitxone-core',
      title: 'Github',
      icon: 'fab fa-github',
    },
  ];

  get path (this: any): string {
    return this.$route.params.pathMatch
  }

  get repo (this: any): RepoType {
    const repo = this.$store.state.repos[this.path] || defaultState
    return repo
  }

  get panes (this: any): PanesType {
    return this.repo.panes
  }

  get paneItems (this: any) {
    return Object.entries(this.panes)
  }

  mounted (this: any) {
    this.$store.commit(SET_REPO, {... this.repo, path: this.path, lastAccess: moment().format('YYYY-MM-DD HH:mm')})
    this.$store.commit(SAVE_REPOS)
  }

  handleMenuClick (this: any, event: Event, item: MenuType) {
    switch(item.operation) {
      case 'command':
        this.$store.commit(INC_TOPZ, {path: this.path})
        this.$store.commit(ADD_PANE, {path: this.path, z: this.repo.topZ, command: item.command || '', post: item.post || ''})
        break
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
  min-height 100vh
</style>
