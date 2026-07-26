import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api, currentUserId } from './api'
import type { Post, User } from './types'

export const queryKeys = {
  feed: ['feed'] as const,
  user: (id: string) => ['user', id] as const,
  userPosts: (id: string) => ['userPosts', id] as const,
  searchUsers: (query: string) => ['searchUsers', query] as const,
}

export const useFeed = () =>
  useQuery({ queryKey: queryKeys.feed, queryFn: () => api<Post[]>('/feed') })

export const useUser = (id: string) =>
  useQuery({ queryKey: queryKeys.user(id), queryFn: () => api<User>(`/users/${id}`) })

export const useUserPosts = (id: string) =>
  useQuery({ queryKey: queryKeys.userPosts(id), queryFn: () => api<Post[]>(`/user/${id}/post`) })

export const useSearchUsers = (query: string) =>
  useQuery({
    queryKey: queryKeys.searchUsers(query),
    queryFn: () => api<User[]>(`/users?user=${encodeURIComponent(query)}`),
    enabled: query.trim().length > 0,
  })

function useFollowAction(userId: string, action: 'follow' | 'unfollow') {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api<void>(`/users/${userId}/${action}`, { method: 'POST' }),
    onMutate: async () => {
      await qc.cancelQueries({ queryKey: queryKeys.user(userId) })
      const prev = qc.getQueryData<User>(queryKeys.user(userId))
      if (prev) {
        qc.setQueryData<User>(queryKeys.user(userId), {
          ...prev,
          isFollowing: action === 'follow',
          followersCount: prev.followersCount + (action === 'follow' ? 1 : -1),
        })
      }
      return { prev }
    },
    onError: (_err, _vars, ctx) => {
      if (ctx?.prev) qc.setQueryData(queryKeys.user(userId), ctx.prev)
    },
    onSettled: () => {
      qc.invalidateQueries({ queryKey: queryKeys.user(userId) })
      qc.invalidateQueries({ queryKey: queryKeys.feed })
    },
  })
}

export const useFollow = (userId: string) => useFollowAction(userId, 'follow')
export const useUnfollow = (userId: string) => useFollowAction(userId, 'unfollow')

function useInvalidateOwnPosts() {
  const qc = useQueryClient()
  return () => {
    const id = currentUserId()
    if (id) qc.invalidateQueries({ queryKey: queryKeys.userPosts(id) })
  }
}

export function useCreatePost() {
  const invalidate = useInvalidateOwnPosts()
  return useMutation({
    mutationFn: (body: { title: string; description: string }) =>
      api<{ id: string }>('/post', { method: 'POST', body: JSON.stringify(body) }),
    onSuccess: invalidate,
  })
}

export function useUpdatePost() {
  const invalidate = useInvalidateOwnPosts()
  return useMutation({
    mutationFn: ({ id, ...body }: { id: string; title: string; description: string }) =>
      api<void>(`/post/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
    onSuccess: invalidate,
  })
}

export function useDeletePost() {
  const invalidate = useInvalidateOwnPosts()
  return useMutation({
    mutationFn: (id: string) => api<void>(`/post/${id}`, { method: 'DELETE' }),
    onSuccess: invalidate,
  })
}

export function useUpdateProfile() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: { name: string; nickname: string; email: string }) =>
      api<void>(`/users/${currentUserId()}`, { method: 'PUT', body: JSON.stringify(body) }),
    onSuccess: () => {
      const id = currentUserId()
      if (id) {
        qc.invalidateQueries({ queryKey: queryKeys.user(id) })
        qc.invalidateQueries({ queryKey: queryKeys.userPosts(id) })
      }
    },
  })
}

export function useChangePassword() {
  return useMutation({
    mutationFn: (body: { currentPassword: string; newPassword: string }) =>
      api<void>(`/users/${currentUserId()}/change-password`, { method: 'POST', body: JSON.stringify(body) }),
  })
}
