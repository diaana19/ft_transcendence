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
import api from '../../services/axiosInstance'
import FriendsList from './FriendsList'
import FollowButton from "../../components/common/FollowButton";
import PostCard from '../posts/PostCard';
import { getPostsByAuthor } from '../posts/postService'
import { sendFriendRequest } from '../../features/user/userService'
import FollowersModal from './FollowersModal'

export default function Profile() {
  const { user: authUser, updateUser } = useAuth();
  const { id, username } = useParams();
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState("posts");

  // The profile is keyed by user id. The /profile/u/:username route (used by
  // @mention links) carries a username instead, so resolve it to an id first.
  const [resolvedId, setResolvedId] = useState(null);
  const userId = resolvedId;

  const [user, setUser] = useState({ displayname: "", email: "", bio: "" });
  const [loading, setLoading] = useState(false);
  const [showEdit, setShowEdit] = useState(false);
  const [form, setForm] = useState(user);
  const [posts, setPosts] = useState([]);
  const [loadingPosts, setLoadingPosts] = useState(false);
  const [isFollowing, setIsFollowing] = useState(false);
  const [uploadingAvatar, setUploadingAvatar] = useState(false)
  const [uploadingBanner, setUploadingBanner] = useState(false)
  const [friendRequested, setFriendRequested] = useState(false)
  const [isFriend, setIsFriend] = useState(false)
  const [followModal, setFollowModal] = useState(null)

  useEffect(() => {
    let active = true;
    if (username) {
      api.get(`/api/users?username=${encodeURIComponent(username)}`)
        .then(({ data }) => { if (active) setResolvedId(data.id); })
        .catch((err) => console.error(err));
    } else {
      setResolvedId(id || authUser?.userId || null);
    }
    return () => { active = false; };
  }, [username, id, authUser?.userId]);

  useEffect(() => {
    if (userId) {
      fetchUser();
      fetchUserPosts();
    }
  }, [userId]);

  const fetchUser = async () => {
    try {
      setLoading(true);
      const { data } = await api.get(`/api/users/${userId}`);
      setUser(data);
      setForm(data);

      if (authUser?.userId !== userId) {
        const followers = await api.get(`/api/users/${userId}/followers`)
        const isAlreadyFollowing = (followers.data?.data || followers.data || [])
          .some(f => f.id === authUser?.userId)
        setIsFollowing(isAlreadyFollowing)

        const friends = await api.get(`/api/users/${userId}/friends`)
        const isAlreadyFriend = (friends.data?.data || friends.data || [])
          .some(f => f.id === authUser?.userId)
        setIsFriend(isAlreadyFriend)
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const fetchUserPosts = async () => {
    try {
      setLoadingPosts(true)
      const data = await getPostsByAuthor(userId)
      setPosts(data)
    } catch (err) {
      console.error(err)
    } finally {
      setLoadingPosts(false)
    }
  }

  const handleUpdate = async () => {
    try {
      const { data: updatedUser } = await api.put(`/api/users/${userId}`, form);
      setUser({ ...user, ...updatedUser });
	  if (authUser?.userId === userId) {
		updateUser({ avatar: updatedUser.avatar, displayname: updatedUser.displayname })
	  }
      setShowEdit(false);
    } catch (err) {
      console.error(err);
    }
  };

  const handleFriendRequest = async () => {
    try {
      await sendFriendRequest(userId)
      setFriendRequested(true)
    } catch (err) {
      console.error('Error sending friend request:', err)
    }
  }

  if (loading) return <p>Loading profile...</p>;

  return (
    <div className="min-h-screen bg-transparent">
      {/* Banner */}
      <div className="h-56 w-full overflow-hidden" style={{ background: 'linear-gradient(120deg, #fde8f0 0%, #ede8fd 50%, #e8f0fd 100%)' }}>
        {user.wallpaper && (
          <img src={user.wallpaper} alt="banner" className="w-full h-full object-cover" />
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
                  ? <img src={user.avatar || '/default-avatar.jpg'} alt="avatar" className="w-full h-full object-cover"
                    onError={(e) => { e.target.src = '/default-avatar.jpg' }} />
                  : <div className="w-full h-full bg-gray-300" />
                }
              </div>
            </div>
          </div>

          {/* Botones top-right */}
          <div className="flex justify-end pt-3 pb-2 gap-2 flex-wrap">
            {authUser?.userId === userId ? (
              <>
                <button
                  onClick={() => setShowEdit(true)}
                  className="text-sm font-semibold px-4 py-1.5 rounded-full transition-colors hover:bg-[#ede8fd]"
                  style={{ border: '1px solid #ede8fd', color: '#534ab7', background: 'white' }}
                >
                  Edit profile
                </button>
                <button
                  onClick={() => navigate('/settings')}
                  className="text-sm font-semibold px-4 py-1.5 rounded-full transition-colors hover:bg-[#ede8fd]"
                  style={{ border: '1px solid #ede8fd', color: '#534ab7', background: 'white' }}
                >
                  Settings
                </button>
              </>
            ) : (
              <div className="flex gap-2">
                <FollowButton onFollow={fetchUser} targetId={userId} isFollowing={isFollowing} />
                <button
                  onClick={handleFriendRequest}
                  disabled={friendRequested || isFriend}
                  className="text-sm font-semibold px-4 py-1.5 rounded-full transition-colors hover:bg-[#faf8f5]"
                  style={{ border: '1px solid #ede8fd', color: '#534ab7', background: 'white' }}
                >
                  {isFriend ? 'Friends ✓' : friendRequested ? 'Requested' : 'Add friend'}
                </button>
                <button
                  onClick={() => navigate(`/messages/${id}`)}
                  className="text-sm font-semibold px-4 py-1.5 rounded-full transition-colors hover:bg-[#faf8f5]"
                  style={{ border: '1px solid #ede8fd', color: '#534ab7', background: 'white' }}
                >
                  Message
                </button>
              </div>
            )}
          </div>

          {/* Info */}
          <div className="mt-12 pb-3">
            <h2 className="text-xl font-bold" style={{ color: '#2c2c2a' }}>{user.displayname}</h2>
            <p className="text-sm" style={{ color: '#afa9ec' }}>@{user.username}</p>
            {user.bio && (
              <p className="mt-2 text-sm" style={{ color: '#5f5e5a', borderLeft: '2px solid #a78bfa', paddingLeft: '8px' }}>
                {user.bio}
              </p>
            )}
            <div className="flex gap-4 mt-3">
              <span className="text-sm cursor-pointer hover:underline">
                <strong style={{ color: '#534ab7' }}>{user.following_count ?? 0} </strong>
                <span onClick={() => setFollowModal('following')} className="ml-1" style={{ color: '#b4b2a9' }}>Following</span>
              </span>
              <span className="text-sm cursor-pointer hover:underline">
                <strong className="text-black">{user.followers_count ?? 0}</strong>
                <span  onClick={() => setFollowModal('followers')} className="text-gray-500 ml-1">Followers</span>
              </span>
            </div>
          </div>
        </div>

        {/* Posts tab */}
        <div className="flex border-b" style={{ borderColor: '#ede8fd' }}>
          {['posts', 'replies', 'media'].map(tab => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className="flex-1 text-center py-4 text-sm font-semibold capitalize transition-colors"
              style={{
                color: activeTab === tab ? '#534ab7' : '#b4b2a9',
                borderBottom: activeTab === tab ? '2px solid #534ab7' : '2px solid transparent',
                background: 'transparent'
              }}
            >
              {tab.charAt(0).toUpperCase() + tab.slice(1)}
            </button>
          ))}
        </div>
      </div>

      {/* Posts list */}
      <div className="mt-6">
        {activeTab === 'posts' && (
          loadingPosts ? (
            <p className="text-center py-8" style={{ color: '#b4b2a9' }}>Loading posts...</p>
          ) : posts.length === 0 ? (
            <p className="text-center py-8" style={{ color: '#b4b2a9' }}>No posts yet</p>
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
          )
        )}
        {activeTab === 'replies' && (
          <p className="text-center py-8" style={{ color: '#b4b2a9' }}>No replies yet</p>
        )}
        {activeTab === 'media' && (
          <p className="text-center py-8" style={{ color: '#b4b2a9' }}>No media yet</p>
        )}
      </div>

      {/* Modal Edit */}
      {showEdit && (
        <div className="fixed inset-0 z-50 flex items-center justify-center" style={{ background: 'rgba(237, 232, 253, 0.4)', backdropFilter: 'blur(4px)' }}>
          <div className="absolute inset-0" onClick={() => setShowEdit(false)} />
          <div className="relative z-10 bg-white rounded-2xl w-full max-w-lg overflow-hidden" style={{ border: '2px solid #ede8fd', boxShadow: '0 8px 32px rgba(167, 139, 250, 0.15)' }}>
            <div className="flex items-center gap-4 px-4 py-3 border-b" style={{ borderColor: '#ede8fd' }}>
              <button onClick={() => setShowEdit(false)} className="rounded-full p-2 text-sm transition-colors hover:bg-[#faf8f5]" style={{ color: '#534ab7' }}>✕</button>
              <h3 className="text-lg font-bold flex-1" style={{ color: '#2c2c2a' }}>Edit profile</h3>
              <button onClick={handleUpdate} className="text-sm font-bold px-4 py-1.5 rounded-full transition-colors" style={{ background: '#534ab7', color: 'white', border: 'none' }}>
                Save
              </button>
            </div>

            <div className="relative h-32" style={{ background: 'linear-gradient(120deg, #fde8f0 0%, #ede8fd 50%, #e8f0fd 100%)' }}>
              {form.wallpaper && <img src={form.wallpaper} alt="banner" className="w-full h-full object-cover" />}
              <label className="absolute inset-0 flex items-center justify-center bg-black/10 cursor-pointer hover:bg-opacity-40 transition-all">
                {uploadingBanner
                  ? <span className="text-white text-xs">Uploading...</span>
                  : <svg viewBox="0 0 24 24" className="w-6 h-6 fill-white">
                    <path d="M12 15.2a3.2 3.2 0 100-6.4 3.2 3.2 0 000 6.4z" />
                    <path d="M9 2L7.17 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2h-3.17L15 2H9zm3 15c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5z" />
                  </svg>
                }
                <input type="file" accept="image/*" className="hidden"
                  onChange={async (e) => {
                    const file = e.target.files[0]
                    if (!file) return
                    setUploadingBanner(true)
                    try {
                      const formData = new FormData()
                      formData.append('file', file)
                      const { data } = await api.post('/api/upload', formData)
                      setForm(prev => ({ ...prev, wallpaper: data.url }))
                    } catch (err) {
                      console.error(err)
                    } finally {
                      setUploadingBanner(false)
                    }
                  }}
                />
              </label>
            </div>

            <div className="px-4">
              <div className="relative -mt-10 mb-4 w-20 h-20">
                <div className="w-20 h-20 rounded-full border-4 border-white bg-gray-300 overflow-hidden">
                  <img src={form.avatar || '/default-avatar.jpg'} alt="avatar" className="w-full h-full object-cover"
                    onError={(e) => { e.target.src = '/default-avatar.jpg' }} />
                </div>
                <label className="absolute inset-0 flex items-center justify-center bg-black/10 rounded-full cursor-pointer hover:bg-opacity-50 transition-all">
                  {uploadingAvatar
                    ? <span className="text-white text-xs">...</span>
                    : <svg viewBox="0 0 24 24" className="w-5 h-5 fill-white">
                      <path d="M12 15.2a3.2 3.2 0 100-6.4 3.2 3.2 0 000 6.4z" />
                      <path d="M9 2L7.17 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2h-3.17L15 2H9zm3 15c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5z" />
                    </svg>
                  }
                  <input type="file" accept="image/*" className="hidden"
                    onChange={async (e) => {
                      const file = e.target.files[0]
                      if (!file) return
                      setUploadingAvatar(true)
                      try {
                        const formData = new FormData()
                        formData.append('file', file)
                        const { data } = await api.post('/api/upload', formData)
                        setForm(prev => ({ ...prev, avatar: data.url }))
                      } catch (err) {
                        console.error(err)
                      } finally {
                        setUploadingAvatar(false)
                      }
                    }}
                  />
                </label>
              </div>

              <input className="w-full border border-gray-300 rounded-lg px-3 py-2 mb-3 text-sm focus:outline-none focus:border-blue-400"
                placeholder="Display name" type="text" value={form.displayname}
                onChange={(e) => setForm({ ...form, displayname: e.target.value })} />
              <input className="w-full border border-gray-300 rounded-lg px-3 py-2 mb-3 text-sm focus:outline-none focus:border-blue-400"
                placeholder="Email" type="email" value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })} />
              <textarea className="w-full border border-gray-300 rounded-lg px-3 py-2 mb-4 text-sm focus:outline-none focus:border-blue-400 resize-none"
                placeholder="Bio" rows={3} value={form.bio}
                onChange={(e) => setForm({ ...form, bio: e.target.value })} />
            </div>
          </div>
        </div>
      )}
	{followModal && (
	<FollowersModal
		userId={userId}
		type={followModal}
		onClose={() => setFollowModal(null)}
	/>
	)}
    </div>
  )
}
