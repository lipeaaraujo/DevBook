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

export default function Login() {
  const navigate = useNavigate()
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const form = useForm({
    initialValues: { email: '', password: '' },
    validate: {
      email: isEmail('Enter a valid email'),
      password: isNotEmpty('Password is required'),
    },
  })

  const submit = form.onSubmit(async (values) => {
    setLoading(true)
    setError(null)
    try {
      const { token } = await api<{ token: string }>('/login', {
        method: 'POST',
        body: JSON.stringify(values),
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
              DevBook
            </Title>
            <TextInput label="Email" placeholder="you@example.com" {...form.getInputProps('email')} />
            <PasswordInput label="Password" {...form.getInputProps('password')} />
            {error && <Alert color="red">{error}</Alert>}
            <Button type="submit" loading={loading}>
              Log in
            </Button>
            <Text size="sm" ta="center" c="dimmed">
              No account yet?{' '}
              <Anchor component={Link} to="/register" size="sm">
                Register
              </Anchor>
            </Text>
          </Stack>
        </form>
      </Card>
    </Center>
  )
}
