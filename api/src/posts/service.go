package posts

type PostService struct{
	repo *PostRepo
}

func NewPostService(repo *PostRepo) *PostService {
	return &PostService{repo: repo}
}
