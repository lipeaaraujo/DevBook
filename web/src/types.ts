export type User = {
  id: string
  name: string
  nickname: string
  email?: string // only returned for the current user
  createdAt: string
  updatedAt?: string
  followersCount: number
  followingCount: number
  isFollowing: boolean
}

export type Post = {
  id: string
  title: string
  description: string
  authorId: string
  authorName: string
  authorNickname: string
  createdAt: string
  updatedAt?: string
}
