import { useParams } from "react-router-dom"
import CommentForm from "./CommentForm"

export default function PostPage() {
  const { id } = useParams()

  return (
    <div>
      {/* Post details aquí */}

      <div id="comments">
        <CommentForm postId={id} />
      </div>
    </div>
  )
}