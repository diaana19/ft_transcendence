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
import { useParams } from "react-router-dom";

export default function Profile() {
	const { id } = useParams();
	const userId = id;

	const [user, setUser] = useState({
	name: "",
	email: "",
	bio: "",
	});

	const [loading, setLoading] = useState(false);
	const [showEdit, setShowEdit] = useState(false);
	const [showDelete, setShowDelete] = useState(false);
	const [form, setForm] = useState(user);

	useEffect(() => {
	if (userId) {
	fetchUser();
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

	const res = await fetch(`/api/users/${userId}`);
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
		},
		body: JSON.stringify(form),
		});

		const updatedUser = await res.json();

		setUser(updatedUser);
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
			});

			alert("User Deleted");
		} catch (err) {
			console.error(err);
		}
	};

	if (loading) return <p>Cargando perfil...</p>;

	return (
	<div className="min-h-screen bg-gray-100">
	{/* Banner */}
	<div className="h-40 bg-blue-500"></div>

	<div className="max-w-2xl mx-auto bg-white shadow">

	<div className="relative px-4">
		<div className="absolute -top-16">
		<div className="w-32 h-32 rounded-full border-4 border-white bg-gray-300"></div>
		</div>

		<div className="flex justify-end pt-4">
		<button
			onClick={() => setShowEdit(true)}
			className="border px-4 py-2 rounded-full font-semibold hover:bg-gray-100"
		>
			Edit profile
		</button>

		<button
			onClick={() => setShowDelete(true)}
			className="ml-2 border px-4 py-2 rounded-full text-red-500 hover:bg-red-50"
		>
			Delete
		</button>
		</div>

		{/* Info */}
		<div className="mt-20 pb-4">
		<h2 className="text-xl font-bold">{user.name}</h2>
		<p className="text-gray-500">@{user.email}</p>

		<p className="mt-3 text-gray-800">{user.bio}</p>
		</div>
	</div>
	</div>

	{showEdit && (
	<div className="fixed inset-0 bg-black bg-opacity-40 flex justify-center items-center">
		<div className="bg-white p-6 rounded-xl w-96 shadow-lg">
		<h3 className="text-lg font-bold mb-4">Edit profile</h3>

		<input
			className="w-full border p-2 rounded mb-2"
			type="text"
			value={form.name}
			onChange={(e) =>
			setForm({ ...form, name: e.target.value })
			}
		/>

		<input
			className="w-full border p-2 rounded mb-2"
			type="email"
			value={form.email}
			onChange={(e) =>
			setForm({ ...form, email: e.target.value })
			}
		/>

		<textarea
			className="w-full border p-2 rounded mb-4"
			value={form.bio}
			onChange={(e) =>
			setForm({ ...form, bio: e.target.value })
			}
		/>

		<div className="flex justify-end gap-2">
			<button
			onClick={() => setShowEdit(false)}
			className="px-4 py-2 rounded bg-gray-200"
			>
			Cancel
			</button>

			<button
			onClick={handleUpdate}
			className="px-4 py-2 rounded bg-blue-500 text-white"
			>
			Save
			</button>
		</div>
		</div>
	</div>
	)}

	{showDelete && (
	<div className="fixed inset-0 bg-black bg-opacity-40 flex justify-center items-center">
		<div className="bg-white p-6 rounded-xl w-80 shadow-lg">
		<h3 className="text-lg font-bold mb-2">Delete account</h3>
		<p className="text-gray-600 mb-4">
			This action is permanent
		</p>

		<div className="flex justify-end gap-2">
			<button
			onClick={() => setShowDelete(false)}
			className="px-4 py-2 rounded bg-gray-200"
			>
			Cancel
			</button>

			<button
			onClick={handleDelete}
			className="px-4 py-2 rounded bg-red-500 text-white"
			>
			Yes, delete
			</button>
		</div>
		</div>
	</div>
	)}
	</div>
	);
	};