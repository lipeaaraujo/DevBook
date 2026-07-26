import { useState } from 'react'
import { Alert, Button, EmptyState, Group, Skeleton, Stack, Title } from '@mantine/core'
import PostCard from '../components/PostCard'
import PostForm from '../components/PostForm'
import { useFeed } from '../hooks'

export default function Feed() {
  const { data, isPending, error, refetch } = useFeed()
  const [composing, setComposing] = useState(false)

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={3}>Feed</Title>
        <Button onClick={() => setComposing(true)}>New post</Button>
      </Group>

      {isPending ? (
        <>
          <Skeleton height={120} radius="md" />
          <Skeleton height={120} radius="md" />
          <Skeleton height={120} radius="md" />
        </>
      ) : error ? (
        <Alert color="red" title="Could not load your feed">
          <Stack align="flex-start" gap="xs">
            {error.message}
            <Button variant="light" color="red" size="xs" onClick={() => refetch()}>
              Retry
            </Button>
          </Stack>
        </Alert>
      ) : data.length === 0 ? (
        <EmptyState
          title="Your feed is empty"
          description="The feed shows posts from people you follow — your own posts never appear here. Use the search bar in the header to find someone and follow them from their profile."
        />
      ) : (
        data.map((post) => <PostCard key={post.id} post={post} />)
      )}

      <PostForm opened={composing} onClose={() => setComposing(false)} />
    </Stack>
  )
}
