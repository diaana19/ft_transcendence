import { colors } from '../../styles/tokens'

// Input is the shared text input with optional label and error message.
function Input({
    label,
    type = 'text',
    value,
    onChange,
    placeholder = '',
    error = null,
    disabled = false,
    fullWidth = true,
    name = '',
}) {
    return (
        <div className={`flex flex-col gap-1 ${fullWidth ? 'w-full' : ''}`}>
            {label && <label className="text-sm font-medium text-gray-700">{label}</label>}
            <input
                type={type}
                name={name}
                value={value}
                onChange={onChange}
                placeholder={placeholder}
                disabled={disabled}
                className={`border rounded-lg px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none transition-colors
          ${error ? 'border-red-400 focus:border-red-500' : 'border-gray-300 focus:border-blue-400'}
          ${disabled ? 'bg-gray-50 cursor-not-allowed' : 'bg-white'}
          ${fullWidth ? 'w-full' : ''}
        `}
            />
            {error && <p className="text-xs text-red-500">{error}</p>}
        </div>
    )
}

export default Input
