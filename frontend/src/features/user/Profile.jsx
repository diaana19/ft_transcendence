import { useEffect, useState } from "react";
import { getUser, updateUser, deleteUser } from "./backend/services/user_service";
import { useAuth } from "../../backend/useAuth";
import { useEffect, useState } from "react";
import { useAuth } from "./useAuth";

export default function Profile() {
  const { token, logout } = useAuth();

  const [user, setUser] = useState({
    name: "",
    email: "",
    bio: "",
  });

  const [loading, setLoading] = useState(false);
  const [showEdit, setShowEdit] = useState(false);
  const [showDelete, setShowDelete] = useState(false);

  const [form, setForm] = useState(user);

  // 🔄 GET USER
  useEffect(() => {
    fetchUser();
  }, []);

  const fetchUser = async () => {
    try {
      setLoading(true);
      const res = await fetch("/api/user/me", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
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

  // ✏️ UPDATE USER
  const handleUpdate = async () => {
    try {
      await fetch("/api/user/update", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(form),
      });

      setUser(form);
      setShowEdit(false);
    } catch (err) {
      console.error(err);
    }
  };

  // ❌ DELETE USER
  const handleDelete = async () => {
    try {
      await fetch("/api/user/delete", {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      logout(); // cerrar sesión después de eliminar
    } catch (err) {
      console.error(err);
    }
  };

  if (loading) return <p>Cargando perfil...</p>;

  return (
    <div style={{ maxWidth: "500px", margin: "auto" }}>
      {/* 🧑 PERFIL */}
      <h2>Perfil</h2>

      <div style={{ border: "1px solid #ccc", padding: "16px", borderRadius: "10px" }}>
        <p><strong>Nombre:</strong> {user.name}</p>
        <p><strong>Email:</strong> {user.email}</p>
        <p><strong>Bio:</strong> {user.bio}</p>

        <button onClick={() => setShowEdit(true)}>
          Editar perfil
        </button>

        <button
          onClick={() => setShowDelete(true)}
          style={{ marginLeft: "10px", color: "red" }}
        >
          Eliminar cuenta
        </button>
      </div>

      {/* ✏️ MODAL EDIT */}
      {showEdit && (
        <div style={overlayStyle}>
          <div style={modalStyle}>
            <h3>Editar perfil</h3>

            <input
              type="text"
              placeholder="Nombre"
              value={form.name}
              onChange={(e) =>
                setForm({ ...form, name: e.target.value })
              }
            />

            <input
              type="email"
              placeholder="Email"
              value={form.email}
              onChange={(e) =>
                setForm({ ...form, email: e.target.value })
              }
            />

            <textarea
              placeholder="Bio"
              value={form.bio}
              onChange={(e) =>
                setForm({ ...form, bio: e.target.value })
              }
            />

            <div style={{ marginTop: "10px" }}>
              <button onClick={handleUpdate}>
                Guardar
              </button>

              <button onClick={() => setShowEdit(false)}>
                Cancelar
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ⚠️ MODAL DELETE */}
      {showDelete && (
        <div style={overlayStyle}>
          <div style={modalStyle}>
            <h3>Eliminar cuenta</h3>
            <p>Esta acción es permanente. No se puede deshacer.</p>

            <button
              onClick={handleDelete}
              style={{ color: "white", backgroundColor: "red" }}
            >
              Sí, eliminar
            </button>

            <button onClick={() => setShowDelete(false)}>
              Cancelar
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

// 🎨 estilos simples tipo modal
const overlayStyle = {
  position: "fixed",
  top: 0,
  left: 0,
  width: "100%",
  height: "100%",
  backgroundColor: "rgba(0,0,0,0.5)",
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
};

const modalStyle = {
  background: "white",
  padding: "20px",
  borderRadius: "10px",
  width: "300px",
  display: "flex",
  flexDirection: "column",
  gap: "10px",
};