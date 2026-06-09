package controllers

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

// PostController handles the post and comment endpoints.
type PostController struct {
	postService         *services.PostService
	notificationService *services.NotificationService
	uploadService       *services.UploadService
}

// NewPostController creates a new PostController with its services.
func NewPostController(
	postService *services.PostService,
	notifService *services.NotificationService,
	uploadService *services.UploadService,
) *PostController {
	return &PostController{
		postService:         postService,
		notificationService: notifService,
		uploadService:       uploadService,
	}
}

// respondOwnedResourceError sends the right HTTP status for an owned resource error.
func respondOwnedResourceError(c *gin.Context, err error, notFoundMsg string) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundMsg})
	case strings.HasPrefix(err.Error(), "you can only"):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// GetTrends returns the most used hashtags in the posts of the last week.
// @Summary   Trending hashtags (most-used in posts over the last week)
// @Tags      posts
// @Produce   json
// @Param     limit query int false "max tags to return (default 10, max 50)"
// @Success   200 {object} map[string]interface{}
// @Failure   500 {object} map[string]string
// @Router    /trends [get]
func (pc *PostController) GetTrends(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	trends, err := pc.postService.GetTrends(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if trends == nil {
		trends = []models.TagCount{}
	}

	c.JSON(http.StatusOK, gin.H{"data": trends})
}

// parseLimitOffset reads the limit and offset query params with safe defaults.
func parseLimitOffset(c *gin.Context) (limit, offset int) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	return limit, offset
}

// GetPosts lists posts with pagination, optionally filtered by hashtag or replied-to.
// @Summary   List posts (limit/offset pagination), optionally filtered by hashtag or replied-to
// @Tags      posts
// @Produce   json
// @Param     tag       query string false "filter to posts carrying this hashtag (with or without #)"
// @Param     repliedBy query string false "list only posts this user id has replied to"
// @Param     limit     query int    false "max items to return (default 10, max 50)"
// @Param     offset    query int    false "items to skip (default 0)"
// @Success   200 {object} map[string]interface{}
// @Failure   500 {object} map[string]string
// @Router    /posts [get]
func (pc *PostController) GetPosts(c *gin.Context) {
	limit, offset := parseLimitOffset(c)
	userID, _ := c.Get("user_id")

	var (
		posts []models.Post
		total int64
		err   error
	)
	switch {
	case c.Query("repliedBy") != "":
		posts, total, err = pc.postService.GetRepliedPosts(c.Query("repliedBy"), limit, offset)
	case utils.NormalizeHashtag(c.Query("tag")) != "":
		posts, total, err = pc.postService.GetPostsByTag(utils.NormalizeHashtag(c.Query("tag")), limit, offset)
	default:
		posts, total, err = pc.postService.GetPosts(limit, offset)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.PostResponse, len(posts))
	for i, p := range posts {
		resp := p.ToResponse()
		if userID != nil {
			value, _ := pc.postService.GetPostReaction(userID.(string), p.ID)
			resp.Liked = value == 1
			resp.Disliked = value == -1
		}
		responses[i] = resp
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   responses,
		"limit":  limit,
		"offset": offset,
		"total":  total,
	})
}

// GetPost returns one post by its id.
// @Summary   Get a single post by id
// @Tags      posts
// @Produce   json
// @Param     id path string true "post id"
// @Success   200 {object} models.PostResponse
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /posts/{id} [get]
func (pc *PostController) GetPost(c *gin.Context) {
	id := c.Param("id")
	post, err := pc.postService.GetPost(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": msgPostNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := post.ToResponse()
	if userID, exists := c.Get("user_id"); exists {
		value, _ := pc.postService.GetPostReaction(userID.(string), id)
		resp.Liked = value == 1
		resp.Disliked = value == -1
	}

	c.JSON(http.StatusOK, resp)
}

// GetPostsByUser lists the posts written by the given user.
// @Summary   List posts authored by a user
// @Tags      posts
// @Produce   json
// @Param     userId path string true "user id"
// @Success   200 {object} map[string]interface{}
// @Failure   500 {object} map[string]string
// @Router    /posts/user/{userId} [get]
func (pc *PostController) GetPostsByUser(c *gin.Context) {
	userID := c.Param("userId")

	posts, err := pc.postService.GetPostsByAuthor(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	currentUserID, _ := c.Get("user_id")
	responses := make([]models.PostResponse, len(posts))
	for i, p := range posts {
		resp := p.ToResponse()
		if currentUserID != nil {
			value, _ := pc.postService.GetPostReaction(currentUserID.(string), p.ID)
			resp.Liked = value == 1
			resp.Disliked = value == -1
		}
		responses[i] = resp
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// CreatePost creates a new post for the logged in user.
// @Summary   Create a new post
// @Tags      posts
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     body body object true "post content and optional media_url"
// @Success   201 {object} models.PostResponse
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Router    /posts [post]
func (pc *PostController) CreatePost(c *gin.Context) {
	var req struct {
		Content   string  `json:"content" binding:"required"`
		MediaURL  *string `json:"media_url"`
		MediaMIME *string `json:"media_mime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	post, err := pc.postService.CreatePost(req.Content, authorID.(string), req.MediaURL, req.MediaMIME)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAuthorNotFound):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "your session is no longer valid, please log in again"})
		case strings.Contains(err.Error(), "content"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.Printf("CreatePost failed for author %s: %v", authorID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create post"})
		}
		return
	}

	c.JSON(http.StatusCreated, post.ToResponse())
}

// UpdatePost updates a post. Only the author can update it.
// @Summary   Update a post
// @Tags      posts
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     id   path string                  true "post id"
// @Param     body body models.UpdatePostInput true "post fields to update"
// @Success   200 {object} models.PostResponse
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   403 {object} map[string]string
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /posts/{id} [put]
func (pc *PostController) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var input models.UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	post, err := pc.postService.UpdatePost(id, input, authorID.(string))
	if err != nil {
		respondOwnedResourceError(c, err, msgPostNotFound)
		return
	}

	c.JSON(http.StatusOK, post.ToResponse())
}

// DeletePost deletes a post. Only the author can delete it.
// @Summary   Delete a post
// @Tags      posts
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "post id"
// @Success   200 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   403 {object} map[string]string
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /posts/{id} [delete]
func (pc *PostController) DeletePost(c *gin.Context) {
	id := c.Param("id")

	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := pc.postService.DeletePost(id, authorID.(string))
	if err != nil {
		respondOwnedResourceError(c, err, msgPostNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// reactionInput holds the reaction value sent to like or dislike.
type reactionInput struct {
	Value int `json:"value" binding:"required"`
}

// bindReactionValue reads and validates the reaction value (1 or -1) from the body.
func bindReactionValue(c *gin.Context) (int, bool) {
	var input reactionInput
	if err := c.ShouldBindJSON(&input); err != nil || (input.Value != 1 && input.Value != -1) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "value must be 1 (like) or -1 (dislike)"})
		return 0, false
	}
	return input.Value, true
}

// React likes or dislikes a post and notifies the author on a like.
// @Summary   Like or dislike a post
// @Tags      posts
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     id   path string true "post id"
// @Param     body body controllers.reactionInput true "reaction value: 1 like, -1 dislike"
// @Success   200 {object} models.ReactionResponse
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /posts/{id}/react [post]
func (pc *PostController) React(c *gin.Context) {
	postID := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pressed, ok := bindReactionValue(c)
	if !ok {
		return
	}

	value, post, err := pc.postService.ReactToPost(userID.(string), postID, pressed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": msgPostNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if value == 1 && post.AuthorID != userID.(string) {
		username, _ := c.Get("username")
		_ = pc.notificationService.SendNotification(
			post.AuthorID,
			post.Author.Username,
			userID.(string),
			username.(string),
			"like",
			username.(string)+" liked your post",
			postID,
		)
	}

	c.JSON(http.StatusOK, models.ReactionResponse{
		PostID:        postID,
		UserReaction:  value,
		LikesCount:    post.LikesCount,
		DislikesCount: post.DislikesCount,
	})
}

// ReactComment likes or dislikes a comment.
// @Summary   Like or dislike a comment
// @Tags      posts
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     id        path string true "post id"
// @Param     commentId path string true "comment id"
// @Param     body      body controllers.reactionInput true "reaction value: 1 like, -1 dislike"
// @Success   200 {object} models.ReactionResponse
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /posts/{id}/comments/{commentId}/react [post]
func (pc *PostController) ReactComment(c *gin.Context) {
	commentID := c.Param("commentId")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pressed, ok := bindReactionValue(c)
	if !ok {
		return
	}

	value, comment, err := pc.postService.ReactToComment(userID.(string), commentID, pressed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ReactionResponse{
		CommentID:     commentID,
		UserReaction:  value,
		LikesCount:    comment.LikesCount,
		DislikesCount: comment.DislikesCount,
	})
}

// GetComments lists the comments of a post.
// @Summary   List comments for a post
// @Tags      posts
// @Produce   json
// @Param     id path string true "post id"
// @Success   200 {object} map[string]interface{}
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /posts/{id}/comments [get]
func (pc *PostController) GetComments(c *gin.Context) {
	postID := c.Param("id")

	comments, err := pc.postService.GetComments(postID)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": msgPostNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	responses := make([]models.CommentResponse, len(comments))
	for i, cm := range comments {
		resp := cm.ToResponse()
		if userID != nil {
			value, _ := pc.postService.GetCommentReaction(userID.(string), cm.ID)
			resp.Liked = value == 1
			resp.Disliked = value == -1
		}
		responses[i] = resp
	}

	c.JSON(http.StatusOK, gin.H{"data": responses, "total": len(responses)})
}

// CreateComment adds a comment to a post, with an optional file, and notifies the author.
// @Summary   Add a comment to a post
// @Tags      posts
// @Security  BearerAuth
// @Accept    multipart/form-data
// @Produce   json
// @Param     id         path     string true  "post id"
// @Param     content    formData string true  "comment content"
// @Param     file       formData file   false "optional image/media attachment"
// @Param     visibility formData string false "file visibility (public/friends/private)"
// @Success   201 {object} models.CommentResponse
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   404 {object} map[string]string
// @Router    /posts/{id}/comments [post]
func (pc *PostController) CreateComment(c *gin.Context) {
	postID := c.Param("id")

	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}
	if len(content) > 280 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content must not exceed 280 characters"})
		return
	}

	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := authorID.(string)

	post, err := pc.postService.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": msgPostNotFound})
		return
	}

	var fileID *string
	fileHeader, err := c.FormFile("file")
	if err == nil && fileHeader != nil {
		visibility := c.DefaultPostForm("visibility", models.FileVisibilityPublic)

		saveFn := func(fh *multipart.FileHeader, dst string) error {
			return c.SaveUploadedFile(fh, dst)
		}

		file, uploadErr := pc.uploadService.SaveFile(fileHeader, userIDStr, visibility, saveFn)
		if uploadErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file upload failed: " + uploadErr.Error()})
			return
		}
		fileID = &file.ID
	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file error: " + err.Error()})
		return
	}
	comment, err := pc.postService.CreateComment(content, userIDStr, postID, fileID)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": msgPostNotFound})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if post.AuthorID != userIDStr {
		username, _ := c.Get("username")
		_ = pc.notificationService.SendCommentNotification(
			post.AuthorID,
			userIDStr,
			username.(string),
			postID,
			content,
		)
	}

	c.JSON(http.StatusCreated, comment.ToResponse())
}

// UpdateComment updates a comment text and optionally its media. Only the author can update it.
// @Summary   Update a comment (text, and optionally its media attachment)
// @Tags      posts
// @Security  BearerAuth
// @Accept    json
// @Accept    multipart/form-data
// @Produce   json
// @Param     id          path     string true  "post id"
// @Param     commentId   path     string true  "comment id"
// @Param     content     formData string false "new content (multipart)"
// @Param     file        formData file   false "new attachment, replaces the current one (multipart)"
// @Param     remove_file formData bool   false "detach the current attachment (multipart or JSON)"
// @Param     body        body     models.UpdateCommentInput false "JSON body for a text-only edit"
// @Success   200 {object} models.CommentResponse
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   403 {object} map[string]string
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /posts/{id}/comments/{commentId} [put]
func (pc *PostController) UpdateComment(c *gin.Context) {
	commentID := c.Param("commentId")

	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := authorID.(string)

	var input models.UpdateCommentInput

	if strings.HasPrefix(c.ContentType(), "application/json") {
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		content := c.PostForm("content")
		if content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
			return
		}
		if len(content) > 280 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content must not exceed 280 characters"})
			return
		}
		input.Content = content
		input.RemoveFile = c.PostForm("remove_file") == "true"

		fileHeader, ferr := c.FormFile("file")
		if ferr == nil && fileHeader != nil {
			visibility := c.DefaultPostForm("visibility", models.FileVisibilityPublic)

			saveFn := func(fh *multipart.FileHeader, dst string) error {
				return c.SaveUploadedFile(fh, dst)
			}

			file, uploadErr := pc.uploadService.SaveFile(fileHeader, userIDStr, visibility, saveFn)
			if uploadErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "file upload failed: " + uploadErr.Error()})
				return
			}
			input.NewFileID = &file.ID
			input.RemoveFile = false
		} else if ferr != http.ErrMissingFile {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file error: " + ferr.Error()})
			return
		}
	}

	comment, err := pc.postService.UpdateComment(commentID, input, userIDStr)
	if err != nil {
		respondOwnedResourceError(c, err, "Comment not found")
		return
	}

	c.JSON(http.StatusOK, comment.ToResponse())
}

// DeleteComment deletes a comment. Only the author can delete it.
// @Summary   Delete a comment
// @Tags      posts
// @Security  BearerAuth
// @Produce   json
// @Param     id        path string true "post id"
// @Param     commentId path string true "comment id"
// @Success   200 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   403 {object} map[string]string
// @Failure   404 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /posts/{id}/comments/{commentId} [delete]
func (pc *PostController) DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")

	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := pc.postService.DeleteComment(commentID, authorID.(string))
	if err != nil {
		respondOwnedResourceError(c, err, "Comment not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
