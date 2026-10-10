import { useState } from 'react'
import {
  Alert,
  Anchor,
  Avatar,
  Button,
  Card,
  EmptyState,
  Group,
  Skeleton,
  Stack,
  Text,
} from '@mantine/core'
import { notifications } from '@mantine/notifications'
import { Link, useParams } from 'react-router-dom'
import PostCard from '../components/PostCard'
import PostForm from '../components/PostForm'
import { useFollow, useMe, useUnfollow, useUser, useUserPosts } from '../hooks'

export default function Profile() {
  const { id = '' } = useParams()
  const { data: me } = useMe()
  const isMe = id === me?.id
  const { data: user, isPending, error } = useUser(id)
  const posts = useUserPosts(id)
  const follow = useFollow(id)
  const unfollow = useUnfollow(id)
  const [composing, setComposing] = useState(false)

  if (isPending) return <Skeleton height={140} radius="md" />
  if (error) return <Alert color="red">{error.message}</Alert>
  if (!user.id) return <EmptyState title="User not found" description="This user does not exist." />

  const onFollowError = (e: Error) => notifications.show({ message: e.message, color: 'red' })

  return (
    <Stack>
      <Card withBorder>
        <Group wrap="nowrap" align="flex-start">
          <Avatar size="lg" name={user.name} color="initials" />
          <Stack gap={4} style={{ flex: 1, minWidth: 0 }}>
            <div>
              <Text fw={700} size="lg">
                {user.name}
              </Text>
              <Text c="dimmed">@{user.nickname}</Text>
            </div>
            <Text size="sm" c="dimmed">
              Joined{' '}
              {new Date(user.createdAt).toLocaleDateString(undefined, {
                month: 'long',
                year: 'numeric',
              })}
            </Text>
            <Group gap="md">
              <Anchor component={Link} to={`/users/${id}/followers`} size="sm" c="inherit">
                <Text component="span" fw={600}>
                  {user.followersCount}
                </Text>{' '}
                followers
              </Anchor>
              <Anchor component={Link} to={`/users/${id}/following`} size="sm" c="inherit">
                <Text component="span" fw={600}>
                  {user.followingCount}
                </Text>{' '}
                following
              </Anchor>
            </Group>
            {isMe && (
              <Group gap="xs" mt={4}>
                <Button variant="default" component={Link} to="/settings">
                  Edit profile
                </Button>
                <Button onClick={() => setComposing(true)}>New post</Button>
              </Group>
            )}
          </Stack>
          {!isMe &&
            (user.isFollowing ? (
              <Button
                variant="default"
                onClick={() => unfollow.mutate(undefined, { onError: onFollowError })}
              >
                Unfollow
              </Button>
            ) : (
              <Button onClick={() => follow.mutate(undefined, { onError: onFollowError })}>
                Follow
              </Button>
            ))}
        </Group>
      </Card>

      {posts.isPending ? (
        <Skeleton height={120} radius="md" />
      ) : posts.error ? (
        <Alert color="red">{posts.error.message}</Alert>
      ) : posts.data.length === 0 ? (
        <EmptyState
          title="No posts yet"
          description={
            isMe ? 'Your posts will show up here.' : `@${user.nickname} has not posted anything yet.`
          }
        />
      ) : (
        posts.data.map((post) => <PostCard key={post.id} post={post} />)
      )}

      {isMe && <PostForm opened={composing} onClose={() => setComposing(false)} />}
    </Stack>
  )
}
