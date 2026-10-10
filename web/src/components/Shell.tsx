import { AppShell, Avatar, Container, Group, Menu, Text, UnstyledButton } from '@mantine/core'
import { useQueryClient } from '@tanstack/react-query'
import { Link, Outlet, useNavigate } from 'react-router-dom'
import { clearToken } from '../api'
import { useMe } from '../hooks'
import UserSearch from './UserSearch'

export default function Shell() {
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { data: me } = useMe()
  return (
    <AppShell header={{ height: 56 }} padding="md">
      <AppShell.Header>
        <Group h="100%" px="md" gap="sm" wrap="nowrap">
          <Text
            component={Link}
            to="/"
            fw={700}
            size="lg"
            style={{ textDecoration: 'none', color: 'inherit', flexShrink: 0 }}
          >
            DevBook
          </Text>
          <div style={{ flex: 1, minWidth: 0, maxWidth: 420, marginInline: 'auto' }}>
            <UserSearch />
          </div>
          <Menu position="bottom-end" width={180}>
            <Menu.Target>
              <UnstyledButton aria-label="Account menu" style={{ flexShrink: 0 }}>
                <Avatar radius="xl" />
              </UnstyledButton>
            </Menu.Target>
            <Menu.Dropdown>
              <Menu.Item component={Link} to={`/users/${me?.id}`} disabled={!me}>
                My profile
              </Menu.Item>
              <Menu.Item component={Link} to="/settings">
                Settings
              </Menu.Item>
              <Menu.Divider />
              <Menu.Item
                color="red"
                onClick={() => {
                  clearToken()
                  qc.clear()
                  navigate('/login')
                }}
              >
                Log out
              </Menu.Item>
            </Menu.Dropdown>
          </Menu>
        </Group>
      </AppShell.Header>
      <AppShell.Main>
        <Container size="sm" px={0}>
          <Outlet />
        </Container>
      </AppShell.Main>
    </AppShell>
  )
}
