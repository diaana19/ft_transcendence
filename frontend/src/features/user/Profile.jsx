/*
** File: Profile.jsx
** Description: User profile page that displays user information and allows editing or deletion
** Responsibilities:
** - Fetch user data by ID from API
** - Display user profile information
** - Provide edit functionality via modal
** - Provide account deletion functionality via confirmation modal
*/

import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";
import FriendsList from './FriendsList'
import FollowButton from "../../components/common/FollowButton";
import PostCard from '../posts/PostCard';

export default function Profile() {
	const { user: authUser, token, logout } = useAuth();
	const { id } = useParams();
  const userId = id || authUser?.userId;
	const navigate = useNavigate();

	const [user, setUser] = useState({
	name: "",
	email: "",
	bio: "",
	});

	const [loading, setLoading] = useState(false);
	const [showEdit, setShowEdit] = useState(false);
	const [showDelete, setShowDelete] = useState(false);
	const [form, setForm] = useState(user);
  const [posts, setPosts] = useState([]);
  const [loadingPosts, setLoadingPosts] = useState(false);
  const [isFollowing, setIsFollowing] = useState(false);
  const [uploadingAvatar, setUploadingAvatar] = useState(false)
  const [uploadingBanner, setUploadingBanner] = useState(false)

	useEffect(() => {
	if (userId) {
	fetchUser();
  fetchUserPosts();
	}
	}, [userId]);

	/*
	** Fetches user data from the API and updates local state
	** params:
	**   none (uses userId from route params)
	** returns:
	**   Promise<void>
	*/
	const fetchUser = async () => {
	try {
	setLoading(true);

	const res = await fetch(`/api/users/${userId}`, {
    headers: { Authorization: `Bearer ${token}` },
    });
	const data = await res.json();

	setUser(data);
	setForm(data);
	} catch (err) {
	console.error(err);
	} finally {
	setLoading(false);
	}
	};

	/*
	** Fetches user posts from the API and updates local state
	** params:
	**   none (uses userId from route params)
	*/

  const fetchUserPosts = async () => {
  try {
    setLoadingPosts(true)
    const res = await fetch(`/api/posts/user/${userId}`)
    const data = await res.json()
    setPosts(data)
  } catch (err) {
    console.error(err)
  } finally {
    setLoadingPosts(false)
  }
}

	/*
	** Updates user profile data on the server
	** params:
	**   none (uses form state and userId from route params)
	** returns:
	**   Promise<void>
	*/
	const handleUpdate = async () => {
	try {
		const res = await fetch(`/api/users/${userId}`, {
		method: "PUT",
		headers: {
			"Content-Type": "application/json",
			Authorization: `Bearer ${token}`
		},
		body: JSON.stringify(form),
		});

		const updatedUser = await res.json();

		setUser({ ...user, ...updatedUser });
		setShowEdit(false);
	} catch (err) {
		console.error(err);
	}
	};

	/*
	** Deletes the user account from the server
	** params:
	**   none (uses userId from route params)
	** returns:
	**   Promise<void>
	*/
	const handleDelete = async () => {
		try {
			await fetch(`/api/users/${userId}`, {
			method: "DELETE",
			headers: { Authorization: `Bearer ${token}` },
			});
			logout();
			navigate("/login");
		} catch (err) {
			console.error(err);
		}
	};

	if (loading) return <p>Cargando perfil...</p>;

return (
  <div className="min-h-screen bg-white">

    {/* Banner */}
    <div className="h-48 w-full overflow-hidden bg-blue-500" >
    {user.banner && (
      <img src={user.banner} alt="banner" className="w-full h-full object-cover" />
    )}
    </div>
    {/* Card */}
    <div className="border-b border-gray-200">
      <div className="relative px-4">

        {/* Avatar */}
        <div className="absolute -top-16">
          <div className="relative w-32 h-32">
            <div className="w-32 h-32 rounded-full border-4 border-white bg-gray-300 overflow-hidden">
              {user.avatar
                ? <img src={user.avatar || '/default-avatar.jpg'} alt="avatar"  className="w-full h-full object-cover"
                    onError={(e) => { e.target.src = '/default-avatar.jpg' }} />
                : <div className="w-full h-full bg-gray-300" />
              }
            </div>
            {authUser?.userId === userId && (
              <label className="absolute bottom-0 right-0 bg-black bg-opacity-60 rounded-full p-2 cursor-pointer hover:bg-opacity-80">
                <svg viewBox="0 0 24 24" className="w-4 h-4 fill-white">
                  <path d="M12 15.2a3.2 3.2 0 100-6.4 3.2 3.2 0 000 6.4z"/>
                  <path d="M9 2L7.17 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2h-3.17L15 2H9zm3 15c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5z"/>
                </svg>
                <input type="file" accept="image/*" className="hidden"
                  onChange={async (e) => {
                    const file = e.target.files[0]
                    if (!file) return
                    {/* Upload Image */}
                    const formData = new FormData()
                    formData.append('file', file)
                    const res = await fetch('/api/upload', {
                      method: 'POST',
                      headers: { Authorization: `Bearer ${token}` },
                      body: formData
                    })
                    const data = await res.json()
                     const updateRes = await fetch(`/api/users/${userId}`, {
                      method: 'PUT',
                      headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${token}`
                      },
                      body: JSON.stringify({ ...user, avatar: data.path })
                    })
                    const updatedUser = await updateRes.json()

                    setUser({ ...user, ...updatedUser , avatar: data.path })
                  }}
                />
              </label>
            )}
          </div>
        </div>

        {/* Botones top-right */}
        <div className="flex justify-end pt-3 pb-2 gap-2">
          {authUser?.userId === userId ? (
            <>
              <button
                onClick={() => setShowEdit(true)}
                className="border border-gray-300 px-4 py-1.5 rounded-full text-sm font-bold hover:bg-gray-100 transition-colors"
              >
                Edit profile
              </button>
              <button
                onClick={() => setShowDelete(true)}
                className="border border-red-300 px-4 py-1.5 rounded-full text-sm font-bold text-red-500 hover:bg-red-50 transition-colors"
              >
                Delete
              </button>
            </>
          ) : (
            <FollowButton targetId={userId} isFollowing={isFollowing} />
          )}
        </div>

        {/* Info — empieza después del avatar */}
        <div className="mt-16 pb-3">
          <h2 className="text-xl font-bold text-black">{user.name || 'No name'}</h2>
          <p className="text-gray-500 text-sm">@{user.email}</p>
          {user.bio && <p className="mt-2 text-sm text-gray-800">{user.bio}</p>}

          {/* Contadores */}
          <div className="flex gap-4 mt-3">
            <span className="text-sm cursor-pointer hover:underline">
              <strong className="text-black">{user.following_count ?? 0}</strong>
              <span className="text-gray-500 ml-1">Following</span>
            </span>
            <span className="text-sm cursor-pointer hover:underline">
              <strong className="text-black">{user.followers_count ?? 0}</strong>
              <span className="text-gray-500 ml-1">Followers</span>
            </span>
          </div>
        </div>
      </div>

      {/* Posts tab — como Twitter */}
      <div className="flex border-b border-gray-200">
        <div className="flex-1 text-center py-4 text-sm font-bold border-b-2 border-blue-400 text-black cursor-pointer">
          Posts
        </div>
      </div>
    </div>

    {/* Posts list */}
    <div>
      {loadingPosts ? (
        <p className="text-center py-8 text-gray-400">Loading posts...</p>
      ) : posts.length === 0 ? (
        <p className="text-center py-8 text-gray-400">No posts yet</p>
      ) : (
        posts.map((post) => (
          <PostCard
            key={post.id}
            post={post}
            currentUserId={authUser?.userId}
            onDelete={(id) => setPosts(posts.filter(p => p.id !== id))}
            onUpdate={(updated) => setPosts(posts.map(p => p.id === updated.id ? updated : p))}
          />
        ))
      )}
    </div>

    {/* Modal Edit */}
    {showEdit && (
      <div className="fixed inset-0 z-50 flex items-center justify-center">
        {/* Overlay */}
      <div className="absolute inset-0 " onClick={() => setShowEdit(false)} />

    {/* Modal*/}
      <div className="relative z-10 bg-white rounded-2xl w-full max-w-lg shadow-xl overflow-hidden">

          {/* Header */}
          <div className="flex items-center gap-4 px-4 py-3 border-b border-gray-200">
            <button onClick={() => setShowEdit(false)} className="hover:bg-gray-100 rounded-full p-2 text-sm">✕</button>
            <h3 className="text-lg font-bold flex-1">Edit profile</h3>
            <button
              onClick={handleUpdate}
              className="bg-black text-white text-sm font-bold px-4 py-1.5 rounded-full hover:bg-gray-800"
            >
              Save
            </button>
          </div>

          {/* Banner clickable */}
          <div className="relative h-32 bg-blue-500 ">
            {form.banner
              ? <img src={form.banner} alt="banner" className="w-full h-full object-cover" />
              : null
            }
            <label className="absolute inset-0 flex items-center justify-center backdrop-blur-sm bg-black/10 cursor-pointer hover:bg-opacity-40 transition-all">
              {uploadingBanner
                ? <span className="text-white text-xs">Uploading...</span>
                : <svg viewBox="0 0 24 24" className="w-6 h-6 fill-white">
                    <path d="M12 15.2a3.2 3.2 0 100-6.4 3.2 3.2 0 000 6.4z"/>
                    <path d="M9 2L7.17 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2h-3.17L15 2H9zm3 15c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5z"/>
                  </svg>
              }
              <input type="file" accept="image/*" className="hidden"
                onChange={async (e) => {
                  const file = e.target.files[0]
                  if (!file) return
                  setUploadingBanner(true)
                  try {
                    // POST — sube la imagen
                    const formData = new FormData()
                    formData.append('file', file)
                    const res = await fetch('/api/upload', {
                      method: 'POST',
                      headers: { Authorization: `Bearer ${token}` },
                      body: formData
                    })
                    const data = await res.json()
                    // Guarda en form para que el Save lo mande con PUT
                    setForm(prev => ({ ...prev, banner: data.path }))
                  } catch (err) {
                    console.error(err)
                  } finally {
                    setUploadingBanner(false)
                  }
                }}
              />
            </label>
          </div>

          {/* Avatar clickable */}
          <div className="px-4">
            <div className="relative -mt-10 mb-4 w-20 h-20">
              <div className="w-20 h-20 rounded-full border-4 border-white bg-gray-300 overflow-hidden">
                <img
                  src={form.avatar || '/default-avatar.jpg'}
                  alt="avatar"
                  className="w-full h-full object-cover"
                  onError={(e) => { e.target.src = '/default-avatar.jpg' }}
                />
              </div>
                <label className="absolute inset-0 flex items-center justify-center backdrop-blur-sm bg-black/10 rounded-full cursor-pointer hover:bg-opacity-50 transition-all">
                  {uploadingAvatar
                    ? <span className="text-white text-xs">...</span>
                    : <svg viewBox="0 0 24 24" className="w-5 h-5 fill-white">
                        <path d="M12 15.2a3.2 3.2 0 100-6.4 3.2 3.2 0 000 6.4z"/>
                        <path d="M9 2L7.17 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2h-3.17L15 2H9zm3 15c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5z"/>
                      </svg>
                  }
                  <input type="file" accept="image/*" className="hidden"
                    onChange={async (e) => {
                      const file = e.target.files[0]
                      if (!file) return
                      setUploadingAvatar(true)
                      try {
                        // POST — sube la imagen
                        const formData = new FormData()
                        formData.append('file', file)
                        const res = await fetch('/api/upload', {
                          method: 'POST',
                          headers: { Authorization: `Bearer ${token}` },
                          body: formData
                        })
                        const data = await res.json()
                        // Guarda en form — el Save manda todo con PUT
                        setForm(prev => ({ ...prev, avatar: data.path }))
                      } catch (err) {
                        console.error(err)
                      } finally {
                        setUploadingAvatar(false)
                      }
                    }}
                  />
                </label>
            </div>

            {/* Campos */}
            <input
              className="w-full border border-gray-300 rounded-lg px-3 py-2 mb-3 text-sm focus:outline-none focus:border-blue-400"
              placeholder="Name"
              type="text"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
            />
            <input
              className="w-full border border-gray-300 rounded-lg px-3 py-2 mb-3 text-sm focus:outline-none focus:border-blue-400"
              placeholder="Email"
              type="email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
            />
            <textarea
              className="w-full border border-gray-300 rounded-lg px-3 py-2 mb-4 text-sm focus:outline-none focus:border-blue-400 resize-none"
              placeholder="Bio"
              rows={3}
              value={form.bio}
              onChange={(e) => setForm({ ...form, bio: e.target.value })}
            />
          </div>
        </div>
      </div>
    )}

    {/* Modal Delete */}
    {showDelete && (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
        <div className="bg-white p-6 rounded-2xl w-80 shadow-xl">
          <h3 className="text-lg font-bold mb-2">Delete account</h3>
          <p className="text-gray-500 text-sm mb-4">This action is permanent and cannot be undone.</p>
          <div className="flex justify-end gap-2">
            <button onClick={() => setShowDelete(false)} className="px-4 py-2 rounded-full border border-gray-300 text-sm font-semibold hover:bg-gray-100">
              Cancel
            </button>
            <button onClick={handleDelete} className="px-4 py-2 rounded-full bg-red-500 text-white text-sm font-bold hover:bg-red-600">
              Yes, delete
            </button>
          </div>
        </div>
      </div>
    )}
  </div>
)
}