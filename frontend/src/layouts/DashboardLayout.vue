<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useUserStore } from '@/stores/user'
import AppSidebar from '@/components/layout/AppSidebar.vue'

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Separator } from '@/components/ui/separator'
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'

const { breadcrumbs } = useUserStore()
</script>

<template>
  <SidebarProvider>
    <AppSidebar />
    <SidebarInset>
      <main class="flex-1 p-8 m-2 overflow-y-auto rounded md:rounded-lg shadow-xl border">
        <header
          class="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12"
        >
          <div class="flex items-center gap-2 px-4">
            <SidebarTrigger class="-ml-1" />
            <Separator orientation="vertical" class="mr-2 h-4" />
            <Breadcrumb>
              <BreadcrumbList>
                <template v-for="(crumb, index) in breadcrumbs" :key="crumb.text">
                  <BreadcrumbItem>
                    <BreadcrumbLink v-if="index < breadcrumbs.length - 1" as-child>
                      <RouterLink :to="crumb.to">{{ crumb.text }}</RouterLink>
                    </BreadcrumbLink>
                    <BreadcrumbPage v-else>
                      {{ crumb.text }}
                    </BreadcrumbPage>
                  </BreadcrumbItem>
                  <BreadcrumbSeparator v-if="index < breadcrumbs.length - 1" />
                </template>
              </BreadcrumbList>
            </Breadcrumb>
          </div>
        </header>
        <div class="p-4 flex-1 overflow-y-auto">
          <slot />
        </div>
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>
