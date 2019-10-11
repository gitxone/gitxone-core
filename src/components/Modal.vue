<template>
  <transition name="modal" appear>
    <div class="modal modal-overlay" @click.self="$emit('close')">
      <div class="modal-content">
        <slot/>
      </div>
    </div>
  </transition>
</template>

<style lang="stylus" scoped>
.modal 
  &-overlay
    display flex
    align-items center
    justify-content center
    position fixed
    z-index 30
    top 0
    left 0
    width 100%
    height 100%
    background rgba(100, 100, 100, 0.5)
  
  &-content
    background-color #000000
    opacity 0.9
    position relative
    z-index auto
    width 70%
    height 100vh

.modal-enter-active, .modal-leave-active
  transition opacity 0.4s

.modal-leave-active
  transition opacity 0.6s ease 0.4s

.modal-enter, .modal-leave-to
  opacity 0


</style>
