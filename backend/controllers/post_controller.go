package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Transcendence/models"
	"github.com/Transcendence/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController struct {
	postService         *services.PostService
	notificationService *services.NotificationService
}

func NewPostController(postService *services.PostService, notifService *services.NotificationService) *PostController {
	return &PostController{postService: postService, notificationService: notifService}
}

func (pc *PostController) GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	posts, total, err := pc.postService.GetPosts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	responses := make([]models.PostResponse, len(posts))
	for i, p := range posts {
		resp := p.ToResponse()
		if userID != nil {
			liked, _ := pc.postService.HasLiked(userID.(string), p.ID)
			resp.Liked = liked
		}
		responses[i] = resp
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  responses,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (pc *PostController) GetPost(c *gin.Context) {
	id := c.Param("id")
	post, err := pc.postService.GetPost(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := post.ToResponse()
	if userID, exists := c.Get("user_id"); exists {
		liked, _ := pc.postService.HasLiked(userID.(string), id)
		resp.Liked = liked
	}

	c.JSON(http.StatusOK, resp)
}

func (pc *PostController) CreatePost(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
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

	post, err := pc.postService.CreatePost(req.Content, authorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post.ToResponse())
}

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
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		case err.Error() == "you can only update your own posts":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, post.ToResponse())
}

func (pc *PostController) DeletePost(c *gin.Context) {
	id := c.Param("id")

	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := pc.postService.DeletePost(id, authorID.(string))
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		case err.Error() == "you can only delete your own posts":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}






func (pc *PostController) ToggleLike(c *gin.Context) {
	postID := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	liked, post, err := pc.postService.ToggleLike(userID.(string), postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if liked && post.AuthorID != userID.(string) {
		username, _ := c.Get("username")
		pc.notificationService.SendNotification(
			post.AuthorID,
			post.Author.Username,
			userID.(string),
			username.(string),
			"like",
			username.(string)+" liked your post",
		)
	}

	c.JSON(http.StatusOK, models.LikeResponse{
		PostID:     postID,
		Liked:      liked,
		LikesCount: post.LikesCount,
	})
}




func (pc *PostController) GetComments(c *gin.Context) {
	postID := c.Param("id")

	comments, err := pc.postService.GetComments(postID)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.CommentResponse, len(comments))
	for i, cm := range comments {
		responses[i] = cm.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"data": responses, "total": len(responses)})
}


func (pc *PostController) CreateComment(c *gin.Context) {
	postID := c.Param("id")

	var input models.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	post, err := pc.postService.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comment, err := pc.postService.CreateComment(input.Content, authorID.(string), postID)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if post.AuthorID != authorID.(string) {
		username, _ := c.Get("username")
		pc.notificationService.SendNotification(
			post.AuthorID,
			post.Author.Username,
			authorID.(string),
			username.(string),
			"comment",
			username.(string)+" commented on your post",
		)
	}
	
	c.JSON(http.StatusCreated, comment.ToResponse())
}


func (pc *PostController) UpdateComment(c *gin.Context) {
	commentID := c.Param("commentId")

	var input models.UpdateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	comment, err := pc.postService.UpdateComment(commentID, input, authorID.(string))
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		case err.Error() == "you can only update your own comments":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, comment.ToResponse())
}


func (pc *PostController) DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")

	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := pc.postService.DeleteComment(commentID, authorID.(string))
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		case err.Error() == "you can only delete your own comments":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
