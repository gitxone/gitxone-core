<template>
<vue-draggable-resizable
  :x="x" :y="y" :z="z"
  :w="width" :h="height"
  :parent="true"
  dragHandle=".git"
  @dragging="handleDrag" 
  @resizing="handleResize"
>
<div 
  :class="`pane ${post && 'pane-' + post} ${error && 'pane-error'} ${processing ? 'processing' : ''}`" 
  @click="handleClick">
  <div class="pane-head">
    <div class="git">
      git
    </div>
    <div class="command">
      <input 
        placeholder="command or 'exit' to close" 
        v-model="command"
        @keyup.enter="handleEnter"
      />
    </div>
    <div class="post-process" v-if="post">
      &amp;&amp; {{ post }}
    </div>
  </div>
  <pre :class="`pane-body`" :title="error">{{ result }}</pre>
</div>
</vue-draggable-resizable>
</template>

<style lang="stylus">
.pane
  background-color #000000
  color #ffffff
  padding 0
  font-family "Courier New", Consolas, monospace
  border solid 1px #ffffff
  position absolute
  opacity 0.8
  display flex
  flex-direction column
  width 100%
  height 100%

  &.processing
    background-color #000044

  &.pane-error
    background-color #440000

  &.pane-clear
    border-color #aaffff

  &.pane-exit
    border-color #ffaaaa

  .pane-head
    display -webkit-flex
    display flex
    flex-direction row
    .git
      padding 5px
      cursor move

    .command
      flex 1
      input
        width 100% 
        background-color transparent
        color #ffffff
        padding 5px
        border none
        font-size 16px
        font-family "Courier New", Consolas, monospace
    .post-process
      padding 5px
  .pane-body
    padding 5px
    line-height 20px
    border-top solid 1px #444444
    font-family "Courier New", Consolas, monospace
    overflow auto
    flex 1

</style>

<script lang="ts">
import { mapMutations, Store } from 'vuex'
import { Prop, Vue } from 'vue-property-decorator'
import Component from 'vue-class-component'
import VueDraggableResizable from 'vue-draggable-resizable'
import 'vue-draggable-resizable/dist/VueDraggableResizable.css'

import {
  SET_PANE,
  DEL_PANE,
  INC_TOPZ,
  DEL_REPO,
  SAVE_REPOS,
} from '@/store/mutationTypes'
import { StoreType, RepoType, PaneType } from '@/store/types'
import defaultState from '@/defaultState.json'

function handleSocketCreated (this: any) {
    this.socket.addEventListener('message', (event: {data: any}) => {
      const data = JSON.parse(event.data)
      if (this.id !== data.id) {
        return
      }
      this.result = data.content || ''
      this.error =  data.error || ''
      switch (this.post) {
        case 'exit':
          this.$store.commit(DEL_PANE, {path: this.path, id: this.id})
          this.loadAll()
          break
        case 'clear':
          this.command = ''
          this.$store.commit(SET_PANE, {path: this.path, id: this.id, command: ''})
          this.loadAll()
          break
      }
      this.$store.commit(SAVE_REPOS)
      this.processing = false
      this.$forceUpdate()
    }, false)

}

@Component({
  components: {
    'vue-draggable-resizable': VueDraggableResizable,
  }
})
export default class VueComponent extends Vue {
  @Prop()
  initialCommand?: string;
  @Prop()
  post?: string;
  @Prop()
  id?: string;
  @Prop()
  path?: string;
  @Prop()
  socket?: WebSocket ;
  @Prop()
  loadAll?: Function;

  result = '';
  error = '';
  processing = false;
  command = this.initialCommand;

  handleEnter (this: any, event: Event) {
    if (this.command == 'exit') {
      this.$store.commit(DEL_PANE, {path: this.path, id: this.id})
    } 
    else {
      if (!this.command) { return }
      this.processing = true
      const params = JSON.stringify({command: this.command, id: this.id})
      this.socket.send(params)
      this.$store.commit(SET_PANE, {path: this.path, id: this.id, command: this.command})
    }
    this.$store.commit(SAVE_REPOS)
  }
  handleResize (this: any, x: number, y: number, width: number, height: number) {
    this.$store.commit(SET_PANE, {path: this.path, id: this.id, x, y, width, height})
    this.$store.commit(SAVE_REPOS)
  }
  handleDrag (this: any, x: number, y: number) {
    this.$store.commit(SET_PANE, {path: this.path, id: this.id, x, y})
    this.$store.commit(SAVE_REPOS)
  }
  handleClick (this: any) {
    this.$store.commit(INC_TOPZ, {path: this.path})
    this.$store.commit(SET_PANE, {path: this.path, id: this.id, z: this.repo.topZ})
    this.$store.commit(SAVE_REPOS)
  }

  get repo (this: any): RepoType {
    return this.$store.state.repos[this.path] || defaultState
  }
  get pane (this: any): PaneType {
    return this.repo.panes[this.id] || {}
  }

  get x (this: any): number {
    const x = this.pane.x
    return isNaN(x) ? 100 : x
  }
  get y (this: any): number {
    const y = this.pane.y
    return isNaN(y) ? 100 : y
  }
  get z (this: any): number {
    const z = this.pane.z
    return z
  }
  get width (this: any): number {
    return this.pane.width || 300
  }
  get height (this: any): number {
    return this.pane.height || 100
  }

  public mounted () {
    handleSocketCreated.call(this)
  }
  public updated () {
    handleSocketCreated.call(this)
  }
}
</script>
