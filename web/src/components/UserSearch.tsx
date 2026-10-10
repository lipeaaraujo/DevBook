import { useState } from 'react'
import { Autocomplete } from '@mantine/core'
import { useDebouncedValue } from '@mantine/hooks'
import { useNavigate } from 'react-router-dom'
import { useMe, useSearchUsers } from '../hooks'

const label = (u: { name: string; nickname: string }) => `${u.name} @${u.nickname}`

export default function UserSearch() {
  const navigate = useNavigate()
  const [value, setValue] = useState('')
  const [debounced] = useDebouncedValue(value, 250)
  const { data } = useSearchUsers(debounced)
  const { data: me } = useMe()
  const users = (data ?? []).filter((u) => u.id !== me?.id)

  return (
    <Autocomplete
      placeholder="Search users"
      aria-label="Search users"
      value={value}
      onChange={setValue}
      data={users.map(label)}
      filter={({ options }) => options}
      onOptionSubmit={(selected) => {
        const user = users.find((u) => label(u) === selected)
        if (user) navigate(`/users/${user.id}`)
        setValue('')
      }}
    />
  )
}
