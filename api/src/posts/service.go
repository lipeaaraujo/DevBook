package posts

import "api/src/apierrors"

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

func (service PostService) GetPosts(title string) ([]Post, error) {
	posts, err := service.repo.Get(title)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (service PostService) GetById(postId string) (Post, error) {
	if (postId == "") {
		return Post{}, apierrors.BadRequest("You need to pass the postId")
	}

	post, err := service.repo.GetById(postId)
	if err != nil {
		return Post{}, err
	}

	return post, err
}

func (service PostService) GetByAuthor(authorId string, title string) ([]Post, error) {
	if (authorId == "") {
		return nil, apierrors.BadRequest("You need to pass the authorId")
	}

	posts, err := service.repo.GetByAuthor(authorId, title)
	if err != nil {
		return nil, err
	}

	return posts, err
}

func (service PostService) UpdatePost(post *Post) (error) {
	if post.Title == "" && post.Description == "" {
		return apierrors.BadRequest("You need to pass the title or Description")
	}

	if err := post.PrepareUpdate(); err != nil {
		return err
	}

	if err := service.repo.Update(post); err != nil {
		return err
	}

	return nil
}

func (service PostService) DeletePost(postId string) (error) {
	if postId != "" {
		return apierrors.BadRequest("You have to pass the postId")
	}

	if err := service.repo.Delete(postId); err != nil {
		return err
	}

	return nil
}
