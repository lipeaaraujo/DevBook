import { useEffect } from 'react'
import { Button, Card, PasswordInput, Stack, TextInput, Title } from '@mantine/core'
import { isEmail, isNotEmpty, matchesField, useForm } from '@mantine/form'
import { notifications } from '@mantine/notifications'
import { currentUserId } from '../api'
import { useChangePassword, useUpdateProfile, useUser } from '../hooks'

export default function Settings() {
  const { data: user } = useUser(currentUserId() ?? '')
  const updateProfile = useUpdateProfile()
  const changePassword = useChangePassword()

  const profileForm = useForm({
    initialValues: { name: '', nickname: '', email: '' },
    validate: {
      name: isNotEmpty('Name is required'),
      nickname: isNotEmpty('Nickname is required'),
      email: isEmail('Enter a valid email'),
    },
  })

  useEffect(() => {
    if (user?.id) {
      profileForm.setValues({ name: user.name, nickname: user.nickname, email: user.email })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user])

  const passwordForm = useForm({
    initialValues: { currentPassword: '', newPassword: '', confirmPassword: '' },
    validate: {
      currentPassword: isNotEmpty('Current password is required'),
      newPassword: isNotEmpty('New password is required'),
      confirmPassword: matchesField('newPassword', 'Passwords do not match'),
    },
  })

  return (
    <Stack>
      <Card withBorder>
        <form
          onSubmit={profileForm.onSubmit((values) =>
            updateProfile.mutate(values, {
              onSuccess: () => notifications.show({ message: 'Profile updated', color: 'green' }),
              onError: (e) => notifications.show({ message: e.message, color: 'red' }),
            }),
          )}
        >
          <Stack>
            <Title order={4}>Profile</Title>
            <TextInput label="Name" {...profileForm.getInputProps('name')} />
            <TextInput label="Nickname" {...profileForm.getInputProps('nickname')} />
            <TextInput label="Email" {...profileForm.getInputProps('email')} />
            <Button type="submit" loading={updateProfile.isPending} style={{ alignSelf: 'flex-start' }}>
              Save profile
            </Button>
          </Stack>
        </form>
      </Card>

      <Card withBorder>
        <form
          onSubmit={passwordForm.onSubmit(({ currentPassword, newPassword }) =>
            changePassword.mutate(
              { currentPassword, newPassword },
              {
                onSuccess: () => {
                  notifications.show({ message: 'Password changed', color: 'green' })
                  passwordForm.reset()
                },
                onError: (e) => notifications.show({ message: e.message, color: 'red' }),
              },
            ),
          )}
        >
          <Stack>
            <Title order={4}>Password</Title>
            <PasswordInput
              label="Current password"
              {...passwordForm.getInputProps('currentPassword')}
            />
            <PasswordInput label="New password" {...passwordForm.getInputProps('newPassword')} />
            <PasswordInput
              label="Confirm new password"
              {...passwordForm.getInputProps('confirmPassword')}
            />
            <Button
              type="submit"
              loading={changePassword.isPending}
              style={{ alignSelf: 'flex-start' }}
            >
              Change password
            </Button>
          </Stack>
        </form>
      </Card>
    </Stack>
  )
}
