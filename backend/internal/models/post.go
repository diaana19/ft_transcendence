package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	ID       string  `gorm:"primaryKey;type:uuid"`
	AuthorID string  `gorm:"type:uuid;not null"`
	Author   User    `gorm:"foreignKey:AuthorID;references:ID"`
	Content  string  `gorm:"type:text;not null"`
	MediaURL *string `gorm:"type:text" json:"media_url,omitempty"`
	MediaMIME *string `gorm:"type:varchar(100)" json:"media_mime,omitempty"`
	// Tags holds the distinct lowercased hashtags extracted from Content at
	// write time (see utils.ExtractHashtags). Stored as a Postgres text[] with a
	// GIN index so trends can be aggregated with unnest() without re-scanning
	// content. Kept denormalized alongside Content so it survives a Redis flush
	// and powers real-time trend counts straight from the source of truth.
	Tags          pq.StringArray `gorm:"type:text[];index:idx_posts_tags,type:gin" json:"tags,omitempty"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	LikesCount    int            `json:"likes_count" gorm:"default:0"`
	DislikesCount int            `json:"dislikes_count" gorm:"default:0"`
	CommentsCount int            `json:"comments_count" gorm:"default:0"`
	Comments      []Reply        `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}

// TagCount is one row of the trends aggregation: a hashtag and how many posts
// used it within the queried window.
type TagCount struct {
	Tag   string `json:"tag"`
	Count int64  `json:"count"`
}

type UpdatePostInput struct {
	Content   string  `json:"content" binding:"required,min=1,max=280"`
	MediaURL  *string `json:"media_url,omitempty"`
	MediaMIME *string `json:"media_mime,omitempty"`
}

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

// PostReaction is a user's like (Value=+1) or dislike (Value=-1) on a post. A
// user has at most one reaction per post (enforced by the composite unique
// index), so liking switches an existing dislike and vice-versa. The table is
// still named "likes" so the original like rows are preserved on migration —
// the added value column defaults to 1, marking every legacy like as a like.
type PostReaction struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	UserID    string `gorm:"type:uuid;not null;uniqueIndex:idx_like_user_post"`
	User      User   `gorm:"foreignKey:UserID;references:ID"`
	PostID    string `gorm:"type:uuid;not null;uniqueIndex:idx_like_user_post"`
	Post      Post   `gorm:"foreignKey:PostID;references:ID"`
	Value     int    `gorm:"not null;default:1"`
	CreatedAt time.Time
}

func (PostReaction) TableName() string { return "likes" }

// ReplyReaction is the reply counterpart of PostReaction: a user's like
// (Value=+1) or dislike (Value=-1) on a single reply, unique per user/reply.
type ReplyReaction struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	UserID    string `gorm:"type:uuid;not null;uniqueIndex:idx_reply_reaction_user_reply"`
	User      User   `gorm:"foreignKey:UserID;references:ID"`
	ReplyID   string `gorm:"type:uuid;not null;uniqueIndex:idx_reply_reaction_user_reply"`
	Reply     Reply  `gorm:"foreignKey:ReplyID;references:ID"`
	Value     int    `gorm:"not null;default:1"`
	CreatedAt time.Time
}

// ReactionResponse is returned by the post and comment react endpoints.
// UserReaction is the caller's resulting reaction: +1 liked, -1 disliked, 0
// cleared. Exactly one of PostID / CommentID is set.
type ReactionResponse struct {
	PostID        string `json:"post_id,omitempty"`
	CommentID     string `json:"comment_id,omitempty"`
	UserReaction  int    `json:"user_reaction"`
	LikesCount    int    `json:"likes_count"`
	DislikesCount int    `json:"dislikes_count"`
}

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

type CreateCommentInput struct {
	Content string `form:"content" binding:"required,min=1,max=280"`
}

type UpdateCommentInput struct {
	Content string `json:"content" binding:"required,min=1,max=280"`
}

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

type Repost struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	PostID    string `gorm:"type:uuid;not null;index"`
	Post      Post   `gorm:"foreignKey:PostID;references:ID"`
	AuthorID  string `gorm:"type:uuid;not null"`
	Author    User   `gorm:"foreignKey:AuthorID;references:ID"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type RepostResponse struct {
	ID        string       `json:"id"`
	PostID    string       `json:"post_id"`
	AuthorID  string       `json:"author_id"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
}

func (r *Repost) ToResponse() RepostResponse {
	return RepostResponse{
		ID:        r.ID,
		PostID:    r.PostID,
		AuthorID:  r.AuthorID,
		Author:    r.Author.ToResponse(),
		CreatedAt: r.CreatedAt,
	}
}
