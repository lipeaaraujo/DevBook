import { useEffect } from 'react'
import { Alert, Button, Group, Modal, Textarea, TextInput, Stack } from '@mantine/core'
import { useForm } from '@mantine/form'
import { notifications } from '@mantine/notifications'
import { useCreatePost, useUpdatePost } from '../hooks'
import type { Post } from '../types'

const TITLE_MAX = 50
const DESCRIPTION_MAX = 2000

export default function PostForm({
  opened,
  onClose,
  post,
}: {
  opened: boolean
  onClose: () => void
  post?: Post
}) {
  const create = useCreatePost()
  const update = useUpdatePost()
  const mutation = post ? update : create

  const form = useForm({
    initialValues: { title: '', description: '' },
    validate: {
      title: (v) =>
        !v.trim() ? 'Title is required' : v.length > TITLE_MAX ? `Max ${TITLE_MAX} characters` : null,
      description: (v) =>
        !v.trim()
          ? 'Description is required'
          : v.length > DESCRIPTION_MAX
            ? `Max ${DESCRIPTION_MAX} characters`
            : null,
    },
  })

  useEffect(() => {
    if (opened) {
      form.setValues({ title: post?.title ?? '', description: post?.description ?? '' })
      form.clearErrors()
      mutation.reset()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [opened])

  const submit = form.onSubmit((values) => {
    const options = {
      onSuccess: () => {
        notifications.show({ message: post ? 'Post updated' : 'Post created', color: 'green' })
        onClose()
      },
    }
    if (post) update.mutate({ id: post.id, ...values }, options)
    else create.mutate(values, options)
  })

  return (
    <Modal opened={opened} onClose={onClose} title={post ? 'Edit post' : 'New post'} centered>
      <form onSubmit={submit}>
        <Stack>
          <TextInput
            label="Title"
            withAsterisk
            maxLength={TITLE_MAX}
            description={`${form.values.title.length}/${TITLE_MAX}`}
            data-autofocus
            {...form.getInputProps('title')}
          />
          <Textarea
            label="Description"
            withAsterisk
            autosize
            minRows={4}
            maxRows={12}
            maxLength={DESCRIPTION_MAX}
            description={`${form.values.description.length}/${DESCRIPTION_MAX}`}
            {...form.getInputProps('description')}
          />
          {mutation.error && <Alert color="red">{mutation.error.message}</Alert>}
          <Group justify="flex-end">
            <Button variant="default" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" loading={mutation.isPending}>
              {post ? 'Save' : 'Post'}
            </Button>
          </Group>
        </Stack>
      </form>
    </Modal>
  )
}
