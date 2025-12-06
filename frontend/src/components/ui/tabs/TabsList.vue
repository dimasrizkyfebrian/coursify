<script setup lang="ts">
import type { TabsListProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { TabsList } from "reka-ui"
import { cn } from "@/lib/utils"

const props = defineProps<TabsListProps & { class?: HTMLAttributes["class"] }>()

const delegatedProps = reactiveOmit(props, "class")
</script>

<template>
  <TabsList
    data-slot="tabs-list"
    v-bind="delegatedProps"
    :class="cn(
      'bg-neutral-100 text-neutral-500 inline-flex h-9 w-fit items-center justify-center rounded-lg p-[3px] dark:bg-neutral-800 dark:text-neutral-400',
      props.class,
    )"
  >
    <slot />
  </TabsList>
</template>
