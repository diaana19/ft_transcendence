import { useEffect, useState } from "react";

export default function Profile() {
  const userId = "1"; // 🔥 usuario fijo temporal

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

  // ✏️ UPDATE USER
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

  // ❌ DELETE USER
  const handleDelete = async () => {
    try {
      await fetch(`/api/users/${userId}`, {
        method: "DELETE",
      });

      alert("Usuario eliminado");
    } catch (err) {
      console.error(err);
    }
  };

  if (loading) return <p>Cargando perfil...</p>;

  return (
    <div style={{ maxWidth: "500px", margin: "auto" }}>
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
              value={form.name}
              onChange={(e) =>
                setForm({ ...form, name: e.target.value })
              }
            />

            <input
              type="email"
              value={form.email}
              onChange={(e) =>
                setForm({ ...form, email: e.target.value })
              }
            />

            <textarea
              value={form.bio}
              onChange={(e) =>
                setForm({ ...form, bio: e.target.value })
              }
            />

            <button onClick={handleUpdate}>Guardar</button>
            <button onClick={() => setShowEdit(false)}>Cancelar</button>
          </div>
        </div>
      )}

      {/* ⚠️ MODAL DELETE */}
      {showDelete && (
        <div style={overlayStyle}>
          <div style={modalStyle}>
            <h3>Eliminar cuenta</h3>
            <p>Esta acción es permanente</p>

            <button
              onClick={handleDelete}
              style={{ backgroundColor: "red", color: "white" }}
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

// 🎨 estilos
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