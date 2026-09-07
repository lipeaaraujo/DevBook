package posts

type PostRepoInterface interface {
	Create(post Post) (string, error)
	Get(title string) ([]Post, error)
	GetById(id string) (Post, error)
	GetByAuthor(authorId string, title string) ([]Post, error)
	GetFromFollowers(userId string) ([]Post, error)
	Update(post Post) error
	Delete(id string) error
}

type PostService struct {
	repo PostRepoInterface
}

func NewPostService(repo PostRepoInterface) *PostService {
	return &PostService{repo: repo}
}

func (service PostService) CreatePost(post Post) (string, error) {
	if err := post.PrepareCreate(); err != nil {
		return "", err
	}

	createdId, err := service.repo.Create(post)
	if err != nil {
		return "", err
	}

	return createdId, nil
}

func (service PostService) GetPosts(title string) ([]Post, error) {
	posts, err := service.repo.Get(title)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (service PostService) GetById(postId string) (Post, error) {
	if postId == "" {
		return Post{}, ErrPostIDRequired
	}

	post, err := service.repo.GetById(postId)
	if err != nil {
		return Post{}, err
	}

	return post, err
}

func (service PostService) GetByAuthor(authorId string, title string) ([]Post, error) {
	if authorId == "" {
		return nil, ErrAuthorIDRequired
	}

	posts, err := service.repo.GetByAuthor(authorId, title)
	if err != nil {
		return nil, err
	}

	return posts, err
}

func (service PostService) GetByFollowers(userId string) ([]Post, error) {
	if userId == "" {
		return nil, ErrInvalidUserID
	}

	posts, err := service.repo.GetFromFollowers(userId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (service PostService) UpdatePost(post Post) error {
	if err := post.PrepareUpdate(); err != nil {
		return err
	}

	if err := service.repo.Update(post); err != nil {
		return err
	}

	return nil
}

func (service PostService) DeletePost(postId string) error {
	if postId == "" {
		return ErrPostIDRequiredForDelete
	}

	if err := service.repo.Delete(postId); err != nil {
		return err
	}

	return nil
}
