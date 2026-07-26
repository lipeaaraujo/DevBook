import { useState } from 'react'
import { ActionIcon, Anchor, Card, Group, Menu, Text } from '@mantine/core'
import { notifications } from '@mantine/notifications'
import { Link } from 'react-router-dom'
import { currentUserId } from '../api'
import { useDeletePost } from '../hooks'
import type { Post } from '../types'
import PostForm from './PostForm'

function relativeTime(iso: string) {
  const diff = (new Date(iso).getTime() - Date.now()) / 1000
  const units: [Intl.RelativeTimeFormatUnit, number][] = [
    ['year', 31536000],
    ['month', 2592000],
    ['week', 604800],
    ['day', 86400],
    ['hour', 3600],
    ['minute', 60],
  ]
  const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' })
  for (const [unit, seconds] of units) {
    if (Math.abs(diff) >= seconds) return rtf.format(Math.round(diff / seconds), unit)
  }
  return 'just now'
}

export default function PostCard({ post }: { post: Post }) {
  const isOwn = post.authorId === currentUserId()
  const [editing, setEditing] = useState(false)
  const deletePost = useDeletePost()

  return (
    <Card withBorder>
      <Group justify="space-between" wrap="nowrap" align="flex-start">
        <div style={{ minWidth: 0 }}>
          <Anchor component={Link} to={`/users/${post.authorId}`} c="inherit" underline="hover">
            <Text component="span" fw={600} size="sm">
              {post.authorName}
            </Text>{' '}
            <Text component="span" c="dimmed" size="sm">
              @{post.authorNickname}
            </Text>
          </Anchor>
          <Text component="span" c="dimmed" size="sm">
            {' · '}
            <time dateTime={post.createdAt} title={new Date(post.createdAt).toLocaleString()}>
              {relativeTime(post.createdAt)}
            </time>
          </Text>
        </div>
        {isOwn && (
          <Menu position="bottom-end" width={140}>
            <Menu.Target>
              <ActionIcon variant="subtle" color="gray" aria-label="Post options">
                ⋯
              </ActionIcon>
            </Menu.Target>
            <Menu.Dropdown>
              <Menu.Item onClick={() => setEditing(true)}>Edit</Menu.Item>
              <Menu.Item
                color="red"
                onClick={() => {
                  if (!window.confirm('Delete this post? This cannot be undone.')) return
                  deletePost.mutate(post.id, {
                    onSuccess: () => notifications.show({ message: 'Post deleted', color: 'green' }),
                    onError: (e) => notifications.show({ message: e.message, color: 'red' }),
                  })
                }}
              >
                Delete
              </Menu.Item>
            </Menu.Dropdown>
          </Menu>
        )}
      </Group>
      <Text fw={600} mt="xs">
        {post.title}
      </Text>
      <Text style={{ whiteSpace: 'pre-wrap', overflowWrap: 'anywhere' }}>{post.description}</Text>
      {isOwn && <PostForm opened={editing} onClose={() => setEditing(false)} post={post} />}
    </Card>
  )
}
