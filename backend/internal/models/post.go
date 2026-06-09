package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Post is a publication made by a user.
type Post struct {
	ID            string         `gorm:"primaryKey;type:uuid"`
	AuthorID      string         `gorm:"type:uuid;not null"`
	Author        User           `gorm:"foreignKey:AuthorID;references:ID"`
	Content       string         `gorm:"type:text;not null"`
	MediaURL      *string        `gorm:"type:text" json:"media_url,omitempty"`
	MediaMIME     *string        `gorm:"type:varchar(100)" json:"media_mime,omitempty"`
	Tags          pq.StringArray `gorm:"type:text[];index:idx_posts_tags,type:gin" json:"tags,omitempty"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	LikesCount    int            `json:"likes_count" gorm:"default:0"`
	DislikesCount int            `json:"dislikes_count" gorm:"default:0"`
	CommentsCount int            `json:"comments_count" gorm:"default:0"`
	Comments      []Reply        `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}

// TagCount holds a tag and how many times it is used.
type TagCount struct {
	Tag   string `json:"tag"`
	Count int64  `json:"count"`
}

// UpdatePostInput holds the fields a user can change in a post.
type UpdatePostInput struct {
	Content   string  `json:"content" binding:"required,min=1,max=280"`
	MediaURL  *string `json:"media_url,omitempty"`
	MediaMIME *string `json:"media_mime,omitempty"`
}

// PostResponse is the post data we send to the client.
type PostResponse struct {
	ID            string       `json:"id"`
	Content       string       `json:"content"`
	MediaURL      *string      `json:"media_url,omitempty"`
	MediaMIME     *string      `json:"media_mime,omitempty"`
	AuthorID      string       `json:"author_id"`
	Author        UserResponse `json:"author"`
	LikesCount    int          `json:"likes_count"`
	DislikesCount int          `json:"dislikes_count"`
	CommentsCount int          `json:"comments_count"`
	Liked         bool         `json:"liked"`
	Disliked      bool         `json:"disliked"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// ToResponse converts the Post into a PostResponse.
func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:            p.ID,
		Content:       p.Content,
		MediaURL:      p.MediaURL,
		MediaMIME:     p.MediaMIME,
		AuthorID:      p.AuthorID,
		Author:        p.Author.ToResponse(),
		LikesCount:    p.LikesCount,
		DislikesCount: p.DislikesCount,
		CommentsCount: p.CommentsCount,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// PostReaction is a like or dislike a user made on a post.
type PostReaction struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	UserID    string `gorm:"type:uuid;not null;uniqueIndex:idx_like_user_post"`
	User      User   `gorm:"foreignKey:UserID;references:ID"`
	PostID    string `gorm:"type:uuid;not null;uniqueIndex:idx_like_user_post"`
	Post      Post   `gorm:"foreignKey:PostID;references:ID"`
	Value     int    `gorm:"not null;default:1"`
	CreatedAt time.Time
}

// TableName returns the database table name for PostReaction.
func (PostReaction) TableName() string { return "likes" }

// ReplyReaction is a like or dislike a user made on a reply.
type ReplyReaction struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	UserID    string `gorm:"type:uuid;not null;uniqueIndex:idx_reply_reaction_user_reply"`
	User      User   `gorm:"foreignKey:UserID;references:ID"`
	ReplyID   string `gorm:"type:uuid;not null;uniqueIndex:idx_reply_reaction_user_reply"`
	Reply     Reply  `gorm:"foreignKey:ReplyID;references:ID"`
	Value     int    `gorm:"not null;default:1"`
	CreatedAt time.Time
}

// ReactionResponse is the reaction data we send to the client.
type ReactionResponse struct {
	PostID        string `json:"post_id,omitempty"`
	CommentID     string `json:"comment_id,omitempty"`
	UserReaction  int    `json:"user_reaction"`
	LikesCount    int    `json:"likes_count"`
	DislikesCount int    `json:"dislikes_count"`
}

// Reply is a comment made on a post.
type Reply struct {
	ID            string  `gorm:"primaryKey;type:uuid"`
	PostID        string  `gorm:"type:uuid;not null;index"`
	Post          Post    `gorm:"foreignKey:PostID;references:ID"`
	AuthorID      string  `gorm:"type:uuid;not null"`
	Author        User    `gorm:"foreignKey:AuthorID;references:ID"`
	Content       string  `gorm:"type:text;not null"`
	FileID        *string `gorm:"type:uuid"`
	File          *File   `gorm:"foreignKey:FileID;references:ID"`
	LikesCount    int     `json:"likes_count" gorm:"default:0"`
	DislikesCount int     `json:"dislikes_count" gorm:"default:0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// CreateCommentInput holds the data to create a comment.
type CreateCommentInput struct {
	Content string `form:"content" binding:"required,min=1,max=280"`
}

// UpdateCommentInput holds the data to update a comment.
type UpdateCommentInput struct {
	Content string `json:"content" form:"content" binding:"required,min=1,max=280"`
	// NewFileID, when set, replaces the comment's attachment. It is populated by
	// the controller after uploading the new file, never bound from the request.
	NewFileID *string `json:"-"`
	// RemoveFile detaches the current attachment when true. Ignored if a new
	// file is provided (a replacement implies keeping an attachment).
	RemoveFile bool `json:"remove_file" form:"remove_file"`
}

// CommentResponse is the comment data we send to the client.
type CommentResponse struct {
	ID            string       `json:"id"`
	PostID        string       `json:"post_id"`
	Content       string       `json:"content"`
	AuthorID      string       `json:"author_id"`
	Author        UserResponse `json:"author"`
	LikesCount    int          `json:"likes_count"`
	DislikesCount int          `json:"dislikes_count"`
	Liked         bool         `json:"liked"`
	Disliked      bool         `json:"disliked"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	FileID        *string      `json:"file_id,omitempty"`
	FileURL       *string      `json:"file_url,omitempty"`
	FileMIME      *string      `json:"file_mime,omitempty"`
}

// ToResponse converts the Reply into a CommentResponse.
func (r *Reply) ToResponse() CommentResponse {
	resp := CommentResponse{
		ID:            r.ID,
		PostID:        r.PostID,
		Content:       r.Content,
		AuthorID:      r.AuthorID,
		Author:        r.Author.ToResponse(),
		LikesCount:    r.LikesCount,
		DislikesCount: r.DislikesCount,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}

	if r.FileID != nil {
		resp.FileID = r.FileID
		url := "/api/files/" + *r.FileID
		resp.FileURL = &url
		if r.File != nil {
			resp.FileMIME = &r.File.MimeType
		}
	}

	return resp
}

// Repost is a post shared again by another user.
type Repost struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	PostID    string `gorm:"type:uuid;not null;index"`
	Post      Post   `gorm:"foreignKey:PostID;references:ID"`
	AuthorID  string `gorm:"type:uuid;not null"`
	Author    User   `gorm:"foreignKey:AuthorID;references:ID"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// RepostResponse is the repost data we send to the client.
type RepostResponse struct {
	ID        string       `json:"id"`
	PostID    string       `json:"post_id"`
	AuthorID  string       `json:"author_id"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
}

// ToResponse converts the Repost into a RepostResponse.
func (r *Repost) ToResponse() RepostResponse {
	return RepostResponse{
		ID:        r.ID,
		PostID:    r.PostID,
		AuthorID:  r.AuthorID,
		Author:    r.Author.ToResponse(),
		CreatedAt: r.CreatedAt,
	}
}
