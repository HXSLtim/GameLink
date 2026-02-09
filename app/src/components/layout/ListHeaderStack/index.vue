<template>
  <view class="list-header-stack">
    <slot name="nav"></slot>
    <slot name="search"></slot>
    <slot name="banner"></slot>
    <slot name="tabs"></slot>
    <slot name="sort"></slot>
    <slot></slot>
  </view>
</template>

<style lang="scss" scoped>
.list-header-stack {
  display: flex;
  flex-direction: column;
}
</style>
