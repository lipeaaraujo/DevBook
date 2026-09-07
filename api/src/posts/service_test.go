package posts

import (
	"errors"
	"strings"
	"testing"
)

type stubPostRepo struct {
	createID            string
	posts               []Post
	post                Post
	createErr           error
	getErr              error
	getByIDErr          error
	getByAuthorErr      error
	getFromFollowersErr error
	updateErr           error
	deleteErr           error
	createdPost         Post
	getTitle            string
	getPostID           string
	getAuthorID         string
	getAuthorTitle      string
	getFollowerUserID   string
	updatedPost         Post
	deletedPostID       string
}

func (r *stubPostRepo) Create(post Post) (string, error) {
	r.createdPost = post
	return r.createID, r.createErr
}

func (r *stubPostRepo) Get(title string) ([]Post, error) {
	r.getTitle = title
	return r.posts, r.getErr
}

func (r *stubPostRepo) GetById(id string) (Post, error) {
	r.getPostID = id
	return r.post, r.getByIDErr
}

func (r *stubPostRepo) GetByAuthor(authorID, title string) ([]Post, error) {
	r.getAuthorID = authorID
	r.getAuthorTitle = title
	return r.posts, r.getByAuthorErr
}

func (r *stubPostRepo) GetFromFollowers(userID string) ([]Post, error) {
	r.getFollowerUserID = userID
	return r.posts, r.getFromFollowersErr
}

func (r *stubPostRepo) Update(post Post) error {
	r.updatedPost = post
	return r.updateErr
}

func (r *stubPostRepo) Delete(id string) error {
	r.deletedPostID = id
	return r.deleteErr
}

func TestPostService_CreatePost(t *testing.T) {
	createErr := errors.New("create failed")
	tests := []struct {
		name    string
		post    Post
		repo    *stubPostRepo
		wantID  string
		wantErr error
	}{
		{name: "creates post", post: Post{Title: "Title", Description: "Description", AuthorId: "author-id"}, repo: &stubPostRepo{createID: "post-id"}, wantID: "post-id"},
		{name: "rejects missing fields", post: Post{}, repo: &stubPostRepo{}, wantErr: ErrPostTitleDescriptionAuthorIDEmpty},
		{name: "rejects long title", post: Post{Title: strings.Repeat("a", 51), Description: "Description", AuthorId: "author-id"}, repo: &stubPostRepo{}, wantErr: ErrPostTitleTooLong},
		{name: "rejects long description", post: Post{Title: "Title", Description: strings.Repeat("a", 2001), AuthorId: "author-id"}, repo: &stubPostRepo{}, wantErr: ErrPostDescriptionTooLong},
		{name: "returns repository error", post: Post{Title: "Title", Description: "Description", AuthorId: "author-id"}, repo: &stubPostRepo{createErr: createErr}, wantErr: createErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewPostService(tt.repo).CreatePost(tt.post)
			if !errors.Is(err, tt.wantErr) || id != tt.wantID {
				t.Fatalf("expected id %q and error %v, got id %q and error %v", tt.wantID, tt.wantErr, id, err)
			}
			if tt.wantErr == nil && tt.repo.createdPost != tt.post {
				t.Fatalf("expected repository post %+v, got %+v", tt.post, tt.repo.createdPost)
			}
			if tt.wantErr != nil && tt.repo.createdPost != (Post{}) && tt.name != "returns repository error" {
				t.Fatal("repository should not have been called")
			}
		})
	}
}

func TestPostService_GetPosts(t *testing.T) {
	wantErr := errors.New("get failed")
	for _, tt := range []struct {
		name string
		repo *stubPostRepo
		want error
	}{
		{name: "returns posts", repo: &stubPostRepo{posts: []Post{{Id: "post-id"}}}},
		{name: "returns repository error", repo: &stubPostRepo{getErr: wantErr}, want: wantErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			posts, err := NewPostService(tt.repo).GetPosts("title")
			if !errors.Is(err, tt.want) || tt.repo.getTitle != "title" || (tt.want == nil && len(posts) != 1) {
				t.Fatalf("unexpected posts, error, or query: %+v, %v, %q", posts, err, tt.repo.getTitle)
			}
		})
	}
}

func TestPostService_GetByID(t *testing.T) {
	wantErr := errors.New("get by id failed")
	for _, tt := range []struct {
		name string
		repo *stubPostRepo
		id   string
		want error
	}{
		{name: "returns post", repo: &stubPostRepo{post: Post{Id: "post-id"}}, id: "post-id"},
		{name: "rejects empty id", repo: &stubPostRepo{}, want: ErrPostIDRequired},
		{name: "returns repository error", repo: &stubPostRepo{getByIDErr: wantErr}, id: "post-id", want: wantErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			post, err := NewPostService(tt.repo).GetById(tt.id)
			if !errors.Is(err, tt.want) || (tt.want == nil && post.Id != tt.id) || (tt.id == "" && tt.repo.getPostID != "") {
				t.Fatalf("unexpected post, error, or repository call: %+v, %v, %q", post, err, tt.repo.getPostID)
			}
		})
	}
}

func TestPostService_GetByAuthor(t *testing.T) {
	wantErr := errors.New("get by author failed")
	for _, tt := range []struct {
		name   string
		repo   *stubPostRepo
		author string
		want   error
	}{
		{name: "returns posts", repo: &stubPostRepo{posts: []Post{{Id: "post-id"}}}, author: "author-id"},
		{name: "rejects empty author", repo: &stubPostRepo{}, want: ErrAuthorIDRequired},
		{name: "returns repository error", repo: &stubPostRepo{getByAuthorErr: wantErr}, author: "author-id", want: wantErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			posts, err := NewPostService(tt.repo).GetByAuthor(tt.author, "title")
			if !errors.Is(err, tt.want) || (tt.want == nil && len(posts) != 1) || (tt.author == "" && tt.repo.getAuthorID != "") || (tt.author != "" && tt.repo.getAuthorTitle != "title") {
				t.Fatalf("unexpected posts, error, or repository call: %+v, %v, %q, %q", posts, err, tt.repo.getAuthorID, tt.repo.getAuthorTitle)
			}
		})
	}
}

func TestPostService_GetByFollowers(t *testing.T) {
	wantErr := errors.New("get feed failed")
	for _, tt := range []struct {
		name string
		repo *stubPostRepo
		id   string
		want error
	}{
		{name: "returns posts", repo: &stubPostRepo{posts: []Post{{Id: "post-id"}}}, id: "user-id"},
		{name: "rejects empty user", repo: &stubPostRepo{}, want: ErrInvalidUserID},
		{name: "returns repository error", repo: &stubPostRepo{getFromFollowersErr: wantErr}, id: "user-id", want: wantErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			posts, err := NewPostService(tt.repo).GetByFollowers(tt.id)
			if !errors.Is(err, tt.want) || (tt.want == nil && len(posts) != 1) || (tt.id == "" && tt.repo.getFollowerUserID != "") {
				t.Fatalf("unexpected posts, error, or repository call: %+v, %v, %q", posts, err, tt.repo.getFollowerUserID)
			}
		})
	}
}

func TestPostService_UpdatePost(t *testing.T) {
	wantErr := errors.New("update failed")
	for _, tt := range []struct {
		name string
		repo *stubPostRepo
		post Post
		want error
	}{
		{name: "updates post", repo: &stubPostRepo{}, post: Post{Id: "post-id", Title: "Title", Description: "Description"}},
		{name: "rejects missing fields", repo: &stubPostRepo{}, post: Post{}, want: ErrPostTitleDescriptionEmpty},
		{name: "rejects long title", repo: &stubPostRepo{}, post: Post{Title: strings.Repeat("a", 51), Description: "Description"}, want: ErrPostTitleTooLong},
		{name: "rejects long description", repo: &stubPostRepo{}, post: Post{Title: "Title", Description: strings.Repeat("a", 2001)}, want: ErrPostDescriptionTooLong},
		{name: "returns repository error", repo: &stubPostRepo{updateErr: wantErr}, post: Post{Id: "post-id", Title: "Title", Description: "Description"}, want: wantErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := NewPostService(tt.repo).UpdatePost(tt.post)
			if !errors.Is(err, tt.want) || (tt.want == nil && tt.repo.updatedPost != tt.post) || (tt.want != nil && tt.name != "returns repository error" && tt.repo.updatedPost != (Post{})) {
				t.Fatalf("unexpected error or repository call: %v, %+v", err, tt.repo.updatedPost)
			}
		})
	}
}

func TestPostService_DeletePost(t *testing.T) {
	wantErr := errors.New("delete failed")
	for _, tt := range []struct {
		name string
		repo *stubPostRepo
		id   string
		want error
	}{
		{name: "deletes post", repo: &stubPostRepo{}, id: "post-id"},
		{name: "rejects empty id", repo: &stubPostRepo{}, want: ErrPostIDRequiredForDelete},
		{name: "returns repository error", repo: &stubPostRepo{deleteErr: wantErr}, id: "post-id", want: wantErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := NewPostService(tt.repo).DeletePost(tt.id)
			if !errors.Is(err, tt.want) || (tt.want == nil && tt.repo.deletedPostID != tt.id) || (tt.id == "" && tt.repo.deletedPostID != "") {
				t.Fatalf("unexpected error or repository call: %v, %q", err, tt.repo.deletedPostID)
			}
		})
	}
}
