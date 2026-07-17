package posts

type PostService struct{
	repo *PostRepo
}

func NewPostService(repo *PostRepo) *PostService {
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
