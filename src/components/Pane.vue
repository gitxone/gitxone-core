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
import Vue from 'vue'
import { mapMutations } from 'vuex'
import VueDraggableResizable from 'vue-draggable-resizable'
import 'vue-draggable-resizable/dist/VueDraggableResizable.css'

import {
  SET_PANE,
  DEL_PANE,
  INC_TOPZ,
  DEL_REPO,
  SAVE_REPOS,
} from '@/store/mutationTypes.ts'

function handleSocketCreated () {
    this.socket.addEventListener('message', (event) => {
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

export default Vue.extend({
  components: {
    'vue-draggable-resizable': VueDraggableResizable,
  },
  props: [
    'initialCommand', 'post', 'id', 'path',
    'socket', 'loadAll',
    'width', 'height', 'x', 'y', 'z',
  ],
  data: ({initialCommand, processing}) => ({
    result: '',
    error: '',
    processing,
    command: initialCommand || '',
  }),
  methods: {
    handleEnter: function (event) {
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
    },
    handleResize: function (x, y, width, height) {
      this.$store.commit(SET_PANE, {path: this.path, id: this.id, x, y, width, height})
      this.$store.commit(SAVE_REPOS)
    },
    handleDrag: function (x, y) {
      this.$store.commit(SET_PANE, {path: this.path, id: this.id, x, y})
      this.$store.commit(SAVE_REPOS)
    },
    handleClick: function () {
      this.$store.commit(INC_TOPZ, {path: this.path})
      this.$store.commit(SET_PANE, {path: this.path, id: this.id, z: this.repo.topZ})
      this.$store.commit(SAVE_REPOS)
    },
  },
  computed: {
    repo: function () {
      return this.$store.state.repos[this.path]
    },
  },
  mounted: function () {
    handleSocketCreated.call(this)
  },
  updated: function () {
    handleSocketCreated.call(this)
  },

})
</script>
