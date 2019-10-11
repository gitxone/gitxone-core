<template>
<vue-draggable-resizable
  :x="x" :y="y" :z="z"
  :w="width" :h="height"
  :parent="true"
  :draggable="draggable"
  :resizable="resizable"
  dragHandle=".git"
  @dragging="handleDrag" 
  @resizing="handleResize"
>
<div 
  :class="`pane ${post && 'pane-' + post} ${error && 'pane-error'} ${processing ? 'processing' : ''} ${connected ? '' : 'not-connected'}`" 
  @click="handleClick">
  <div class="pane-head">
    <div class="git">
      git
    </div>
    <div class="command">
      <input 
        v-model="command"
        @keyup="handleComplete"
        @keyup.enter="handleEnter"
        @keydown.tab="handleComplement"
        autocomplete="on"
        :placeholder="connected ? 'Command or \'exit\' to close' : 'Wait a moment...'" 
        :list="`comp-list-${id}`"
        :disabled="!connected"
      />
      <datalist :id="`comp-list-${id}`">
        <option 
          v-for="c in candidates" 
          :key="c"
          :value="getCompletion(c)"
        />
      </datalist>
    </div>
    <div class="post-process" v-if="post">
      &amp;&amp; {{ post }}
    </div>
  </div>
  <pre
    :class="`pane-body`"
    :title="error"
    @dblclick="$emit('modalOpen', result)"
    title="Double click to display the text on a modal."
  >{{ result }}</pre>
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

  &.not-connected
    background-color #333333

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
  cursor pointer

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
  INC_CONN,
} from '@/store/mutationTypes'
import { StoreType, RepoType, PaneType } from '@/store/types'
import defaultState from '@/defaultState.json'

function handleExecMessage (data: any) {
  this.result = data.content || ''
  this.error =  data.error || ''
  switch (this.post) {
    case 'exit':
      if (!this.error) {
        this.$store.commit(DEL_PANE, {path: this.path, id: this.id})
      }
      this.sync()
      break
    case 'clear':
      if (!this.error) {
        this.command = ''
        this.$store.commit(SET_PANE, {path: this.path, id: this.id, command: ''})
        this.$store.commit(SAVE_REPOS)
      }
      this.sync()
      break
  }
  this.processing = false
  this.$forceUpdate()
}

function handleCompleteMessage (data: any) {
  this.candidates = data.content.filter((x, i, that) => that.indexOf(x) === i)
}



function makeSocket(this: any) {
  const port = process.env.serverPort || location.port
  const socket = new WebSocket(`ws://${location.hostname}:${port}/git?path=${this.path}`)

  socket.addEventListener('open', (event: Event) => {
    this.connected = true
    if (!this.post) {
      this.send()
    }
  })

  socket.addEventListener('close', (event: Event) => {
    this.connected = false
    setTimeout(() => makeSocket.call(this), 500)
  })

  socket.addEventListener('message', (event: {data: any}) => {
    const data = JSON.parse(event.data)
    switch (data.type) {
      case "exec":
        handleExecMessage.call(this, data)
        break
      case "complete":
        handleCompleteMessage.call(this, data)
        break
    }
  }, false)

  this.socket = socket
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
  draggable?: bool;
  @Prop()
  resizable?: bool;

  connected = false;
  candidates = [];
  result = '';
  error = '';
  processing = false;
  command = this.initialCommand;
  conn = this.$store.state.conn || 0
  socket = null

  handleComplete (this: any) {
    this.complete()
  }
  handleComplement (this: any, event: Event) {
    if (this.command === '') {
      this.complete()
      event.preventDefault()
      return
    }
    const lastToken = this.command.split(' ').pop()
    const cands = this.candidates.filter((c: string) => c.indexOf(lastToken) === 0)
    const common = this.getCommonString(cands)
    this.command = this.getCompletion(common)
    event.preventDefault()
  }

  handleEnter (this: any, event: Event) {
    if (this.command == 'exit') {
      this.$store.commit(DEL_PANE, {path: this.path, id: this.id})
    } 
    else {
      if (!this.command) { return }
      this.processing = true
      this.send()
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
    return isNaN(x) ? 100 : x >= 0 ? x : 0
  }
  get y (this: any): number {
    const y = this.pane.y
    return isNaN(y) ? 100 : y >= 0 ? y : 0
  }
  get z (this: any): number {
    const z = this.pane.z
    return z
  }
  
  get width (this: any): number {
    return this.pane.width || 400
  }
  get height (this: any): number {
    return this.pane.height || 150
  }

  public connSerial (this: any): number {
    return this.$store.state.conn || 0
  }

  public getCommonString (paths: string[]): string {
    let index = 0
    let common = ""
    while (true) {
      const chars = paths.map(p => p[index])
      index ++
      if (chars[0] && chars.every(c => c === chars[0])) {
        common += chars[0]
      } else {
        break
      }
    }
    return common
  }

  public sync(this: any) {
    this.conn ++
    this.$store.commit(INC_CONN)
  }

  public send (this: any) {
    if (!this.command) { return }
    const type = 'exec'
    const params = JSON.stringify({command: this.command, type})
    this.socket.send(params)
    this.conn = this.connSerial()
  }
  
  public complete (this: any, event: Event) {
    const type = 'complete'
    const params = JSON.stringify({command: this.command, type})
    this.socket.send(params)
  }

  public getCompletion(this: any, common: string): string {
    const tokens = this.command.split(' ')
    const lastToken = tokens[tokens.length - 1]
    if (tokens.length == 1) {
      return common || lastToken
    }
    return `${tokens.slice(0, -1).join(' ')} ${common || lastToken}`
  }

  public created () {
    makeSocket.call(this)
  }

  public updated () {
    if (this.connSerial() > this.conn && !this.post) {
      this.send()
    }
  }

}
</script>
