<script setup lang="ts">
import { RouterLink } from 'vue-router'

import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { ChevronRight } from 'lucide-vue-next'

import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubItem,
  SidebarMenuSubButton,
} from '@/components/ui/sidebar'

defineProps<{
  items: any[]
}>()
</script>

<template>
  <SidebarGroup>
    <SidebarGroupLabel>Platform</SidebarGroupLabel>
    <SidebarMenu>
      <template v-for="link in items" :key="link.title">
        <Collapsible v-if="link.children" :default-open="true" as-child class="group/collapsible">
          <SidebarMenuItem>
            <CollapsibleTrigger as-child>
              <SidebarMenuButton>
                <component :is="link.icon" class="w-5 h-5 mr-2" />
                <span>{{ link.title }}</span>
                <ChevronRight
                  class="ml-auto h-4 w-4 transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"
                />
              </SidebarMenuButton>
            </CollapsibleTrigger>
            <CollapsibleContent>
              <SidebarMenuSub>
                <SidebarMenuSubItem v-for="child in link.children" :key="child.title">
                  <SidebarMenuSubButton as-child>
                    <RouterLink :to="child.to">
                      {{ child.title }}
                    </RouterLink>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              </SidebarMenuSub>
            </CollapsibleContent>
          </SidebarMenuItem>
        </Collapsible>

        <SidebarMenuItem v-else as-child>
          <RouterLink :to="link.to">
            <SidebarMenuButton>
              <component :is="link.icon" class="w-5 h-5 mr-2" />
              <span>{{ link.title }}</span>
            </SidebarMenuButton>
          </RouterLink>
        </SidebarMenuItem>
      </template>
    </SidebarMenu>
  </SidebarGroup>
</template>
