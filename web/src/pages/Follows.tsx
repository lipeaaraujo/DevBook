import { Alert, Anchor, Avatar, Card, EmptyState, Group, Skeleton, Stack, Text, Title } from '@mantine/core'
import { Link, useParams } from 'react-router-dom'
import { useFollowList, useUser } from '../hooks'

export default function Follows({ kind }: { kind: 'followers' | 'following' }) {
  const { id = '' } = useParams()
  const { data: user } = useUser(id)
  const { data, isPending, error } = useFollowList(id, kind)

  return (
    <Stack>
      <div>
        <Title order={3}>{kind === 'followers' ? 'Followers' : 'Following'}</Title>
        {user && (
          <Anchor component={Link} to={`/users/${id}`} size="sm" c="dimmed">
            @{user.nickname}
          </Anchor>
        )}
      </div>

      {isPending ? (
        <Skeleton height={64} radius="md" />
      ) : error ? (
        <Alert color="red">{error.message}</Alert>
      ) : data.length === 0 ? (
        <EmptyState
          title={kind === 'followers' ? 'No followers yet' : 'Not following anyone yet'}
        />
      ) : (
        data.map((u) => (
          <Card
            key={u.id}
            withBorder
            component={Link}
            to={`/users/${u.id}`}
            style={{ textDecoration: 'none', color: 'inherit' }}
          >
            <Group wrap="nowrap">
              <Avatar name={u.name} color="initials" />
              <div style={{ minWidth: 0 }}>
                <Text fw={600}>{u.name}</Text>
                <Text c="dimmed" size="sm">
                  @{u.nickname}
                </Text>
              </div>
            </Group>
          </Card>
        ))
      )}
    </Stack>
  )
}
