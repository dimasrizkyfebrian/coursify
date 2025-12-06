<script setup lang="ts">
import type { TabsTriggerProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { TabsTrigger, useForwardProps } from "reka-ui"
import { cn } from "@/lib/utils"

const props = defineProps<TabsTriggerProps & { class?: HTMLAttributes["class"] }>()

const delegatedProps = reactiveOmit(props, "class")

const forwardedProps = useForwardProps(delegatedProps)
</script>

<template>
  <TabsTrigger
    data-slot="tabs-trigger"
    v-bind="forwardedProps"
    :class="cn(
      'data-[state=active]:bg-white dark:data-[state=active]:text-neutral-950 focus-visible:border-neutral-950 focus-visible:ring-neutral-950/50 focus-visible:outline-ring dark:data-[state=active]:border-neutral-200 dark:data-[state=active]:bg-neutral-200/30 text-neutral-950 dark:text-neutral-500 inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center gap-1.5 rounded-md border border-neutral-200 border-transparent px-2 py-1 text-sm font-medium whitespace-nowrap transition-[color,box-shadow] focus-visible:ring-[3px] focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 data-[state=active]:shadow-sm [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*=\'size-\'])]:size-4 dark:data-[state=active]:bg-neutral-950 dark:dark:data-[state=active]:text-neutral-50 dark:focus-visible:border-neutral-300 dark:focus-visible:ring-neutral-300/50 dark:dark:data-[state=active]:border-neutral-800 dark:dark:data-[state=active]:bg-neutral-800/30 dark:text-neutral-50 dark:dark:text-neutral-400 dark:border-neutral-800',
      props.class,
    )"
  >
    <slot />
  </TabsTrigger>
</template>
