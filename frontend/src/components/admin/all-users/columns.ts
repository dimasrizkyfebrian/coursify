import type { ColumnDef } from '@tanstack/table-core'
import { h } from 'vue'
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import AllUsersDropdownAction from './AllUsersDropdownAction.vue'

export interface User {
  id: string
  full_name: string
  email: string
  role: 'admin' | 'instructor' | 'student'
  status: 'pending' | 'active' | 'rejected'
}

export const columns: ColumnDef<User>[] = [
  {
    id: 'select',
    header: ({ table }) =>
      h(Checkbox, {
        modelValue:
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && 'indeterminate'),
        'onUpdate:modelValue': (value) => table.toggleAllPageRowsSelected(!!value),
        ariaLabel: 'Select all',
      }),
    cell: ({ row }) =>
      h(Checkbox, {
        modelValue: row.getIsSelected(),
        'onUpdate:modelValue': (value) => row.toggleSelected(!!value),
        ariaLabel: 'Select row',
      }),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: 'full_name',
    header: 'Full Name',
  },
  {
    accessorKey: 'email',
    header: 'Email',
  },
  {
    accessorKey: 'role',
    header: 'Role',
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => {
      const status = row.getValue('status') as string
      const variant =
        status === 'active' ? 'default' : status === 'pending' ? 'outline' : 'destructive'
      return h(Badge, { variant }, () => status)
    },
  },
  {
    id: 'actions',
    cell: ({ row }) => h(AllUsersDropdownAction, { user: row.original }),
  },
]
