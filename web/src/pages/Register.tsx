import { useState } from 'react'
import {
  Alert,
  Anchor,
  Button,
  Card,
  Center,
  PasswordInput,
  Stack,
  Text,
  TextInput,
  Title,
} from '@mantine/core'
import { isEmail, isNotEmpty, useForm } from '@mantine/form'
import { Link, useNavigate } from 'react-router-dom'
import { api, setToken } from '../api'

export default function Register() {
  const navigate = useNavigate()
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const form = useForm({
    initialValues: { name: '', nickname: '', email: '', password: '' },
    validate: {
      name: isNotEmpty('Name is required'),
      nickname: isNotEmpty('Nickname is required'),
      email: isEmail('Enter a valid email'),
      password: isNotEmpty('Password is required'),
    },
  })

  const submit = form.onSubmit(async (values) => {
    setLoading(true)
    setError(null)
    try {
      await api('/users', { method: 'POST', body: JSON.stringify(values) })
      const { token } = await api<{ token: string }>('/login', {
        method: 'POST',
        body: JSON.stringify({ email: values.email, password: values.password }),
      })
      setToken(token)
      navigate('/')
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setLoading(false)
    }
  })

  return (
    <Center mih="100dvh" p="md">
      <Card withBorder w="100%" maw={400} p="lg">
        <form onSubmit={submit}>
          <Stack>
            <Title order={2} ta="center">
              Create your account
            </Title>
            <TextInput label="Name" {...form.getInputProps('name')} />
            <TextInput label="Nickname" {...form.getInputProps('nickname')} />
            <TextInput label="Email" placeholder="you@example.com" {...form.getInputProps('email')} />
            <PasswordInput label="Password" {...form.getInputProps('password')} />
            {error && <Alert color="red">{error}</Alert>}
            <Button type="submit" loading={loading}>
              Register
            </Button>
            <Text size="sm" ta="center" c="dimmed">
              Already have an account?{' '}
              <Anchor component={Link} to="/login" size="sm">
                Log in
              </Anchor>
            </Text>
          </Stack>
        </form>
      </Card>
    </Center>
  )
}
