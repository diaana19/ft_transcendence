/*
** File: FileUploader.jsx
** Description: Component for uploading files to the server
** Responsibilities:
** - Render file input and upload button
** - Validate file size before sending (max 5MB)
** - Send file to API using multipart/form-data
** - Show upload progress and result
*/

import { useState } from 'react'
import axios from 'axios'

function FileUploader({ onUploadSuccess }) {
  const [file, setFile] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [success, setSuccess] = useState(null)

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0]
    if (selectedFile && selectedFile.size > 5 * 1024 * 1024) {
      setError('File too large. Maximum size is 5MB')
      setFile(null)
      return
    }
    setFile(selectedFile)
    setError(null)
  }

  const handleUpload = async () => {
    if (!file) {
      setError('Please select a file')
      return
    }
    setLoading(true)
    setError(null)
    setSuccess(null)
    try {
      const formData = new FormData()
      formData.append('file', file)
      const response = await axios.post('/api/upload', formData, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      })
      setSuccess(response.data.message)
      onUploadSuccess(response.data.path)
      setFile(null)
    } catch (err) {
      setError(err.response?.data?.error || 'Something went wrong')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex flex-col items-center gap-4">

      {error && (
        <p className="text-red-500 text-sm bg-red-50 p-3 rounded-lg w-full">
          {error}
        </p>
      )}

      {success && (
        <p className="text-green-500 text-sm bg-green-50 p-3 rounded-lg w-full">
          {success}
        </p>
      )}

      <label className="cursor-pointer w-full">
        <div className="border-2 border-dashed border-gray-300 rounded-2xl p-6 flex flex-col items-center gap-2 hover:border-blue-400 transition-colors">
          {file ? (
            <img
              src={URL.createObjectURL(file)}
              alt="preview"
              className="w-32 h-32 rounded-full object-cover"
            />
          ) : (
            <div className="w-32 h-32 rounded-full bg-gray-100 flex items-center justify-center">
              <span className="text-gray-400 text-sm">Click to upload</span>
            </div>
          )}
          {file && (
            <p className="text-gray-500 text-sm">
              {file.name} — {(file.size / 1024 / 1024).toFixed(2)} MB
            </p>
          )}
        </div>
        <input
          type="file"
          accept="image/*"
          onChange={handleFileChange}
          className="hidden"
        />
      </label>

      <button
        onClick={handleUpload}
        disabled={loading || !file}
        className="px-4 py-2 bg-blue-400 hover:bg-blue-500 text-white font-bold rounded-full disabled:opacity-50 w-full"
      >
        {loading ? 'Uploading...' : 'Save photo'}
      </button>

    </div>
  )
}

export default FileUploader